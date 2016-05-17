package v2

import "errors"

var (
	ErrBinding = errors.New("Couldn't bind JSON to model")
)
