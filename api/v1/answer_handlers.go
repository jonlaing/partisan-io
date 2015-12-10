package v1

import (
	"errors"
	"net/http"
	"partisan/auth"
	"partisan/db"
	m "partisan/models"

	"partisan/Godeps/_workspace/src/github.com/gin-gonic/gin"
)

// AnswersUpdate updates the coordinates of user based on question answers
func AnswersUpdate(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		return handleError(err, c)
	}

	a := m.Answer{}

	if err := c.BindJSON(&a); err != nil {
		return handleError(&ErrBinding{err}, c)
	}

	if len(a.Map) == 0 {
		return handleError(&ErrBinding{errors.New("Answer doesn't have map. Probably an error in binding")}, c)
	}

	err = user.PoliticalMap.Add(a.Map, a.Agree)
	if err != nil {
		return handleError(err, c)
	}

	x, y := user.PoliticalMap.Center()
	user.CenterX = x
	user.CenterY = y

	if err := db.Save(&user).Error; err != nil {
		handleError(&ErrDBInsert{err}, c)
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}
