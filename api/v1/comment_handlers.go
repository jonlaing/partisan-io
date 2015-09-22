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
	Comment m.Comment `json:"comment"`
	User    m.User    `json:"user"`
	Count   int       `json:"count"`
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
	var userIDs []uint64
	var users []m.User

	resp := []CommentResp{}

	rID, rType, err := getRecord(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	if err := db.Where("record_type = ? AND record_id = ?", rType, rID).Find(&comments).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	if len(comments) == 0 {
		c.JSON(http.StatusOK, gin.H{"comments": resp})
		return
	}

	// get all User IDs
	for _, comment := range comments {
		userIDs = append(userIDs, comment.UserID)
	}

	// Batch find all users
	db.Where("id IN (?)", userIDs).Find(&users).Debug()

	// compile all users and comments
	// in this order because there will be more comments than users most likely
	for _, comment := range comments {
		for _, user := range users {
			if comment.UserID == user.ID {
				resp = append(resp, CommentResp{comment, user, 0})
			}
		}
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

	rID, rType, err := getRecord(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	var count int
	if err := db.Model(m.Comment{}).Where("record_type = ? AND record_id = ?", rType, rID).Count(&count).Error; err != nil {
		count = 0
	}

	c.JSON(http.StatusOK, gin.H{"record_type": rType, "record_id": rID, "comment_count": count})
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
        db.Model(m.Comment{}).Where("record_id = ? AND record_type = ?", comment.RecordID, comment.RecordType).Count(&count)

	commentResp := CommentResp{
		Comment: comment,
		User:    user,
                Count: count,
	}

	// Create feed item
	feedItem := m.FeedItem{
		Action:     "comment",
		UserID:     user.ID,
		RecordType: comment.RecordType,
		RecordID:   comment.RecordID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	db.Create(&feedItem) // Don't need to error check

	c.JSON(http.StatusCreated, commentResp)
}
