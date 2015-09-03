package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// Post is the primary user created content. It can be just about anything.
type Post struct {
	ID        uint64    `json:"id" gorm:"primary_key"`      // Primary key
	UserID    uint64    `json:"user_id" binding:"required"` // ID of user that created Post
	Body      string    `json:"body" binding:"required"`    // The text based body
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	FeedItem  FeedItem  `json:"-" gorm:"polymorphic:FeedItem"`
}

// PostResponse is the response schema
type PostResponse struct {
	Post         Post `json:"post"`
	User         User `json:"user"`
	LikeCount    int  `json:"like_count"`
	DislikeCount int  `json:"dislike_count"`
	CommentCount int  `json:"comment_count"`
}

// PostsIndex display all posts
func PostsIndex(c *gin.Context) {
	db, err := initDB()
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

	posts := []Post{}
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

	db, err := initDB()
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

	post := Post{}

	if err := c.BindJSON(&post); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	post.UserID = userID.(uint64)

	if err := db.Create(&post).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	c.JSON(http.StatusCreated, post)
}

// PostsShow show a post
func PostsShow(c *gin.Context) {
	db, err := initDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	post := Post{
          ID: 123,
		Body:      "this is how we do it! (uhuh)",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	user := User{
          Username: "Franny_Frumpernickle",
        }
	// id := c.Params.ByName("id")

	// if err := db.First(&post, id).Related(&user).Error; err != nil {
	// 	c.AbortWithError(http.StatusNotFound, err)
	// 	return
	// }

	resp := PostResponse{
		Post: post,
		User: user,
	}

	c.JSON(http.StatusOK, resp)
}

// PostsUpdate update a post
func PostsUpdate(c *gin.Context) {
	// type conversion from getting user from context can cause panic
	defer func() {
		if r := recover(); r != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}()

	db, err := initDB()
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

	post := Post{}
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

	db, err := initDB()
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

	post := Post{}
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
