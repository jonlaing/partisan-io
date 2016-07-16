package v2

import (
	"net/http"

	"partisan/auth"
	q "partisan/questions"

	"github.com/gin-gonic/gin"
)

var questionSets q.QuestionSets

func init() {
	questionSets = q.QuestionSets{}
	questionSets = append(questionSets, q.InitialQuestions...)
	questionSets = append(questionSets, q.LeftWingQuestions...)
	questionSets = append(questionSets, q.RightWingQuestions...)
	questionSets = append(questionSets, q.AuthoritarianQuestions...)
	questionSets = append(questionSets, q.LibertarianQuestions...)
	questionSets = append(questionSets, q.AuthoritarianSocialistQuestions...)
	questionSets = append(questionSets, q.LibertarianSocialistQuestions...)
	questionSets = append(questionSets, q.AuthoritarianCapitalistQuestions...)
	questionSets = append(questionSets, q.LibertarianCapitalistQuestions...)
}

// QuestionIndex finds a random QuestionSet, shuffles, and shows it
func QuestionIndex(c *gin.Context) {
	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	set, err := questionSets.NextSet(user.CenterX, user.CenterY)
	if err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	c.JSON(http.StatusOK, set)
}
