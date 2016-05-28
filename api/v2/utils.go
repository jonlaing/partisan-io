package v2

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func getPage(c *gin.Context) int {
	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
	if err != nil {
		return 0
	}

	return page
}
