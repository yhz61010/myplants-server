package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents an application user
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"unique;not null"`
	Password  string         `json:"-" gorm:"not null"`
	Avatar    string         `json:"avatar"`
	Bio       string         `json:"bio" gorm:"type:text"`
	IsAdmin   bool           `json:"isAdmin" gorm:"default:false"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
