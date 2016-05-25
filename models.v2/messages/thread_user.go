package messages

import (
	"time"

	"partisan/models.v2/users"
)

// MessageThreadUser is a join table between Message and MessageThread
type ThreadUser struct {
	ID        string    `json:"id" gorm:"primary_key" sql:"type:uuid;default:uuid_generate_v4()"`
	UserID    string    `json:"user_id" sql:"type:uuid"`
	ThreadID  string    `json:"thread_id" sql:"type:uuid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User users.User `json:"user" sql:"-"`
}

type ThreadUsers []ThreadUser

func (mtus ThreadUsers) GetUserIDs() (userIDs []string) {
	for _, mtu := range mtus {
		userIDs = append(userIDs, mtu.UserID)
	}

	return
}
