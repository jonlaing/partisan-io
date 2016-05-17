package models

import "time"

type PreSignUp struct {
	ID        uint64    `form:"id" json:"id" gorm:"primary_key"`
	Email     string    `form:"email" json:"email" sql:"not null;unique_index" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
