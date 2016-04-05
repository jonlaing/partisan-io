package dao

import (
	m "partisan/models"
	"time"

	"github.com/jinzhu/gorm"
)

func GetFeedByUserIDs(currentUserID uint64, userIDs []uint64, page int, db *gorm.DB) (feedItems []m.FeedItem, err error) {
	offset := (page - 1) * 25

	// TODO: limit feed so a particular record only comes up once
	if err = db.Where("user_id IN (?) AND action = ?", userIDs, "post").Order("created_at desc").Limit(25).Offset(offset).Find(&feedItems).Error; err != nil {
		return feedItems, &ErrNotFound{err}
	}

	records, _ := GetRelatedPostReponse(currentUserID, m.FeedItems(feedItems), db)
	for i := 0; i < len(feedItems); i++ {
		feedItems[i].Record = records[feedItems[i].RecordID]
	}

	return feedItems, nil
}

func GetFeedByUserIDsAfter(currentUserID uint64, userIDs []uint64, after time.Time, db *gorm.DB) (feedItems []m.FeedItem, err error) {
	// TODO: limit feed so a particular record only comes up once
	if err = db.Where("user_id IN (?) AND action = ? AND created_at >= ?::timestamp", userIDs, "post", after).Order("created_at desc").Find(&feedItems).Error; err != nil {
		return
	}

	records, _ := GetRelatedPostReponse(currentUserID, m.FeedItems(feedItems), db)
	for i := 0; i < len(feedItems); i++ {
		feedItems[i].Record = records[feedItems[i].RecordID]
	}

	return
}
