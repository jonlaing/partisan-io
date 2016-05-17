package users

import "errors"

var (
	ErrPasswordConfirmMatch = errors.New("Password and Password Confirm don't match. Try retyping.")
	ErrEmailValidation      = errors.New("Email looks malformed. Check for typos.")
	ErrUsernameValidation   = errors.New("Username can only have letters, numbers, dashes and underscores, and can't be longer than 16 characters. Ex: my_username123")
	ErrNoAPIKey             = errors.New("There is no APIKey for this user")
	ErrAPIKeyExpired        = errors.New("This APIKey has expired")
)
