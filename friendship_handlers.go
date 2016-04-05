package main

import (
	"net/http"
	"partisan/auth"

	"github.com/gin-gonic/gin"
)

func FriendsIndex(c *gin.Context) {
	currentUser, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	c.HTML(http.StatusOK, "friends", gin.H{
		"title": "Friends",
		"data": gin.H{
			"user": currentUser,
		},
	})
}
