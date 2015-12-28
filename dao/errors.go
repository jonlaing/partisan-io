package dao

import "fmt"

type ErrBasic struct {
	Err error
}

func (e *ErrBasic) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}

	return ""
}

type ErrNotFound struct {
	Err error
}

func (e *ErrNotFound) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}

	return ""
}

type ErrNoFriendship struct {
	UserID   uint64
	FriendID uint64
}

func (e *ErrNoFriendship) Error() string {
	return fmt.Sprintf("Couldn't find friendship between User: %d and Friend: %d", e.UserID, e.FriendID)
}

type ErrNoMatches struct {
	Err error
}

func (e *ErrNoMatches) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}

	return ""
}

type ErrThreadNotFound struct {
	Err error
}

func (e *ErrThreadNotFound) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}

	return ""
}
