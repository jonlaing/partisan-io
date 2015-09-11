package v1

import (
	// "encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"partisan/db"
	"partisan/matcher"
	m "partisan/models"
	"partisan/auth"
)

// ProfileResp is the JSON response for a show
type ProfileResp struct {
	Profile m.Profile `json:"profile"`
	User    m.User    `json:"user"`
	Match   float64   `json:"match"`
	Enemy   float64   `json:"enemy"`
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
	db, err := db.InitDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	var user, currentUser m.User

	userID := c.Param("user_id")
	// no param in route
	if len(userID) == 0 {
		user, err = auth.CurrentUser(c, &db)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}
	} else {
		if err := db.First(&user, userID).Error; err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}

		currentUser, err = auth.CurrentUser(c, &db)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}
	}

	profile := m.Profile{}
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
		match, _ := matcher.Match(user.PoliticalMap, currentUser.PoliticalMap)
		enemy, _ := matcher.Enemy(user.PoliticalMap, currentUser.PoliticalMap)

		resp.Match = match
		resp.Enemy = enemy
	}

	c.JSON(http.StatusOK, resp)
}
