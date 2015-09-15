package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"partisan/db"
	m "partisan/models"
	"time"
	"partisan/auth"
)

// PostResponse is the response schema
type PostResponse struct {
	Post m.Post `json:"post"`
	User m.User `json:"user"`
}

// PostsIndex display all posts
func PostsIndex(c *gin.Context) {
	db, err := db.InitDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

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
	// type conversion from getting user from context can cause panic
	defer func() {
		if r := recover(); r != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}()

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

	post := m.Post{
		UserID:    user.ID,
		Body:      c.Request.PostFormValue("body"),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := db.Create(&post).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	postRes := PostResponse{
		Post: post,
		User: user,
	}

	// Create feed item
	feedItem := m.FeedItem{
          UserID: user.ID,
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

	db, err := db.InitDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

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

	db, err := db.InitDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

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
