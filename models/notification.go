package models

import (
	"errors"
	"fmt"
	"partisan/Godeps/_workspace/src/github.com/jinzhu/gorm"
	"time"
)

// Notifier is an interface for records that can produce notifications
type Notifier interface {
	GetID() uint64
	GetType() string
	GetRecordUserID(*gorm.DB) (uint64, error) // kinda sucks this require an extra DB call, but c'est la vie
}

// Notification is what it sounds like
type Notification struct {
	ID           uint64    `json:"id" gorm:"primary_key"` // Primary key
	UserID       uint64    `json:"user_id"`
	TargetUserID uint64    `json:"target_user_id" binding:"required"` // ID of user to recieve the notification
	RecordID     uint64    `json:"record_id"`
	RecordType   string    `json:"record_type"`
	Seen         bool      `json:"seen"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// NewNotification initializes and saves a new notification
func NewNotification(n Notifier, initiatedUserID uint64, db *gorm.DB) error {
	targetUserID, err := n.GetRecordUserID(db)
	if err != nil {
		return err
	}

	if initiatedUserID == targetUserID {
		return fmt.Errorf("Can't send notification to the same user: %d", targetUserID)
	}

	notif := Notification{
		UserID:       initiatedUserID,
		TargetUserID: targetUserID,
		RecordID:     n.GetID(),
		RecordType:   n.GetType(),
		Seen:         false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return db.Save(&notif).Error
}

// GetRecord returns the related record of the Notification
func (n *Notification) GetRecord(db *gorm.DB) (Notifier, error) {
	switch n.RecordType {
	case "like":
		var l Like
		err := db.Find(&l, n.RecordID).Error
		return &l, err
	case "comment":
		var c Comment
		err := db.Find(&c, n.RecordID).Error
		return &c, err
	case "friendships":
		var f Friendship
		err := db.Find(&f, n.RecordID).Error
		return &f, err
	}

	var r Notifier
	return r, errors.New("Couldn't find record, possibly unsupported record type")
}
