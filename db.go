package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

func initDB() (db gorm.DB, err error) {
	db, err = gorm.Open("postgres", "user=partisan dbname=partisan password=bakunin1234 sslmode=disable")
	return
}
