package auth

import "errors"

var (
	ErrNoAppToken      = errors.New("No X-APP-TOKEN was found. This app is not authorized to use Partisan API")
	ErrAppTokenRevoked = errors.New("This App Token has been revoked.")
	ErrNoToken         = errors.New("No X-AUTH-TOKEN header was found")
)
