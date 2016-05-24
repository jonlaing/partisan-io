package v2

import (
	"net/http"
	"net/url"
	"partisan/auth"
	"partisan/db"

	"partisan/models.v2/hashtags"

	"github.com/gin-gonic/gin"
)

func HashtagShow(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	q := c.Query("q")
	search, _ := url.QueryUnescape(q)

	page := getPage(c)

	hashtagSearches := hashtags.ExtractTags(search)
	posts, err := hashtags.ListPostsByHashtags(user.ID, hashtagSearches, page*25, db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}
