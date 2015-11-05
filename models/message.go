package models

import "time"

// Message is a direct message. Its "to" User is set by the MessageThread and MessageThreadUser
type Message struct {
	ID        uint64    `json:"id" gorm:"primary_key"`      // Primary key
	UserID    uint64    `json:"user_id" binding:"required"` // ID of user that created Message
	ThreadID  uint64    `json:"thread_id" binding:"required"`
	User      User      `sql:"-"`
	Body      string    `json:"body" binding:"required"` // The text based body
	Read      bool      `json:"read"`                    // has the message been read?
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Messages []Message

// MessageThread holds all messages between two (or more, eventually) Users
type MessageThread struct {
	ID        uint64    `json:"id" gorm:"primary_key"` // Primary key
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetUserIDs satisfies dao.UserIDerSlice
func (ms Messages) GetUserIDs() (userIDs []uint64) {
	for _, m := range ms {
		userIDs = append(userIDs, m.UserID)
	}

	return
}
