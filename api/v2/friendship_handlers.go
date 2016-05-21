package v2

import (
	"net/http"
	"partisan/auth"
	"partisan/db"
	"partisan/matcher"

	"partisan/models.v2/friendships"
	"partisan/models.v2/notifications"
	"partisan/models.v2/users"

	"github.com/gin-gonic/gin"
)

func FriendshipIndex(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	fs, err := friendships.GetByUserID(user.ID, db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	friendIDs := collectUserIDs(user.ID, fs)
	friends, err := users.ListByIDs(friendIDs, db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	collectUsers(userID, &fs, friends)

	c.JSON(http.StatusOK, gin.H{"friendships": fs})
}

func FriendshipCreate(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CreateUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	var binding friendships.CreatorBinding
	if err := c.BindJSON(&binding); err != nil {
		c.AbortWithError(http.StatusNotAcceptable, ErrBinding)
		return
	}

	if friendships.Exists(user.ID, binding.FriendID, db) {
		c.AbortWithError(http.StatusConflict, ErrAlreadyExists)
		return
	}

	f, err := friendships.New(user.ID, binding)
	if err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	if err := db.Save(&f).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	n, errs := notifications.New(user.ID, binding.FriendID, f)
	if len(errs) == 0 {
		db.Save(&n)
	}

	c.JSON(http.StatusOK, gin.H{"friendship": f})
}

func FriendshipUpdate(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CreateUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	friendID := c.Param("record_id")

	f, err := friendships.GetByUserIDs(user.ID, friendID, db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	var binding friendships.UpdaterBinding
	if err := c.BindJSON(&binding); err != nil {
		c.AbortWithError(http.StatusNotAcceptable, ErrBinding)
		return
	}

	if err := f.Update(binding); err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	if err := db.Save(&f).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	n, errs := notifications.New(user.ID, friendID, f)
	if len(errs) == 0 {
		db.Save(&n)
	}

	c.JSON(http.StatusOK, gin.H{"friendship": f})
}

func FriendshipDestroy(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	friendID := c.Param("record_id")

	f, err := friendships.GetByUserIDs(user.ID, friendID, db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if !f.CanDelete(user.ID) {
		c.AborthWithError(http.StatusUnauthorized, err)
		return
	}

	if err := db.Delete(&f).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func collectUserIDs(userID string, fs []Friendship) (ids []string) {
	for _, f := range fs {
		if f.UserID == userID {
			ids = append(ids, f.FriendID)
		}

		if f.FriendID == userID {
			ids = append(ids, f.UserID)
		}
	}

	return
}

func collectUsers(user users.User, fs *friendships.Friendships, friends []users.User) {
	friendships := []Frienship(*fs)

	for i := range friendships {
		for _, friend := range friends {
			if friendhips[i].UserID == friend.ID || friendhips[i].FriendID == friend.ID {
				friendhips[i].User = friend
				if match, err := matcher.Match(user.PoliticalMap, friend.PoliticalMap); err == nil {
					friendhips[i].Match = float64(int(match*1000)) / 10
				}
			}
		}
	}

	*fs = Friendships(friendships)
}
