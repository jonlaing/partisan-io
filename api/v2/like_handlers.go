package v2

import (
	"net/http"
	"partisan/auth"
	"partisan/db"
	"partisan/logger"

	"partisan/models.v2/notifications"
	"partisan/models.v2/posts"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func LikeCreate(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	id := c.Param("record_id")
	post, err := posts.GetByID(id, user.ID, db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	// If it's already been liked, get rid of it
	if post.Liked {
		deleteLike(&post, user.ID, c, db)
		return
	}

	// If this is a new like, create it. Really don't need too much from
	// the user, so build the binding manually
	binding := posts.CreatorBinding{
		Action:     string(posts.ALike),
		ParentType: string(post.PostParentType()),
		ParentID:   post.ID,
	}

	like, errs := posts.New(user.ID, binding)
	if len(errs) > 0 {
		c.AbortWithError(http.StatusNotAcceptable, errs)
		return
	}

	if err := db.Save(&like).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	post.LikeCount++
	post.Liked = true

	if user.ID != post.UserID {
		n, errs := notifications.New(user.ID, post.UserID, like)
		if len(errs) == 0 {
			db.Save(&n)

			if pn, err := n.NewPushNotification(db); err == nil {
				pushNotif := pn.Prepare()
				pushClient.Send(pushNotif)
			} else {
				logger.Error.Println("Error sending push notif:", err)
			}
		}
	}

	c.JSON(http.StatusCreated, gin.H{string(post.Action): post})
}

func deleteLike(post *posts.Post, userID string, c *gin.Context, db *gorm.DB) {
	like, err := post.GetLikeByUserID(userID, db)
	if err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	if like.CanDelete(userID, db) {
		if err := db.Delete(&like).Error; err != nil {
			c.AbortWithError(http.StatusNotAcceptable, err)
			return
		}

		post.LikeCount--
		post.Liked = false
		c.JSON(http.StatusOK, gin.H{string(post.Action): post})
	}
}
