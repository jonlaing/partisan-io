package models

import (
	"time"
)

// ImageAttachment is used to attach images to a Post (or something else, when we have it)
type ImageAttachment struct {
	ID         uint64    `json:"id" gorm:"primary_key"`      // Primary key
	UserID     uint64    `json:"user_id" binding:"required"` // ID of user that created Post
	RecordID   uint64    `json:"record_id" binding:"required"`
	RecordType string    `json:"record_type" binding:"required"`
	URL   string    `json:"image_url"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
