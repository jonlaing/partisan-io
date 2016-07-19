package v2

import (
	"net/http"
	"partisan/auth"
	"partisan/db"

	"partisan/models.v2/answers"

	"github.com/gin-gonic/gin"
)

// AnswersUpdate updates the coordinates of user based on question answers
func AnswersUpdate(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	var a answers.Answer
	if err := c.BindJSON(&a); err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	if len(a.Map) == 0 {
		c.AbortWithError(http.StatusNotAcceptable, answers.ErrMap)
		return
	}

	dx, dy, err = user.PoliticalMap.Add(a.Map, a.Mask, a.Agree)
	if err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}
	x, y := user.PoliticalMap.Center()

	user.CenterX = x
	user.CenterY = y
	user.DeltaX = dx
	user.DeltaY = dy

	if err := db.Save(&user).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}
