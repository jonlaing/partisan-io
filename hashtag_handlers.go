package main

import (
	"net/http"
	"partisan/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"partisan/auth"
)

// HashtagShow renders the HTML for the hashtag search
func HashtagShow(c *gin.Context) {
	user, _ := auth.CurrentUser(c)

	search := c.Query("q")

	c.HTML(http.StatusOK, "hashtags", gin.H{
		"title": "Search Hashtags:" + search,
		"data": gin.H{
			"user":   user,
			"search": search,
		},
	})
}
