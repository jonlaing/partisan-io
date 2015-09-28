package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"partisan/auth"
	"partisan/db"
	m "partisan/models"
	"strconv"
	"time"
)

// FriendshipIndex returns all friends as a slice of m.User (in JSON)
func FriendshipIndex(c *gin.Context) {
	db, err := db.InitDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	user, _ := auth.CurrentUser(c, &db)

	friendIDs, err := FriendIDs(user, c, &db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	var friends []m.User
	if err := db.Where("user_id IN ?", friendIDs).Find(&friends).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, friends)
}

// FriendshipShow shows a friendship
func FriendshipShow(c *gin.Context) {
	db, err := db.InitDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	user, _ := auth.CurrentUser(c, &db)

	fID := c.Param("friend_id")
	friendID, err := strconv.ParseUint(fID, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var f1, f2 m.Friendship

	if err := db.Where("user_id = ? AND friend_id = ?", user.ID, friendID).Find(&f1).Error; err == nil {
		c.JSON(http.StatusOK, f1)
		return
	}

	if err := db.Where("user_id = ? AND friend_id = ?", friendID, user.ID).Find(&f2).Error; err == nil {
		c.JSON(http.StatusOK, f2)
		return
	}

	c.AbortWithError(http.StatusNotFound, fmt.Errorf("Couldn't find friendship between User: %d and Friend: %d", user.ID, friendID))
	return
}

// FriendshipCreate handles making a new friendship
func FriendshipCreate(c *gin.Context) {
	db, err := db.InitDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	user, _ := auth.CurrentUser(c, &db)

	fID := c.PostForm("friend_id")
	friendID, err := strconv.ParseUint(fID, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	f := m.Friendship{
		UserID:    user.ID,
		FriendID:  friendID,
		Confirmed: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Create(&f).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	m.NewNotification(&f, user.ID, &db)

	c.JSON(http.StatusCreated, f)
}

// FriendshipConfirm allows a user to accept a friend request
func FriendshipConfirm(c *gin.Context) {
	db, err := db.InitDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	user, _ := auth.CurrentUser(c, &db)

	fID := c.PostForm("friend_id")
	friendID, err := strconv.ParseUint(fID, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	f := m.Friendship{}

	// only the friend can confirm, so we put friendID in the user slot and userID in the friend slot
	if err := db.Where("friend_id = ? AND user_id = ?", user.ID, friendID).First(&f).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	f.Confirmed = true

	if err := db.Save(&f).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	m.NewNotification(&f, user.ID, &db)

	c.JSON(http.StatusOK, f)
}

// FriendshipDestroy unfriends
func FriendshipDestroy(c *gin.Context) {
	db, err := db.InitDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	user, _ := auth.CurrentUser(c, &db)

	fID := c.Param("friend_id")
	friendID, err := strconv.ParseUint(fID, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// We have to look for two possible friendships
	f1 := m.Friendship{}
	f2 := m.Friendship{}

	if err := db.Find(&f1, user.ID).Error; err == nil {
		if err := db.Delete(&f1).Error; err != nil {
			c.AbortWithError(http.StatusNotAcceptable, err)
			return
		}
	}

	if err := db.Find(&f2, friendID).Error; err == nil {
		if err := db.Delete(&f2).Error; err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// FriendIDs returns all userIDs associated with a user's friends
func FriendIDs(user m.User, c *gin.Context, db *gorm.DB) (friendIDs []uint64, err error) {
	var friendships []m.Friendship
	if err = db.Where("user_id = ? OR friend_id = ?", user.ID, user.ID).Find(&friendships).Error; err != nil {
		return
	}

	if len(friendships) < 1 {
		err = fmt.Errorf("Couldn't find any friends associated with User: %d", user.ID)
		return
	}

	for _, v := range friendships {
		if v.UserID != user.ID {
			friendIDs = append(friendIDs, v.UserID)
		} else {
			friendIDs = append(friendIDs, v.FriendID)
		}
	}

	return
}

// ConfirmedFriendIDs returns all userIDs associated with a user's CONFIRMED friends
func ConfirmedFriendIDs(user m.User, c *gin.Context, db *gorm.DB) (friendIDs []uint64, err error) {
	var friendships []m.Friendship
	if err = db.Where("user_id = ? OR friend_id = ? AND confirmed = ?", user.ID, user.ID, true).Find(&friendships).Error; err != nil {
		return
	}

	if len(friendships) < 1 {
		err = fmt.Errorf("Couldn't find any friends associated with User: %d", user.ID)
		return
	}

	for _, v := range friendships {
		if v.UserID != user.ID {
			friendIDs = append(friendIDs, v.UserID)
		} else {
			friendIDs = append(friendIDs, v.FriendID)
		}
	}

	return
}
