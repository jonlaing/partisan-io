package main

import (
	"net/http"
	"partisan/auth"

	"github.com/gin-gonic/gin"
)

// MatchesIndex shows the matches screen
func MatchesIndex(c *gin.Context) {
	currentUser, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	c.HTML(http.StatusOK, "matches", gin.H{
		"title": "Matches",
		"data": gin.H{
			"user": currentUser,
		},
	})
}
