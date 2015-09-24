package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"partisan/db"
	m "partisan/models"
)

// HashtagShow shows a list of Posts (and Comments) that contain a particular hashtag
func HashtagShow(c *gin.Context) {
	db, err := db.InitDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	search := c.Query("q")

	hashtagSearches := m.ExtractTags(search)

	var postIDs uint64
	if err := db.Model(m.Taxonomy{}).
		Joins("inner join hashtags on taxonomies.hashtag_id = hashtags.id").
		Where("tag IN (?) AND record_type = ?", hashtagSearches, "post").
		Pluck("record_id", &postIDs).Error; err != nil {

		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	var posts []m.Post
	if err := db.Find(&posts, postIDs).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, posts)
}
