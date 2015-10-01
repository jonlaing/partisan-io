package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"partisan/auth"
	"partisan/db"
	m "partisan/models"
	"time"
)

// FlagCreate flags a record
func FlagCreate(c *gin.Context) {
	db := db.GetDB(c)

	user, _ := auth.CurrentUser(c)

	var flag m.Flag
	if err := c.BindJSON(&flag); err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	flag.UserID = user.ID
	flag.CreatedAt = time.Now()
	flag.UpdatedAt = time.Now()

	if err := db.Save(&flag).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "flagged"})
}
