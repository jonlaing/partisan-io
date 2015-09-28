package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"partisan/auth"
	"partisan/db"
	m "partisan/models"
	"strconv"
	"time"
)

// CommentResp is the format for JSON responses
type CommentResp struct {
	Comment    m.Comment         `json:"comment"`
	Attachment m.ImageAttachment `json:"image_attachment"`
	User       m.User            `json:"user"`
	LikeCount  int               `json:"like_count"`
	Liked      bool              `json:"liked"`
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
	var attachments []m.ImageAttachment

	resp := []CommentResp{}

	pID := c.Param("record_id")
	user, _ := auth.CurrentUser(c)

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
	db.Where("id IN (?)", userIDs).Find(&users)

	// Batch find all attachments
	db.Where("record_type = ? AND record_id IN (?)", "comment", commentIDs).Find(&attachments)

	// compile all users and comments
	// in this order because there will be more comments than users most likely
	for _, comment := range comments {
		cUser, _ := findMatchingCommentUser(comment, users)
		like, _ := findMatchingCommentLikes(comment, likes)
		attachment, _ := m.GetAttachment(comment.ID, attachments)

		resp = append(resp, CommentResp{
			Comment:    comment,
			User:       cUser,
			LikeCount:  like.Count,
			Liked:      like.UserCount == 1,
			Attachment: attachment,
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

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	postIDString := c.Request.FormValue("post_id")
	postID, err := strconv.ParseUint(postIDString, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	comment := m.Comment{
		UserID:    user.ID,
		PostID:    postID,
		Body:      c.Request.FormValue("body"),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Create(&comment).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	m.FindAndCreateHashtags(&comment, &db)
	m.FindAndCreateUserTags(&comment, &db)
	m.NewNotification(&comment, user.ID, &db)

	var count int
	db.Model(m.Comment{}).Where("post_id = ?", comment.PostID).Count(&count)

	commentResp := CommentResp{
		Comment:   comment,
		User:      user,
		LikeCount: 0,
		Liked:     false,
	}

	// Doing it this way because we don't know if a user will try
	// to attach an image. This way we can fail elegantly
	if err = m.AttachImage(c, &db, &commentResp); err != nil {
		// only errs with catostrophic failure,
		// silently fails if no attachment is present
		c.AbortWithError(http.StatusInternalServerError, err)
		return
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

// GetID satisfies m.ImageAttacher interface
func (cr *CommentResp) GetID() uint64 {
	return cr.Comment.ID
}

// GetUserID satisfies m.ImageAttacher interface
func (cr *CommentResp) GetUserID() uint64 {
	return cr.User.ID
}

// GetType satisfies m.ImageAttacher interface
func (cr *CommentResp) GetType() string {
	return "comment"
}

// AttachImage satisfies m.ImageAttacher interface
func (cr *CommentResp) AttachImage(i m.ImageAttachment) error {
	cr.Attachment = i
	return nil
}
