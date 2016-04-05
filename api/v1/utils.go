package v1

import (
	"fmt"
	"partisan/dao"
	m "partisan/models"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
)

// getRecord looks in the `GET` params for a post or comment ID.
// i.e. `/api/v1/(post|comment)/:record_id`
func getRecord(c *gin.Context) (rID uint64, rType string, err error) {
	url := c.Request.RequestURI
	re := regexp.MustCompile("(post|comment)")
	rType = re.FindString(url)

	recordID, ok := c.Params.Get("record_id")
	if ok {
		rID, err = strconv.ParseUint(recordID, 10, 64)
		if err != nil {
			return rID, rType, &ErrParseID{err}
		}

		return rID, rType, nil
	}

	return rID, rType, &ErrParseID{fmt.Errorf("Couldn't parse Params: %v", url)}
}

func findMatchingPostLikes(post m.Post, likes []m.RecordLikes) (int, bool, bool) {
	for _, like := range likes {
		if like.RecordID == post.ID {
			return like.Count, like.UserCount == 1, true
		}
	}
	return 0, false, false
}

func findMatchingCommentCount(post m.Post, comments []dao.PostComments) (int, bool) {
	for _, comment := range comments {
		if comment.RecordID == post.ID {
			return comment.Count, true
		}
	}
	return 0, false
}

func getPage(c *gin.Context) int {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		return 1
	}

	return page
}
