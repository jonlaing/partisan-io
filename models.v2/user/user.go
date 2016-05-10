package user

import (
	"partisan/matcher"
	"regexp"
	"time"

	"github.com/jasonmoo/geo"
	"github.com/nu7hatch/gouuid"

	models "partisan/models.v2"

	"golang.org/x/crypto/bcrypt"
)

// User the user model
type User struct {
	ID                 uint64               `json:"id" gorm:"primary_key"`
	Username           string               `json:"username" sql:"not null;unique_index" binding:"required"`
	Email              string               `json:"email" sql:"not null;unique_index" binding:"required"`
	Gender             string               `json:"gender"`
	Birthdate          time.Time            `json:"birthdate"`
	AvatarURL          string               `json:"avatar_url"`
	AvatarThumbnailURL string               `json:"avatar_thumbnail_url"`
	PoliticalMap       matcher.PoliticalMap `json:"-" sql:"type:varchar(255)"`
	CenterX            int                  `json:"-"`
	CenterY            int                  `json:"-"`
	PostalCode         string               `json:"-"`
	Location           string               `json:"location"`
	Longitude          float64              `json:"-"`
	Latitude           float64              `json:"-"`
	CreatedAt          time.Time            `json:"created_at"`
	UpdatedAt          time.Time            `json:"updated_at"`
	APIKey             string               `json:"-"`
	APIKeyExp          time.Time            `json:"-"`
	PasswordHash       []byte               `json:"-"`
	Type               int                  `json:"-"`
}

// CreatorBinding is a struct for fields neeeded to create a user via JSON Binding
type CreatorBinding struct {
	Username        string `json:"username" binding:"required"`
	Email           string `json:"email" binding:"required"`
	PostalCode      string `json:"postal_code"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"-"`
}

// UpdaterBinding is a struct for fields neeeded to update a user via JSON Binding
type UpdaterBinding struct {
	Gender     string    `json:"gender"`
	Birthdate  time.Time `json:"birthdate"`
	PostalCode string    `json:"-"`
}

// New initializes a new User based on the CreatorBinding. Performs Validation and generates location and password.
func New(b CreatorBinding) (u User, errs models.ValidationErrors) {
	errs = make(models.ValidationErrors)

	u.Username = b.Username
	u.Email = b.Email
	u.PostalCode = b.PostalCode
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

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

	u.Gender = b.Gender
	u.Birthdate = b.Birthdate
	u.PostalCode = b.PostalCode
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
	if !emailRegex.MatchString(u.Email) {
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
	if u.APIKey == "" {
		return ErrNoAPIKey
	}

	// If more than a week has past since the User's last activity, then
	// this key is expired
	if time.Now().After(u.APIKeyExp) {
		return ErrAPIKeyExpired
	}

	return nil
}
