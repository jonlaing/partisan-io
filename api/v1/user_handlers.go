package v1

import (
	"crypto/rand"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"golang.org/x/crypto/bcrypt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"partisan/auth"
	"partisan/db"
	"partisan/matcher"
	m "partisan/models"
	"regexp"
	"time"
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

// UserAvatarUpload handles uploading a user's avatar
func UserAvatarUpload(c *gin.Context) {
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

	tmpFile, _, err := c.Request.FormFile("avatar")
	if err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}
	defer tmpFile.Close()

	// DETECT FILE TYPE
	buf := make([]byte, 512) // why 512 bytes ? see http://golang.org/pkg/net/http/#DetectContentType
	_, err = tmpFile.Read(buf)
	if err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	filetype := http.DetectContentType(buf)

	// Resetting read of tmpFile (otherwise we'd copy an incomplete file)
	_, err = tmpFile.Seek(0, 0)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tmpImg, _, err := image.Decode(tmpFile)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var avatar image.Image
	bounds := tmpImg.Bounds().Max
	fmt.Println(bounds)
	if bounds.X > bounds.Y {
		ratio := float64(bounds.X) / float64(bounds.Y)
		avatar = resize.Resize(uint(ratio*float64(250)), 250, tmpImg, resize.Bicubic)
	} else if bounds.Y > bounds.X {
		ratio := float64(bounds.Y) / float64(bounds.X)
		avatar = resize.Resize(250, uint(ratio*float64(250)), tmpImg, resize.Bicubic)
	} else {
		avatar = resize.Thumbnail(250, 250, tmpImg, resize.Bicubic)
	}

	var file *os.File
	var filename string

	switch filetype {
	case "image/jpeg", "image/jpg":
		filename = "/localfiles/img/" + imgFilename() + ".jpg"
		file, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		defer file.Close()
		jpeg.Encode(file, avatar, &jpeg.Options{
			Quality: 100,
		})
		break
	case "image/png":
		filename = "./localfiles/img/" + imgFilename() + ".png"
		file, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		defer file.Close()
		png.Encode(file, avatar)
		break
	case "image/gif":
		// NOT SUPPORTING GIF YET
		fallthrough
	default:
		c.AbortWithError(http.StatusNotAcceptable, fmt.Errorf("Unacceptable file format: %s", filetype))
		return
		break
	}

	currentUser.AvatarURL = filename

	if err := db.Save(&currentUser).Error; err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, currentUser)
}

// generate a random filename
func imgFilename() string {
	dictionary := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, 16)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}
