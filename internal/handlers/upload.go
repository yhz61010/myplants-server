package handlers

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	upyun "github.com/upyun/go-sdk/v3/upyun"
)

// UploadImage handles POST /api/upload
// Accepts multipart/form-data with field "file". Returns JSON {url: string}
func UploadImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	defer file.Close()

	// sanitize filename and compute a short hash to avoid collisions
	fname := sanitizeFilename(header)
	now := time.Now().UTC().Format("20060102")
	key := fmt.Sprintf("%s/%s_%s", now, generateHash(fname), fname)

	// Write upload to a temporary file to avoid buffering large files in memory
	ext := filepath.Ext(header.Filename)
	tmp, err := os.CreateTemp("", "upload-*"+ext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create temp file"})
		return
	}
	tmpPath := tmp.Name()
	// ensure temp file closed
	if _, err := io.Copy(tmp, file); err != nil {
		tmp.Close()
		os.Remove(tmpPath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to write temp file"})
		return
	}
	tmp.Close()

	// compute MD5 of temp file for integrity check
	var md5Hex string
	var md5Base64 string
	if f, err := os.Open(tmpPath); err == nil {
		hasher := md5.New()
		if _, err := io.Copy(hasher, f); err == nil {
			sum := hasher.Sum(nil)
			md5Hex = hex.EncodeToString(sum)
			md5Base64 = base64.StdEncoding.EncodeToString(sum)
		}
		f.Close()
	}

	// If Upyun credentials present, attempt upload using local path (resumable)
	upyunBucket := os.Getenv("UPYUN_BUCKET")
	upyunOperator := os.Getenv("UPYUN_OPERATOR")
	upyunPassword := os.Getenv("UPYUN_PASSWORD")

	if upyunBucket != "" && upyunOperator != "" && upyunPassword != "" {
		up := upyun.NewUpYun(&upyun.UpYunConfig{
			Bucket:   upyunBucket,
			Operator: upyunOperator,
			Password: upyunPassword,
		})

		headers := map[string]string{}
		if ct := header.Header.Get("Content-Type"); ct != "" {
			headers["Content-Type"] = ct
		}
		if md5Base64 != "" {
			// HTTP Content-MD5 is base64 of MD5
			headers["Content-MD5"] = md5Base64
		}
		if md5Hex != "" {
			// also expose hex for debugging
			headers["X-Content-MD5-Hex"] = md5Hex
		}

		cfg := &upyun.PutObjectConfig{
			Path:      "/" + key,
			LocalPath: tmpPath,
			Headers:   headers,
		}

		if err := up.Put(cfg); err != nil {
			// cleanup temp
			os.Remove(tmpPath)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "upyun upload failed", "detail": err.Error()})
			return
		}

		// cleanup temp
		os.Remove(tmpPath)
		publicURL := fmt.Sprintf("https://myplants.leovp.com/%s", key)
		c.JSON(http.StatusOK, gin.H{"url": publicURL})
		return
	}

	// Fallback: move temp file into uploads dir and serve from domain
	uploadsDir := "uploads"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		os.Remove(tmpPath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create upload dir"})
		return
	}
	dstPath := filepath.Join(uploadsDir, key)
	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		os.Remove(tmpPath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create dir"})
		return
	}
	if err := os.Rename(tmpPath, dstPath); err != nil {
		// fallback to copy
		in, err := os.Open(tmpPath)
		if err != nil {
			os.Remove(tmpPath)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to move file"})
			return
		}
		out, err := os.Create(dstPath)
		if err != nil {
			in.Close()
			os.Remove(tmpPath)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
			return
		}
		if _, err := io.Copy(out, in); err != nil {
			in.Close()
			out.Close()
			os.Remove(tmpPath)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to write file"})
			return
		}
		in.Close()
		out.Close()
		os.Remove(tmpPath)
	}

	publicURL := fmt.Sprintf("https://myplants.leovp.com/%s", key)
	c.JSON(http.StatusOK, gin.H{"url": publicURL})
}

func sanitizeFilename(h *multipart.FileHeader) string {
	name := filepath.Base(h.Filename)
	// keep only letters, numbers, dash, underscore and dot
	name = strings.Map(func(r rune) rune {
		switch {
		case (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9'):
			return r
		case r == '.' || r == '-' || r == '_':
			return r
		default:
			return '_'
		}
	}, name)
	return name
}

func generateHash(s string) string {
	h := md5.Sum([]byte(fmt.Sprintf("%s:%d", s, time.Now().UnixNano())))
	return hex.EncodeToString(h[:])[:12]
}
