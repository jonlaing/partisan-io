package v2

import (
	"net/http"
	"partisan/auth"
	"partisan/db"

	"partisan/models.v2/attachments"
	"partisan/models.v2/friendships"
	"partisan/models.v2/hashtags"
	"partisan/models.v2/posts"

	"github.com/gin-gonic/gin"
)

func PostIndex(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	friendIDs, err := friendships.ListConfirmedIDsByUserID(user.ID, db)
	if err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	if len(friendIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{"posts": posts.Posts{}})
		return
	}

	page := getPage(c)

	posts, err := posts.GetFeedByUserIDs(user.ID, friendIDs, page*25, db)
	if err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

func PostShow(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{"post": post})
}

func PostCreate(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	var binding posts.CreatorBinding
	if err := c.Bind(&binding); err != nil {
		if err := c.BindJSON(&binding); err != nil {
			c.AbortWithError(http.StatusNotAcceptable, ErrBinding)
			return
		}
	}

	post, errs := posts.New(user.ID, binding)
	if len(errs) > 0 {
		c.AbortWithError(http.StatusNotAcceptable, errs)
		return
	}

	if err := db.Save(&post).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	tmpFile, _, err := c.Request.FormFile("attachments")
	if err == nil {
		defer tmpFile.Close()
		attachment, err := attachments.NewImage(user.ID, post.ID, tmpFile)
		if err != nil {
			c.AbortWithError(http.StatusNotAcceptable, err)
			db.Delete(&post) // clean up post
			return
		}

		if err := db.Save(&attachment).Error; err != nil {
			c.AbortWithError(http.StatusNotAcceptable, err)
			db.Delete(&post) // clean up post
			return
		}
	}

	hashtags.FindAndCreate(post, db)

	c.JSON(http.StatusCreated, gin.H{"post": post})
}

func PostUpdate(c *gin.Context) {
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

	var binding posts.UpdaterBinding
	if err := c.Bind(&binding); err != nil {
		if err := c.BindJSON(&binding); err != nil {
			c.AbortWithError(http.StatusNotAcceptable, ErrBinding)
			return
		}
	}

	errs := post.Update(user.ID, binding)
	if len(errs) > 0 {
		c.AbortWithError(http.StatusNotAcceptable, errs)
		return
	}

	if err := db.Save(&post).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	hashtags.FindAndCreate(post, db)

	c.JSON(http.StatusOK, gin.H{"post": post})
}

func PostDestroy(c *gin.Context) {
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

	if !post.CanDelete(user.ID, db) {
		c.AbortWithError(http.StatusUnauthorized, ErrCannotDelete)
		return
	}

	if err := db.Delete(&post).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	db.Delete(posts.Post{}, "parent_id = ?", id)

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
