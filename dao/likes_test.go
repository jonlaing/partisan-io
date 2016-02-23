package dao

import (
	m "partisan/models"
	"testing"
)

func TestGetRelatedLikes(t *testing.T) {
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

	rls, err := GetRelatedLikes(uint64(1), posts, &db)
	if err != nil {
		t.Error("Error getting related comment counts from posts:", err)
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

	db.Delete(&likes)
}
