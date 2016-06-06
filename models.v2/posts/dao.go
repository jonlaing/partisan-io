package posts

import (
	"time"

	"partisan/models.v2/attachments"
	"partisan/models.v2/users"

	"github.com/jinzhu/gorm"
)

func GetByID(id string, userID string, db *gorm.DB) (p Post, err error) {
	err = db.Where("id = ?", id).Find(&p).Error
	if p.Action == APost {
		p.GetCommentCount(db)
	}
	p.GetLikeCount(userID, db)
	p.GetUser(db)
	p.Attachments, _ = attachments.GetByPostID(p.ID, db)
	p.GetParent(userID, db)

	return
}

func GetFeedByUserIDs(currentUserID string, userIDs []string, offset int, db *gorm.DB) (ps Posts, err error) {
	err = db.Joins("left join flags on flags.record_id = posts.id").
		Joins("left join flags as pflags on pflags.record_id = posts.parent_id").
		Where("posts.user_id IN (?)", append(userIDs, currentUserID)).
		Where("flags.user_id != ? OR flags.record_id IS NULL", currentUserID).
		Where("pflags.user_id != ? OR pflags.record_id IS NULL", currentUserID).
		Where("posts.parent_type IN (?)", []ParentType{PTNoType, PTPost}).
		Order("posts.created_at desc").
		Limit(25).
		Offset(offset).
		Find(&ps).Error

	ps.Unique()
	ps.GetRelations(currentUserID, db)
	ps.GetUsers(db)
	ps.CollectParents(currentUserID, db)

	return
}

func ListByIDs(currentUserID string, ids []string, offset int, db *gorm.DB) (ps Posts, err error) {
	err = db.Joins("left join flags on flags.record_id = posts.id").
		Where("flags.user_id != ? OR flags.record_id IS NULL", currentUserID).
		Where("posts.id IN (?)", ids).
		Where("posts.parent_type IN (?)", []ParentType{PTNoType, PTPost}).
		Order("posts.created_at desc").
		Limit(25).
		Offset(offset).
		Find(&ps).Error

	ps.Unique()
	ps.GetRelations(currentUserID, db)
	ps.GetUsers(db)

	return
}

func GetFeedByUserIDsAfter(currentUserID string, userIDs []string, after time.Time, db *gorm.DB) (ps Posts, err error) {
	err = db.Joins("left join flags on flags.record_id = posts.id").
		Where("posts.user_id IN (?)", userIDs).
		Where("flags.user_id != ? OR flags.record_id IS NULL", currentUserID).
		Where("posts.parent_type IN (?)", []ParentType{PTNoType, PTPost}).
		Where("posts.created_at >= ?::timestamp", after).
		Order("posts.created_at desc").
		Limit(25).
		Find(&ps).Error

	ps.Unique()
	ps.GetRelations(currentUserID, db)
	ps.GetUsers(db)

	return
}

// GetParent queries the database for the parent of a post
func (p *Post) GetParent(currentUserID string, db *gorm.DB) (err error) {
	if p.Action == APost {
		err = ErrParentQuery
		return
	}

	var parent Post
	err = db.Where("id = ?", p.ParentID.String).Find(&parent).Error
	if err == nil {
		parent.GetUser(db)
		parent.Attachments, _ = attachments.GetByPostID(parent.ID, db)
		p.Parent = parent
	}

	return
}

// CollectParents queries the database for the parent of a group of posts
func (ps *Posts) CollectParents(currentUserID string, db *gorm.DB) (err error) {
	ids := collectParentIDs(ps)

	var parents Posts
	err = db.Where("id IN (?)", ids).Find(&parents).Error
	if err != nil {
		return err
	}

	parents.GetRelations(currentUserID, db)
	parents.GetUsers(db)

	posts := []Post(*ps)

	for i := range posts {
		for _, parent := range parents {
			if parent.ID == posts[i].ParentID.String {
				posts[i].Parent = parent
			}
		}
	}

	*ps = Posts(posts)

	return
}

func (p Post) GetChildren(currentUserID string, db *gorm.DB) (c Posts, err error) {
	err = db.Joins("left join flags on flags.record_id = posts.id").
		Where("posts.parent_id = ?", p.ID).
		Where("flags.user_id != ? OR flags.record_id IS NULL", currentUserID).
		Order("posts.created_at desc").
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
		Find(&cs).Error

	cs.GetRelations(userID, db)
	cs.GetUsers(db)

	return
}

func (p *Post) GetCommentCount(db *gorm.DB) error {
	return db.Table("posts").
		Where("parent_id = ?", p.ID).
		Where("action = ?", AComment).
		Count(&p.CommentCount).Error
}

func (p *Post) GetUser(db *gorm.DB) error {
	return db.Where("id = ?", p.UserID).Find(&p.User).Error
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

func (ps *Posts) GetUsers(db *gorm.DB) error {
	users, err := users.ListRelated(ps, db)
	if err != nil {
		return err
	}

	posts := []Post(*ps)
	for i := range posts {
		for _, u := range users {
			if posts[i].UserID == u.ID {
				posts[i].User = u
			}
		}
	}

	*ps = Posts(posts)
	return nil
}

func collectIDs(ps *Posts) []string {
	var ids = make([]string, 0)

	for _, p := range *ps {
		ids = append(ids, p.ID)
	}

	return ids
}

func collectParentIDs(ps *Posts) []string {
	var ids = make([]string, 0)

	for _, p := range *ps {
		if p.Action != APost && p.ParentID.Valid {
			ids = append(ids, p.ParentID.String)
		}
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
