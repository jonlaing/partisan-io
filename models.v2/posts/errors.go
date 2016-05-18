package posts

import "errors"

var (
	ErrNoUpdate       = errors.New("Cannot update this post")
	ErrAction         = errors.New("Unknown Post.Action")
	ErrParentType     = errors.New("Post cannot be a child to this type")
	ErrUUIDFormat     = errors.New("This UUID is not formatted correctly")
	ErrLikeParent     = errors.New("Invalid parent type for like post")
	ErrLikeBody       = errors.New("Post with Like action cannot have a body")
	ErrMustHaveParent = errors.New("Posts with Comment and Like actions must have a parent")
	ErrCommentParent  = errors.New("Posts with Comment action can only be children of Post types")
	ErrScanAction     = errors.New("Could not scan action")
	ErrScanParentType = errors.New("Could not scan parent type")
	ErrParentQuery    = errors.New("Cannot get this post's parent in this way. This only works for comments and likes")
)
