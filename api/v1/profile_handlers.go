package v1

import (
	"net/http"
	"partisan/auth"
	"partisan/db"
	"partisan/matcher"
	m "partisan/models"
	"strconv"

	"partisan/Godeps/_workspace/src/github.com/gin-gonic/gin"
)

// ProfileResp is the JSON response for a show
type ProfileResp struct {
	Profile m.Profile `json:"profile"`
	User    m.User    `json:"user"`
	Match   float64   `json:"match"`
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
	db := db.GetDB(c)

	var user, currentUser m.User
	var err error

	userID := c.Param("user_id")
	// no param in route
	if len(userID) == 0 {
		user, err = auth.CurrentUser(c)
		if err != nil {
			handleError(err, c)
			return
		}
	} else {
		if err := db.First(&user, userID).Error; err != nil {
			handleError(&ErrDBNotFound{err}, c)
			return
		}

		var err error
		currentUser, err = auth.CurrentUser(c)
		if err != nil {
			handleError(err, c)
			return
		}
	}

	profile := m.Profile{}
	if err := db.Where("user_id = ?", user.ID).First(&profile).Error; err != nil {
		handleError(&ErrDBNotFound{err}, c)
		return
	}

	resp := ProfileResp{
		Profile: profile,
		User:    user,
	}

	// If this is not the current user, do the match
	if user.ID != currentUser.ID {
		match, _ := matcher.Match(user.PoliticalMap, currentUser.PoliticalMap)
		resp.Match = match
	}

	c.JSON(http.StatusOK, resp)
}

// ProfileUpdate updates the profile
func ProfileUpdate(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		handleError(err, c)
		return
	}

	var profile m.Profile
	if err := db.Where("user_id = ?", user.ID).Find(&profile).Error; err != nil {
		handleError(&ErrDBNotFound{err}, c)
		return
	}

	// Update looking for
	lookingForS := c.PostForm("looking_for")
	lookingFor, err := strconv.Atoi(lookingForS)
	if err == nil {
		profile.LookingFor = lookingFor
	}

	profile.Summary = c.DefaultPostForm("summary", profile.Summary)

	if err := db.Save(&profile).Error; err != nil {
		handleError(&ErrDBInsert{err}, c)
		return
	}

	c.JSON(http.StatusOK, profile)
}
