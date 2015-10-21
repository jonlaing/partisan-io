package dao

import (
	m "partisan/models"

	"partisan/Godeps/_workspace/src/github.com/jinzhu/gorm"
)

func GetFeedByUserIDs(currentUserID uint64, userIDs []uint64, db *gorm.DB) (feedItems []m.FeedItem, err error) {
	// TODO: limit feed so a particular record only comes up once
	if err = db.Where("user_id IN (?) AND action = ?", userIDs, "post").Order("created_at desc").Limit(25).Find(&feedItems).Error; err != nil {
		return
	}

	records, _ := GetRelatedPostReponse(currentUserID, m.FeedItems(feedItems), db)
	for i := 0; i < len(feedItems); i++ {
		feedItems[i].Record = records[feedItems[i].RecordID]
	}

	return
}
