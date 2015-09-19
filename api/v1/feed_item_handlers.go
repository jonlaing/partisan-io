package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"partisan/auth"
	"partisan/db"
	m "partisan/models"
)

// FeedIndex shows all Feed Items for a particular user
func FeedIndex(c *gin.Context) {
	db, err := db.InitDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	db.LogMode(true)

	user, _ := auth.CurrentUser(c, &db)

	friendIDs, err := ConfirmedFriendIDs(user, c, &db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	friendIDs = append(friendIDs, user.ID)

	feedItems := []m.FeedItem{}
	if err := db.Where("user_id IN (?)", friendIDs).Order("created_at desc").Limit(50).Find(&feedItems).Error; err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var postIDs, postUserIDs []uint64
	for _, item := range feedItems {
		if item.RecordType == "post" {
			postIDs = append(postIDs, item.RecordID)
		}
	}

	var posts []m.Post
	if err := db.Where("id IN (?)", postIDs).Find(&posts).Error; err == nil {
		for _, post := range posts {
			postUserIDs = append(postUserIDs, post.UserID)
		}
	}

	var users []m.User
	db.Where("id IN (?)", postUserIDs).Find(&users)

	var attachments []m.ImageAttachment
	db.Where("record_type = ? AND record_id IN (?)", "post", postIDs).Find(&attachments)

	for i := 0; i < len(feedItems); i++ {
		collectPosts(&feedItems[i], posts, users, attachments)
	}

	c.JSON(http.StatusOK, gin.H{"feed_items": feedItems})
}

func collectPosts(f *m.FeedItem, posts []m.Post, users []m.User, attachments []m.ImageAttachment) {
	for _, post := range posts {
		if f.RecordID == post.ID {
			user, _ := findMatchingUser(post, users)
			attachment, _ := findMatchingAttachment(post, attachments)

			f.Record = PostResponse{
				Post:       post,
				User:       user,
				Attachment: attachment,
			}
		}
	}
}
func findMatchingUser(post m.Post, users []m.User) (m.User, bool) {
	for _, user := range users {
		if user.ID == post.UserID {
			return user, true
		}
	}
	return m.User{}, false
}

func findMatchingAttachment(post m.Post, attachments []m.ImageAttachment) (m.ImageAttachment, bool) {
	for _, attachment := range attachments {
		if attachment.RecordID == post.ID {
			return attachment, true
		}
	}
	return m.ImageAttachment{}, false
}
