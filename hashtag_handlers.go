package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"partisan/auth"
	"partisan/db"
)

// HashtagShow renders the HTML for the hashtag search
func HashtagShow(c *gin.Context) {
	db, err := db.InitDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	user, _ := auth.CurrentUser(c)

	search := c.Query("q")

	c.HTML(http.StatusOK, "hashtags", gin.H{
		"title": "Search Hashtags:" + search,
		"data": gin.H{
			"user":   user,
			"search": search,
		},
	})
}
