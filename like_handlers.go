package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"partisan/db"
	m "partisan/models"
)

// LikeShow redirects to the post of the comment and the inline anchor
func LikeShow(c *gin.Context) {
	db, err := db.InitDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	likeID := c.Param("record_id")

	var like m.Like
	if err := db.Find(&like, likeID).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if like.RecordType == "post" || like.RecordType == "posts" {
		route := fmt.Sprintf("/posts/%d", like.RecordID)
		c.Redirect(http.StatusMovedPermanently, route)
		return
	}

	c.AbortWithStatus(http.StatusNotAcceptable)
	return
}
