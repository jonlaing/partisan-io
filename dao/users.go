package dao

import (
	"errors"
	m "partisan/models"

	"partisan/Godeps/_workspace/src/github.com/jinzhu/gorm"
)

type UserIDerSlice interface {
	GetUserIDs() []uint64
}

type UserIDer interface {
	GetUserID() uint64
}

func GetUsersByIDs(userIDs []uint64, db *gorm.DB) (users []m.User, err error) {
	err = db.Where("id IN (?)", userIDs).Find(&users).Error
	if err != nil {
		return
	}

	if len(users) == 0 {
		return users, errors.New("No Users found")
	}

	return
}

func GetRelatedUsers(rs UserIDerSlice, db *gorm.DB) (users []m.User, err error) {
	userIDs := rs.GetUserIDs()
	err = db.Where("id IN (?)", userIDs).Find(&users).Error
	return
}

func GetMatchingUser(r UserIDer, users []m.User) (m.User, bool) {
	for _, user := range users {
		if user.ID == r.GetUserID() {
			return user, true
		}
	}
	return m.User{}, false
}
