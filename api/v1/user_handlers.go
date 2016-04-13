package v1

import (
	"fmt"
	"net/http"
	"partisan/auth"
	"partisan/dao"
	"partisan/db"
	"partisan/imager"
	"partisan/logger"
	"partisan/matcher"
	m "partisan/models"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// UserCreate is the sign up route
func UserCreate(c *gin.Context) {
	db := db.GetDB(c)

	validationErrs := make(map[string]string)

	var user m.User
	if err := c.Bind(&user); err != nil {
		handleError(&ErrBinding{err}, c)
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

	// check for non unique user
	var uniqueUserTest m.User
	db.Where("email = ?", user.Email).Or("username = ?", user.Username).Find(&uniqueUserTest)
	if uniqueUserTest.Email == user.Email {
		validationErrs["email"] = "Email is already in use"
	}

	if uniqueUserTest.Username == user.Username {
		validationErrs["username"] = "Username already in use"
	}

	if len(validationErrs) > 0 {
		c.JSON(http.StatusNotAcceptable, validationErrs)
		return
	}

	// Validation checks out, create the user
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		handleError(err, c)
		return
	}

	user.PasswordHash = hash

	if err := db.Create(&user).Error; err != nil {
		handleError(&ErrDBInsert{err}, c)
		return
	}

	profile := m.Profile{
		UserID: user.ID,
	}
	if err := db.Create(&profile).Error; err != nil {
		handleError(&ErrDBInsert{err}, c)
		return
	}

	// if err := emailer.SendWelcomeEmail(user.Username, user.Email); err != nil {
	// 	logger.Error.Println(err)
	// }

	logger.Warning.Println("EMAILS ARE NOT SENDING")

	token, _ := auth.Login(user, c)

	c.JSON(http.StatusCreated, gin.H{"token": token, "user": user, "api_key": user.APIKey})
}

// UserShow shows shit about the current user
func UserShow(c *gin.Context) {
	currentUser, err := auth.CurrentUser(c)
	if err != nil {
		handleError(err, c)
		return
	}

	c.JSON(http.StatusOK, currentUser)
}

// UserUpdate will update the user
func UserUpdate(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		handleError(err, c)
		return
	}

	user.Gender = c.DefaultPostForm("gender", user.Gender)

	birthdate := c.PostForm("birthdate")
	if birthdate != "" {
		user.Birthdate, err = time.Parse("2006-01-02", birthdate)
		if err != nil {
			logger.Error.Println(err)
		}
	}

	postalCode := c.PostForm("postal_code")
	if postalCode != "" {
		user.PostalCode = postalCode
		user.GetLocation()
	}

	if err := db.Save(&user).Error; err != nil {
		handleError(&ErrDBInsert{err}, c)
		return
	}

	c.JSON(http.StatusOK, user)
}

// UserMatch will return the match % of the signed in user, and the user in the path
func UserMatch(c *gin.Context) {
	db := db.GetDB(c)

	currentUser, err := auth.CurrentUser(c)
	if err != nil {
		handleError(err, c)
		return
	}

	mUser := m.User{} // User to match
	userID := c.Params.ByName("user_id")

	if err := db.First(&mUser, userID).Error; err != nil {
		handleError(&ErrDBInsert{err}, c)
		return
	}

	match, _ := matcher.Match(currentUser.PoliticalMap, mUser.PoliticalMap)

	c.JSON(http.StatusOK, gin.H{"match": match})
}

// UserAvatarUpload handles uploading a user's avatar
func UserAvatarUpload(c *gin.Context) {
	db := db.GetDB(c)

	currentUser, err := auth.CurrentUser(c)
	if err != nil {
		handleError(err, c)
		return
	}

	tmpFile, _, err := c.Request.FormFile("avatar")
	if err != nil {
		handleError(&ErrNoFile{err}, c)
		return
	}
	defer tmpFile.Close()

	processor := imager.ImageProcessor{File: tmpFile}

	// channels to process images on different threads
	var fullPath string
	var thumbPath string

	if err := processor.Resize(1500); err != nil {
		handleError(err, c)
		return
	}
	fullPath, err = processor.Save("/localfiles/img")
	if err != nil {
		handleError(err, c)
		return
	}

	// Save the thumbnail
	if err := processor.Thumbnail(250); err != nil {
		handleError(err, c)
		return
	}
	thumbPath, err = processor.Save("/localfiles/img/thumb")
	if err != nil {
		handleError(err, c)
		return
	}

	currentUser.AvatarURL = fullPath
	currentUser.AvatarThumbnailURL = thumbPath

	if err := db.Save(&currentUser).Error; err != nil {
		handleError(&ErrDBInsert{err}, c)
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

	u, err := auth.CurrentUser(c)
	if err != nil {
		handleError(err, c)
		return
	}

	username := c.Query("tag")

	friendIDs, err := dao.FriendIDs(u, true, db)
	if err != nil || len(friendIDs) == 0 {
		handleError(&ErrDBNotFound{err}, c)
		return
	}

	var suggestions []string
	db.Model(m.User{}).Where("id IN (?) AND username LIKE ?", friendIDs, "%"+username+"%").Pluck("username", &suggestions)

	c.JSON(http.StatusOK, gin.H{"suggestions": suggestions})
}
