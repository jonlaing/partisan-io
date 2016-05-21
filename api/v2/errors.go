package v2

import "errors"

var (
	ErrBinding       = errors.New("Couldn't bind JSON to model")
	ErrCannotUpdate  = errors.New("You do not have the permissions to update this record")
	ErrCannotDelete  = errors.New("You do not have the permissions to delete this record")
	ErrAlreadyExists = errors.New("Cannot create resource as it already exists")
)
