package pwreset

import "time"

type PasswordReset struct {
	ID         string `gorm:"primary_key" sql:"type:uuid;default:uuid_generate_v4()"`
	UserID     string `sql:"type:uuid"`
	Expiration time.Time
}

func New(userID string) PasswordReset {
	return PasswordReset{
		UserID:     userID,
		Expiration: time.Now().Add(24 * time.Hour),
	}
}

func (p PasswordReset) IsValid() bool {
	return time.Now().Before(p.Expiration)
}
