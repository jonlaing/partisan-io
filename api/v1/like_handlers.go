package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"partisan/auth"
	"partisan/db"
	m "partisan/models"
)

// LikeCount shows the like count for a particular record
func LikeCount(c *gin.Context) {
	db, err := db.InitDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	user, _ := auth.CurrentUser(c, &db)
	rID, rType, err := getRecord(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	var count int
	db.Model(m.Like{}).Where("record_type = ? AND record_id = ?", rType, rID).Count(&count)

	// check if this user like this record
	var userCount int
	db.Model(m.Like{}).Where("record_type = ? AND record_id = ? AND user_id= ? ", rType, rID, user.ID).Count(&userCount)

	c.JSON(http.StatusOK, gin.H{"record_type": rType, "record_id": rID, "like_count": count, "liked": userCount > 0})
	return
}

// LikeCreate creates a Like for a particular record
func LikeCreate(c *gin.Context) {
	db, err := db.InitDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	var like m.Like
	var count int
	var userCount int
	user, _ := auth.CurrentUser(c, &db)

	rID, rType, err := getRecord(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	like = m.Like{
		UserID:     user.ID,
		RecordID:   rID,
		RecordType: rType,
		IsDislike:  false,
	}

	// get the current count of likes
	if err := db.Model(m.Like{}).Where("record_type = ? AND record_id = ?", rType, rID).Count(&count).Error; err != nil {
		count = 0
	}

	// get the number of likes associated with a user
	if err := db.Model(m.Like{}).Where("record_type = ? AND record_id = ? AND user_id = ?", rType, rID, user.ID).Count(&userCount).Error; err != nil {
		userCount = 0
	}

	// if the there is no like record associated with the user
	if userCount < 1 {
		// create a like record
		if err := db.Create(&like).Error; err != nil {
			c.AbortWithError(http.StatusNotAcceptable, err)
			return
		}
                m.NewNotification(&like, user.ID, &db)
		// return the old count + 1
		c.JSON(http.StatusCreated, gin.H{"record_type": rType, "record_id": rID, "like_count": count + 1, "liked": true})
	} else {
		if err := db.Where("record_type = ? AND record_id = ? AND user_id = ?", rType, rID, user.ID).Delete(m.Like{}).Error; err != nil {
			fmt.Println("problem deleting like:", err)
		}
		if err := db.Where("record_type = ? AND record_id = ? AND user_id = ?", rType, rID, user.ID).Delete(m.Notification{}).Error; err != nil {
			fmt.Println("problem deleting notificatino:", err)
		}
		c.JSON(http.StatusCreated, gin.H{"record_type": rType, "record_id": rID, "like_count": count - 1, "liked": false})
	}
	return
}
