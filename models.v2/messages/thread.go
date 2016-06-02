package messages

import (
	"time"

	"github.com/nu7hatch/gouuid"

	models "partisan/models.v2"
)

const (
	SOpen      = "opened"
	SBlocked   = "blocked"
	SEncrypted = "encrypted"
)

// Thread holds all messages between two (or more, eventually) Users
type Thread struct {
	ID        string    `json:"id" gorm:"primary_key" sql:"type:uuid;default:uuid_generate_v4()"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Users       ThreadUsers `json:"users,omitempty" sql:"-"`
	LastMessage Message     `json:"last_message,omitempty" sql:"-"`
}

type Threads []Thread

type ThreadCreatorBinding struct {
	UserIDs []string `json:"user_ids" binding:"required"`
}

func NewThread(userID string, b ThreadCreatorBinding) (t Thread, errs models.ValidationErrors) {
	errs = make(models.ValidationErrors)

	id, err := uuid.NewV4()
	if err != nil {
		errs["id"] = err
		return
	}

	t = Thread{ID: id.String(), Status: SOpen, CreatedAt: time.Now(), UpdatedAt: time.Now()}

	b.UserIDs = append(b.UserIDs, userID)

	var mtus []ThreadUser
	for _, uid := range b.UserIDs {
		mtus = append(mtus, ThreadUser{
			ThreadID: id.String(),
			UserID:   uid,
		})
	}

	ids := make(map[string]bool)
	for _, mtu := range mtus {
		if _, ok := ids[mtu.UserID]; ok {
			errs["user_id"] = ErrMessageSelf
			return
		}

		ids[mtu.UserID] = true
	}

	t.Users = ThreadUsers(mtus)

	return
}
