package v1

import (
	"net/http"
	"partisan/auth"
	"partisan/dao"
	"partisan/matcher"
	"partisan/questions"

	"github.com/gin-gonic/gin"
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
	return e.Err.Error()
}

type ErrBinding struct {
	ErrBasic
}

type ErrDBInsert struct {
	ErrBasic
}

type ErrDBDelete struct {
	ErrBasic
}

type ErrDBNotFound struct {
	ErrBasic
}

type ErrParseID struct {
	ErrBasic
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

type ErrNoFile struct {
	ErrBasic
}
