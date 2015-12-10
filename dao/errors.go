package dao

import "fmt"

type ErrBasic struct {
	Err error
}

func (e *ErrBasic) Error() string {
	return e.Err.Error()
}

type ErrNotFound struct {
	ErrBasic
}

type ErrNoFriendship struct {
	UserID   uint64
	FriendID uint64
}

func (e *ErrNoFriendship) Error() string {
	return fmt.Sprintf("Couldn't find friendship between User: %d and Friend: %d", e.UserID, e.FriendID)
}

type ErrNoMatches struct {
	ErrBasic
}

type ErrThreadNotFound struct {
	ErrBasic
}
