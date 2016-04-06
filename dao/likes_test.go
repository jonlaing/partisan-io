package dao

import (
	m "partisan/models"
	"testing"
)

func TestGetMultipleRelatedLikes(t *testing.T) {
	var posts m.Posts
	var likes []m.Like
	var idList []uint64
	for i := 0; i <= 10; i++ {
		idList = append(idList, uint64(i))
		posts = append(posts, m.Post{ID: uint64(i)})
		for j := 0; j < 5; j++ {
			like := m.Like{RecordID: uint64(i), RecordType: "post"}
			db.Create(&like)
			likes = append(likes, like)
		}
	}
	defer db.Delete(&likes)

	rls, err := GetMultipleRelatedLikes(uint64(1), posts, db)
	if err != nil {
		t.Error("Error getting related like counts from posts:", err)
		return
	}

	if len(rls) < 10 {
		t.Error("Not enough record likes:", len(rls))
	}

	for _, rl := range rls {
		found := false
		for _, v := range idList {
			if rl.RecordID == v && rl.Count == 5 {
				found = true
			}
		}
		if !found {
			t.Error("Couldn't find id:", rl.RecordID, "in", idList, "with count of 5:", rl.Count)
		}
	}
}

func TestGetRelatedLikes(t *testing.T) {
	post := m.Post{ID: 1}
	var likes []m.Like

	for i := 0; i <= 10; i++ {
		like := m.Like{UserID: uint64(i), RecordID: 1, RecordType: "post"}
		db.Create(&like)
		likes = append(likes, like)
	}
	defer db.Delete(&likes)

	c, l, err := GetRelatedLikes(1, &post, db)
	if err != nil {
		t.Error("Error getting related likes for post:", err)
		return
	}

	if c < 10 {
		t.Error("Expected 10 likes, got:", c)
		return
	}

	if !l {
		t.Error("Expected to have liked post")
		return
	}
}
