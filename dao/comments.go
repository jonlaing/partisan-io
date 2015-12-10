package dao

import (
	m "partisan/models"

	"partisan/Godeps/_workspace/src/github.com/jinzhu/gorm"
)

// PostComments stores like data for ease
type PostComments struct {
	RecordID uint64
	Count    int
}

func GetRelatedComments(ps []m.Post, db *gorm.DB) ([]PostComments, error) {
	var postIDs []uint64
	for _, p := range ps {
		postIDs = append(postIDs, p.ID)
	}

	return GetRelatedCommentsByIDs(postIDs, db)
}

func GetRelatedCommentsByIDs(postIDs []uint64, db *gorm.DB) ([]PostComments, error) {
	var comments []PostComments

	rows, err := db.Raw("SELECT count(*), post_id FROM \"comments\"  WHERE (post_id IN (?)) group by post_id", postIDs).Rows()
	defer rows.Close()
	if err != nil {
		return []PostComments{}, err
	}

	for rows.Next() {
		var count int
		var rID uint64

		rows.Scan(&count, &rID)
		comments = append(comments, PostComments{Count: count, RecordID: rID})
	}

	return comments, nil
}
