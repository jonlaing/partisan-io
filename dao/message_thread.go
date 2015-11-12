package dao

import (
	"errors"
	"fmt"
	m "partisan/models"

	"partisan/Godeps/_workspace/src/github.com/jinzhu/gorm"
)

type MessageThreadUnreciprocated struct {
	ThreadID uint64
}

func (e *MessageThreadUnreciprocated) Error() string {
	return fmt.Sprintf("Need at least two MessageThreadUsers for Thread: %d", e.ThreadID)
}

// GetMessageThread will get the message or throw an error if nothing is found
func GetMessageThread(threadID uint64, db *gorm.DB) (thread m.MessageThread, err error) {
	err = db.Where("id = ?", threadID).First(&thread).Error
	return
}

func GetMessageThreads(userID uint64, db *gorm.DB) (threads []m.MessageThread, err error) {
	err = db.Joins("INNER JOIN message_thread_users on message_thread_users.thread_id = message_threads.id").
		Where("message_thread_users.user_id = ?", userID).
		Order("message_threads.updated_at DESC").
		Find(&threads).Error

	return
}

func GetMessageThreadUsers(userID uint64, db *gorm.DB) (mtus []m.MessageThreadUser, err error) {
	threadIDs, err := GetMessageThreadIDs(userID, db)
	if err != nil {
		return
	}

	if len(threadIDs) == 0 {
		return mtus, errors.New("No threads found")
	}

	err = db.Where("user_id != ? AND thread_id IN (?)", userID, threadIDs).Find(&mtus).Error
	if err != nil {
		return
	}

	if len(mtus) == 0 {
		return mtus, errors.New("No MessageThreadUsers found")
	}

	users, err := GetRelatedUsers(m.MessageThreadUsers(mtus), db)
	if err != nil {
		return
	}

	for k := range mtus {
		for _, u := range users {
			if u.ID == mtus[k].UserID {
				mtus[k].User = u
			}
		}
	}

	return
}

func GetMessageThreadUserIDs(userID uint64, db *gorm.DB) (userIDs []uint64, err error) {
	threadIDs, err := GetMessageThreadIDs(userID, db)
	if err != nil {
		return
	}

	if len(threadIDs) == 0 {
		return userIDs, errors.New("No threads found")
	}

	err = db.Where("user_id != ? AND thread_id IN (?)", userID, threadIDs).Pluck("user_id", &userIDs).Error
	if err != nil {
		return
	}

	if len(userIDs) == 0 {
		return userIDs, errors.New("No MessageThreadUsers found")
	}

	return
}

func GetMessageThreadIDs(userID uint64, db *gorm.DB) (threadIDs []uint64, err error) {
	err = db.Model(m.MessageThreadUser{}).
		Where("user_id = ?", userID).
		Pluck("thread_id", &threadIDs).Error

	return
}

func MessageThreadHasUnread(userID, threadID uint64, db *gorm.DB) (bool, error) {
	var count int
	err := db.Model(m.Message{}).Where("user_id != ? AND thread_id = ? AND read = ?", userID, threadID, false).Count(&count).Error

	return count > 0, err
}

func MessageThreadHasUser(userID, threadID uint64, db *gorm.DB) (bool, error) {
	var count int
	err := db.Model(m.MessageThreadUser{}).Where("thread_id = ? AND user_id = ?", threadID, userID).Count(&count).Error
	return count > 0, err
}

func GetMessageThreadByUsers(userID, toID uint64, db *gorm.DB) (thread m.MessageThread, err error) {
	var count int

	err = db.Joins("INNER JOIN message_thread_users ON message_thread_users.thread_id = message_threads.id").
		Where("message_thread_users.user_id IN (?)", []uint64{userID, toID}).
		First(&thread).Error
	if err != nil {
		return
	}

	err = db.Model(m.MessageThreadUser{}).Where("thread_id = ? AND user_id IN (?)", thread.ID, []uint64{userID, toID}).Count(&count).Error
	if err != nil {
		return
	}

	if count < 2 {
		return thread, &MessageThreadUnreciprocated{thread.ID}
	}

	return
}
