package main

import (
	// "encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Profile of a user
type Profile struct {
	ID            uint64 `form:"id" json:"id" gorm:"primary_key"`
	UserID        uint64 `json:"user_id"`
	LookingFor    int    `form:"looking_for" json:"looking_for"`
	Gender        string `form: "gender" json:"gender"`
	Summary       string `form:"summary" json:"summary"`
	CoverPhotoURL string `json:"cover_photo_url"`
}

// ProfileResp is the JSON response for a show
type ProfileResp struct {
	Profile Profile `json:"profile"`
	User    User    `json:"user"`
	Match   float64 `json:"match"`
	Enemy   float64 `json:"enemy"`
}

const (
	// ForFriends is value for LookingFor bitfield
	ForFriends int = 1 << iota

	// ForLove is value for LookingFor bitfield
	ForLove

	// ForEnemies is value for LookingFor bitfield
	ForEnemies
)

// ProfileShow is API for showing profile info
func ProfileShow(c *gin.Context) {
	db, err := initDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	var user, currentUser User

	userID := c.Param("user_id")
	// no param in route
	if len(userID) == 0 {
		user, err = CurrentUser(c, &db)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}
	} else {
		if err := db.First(&user, userID).Error; err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}

		currentUser, err = CurrentUser(c, &db)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}
	}

	profile := Profile{}
	if err := db.Where("user_id = ?", user.ID).First(&profile).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	resp := ProfileResp{
		Profile: profile,
		User:    user,
	}

	// If this is not the current user, do the match
	if user.ID != currentUser.ID {
		match, _ := Match(user.PoliticalMap, currentUser.PoliticalMap)
		enemy, _ := Enemy(user.PoliticalMap, currentUser.PoliticalMap)

		resp.Match = match
		resp.Enemy = enemy
	}

	c.JSON(http.StatusOK, resp)
}

// ProfileHTMLShow renders HTML
func ProfileHTMLShow(c *gin.Context) {
	db, err := initDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	currentUser, err := CurrentUser(c, &db)
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

	user := User{}
	if err := db.First(&user, userID).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	profile := Profile{}
	if err := db.Where("user_id = ?", user.ID).First(&profile).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	match, _ := Match(user.PoliticalMap, currentUser.PoliticalMap)
	enemy, _ := Enemy(user.PoliticalMap, currentUser.PoliticalMap)

	fmt.Println("dafuq")
	c.HTML(http.StatusOK, "profile_show.html",
		gin.H{
			"profile": profile,
			"user":    user,
			"match":   fmt.Sprintf("%.f", match*100),
			"enemy":   fmt.Sprintf("%.f", enemy*100),
		})
}
