package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"myplants-server/internal/database"
	"myplants-server/internal/models"
)

// CreateContentRequest payload
type CreateContentRequest struct {
	Type       string   `json:"type" binding:"required"`
	UserID     string   `json:"userId" binding:"required"`
	Title      string   `json:"title" binding:"required"`
	Content    string   `json:"content"`
	Images     []string `json:"images"`
	Tags       []string `json:"tags"`
	CreateTime string   `json:"createTime"`
}

// UpdateContentRequest payload
type UpdateContentRequest struct {
	Title   *string   `json:"title"`
	Content *string   `json:"content"`
	Images  *[]string `json:"images"`
	Tags    *[]string `json:"tags"`
}

// CreateContent handles POST /api/contents
func CreateContent(c *gin.Context) {
	var req CreateContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var createdAt time.Time
	if req.CreateTime != "" {
		if t, err := time.Parse(time.RFC3339, req.CreateTime); err == nil {
			createdAt = t
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid createTime format"})
			return
		}
	} else {
		createdAt = time.Now()
	}

	content := models.Content{
		Type:      req.Type,
		UserID:    req.UserID,
		Title:     req.Title,
		Content:   req.Content,
		Images:    req.Images,
		Tags:      req.Tags,
		CreatedAt: createdAt,
	}

	if err := database.GetDB().Create(&content).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create content"})
		return
	}

	c.JSON(http.StatusCreated, content)
}

// CreateDiary handles POST /api/diaries and uses authenticated user as owner
func CreateDiary(c *gin.Context) {
	var req struct {
		Title      string   `json:"title" binding:"required"`
		Content    string   `json:"content"`
		Images     []string `json:"images"`
		Tags       []string `json:"tags"`
		CreateTime string   `json:"createTime"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// get authenticated user id from context
	uid, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userIdStr := strconv.FormatUint(uint64(uid.(uint)), 10)

	var createdAt time.Time
	if req.CreateTime != "" {
		if t, err := time.Parse(time.RFC3339, req.CreateTime); err == nil {
			createdAt = t
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid createTime format"})
			return
		}
	} else {
		createdAt = time.Now()
	}

	content := models.Content{
		Type:      "diary",
		UserID:    userIdStr,
		Title:     req.Title,
		Content:   req.Content,
		Images:    req.Images,
		Tags:      req.Tags,
		CreatedAt: createdAt,
	}

	if err := database.GetDB().Create(&content).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create diary"})
		return
	}
	c.JSON(http.StatusCreated, content)
}

// ListDiaries handles GET /api/diaries
func ListDiaries(c *gin.Context) {
	// reuse ListContents logic by forcing type=diary
	q := c.Query("query")
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	db := database.GetDB()
	var items []models.Content
	query := db.Model(&models.Content{}).Where("type = ?", "diary")
	if q != "" {
		like := "%" + q + "%"
		query = query.Where("title LIKE ? OR tags LIKE ?", like, like)
	}
	var total int64
	query.Count(&total)

	if err := query.Order("created_at desc").Limit(limit).Offset(offset).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "total": total})
}

// GetDiary handles GET /api/diaries/:id
func GetDiary(c *gin.Context) {
	id := c.Param("id")
	var content models.Content
	if err := database.GetDB().First(&content, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	if content.Type != "diary" {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, content)
}

// UpdateDiary handles PUT /api/diaries/:id (only owner can update)
func UpdateDiary(c *gin.Context) {
	id := c.Param("id")
	var req UpdateContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var content models.Content
	db := database.GetDB()
	if err := db.First(&content, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	if content.Type != "diary" {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	uid, exists := c.Get("userId")
	if !exists || strconv.FormatUint(uint64(uid.(uint)), 10) != content.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if req.Title != nil {
		content.Title = *req.Title
	}
	if req.Content != nil {
		content.Content = *req.Content
	}
	if req.Images != nil {
		content.Images = *req.Images
	}
	if req.Tags != nil {
		content.Tags = *req.Tags
	}

	if err := db.Save(&content).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update"})
		return
	}
	c.JSON(http.StatusOK, content)
}

// DeleteDiary handles DELETE /api/diaries/:id (only owner)
func DeleteDiary(c *gin.Context) {
	id := c.Param("id")
	var content models.Content
	db := database.GetDB()
	if err := db.First(&content, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	if content.Type != "diary" {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	uid, exists := c.Get("userId")
	if !exists || strconv.FormatUint(uint64(uid.(uint)), 10) != content.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if err := db.Delete(&models.Content{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete"})
		return
	}
	c.Status(http.StatusNoContent)
}

// GetContent handles GET /api/contents/:id
func GetContent(c *gin.Context) {
	id := c.Param("id")
	var content models.Content
	if err := database.GetDB().First(&content, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, content)
}

// UpdateContent handles PUT /api/contents/:id
func UpdateContent(c *gin.Context) {
	id := c.Param("id")
	var req UpdateContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var content models.Content
	db := database.GetDB()
	if err := db.First(&content, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	if req.Title != nil {
		content.Title = *req.Title
	}
	if req.Content != nil {
		content.Content = *req.Content
	}
	if req.Images != nil {
		content.Images = *req.Images
	}
	if req.Tags != nil {
		content.Tags = *req.Tags
	}

	if err := db.Save(&content).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update"})
		return
	}

	c.JSON(http.StatusOK, content)
}

// DeleteContent handles DELETE /api/contents/:id
func DeleteContent(c *gin.Context) {
	id := c.Param("id")
	if err := database.GetDB().Delete(&models.Content{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete"})
		return
	}
	c.Status(http.StatusNoContent)
}

// ListContents handles GET /api/contents
// supports query, limit, offset, and type filter
func ListContents(c *gin.Context) {
	q := c.Query("query")
	typeFilter := c.Query("type")
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	db := database.GetDB()
	var items []models.Content
	query := db.Model(&models.Content{})
	if typeFilter != "" {
		query = query.Where("type = ?", typeFilter)
	}
	if q != "" {
		like := "%" + q + "%"
		// match title OR tags (tags stored as JSON string)
		query = query.Where("title LIKE ? OR tags LIKE ?", like, like)
	}
	var total int64
	query.Count(&total)

	if err := query.Order("created_at desc").Limit(limit).Offset(offset).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items, "total": total})
}
