package main

import (
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
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

var store = sessions.NewCookieStore([]byte("aoishoi1293220hdacns92309")) // TODO: get better random string

// Auth is the authentication middleware
func Auth() gin.HandlerFunc {
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
				c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("No session"))
				return
			}

			userID, _ = strconv.ParseUint(fmt.Sprintf("%d", sess.Get("user_id")), 10, 0)
			// apiKey = fmt.Sprintf("%s", sess.Get("api_key"))
		}

		// Check this is the right user with correct API key
		db, err := initDB()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		defer db.Close()

		user := User{}

		// look up by user id and api key, if you can't find it, the api key is no good
		if err := db.Where(User{ID: userID}).First(&user).Error; err != nil {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Couldn't find user"))
			return
		}

		// if user.APIKeyExp.Before(time.Now()) {
		// 	c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("API key is expired"))
		// 	return
		// }

		c.Set("user", user)
		c.Next() // continue on to next endpoint
	}
}

// LoginHandler Handle logging in
func LoginHandler(c *gin.Context) {
	// There's some potential for functions in here to panic, so I'm trying to recover
	defer func() {
		if r := recover(); r != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("%v", r))
		}
	}()

	// init db
	db, err := initDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	email := c.Request.PostFormValue("email")
	password := c.Request.PostFormValue("password")

	user := User{}

	if err := db.Where(User{Email: email}).First(&user).Error; err != nil {
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

	tokenString, _ := login(user, c)

	// User.APIKey is not exported via JSON so we have to manually add it
	c.JSON(http.StatusOK, gin.H{"token": tokenString, "user": user, "api_key": user.APIKey})
}

// LogoutHandler logs the user out (duh)
func LogoutHandler(c *gin.Context) {
	logout(c)

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func login(user User, c *gin.Context) (string, error) {
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

func logout(c *gin.Context) {
	sess := sessions.Default(c)
	sess.Clear()
	sess.Save()
}
