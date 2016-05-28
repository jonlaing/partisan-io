package messages

import (
	"errors"
	"sort"
	"time"

	"partisan/models.v2/users"

	"github.com/jinzhu/gorm"
)

func UnreadCount(userID string, db *gorm.DB) (count int, err error) {
	threadIDs, err := GetThreadIDs(userID, db)
	if err != nil {
		return
	}

	if len(threadIDs) == 0 {
		err = errors.New("No threads found")
		return
	}

	err = db.Model(Message{}).Where("user_id != ? AND thread_id IN (?) AND read = ?", userID, threadIDs, false).Count(&count).Error

	return
}

func GetMessages(threadID string, db *gorm.DB) (msgs Messages, err error) {
	err = db.Where("thread_id = ?", threadID).Limit(200).Order("created_at DESC").Find(&msgs).Error
	if err != nil {
		return
	}

	// If nothing was found, just return
	if len(msgs) == 0 {
		return
	}

	// otherwise, link up the users with the mssages
	msgs.GetUsers(db)

	sort.Sort(msgs)

	return
}

func GetMessagesAfter(threadID string, after time.Time, db *gorm.DB) (msgs Messages, err error) {
	err = db.Where("thread_id = ? AND created_at >= ?::timestamp", threadID, after).Order("created_at DESC").Find(&msgs).Error
	if err != nil {
		return
	}

	// If nothing was found, just return
	if len(msgs) == 0 {
		return
	}

	// otherwise, link up the users with the mssages
	msgs.GetUsers(db)

	sort.Sort(msgs)

	return
}

func MarkAllMessagesRead(userID, threadID string, db *gorm.DB) error {
	return db.Table("messages").Where("user_id != ? AND thread_id = ?", userID, threadID).UpdateColumn("read", true).Error
}

func (ms *Messages) GetUsers(db *gorm.DB) error {
	users, err := users.ListRelated(ms, db)
	if err != nil {
		return err
	}

	if len(users) == 0 {
		return ErrNoUsers
	}

	msgs := []Message(*ms)
	for i := range msgs {
		for _, u := range users {
			if msgs[i].UserID == u.ID {
				msgs[i].User = u
			}
		}
	}

	*ms = Messages(msgs)
	return nil
}
