package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"partisan/auth"
	"partisan/db"
)

// QuestionsIndex shows the question screen
func QuestionsIndex(c *gin.Context) {
	db, err := db.InitDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	user, err := auth.CurrentUser(c, &db)
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
