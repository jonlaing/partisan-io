package v1

import (
	"net/http"
	"partisan/db"
	"partisan/models"

	"github.com/gin-gonic/gin"
)

func PreSignUpCreate(c *gin.Context) {
	db := db.GetDB(c)

	var signup models.PreSignUp
	if err := c.BindJSON(&signup); err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	if err := db.Create(&signup).Error; err != nl {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Signed Up!"})
}
