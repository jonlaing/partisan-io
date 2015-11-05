package dao

import (
	"errors"
	"time"

	m "partisan/models"

	"partisan/Godeps/_workspace/src/github.com/jinzhu/gorm"
)

func MessageUnreadCount(userID uint64, db *gorm.DB) (count int, err error) {
	threadIDs, err := GetMessageThreadIDs(userID, db)
	if err != nil {
		return
	}

	if len(threadIDs) == 0 {
		err = errors.New("No threads found")
		return
	}

	err = db.Model(m.Message{}).Where("user_id != ? AND thread_id IN (?) AND read = ?", userID, threadIDs, false).Count(&count).Error

	return
}

func GetMessages(threadID uint64, db *gorm.DB) (msgs []m.Message, err error) {
	err = db.Where("thread_id = ?", threadID).Find(&msgs).Error
	if err != nil {
		return
	}

	// If nothing was found, just return
	if len(msgs) == 0 {
		return
	}

	// otherwise, link up the users with the mssages
	users, err := GetRelatedUsers(m.Messages(msgs), db)
	if err != nil {
		return
	}

	collectMessageUsers(msgs, users)

	return
}

func GetMessagesAfter(threadID uint64, after time.Time, db *gorm.DB) (msgs []m.Message, err error) {
	err = db.Where("thread_id = ? AND created_at >= ?::timestamp", threadID, after).Find(&msgs).Error
	if err != nil {
		return
	}

	// If nothing was found, just return
	if len(msgs) == 0 {
		return
	}

	// otherwise, link up the users with the mssages
	users, err := GetRelatedUsers(m.Messages(msgs), db)
	if err != nil {
		return
	}

	collectMessageUsers(msgs, users)

	return
}

func collectMessageUsers(msgs []m.Message, users []m.User) {
	for k := range msgs {
		for _, u := range users {
			if u.ID == msgs[k].UserID {
				msgs[k].User = u
			}
		}
	}
}
