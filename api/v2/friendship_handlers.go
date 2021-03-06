package v2

import (
	"net/http"
	"partisan/auth"
	"partisan/db"
	"partisan/logger"
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

	fs, err := friendships.ListConfirmedByUserID(user.ID, db)
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

	collectUsers(user, &fs, friends)

	c.JSON(http.StatusOK, gin.H{"friendships": fs})
}

func FriendshipShow(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	friendID := c.Param("friend_id")

	f, err := friendships.GetByUserIDs(user.ID, friendID, db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"friendship": f})
}

func FriendshipCreate(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
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

	friend, err := users.GetByID(binding.FriendID, db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	// the friend is not looking for enemies, disallow if
	// match % isn't high enough
	if friend.LookingFor < 4 && !isPartisan(user, friend) {
		m, _ := matcher.Match(user.PoliticalMap, friend.PoliticalMap)
		if m < 0.75 {
			c.AbortWithError(http.StatusForbidden, ErrEnemy)
			return
		}
	}

	f, errs := friendships.New(user.ID, binding)
	if len(errs) > 0 {
		c.AbortWithError(http.StatusNotAcceptable, errs)
		return
	}

	if err := db.Save(&f).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	n, errs := notifications.New(user.ID, binding.FriendID, f)
	if len(errs) == 0 {
		db.Save(&n)

		if pn, err := n.NewPushNotification(db); err == nil {
			pushNotif := pn.Prepare()
			pushClient.Send(pushNotif)
		} else {
			logger.Error.Println("Error sending push notif:", err)
		}
	}

	c.JSON(http.StatusCreated, gin.H{"friendship": f})
}

func FriendshipUpdate(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	friendID := c.Param("friend_id")

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

	if errs := f.Update(binding); len(errs) > 0 {
		c.AbortWithError(http.StatusNotAcceptable, errs)
		return
	}

	if err := db.Save(&f).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	n, errs := notifications.New(user.ID, friendID, f)
	if len(errs) == 0 {
		db.Save(&n)

		if pn, err := n.NewPushNotification(db); err == nil {
			pushNotif := pn.Prepare()
			pushClient.Send(pushNotif)
		} else {
			logger.Error.Println("Error sending push notif:", err)
		}
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

	friendID := c.Param("user_id")

	f, err := friendships.GetByUserIDs(user.ID, friendID, db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if !f.CanDelete(user.ID) {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	if err := db.Delete(&f).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func collectUserIDs(userID string, fs friendships.Friendships) (ids []string) {
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
	xfs := []friendships.Friendship(*fs)

	for i := range xfs {
		for _, friend := range friends {
			if xfs[i].UserID == friend.ID || xfs[i].FriendID == friend.ID {
				xfs[i].Friend = friend
				if match, err := matcher.Match(user.PoliticalMap, friend.PoliticalMap); err == nil {
					xfs[i].Match = float64(int(match*1000)) / 10
				}
			}
		}
	}

	*fs = friendships.Friendships(xfs)
}

// We wanna allow everyone to be able to friend partisan
func isPartisan(us ...users.User) bool {
	for _, u := range us {
		if u.Email == "jon@partisan.io" {
			return true
		}
	}

	return false
}
