package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

// Notifier is an interface for records that can produce notifications
type Notifier interface {
	GetID() uint64
	Type() string
}

// Notification is what it sounds like
type Notification struct {
	ID           uint64    `json:"id" gorm:"primary_key"`             // Primary key
	TargetUserID uint64    `json:"target_user_id" binding:"required"` // ID of user to recieve the notification
	RecordID     uint64    `json:"record_id"`
	RecordType   string    `json:"record_type"`
	Seen         bool      `json:"seen"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// NewNotification initializes and saves a new notification
func NewNotification(n Notifier, targetUserID uint64, db *gorm.DB) error {
	notif := Notification{
		TargetUserID: targetUserID,
		RecordID:     n.GetID(),
		RecordType:   n.Type(),
		Seen:         false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return db.Save(&notif).Error
}
