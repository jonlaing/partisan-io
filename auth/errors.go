package auth

import "errors"

var (
	ErrNoToken = errors.New("No X-AUTH-TOKEN header was found")
)
