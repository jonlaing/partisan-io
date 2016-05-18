// Package posts contains functionality to create, find, and update posts, comments and likes.
package posts

import (
	"database/sql"
	"time"

	models "partisan/models.v2"

	"github.com/jinzhu/gorm"
	"github.com/nu7hatch/gouuid"

	"partisan/models.v2/attachments"
	"partisan/models.v2/users"
)

// Post is the representation of posts, comments and likes. See ./types.go for more information
// About what kind of type a post is.
type Post struct {
	ID           string                   `json:"id" gorm:"primary_key" sql:"type:uuid;default:uuid_generate_v4()"`
	UserID       string                   `json:"user_id" sql:"type:uuid"`
	ParentType   ParentType               `json:"parent_type"`
	ParentID     sql.NullString           `json:"-" sql:"type:uuid;default:null"`
	Body         string                   `json:"body"`
	Action       Action                   `json:"action"`
	CreatedAt    time.Time                `json:"created_at"`
	UpdatedAt    time.Time                `json:"updated_at"`
	User         users.User               `json:"user" sql:"-"`
	LikeCount    int                      `json:"like_count" sql:"-"`
	Liked        bool                     `json:"liked" sql:"-"`
	CommentCount int                      `json:"child_count" sql:"-"`
	Attachments  []attachments.Attachment `json:"attachments" sql:"-"`
}

// Posts is a list type of post
type Posts []Post

// CreatorBinding is a struct to use for binding JSON requests to a new Post.
type CreatorBinding struct {
	UserID     string     `json:"user_id" binding:"required"`
	ParentType ParentType `json:"parent_type" binding:"required"`
	ParentID   string     `json:"parent_id" binding:"required"`
	Body       string     `json:"body"`
	Action     Action     `json:"action" binding:"required"`
}

// UpdaterBinding is a struct to use for binding JSON requests to update a Post.
type UpdaterBinding struct {
	Body string `json:"body"`
}

// New uses a CreatorBinding to initialize a new Post and validate it. It does not
// save the Post to the database. This should always be used rather than creating a
// Post manually from user input.
func New(b CreatorBinding) (p Post, errs models.ValidationErrors) {
	p.UserID = b.UserID
	if b.ParentID != "" {
		p.ParentID.String = b.ParentID
		p.ParentID.Valid = true
		p.ParentType = b.ParentType
	}
	if b.Action != ALike {
		p.Body = b.Body
	}

	p.Action = b.Action
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

	errs = p.Validate()

	return
}

// Update uses am UpdaterBinding to update an existing Post and validate it.
// It does not save the Post to the databse. This should always be used rather
// than updating a Post manually from user input.
func (p *Post) Update(b UpdaterBinding) models.ValidationErrors {
	if p.Action != ALike {
		p.Body = b.Body
	}
	return p.Validate()
}

// CanDelete is a helper function to determine wheter a user should be able to
// delete a Post
func (p Post) CanDelete(userID string, db *gorm.DB) bool {
	if userID == p.UserID {
		return true
	}

	// if the user asking is the admin of a group or event
	// if p.ParentType == PTGroup || p.ParentType == PTEvent {
	// 	parent, err := p.GetParent(db)
	// 	if err != nil {
	// 		return false
	// 	}

	// 	if parent.UserID == userID {
	// 		return true
	// 	}
	// }

	return false
}

// Validate validates a Post based on its properties
func (p Post) Validate() models.ValidationErrors {
	errs := make(models.ValidationErrors)

	if _, err := uuid.ParseHex(p.UserID); err != nil {
		errs["user_id"] = ErrUUIDFormat
	}

	if p.ParentID.Valid {
		if _, err := uuid.ParseHex(p.ParentID.String); err != nil {
			errs["parent_id"] = ErrUUIDFormat
		}
	}

	if !validAction(p.Action) {
		errs["action"] = ErrAction
	}

	if !validParentType(p.ParentType) {
		errs["parent_type"] = ErrParentType
	}

	if p.ParentID.Valid && p.Action == APost {
		if !validPostParentType(p.ParentType) {
			errs["parent_type"] = ErrParentType
		}
	}

	if p.Action == AComment || p.Action == ALike {
		if !p.ParentID.Valid {
			errs["parent_id"] = ErrMustHaveParent
		}
	}

	if p.Action == ALike {
		if p.ParentType != PTPost && p.ParentType != PTComment {
			errs["parent_type"] = ErrLikeParent
		}

		if p.Body != "" {
			errs["body"] = ErrLikeBody
		}
	}

	if p.Action == AComment {
		if p.ParentType != PTPost {
			errs["parent_type"] = ErrCommentParent
		}
	}

	return errs
}

// Unique removes any post from the list that has its parent already in the list
func (ps *Posts) Unique() {
	posts := []Post(*ps)
	delIDs := []string{}

	for _, p1 := range posts {
		for _, p2 := range posts {
			if p1.ID == p2.ParentID.String {
				delIDs = append(delIDs, p2.ID)
			}
		}
	}

	for i := len(posts) - 1; i >= 0; i-- {
		shouldDelete := false
		for _, id := range delIDs {
			if posts[i].ID == id {
				shouldDelete = true
			}
		}

		if shouldDelete {
			posts = append(posts[:i], posts[i+1:]...)
		}
	}

	*ps = Posts(posts)
}
