package messages

import "errors"

var (
	ErrMessageSelf          = errors.New("You can't message yourself. Get some friends")
	ErrNoUsers              = errors.New("No Users were found for this thread")
	ErrThreadsNotFound      = errors.New("No threads were found")
	ErrThreadUsersNotFound  = errors.New("No thread users were found")
	ErrThreadUnreciprocated = errors.New("This thread is unreciprocated")
)
