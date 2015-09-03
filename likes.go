package main

import (
	// "github.com/gin-gonic/gin"
	// "github.com/jinzhu/gorm"
	// "net/http"
)

type Like struct {
	ID         uint `gorm:"primary_key"`
	UserID     uint
	RecordID   uint
	RecordType uint
}

type Dislike struct {
	Like
}
