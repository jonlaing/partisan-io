package dao

import (
	m "partisan/models"

	"github.com/jinzhu/gorm"
)

type ImageAttachers interface {
	GetRecordTypeAndIDs() (t string, ids []uint64, err error)
}

type ImageAttacher interface {
	GetID() uint64
	GetType() string
}

func GetMultipleRelatedAttachments(rs ImageAttachers, db *gorm.DB) (attachments []m.ImageAttachment, err error) {
	t, ids, err := rs.GetRecordTypeAndIDs()
	if err != nil {
		return []m.ImageAttachment{}, err
	}

	err = db.Where("record_type = ? AND record_id IN (?)", t, ids).Find(&attachments).Error

	return
}

func GetRelatedAttachments(rs ImageAttacher, db *gorm.DB) (attachments []m.ImageAttachment, err error) {
	id := rs.GetID()
	t := rs.GetType()

	err = db.Where("record_type = ? AND record_id = ?", t, id).Find(&attachments).Error

	return
}
