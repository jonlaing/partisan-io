package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// FeedItem is the record of a user interacting with a Post, this is used to build a feed
type FeedItem struct {
	ID         uint64      `json:"id" gorm:"primary_key"`      // Primary key
	UserID     uint64      `json:"user_id" binding:"required"` // ID of user that created Post
	Action     string      `json:"action" binding:"required"`
	RecordType string      `json:"record_type" binding:"required"`
	RecordID   uint64      `json:"record_id" binding:"required"`
	Record     interface{} `json:"record,omitempty" sql:"-"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

// FeedIndex shows all Feed Items for a particular user
func FeedIndex(c *gin.Context) {
	db, err := initDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	feedItems := []FeedItem{}
	if err := db.Find(&feedItems).Error; err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// TODO: This is an N+1 problem...
	for k, item := range feedItems {
		if item.RecordType == "post" {
			post := Post{}
			user := User{}
			db.First(&post, item.RecordID).Related(&user) // right here

			feedItems[k].Record = PostResponse{
				Post: post,
				User: user,
			}
		}
	}

	c.JSON(http.StatusOK, feedItems)
}
