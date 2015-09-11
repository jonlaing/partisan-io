package models

import(
  "time"
)

// FeedItem is the record of a user interacting with a Post, this is used to build a feed
type FeedItem struct {
	ID         uint64      `json:"id" gorm:"primary_key"`      // Primary key
	UserID     uint64      `json:"user_id" binding:"required"` // ID of user that created Post
	Action     string      `json:"action" binding:"required"`
	RecordType string      `json:"record_type" binding:"required"`
	RecordID   uint64      `json:"record_id" binding:"required"`
	Record     interface{} `json:"record,omitempty" sql:"-"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}