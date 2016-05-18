package attachments

import "github.com/jinzhu/gorm"

func GetByPostID(postID string, db *gorm.DB) (as []Attachment, err error) {
	err = db.Where("post_id = ?", postID).Find(&as).Error
	return
}

func GetByPostIDs(postIDs []string, db *gorm.DB) (as []Attachment, err error) {
	err = db.Where("post_id IN (?)", postIDs).Find(&as).Error
	return
}
