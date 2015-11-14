package models

import (
	"regexp"
	"time"

	"partisan/Godeps/_workspace/src/github.com/jinzhu/gorm"
)

// UserTag is a relational model between a user and a piece of content.
// Users can tag other users in posts and comments.
type UserTag struct {
	ID         uint64    `json:"id" gorm:"primary_key"` // Primary key
	UserID     uint64    `json:"user_id"`
	RecordID   uint64    `json:"record_id"`
	RecordType string    `json:"record_type"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// UserTagger is an interface for any record that can tag a user
type UserTagger interface {
	GetID() uint64
	GetUserID() uint64
	GetType() string
	GetContent() string
}

// FindUserTags searches a UserTagger for instances of a user tag
// in the form of `@username`.
func FindUserTags(r UserTagger, db *gorm.DB) (tags []UserTag) {
	// This regex will match emails too (ex: @gmail.com). That is intentional, and we'll
	// filter those out later
	re := regexp.MustCompile("@([a-zA-Z0-9_.]+)")

	matches := re.FindAllStringSubmatch(r.GetContent(), -1)

	// Protect against someone trying to crash the server by overloading
	// with usertags. 20 seems like a pretty lenient cap.
	if len(matches) > 20 {
		matches = matches[:20]
	}

	for _, match := range matches {
		// This regex searches for emails. Really the `.` and a character after is all
		// we need to determine that a tag is an illegal username, and we'll just skip it.
		if ok, _ := regexp.Match("\\.[a-zA-Z]", []byte(match[0])); ok {
			continue
		}

		var ids []uint64 // Pluck only works with slices, but we only need the first

		// It's okay to do an N+1 here because there will only ever be a handful of tags in each piece
		// of content. Not enough to cause a problem.
		if err := db.Model(User{}).Where("username = ?", match[1]).Pluck("id", &ids).Error; err != nil {
			continue
		}

		tags = append(tags, UserTag{
			UserID:     ids[0],
			RecordID:   r.GetID(),
			RecordType: r.GetType(),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		})
	}

	return

}

// FindAndCreateUserTags searches through a UserTagger, finds user tags, and saves them to the db.
func FindAndCreateUserTags(r UserTagger, db *gorm.DB) (errs []error) {
	tags := FindUserTags(r, db)

	for _, tag := range tags {
		if err := db.Save(&tag).Error; err != nil {
			errs = append(errs, err)
		}
		NewNotification(&tag, r.GetUserID(), db)
	}

	return
}

// GetID satisfies the Notifier interface
func (t *UserTag) GetID() uint64 {
	return t.ID
}

// GetType satisfies the Notifier interface
func (t *UserTag) GetType() string {
	return "user_tag"
}

// GetRecordUserID satisfies the Notifier interface
func (t *UserTag) GetRecordUserID(*gorm.DB) (uint64, error) {
	return t.UserID, nil
}
