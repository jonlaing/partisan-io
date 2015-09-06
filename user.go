package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
	"time"
)

// Record is an interace for everything that could be owned by a User
// This is used for permissions
type Record interface {
	User() uint64
}

// User the user model
type User struct {
	ID              uint64       `json:"id" gorm:"primary_key"`
	Username        string       `json:"username" sql:"not null,unique" binding:"required"`
	FullName        string       `json:"full_name" binding:"required"`
	Email           string       `json:"email" sql:"not null,unique" binding:"required"`
	AvatarURL       string       `json:"avatar_url"`
	PoliticalMap    PoliticalMap `json:"political_map"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
	APIKey          string       `json:"-"`
	APIKeyExp       time.Time    `json:"-"`
	PasswordHash    []byte       `json:"-"`
	Password        string       `json:"password" sql:"-" binding:"required"`
	PasswordConfirm string       `json:"password_confirm" sql:"-" binding:"required"`
}

// UserCreate is the sign up route
func UserCreate(c *gin.Context) {
	db, err := initDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	validationErrs := make(map[string]string)

	user := User{
		Email:           c.Request.PostFormValue("email"),
		Username:        c.Request.PostFormValue("username"),
		FullName:        c.Request.PostFormValue("full_name"),
		Password:        c.Request.PostFormValue("password"),
		PasswordConfirm: c.Request.PostFormValue("password_confirm"),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// VALIDATE
	if user.Password != user.PasswordConfirm {
		validationErrs["password_confirm"] = "Password and Password Confirm don't match. Try retyping."
	}

	emailRegex := regexp.MustCompile("(?i)[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?")
	if !emailRegex.MatchString(user.Email) {
		validationErrs["email"] = "Email looks malformed. Check for typos."
	}

	usernameRegex := regexp.MustCompile("(?i)^[a-z0-9-_]+$")
	if !usernameRegex.MatchString(user.Username) {
		validationErrs["username"] = "Username can only have letters, numbers, dashes and underscores. Ex: my_username123"
	}

	if len(validationErrs) > 0 {
		c.JSON(http.StatusNotAcceptable, validationErrs)
		return
	}

	// Validation checks out, create the user
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	user.PasswordHash = hash

	if err := db.Create(&user).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	login(user, c)

	c.JSON(http.StatusCreated, user)
}

// UserShow shows shit about the current user
func UserShow(c *gin.Context) {
	db, err := initDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	currentUser, err := CurrentUser(c, &db)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	c.JSON(http.StatusOK, currentUser)
}

// UserMatch will return the match % of the signed in user, and the user in the path
func UserMatch(c *gin.Context) {
	db, err := initDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	currentUser, err := CurrentUser(c, &db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	mUser := User{} // User to match
	userID := c.Params.ByName("user_id")

	if err := db.First(&mUser, userID).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	match, err := Match(currentUser.PoliticalMap, mUser.PoliticalMap)
	if err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
	}

	c.JSON(http.StatusOK, gin.H{"match": match})
}

//////////////////////////////////////////

// IsOwner checks if this User is the owner of the record
func (u User) IsOwner(r Record) bool {
	return u.ID == r.User()
}

// Friends returns all the User's friends
func (u User) Friends(db *gorm.DB) []User {
	friendIDs := u.FriendIDs(db)
	users := []User{}
	db.Where(friendIDs).Find(&users)
	return users
}

// FriendIDs returns all IDs of User's Friends
func (u User) FriendIDs(db *gorm.DB) []uint64 {
	var friendIDs []uint64
	db.Table("friendships").Select("friend_id").Where("confirmed = ?", true).Scan(&friendIDs)
	return friendIDs
}

// CurrentUser gets the current user from the session
func CurrentUser(c *gin.Context, db *gorm.DB) (User, error) {
	user, ok := c.Get("user")
	if !ok {
		return user.(User), fmt.Errorf("User ID not set")
	}

	return user.(User), nil
}
