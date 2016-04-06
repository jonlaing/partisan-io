package dao

import (
	m "partisan/models"
	"testing"
)

func TestGetRelatedComments(t *testing.T) {
	var posts []m.Post
	var comments []m.Comment
	var idList []uint64
	for i := 0; i <= 10; i++ {
		idList = append(idList, uint64(i))
		posts = append(posts, m.Post{ID: uint64(i)})
		for j := 0; j < 5; j++ {
			comment := m.Comment{PostID: uint64(i)}
			db.Create(&comment)
			comments = append(comments, comment)
		}
	}

	pcs, err := GetRelatedComments(posts, db)
	if err != nil {
		t.Error("Error getting related comment counts from posts:", err)
		return
	}

	if len(pcs) < 10 {
		t.Error("Not enough post comments:", len(pcs))
	}

	for _, pc := range pcs {
		found := false
		for _, v := range idList {
			if pc.RecordID == v && pc.Count == 5 {
				found = true
			}
		}
		if !found {
			t.Error("Couldn't find id:", pc.RecordID, "in", idList, "with count of 5:", pc.Count)
		}
	}

	db.Delete(&comments)
}
