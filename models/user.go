package models

import (
	"partisan/matcher"
	"time"

	"github.com/jasonmoo/geo"
)

const (
	// UTBasic is a basic user
	UTBasic = iota
	// UTAdvertiser is an advertiser
	UTAdvertiser
	// UTAdmin is an admin
	UTAdmin
)

// UserIDer is an interface for anything that can be owned by a user
type UserIDer interface {
	GetUserID() uint64
}

// User the user model
type User struct {
	ID                 uint64               `form:"id" json:"id" gorm:"primary_key"`
	Username           string               `form:"username" json:"username" sql:"not null;unique_index" binding:"required"`
	Email              string               `form:"email" json:"email" sql:"not null;unique_index" binding:"required"`
	Gender             string               `form:"gender" json:"gender"`
	Birthdate          time.Time            `form:"birthdate" json:"birthdate"`
	AvatarURL          string               `form:"avatar_url" json:"avatar_url"`
	AvatarThumbnailURL string               `form:"avatar_thumbnail_url" json:"avatar_thumbnail_url"`
	PoliticalMap       matcher.PoliticalMap `json:"political_map" sql:"type:varchar(255)"`
	CenterX            int                  `json:"center_x"`
	CenterY            int                  `json:"center_y"`
	PostalCode         string               `form:"postal_code" json:"postal_code"`
	Location           string               `json:"location"`
	Longitude          float64              `json:"-"`
	Latitude           float64              `json:"-"`
	CreatedAt          time.Time            `json:"created_at"`
	UpdatedAt          time.Time            `json:"updated_at"`
	APIKey             string               `json:"-"`
	APIKeyExp          time.Time            `json:"-"`
	PasswordHash       []byte               `json:"-"`
	Password           string               `form:"password" json:"password" sql:"-" binding:"required"`
	PasswordConfirm    string               `form:"password_confirm" json:"password_confirm" sql:"-" binding:"required"`
	Type               int                  `json:"-"`
}

// GetLocation finds the latitude/longitude by postal code
func (u *User) GetLocation() error {
	address, err := geo.Geocode(u.PostalCode)
	if err != nil {
		return err
	}

	u.Location = address.Address
	u.Latitude = address.Lat
	u.Longitude = address.Lng

	return nil
}

// CanDelete determines whether a user is allowed to delete a record
func (u *User) CanDelete(r UserIDer) bool {
	if u.Type == UTAdmin {
		// admins can always delete or modify
		return true
	}

	// otherwise they can only delete if it's theirs
	return r.GetUserID() == u.ID
}
