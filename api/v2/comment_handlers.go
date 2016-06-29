package v2

import (
	"net/http"
	"partisan/auth"
	"partisan/db"
	"partisan/logger"

	"partisan/models.v2/attachments"
	"partisan/models.v2/hashtags"
	"partisan/models.v2/notifications"
	"partisan/models.v2/posts"
	"partisan/models.v2/user_tags"

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

	comment, errs := posts.New(user.ID, binding)
	if len(errs) > 0 {
		c.AbortWithError(http.StatusNotAcceptable, errs)
		return
	}

	if err := db.Save(&comment).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, errs)
		return
	}

	tmpFile, _, err := c.Request.FormFile("attachments")
	if err == nil {
		defer tmpFile.Close()
		attachment, err := attachments.NewImage(user.ID, comment.ID, tmpFile)
		if err != nil {
			c.AbortWithError(http.StatusNotAcceptable, err)
			db.Delete(&comment) // clean up comment
			return
		}

		if err := db.Save(&attachment).Error; err != nil {
			c.AbortWithError(http.StatusNotAcceptable, err)
			db.Delete(&comment) // clean ucomment
			return
		}
	} else {
		logger.Error.Println(err)
	}

	if user.ID != post.UserID {
		n, errs := notifications.New(user.ID, post.UserID, comment)
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

	hashtags.FindAndCreate(comment, db)
	if err := usertags.Extract(comment, db, &pushClient); err != nil {
		logger.Error.Println(err)
	}

	c.JSON(http.StatusCreated, gin.H{"comment": comment})
}
