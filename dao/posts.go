package dao

import (
	"errors"
	m "partisan/models"

	"partisan/Godeps/_workspace/src/github.com/jinzhu/gorm"
)

// PostResponse is the response schema
type PostResponse struct {
	Post         m.Post            `json:"post"`
	Attachment   m.ImageAttachment `json:"image_attachment"`
	User         m.User            `json:"user"`
	LikeCount    int               `json:"like_count"`
	Liked        bool              `json:"liked"`
	CommentCount int               `json:"comment_count"`
}

type PostIDer interface {
	GetPostIDs() []uint64
}

func GetRelatedPosts(rs PostIDer, db *gorm.DB) (posts m.Posts, err error) {
	ids := rs.GetPostIDs()

	if len(ids) < 1 {
		err = errors.New("Need at least one ID in list")
		return
	}

	err = db.Where("id IN (?)", ids).Find(&posts).Error
	return
}

func GetRelatedPostReponse(currentUserID uint64, rs PostIDer, db *gorm.DB) (resp map[uint64]PostResponse, err error) {
	ps, err := GetRelatedPosts(rs, db)
	if err != nil {
		return resp, err
	}

	users, err := GetRelatedUsers(ps, db)
	if err != nil {
		return
	}

	attachments, err := GetRelatedAttachments(ps, db)
	if err != nil {
		return
	}

	postLikes, _ := GetRelatedLikes(currentUserID, ps, db)

	postComments, _ := GetRelatedComments(ps, db)

	resp = collectPostResponses(ps, users, attachments, postLikes, postComments)
	return
}

// TODO: benchmark this jawn
func collectPostResponses(posts []m.Post, users []m.User, attachments []m.ImageAttachment, likes []RecordLikes, comments []PostComments) (pr map[uint64]PostResponse) {
	pr = make(map[uint64]PostResponse, len(posts))

	for _, post := range posts {
		user, _ := GetMatchingUser(&post, users)
		attachment, _ := findRelatedPostAttachment(post, attachments)
		likeCount, liked, _ := findMatchingPostLikes(post, likes)
		commentCount, _ := findMatchingPostCommentCount(post, comments)

		pr[post.ID] = PostResponse{
			Post:         post,
			User:         user,
			Attachment:   attachment,
			LikeCount:    likeCount,
			Liked:        liked,
			CommentCount: commentCount,
		}
	}

	return
}

func findMatchingPostUser(post m.Post, users []m.User) (m.User, bool) {
	for _, user := range users {
		if user.ID == post.UserID {
			return user, true
		}
	}
	return m.User{}, false
}

func findMatchingPostLikes(post m.Post, likes []RecordLikes) (int, bool, bool) {
	for _, like := range likes {
		if like.RecordID == post.ID {
			return like.Count, like.UserCount == 1, true
		}
	}
	return 0, false, false
}

func findMatchingPostCommentCount(post m.Post, comments []PostComments) (int, bool) {
	for _, comment := range comments {
		if comment.RecordID == post.ID {
			return comment.Count, true
		}
	}
	return 0, false
}

func findRelatedPostAttachment(post m.Post, attachments []m.ImageAttachment) (m.ImageAttachment, bool) {
	for _, attachment := range attachments {
		if attachment.RecordType == "post" && attachment.RecordID == post.ID {
			return attachment, true
		}
	}
	return m.ImageAttachment{}, false
}
