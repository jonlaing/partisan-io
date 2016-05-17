package users

import "github.com/jinzhu/gorm"

type UserIDerSlice interface {
	GetUserIDs() []string
}

type UserIDer interface {
	GetUserID() string
}

func GetByID(id string, db *gorm.DB) (u User, err error) {
	err = db.Where("id = ?", id).Find(&u).Error
	return u, err
}

func GetByUsername(username string, db *gorm.DB) (u User, err error) {
	err = db.Where("username = ?", username).Find(&u).Error
	return u, err
}

func GetByEmail(email string, db *gorm.DB) (u User, err error) {
	err = db.Where("email = ?", email).Find(&u).Error
	return u, err
}

func GetByAPIKey(key string, db *gorm.DB) (u User, err error) {
	err = db.Where("api_key = ?", key).Find(&u).Error
	if err == nil {
		err = u.ValidateAPIKey()
		return
	}

	return
}

func GetRelated(r UserIDer, us []User) (User, bool) {
	for _, u := range us {
		if u.ID == r.GetUserID() {
			return u, true
		}
	}
	return User{}, false
}

func ListByIDs(ids []string, db *gorm.DB) (us []User, err error) {
	err = db.Where("id IN (?)", ids).Find(&us).Error
	return
}

func ListRelated(rs UserIDerSlice, db *gorm.DB) (us []User, err error) {
	ids := rs.GetUserIDs()
	err = db.Where("id IN (?)", ids).Find(&us).Error
	return
}
