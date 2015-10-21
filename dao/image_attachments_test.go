package dao

import (
	m "partisan/models"
	"testing"
)

func TestGetRelatedAttachments(t *testing.T) {
	var posts m.Posts
	var idList []uint64
	for i := 0; i <= 10; i++ {
		idList = append(idList, uint64(i))
		posts = append(posts, m.Post{ID: uint64(i)})
		db.Create(&m.ImageAttachment{RecordID: uint64(i), RecordType: "post"})
	}

	as, err := GetRelatedAttachments(posts, &db)
	if err != nil {
		t.Error("Error getting related comment counts from posts:", err)
		return
	}

	if len(as) < 10 {
		t.Error("Not enough record likes:", len(as))
	}

	for _, a := range as {
		found := false
		for _, v := range idList {
			if a.RecordID == v {
				found = true
			}
		}
		if !found {
			t.Error("Couldn't find id:", a.RecordID, "in", idList)
		}
	}

	db.Delete(&as)
}
