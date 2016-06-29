package messages

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/nu7hatch/gouuid"

	models "partisan/models.v2"
	"partisan/models.v2/notifications"
	"partisan/models.v2/users"
)

// Message is a direct message. Its "to" User is set by the MessageThread and MessageThreadUser
type Message struct {
	ID        string    `json:"id" gorm:"primary_key" sql:"type:uuid;default:uuid_generate_v4()"`
	ThreadID  string    `json:"thread_id" sql:"type:uuid"`
	UserID    string    `json:"user_id" sql:"type:uuid"` // User that sent the message
	Body      string    `json:"body"`                    // The text based body
	Read      bool      `json:"read"`                    // has the message been read?
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User users.User `json:"user" sql:"-"`
}

type Messages []Message

// GetUserIDs satisfies dao.UserIDerSlice
func (ms Messages) GetUserIDs() (userIDs []string) {
	for _, m := range ms {
		userIDs = append(userIDs, m.UserID)
	}

	return
}

func (ms Messages) Len() int {
	return len(ms)
}

func (ms Messages) Less(i, j int) bool {
	return ms[i].CreatedAt.Before(ms[j].CreatedAt)
}

func (ms Messages) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}

type MessageCreatorBinding struct {
	ThreadID string `json:"thread_id"`
	UserID   string `json:"user_id"`
	Body     string `json:"body" binding:"required"`
}

func NewMessage(b MessageCreatorBinding) (Message, models.ValidationErrors) {
	m := Message{
		ThreadID:  b.ThreadID,
		UserID:    b.UserID,
		Body:      b.Body,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return m, m.Validate()
}

func (m Message) Validate() (errs models.ValidationErrors) {
	errs = make(models.ValidationErrors)

	if _, err := uuid.ParseHex(m.UserID); err != nil {
		errs["user_id"] = models.ErrUUIDFormat
	}

	if _, err := uuid.ParseHex(m.ThreadID); err != nil {
		errs["friend_id"] = models.ErrUUIDFormat
	}

	return
}

func (m Message) NewPushNotifications(db *gorm.DB) (pns []notifications.PushNotification, err error) {
	thread, err := GetThread(m.ThreadID, db)
	if err != nil {
		return pns, err
	}

	var from users.User
	for _, mtu := range thread.Users {
		if mtu.UserID == m.UserID {
			from = mtu.User
		}
	}

	for _, mtu := range thread.Users {
		if mtu.UserID != m.UserID {
			pns = append(pns, notifications.PushNotification{
				DeviceToken:    mtu.User.DeviceToken,
				Message:        fmt.Sprintf("@%s: %s", from.Username, m.Body),
				NotificationID: m.ID,
			})
		}
	}

	return pns, nil
}
