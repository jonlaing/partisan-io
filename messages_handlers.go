package main

import (
	"net/http"
	"partisan/auth"

	"github.com/gin-gonic/gin"
)

// MessagesIndex shows the messages screen
func MessagesIndex(c *gin.Context) {
	currentUser, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	c.HTML(http.StatusOK, "messages", gin.H{
		"title": "Messages",
		"data": gin.H{
			"user": currentUser,
		},
	})
}
