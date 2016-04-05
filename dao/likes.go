package dao

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// RecordLikes stores like data for ease
type RecordLikes struct {
	RecordID  uint64
	Count     int
	UserCount int
}

type Liker interface {
	GetRecordTypeAndIDs() (t string, ids []uint64, err error)
}

func GetRelatedLikes(uID uint64, rs Liker, db *gorm.DB) (likes []RecordLikes, err error) {
	t, ids, err := rs.GetRecordTypeAndIDs()

	prepare := fmt.Sprintf("SELECT count(*), sum(case when user_id = ? then 1 else 0 end), record_id FROM \"likes\"  WHERE (record_type = '%s' AND record_id IN (?)) GROUP BY record_id", t)
	rows, err := db.Raw(prepare, uID, ids).Rows()
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
