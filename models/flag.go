package models

import (
	"time"
)

const (
	FlNoReason int = iota
	FlOffensive
	FlCopyright
	FlSpam
	FlOther
)

// Flag alerts admins of potentially problematic behavior of a user,
// or problematic content in a post/comment
type Flag struct {
	ID         uint64    `json:"id" gorm:"primary_key"` // Primary key
	UserID     uint64    `json:"user_id"`               // ID of user that created Post
	RecordID   uint64    `json:"record_id"`
	RecordType string    `json:"record_type"`
	Reason     int       `json:"reason"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
