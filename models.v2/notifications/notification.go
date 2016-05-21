package notifications

import (
	"database/sql"
	"time"

	"github.com/nu7hatch/gouuid"

	models "partisan/models.v2"
	"partisan/models.v2/users"
)

type Notifier interface {
	GetID() sql.NullString
	GetAction() string
}

type Notification struct {
	ID        string         `json:"id" gorm:"primary_key" sql:"type:uuid;default:uuid_generate_v4()"`
	UserID    string         `json:"user_id" sql:"type:uuid"` // ID of first user
	ToID      string         `json:"to_id" sql:"type:uuid"`   // ID of second user
	RecordID  sql.NullString `json:"record_id"`
	Action    Action         `json:"action"`
	Read      bool           `json:"read"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`

	User users.User `json:"user" sql:"-"`
}

type Notifications []Notification

func New(userID, toID string, r Notifier) (n Notification, errs models.ValidationErrors) {
	n = Notification{
		UserID:    userID,
		ToID:      toID,
		RecordID:  r.GetID(),
		Action:    Action(r.GetAction()),
		CreatedAt: time.Now(),
		UpdatedAt: r.CreatedAt,
	}

	errs = n.Validate()
	return
}

func (n *Notification) Validate() (errs models.ValidationErrors) {
	errs = make(models.ValidationErrors)

	if _, err := uuid.ParseHex(f.UserID); err != nil {
		errs["user_id"] = models.ErrUUIDFormat
	}

	if _, err := uuid.ParseHex(f.ToID); err != nil {
		errs["to_id"] = models.ErrUUIDFormat
	}

	if f.UserID == f.ToID {
		errs["user_id"] = ErrNotifySelf
	}

	return errs
}
