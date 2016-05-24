package v2

import (
	"net/http"
	"partisan/auth"
	"partisan/db"

	"partisan/models.v2/matches"

	"github.com/gin-gonic/gin"
)

func MatchIndex(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	var search matches.SearchBinding
	if err := c.BindJSON(&search); err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	matches, err := matches.List(user, search, db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"matches": matches})
}
