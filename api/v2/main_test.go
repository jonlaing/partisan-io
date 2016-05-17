package v2

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"

	"partisan/auth"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"partisan/models.v2/users"
)

var testDB *gorm.DB
var testRouter *gin.Engine

func init() {
	var err error
	connString := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", os.Getenv("DB_TEST_USER"), os.Getenv("DB_TEST_NAME"), os.Getenv("DB_TEST_PW"))
	testDB, err = gorm.Open("postgres", connString)
	if err != nil {
		panic(err)
	}

	if err := testDB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		panic(err)
	}

	testDB.AutoMigrate(users.User{})

	testRouter = gin.Default()
	testRouter.Use(addTestDB())
}

func addTestDB() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := testDB.DB().Ping(); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Set("db", testDB)
	}
}

func performRequest(r http.Handler, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func login(user *users.User) gin.HandlerFunc {
	return func(c *gin.Context) {
		tok, _ := auth.Login(user, c)
		c.Request.Header.Set("X-Auth-Token", tok)
		auth.Auth()(c)
	}
}
