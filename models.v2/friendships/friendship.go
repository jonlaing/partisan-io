package friendships

import (
	"database/sql"
	"time"

	"github.com/nu7hatch/gouuid"

	models "partisan/models.v2"
	"partisan/models.v2/users"
)

// Friendship is a joining table between two users who are friends. For each
// friendship there are two Friendship records: one for each user
type Friendship struct {
	ID        string     `json:"id" gorm:"primary_key" sql:"type:uuid;default:uuid_generate_v4()"`
	UserID    string     `json:"user_id" sql:"type:uuid"`   // ID of first user
	FriendID  string     `json:"friend_id" sql:"type:uuid"` // ID of second user
	Friend    users.User `json:"user" sql:"-"`
	Match     float64    `json:"match" sql:"-"`
	Confirmed bool       `json:"confirmed"` // Whether the second user has confirmed
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type Friendships []Friendship

type CreatorBinding struct {
	FriendID string `json:"friend_id" binding:"required"`
}

type UpdaterBinding struct {
	Confirmed bool `json:"confirmed" binding:"required"`
}

func New(userID string, b CreatorBinding) (f Friendship, errs models.ValidationErrors) {
	f.UserID = userID
	f.FriendID = b.FriendID
	f.Confirmed = false
	f.CreatedAt = time.Now()
	f.UpdatedAt = f.CreatedAt

	errs = f.Validate()
	return
}

func (f *Friendship) Update(b UpdaterBinding) models.ValidationErrors {
	f.Confirmed = b.Confirmed
	return f.Validate()
}

func (f Friendship) Validate() (errs models.ValidationErrors) {
	errs = make(models.ValidationErrors)

	if _, err := uuid.ParseHex(f.UserID); err != nil {
		errs["user_id"] = models.ErrUUIDFormat
	}

	if _, err := uuid.ParseHex(f.FriendID); err != nil {
		errs["friend_id"] = models.ErrUUIDFormat
	}

	if f.UserID == f.FriendID {
		errs["friend_id"] = ErrFriendSelf
	}

	return errs
}

func (f Friendship) CanDelete(userID) bool {
	return f.UserID == userID || f.FriendID == userID
}

func (f Friendship) GetID() sql.NullString {
	if f.Confirmed {
		return sql.NullString{f.ToID, true}
	}

	return sql.NullString{f.UserID, true}
}

func (f Friendship) GetAction() string {
	if Confirmed {
		return string(AFriendAccept)
	}

	return string(AFriendRequest)
}
