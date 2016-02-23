package models

// Profile of a user
type Profile struct {
	ID            uint64 `form:"id" json:"id" gorm:"primary_key"`
	UserID        uint64 `json:"user_id"`
	LookingFor    int    `form:"looking_for" json:"looking_for"`
	Summary       string `form:"summary" json:"summary"`
	CoverPhotoURL string `json:"cover_photo_url"`
}
