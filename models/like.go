package models

import (
	"github.com/jinzhu/gorm"
)

// Like is polymorphic
type Like struct {
	ID         uint64 `gorm:"primary_key"`
	UserID     uint64
	RecordID   uint64
	RecordType string
	IsDislike  bool // so that we can use the same table for both likes and dislikes, not currently in use
}

// GetID satifies Notifier interface
func (l *Like) GetID() uint64 {
	return l.ID
}

// Type satisifies Notifier interface
func (l *Like) Type() string {
	return "like"
}

// GetRecordUserID returns the user ID of the record being liked. Satisfies Notifier interface.
func (l *Like) GetRecordUserID(db *gorm.DB) (uint64, error) {
	var notifUserIDs []uint64
	err := db.Table(l.RecordType).Where("id = ?", l.RecordID).Pluck("user_id", &notifUserIDs).Error

	if len(notifUserIDs) < 1 {
		return 0, err
	}

	return notifUserIDs[0], err
}
