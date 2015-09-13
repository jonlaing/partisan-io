package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"partisan/auth"
	"partisan/db"
	m "partisan/models"
	"strconv"
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
	postID := c.Param("post_id")
	if pID, err := strconv.ParseUint(postID, 10, 64); err == nil {
		var count int
		db.Where("record_type = ? AND record_id = ?", "post", pID).Find(&[]m.Like{}).Count(&count)

		var userCount int
		db.Where("record_type = ? AND record_id = ? AND user_id= ? ", "post", pID, user.ID).Find(&[]m.Like{}).Count(&userCount)

		c.JSON(http.StatusOK, gin.H{"record_type": "post", "record_id": pID, "like_count": count, "liked": userCount > 0})
		return
	}
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

	postID := c.Param("post_id")
	if pID, err := strconv.ParseUint(postID, 10, 64); err == nil {
		like = m.Like{
			UserID:     user.ID,
			RecordID:   pID,
			RecordType: "post",
			IsDislike:  false,
		}

		// get the current count of likes
		if err := db.Where("record_type = ? AND record_id = ?", "post", pID).Find(&[]m.Like{}).Count(&count).Error; err != nil {
			count = 0
		}

		// get the number of likes associated with a user
		if err := db.Where("record_type = ? AND record_id = ? AND user_id = ?", "post", pID, user.ID).Find(&[]m.Like{}).Count(&userCount).Error; err != nil {
			userCount = 0
		}

		// if the there is no like record associated with the user
		if userCount < 1 {
			// create a like record
			if err := db.Create(&like).Error; err != nil {
				c.AbortWithError(http.StatusNotAcceptable, err)
				return
			}
			// return the old count + 1
			c.JSON(http.StatusCreated, gin.H{"record_type": "post", "record_id": pID, "like_count": count + 1, "liked": true})
		} else {
			if err := db.Where("record_type = ? AND record_id = ? AND user_id = ?", "post", pID, user.ID).Delete(m.Like{}).Error; err != nil {
				fmt.Println("problem deleting like:", err)
			}
			c.JSON(http.StatusCreated, gin.H{"record_type": "post", "record_id": pID, "like_count": count - 1, "liked": false})
		}
		return
	}
}
