package v1

import (
	"net/http"
	"partisan/auth"
	"partisan/dao"
	"partisan/db"
	"strconv"

	"partisan/Godeps/_workspace/src/github.com/gin-gonic/gin"
)

// PostComments stores like data for ease
type PostComments struct {
	RecordID uint64
	Count    int
}

// FeedIndex shows all Feed Items for a particular user
func FeedIndex(c *gin.Context) {
	db := db.GetDB(c)

	user, _ := auth.CurrentUser(c)

	friendIDs, err := dao.ConfirmedFriendIDs(user, db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	friendIDs = append(friendIDs, user.ID)

	feedItems, err := dao.GetFeedByUserIDs(user.ID, friendIDs, db)
	if err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"feed_items": feedItems})
}

func FeedShow(c *gin.Context) {
	db := db.GetDB(c)
	user, _ := auth.CurrentUser(c)

	uID := c.Param("user_id")
	userID, err := strconv.Atoi(uID)
	if err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
	}

	feedItems, err := dao.GetFeedByUserIDs(user.ID, []uint64{uint64(userID)}, db)
	if err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"feed_items": feedItems})
}
