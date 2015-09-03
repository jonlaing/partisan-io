package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var userID uint64
var postID uint64

func seedUsers(db *gorm.DB) {
	user := User{
		Username:  "user1",
		FullName:  "User One",
		Email:     "user1@email.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	if err != nil {
		fmt.Println("Couldn't make hash:", err)
		return
	}
	user.PasswordHash = hash

	if err := db.Create(&user).Error; err != nil {
		fmt.Println("Coulnd't create user:", err)
	}

	userID = user.ID
}

func seedPosts(db *gorm.DB) {
	post := Post{
		UserID:    1,
		Body:      "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer nec pulvinar quam. Aliquam erat volutpat. Vestibulum tempus egestas nullam.",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Create(&post).Error; err != nil {
		fmt.Println("Couldn't create post:", err)
	}

	postID = post.ID
}

func seedFeed(db *gorm.DB) {
	feedItem := FeedItem{
		UserID:     123,
		Action:     "post",
		RecordType: "post",
		RecordID:   postID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

        if err := db.Create(&feedItem).Error; err != nil {
          fmt.Println("Couldn't create feed item:", err)
        }
}
