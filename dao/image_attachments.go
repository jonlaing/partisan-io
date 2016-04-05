package dao

import (
	m "partisan/models"

	"github.com/jinzhu/gorm"
)

type ImageAttacher interface {
	GetRecordTypeAndIDs() (t string, ids []uint64, err error)
}

func GetRelatedAttachments(rs ImageAttacher, db *gorm.DB) (attachments []m.ImageAttachment, err error) {
	t, ids, err := rs.GetRecordTypeAndIDs()
	if err != nil {
		return []m.ImageAttachment{}, err
	}

	err = db.Where("record_type = ? AND record_id IN (?)", t, ids).Find(&attachments).Error

	return
}
