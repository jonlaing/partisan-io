package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"partisan/auth"
	"partisan/db"
	m "partisan/models"
)

// PostLikes stores like data for ease
type PostLikes struct {
	RecordID  uint64
	Count     int
	UserCount int
}

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

	postLikes, err := getLikes(user.ID, postIDs, &db)
	if err != nil {
		fmt.Println(err)
	}

	postComments, err := getComments(postIDs, &db)
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < len(feedItems); i++ {
		collectPosts(&feedItems[i], posts, users, attachments, postLikes, postComments)
	}

	c.JSON(http.StatusOK, gin.H{"feed_items": feedItems})
}

func collectPosts(f *m.FeedItem, posts []m.Post, users []m.User, attachments []m.ImageAttachment, likes []PostLikes, comments []PostComments) {
	for _, post := range posts {
		if f.RecordID == post.ID {
			user, _ := findMatchingUser(post, users)
			attachment, _ := findMatchingAttachment(post, attachments)
			likeCount, liked, _ := findMatchingLikes(post, likes)
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

func findMatchingLikes(post m.Post, likes []PostLikes) (int, bool, bool) {
	for _, like := range likes {
		if like.RecordID == post.ID {
			return like.Count, like.UserCount == 1, true
		}
	}
	return 0, false, false
}

func findMatchingCommentCount(post m.Post, comments []PostComments) (int, bool) {
	for _, comment := range comments {
		if comment.RecordID == post.ID {
			return comment.Count, true
		}
	}
	return 0, false
}

func getLikes(uID uint64, postIDs []uint64, db *gorm.DB) ([]PostLikes, error) {
	var likes []PostLikes

	rows, err := db.Raw("SELECT count(*), sum(case when user_id = ? then 1 else 0 end), record_id FROM \"likes\"  WHERE (record_type = 'posts' AND record_id IN (?)) GROUP BY record_id", uID, postIDs).Rows()
	defer rows.Close()
	if err != nil {
		return []PostLikes{}, err
	}

	for rows.Next() {
		var count, userCount int
		var rID uint64
		rows.Scan(&count, &userCount, &rID)
		likes = append(likes, PostLikes{Count: count, UserCount: userCount, RecordID: rID})
	}

	return likes, nil
}

func getComments(postIDs []uint64, db *gorm.DB) ([]PostComments, error) {
	var comments []PostComments

	rows, err := db.Raw("SELECT count(*), record_id FROM \"comments\"  WHERE (record_type = 'posts' AND record_id IN (?)) group by record_id", postIDs).Rows()
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
