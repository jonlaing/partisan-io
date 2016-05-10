package post

import (
	"time"

	"partisan/models.v2/user"
)

type Action int

const (
	NoAction Action = iota
	NewPost
	Commnet
	Like
)

type Post struct {
	ID        uint64    `json:"id" gorm:"primary_key"`
	UserID    uint64    `json:"user_id"`
	User      user.User `json:"user" sql:"-"`
	ParentID  uint64    `json:"-"`
	Body      string    `json:"body"`
	Action    Action    `json:"action"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreatorBinding struct {
	UserID   uint64 `json:"user_id" binding:"required"`
	ParentID uint64 `json:"parent_id" binding:"required"`
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
