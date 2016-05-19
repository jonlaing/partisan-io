package v2

import (
	"net/http"
	"partisan/auth"
	"partisan/db"

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

	id := c.Params("record_id")
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

	// If this is a new like, create it
	var binding CreatorBinding
	if err := c.BindJSON(&binding); err != nil {
		c.AbortWithError(http.StatusNotAcceptable)
		return
	}

	binding.Action = string(posts.ALike)
	binding.ParentType = string(post.PostParentType())
	binding.ParentID = post.ID

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

	c.JSON(http.StatusOK, gin.H{string(post.Action): post})
}

func deleteLike(post *posts.Post, userID string, c *gin.Context, db *gorm.DB) {
	like, err := post.GetLikeByUserID(user.ID, db)
	if err != nil {
		c.AbortWithError(http.StatusNotAcceptable)
		return
	}

	if post.CanDelete(user.ID, db) {
		if err := db.Delete(&post).Error; err != nil {
			c.AbortWithError(http.StatusNotAcceptable)
			return
		}

		post.LikeCount--
		post.Liked = false
		c.JSON(http.StatusOK, gin.H{string(post.Action): post})
	}
}
