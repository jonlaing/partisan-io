package v1

import (
	"fmt"
	m "partisan/models"
	"regexp"
	"strconv"

	"partisan/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"partisan/Godeps/_workspace/src/github.com/jinzhu/gorm"
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

func findMatchingPostUser(post m.Post, users []m.User) (m.User, bool) {
	for _, user := range users {
		if user.ID == post.UserID {
			return user, true
		}
	}
	return m.User{}, false
}

func fineMatchingPostLikes(post m.Post, likes []m.RecordLikes) (int, bool, bool) {
	for _, like := range likes {
		if like.RecordID == post.ID {
			return like.Count, like.UserCount == 1, true
		}
	}
	return 0, false, false
}

func findMatchingCommentCount(post m.Post, comments []PostComments) (int, bool) {
	for _, comment := range comments {
		if comment.RecordID == post.ID {
			return comment.Count, true
		}
	}
	return 0, false
}

func getPostComments(postIDs []uint64, db *gorm.DB) ([]PostComments, error) {
	var comments []PostComments

	rows, err := db.Raw("SELECT count(*), post_id FROM \"comments\"  WHERE (post_id IN (?)) group by post_id", postIDs).Rows()
	defer rows.Close()
	if err != nil {
		return []PostComments{}, err
	}

	for rows.Next() {
		var count int
		var rID uint64

		rows.Scan(&count, &rID)
		comments = append(comments, PostComments{Count: count, RecordID: rID})
	}

	return comments, nil
}

func getPage(c *gin.Context) int {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		return 1
	}

	return page
}
