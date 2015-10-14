package v1

import (
	"fmt"
	"net/http"
	"partisan/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"partisan/auth"
	"partisan/db"
	m "partisan/models"
)

// AnswersUpdate updates the coordinates of user based on question answers
func AnswersUpdate(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	a := m.Answer{}

	if err := c.BindJSON(&a); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if len(a.Map) == 0 {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Answer doesn't have map. Probably an error in binding"))
		return
	}

	err = user.PoliticalMap.Add(a.Map, a.Agree)
	if err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
	}

	x, y := user.PoliticalMap.Center()
	user.CenterX = x
	user.CenterY = y

	if err := db.Save(&user).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}
