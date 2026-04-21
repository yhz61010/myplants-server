package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// Content represents a unified plant diary or plant knowledge entry
type Content struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Type      string         `json:"type" gorm:"not null"` // diary, plant, etc.
	UserID    string         `json:"userId" gorm:"not null"`
	Title     string         `json:"title" gorm:"not null"`
	Content   string         `json:"content" gorm:"type:text"`
	ImagesStr string         `json:"-" gorm:"column:images"` // Store as JSON string
	TagsStr   string         `json:"-" gorm:"column:tags"`   // Store as JSON string
	Images    []string       `json:"images" gorm:"-"`        // Computed field
	Tags      []string       `json:"tags" gorm:"-"`          // Computed field
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// BeforeSave converts slices to JSON strings
func (d *Content) BeforeSave(tx *gorm.DB) error {
	if data, err := json.Marshal(d.Images); err == nil {
		d.ImagesStr = string(data)
	}
	if data, err := json.Marshal(d.Tags); err == nil {
		d.TagsStr = string(data)
	}
	return nil
}

// AfterFind converts JSON strings back to slices
func (d *Content) AfterFind(tx *gorm.DB) error {
	if err := json.Unmarshal([]byte(d.ImagesStr), &d.Images); err != nil {
		d.Images = []string{}
	}
	if err := json.Unmarshal([]byte(d.TagsStr), &d.Tags); err != nil {
		d.Tags = []string{}
	}
	return nil
}
