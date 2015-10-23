package models

import (
	"time"

	"partisan/Godeps/_workspace/src/github.com/jinzhu/gorm"
)

// Friendship is a joining table between two users who are friends. For each
// friendship there are two Friendship records: one for each user
type Friendship struct {
	ID        uint64    `json:"id" gorm:"primary_key"`    // Primary key
	UserID    uint64    `json:"user_id" sql:"not null"`   // ID of first user
	FriendID  uint64    `json:"friend_id" sql:"not null"` // ID of second user
	Confirmed bool      `json:"confirmed"`                // Whether the second user has confirmed
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetID satisfies Notifier interface
func (f *Friendship) GetID() uint64 {
	return f.ID
}

// Kinda wish I documented why I needed it to work this way...
func (f *Friendship) GetRecordUserID(db *gorm.DB) (uint64, error) {
	if !f.Confirmed {
		return f.FriendID, nil
	} else {
		return f.UserID, nil
	}
}

// GetType satisfies Notifier interface
func (f *Friendship) GetType() string {
	return "friendship"
}
