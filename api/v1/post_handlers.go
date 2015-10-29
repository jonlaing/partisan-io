package v1

import (
	"fmt"
	"net/http"
	"partisan/auth"
	"partisan/db"
	m "partisan/models"
	"time"

	"partisan/Godeps/_workspace/src/github.com/gin-gonic/gin"
)

// PostResponse is the response schema
type PostResponse struct {
	Post         m.Post            `json:"post"`
	Attachment   m.ImageAttachment `json:"image_attachment"`
	User         m.User            `json:"user"`
	LikeCount    int               `json:"like_count"`
	Liked        bool              `json:"liked"`
	CommentCount int               `json:"comment_count"`
}

// PostsIndex display all posts
func PostsIndex(c *gin.Context) {
	db := db.GetDB(c)

	userID, ok := c.Get("user_id")
	if !ok {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("User ID not set"))
		return
	}

	posts := []m.Post{}
	if err := db.Where("user_id = ?", userID).Find(&posts).Error; err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, posts)
}

// PostsCreate create a post
func PostsCreate(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	postBody := c.Request.FormValue("body")

	post := m.Post{
		UserID:    user.ID,
		Body:      postBody,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := db.Create(&post).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	m.FindAndCreateHashtags(&post, db)
	m.FindAndCreateUserTags(&post, db)

	postRes := PostResponse{
		Post: post,
		User: user,
	}

	// Doing it this way because we don't know if a user will try
	// to attach an image. This way we can fail elegantly
	if err := m.AttachImage(c, &postRes); err != nil {
		// only errs with catostrophic failure,
		// silently fails if no attachment is present
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Create feed item
	feedItem := m.FeedItem{
		UserID:     user.ID,
		Action:     "post",
		RecordType: "post",
		RecordID:   post.ID,
		Record:     postRes,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := db.Create(&feedItem).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	c.JSON(http.StatusOK, feedItem)
}

// // PostsShow show a post
// func PostsShow(c *gin.Context) {
// 	db, err := db.InitDB()
// 	if err != nil {
// 		c.AbortWithError(http.StatusInternalServerError, err)
// 		return
// 	}
// 	defer db.Close()

// 	post := m.Post{
// 		ID:        123,
// 		Body:      "this is how we do it! (uhuh)",
// 		CreatedAt: time.Now(),
// 		UpdatedAt: time.Now(),
// 	}
// 	user := m.User{
// 		Username: "Franny_Frumpernickle",
// 	}
// 	// id := c.Params.ByName("id")

// 	// if err := db.First(&post, id).Related(&user).Error; err != nil {
// 	// 	c.AbortWithError(http.StatusNotFound, err)
// 	// 	return
// 	// }

// 	resp := PostResponse{
// 		Post: post,
// 		User: user,
// 	}

// 	c.JSON(http.StatusOK, resp)
// }

// PostsUpdate update a post
func PostsUpdate(c *gin.Context) {
	// type conversion from getting user from context can cause panic
	defer func() {
		if r := recover(); r != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}()

	db := db.GetDB(c)

	userID, ok := c.Get("user_id")
	if !ok {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("User ID not set"))
		return
	}

	post := m.Post{}
	id := c.Params.ByName("id")

	if err := db.First(&post, id).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if post.UserID != userID.(uint64) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if err := c.BindJSON(&post); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := db.Save(&post).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	c.JSON(http.StatusOK, post)
}

// PostsDestroy update a post
func PostsDestroy(c *gin.Context) {
	// type conversion from getting user from context can cause panic
	defer func() {
		if r := recover(); r != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}()

	db := db.GetDB(c)

	userID, ok := c.Get("user_id")
	if !ok {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("User ID not set"))
		return
	}

	post := m.Post{}
	id := c.Params.ByName("id")

	if err := db.First(&post, id).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if post.UserID != userID.(uint64) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if err := db.Delete(&post).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// GetID satisfies m.ImageAttacher interface
func (pr *PostResponse) GetID() uint64 {
	return pr.Post.ID
}

// GetUserID satisfies m.ImageAttacher interface
func (pr *PostResponse) GetUserID() uint64 {
	return pr.User.ID
}

// GetType satisfies m.ImageAttacher interface
func (pr *PostResponse) GetType() string {
	return "post"
}

// AttachImage satisfies m.ImageAttacher interface
func (pr *PostResponse) AttachImage(i m.ImageAttachment) error {
	pr.Attachment = i
	return nil
}
