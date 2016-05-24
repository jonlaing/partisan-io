package v2

import (
	"net/http"
	"partisan/auth"
	"partisan/db"
	"time"

	"partisan/models.v2/flags"

	"github.com/gin-gonic/gin"
)

// FlagCreate flags a record
func FlagCreate(c *gin.Context) {
	db := db.GetDB(c)

	user, _ := auth.CurrentUser(c)

	var flag flags.Flag
	if err := c.BindJSON(&flag); err != nil {
		c.AbortWithError(http.StatusNotAcceptable, ErrBinding)
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
