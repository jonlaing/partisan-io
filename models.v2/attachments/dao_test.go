package attachments

import (
	"fmt"
	"os"

	_ "github.com/lib/pq"

	"github.com/jinzhu/gorm"
)

var testdb *gorm.DB

func init() {
	var err error
	connString := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", os.Getenv("DB_TEST_USER"), os.Getenv("DB_TEST_NAME"), os.Getenv("DB_TEST_PW"))
	testdb, err = gorm.Open("postgres", connString)
	if err != nil {
		panic(err)
	}

	if err := testdb.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		panic(err)
	}

	testdb.AutoMigrate(Attachment{})

}
