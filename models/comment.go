package models

import (
	"partisan/Godeps/_workspace/src/github.com/jinzhu/gorm"
	"time"
)

// Comment is for blah
type Comment struct {
	ID        uint64    `gorm:"primary_key" json:"id"`
	PostID    uint64    `json:"post_id"`
	UserID    uint64    `json:"user_id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetID satisfies Notifier and Hashtagger interface
func (c *Comment) GetID() uint64 {
	return c.ID
}

// GetUserID satisfies Notifier and Hashtagger interface
func (c *Comment) GetUserID() uint64 {
	return c.UserID
}

// GetType satisfies the Notifier and Hashtagger interface
func (c *Comment) GetType() string {
	return "comment"
}

// GetContent satisfies the Hashtagger interface
func (c *Comment) GetContent() string {
	return c.Body
}

// GetRecordUserID returns the user ID of the record being commented upon. Satisfies Notifier interface.
func (c *Comment) GetRecordUserID(db *gorm.DB) (uint64, error) {
	var notifUserIDs []uint64
	err := db.Model(Post{}).Where("id = ?", c.PostID).Pluck("user_id", &notifUserIDs).Error

	if len(notifUserIDs) < 1 {
		return 0, err
	}

	return notifUserIDs[0], err
}
