package v2

import (
	"net/http"

	"partisan/auth"
	q "partisan/questions"

	"github.com/gin-gonic/gin"
)

// QuestionIndex finds a random QuestionSet, shuffles, and shows it
func QuestionIndex(c *gin.Context) {
	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	set, err := q.Sets.NextSet(user.CenterX, user.CenterY, user.DeltaX, user.DeltaY)
	if err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	c.JSON(http.StatusOK, set)
}
