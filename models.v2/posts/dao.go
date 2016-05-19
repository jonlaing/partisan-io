package posts

import (
	"time"

	"partisan/models.v2/attachments"

	"github.com/jinzhu/gorm"
)

func GetByID(id string, userID string, db *gorm.DB) (p Post, err error) {
	err = db.Where("id = ?", id).Find(&p).Error
	if p.Action == APost {
		p.GetCommentCount(db)
	}
	p.GetLikeCount(userID, db)
	return
}

func GetFeedByUserIDs(currentUserID string, userIDs []string, offset int, db *gorm.DB) (ps Posts, err error) {
	// TODO: Get flagging in there
	// err = db.Joins("left join flags on flags.record_id = posts.id").
	// 	Where("posts.user_id IN (?)", userIDs).
	// 	Where("flags.user_id != ? OR flags.record_id IS NULL", currentUserID).
	// 	Order("posts.created_at desc").
	// 	Limit(25).
	// 	Offset(offset).
	// 	Find(&ps).Error

	err = db.Where("user_id IN (?)", userIDs).
		Where("parent_type IN (?)", []ParentType{PTNoType, PTPost}).
		Order("created_at desc").
		Limit(25).
		Offset(offset).
		Find(&ps).Error

	ps.Unique()

	return
}

func GetFeedByUserIDsAfter(currentUserID string, userIDs []string, after time.Time, db *gorm.DB) (ps Posts, err error) {
	// TODO: Get flagging in there
	// err = db.Joins("left join flags on flags.record_id = posts.id").
	// 	Where("posts.user_id IN (?)", userIDs).
	// 	Where("flags.user_id != ? OR flags.record_id IS NULL", currentUserID).
	// 	Where("posts.create_at >= ?::timestamp", after).
	// 	Order("posts.created_at desc").
	// 	Limit(25).
	// 	Find(&ps).Error
	err = db.Where("user_id IN (?)", userIDs).
		Where("parent_type IN (?)", []ParentType{PTNoType, PTPost}).
		Where("created_at >= ?::timestamp", after).
		Order("created_at desc").
		Limit(25).
		Find(&ps).Error

	ps.Unique()

	return
}

// GetParent queries the database for the parent of a post
func (p Post) GetParent(db *gorm.DB) (parent Post, err error) {
	if p.Action == APost {
		err = ErrParentQuery
		return
	}

	err = db.Where("id = ?", p.ParentID.String).Find(&parent).Error
	return
}

func (p Post) GetChildren(currentUserID string, db *gorm.DB) (c Posts, err error) {
	// TODO: Get flagging in there
	// err = db.Joins("left join flags on flags.record_id = posts.id").
	// 	Where("posts.parent_id = ?", p.ID).
	// 	Where("flags.user_id != ? OR flags.record_id IS NULL", currentUserID).
	// 	Order("posts.created_at desc").
	// 	Find(&c).Error

	err = db.Where("parent_id = ?", p.ID).
		Order("created_at desc").
		Find(&c).Error

	c.GetRelations(currentUserID, db)

	return
}

func (p *Post) GetLikeByUserID(userID string, db *gorm.DB) (l Post, err error) {
	err = db.Where("parent_id = ?", p.ID).
		Where("user_id = ?", userID).
		Where("action = ?", ALike).
		Find(&l).Error

	return
}

func (p *Post) GetLikeCount(userID string, db *gorm.DB) error {
	var likes []Post
	err := db.Table("posts").
		Where("parent_id = ?", p.ID).
		Where("action = ?", ALike).
		Find(&likes).Error

	p.LikeCount = len(likes)

	for _, l := range likes {
		if l.UserID == userID {
			p.Liked = true
		}
	}

	return err
}

func (p *Post) GetComments(userID string, db *gorm.DB) (cs Posts, err error) {
	err = db.Where("parent_id = ?", p.ID).
		Where("action = ?", AComment).
		Find(&c).Error

	cs.GetRelations(userID, db)

	return
}

func (p *Post) GetCommentCount(db *gorm.DB) error {
	return db.Table("posts").
		Where("parent_id = ?", p.ID).
		Where("action = ?", AComment).
		Count(&p.CommentCount).Error
}

func (ps *Posts) GetChildCountList(userID string, db *gorm.DB) error {
	postIDs := collectIDs(ps)
	posts := []Post(*ps)

	var children Posts
	err := db.Where("parent_id IN (?)", postIDs).Find(&children).Error
	if err != nil {
		return err
	}

	for i := range posts {
		// resetting it just incase someone called this twice for some dumb reason
		posts[i].CommentCount = 0
		posts[i].LikeCount = 0
		posts[i].Liked = false

		for _, child := range children {
			if child.ParentID.String == posts[i].ID {
				if child.Action == AComment {
					posts[i].CommentCount++
				}

				if child.Action == ALike {
					posts[i].LikeCount++
					if child.UserID == userID {
						posts[i].Liked = true
					}
				}
			}
		}
	}

	*ps = Posts(posts)

	return nil
}

func (ps *Posts) GetRelations(userID string, db *gorm.DB) {
	postIDs := collectIDs(ps)
	posts := []Post(*ps)

	as, _ := attachments.GetByPostIDs(postIDs, db)
	for i := range posts {
		if pas, ok := getAttachments(posts[i], as); ok {
			posts[i].Attachments = pas
		}
	}

	ps.GetChildCountList(userID, db)

	*ps = Posts(posts)
}

func collectIDs(ps *Posts) []string {
	var ids = make([]string, 0)

	for _, p := range *ps {
		ids = append(ids, p.ID)
	}

	return ids
}

func getAttachments(p Post, as []attachments.Attachment) ([]attachments.Attachment, bool) {
	var pas = make([]attachments.Attachment, 0)

	for _, a := range as {
		if a.PostID == p.ID {
			pas = append(pas, a)
		}
	}

	return pas, len(pas) > 0
}
