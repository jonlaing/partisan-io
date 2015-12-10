package v1

import (
	"fmt"
	"net/http"
	"partisan/auth"
	"partisan/dao"
	"partisan/db"
	"partisan/matcher"
	m "partisan/models"
	"strconv"
	"time"

	"partisan/Godeps/_workspace/src/github.com/gin-gonic/gin"
)

// FriendResp is the JSON schema we respond with
type FriendResp struct {
	User  m.User  `json:"user"`
	Match float64 `json:"match"`
}

// FriendshipIndex returns all friends as a slice of m.User (in JSON)
func FriendshipIndex(c *gin.Context) {
	db := db.GetDB(c)

	user, _ := auth.CurrentUser(c)

	friends, err := dao.ConfirmedFriends(user, db)
	if err != nil {
		return handleError(err, c)
	}

	var friendships []FriendResp
	for _, f := range friends {
		if match, err := matcher.Match(user.PoliticalMap, f.PoliticalMap); err != nil {
			fmt.Println(err)
		} else {
			match = float64(int(match*1000)) / 10
			friendships = append(friendships, FriendResp{User: f, Match: match})
		}
	}

	c.JSON(http.StatusOK, friendships)
}

// FriendshipShow shows a friendship
func FriendshipShow(c *gin.Context) {
	db := db.GetDB(c)

	user, _ := auth.CurrentUser(c)

	fID := c.Param("friend_id")
	friendID, err := strconv.ParseUint(fID, 10, 64)
	if err != nil {
		return handleError(&ErrParseID{err}, c)
	}

	var friendship m.Friendship
	if friendship, err = dao.GetFriendship(user, friendID, db); err != nil {
		return handleError(err, c)
	}

	c.JSON(http.StatusOK, friendship)
}

// FriendshipCreate handles making a new friendship
func FriendshipCreate(c *gin.Context) {
	db := db.GetDB(c)

	user, _ := auth.CurrentUser(c)

	fID := c.PostForm("friend_id")
	friendID, err := strconv.ParseUint(fID, 10, 64)
	if err != nil {
		return handleError(&ErrParseID{err}, c)
	}

	f := m.Friendship{
		UserID:    user.ID,
		FriendID:  friendID,
		Confirmed: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Create(&f).Error; err != nil {
		return handleError(&ErrDBInsert{err}, c)
	}

	m.NewNotification(&f, user.ID, db)

	c.JSON(http.StatusCreated, f)
}

// FriendshipConfirm allows a user to accept a friend request
func FriendshipConfirm(c *gin.Context) {
	db := db.GetDB(c)

	user, _ := auth.CurrentUser(c)

	fID := c.PostForm("friend_id")
	friendID, err := strconv.ParseUint(fID, 10, 64)
	if err != nil {
		return handleError(&ErrParseID{err}, c)
	}

	var f m.Friendship
	// only the friend can confirm, so we put friendID in the user slot and userID in the friend slot
	if err := db.Where("friend_id = ? AND user_id = ?", user.ID, friendID).First(&f).Error; err != nil {
		return handleError(&ErrDBNotFound{err}, c)
	}

	f.Confirmed = true

	if err := db.Save(&f).Error; err != nil {
		return handleError(&ErrDBInsert{err}, c)
	}

	m.NewNotification(&f, user.ID, db)

	c.JSON(http.StatusOK, f)
}

// FriendshipDestroy unfriends
func FriendshipDestroy(c *gin.Context) {
	db := db.GetDB(c)

	user, _ := auth.CurrentUser(c)

	fID := c.Param("friend_id")
	friendID, err := strconv.ParseUint(fID, 10, 64)
	if err != nil {
		return handleError(&ErrParseID{err}, c)
	}

	// We have to look for two possible friendships
	f1 := m.Friendship{}
	f2 := m.Friendship{}

	if err := db.Find(&f1, user.ID).Error; err == nil {
		if err := db.Delete(&f1).Error; err != nil {
			return handleError(&ErrDBDelete{err}, c)
		}
	}

	if err := db.Find(&f2, friendID).Error; err == nil {
		if err := db.Delete(&f2).Error; err != nil {
			return handleError(&ErrDBNotFound{err}, c)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
