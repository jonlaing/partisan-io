package users

import (
	"mime/multipart"
	"os"
	"partisan/imager"
	"partisan/location"
	"partisan/matcher"
	"regexp"
	"time"

	"github.com/jasonmoo/geo"
	"github.com/nu7hatch/gouuid"

	models "partisan/models.v2"

	"golang.org/x/crypto/bcrypt"
)

// LookingFor is a bitfield for determining what sort of
// interaction the user is looking for
type LookingFor int

const (
	// LNone is the zero value
	LNone LookingFor = 0
	// LFriends - the user is looking for friends
	LFriends LookingFor = 1 << iota
	// LLove - the user is looking for love
	LLove
	// LEnemies - the user is looking for enemies
	LEnemies
)

// User the user model
type User struct {
	ID                 string               `json:"id" gorm:"primary_key" sql:"type:uuid;default:uuid_generate_v4()"`
	Username           string               `json:"username" sql:"not null;unique_index" binding:"required"`
	Email              string               `json:"email" sql:"not null;unique_index" binding:"required"`
	Gender             string               `json:"gender"`
	Birthdate          time.Time            `json:"birthdate"`
	AvatarURL          string               `json:"avatar_url"`
	AvatarThumbnailURL string               `json:"avatar_thumbnail_url"`
	PoliticalMap       matcher.PoliticalMap `json:"-" sql:"type:varchar(255)"`
	CenterX            int                  `json:"-"`
	CenterY            int                  `json:"-"`
	DeltaX             int                  `json:"-"` // For choosing new question sets
	DeltaY             int                  `json:"-"` // For choosing new question sets
	PostalCode         string               `json:"postal_code"`
	Location           string               `json:"location"`
	Longitude          float64              `json:"-"`
	Latitude           float64              `json:"-"`
	LookingFor         LookingFor           `json:"looking_for"`
	Summary            string               `json:"summary"`
	CreatedAt          time.Time            `json:"created_at"`
	UpdatedAt          time.Time            `json:"updated_at"`
	APIKey             string               `json:"-" sql:"type:uuid;default:uuid_generate_v4()"`
	APIKeyExp          time.Time            `json:"-"`
	PasswordHash       []byte               `json:"-"`
	DeviceToken        string               `json:"-"`
}

