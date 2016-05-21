package friendships

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/nu7hatch/gouuid"
)

var testdb *gorm.DB
var userID, friendID1, friendID2 string
var unconfirmed, confirmed, dontfind Friendship

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

	testdb.AutoMigrate(Friendship{})
}

func TestMain(m *testing.M) {
	tuuid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	userID = tuuid.String()

	tuuid, err = uuid.NewV4()
	if err != nil {
		panic(err)
	}
	friendID1 = tuuid.String()

	tuuid, err = uuid.NewV4()
	if err != nil {
		panic(err)
	}
	friendID2 = tuuid.String()

	unconfirmed = Friendship{
		UserID:    userID,
		FriendID:  friendID1,
		Confirmed: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := testdb.Save(&unconfirmed).Error; err != nil {
		panic(err)
	}
	defer testdb.Delete(&unconfirmed)

	confirmed = Friendship{
		UserID:    friendID2,
		FriendID:  userID,
		Confirmed: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := testdb.Save(&confirmed).Error; err != nil {
		panic(err)
	}
	defer testdb.Delete(&confirmed)

	dontfind = Friendship{
		UserID:    friendID2,
		FriendID:  friendID1,
		Confirmed: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := testdb.Save(&dontfind).Error; err != nil {
		panic(err)
	}
	defer testdb.Delete(&dontfind)

	m.Run()
}

func TestGetByUserID(t *testing.T) {
	fs, err := ListByUserID(userID, testdb)
	if err != nil {
		t.Error("Unexpected error")
	}

	if len(fs) != 2 {
		t.Error("Expected 2 friendships, got:", len(fs))
	}

	foundFirst := false
	foundSecond := false
	for _, f := range fs {
		if f.ID == unconfirmed.ID {
			foundFirst = true
		}

		if f.ID == confirmed.ID {
			foundSecond = true
		}
	}

	if !foundFirst {
		t.Error("Expected to find unconfirmed friendship in list")
	}

	if !foundSecond {
		t.Error("Expected to find confirmed friendship in list")
	}
}

func TestGetIDsByUserID(t *testing.T) {
	uuids, err := GetIDsByUserID(userID, testdb)
	if err != nil {
		t.Error("Unexpected error")
	}

	if len(uuids) != 2 {
		t.Error("Expected 2 uuids, got:", len(uuids))
	}

	foundFirst := false
	foundSecond := false
	for _, id := range uuids {
		if id == unconfirmed.UserID || id == unconfirmed.FriendID {
			foundFirst = true
		}

		if id == confirmed.UserID || id == confirmed.FriendID {
			foundSecond = true
		}
	}

	if !foundFirst {
		t.Error("Expected to find unconfirmed friendship in list")
	}

	if !foundSecond {
		t.Error("Expected to find confirmed friendship in list")
	}
}

func TestGetConfirmedByUserID(t *testing.T) {
	fs, err := GetConfirmedByUserID(userID, testdb)
	if err != nil {
		t.Error("Unexpected error")
	}

	if len(fs) != 1 {
		t.Error("Expected 1 friendship, got:", len(fs))
	}

	if fs[0].ID != confirmed.ID {
		t.Error("Expected friendship to be confirmed")
	}
}

func TestGetConfirmedIDsByUserID(t *testing.T) {
	uuids, err := GetConfirmedIDsByUserID(userID, testdb)
	if err != nil {
		t.Error("Unexpected error")
	}

	if len(uuids) != 1 {
		t.Error("Expected 1 uuid, got:", len(uuids))
	}

	if uuids[0] != confirmed.UserID && uuids[0] != confirmed.FriendID {
		t.Error("Expected friendship to be confirmed")
	}
}
