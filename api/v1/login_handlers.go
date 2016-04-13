package v1

import (
	"fmt"
	"net/http"
	"partisan/auth"
	"partisan/db"
	m "partisan/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type loginFields struct {
	Email    string
	Password string
}

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

	var login loginFields

	if err := c.BindJSON(&login); err != nil {
		login.Email = c.Request.PostFormValue("email")
		login.Password = c.Request.PostFormValue("password")
	}

	user := m.User{}

	if err := db.Where(m.User{Email: login.Email}).First(&user).Error; err != nil {
		handleError(&ErrUserNotFound{}, c)
		return
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(login.Password)); err != nil {
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