// CreatorBinding is a struct for fields neeeded to create a user via JSON Binding
type CreatorBinding struct {
	Username        string `json:"username" binding:"required"`
	Email           string `json:"email" binding:"required"`
	PostalCode      string `json:"postal_code" binding:"required"`
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"password_confirm" binding:"required"`
	DeviceToken     string `json:"device_token"`
}

// UpdaterBinding is a struct for fields neeeded to update a user via JSON Binding
type UpdaterBinding struct {
	Gender     string     `json:"gender"`
	Birthdate  string     `json:"birthdate"`
	LookingFor LookingFor `json:"looking_for"`
	Summary    string     `json:"summary"`
	PostalCode string     `json:"postal_code"`
}

// New initializes a new User based on the CreatorBinding. Performs Validation and generates location and password.
func New(b CreatorBinding) (u User, errs models.ValidationErrors) {
	errs = make(models.ValidationErrors)

	u.Username = b.Username
	u.Email = b.Email
	u.PostalCode = b.PostalCode
	u.CreatedAt = time.Now()
	u.UpdatedAt = u.CreatedAt
	u.DeviceToken = b.DeviceToken

	if err := u.GeneratePasswordHash(b.Password, b.PasswordConfirm); err != nil {
		errs["password"] = err
	}

	if verrs := u.Validate(); len(verrs) > 0 {
		for k, e := range verrs {
			errs[k] = e
		}
	}

	if err := u.GetLocation(); err != nil {
		errs["postal_code"] = err
	}

	return
}

// Update updates a User based on the UpdaterBinding. Performs Validation and generates location. Does not save to the db.
func (u *User) Update(b UpdaterBinding) (errs models.ValidationErrors) {
	errs = make(models.ValidationErrors)

	if b.Gender != "" {
		u.Gender = b.Gender
	}

	if b.Birthdate != "" {
		var err error
		u.Birthdate, err = time.Parse("2006-01-02", b.Birthdate)
		if err != nil {
			errs["birthdate"] = err
		}
	}

	if b.PostalCode != "" {
		u.PostalCode = b.PostalCode
	}

	if b.LookingFor != LNone {
		u.LookingFor = b.LookingFor
	}

	if b.Summary != "" {
		u.Summary = b.Summary
	}

	u.UpdatedAt = time.Now()

	if err := u.GetLocation(); err != nil {
		errs["postal_code"] = err
	}

	return
}

// GeneratePasswordHash generates a new password hash from the provided password and confirmation
func (u *User) GeneratePasswordHash(password, confirm string) error {
	if password != confirm {
		return ErrPasswordConfirmMatch
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	u.PasswordHash = hash
	return err
}

// Validate validates the user to ensure it meets the requirements of the system
func (u *User) Validate() (errs models.ValidationErrors) {
	errs = make(models.ValidationErrors)

	emailRegex := regexp.MustCompile("(?i)[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?")
	if len(u.Email) > 255 || !emailRegex.MatchString(u.Email) {
		errs["email"] = ErrEmailValidation
	}

	usernameRegex := regexp.MustCompile("(?i)^[a-z0-9-_]+$")
	if len(u.Username) < 3 || len(u.Username) > 16 || !usernameRegex.MatchString(u.Username) {
		errs["username"] = ErrUsernameValidation
	}

	return errs
}

// GetLocation finds the latitude/longitude by postal code
func (u *User) GetLocation() error {
	location, err := geo.Geocode(u.PostalCode)
	if err != nil {
		return err
	}

	u.Location = location.Address
	u.Latitude = location.Lat
	u.Longitude = location.Lng

	return nil
}

// GenAPIKey generates a new APIKey and expiration date one week from now
func (u *User) GenAPIKey() error {
	k, err := uuid.NewV4()
	u.APIKey = k.String()
	u.APIKeyExp = time.Now().Add(168 * time.Hour)
	return err
}

// UpdateAPIKey is meant to be used when a user makes a request, to update the expiration
// of their APIKey. An APIKey will expire after 1 week of inactivity.
func (u *User) UpdateAPIKey() error {
	err := u.ValidateAPIKey() // check if the current APIKey is valid
	u.APIKeyExp = time.Now().Add(168 * time.Hour)
	return err
}

// ValidateAPIKey checks if an APIKey exists, and if it's expired
func (u *User) ValidateAPIKey() error {
	if u.APIKey == new(uuid.UUID).String() || u.APIKey == "" {
		return ErrNoAPIKey
	}

	// If more than a week has passed since the User's last activity, then
	// this key is expired
	if time.Now().After(u.APIKeyExp) {
		return ErrAPIKeyExpired
	}

	return nil
}

func (u *User) DestroyAPIKey() {
	u.APIKey = new(uuid.UUID).String()
	u.APIKeyExp = time.Time{}
}

// CheckPassword checks the provided password with the password hash, returns an error if they don't match
func (u User) CheckPassword(pw string) error {
	return bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(pw))
}

func (u *User) AttachAvatar(f multipart.File) error {
	var err error
	var fullPath, thumbPath string
	isS3 := false

	if len(os.Getenv("AWS_ACCESS_KEY_ID")) > 0 {
		isS3 = true
	}

	processor := imager.ImageProcessor{File: f}

	if err := processor.Resize(1080); err != nil {
		return err
	}

	if isS3 {
		fullPath, err = processor.Save("/img")
	} else {
		fullPath, err = processor.Save("/localfiles/img")
	}
	if err != nil {
		return err
	}

	// Save the thumbnail
	if err := processor.Thumbnail(250); err != nil {
		return err
	}

	if isS3 {
		thumbPath, err = processor.Save("/img/thumb")
	} else {
		thumbPath, err = processor.Save("/localfiles/img/thumb")
	}
	if err != nil {
		return err
	}

	u.AvatarURL = fullPath
	u.AvatarThumbnailURL = thumbPath

	return nil
}

func (u User) GetID() string {
	return u.ID
}

func (u User) GetSubscriberType() string {
	return "user"
}

func (u User) GetPoliticalMap() matcher.PoliticalMap {
	return u.PoliticalMap
}

func (u User) Distance(lat, long float64) float64 {
	p1 := location.NewPoint(u.Latitude, u.Longitude)
	p2 := location.NewPoint(lat, long)

	return location.Distance(p1, p2)
}
