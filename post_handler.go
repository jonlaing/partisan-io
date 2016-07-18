package main

import (
	"net/http"
	"partisan/db"

	"github.com/gin-gonic/gin"
	"partisan/models.v2/posts"
)

func PostShow(c *gin.Context) {
	db := db.GetDB(c)

	postID := c.Param("post_id")

	post, err := posts.GetByID(postID, "", db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.HTML(http.StatusOK, "post", gin.H{"data": gin.H{}, "post": post})
}
