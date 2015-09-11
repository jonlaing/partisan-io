package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

// InitDB establishes a connnection to the db
func InitDB() (db gorm.DB, err error) {
	db, err = gorm.Open("postgres", "user=partisan dbname=partisan password=bakunin1234 sslmode=disable")
	return
}
