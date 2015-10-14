package models

import (
	"partisan/Godeps/_workspace/src/github.com/jasonmoo/geo"
	"partisan/Godeps/_workspace/src/github.com/jinzhu/gorm"
	"partisan/matcher"
	"time"
)

// User the user model
type User struct {
	ID                 uint64               `form:"id" json:"id" gorm:"primary_key"`
	Username           string               `form:"username" json:"username" sql:"not null,unique" binding:"required"`
	Email              string               `form:"email" json:"email" sql:"not null,unique" binding:"required"`
	Gender             string               `form:"gender" json:"gender"`
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
}

// Friends returns all the User's friends
func (u User) Friends(db *gorm.DB) []User {
	friendIDs := u.FriendIDs(db)
	users := []User{}
	db.Where(friendIDs).Find(&users)
	return users
}

// FriendIDs returns all IDs of User's Friends
func (u User) FriendIDs(db *gorm.DB) []uint64 {
	var friendIDs []uint64
	db.Table("friendships").Select("friend_id").Where("confirmed = ?", true).Scan(&friendIDs)
	return friendIDs
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
