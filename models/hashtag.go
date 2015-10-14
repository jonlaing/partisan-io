package models

import (
	"fmt"
	"partisan/Godeps/_workspace/src/github.com/jinzhu/gorm"
	"regexp"
	"strings"
	"time"
)

// Hashtag is attached to posts and comments via Taxonomy
type Hashtag struct {
	ID        uint64    `json:"id" gorm:"primary_key"` // Primary key
	Tag       string    `json:"tag" sql:"unique_index"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Taxonomy links a record to a hashtag
type Taxonomy struct {
	ID         uint64 `json:"id" gorm:"primary_key"` // Primary key
	RecordID   uint64 `json:"id"`
	RecordType string `json:"record_type"`
	HashtagID  uint64 `json:"hashtag_id"`
}

// Hashtagger is an interface for records that can be hashtagged
type Hashtagger interface {
	GetID() uint64
	GetType() string
	GetContent() string
}

// CreateHashtag checks for uniqueness of a tag, creates it if it does not exist,
// and links a Hashtagger via Taxonomy (regardless of uniquenss of tag)
// it does, however, check for uniqueness of the taxonomy
func CreateHashtag(r Hashtagger, tag string, db *gorm.DB) error {
	var hashtag Hashtag
	if err := db.Where(Hashtag{Tag: strings.ToLower(tag)}).FirstOrInit(&hashtag).Error; err != nil {
		return err
	}

	if hashtag.CreatedAt.IsZero() {
		hashtag.CreatedAt = time.Now()
	}

	hashtag.UpdatedAt = time.Now()

	if err := db.Save(&hashtag).Error; err != nil {
		return err
	}

	var taxonomy Taxonomy
	if err := db.Where(Taxonomy{RecordID: r.GetID(), RecordType: r.GetType(), HashtagID: hashtag.ID}).FirstOrInit(&taxonomy).Error; err != nil {
		return err
	}

	return db.Save(&taxonomy).Error
}

// FindAndCreateHashtags searches the content of the Hashtagger for hashtags and saves them
func FindAndCreateHashtags(r Hashtagger, db *gorm.DB) {
	hashtags := ExtractTags(r.GetContent())
	for _, hashtag := range hashtags {
		if err := CreateHashtag(r, hashtag, db); err != nil {
			fmt.Println(err)
		}
	}
}

// ExtractTags finds hashtags in a string
func ExtractTags(text string) []string {
	hashtagSearch := regexp.MustCompile("#([a-zA-Z]+)")
	hashtagSlices := hashtagSearch.FindAllStringSubmatch(text, -1) // -1 means find all instances

	var hashtags []string

	for _, hashtag := range hashtagSlices {
		hashtags = append(hashtags, hashtag[1])
	}

	return hashtags
}
