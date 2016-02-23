package models

import (
	"fmt"

	"partisan/Godeps/_workspace/src/github.com/jinzhu/gorm"
)

// Like is polymorphic
type Like struct {
	ID         uint64 `gorm:"primary_key"`
	UserID     uint64
	RecordID   uint64
	RecordType string
	IsDislike  bool // so that we can use the same table for both likes and dislikes, not currently in use
}

// RecordLikes stores like data for ease
type RecordLikes struct {
	RecordID  uint64
	Count     int
	UserCount int
}

// GetID satifies Notifier interface
func (l *Like) GetID() uint64 {
	return l.ID
}

// GetType satisifies Notifier interface
func (l *Like) GetType() string {
	return "like"
}

// GetRecordUserID returns the user ID of the record being liked. Satisfies Notifier interface.
func (l *Like) GetRecordUserID(db *gorm.DB) (uint64, error) {
	var notifUserIDs []uint64
	table := fmt.Sprintf("%ss", l.RecordType) // yeah, this is pretty naive, but no need to over engineer at this point
	err := db.Table(table).Where("id = ?", l.RecordID).Pluck("user_id", &notifUserIDs).Error

	if len(notifUserIDs) < 1 {
		return 0, err
	}

	return notifUserIDs[0], err
}

// GetLikes retrives Likes for a group of records
// PENDING DEPRECATION! See partisan/dao
func GetLikes(uID uint64, recordType string, recordIDs []uint64, db *gorm.DB) ([]RecordLikes, error) {
	fmt.Println("Warning: partisan/models.GetLikes is pending deprecation. See partisan/dao")
	var likes []RecordLikes

	// gorm seems to be having trouble with too many variables
	prepare := fmt.Sprintf("SELECT count(*), sum(case when user_id = ? then 1 else 0 end), record_id FROM \"likes\"  WHERE (record_type = '%s' AND record_id IN (?)) GROUP BY record_id", recordType)
	rows, err := db.Raw(prepare, uID, recordIDs).Rows()
	defer rows.Close()
	if err != nil {
		return []RecordLikes{}, err
	}

	for rows.Next() {
		var count int     // number of likes
		var userCount int // number of likes associated with particular user, should be between 0 and 1
		var rID uint64    // record ID
		rows.Scan(&count, &userCount, &rID)
		likes = append(likes, RecordLikes{Count: count, UserCount: userCount, RecordID: rID})
	}

	return likes, nil
}
