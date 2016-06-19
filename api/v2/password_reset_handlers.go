package v2

import (
	"fmt"
	"net/http"
	"partisan/db"
	"partisan/emailer"
	"strings"

	"partisan/models.v2/password_reset"
	"partisan/models.v2/users"

	"github.com/gin-gonic/gin"
)

type PWInitBinding struct {
	Email string `json:"email"`
}

type PWBinding struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

func PasswordResetCreate(c *gin.Context) {
	db := db.GetDB(c)

	var binding PWInitBinding
	if err := c.BindJSON(&binding); err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	user, err := users.GetByEmail(binding.Email, db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	var reset pwreset.PasswordReset
	if err := db.Where("user_id = ?", user.ID).Find(&reset).Error; err == nil {
		db.Delete(&reset) // if there's an existing password reset, kill it
	}

	reset = pwreset.New(user.ID)
	if err := db.Create(&reset).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	if err := emailer.SendPasswordReset(user.Email, reset.ID); err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"email": hideEmail(user.Email)})
}

func PasswordResetShow(c *gin.Context) {
	db := db.GetDB(c)

	resetID := c.Param("reset_id")

	var reset pwreset.PasswordReset
	if err := db.Where("id = ?", resetID).Find(&reset).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if !reset.IsValid() {
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}

	c.JSON(http.StatusOK, gin.H{"reset_id": reset.ID})
}

func PasswordResetUpdate(c *gin.Context) {
	db := db.GetDB(c)

	resetID := c.Param("reset_id")

	var binding PWBinding
	if err := c.BindJSON(&binding); err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	var reset pwreset.PasswordReset
	if err := db.Where("id = ?", resetID).Find(&reset).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if !reset.IsValid() {
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}

	user, err := users.GetByID(reset.UserID, db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if user.Email != binding.Email {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if err := user.GeneratePasswordHash(binding.Password, binding.PasswordConfirm); err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	if err := db.Save(&user).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	db.Delete(&reset)

	c.Status(http.StatusOK)
}

func hideEmail(email string) (hidden string) {
	i := strings.Index(email, "@")
	buf := ""

	if i > 4 {
		for j := 0; j < i-4; j++ {
			buf += "*"
		}
		hidden = fmt.Sprintf("%s%s%s", email[:2], buf, email[i-2:])
	} else {
		for j := 0; j < i; j++ {
			buf += "*"
		}
		hidden = fmt.Sprintf("%s%s%s", email[:1], buf, email[i:])
	}

	return
}
