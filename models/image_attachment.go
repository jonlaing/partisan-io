package models

import (
	"partisan/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"partisan/db"
	"partisan/imager"
	"time"
)

// ImageAttachment is used to attach images to a Post (or something else, when we have it)
type ImageAttachment struct {
	ID         uint64    `json:"id" gorm:"primary_key"`      // Primary key
	UserID     uint64    `json:"user_id" binding:"required"` // ID of user that created ImageAttachment
	RecordID   uint64    `json:"record_id" binding:"required"`
	RecordType string    `json:"record_type" binding:"required"`
	URL        string    `json:"image_url"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ImageAttacher is a record that can have an image attached to it
type ImageAttacher interface {
	GetID() uint64
	GetUserID() uint64
	GetType() string
	AttachImage(ImageAttachment) error
}

// AttachImage attaches an image to an ImageAttacher
func AttachImage(c *gin.Context, r ImageAttacher) error {
	db := db.GetDB(c)

	tmpFile, _, err := c.Request.FormFile("attachment")
	if err != nil {
		return nil // ignore missing file
	}
	defer tmpFile.Close()

	processor := imager.ImageProcessor{File: tmpFile}

	// Save the full-size
	var path string
	if err := processor.Resize(1500); err != nil {
		return err
	}
	path, err = processor.Save("/localfiles/img")
	if err != nil {
		return err
	}

	a := ImageAttachment{
		UserID:     r.GetUserID(),
		RecordID:   r.GetID(),
		RecordType: r.GetType(),
		URL:        path,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err = db.Save(&a).Error; err == nil {
		r.AttachImage(a)
		return nil
	}

	return err
}

// GetAttachment find related attachment for a record
func GetAttachment(recordID uint64, attachments []ImageAttachment) (ImageAttachment, bool) {
	for _, attachment := range attachments {
		if attachment.RecordID == recordID {
			return attachment, true
		}
	}
	return ImageAttachment{}, false
}
