package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"partisan/auth"
)

// QuestionsIndex shows the question screen
func QuestionsIndex(c *gin.Context) {
	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	c.HTML(http.StatusOK, "questions", gin.H{
		"title": "Answer Questions",
		"data": gin.H{
			"user": user,
		},
	})
}
