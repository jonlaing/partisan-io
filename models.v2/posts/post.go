package posts

import (
	"time"

	"github.com/jinzhu/gorm"

	"partisan/models.v2/user"
)

type Action int

const (
	NoAction Action = iota
	NewPost
	Comment
	Like
)

type Post struct {
	ID        string    `json:"id" gorm:"primary_key" sql:"type:uuid;default:uuid_generate_v4()"`
	UserID    string    `json:"user_id"`
	User      user.User `json:"user" sql:"-"`
	ParentID  string    `json:"-"`
	Body      string    `json:"body"`
	Action    Action    `json:"action"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreatorBinding struct {
	UserID   string `json:"user_id" binding:"required"`
	ParentID string `json:"parent_id" binding:"required"`
	Body     string `json:"body" binding:"required"`
	Action   Action `json:"action" binding:"required"`
}

type UpdaterBinding struct {
	Body string `json:"body" binding:"required"`
}

func New(b CreatorBinding) (p Post, err error) {
	p.UserID = b.UserID
	p.ParentID = b.ParentID
	p.Body = b.Body
	p.Action = b.Action
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

	return
}

func (p Post) GetParent(db *gorm.DB) (parent Post, err error) {
	err = db.Find(&parent, p.ParentID).Error
	return
}
