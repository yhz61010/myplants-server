package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"myplants-server/internal/auth"
	"myplants-server/internal/database"
	"myplants-server/internal/models"
)

// RegisterRequest represents the request payload for user registration
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
}

// LoginRequest represents the request payload for user login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register handles POST /api/auth/register
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	var existingUser models.User
	if err := database.GetDB().Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create user
	user := models.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Avatar:   req.Avatar,
		Bio:      req.Bio,
	}

	if err := database.GetDB().Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Return user without password
	user.Password = ""
	c.JSON(http.StatusCreated, user)
}

// Login handles POST /api/auth/login
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user
	var user models.User
	if err := database.GetDB().Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token via internal/auth
	tokenString, err := auth.SignToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"avatar":   user.Avatar,
			"bio":      user.Bio,
			"isAdmin":  user.IsAdmin,
		},
	})
}

// ListUsers handles GET /api/users with pagination
func ListUsers(c *gin.Context) {
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
	var users []models.User
	var total int64
	db.Model(&models.User{}).Count(&total)
	if err := db.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	// strip passwords
	for i := range users {
		users[i].Password = ""
	}
	c.JSON(http.StatusOK, gin.H{"items": users, "total": total})
}

// GetUser handles GET /api/users/:id
func GetUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := database.GetDB().First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	user.Password = ""
	c.JSON(http.StatusOK, user)
}

// UpdateUser handles PUT /api/users/:id
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	// only allow owner to update
	uid, _ := c.Get("userId")
	if uid == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if fmt.Sprintf("%v", uid) != id && fmt.Sprintf("%v", uint(uid.(uint))) != id {
		// allow numeric mismatch handling
		if uid.(uint) != 0 {
			if strconv.FormatUint(uint64(uid.(uint)), 10) != id {
				c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
				return
			}
		}
	}

	var req struct {
		Avatar   *string `json:"avatar"`
		Bio      *string `json:"bio"`
		Password *string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	db := database.GetDB()
	if err := db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	if req.Avatar != nil {
		user.Avatar = *req.Avatar
	}
	if req.Bio != nil {
		user.Bio = *req.Bio
	}
	if req.Password != nil {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}
		user.Password = string(hashed)
	}

	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update"})
		return
	}
	user.Password = ""
	c.JSON(http.StatusOK, user)
}

// DeleteUser handles DELETE /api/users/:id
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	uid, _ := c.Get("userId")
	if uid == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if strconv.FormatUint(uint64(uid.(uint)), 10) != id {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	if err := database.GetDB().Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete"})
		return
	}
	c.Status(http.StatusNoContent)
}

// --- Admin helpers and handlers (merged from admin.go) ---

// AdminUpdateUserRequest for admin user updates
type AdminUpdateUserRequest struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
	Avatar   *string `json:"avatar"`
	Bio      *string `json:"bio"`
	IsAdmin  *bool   `json:"isAdmin"`
}

// AdminUpdateDiaryRequest for admin diary updates
type AdminUpdateDiaryRequest struct {
	Title    *string   `json:"title"`
	Content  *string   `json:"content"`
	Images   *[]string `json:"images"`
	Tags     *[]string `json:"tags"`
	IsPublic *bool     `json:"isPublic"`
}

// AdminCreateDiaryRequest for admin diary creation
type AdminCreateDiaryRequest struct {
	Title      string   `json:"title" binding:"required"`
	Content    string   `json:"content"`
	Images     []string `json:"images"`
	Tags       []string `json:"tags"`
	IsPublic   bool     `json:"isPublic"`
	UserID     string   `json:"userId"` // Optional: specify user, defaults to admin
	CreateTime string   `json:"createTime"`
}

// AdminMiddleware checks if user is admin
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, exists := c.Get("isAdmin")
		if !exists || !isAdmin.(bool) {
			c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// AdminListUsers handles GET /api/admin/users
func AdminListUsers(c *gin.Context) {
	var users []models.User
	if err := database.GetDB().Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	// strip passwords
	for i := range users {
		users[i].Password = ""
	}
	c.JSON(http.StatusOK, gin.H{"items": users, "total": len(users)})
}

// AdminGetUser handles GET /api/admin/users/:id
func AdminGetUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := database.GetDB().First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	user.Password = ""
	c.JSON(http.StatusOK, user)
}

// AdminUpdateUser handles PUT /api/admin/users/:id
func AdminUpdateUser(c *gin.Context) {
	id := c.Param("id")
	var req AdminUpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	db := database.GetDB()
	if err := db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	if req.Username != nil {
		user.Username = *req.Username
	}
	if req.Password != nil {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}
		user.Password = string(hashed)
	}
	if req.Avatar != nil {
		user.Avatar = *req.Avatar
	}
	if req.Bio != nil {
		user.Bio = *req.Bio
	}
	if req.IsAdmin != nil {
		user.IsAdmin = *req.IsAdmin
	}

	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update"})
		return
	}
	user.Password = ""
	c.JSON(http.StatusOK, user)
}

// AdminDeleteUser handles DELETE /api/admin/users/:id
func AdminDeleteUser(c *gin.Context) {
	id := c.Param("id")

	// Prevent deleting yourself
	uid, _ := c.Get("userId")
	if strconv.FormatUint(uint64(uid.(uint)), 10) == id {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot delete yourself"})
		return
	}

	if err := database.GetDB().Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete"})
		return
	}
	c.Status(http.StatusNoContent)
}

// AdminCreateDiary handles POST /api/admin/diaries
func AdminCreateDiary(c *gin.Context) {
	var req AdminCreateDiaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get admin user ID if no userId specified
	userID := req.UserID
	if userID == "" {
		uid, _ := c.Get("userId")
		userID = strconv.FormatUint(uint64(uid.(uint)), 10)
	}

	content := models.Content{
		Type:     "diary",
		UserID:   userID,
		Title:    req.Title,
		Content:  req.Content,
		Images:   req.Images,
		Tags:     req.Tags,
		IsPublic: req.IsPublic,
	}

	if err := database.GetDB().Create(&content).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create diary"})
		return
	}
	c.JSON(http.StatusCreated, content)
}

// AdminUpdateDiary handles PUT /api/admin/diaries/:id
func AdminUpdateDiary(c *gin.Context) {
	id := c.Param("id")
	var req AdminUpdateDiaryRequest
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
		c.JSON(http.StatusNotFound, gin.H{"error": "not a diary"})
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
	if req.IsPublic != nil {
		content.IsPublic = *req.IsPublic
	}

	if err := db.Save(&content).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update"})
		return
	}
	c.JSON(http.StatusOK, content)
}

// AdminDeleteDiary handles DELETE /api/admin/diaries/:id
func AdminDeleteDiary(c *gin.Context) {
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
		c.JSON(http.StatusNotFound, gin.H{"error": "not a diary"})
		return
	}

	if err := db.Delete(&models.Content{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete"})
		return
	}
	c.Status(http.StatusNoContent)
}
