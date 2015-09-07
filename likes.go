package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Like is polymorphic
type Like struct {
	ID         uint64 `gorm:"primary_key"`
	UserID     uint64
	RecordID   uint64
	RecordType string
}

// Dislike is polymorphic
type Dislike struct {
	ID         uint64 `gorm:"primary_key"`
	UserID     uint64
	RecordID   uint64
	RecordType string
}

// LikeCreate creates a Like for a particular record
func LikeCreate(c *gin.Context) {
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

	var like Like
	user, _ := CurrentUser(c, &db)

	postID, ok := c.Get("post_id")
	if ok {
		like = Like{
			UserID:     user.ID,
			RecordID:   postID.(uint64),
			RecordType: "post",
		}

		if err := db.Create(&like).Error; err != nil {
			c.AbortWithError(http.StatusNotAcceptable, err)
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "liked"})
		return
	}
}

// DislikeCreate creates a Like for a particular record
func DislikeCreate(c *gin.Context) {
	db, err := initDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	var dislike Dislike
	user, _ := CurrentUser(c, &db)

	postID, ok := c.Get("post_id")
	if ok {
		dislike = Dislike{
			UserID:     user.ID,
			RecordID:   postID.(uint64),
			RecordType: "post",
		}

		if err := db.Create(&dislike).Error; err != nil {
			c.AbortWithError(http.StatusNotAcceptable, err)
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "disliked"})
		return
	}
}
