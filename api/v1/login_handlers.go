package v1

import (
	"fmt"
	"net/http"
	"partisan/auth"
	"partisan/db"
	m "partisan/models"

	"partisan/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"partisan/Godeps/_workspace/src/golang.org/x/crypto/bcrypt"
)

// LoginHandler Handle logging in
func LoginHandler(c *gin.Context) {
	// There's some potential for functions in here to panic, so I'm trying to recover
	defer func() {
		if r := recover(); r != nil {
			handleError(fmt.Errorf("%v", r), c)
			return
		}
	}()

	db := db.GetDB(c)

	email := c.Request.PostFormValue("email")
	password := c.Request.PostFormValue("password")

	user := m.User{}

	if err := db.Where(m.User{Email: email}).First(&user).Error; err != nil {
		handleError(&ErrUserNotFound{}, c)
		return
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		handleError(&ErrPasswordMatch{}, c)
		return
	}

	// user.APIKey = uuid.NewRandom().String()
	// user.APIKeyExp = time.Now().Add(time.Hour * 72)

	db.Save(&user)

	tokenString, _ := auth.Login(user, c)

	// User.APIKey is not exported via JSON so we have to manually add it
	c.JSON(http.StatusOK, gin.H{"token": tokenString, "user": user, "api_key": user.APIKey})
}

// LogoutHandler logs the user out (duh)
func LogoutHandler(c *gin.Context) {
	auth.Logout(c)

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}
