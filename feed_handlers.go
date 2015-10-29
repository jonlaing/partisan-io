package main

import (
	"net/http"
	"partisan/auth"
	"partisan/db"
	m "partisan/models"

	"partisan/Godeps/_workspace/src/github.com/gin-gonic/gin"
)

// FeedIndex renders HTML
func FeedIndex(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.Redirect(http.StatusUnauthorized, "/login")
		return
	}

	profile := m.Profile{}
	if err := db.Where("user_id = ?", user.ID).First(&profile).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.HTML(http.StatusOK, "feed", gin.H{
		"title": "My Feed",
		"data": gin.H{
			"user":    user,
			"profile": profile,
		},
	})
}
