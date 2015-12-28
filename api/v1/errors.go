package v1

import (
	"net/http"
	"partisan/auth"
	"partisan/dao"
	"partisan/matcher"
	"partisan/questions"

	"partisan/Godeps/_workspace/src/github.com/gin-gonic/gin"
)

func handleError(err error, c *gin.Context) {
	switch err.(type) {
	case *ErrBinding:
		c.AbortWithError(http.StatusBadRequest, err)
	case *ErrDBInsert:
		c.AbortWithError(http.StatusNotAcceptable, err)
	case *ErrDBDelete:
		c.AbortWithError(http.StatusNotAcceptable, err)
	case *ErrDBNotFound:
		c.AbortWithError(http.StatusNotFound, err)
	case *ErrParseID:
		c.AbortWithError(http.StatusNotAcceptable, err)
	case *ErrUserNotFound, *ErrPasswordMatch:
		c.AbortWithError(http.StatusUnauthorized, err)
	case *auth.ErrNoUser:
		c.AbortWithError(http.StatusUnauthorized, err)
	case *dao.ErrNotFound:
		c.AbortWithError(http.StatusNotFound, err)
	case *dao.ErrNoFriendship:
		c.AbortWithError(http.StatusNotFound, err)
	case *dao.ErrThreadNotFound:
		c.AbortWithError(http.StatusNotFound, err)
	case *matcher.ErrOutOfRange:
		c.AbortWithError(http.StatusNotAcceptable, err)
	case *questions.ErrNoneValid:
		c.AbortWithError(http.StatusNotFound, err)
	default:
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}

type ErrBasic struct {
	Err error
}

func (e *ErrBasic) Error() string {
	if e == nil {
		return ""
	}

	return e.Err.Error()
}

type ErrBinding struct {
	Err error
}

func (e *ErrBinding) Error() string {
	if e == nil {
		return ""
	}

	return e.Err.Error()
}

type ErrDBInsert struct {
	Err error
}

func (e *ErrDBInsert) Error() string {
	if e == nil {
		return ""
	}

	return e.Err.Error()
}

type ErrDBDelete struct {
	Err error
}

func (e *ErrDBDelete) Error() string {
	if e == nil {
		return ""
	}

	return e.Err.Error()
}

type ErrDBNotFound struct {
	Err error
}

func (e *ErrDBNotFound) Error() string {
	if e == nil {
		return ""
	}

	return e.Err.Error()
}

type ErrParseID struct {
	Err error
}

func (e *ErrParseID) Error() string {
	if e == nil {
		return ""
	}

	return e.Err.Error()
}

type ErrNoFile struct {
	Err error
}

func (e *ErrNoFile) Error() string {
	if e == nil {
		return ""
	}

	return e.Err.Error()
}

type ErrUserNotFound struct{}

func (e *ErrUserNotFound) Error() string {
	return "User not found"
}

type ErrPasswordMatch struct{}

func (e *ErrPasswordMatch) Error() string {
	return "Password didn't match"
}

type ErrNoUserID struct{}

func (e *ErrNoUserID) Error() string {
	return "No User ID specified"
}

type ErrNoThreadID struct{}

func (e *ErrNoThreadID) Error() string {
	return "No Thread ID specified"
}
