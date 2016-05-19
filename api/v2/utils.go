package v2

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func getPage(c *gin.Context) int {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		return 1
	}

	return page
}
