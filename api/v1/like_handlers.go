package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
        m "partisan/models"
        "partisan/db"
	"partisan/auth"
)


// LikeCreate creates a Like for a particular record
func LikeCreate(c *gin.Context) {
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

	var like m.Like
	user, _ := auth.CurrentUser(c, &db)

	postID := c.Param("post_id")
	if pID, err := strconv.ParseUint(postID, 10, 64); err == nil {
		like = m.Like{
			UserID:     user.ID,
			RecordID:   pID,
			RecordType: "post",
		}

		if err := db.Create(&like).Error; err != nil {
			c.AbortWithError(http.StatusNotAcceptable, err)
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "liked"})
		return
	}
}

// DislikeCreate creates a Dislike for a particular record
func DislikeCreate(c *gin.Context) {
	db, err := db.InitDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	var dislike m.Dislike
	user, _ := auth.CurrentUser(c, &db)

	postID, ok := c.Get("post_id")
	if ok {
		dislike = m.Dislike{
			UserID:     user.ID,
			RecordID:   postID.(uint64),
			RecordType: "post",
		}

		if err := db.Create(&dislike).Error; err != nil {
			c.AbortWithError(http.StatusNotAcceptable, err)
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "disliked"})
		return
	}
}
