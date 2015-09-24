package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"partisan/auth"
	"partisan/db"
	m "partisan/models"
)

// PostComments stores like data for ease
type PostComments struct {
	RecordID uint64
	Count    int
}

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

	// TODO: limit feed so a particular record only comes up once
	feedItems := []m.FeedItem{}
	if err := db.Where("user_id IN (?) AND action = ?", friendIDs, "post").Order("created_at desc").Limit(50).Find(&feedItems).Error; err != nil {
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

	postLikes, err := m.GetLikes(user.ID, "post", postIDs, &db)
	if err != nil {
		fmt.Println(err)
	}

	postComments, err := getPostComments(postIDs, &db)
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < len(feedItems); i++ {
		collectPosts(&feedItems[i], posts, users, attachments, postLikes, postComments)
	}

	c.JSON(http.StatusOK, gin.H{"feed_items": feedItems})
}

func collectPosts(f *m.FeedItem, posts []m.Post, users []m.User, attachments []m.ImageAttachment, likes []m.RecordLikes, comments []PostComments) {
	for _, post := range posts {
		if f.RecordID == post.ID {
			user, _ := findMatchingPostUser(post, users)
			attachment, _ := m.GetAttachment(post.ID, attachments)
			likeCount, liked, _ := fineMatchingPostLikes(post, likes)
			commentCount, _ := findMatchingCommentCount(post, comments)

			f.Record = PostResponse{
				Post:         post,
				User:         user,
				Attachment:   attachment,
				LikeCount:    likeCount,
				Liked:        liked,
				CommentCount: commentCount,
			}
		}
	}
}
