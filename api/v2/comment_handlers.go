package v2

import (
	"net/http"
	"partisan/auth"
	"partisan/db"

	"partisan/models.v2/posts"

	"github.com/gin-gonic/gin"
)

func CommentIndex(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	postID := c.Param("record_id")
	post, err := posts.GetByID(postID, user.ID, db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	// just an extra check to make sure we're not treating comments or likes like posts
	if post.Action != posts.APost {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	comments, err := post.GetComments(user.ID, db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"comments": comments})
}

func CommentCreate(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	postID := c.Param("record_id")
	post, err := posts.GetByID(postID, user.ID, db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	// just an extra check to make sure we're not treating comments or likes like posts
	if post.Action != posts.APost {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	var binding posts.CreatorBinding
	if err := c.Bind(&binding); err != nil {
		if err := c.BindJSON(&binding); err != nil {
			c.AbortWithError(http.StatusNotAcceptable, ErrBinding)
			return
		}
	}

	// forcing these values to prevent any funniness
	binding.ParentID = post.ID
	binding.ParentType = string(post.PostParentType())
	binding.Action = string(posts.AComment)

	comment, errs := posts.New(&binding)
	if len(errs) > 0 {
		c.AbortWithError(http.StatusNotAcceptable, errs)
		return
	}

	if err := db.Save(&comment).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, errs)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"comment": comment})
}
