package notifications

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/nu7hatch/gouuid"

	models "partisan/models.v2"
	"partisan/models.v2/users"
)

type Notifier interface {
	GetID() string
	GetAction() string
}

type BulkNotifier interface {
	Notifier
	GetNotifUserIDs() []string
}

type Notification struct {
	ID        string    `json:"id" gorm:"primary_key" sql:"type:uuid;default:uuid_generate_v4()"`
	UserID    string    `json:"user_id" sql:"type:uuid"` // ID of first user
	ToID      string    `json:"to_id" sql:"type:uuid"`   // ID of second user
	RecordID  string    `json:"record_id"`
	Action    Action    `json:"action"`
	Read      bool      `json:"read"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

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
		UpdatedAt: time.Now(),
	}

	errs = n.Validate()
	return
}

func NewBulk(userID string, b BulkNotifier) (ns []Notification, errs models.ValidationErrors) {
	for _, toID := range b.GetNotifUserIDs() {
		if userID != toID {
			n := Notification{
				UserID:    userID,
				ToID:      toID,
				RecordID:  b.GetID(),
				Action:    Action(b.GetAction()),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			if errs = n.Validate(); len(errs) > 0 {
				return ns, errs
			}

			ns = append(ns, n)
		}
	}

	return
}

func (n *Notification) Validate() (errs models.ValidationErrors) {
	errs = make(models.ValidationErrors)

	if _, err := uuid.ParseHex(n.UserID); err != nil {
		errs["user_id"] = models.ErrUUIDFormat
	}

	if _, err := uuid.ParseHex(n.ToID); err != nil {
		errs["to_id"] = models.ErrUUIDFormat
	}

	if n.UserID == n.ToID {
		errs["user_id"] = ErrNotifySelf
	}

	return errs
}

func BulkNotify(userID string, b BulkNotifier, db *gorm.DB) {
	bufDB := db.New()

	go func() {
		ns, errs := NewBulk(userID, b)
		if len(errs) > 0 {
			return
		}

		for _, n := range ns {
			bufDB.Create(&n)
		}
	}()
}
