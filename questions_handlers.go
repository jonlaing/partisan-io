package main

import (
	"net/http"
	"partisan/auth"

	"github.com/gin-gonic/gin"
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
