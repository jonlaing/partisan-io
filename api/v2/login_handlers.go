package v2

import (
	"net/http"
	"partisan/auth"
	"partisan/db"

	"partisan/models.v2/users"

	"github.com/gin-gonic/gin"
)

type loginFields struct {
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
	DeviceToken string `json:"device_token"`
}

func LoginHandler(c *gin.Context) {
	db := db.GetDB(c)

	var fields loginFields
	if err := c.BindJSON(&fields); err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	user, err := users.GetByEmail(fields.Email, db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if err := user.CheckPassword(fields.Password); err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	token, err := auth.Login(&user, fields.DeviceToken, c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	if err := db.Save(&user).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}

func LogoutHandler(c *gin.Context) {
	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	auth.Logout(&user, c)
	db.Save(&user)
}
