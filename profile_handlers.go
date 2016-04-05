package main

import (
	"fmt"
	"net/http"
	"partisan/auth"
	"partisan/db"
	"partisan/matcher"
	m "partisan/models"

	"github.com/gin-gonic/gin"
)

// ProfileShow renders HTML
func ProfileShow(c *gin.Context) {
	db := db.GetDB(c)

	currentUser, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	username := c.Param("username")
	// no param in route
	if len(username) == 0 {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("Couldn't find Username"))
		return
	}

	user := m.User{}
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if user.ID == currentUser.ID {
		c.Redirect(http.StatusFound, "/feed")
		return
	}

	profile := m.Profile{}
	if err := db.Where("user_id = ?", user.ID).First(&profile).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	match, _ := matcher.Match(user.PoliticalMap, currentUser.PoliticalMap)

	c.HTML(http.StatusOK, "profile_show",
		gin.H{
			"title": "@" + user.Username + "'s Profile",
			"data": gin.H{
				"profile":      profile,
				"user":         user,
				"current_user": currentUser,
				"match":        fmt.Sprintf("%.f", match*100),
			},
		})
}

// ProfileEdit shows edit form for current user
func ProfileEdit(c *gin.Context) {
	db := db.GetDB(c)

	currentUser, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	profile := m.Profile{}
	if err := db.Where("user_id = ?", currentUser.ID).First(&profile).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.HTML(http.StatusOK, "profile_edit", gin.H{
		"title": "Edit My Profile",
		"data": gin.H{
			"user":    currentUser,
			"profile": profile,
		},
	})
}
