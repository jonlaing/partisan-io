package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"partisan/db"
	"time"

	"partisan/models.v2/tickets"
	"partisan/models.v2/users"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/nu7hatch/gouuid"

	"github.com/gin-gonic/gin"
)

const (
	hmacKeyPath = "keys/hmac.key"
)

// ErrNoUser signifies that the user couldn't be found in the context
var ErrNoUser = errors.New("User ID not set")

var hmacKey []byte

func init() {
	var err error
	hmacKey, err = ioutil.ReadFile(hmacKeyPath)
	if err != nil {
		// hmac doesn't exist, let's create one!
		secret, err := uuid.NewV4()
		if err != nil {
			panic(fmt.Sprintf("Couldn't generate hmac key: %v", err))
		}

		h := hmac.New(sha256.New, []byte(secret.String()))

		// create the file to write to
		file, err := os.Create(hmacKeyPath)
		if err != nil {
			panic(fmt.Sprintf("Couldn't open hmac key file for writing: %v", err))
		}
		defer file.Close() // definitely make sure to close it

		// write the hash to the file
		if _, err := io.Copy(h, file); err != nil {
			panic(fmt.Sprintf("Couldn't write hmac key to file: %v", err))
		}

		// try reading again, if not, we fucked up
		hmacKey, err = ioutil.ReadFile(hmacKeyPath)
		if err != nil {
			panic("couldn't read hmac file")
		}
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

// Auth is the authentication middleware
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := db.GetDB(c)
		var user users.User
		var err error

		tokn, okTok := c.Request.Header["X-Auth-Token"]

		if okTok {
			user, err = userWithToken(tokn[0], db)
			if err != nil {
				c.AbortWithError(http.StatusUnauthorized, err)
			}
		} else {
			if key, ok := c.GetQuery("key"); allowedWithKey(c) && ok {
				user, err = userWithKey(key, db)
				if err != nil {
					c.AbortWithError(http.StatusUnauthorized, err)
					return
				}
			} else {
				c.AbortWithError(http.StatusUnauthorized, ErrNoToken)
				return
			}
		}

		err = user.UpdateAPIKey()
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		if err := db.Save(&user).Error; err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		c.Set("user", user)
		c.Next() // continue on to next endpoint
	}
}

// Login a user
func Login(user *users.User, c *gin.Context) (string, error) {
	if err := user.UpdateAPIKey(); err != nil {
		if err = user.GenAPIKey(); err != nil {
			return "", err
		}
	}

	token := jwt.New(jwt.SigningMethodHS256)

	// Set some claims
	token.Claims["api_key"] = user.APIKey

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(hmacKey)

	return tokenString, err
}

// Logout a user
func Logout(u *users.User, c *gin.Context) {
	u.DestroyAPIKey()
	c.Status(http.StatusOK)
}

// CurrentUser gets the current user from the session
func CurrentUser(c *gin.Context) (u users.User, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ErrNoUser
		}
	}()

	user, ok := c.Get("user")
	if !ok {
		return users.User{}, ErrNoUser
	}

	u = user.(users.User)

	return
}

func userWithToken(tokn string, db *gorm.DB) (users.User, error) {
	token, err := jwt.Parse(tokn, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return hmacKey, nil
	})

	if err != nil {
		return users.User{}, err
	}

	key, err := uuid.ParseHex(token.Claims["api_key"].(string))
	if err != nil {
		return users.User{}, err
	}
	apiKey := key.String()

	// Check this is the right user with correct API key

	user, err := users.GetByAPIKey(apiKey, db)
	if err != nil {
		return users.User{}, err
	}

	return user, nil
}

func userWithKey(key string, db *gorm.DB) (users.User, error) {
	ticket, err := tickets.GetByID(key, db)
	if err != nil {
		return users.User{}, err
	}
	return users.GetByID(ticket.UserID, db)
}

func allowedWithKey(c *gin.Context) bool {
	allowed := []string{
		"partisan/api/v2.NotificationsCount",
		"partisan/api/v2.MessageSocket",
	}
	for _, v := range allowed {
		if c.HandlerName() == v {
			return true
		}
	}

	return false
}
