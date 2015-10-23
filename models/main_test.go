package models

import (
	"fmt"
	"os"

	"partisan/Godeps/_workspace/src/github.com/jinzhu/gorm"
)

var testDB gorm.DB

func init() {
	var err error
	connString := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", os.Getenv("DB_TEST_USER"), os.Getenv("DB_TEST_NAME"), os.Getenv("DB_TEST_PW"))
	testDB, err = gorm.Open("postgres", connString)
	if err != nil {
		panic(err)
	}

	testDB.AutoMigrate(
		&Post{},
		&User{},
		&Friendship{},
		&FeedItem{},
		&Like{},
		&Profile{},
		&Comment{},
		&ImageAttachment{},
		&Notification{},
		&Hashtag{},
		&Taxonomy{},
		&Flag{},
		&UserTag{},
	)
}
