package messages

import (
	"partisan/models.v2/users"

	"github.com/jinzhu/gorm"
)

// GetThread will get the thread or throw an error if nothing is found
func GetThread(threadID string, db *gorm.DB) (thread Thread, err error) {
	err = db.Where("id = ?", threadID).First(&thread).Error
	return
}

// ListThreads returns a list of threads
func ListThreads(userID string, db *gorm.DB) (threads []Thread, err error) {
	err = db.Joins("INNER JOIN thread_users on thread_users.thread_id = threads.id").
		Where("thread_users.user_id = ?", userID).
		Order("threads.updated_at DESC").
		Find(&threads).Error
	return
}

// GetThreadUsers returns all the ThreadUsers for a user
func GetThreadUsers(userID string, db *gorm.DB) (mtus ThreadUsers, err error) {
	threadIDs, err := GetThreadIDs(userID, db)
	if err != nil {
		return
	}

	if len(threadIDs) == 0 {
		return mtus, ErrThreadsNotFound
	}

	err = db.Where("user_id != ? AND thread_id IN (?)", userID, threadIDs).Find(&mtus).Error
	if err != nil {
		return
	}

	if len(mtus) == 0 {
		return mtus, ErrThreadUsersNotFound
	}

	users, err := users.ListRelated(mtus, db)
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
func GetUserIDs(userID string, db *gorm.DB) (userIDs []string, err error) {
	threadIDs, err := GetThreadIDs(userID, db)
	if err != nil {
		return
	}

	if len(threadIDs) == 0 {
		return userIDs, ErrThreadsNotFound
	}

	err = db.Where("user_id != ? AND thread_id IN (?)", userID, threadIDs).Pluck("user_id", &userIDs).Error
	if err != nil {
		return
	}

	if len(userIDs) == 0 {
		err = ErrThreadUsersNotFound
	}

	return
}

// GetThreadIDs returns all the threadIDs for a particular user
func GetThreadIDs(userID string, db *gorm.DB) (threadIDs []string, err error) {
	err = db.Model(ThreadUser{}).
		Where("user_id = ?", userID).
		Pluck("thread_id", &threadIDs).Error
	return
}

// HasUnread returns true if there are unread messages for a thread
func HasUnread(userID, threadID string, db *gorm.DB) (bool, error) {
	var count int
	err := db.Model(Message{}).Where("user_id != ? AND thread_id = ? AND read = ?", userID, threadID, false).Count(&count).Error
	return count > 0, err
}

// HasUser returns true if a user is participating in a particular thread
func HasUser(userID, threadID string, db *gorm.DB) (bool, error) {
	var count int
	err := db.Model(ThreadUser{}).Where("thread_id = ? AND user_id = ?", threadID, userID).Count(&count).Error
	return count > 0, err
}

// GetByUsers returns a thread based on a group of users
func GetByUsers(db *gorm.DB, userIDs ...string) (thread Thread, err error) {
	err = db.Joins("INNER JOIN thread_users ON thread_users.thread_id = threads.id").
		Where("thread_users.user_id IN (?)", userIDs).
		First(&thread).Error
	if err != nil {
		return
	}

	var users []ThreadUser
	err = db.Where("thread_id = ?", thread.ID).Find(&users).Error
	if err != nil {
		return
	}

	count := 0
	for _, u := range users {
		for _, id := range userIDs {
			if id == u.UserID {
				count++
			}
		}
	}

	if count == 0 {
		err = ErrThreadsNotFound
		return
	}

	if count != len(userIDs) {
		err = ErrThreadUnreciprocated
		return
	}

	return
}

// LastMessage returns the most recent message in the thread
func LastMessage(threadID string, db *gorm.DB) (msg Message, err error) {
	err = db.Where("thread_id = ?", threadID).Limit(1).Order("created_at DESC").Find(&msg).Error
	return
}
