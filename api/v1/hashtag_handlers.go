package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"partisan/auth"
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

	user, _ := auth.CurrentUser(c)

	q := c.Query("q")
	search, _ := url.QueryUnescape(q)

	hashtagSearches := m.ExtractTags(search)

	var postIDs []uint64
	if err := db.Model(m.Taxonomy{}).
		Joins("inner join hashtags on taxonomies.hashtag_id = hashtags.id").
		Where("tag IN (?) AND record_type = ?", hashtagSearches, "post").
		Order("created_at DESC").
		Pluck("record_id", &postIDs).Error; err != nil {

		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	var posts []m.Post
	var postUserIDs []uint64
	if err := db.Find(&posts, postIDs).Error; err == nil {
		for _, post := range posts {
			postUserIDs = append(postUserIDs, post.UserID)
		}
	} else {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	var users []m.User
	db.Where("id IN (?)", postUserIDs).Find(&users)

	var attachments []m.ImageAttachment
	db.Where("record_type = ? AND record_id IN (?)", "post", postIDs).Find(&attachments)

	postLikes, err := m.GetLikes(user.ID, "post", postIDs, &db)
	if err != nil {
		fmt.Println(err)
	}

	postComments, err := getPostComments(postIDs, &db)
	if err != nil {
		fmt.Println(err)
	}

	var resp []PostResponse
	for _, post := range posts {
		user, _ := findMatchingPostUser(post, users)
		attachment, _ := m.GetAttachment(post.ID, attachments)
		likeCount, liked, _ := fineMatchingPostLikes(post, postLikes)
		commentCount, _ := findMatchingCommentCount(post, postComments)

		resp = append(resp, PostResponse{
			Post:         post,
			User:         user,
			Attachment:   attachment,
			LikeCount:    likeCount,
			Liked:        liked,
			CommentCount: commentCount,
		})
	}

	c.JSON(http.StatusOK, resp)
}
