package hashtags

import (
	"github.com/jinzhu/gorm"
	"partisan/models.v2/posts"
)

func ListPostsByHashtags(tags []string, userID string, offset int, db *gorm.DB) (ps posts.Posts, err error) {
	var postIDs []string
	err = db.Model(Taxonomy{}).
		Joins("inner join hashtags on taxonomies.hashtag_id = hashtags.id").
		Where("tag IN (?) AND record_type = ?", tags, "post").
		Order("created_at DESC").
		Pluck("record_id", &postIDs).Error
	if err != nil {
		return
	}

	return posts.ListByIDs(userID, postIDs, offset, db)
}
