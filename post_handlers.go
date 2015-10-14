package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"partisan/auth"
	"partisan/db"
	m "partisan/models"
)

// PostShow shows a post
func PostShow(c *gin.Context) {
	db := db.GetDB(c)

	user, _ := auth.CurrentUser(c)

	postID := c.Param("record_id")

	var post m.Post
	var pUser m.User
	if err := db.Find(&post, postID).Related(&pUser).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	var attachment m.ImageAttachment
	db.Where("record_type = ? AND record_id = ?", "post", post.ID).Find(&attachment)

	likeCount := 0
	if err := db.Model(m.Like{}).Where("record_type = ? AND record_id = ?", "posts", post.ID).Count(&likeCount).Error; err != nil {
		fmt.Println(err)
	}

	c.HTML(http.StatusOK, "post", gin.H{
		"data": gin.H{
			"post":       post,
			"post_user":  pUser,
			"user":       user,
			"attachment": attachment,
			"like_count": likeCount,
		},
	})
}
