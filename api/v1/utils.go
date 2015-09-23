package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
        "regexp"
)

func getRecord(c *gin.Context) (rID uint64, rType string, err error) {
	url := c.Request.RequestURI
	re := regexp.MustCompile("(post|comment)")
	rType = re.FindString(url)

	recordID, ok := c.Params.Get("record_id")
	if ok {
		rID, err = strconv.ParseUint(recordID, 10, 64)
		return rID, rType, err
	}

	return rID, rType, fmt.Errorf("Couldn't parse Params: %v", err)
}
