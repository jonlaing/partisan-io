package v2

import (
	"net/http"
	"partisan/auth"
	"partisan/db"

	"github.com/gin-gonic/gin"
)

func PostIndex(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	// TODO: DEFINITELY NOT DONE!!
}
