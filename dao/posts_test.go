package dao

import (
	m "partisan/models"
	"testing"
)

func TestGetRelatedPosts(t *testing.T) {
	var feedItems m.FeedItems
	var idList []uint64
	for i := 0; i <= 10; i++ {
		idList = append(idList, uint64(i))
		feedItems = append(feedItems, m.FeedItem{RecordType: "post", RecordID: uint64(i)})
		db.Create(&m.Post{ID: uint64(i)})
	}

	posts, err := GetRelatedPosts(feedItems, db)
	if err != nil {
		t.Error("Error getting related posts from feed items:", err)
		return
	}

	if len(posts) < 10 {
		t.Error("Not enough posts:", len(posts))
	}

	for _, p := range posts {
		found := false
		for _, v := range idList {
			if p.ID == v {
				found = true
			}
		}
		if !found {
			t.Error("Couldn't find id:", p.ID, "in", idList)
		}
	}

	db.Delete(&posts)
}

func TestFindMatchingPostLikes(t *testing.T) {
	post := m.Post{ID: 5}
	var rls []RecordLikes
	for i := 0; i < 10; i++ {
		rls = append(rls, RecordLikes{RecordID: uint64(i), Count: i})
	}

	count, _, ok := findMatchingPostLikes(post, rls)
	if !ok {
		t.Fail()
	}

	if count != 5 {
		t.Error("Wrong number of likes. Should be 5, was:", count)
	}
}

func TestFindMatchingPostCommentCount(t *testing.T) {
	post := m.Post{ID: 5}
	var pcs []PostComments
	for i := 0; i < 10; i++ {
		pcs = append(pcs, PostComments{RecordID: uint64(i), Count: i})
	}

	count, ok := findMatchingPostCommentCount(post, pcs)
	if !ok {
		t.Fail()
	}

	if count != 5 {
		t.Error("Wrong number of comments. Should be 5, was:", count)
	}
}

func TestFindRelatedPostAttachment(t *testing.T) {
	post := m.Post{ID: 5}
	var as []m.ImageAttachment
	for i := 0; i <= 10; i++ {
		as = append(as, m.ImageAttachment{RecordType: "post", RecordID: uint64(i)})
	}

	a, ok := findRelatedPostAttachment(post, as)
	if !ok {
		t.Fail()
	}

	if a.RecordID != post.ID {
		t.Fail()
	}
}

func TestCollectPostResponses(t *testing.T) {
	var posts []m.Post
	var users []m.User
	var attachments []m.ImageAttachment
	var likes []RecordLikes
	var comments []PostComments

	for i := 0; i < 10; i++ {
		posts = append(posts, m.Post{ID: uint64(i), UserID: uint64(i)})
		users = append(users, m.User{ID: uint64(i)})
		attachments = append(attachments, m.ImageAttachment{RecordID: uint64(i), RecordType: "post"})
		likes = append(likes, RecordLikes{RecordID: uint64(i)})
		comments = append(comments, PostComments{RecordID: uint64(i)})
	}

	prs := collectPostResponses(posts, users, attachments, likes, comments)

	for _, p := range posts {
		if _, ok := prs[p.ID]; !ok {
			t.Error("Couldn't match post:", p.ID)
		}
	}

}
