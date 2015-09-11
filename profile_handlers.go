package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"partisan/auth"
	"partisan/db"
	"partisan/matcher"
	m "partisan/models"
)

// ProfileShow renders HTML
func ProfileShow(c *gin.Context) {
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

	userID := c.Param("user_id")
	// no param in route
	if len(userID) == 0 {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("Couldn't find User ID"))
		return
	}

	user := m.User{}
	if err := db.First(&user, userID).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	profile := m.Profile{}
	if err := db.Where("user_id = ?", user.ID).First(&profile).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	match, _ := matcher.Match(user.PoliticalMap, currentUser.PoliticalMap)
	enemy, _ := matcher.Enemy(user.PoliticalMap, currentUser.PoliticalMap)

	fmt.Println("dafuq")
	c.HTML(http.StatusOK, "profile_show.html",
		gin.H{
			"profile": profile,
			"user":    user,
			"match":   fmt.Sprintf("%.f", match*100),
			"enemy":   fmt.Sprintf("%.f", enemy*100),
		})
}
