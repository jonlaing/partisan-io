package dao

import (
	m "partisan/models"

	"partisan/Godeps/_workspace/src/github.com/jinzhu/gorm"
)

func GetThread(threadID uint64, db *gorm.DB) (thread m.MessageThread, err error) {
	err = db.Where("id = ?", threadID).Find(&thread).Error
	return
}

func GetMessageThreads(userID uint64, db *gorm.DB) (threads []m.MessageThread, err error) {
	err = db.Joins("LEFT JOIN messages on messages.thread_id = message_threads.id").
		Where("messages.user_id = ?", userID).
		Find(&threads).Error

	return
}

func GetMessageThreadIDs(userID uint64, db *gorm.DB) (threadIDs []uint64, err error) {
	err = db.Model(m.MessageThread{}).
		Joins("LEFT JOIN messages on messages.thread_id = message_threads.id").
		Where("messages.user_id = ?", userID).
		Pluck("message_threads.id", &threadIDs).Error

	return
}

func MessageThreadHasUnread(threadID uint64, db *gorm.DB) (bool, error) {
	var count int
	err := db.Model(m.Message{}).Where("thread_id = ? AND read = ?", threadID, false).Count(&count).Error
	return count > 0, err
}

func MessageThreadHasUser(userID, threadID uint64, db *gorm.DB) (bool, error) {
	var count int
	err := db.Model(m.Message{}).Where("thread_id = ? AND user_id = ?", threadID, userID).Count(&count).Error
	return count > 0, err
}
