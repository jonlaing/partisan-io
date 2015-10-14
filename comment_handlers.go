package main

import (
	"fmt"
	"net/http"
	"partisan/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"partisan/db"
	m "partisan/models"
)

// CommentShow redirects to the post of the comment and the inline anchor
func CommentShow(c *gin.Context) {
	db := db.GetDB(c)

	commentID := c.Param("record_id")

	var comment m.Comment
	if err := db.Find(&comment, commentID).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	route := fmt.Sprintf("/posts/%d#comment-%d", comment.PostID, comment.ID)
	c.Redirect(http.StatusMovedPermanently, route)
	return
}
