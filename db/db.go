package db

import (
	"net/http"
	"partisan/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"partisan/Godeps/_workspace/src/github.com/jinzhu/gorm"
	_ "partisan/Godeps/_workspace/src/github.com/lib/pq"
)

var Database gorm.DB

func init() {
	var err error
	Database, err = gorm.Open("postgres", "user=partisan dbname=partisan password=bakunin1234 sslmode=disable")
	if err != nil {
		panic(err)
	}
}

// DB is middleware to get the database
func DB() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := Database.DB().Ping(); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Set("db", &Database)
	}
}

func GetDB(c *gin.Context) *gorm.DB {
	db, ok := c.Get("db")
	if !ok {
		panic("couldn't get database")
	}

	return db.(*gorm.DB)
}
