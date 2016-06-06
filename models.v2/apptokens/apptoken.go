package apptokens

import (
	"time"

	"github.com/jinzhu/gorm"
)

type AppToken struct {
	ID        string    `json:"id" gorm:"primary_key" sql:"type:uuid;default:uuid_generate_v4()"`
	Name      string    `json:"string"`
	Email     string    `json:"email" sql:"not null;unique_index" binding:"required"`
	Website   string    `json:"website"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func GetByID(id string, db *gorm.DB) (t AppToken, err error) {
	err = db.Where("id = ?", id).Find(&t).Error
	return
}
