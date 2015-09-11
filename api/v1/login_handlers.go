package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"partisan/db"
	m "partisan/models"
	"partisan/auth"
)

// LoginHandler Handle logging in
func LoginHandler(c *gin.Context) {
	// There's some potential for functions in here to panic, so I'm trying to recover
	defer func() {
		if r := recover(); r != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("%v", r))
		}
	}()

	// init db
	db, err := db.InitDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	email := c.Request.PostFormValue("email")
	password := c.Request.PostFormValue("password")

	user := m.User{}

	if err := db.Where(m.User{Email: email}).First(&user).Error; err != nil {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Couldn't find user"))
		return
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Password didn't match"))
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
