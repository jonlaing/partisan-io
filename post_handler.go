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

	fbOpengraph := map[string]string{
		"title":     post.Body,
		"url":       "https://www.partisan.io/posts/" + post.ID,
		"site_name": "Partisan.IO",
	}

	if len(post.Attachments) > 0 {
		fbOpengraph["image"] = post.Attachments[0].URL
	}

	c.HTML(http.StatusOK, "post", gin.H{"title": post.Body, "og": fbOpengraph, "data": gin.H{}, "post": post})
}
