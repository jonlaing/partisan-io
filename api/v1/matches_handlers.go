package v1

import (
	"github.com/gin-gonic/gin"
	"partisan/auth"
	"partisan/db"
	m "partisan/models"
        "net/http"
)

// MatchesIndex returns a list of matches orderd by location and match percentage
func MatchesIndex(c *gin.Context) {
	db, err := db.InitDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()
}
