package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"partisan/db"
	m "partisan/models"
	"time"
)

// FriendshipJSON is used for processing incomming JSON data from request
type FriendshipJSON struct {
	UserID   uint64 `json:"user_id"`
	FriendID uint64 `json:"friend_id"`
}

// FriendshipCreate handles making a new friendship
func FriendshipCreate(c *gin.Context) {
	// type conversion from getting user from context can cause panic
	defer func() {
		if r := recover(); r != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}()

	db, err := db.InitDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	userID, ok := c.Get("user_id")
	if !ok {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("User ID not set"))
		return
	}

	var fJSON FriendshipJSON
	if err := c.BindJSON(&fJSON); err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
	}

	if fJSON.UserID != userID.(uint64) {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Friend request userID must match logged-in user ID"))
		return
	}

	f1 := m.Friendship{
		UserID:    fJSON.UserID,
		FriendID:  fJSON.FriendID,
		Confirmed: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// the second gets created for ease of searching db for record
	f2 := m.Friendship{
		UserID:    fJSON.FriendID,
		FriendID:  fJSON.UserID,
		Confirmed: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Create(&f1).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	if err := db.Create(&f2).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	c.JSON(http.StatusCreated, f1)
}

// FriendshipConfirm allows a user to accept a friend request
func FriendshipConfirm(c *gin.Context) {
	// type conversion from getting user from context can cause panic
	defer func() {
		if r := recover(); r != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}()

	db, err := db.InitDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	userID, ok := c.Get("user_id")
	if !ok {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("User ID not set"))
		return
	}

	var fJSON FriendshipJSON
	if err := c.BindJSON(&fJSON); err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
	}

	if fJSON.UserID != userID.(uint64) {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Friend request userID must match logged-in user ID"))
		return
	}

	f1 := m.Friendship{}
	f2 := m.Friendship{}

	if err := db.Where(&m.Friendship{UserID: fJSON.UserID, FriendID: fJSON.FriendID}).First(&f1).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if err := db.Where(&m.Friendship{UserID: fJSON.FriendID, FriendID: fJSON.UserID}).First(&f2).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	f1.Confirmed = true
	f2.Confirmed = true

	if err := db.Save(&f1).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	if err := db.Save(&f2).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "confirmed"})
}

// FriendshipDestroy unfriends
func FriendshipDestroy(c *gin.Context) {
	// type conversion from getting user from context can cause panic
	defer func() {
		if r := recover(); r != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}()

	db, err := db.InitDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	userID, ok := c.Get("user_id")
	if !ok {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("User ID not set"))
		return
	}

	var fJSON FriendshipJSON
	if err := c.BindJSON(&fJSON); err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
	}

	if fJSON.UserID != userID.(uint64) {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Friend request userID must match logged-in user ID"))
		return
	}

	f1 := m.Friendship{}
	f2 := m.Friendship{}

	if err := db.Find(&f1, fJSON.UserID).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if err := db.Find(&f2, fJSON.FriendID).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if err := db.Delete(&f1).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	if err := db.Delete(&f2).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
