package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"partisan/db"
	"partisan/matcher"
	m "partisan/models"
	"regexp"
	"time"
	"partisan/auth"
)

// UserCreate is the sign up route
func UserCreate(c *gin.Context) {
	db, err := db.InitDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	validationErrs := make(map[string]string)

	var user m.User
	if err := c.Bind(&user); err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

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

	if err := user.GetLocation(); err != nil {
		validationErrs["postal_code"] = fmt.Sprintf("Error validating postal code. %s", err.Error())
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

	profile := m.Profile{
		UserID: user.ID,
	}
	if err := db.Create(&profile).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	auth.Login(user, c)

	c.JSON(http.StatusCreated, user)
}

// UserShow shows shit about the current user
func UserShow(c *gin.Context) {
	db, err := db.InitDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	currentUser, err := auth.CurrentUser(c, &db)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	c.JSON(http.StatusOK, currentUser)
}

// UserMatch will return the match % of the signed in user, and the user in the path
func UserMatch(c *gin.Context) {
	db, err := db.InitDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	currentUser, err := auth.CurrentUser(c, &db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	mUser := m.User{} // User to match
	userID := c.Params.ByName("user_id")

	if err := db.First(&mUser, userID).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	match, err := matcher.Match(currentUser.PoliticalMap, mUser.PoliticalMap)
	if err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
	}

	c.JSON(http.StatusOK, gin.H{"match": match})
}

//////////////////////////////////////////

