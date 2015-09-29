package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"partisan/db"
)

// UserTagShow is for notifications. It will find the tag, then redirect to the related record
func UserTagShow(c *gin.Context) {
	db, err := db.InitDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	c.JSON(http.StatusOK, gin.H{})
}
