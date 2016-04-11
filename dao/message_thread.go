package dao

import (
	"errors"
	m "partisan/models"

	"github.com/jinzhu/gorm"
)

type MessageThreadDAO struct {
	DB *gorm.DB
}

// Get will get the thread or throw an error if nothing is found
func (mt *MessageThreadDAO) Get(threadID uint64) (thread m.MessageThread, err error) {
	err = mt.DB.Where("id = ?", threadID).First(&thread).Error
	if err != nil {
		return thread, &ErrThreadNotFound{err}
	}

	return
}

// List returns a list of threads
func (mt *MessageThreadDAO) List(userID uint64) (threads []m.MessageThread, err error) {
	err = mt.DB.Joins("INNER JOIN message_thread_users on message_thread_users.thread_id = message_threads.id").
		Where("message_thread_users.user_id = ?", userID).
		Order("message_threads.updated_at DESC").
		Find(&threads).Error

	if err != nil && err == gorm.ErrRecordNotFound {
		err = &ErrThreadNotFound{}
	}

	return
}

// GetMessageThreadUsers returns all the MessageThreadUsers for a user
func (mt *MessageThreadDAO) GetMessageThreadUsers(userID uint64) (mtus []m.MessageThreadUser, err error) {
	threadIDs, err := mt.GetIDs(userID)
	if err != nil {
		return
	}

	if len(threadIDs) == 0 {
		return mtus, &ErrThreadNotFound{errors.New("No threads found")}
	}

	err = mt.DB.Where("user_id != ? AND thread_id IN (?)", userID, threadIDs).Find(&mtus).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return
	}

	if len(mtus) == 0 {
		return mtus, errors.New("No MessageThreadUsers found")
	}

	users, err := GetRelatedUsers(m.MessageThreadUsers(mtus), mt.DB)
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

// GetUserIDs returns a list of user ids with whom a user has a thread
func (mt *MessageThreadDAO) GetUserIDs(userID uint64) (userIDs []uint64, err error) {
	threadIDs, err := mt.GetIDs(userID)
	if err != nil {
		return
	}

	if len(threadIDs) == 0 {
		return userIDs, &ErrThreadNotFound{errors.New("No threads found")}
	}

	err = mt.DB.Where("user_id != ? AND thread_id IN (?)", userID, threadIDs).Pluck("user_id", &userIDs).Error
	if err != nil {
		return
	}

	if len(userIDs) == 0 {
		return userIDs, errors.New("No MessageThreadUsers found")
	}

	return
}

// GetIDs returns all the threadIDs for a particular user
func (mt *MessageThreadDAO) GetIDs(userID uint64) (threadIDs []uint64, err error) {
	err = mt.DB.Model(m.MessageThreadUser{}).
		Where("user_id = ?", userID).
		Pluck("thread_id", &threadIDs).Error

	if err != nil && err == gorm.ErrRecordNotFound {
		err = &ErrThreadNotFound{}
	}

	return
}

// HasUnread returns true if there are unread messages for a thread
func (mt *MessageThreadDAO) HasUnread(userID, threadID uint64) (bool, error) {
	var count int
	err := mt.DB.Model(m.Message{}).Where("user_id != ? AND thread_id = ? AND read = ?", userID, threadID, false).Count(&count).Error

	return count > 0, err
}

// HasUser returns true if a user is participating in a particular thread
func (mt *MessageThreadDAO) HasUser(userID, threadID uint64) (bool, error) {
	var count int
	err := mt.DB.Model(m.MessageThreadUser{}).Where("thread_id = ? AND user_id = ?", threadID, userID).Count(&count).Error
	if err != nil {
		return false, &ErrThreadNotFound{err}
	}

	return count > 0, nil
}

// GetByUsers returns a thread based on a group of users
func (mt *MessageThreadDAO) GetByUsers(userIDs ...uint64) (thread m.MessageThread, err error) {
	var count int

	err = mt.DB.Joins("INNER JOIN message_thread_users ON message_thread_users.thread_id = message_threads.id").
		Where("message_thread_users.user_id IN (?)", userIDs).
		First(&thread).Error
	if err != nil {
		return
	}

	err = mt.DB.Model(m.MessageThreadUser{}).Where("thread_id = ? AND user_id IN (?)", thread.ID, userIDs).Count(&count).Error
	if err != nil {
		return
	}

	if count < 2 {
		return thread, &MessageThreadUnreciprocated{thread.ID}
	}

	return
}

// LastMessage returns the most recent message in the thread
func (mt *MessageThreadDAO) LastMessage(threadID uint64) (msg m.Message, err error) {
	err = mt.DB.Where("thread_id = ?", threadID).Limit(1).Order("created_at DESC").Find(&msg).Error

	return msg, err
}
