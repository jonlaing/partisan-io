package attachments

import (
	"mime/multipart"
	"os"
	"partisan/imager"
	"time"
)

// Attachment types
const (
	// TNone - zero value
	TNone int = iota
	// TLink - An external link
	TLink
	// TImage - a still image
	TImage
	// TGif - a .gif image
	TGif
	// TVideo - an external video
	TVideo
)

type Attachment struct {
	ID        string    `json:"id" gorm:"primary_key" sql:"type:uuid;default:uuid_generate_v4()"`
	UserID    string    `json:"user_id" sql:"type:uuid"`
	PostID    string    `json:"post_id" sql:"type:uuid"`
	Type      int       `json:"type"`
	URL       string    `json:"url" sql:"not_null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewImage(userID, postID string, f multipart.File) (Attachment, error) {
	var err error
	processor := imager.ImageProcessor{File: f}

	// Save the full-size
	var path string
	if err = processor.Resize(1500); err != nil {
		return Attachment{}, err
	}

	// check if we're using AWS
	if len(os.Getenv("AWS_ACCESS_KEY_ID")) > 0 {
		path, err = processor.Save("/img")
		if err != nil {
			return Attachment{}, err
		}
	} else {
		// if not, save locally
		path, err = processor.Save("/localfiles/img")
		if err != nil {
			return Attachment{}, err
		}
	}

	a := Attachment{
		UserID:    userID,
		PostID:    postID,
		Type:      TImage,
		URL:       path,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return a, nil
}
