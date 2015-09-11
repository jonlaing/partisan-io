package models

import (
  "time"
)

// Post is the primary user created content. It can be just about anything.
type Post struct {
	ID        uint64    `json:"id" gorm:"primary_key"`      // Primary key
	UserID    uint64    `json:"user_id" binding:"required"` // ID of user that created Post
	Body      string    `json:"body" binding:"required"`    // The text based body
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
