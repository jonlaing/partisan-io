package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"partisan/auth"
	"partisan/db"
	m "partisan/models"
	"time"
)

// CommentResp is the format for JSON responses
type CommentResp struct {
	Comment   m.Comment `json:"comment"`
	User      m.User    `json:"user"`
	LikeCount int       `json:"like_count"`
	Liked     bool      `json:"liked"`
}

// CommentsIndex shows a list of comments based on a record_id
func CommentsIndex(c *gin.Context) {
	db, err := db.InitDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	var comments []m.Comment
	var commentIDs []uint64
	var userIDs []uint64
	var users []m.User

	resp := []CommentResp{}

	pID := c.Param("record_id")
	user, _ := auth.CurrentUser(c, &db)

	if err := db.Where("post_id = ?", pID).Find(&comments).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	if len(comments) == 0 {
		c.JSON(http.StatusOK, gin.H{"comments": resp})
		return
	}

	// get all User IDs
	for _, comment := range comments {
		commentIDs = append(commentIDs, comment.ID)
		userIDs = append(userIDs, comment.UserID)
	}

	likes, err := m.GetLikes(user.ID, "comment", commentIDs, &db)

	// Batch find all users
	db.Where("id IN (?)", userIDs).Find(&users).Debug()

	// compile all users and comments
	// in this order because there will be more comments than users most likely
	for _, comment := range comments {
		cUser, _ := findMatchingCommentUser(comment, users)
		like, _ := findMatchingCommentLikes(comment, likes)

		resp = append(resp, CommentResp{
			Comment:   comment,
			User:      cUser,
			LikeCount: like.Count,
			Liked:     like.UserCount == 1,
		})
	}

	c.JSON(http.StatusOK, gin.H{"comments": resp})
	return
}

// CommentsCount returns a count of all comments for a particular record
func CommentsCount(c *gin.Context) {
	db, err := db.InitDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	pID := c.Param("record_id")

	var count int
	if err := db.Model(m.Comment{}).Where("post_id = ?", pID).Count(&count).Error; err != nil {
		count = 0
	}

	c.JSON(http.StatusOK, gin.H{"post_id": pID, "comment_count": count})
	return
}

// CommentsCreate creates a comment
func CommentsCreate(c *gin.Context) {
	db, err := db.InitDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	user, err := auth.CurrentUser(c, &db)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	var comment m.Comment
	if err := c.BindJSON(&comment); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	comment.UserID = user.ID
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()

	if err := db.Create(&comment).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	m.NewNotification(&comment, user.ID, &db)

	var count int
	db.Model(m.Comment{}).Where("post_id = ?", comment.PostID).Count(&count)

	commentResp := CommentResp{
		Comment:   comment,
		User:      user,
		LikeCount: 0,
		Liked:     false,
	}

	// Create feed item
	feedItem := m.FeedItem{
		Action:     "comment",
		UserID:     user.ID,
		RecordType: "post",
		RecordID:   comment.PostID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	db.Create(&feedItem) // Don't need to error check

	c.JSON(http.StatusCreated, commentResp)
}

func findMatchingCommentUser(comment m.Comment, users []m.User) (m.User, bool) {
	for _, user := range users {
		if comment.UserID == user.ID {
			return user, true
		}
	}

	return m.User{}, false
}

func findMatchingCommentLikes(comment m.Comment, likes []m.RecordLikes) (m.RecordLikes, bool) {
	for _, like := range likes {
		if like.RecordID == comment.ID {
			return like, true
		}
	}
	return m.RecordLikes{}, false
}
