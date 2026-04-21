package handlers

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"strconv"

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

	// If Upyun credentials present, attempt upload, otherwise save locally
	upyunBucket := os.Getenv("UPYUN_BUCKET")
	upyunOperator := os.Getenv("UPYUN_OPERATOR")
	upyunPassword := os.Getenv("UPYUN_PASSWORD")

	if upyunBucket != "" && upyunOperator != "" && upyunPassword != "" {
		// Use official UpYun SDK for upload
		up := upyun.NewUpYun(&upyun.UpYunConfig{
			Bucket:   upyunBucket,
			Operator: upyunOperator,
			Password: upyunPassword,
		})

		// Build PutObjectConfig with Reader. For non-*os.File readers set Content-Length header.
		headers := map[string]string{}
		if header.Size > 0 {
			headers["Content-Length"] = strconv.FormatInt(header.Size, 10)
		}
		if ct := header.Header.Get("Content-Type"); ct != "" {
			headers["Content-Type"] = ct
		}

		cfg := &upyun.PutObjectConfig{
			Path:    "/" + key,
			Reader:  file,
			Headers: headers,
		}

		if err := up.Put(cfg); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "upyun upload failed", "detail": err.Error()})
			return
		}

		publicURL := fmt.Sprintf("https://myplants.leovp.com/%s", key)
		c.JSON(http.StatusOK, gin.H{"url": publicURL})
		return
	}

	// Fallback: save locally to ./uploads and return domain-based URL
	uploadsDir := "uploads"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create upload dir"})
		return
	}
	dstPath := filepath.Join(uploadsDir, key)
	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create dir"})
		return
	}
	out, err := os.Create(dstPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}
	defer out.Close()
	if _, err := io.Copy(out, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to write file"})
		return
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
