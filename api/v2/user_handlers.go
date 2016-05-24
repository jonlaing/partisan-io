package v2

import (
	"net/http"

	"partisan/auth"
	"partisan/db"

	"github.com/gin-gonic/gin"
	"partisan/models.v2/users"
)

func UserCreate(c *gin.Context) {
	db := db.GetDB(c)

	var binding users.CreatorBinding
	if err := c.BindJSON(&binding); err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	user, errs := users.New(binding)
	if len(errs) > 0 {
		c.AbortWithError(http.StatusNotAcceptable, errs)
		return
	}

	token, err := auth.Login(&user, c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	if err := db.Save(&user).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user, "token": token})
}

func UserShow(c *gin.Context) {
	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func UserUpdate(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	var binding users.UpdaterBinding
	if err := c.BindJSON(&binding); err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	if errs := user.Update(binding); len(errs) > 0 {
		c.AbortWithError(http.StatusNotAcceptable, errs)
		return
	}

	if err := db.Save(&user).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func UserAvatarUpload(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	tmpFile, _, err := c.Request.FormFile("avatar")
	if err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}
	defer tmpFile.Close()

	if err := user.AttachAvatar(tmpFile); err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	if err := db.Save(&user).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}