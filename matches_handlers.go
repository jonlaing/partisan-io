package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"partisan/auth"
)

// MatchesIndex shows the matches screen
func MatchesIndex(c *gin.Context) {
	currentUser, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	c.HTML(http.StatusOK, "profile_edit", gin.H{
		"title": "Matches",
		"data": gin.H{
			"user": currentUser,
		},
	})
}
