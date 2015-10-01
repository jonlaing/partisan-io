package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"partisan/auth"
)

// FeedIndex renders HTML
func FeedIndex(c *gin.Context) {
	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	c.HTML(http.StatusOK, "feed", gin.H{
		"title": "My Feed",
		"data": gin.H{
			"user": user,
		},
	})
}
