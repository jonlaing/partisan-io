package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

// Comment is for blah
type Comment struct {
	ID         uint64    `gorm:"primary_key" json:"id"`
	RecordType string    `json:"record_type"`
	RecordID   uint64    `json:"record_id"`
	UserID     uint64    `json:"user_id"`
	Body       string    `json:"body"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// GetID satisfies Notifier interface
func (c *Comment) GetID() uint64 {
	return c.ID
}

// Type satisfies the Notifier interrace
func (c *Comment) Type() string {
	return "comment"
}

// GetRecordUserID returns the user ID of the record being commented upon. Satisfies Notifier interface.
func (c *Comment) GetRecordUserID(db *gorm.DB) (uint64, error) {
	var notifUserIDs []uint64
	err := db.Table(c.RecordType).Where("id = ?", c.RecordID).Pluck("user_id", &notifUserIDs).Error

	if len(notifUserIDs) < 1 {
		return 0, err
	}

	return notifUserIDs[0], err
}
