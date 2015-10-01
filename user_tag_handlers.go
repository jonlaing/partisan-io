package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// UserTagShow is for notifications. It will find the tag, then redirect to the related record
func UserTagShow(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
