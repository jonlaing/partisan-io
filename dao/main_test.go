package dao

import (
	"fmt"
	"os"

	m "partisan/models"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func init() {
	var err error
	connString := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", os.Getenv("DB_TEST_USER"), os.Getenv("DB_TEST_NAME"), os.Getenv("DB_TEST_PW"))
	db, err = gorm.Open("postgres", connString)
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(
		&m.Post{},
		&m.User{},
		&m.Friendship{},
		&m.FeedItem{},
		&m.Like{},
		&m.Profile{},
		&m.Comment{},
		&m.ImageAttachment{},
		&m.Notification{},
		&m.Hashtag{},
		&m.Taxonomy{},
		&m.Flag{},
		&m.UserTag{},
		&m.Message{},
		&m.MessageThread{},
		&m.MessageThreadUser{},
	)
}
