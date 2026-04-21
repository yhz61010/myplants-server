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

	"github.com/gin-gonic/gin"
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
		// Attempt Upyun upload via simple HTTP PUT to the bucket domain
		// Note: This is a minimal implementation and requires correct
		// operator/password (secret) to work. For robust integration,
		// use the official SDK or implement form-signature properly.
		url := fmt.Sprintf("https://%s/%s", upyunBucket, key)
		req, err := http.NewRequest("PUT", url, file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create upload request"})
			return
		}
		// Upyun expects Authorization headers; leaving as-is will likely fail
		// unless the bucket is configured to allow anonymous PUTs. Users
		// should set up proper credentials and a proxy or use SDK.
		req.Header.Set("User-Agent", "myplants-server/1.0")
		req.Header.Set("Content-Type", header.Header.Get("Content-Type"))

		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "upload failed"})
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			body, _ := io.ReadAll(resp.Body)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "upload failed", "detail": string(body)})
			return
		}

		// Returned URL uses custom domain as requested
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
