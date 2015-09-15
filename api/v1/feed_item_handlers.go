package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"partisan/auth"
	"partisan/db"
	m "partisan/models"
)

// FeedIndex shows all Feed Items for a particular user
func FeedIndex(c *gin.Context) {
	db, err := db.InitDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	db.LogMode(true)

	user, _ := auth.CurrentUser(c, &db)

	friendIDs, err := ConfirmedFriendIDs(user, c, &db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

        friendIDs = append(friendIDs, user.ID)

	feedItems := []m.FeedItem{}
	if err := db.Where("user_id IN (?)", friendIDs).Order("created_at desc").Find(&feedItems).Error; err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// TODO: This is an N+1 problem...
	for k, item := range feedItems {
		if item.RecordType == "post" {
			post := m.Post{}
			user := m.User{}
			db.First(&post, item.RecordID).Related(&user) // right here

			feedItems[k].Record = PostResponse{
				Post: post,
				User: user,
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"feed_items": feedItems})
}
