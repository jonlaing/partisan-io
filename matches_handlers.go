package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"partisan/auth"
	"partisan/db"
)

// MatchesIndex shows the matches screen
func MatchesIndex(c *gin.Context) {
	db, err := db.InitDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	currentUser, err := auth.CurrentUser(c, &db)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	c.HTML(http.StatusOK, "profile_edit", gin.H{
		"title": "Matches",
		"data": gin.H{
			"user":  currentUser,
		},
	})
}
