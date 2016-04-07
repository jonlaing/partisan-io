package models

import (
	"time"

	uuid "github.com/nu7hatch/gouuid"
)

// SocketTicket exists because the Javascript WebSocket API doesn't allow
// headers to be modified. This isn't a problem on the web, as we have cookies,
// but on mobile we usuall send a header with an authentication token.
// The SocketTicket acts as a temporary authentication to allow authenticated
// websocket transactions.
type SocketTicket struct {
	ID      uint64    `json:"id" gorm:"primary_key"`
	UserID  uint64    `json:"user_id"`
	Key     string    `json:"key" sql:"not null;unique_index"`
	Expires time.Time `json:"expires"`
}

// NewSocketTicket generates a new socket with a key
func NewSocketTicket(userID uint64) (SocketTicket, error) {
	key, err := uuid.NewV4()
	if err != nil {
		return SocketTicket{}, err
	}

	return SocketTicket{
		UserID:  userID,
		Key:     key.String(),
		Expires: time.Now().Add(24 * time.Hour),
	}, nil
}

// IsValid returns whether a ticket is still good. They expire after 24 hours
func (s *SocketTicket) IsValid() bool {
	return s.Expires.After(time.Now())
}
