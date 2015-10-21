package dao

import (
	"testing"

	m "partisan/models"
)

func TestGetRelatedUsers(t *testing.T) {
	var posts []m.Post
	var idList []uint64
	for i := 0; i <= 10; i++ {
		idList = append(idList, uint64(i))
		posts = append(posts, m.Post{UserID: uint64(i)})
		db.Create(&m.User{ID: uint64(i)})
	}

	users, err := GetRelatedUsers(m.Posts(posts), &db)
	if err != nil {
		t.Error("Error getting related users from posts:", err)
		return
	}

	if len(users) < 10 {
		t.Error("Not enough users:", len(users))
	}

	for _, u := range users {
		found := false
		for _, v := range idList {
			if u.ID == v {
				found = true
			}
		}
		if !found {
			t.Error("Couldn't find id:", u.ID, "in", idList)
		}
	}

	db.Delete(&users)
}

func TestGetMatchingUser(t *testing.T) {
	post := m.Post{UserID: 1}
	var users []m.User
	for i := 0; i < 10; i++ {
		users = append(users, m.User{ID: uint64(i)})
	}

	u, ok := GetMatchingUser(&post, users)
	if ok == false {
		t.Fail()
	}

	if u.ID != post.UserID {
		t.Error("IDs don't match:", u.ID, ",", post.UserID)
	}
}
