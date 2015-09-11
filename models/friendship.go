package models

import(
  "time"
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