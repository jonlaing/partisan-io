package tickets

import (
	"time"

	"github.com/jinzhu/gorm"
)

// SocketTicket exists because the Javascript WebSocket API doesn't allow
// headers to be modified. This isn't a problem on the web, as we have cookies,
// but on mobile we usuall send a header with an authentication token.
// The SocketTicket acts as a temporary authentication to allow authenticated
// websocket transactions.
type SocketTicket struct {
	ID      string    `json:"id" gorm:"primary_key" sql:"type:uuid;default:uuid_generate_v4()"`
	UserID  string    `json:"user_id"`
	Expires time.Time `json:"expires"`
}

// NewSocketTicket generates a new socket with a key
func NewSocketTicket(userID string) SocketTicket {
	return SocketTicket{
		UserID:  userID,
		Expires: time.Now().Add(24 * time.Hour),
	}
}

// IsValid returns whether a ticket is still good. They expire after 24 hours
func (s *SocketTicket) IsValid() bool {
	return s.Expires.After(time.Now())
}

func GetByID(id string, db *gorm.DB) (t SocketTicket, err error) {
	err = db.Where("id = ?", id).Find(&t).Error
	return
}

func GetByUserID(userID string, db *gorm.DB) (t SocketTicket, err error) {
	err = db.Where("user_id = ?", userID).Find(&t).Error
	return
}
