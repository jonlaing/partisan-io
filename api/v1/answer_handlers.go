package v1

import (
	"errors"
	"net/http"
	"partisan/auth"
	"partisan/db"
	m "partisan/models"

	"github.com/gin-gonic/gin"
)

// AnswersUpdate updates the coordinates of user based on question answers
func AnswersUpdate(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		handleError(err, c)
		return
	}

	a := m.Answer{}

	if err := c.BindJSON(&a); err != nil {
		handleError(&ErrBinding{err}, c)
		return
	}

	if len(a.Map) == 0 {
		handleError(&ErrBinding{errors.New("Answer doesn't have map. Probably an error in binding")}, c)
		return
	}

	err = user.PoliticalMap.Add(a.Map, a.Agree)
	if err != nil {
		handleError(err, c)
		return
	}

	x, y := user.PoliticalMap.Center()
	user.CenterX = x
	user.CenterY = y

	if err := db.Save(&user).Error; err != nil {
		handleError(&ErrDBInsert{err}, c)
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}
