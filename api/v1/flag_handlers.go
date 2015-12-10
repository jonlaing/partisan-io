package v1

import (
	"net/http"
	"partisan/auth"
	"partisan/db"
	m "partisan/models"
	"time"

	"partisan/Godeps/_workspace/src/github.com/gin-gonic/gin"
)

// FlagCreate flags a record
func FlagCreate(c *gin.Context) {
	db := db.GetDB(c)

	user, _ := auth.CurrentUser(c)

	var flag m.Flag
	if err := c.BindJSON(&flag); err != nil {
		return handleError(&ErrBinding{err}, c)
	}

	flag.UserID = user.ID
	flag.CreatedAt = time.Now()
	flag.UpdatedAt = time.Now()

	if err := db.Save(&flag).Error; err != nil {
		return handleError(&ErrDBInsert{err}, c)
	}

	c.JSON(http.StatusCreated, gin.H{"message": "flagged"})
}
