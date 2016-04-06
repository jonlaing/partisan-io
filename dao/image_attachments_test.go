package dao

import (
	m "partisan/models"
	"testing"
)

func TestGetMultipleRelatedAttachments(t *testing.T) {
	var posts m.Posts
	var idList []uint64
	for i := 0; i <= 10; i++ {
		idList = append(idList, uint64(i))
		posts = append(posts, m.Post{ID: uint64(i)})
		db.Create(&m.ImageAttachment{RecordID: uint64(i), RecordType: "post"})
	}

	as, err := GetMultipleRelatedAttachments(posts, db)
	defer db.Delete(&as)

	if err != nil {
		t.Error("Error getting related attachments from posts:", err)
		return
	}

	if len(as) < 10 {
		t.Error("Not enough attachments:", len(as))
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

}

func TestGetRelatedAttachments(t *testing.T) {
	post := m.Post{ID: 1}
	attachment := m.ImageAttachment{RecordID: 1, RecordType: "post"}
	db.Create(&attachment)
	defer db.Delete(&attachment)

	as, err := GetRelatedAttachments(&post, db)
	if err != nil {
		t.Error("Error getting related attachments from post:", err)
		return
	}

	if len(as) != 1 {
		t.Error("Not enough attachments:", len(as))
		return
	}

	if as[0].RecordID != uint64(1) {
		t.Error("Wrong attachment")
		return
	}
}
