package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"myplants-server/internal/database"
	"myplants-server/internal/models"
)

// CreateDiaryRequest represents the request payload for creating a diary
type CreateDiaryRequest struct {
	UserID     string   `json:"userId" binding:"required"`
	Title      string   `json:"title" binding:"required"`
	Content    string   `json:"content"`
	Images     []string `json:"images"`
	Tags       []string `json:"tags"`
	CreateTime string   `json:"createTime"` // Assuming it's a string, e.g., "2023-01-01T00:00:00Z"
}

// CreateDiary handles POST /api/diaries
func CreateDiary(c *gin.Context) {
	var req CreateDiaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse createTime if provided, otherwise use current time
	var createdAt time.Time
	if req.CreateTime != "" {
		var err error
		createdAt, err = time.Parse(time.RFC3339, req.CreateTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid createTime format"})
			return
		}
	} else {
		createdAt = time.Now()
	}

	// Convert slices to JSON strings
	imagesJSON, _ := json.Marshal(req.Images)
	tagsJSON, _ := json.Marshal(req.Tags)

	// Create diary entry
	diary := models.Diary{
		UserID:    req.UserID,
		Title:     req.Title,
		Content:   req.Content,
		ImagesStr: string(imagesJSON),
		TagsStr:   string(tagsJSON),
		CreatedAt: createdAt,
	}

	// Save to database
	if err := database.GetDB().Create(&diary).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create diary"})
		return
	}

	// Load the diary to populate computed fields
	database.GetDB().First(&diary, diary.ID)

	// Manually populate computed fields
	if err := json.Unmarshal([]byte(diary.ImagesStr), &diary.Images); err != nil {
		diary.Images = []string{}
	}
	if err := json.Unmarshal([]byte(diary.TagsStr), &diary.Tags); err != nil {
		diary.Tags = []string{}
	}

	c.JSON(http.StatusCreated, diary)
}
