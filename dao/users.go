package dao

import (
	m "partisan/models"

	"partisan/Godeps/_workspace/src/github.com/jinzhu/gorm"
)

type UserIDerSlice interface {
	GetUserIDs() []uint64
}

type UserIDer interface {
	GetUserID() uint64
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
