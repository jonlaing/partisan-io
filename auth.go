package main

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
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

// Auth is the authentication middleware
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokn, ok := c.Request.Header["X-Auth-Token"]
		if !ok {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("No X-Auth-Token"))
			return
		}

		token, err := jwt.Parse(tokn[0], func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return hmacKey, nil
		})

		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		// Check this is the right user with correct API key
		db, err := initDB()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		defer db.Close()

		// try to convert claims to respective types
		userID, _ := strconv.ParseUint(fmt.Sprintf("%d", token.Claims["user_id"]), 10, 0)
		apiKey := fmt.Sprintf("%s", token.Claims["api_key"])

		user := User{}

		// look up by user id and api key, if you can't find it, the api key is no good
		if err := db.Where(User{ID: userID, APIKey: apiKey}).First(&user).Error; err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		if user.APIKeyExp.Before(time.Now()) {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("API key is expired"))
			return
		}

		c.Set("user_id", user.ID)
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

	var login LoginJSON
	c.BindJSON(&login)

	user := User{}

	if err := db.Where(User{Email: login.Email}).First(&user).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(login.Password)); err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	user.APIKey = uuid.NewRandom().String()
	user.APIKeyExp = time.Now().Add(time.Hour * 72)

	db.Save(&user)

	token := jwt.New(jwt.SigningMethodHS256)

	// Set some claims
	token.Claims["user_id"] = user.ID
	token.Claims["api_key"] = user.APIKey
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(hmacKey)
	if err != nil {
		fmt.Println("this is where i fucked up")
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	// User.APIKey is not exported via JSON so we have to manually add it
	c.JSON(http.StatusOK, gin.H{"token": tokenString, "user": user, "api_key": user.APIKey})
}
