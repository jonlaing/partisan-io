package auth

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"partisan/db"
	m "partisan/models"
	"strconv"
	"time"

	jwt "partisan/Godeps/_workspace/src/github.com/dgrijalva/jwt-go"
	"partisan/Godeps/_workspace/src/github.com/gin-gonic/contrib/sessions"
	"partisan/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"partisan/Godeps/_workspace/src/github.com/gorilla/securecookie"
)

const (
	hmacKeyPath = "keys/hmac.key"
)

var hmacKey []byte

func init() {
	var err error
	hmacKey, err = ioutil.ReadFile(hmacKeyPath)
	if err != nil {
		panic(err)
	}
}

// UserSession is a database record that holds user sessions (go figure)
type UserSession struct {
	ID         uint `gorm:"primary_key"`
	UserID     uint
	SigningKey []byte
	CreatedAt  time.Time
	ExpiresAt  time.Time
}

// LoginJSON is the expected schema from a login form
type LoginJSON struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

var store = sessions.NewCookieStore(securecookie.GenerateRandomKey(16))

// Auth is the authentication middleware
func Auth(redirectPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		sess := sessions.Default(c)
		tokn, okTok := c.Request.Header["X-Auth-Token"]

		var userID uint64
		// var apiKey string

		if okTok {
			token, err := jwt.Parse(tokn[0], func(token *jwt.Token) (interface{}, error) {
				// Don't forget to validate the alg is what you expect:
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}

				return hmacKey, nil
			})

			if err != nil {
				c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Problems making token"))
				return
			}

			userID, _ = strconv.ParseUint(fmt.Sprintf("%d", token.Claims["user_id"]), 10, 0)
			// apiKey = fmt.Sprintf("%s", token.Claims["api_key"])
		} else {
			if sess.Get("user_id") == nil {
				c.Redirect(http.StatusFound, redirectPath)
				return
			}

			userID, _ = strconv.ParseUint(fmt.Sprintf("%d", sess.Get("user_id")), 10, 0)
			// apiKey = fmt.Sprintf("%s", sess.Get("api_key"))
		}

		// Check this is the right user with correct API key
		db := db.GetDB(c)
		user := m.User{}
		if err := db.First(&user, userID).Error; err != nil {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Couldn't find user"))
			return
		}

		c.Set("user", user)
		c.Next() // continue on to next endpoint
	}
}

// Login a user
func Login(user m.User, c *gin.Context) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	// Set some claims
	token.Claims["user_id"] = user.ID
	token.Claims["api_key"] = user.APIKey
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(hmacKey)

	sess := sessions.Default(c)
	sess.Set("user_id", user.ID)
	sess.Save()

	return tokenString, err
}

// Logout a user
func Logout(c *gin.Context) {
	sess := sessions.Default(c)
	sess.Clear()
	sess.Save()
}

// CurrentUser gets the current user from the session
func CurrentUser(c *gin.Context) (m.User, error) {
	user, ok := c.Get("user")
	if !ok {
		return user.(m.User), fmt.Errorf("User ID not set")
	}

	return user.(m.User), nil
}
