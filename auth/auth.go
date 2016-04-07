package auth

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"partisan/db"
	m "partisan/models"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
)

const (
	hmacKeyPath = "keys/hmac.key"
)

type ErrNoUser struct{}

func (err *ErrNoUser) Error() string {
	return "User ID not set"
}

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

		var userID int // i know, i know, but trying to cast all of these things to something sensible was killing me

		if okTok {
			token, err := jwt.Parse(tokn[0], func(token *jwt.Token) (interface{}, error) {
				// Don't forget to validate the alg is what you expect:
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}

				return hmacKey, nil
			})

			if err != nil {
				c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Problems parsing token"))
				return
			}

			userID, err = strconv.Atoi(fmt.Sprintf("%.f", token.Claims["user_id"]))
			if err != nil {
				c.AbortWithError(http.StatusUnauthorized, err)
				return
			}
		} else {
			if sess.Get("user_id") == nil {
				// check if there's a key in the request, and if so check if that route is
				// whitelisted to use keys. Should only be used for sockets that talk to mobile
				if len(c.Query("key")) > 0 && allowedWithKey(c) {
					c.Next()
					return
				}

				c.Redirect(http.StatusFound, redirectPath)
				return
			}

			userID, _ = strconv.Atoi(fmt.Sprintf("%d", sess.Get("user_id")))
		}

		// Check this is the right user with correct API key
		db := db.GetDB(c)
		user := m.User{}
		if err := db.First(&user, userID).Error; err != nil {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Couldn't find user: %d", userID))
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
	fmt.Println(user.ID)

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(hmacKey)

	sess := sessions.Default(c)
	sess.Set("user_id", user.ID)
	sess.Save()

	fmt.Println(tokenString)

	return tokenString, err
}

// Logout a user
func Logout(c *gin.Context) {
	sess := sessions.Default(c)
	sess.Clear()
	sess.Save()
}

// CurrentUser gets the current user from the session
func CurrentUser(c *gin.Context) (u m.User, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = &ErrNoUser{}
		}
	}()

	user, ok := c.Get("user")
	if !ok {
		return m.User{}, &ErrNoUser{}
	}

	u = user.(m.User)

	return
}

func allowedWithKey(c *gin.Context) bool {
	allowed := []string{"partisan/api/v1.NotificationsCount"}
	for _, v := range allowed {
		if c.HandlerName() == v {
			return true
		}
	}

	return false
}
