package users

import (
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"github.com/jinzhu/gorm"
	"github.com/nu7hatch/gouuid"
)

var testdb *gorm.DB
var u User
var id, apiKey string
var exp time.Time

type useridslicer struct{}

func (u *useridslicer) GetUserIDs() []string {
	return []string{id}
}

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

	testdb.AutoMigrate(User{})

}

func TestMain(m *testing.M) {
	key, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	apiKey = key.String()

	exp = time.Now().Add(168 * time.Hour)

	u = User{
		Username:  "user1",
		Email:     "user1@email.com",
		APIKey:    apiKey,
		APIKeyExp: exp,
	}

	if err := testdb.Save(&u).Error; err != nil {
		panic(err)
	}
	defer testdb.Delete(&u)

	id = u.ID

	m.Run()
}

func TestGetByID(t *testing.T) {
	if id == "" {
		t.Error("Expected ID to have value")
		return
	}

	user, err := GetByID(id, testdb)
	if err != nil {
		t.Error("Unexpected error:", err)
		return
	}

	if user.ID != u.ID {
		t.Error("Expected IDs to match:", user.ID, u.ID)
	}
}

func TestGetByUsername(t *testing.T) {
	user, err := GetByUsername("user1", testdb)
	if err != nil {
		t.Error("Unexpected error:", err)
		return
	}

	if user.ID != u.ID {
		t.Error("Expected IDs to match:", user.ID, u.ID)
	}
}

func TestGetByEmail(t *testing.T) {
	user, err := GetByEmail("user1@email.com", testdb)
	if err != nil {
		t.Error("Unexpected error:", err)
		return
	}

	if user.ID != u.ID {
		t.Error("Expected IDs to match:", user.ID, u.ID)
	}
}

func TestGetByAPIKey(t *testing.T) {
	user, err := GetByAPIKey(apiKey, testdb)
	if err != nil {
		t.Error("Unexpected error:", err)
		return
	}

	if user.ID != u.ID {
		t.Error("Expected IDs to match:", user.ID, u.ID)
	}
}

func TestListByIDs(t *testing.T) {
	user, err := ListByIDs([]string{id}, testdb)
	if err != nil {
		t.Error("Unexpected error:", err)
		return
	}

	if user[0].ID != u.ID {
		t.Error("Expected IDs to match:", user[0].ID, u.ID)
	}
}

func TestListRelated(t *testing.T) {
	slicer := useridslicer{}

	user, err := ListRelated(&slicer, testdb)
	if err != nil {
		t.Error("Unexpected error:", err)
		return
	}

	if user[0].ID != u.ID {
		t.Error("Expected IDs to match:", user[0].ID, u.ID)
	}
}
