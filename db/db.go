package db

import (
	"fmt"
	"net/http"
	"os"

	"partisan/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"partisan/Godeps/_workspace/src/github.com/jinzhu/gorm"
	_ "partisan/Godeps/_workspace/src/github.com/lib/pq"
)

var Database gorm.DB

func init() {
	var err error
	if url := os.Getenv("DATABASE_URL"); len(url) > 0 {
		Database, err = gorm.Open("postgres", url)
	} else {
		connString := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PW"))
		Database, err = gorm.Open("postgres", connString)
	}

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
