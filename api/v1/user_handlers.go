package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"partisan/auth"
	"partisan/db"
	"partisan/emailer"
	"partisan/imager"
	"partisan/matcher"
	m "partisan/models"
	"regexp"
	"time"
)

// UserCreate is the sign up route
func UserCreate(c *gin.Context) {
	db := db.GetDB(c)

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

	if err := emailer.SendWelcomeEmail(user.Username, user.Email); err != nil {
		fmt.Println(err)
	}

	auth.Login(user, c)

	c.JSON(http.StatusCreated, user)
}

// UserShow shows shit about the current user
func UserShow(c *gin.Context) {
	currentUser, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	c.JSON(http.StatusOK, currentUser)
}

// UserUpdate will update the user
func UserUpdate(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	user.Gender = c.DefaultPostForm("gender", user.Gender)

	postalCode := c.PostForm("postal_code")
	if postalCode != "" {
		user.PostalCode = postalCode
		user.GetLocation()
	}

	if err := db.Save(&user).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// UserMatch will return the match % of the signed in user, and the user in the path
func UserMatch(c *gin.Context) {
	db := db.GetDB(c)

	currentUser, err := auth.CurrentUser(c)
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

// UserAvatarUpload handles uploading a user's avatar
func UserAvatarUpload(c *gin.Context) {
	db := db.GetDB(c)

	currentUser, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	tmpFile, _, err := c.Request.FormFile("avatar")
	if err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}
	defer tmpFile.Close()

	processor := imager.ImageProcessor{File: tmpFile}

	// Save the full-size
	var fullPath string
	if err := processor.Resize(1500); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	fullPath, err = processor.Save("/localfiles/img")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	currentUser.AvatarURL = fullPath

	// Save the thumbnail
	var thumbPath string
	if err := processor.Thumbnail(250); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	thumbPath, err = processor.Save("/localfiles/img/thumb")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	currentUser.AvatarThumbnailURL = thumbPath

	if err := db.Save(&currentUser).Error; err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, currentUser)
}

// UserCheckUnique checks a username for uniqueness
func UserCheckUnique(c *gin.Context) {
	db := db.GetDB(c)

	username := c.Query("username")

	var count int
	if err := db.Model(m.User{}).Where("username = ?", username).Count(&count).Error; err == nil {
		if count > 0 {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{})
}

// UsernameSuggest returns a short list of possible friends based on a tag
func UsernameSuggest(c *gin.Context) {
	db := db.GetDB(c)

	username := c.Query("tag")

	var suggestions []string
	db.Model(m.User{}).Where("username LIKE ?", "%"+username+"%").Pluck("username", &suggestions)

	c.JSON(http.StatusOK, gin.H{"suggestions": suggestions})
}
