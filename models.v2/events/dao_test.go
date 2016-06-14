package events

import (
	"fmt"
	"os"
	"partisan/matcher"
	"testing"
	"time"

	"partisan/models.v2/users"

	_ "github.com/lib/pq"

	"github.com/jinzhu/gorm"
)

var testdb *gorm.DB
var testLook, testGuest, testHost users.User
var testEvent Event
var testID string

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

	testdb.AutoMigrate(Event{}, EventSubscription{})
}

func TestMain(m *testing.M) {
	// CREATE TEST HOST
	uBinding := users.CreatorBinding{
		Username:        "testhost",
		Email:           "testhost@email.com",
		PostalCode:      "11233",
		Password:        "password",
		PasswordConfirm: "password",
	}

	testHost, _ = users.New(uBinding)
	testHost.PoliticalMap = matcher.PoliticalMap{
		0, 0, 0, 0,
		0, 0, 0, 0,
		4, 4, 0, 0,
		4, 4, 0, 0,
	}
	x, y := testHost.PoliticalMap.Center()
	testHost.CenterX = x
	testHost.CenterY = y

	if err := testdb.Create(&testHost).Error; err != nil {
		panic(err)
	}
	defer testdb.Delete(&testHost)

	// CREATE TEST GUEST
	uBinding = users.CreatorBinding{
		Username:        "testguest",
		Email:           "testguest@email.com",
		PostalCode:      "11233",
		Password:        "password",
		PasswordConfirm: "password",
	}
	testGuest, _ = users.New(uBinding)
	testGuest.PoliticalMap = matcher.PoliticalMap{
		0, 0, 0, 0,
		0, 0, 0, 0,
		4, 4, 0, 0,
		4, 4, 0, 0,
	}
	x, y = testGuest.PoliticalMap.Center()
	testGuest.CenterX = x
	testGuest.CenterY = y

	if err := testdb.Create(&testGuest).Error; err != nil {
		panic(err)
	}
	defer testdb.Delete(&testGuest)

	// CREATE TEST USER NOT SUBSCRIBED
	uBinding = users.CreatorBinding{
		Username:        "testlook",
		Email:           "testlook@email.com",
		PostalCode:      "11233",
		Password:        "password",
		PasswordConfirm: "password",
	}
	testLook, _ = users.New(uBinding)
	testLook.PoliticalMap = matcher.PoliticalMap{
		0, 0, 0, 0,
		0, 0, 0, 0,
		4, 4, 0, 0,
		4, 4, 0, 0,
	}
	x, y = testLook.PoliticalMap.Center()
	testLook.CenterX = x
	testLook.CenterY = y

	if err := testdb.Create(&testLook).Error; err != nil {
		panic(err)
	}
	defer testdb.Delete(&testLook)

	// CREATE EVENT
	b := CreatorBinding{
		StartDate: time.Now().Add(time.Hour),
		EndDate:   time.Now().Add(time.Hour * 2),
		Location:  "197 Hull St. #3 Brooklyn, NY 11233",
		Summary:   "test event",
	}

	testEvent, sub, errs := New(testHost, b)
	if len(errs) > 0 {
		panic(errs)
	}

	if err := testdb.Create(&testEvent).Error; err != nil {
		panic(err)
	}
	defer testdb.Delete(&testEvent)

	if err := testdb.Create(&sub).Error; err != nil {
		panic(err)
	}
	defer testdb.Delete(&sub)

	testID = testEvent.ID

	// SUBSCRIBE GUEST
	guest, _ := testEvent.NewGuest(testGuest, RTGoing)
	if err := testdb.Create(&guest).Error; err != nil {
		panic(err)
	}
	defer testdb.Delete(&guest)

	m.Run()
}

func TestGetByID(t *testing.T) {
	e, err := GetByID(testID, testLook, testdb)
	if err != nil {
		t.Error("Unexpected error:", err)
		return
	}

	if e.ID != testID {
		t.Error("Expected ID to be:", testID, "got:", e.ID)
	}

	if e.RSVP != RTNone {
		t.Error("Expected user to not be subscribed")
	}

	found := false
	for _, host := range e.Hosts {
		if host.GetID() == testHost.ID {
			found = true
		}
	}

	if !found {
		t.Error("Expected to find testHost in hosts")
	}

	if e.GoingCount < 1 {
		t.Error("Expected at least one person going")
	}

	if e.MaybeCount > 0 {
		t.Error("Didn't expect anyone to be maybe")
	}
}

func TestSearchForUser(t *testing.T) {
	es, err := SearchForUser(testLook, 0, testdb)
	if err != nil {
		t.Error("Unexpected error:", err)
		return
	}

	if len(es) == 0 {
		t.Error("Expected at least one event")
		return
	}

	found := false
	for _, e := range es {
		if e.ID == testID {
			found = true
		}
	}

	if !found {
		t.Error("Expected test event to show")
	}
}

func TestGetByHost(t *testing.T) {
	es, err := GetByHost(testHost, 0, testdb)
	if err != nil {
		t.Error("Unexpected error:", err)
		return
	}

	if len(es) == 0 {
		t.Error("Expected at least one event")
		return
	}

	found := false
	for _, e := range es {
		if e.ID == testID {
			found = true
		}

		foundHost := false
		for _, host := range e.Hosts {
			if host.GetID() == testHost.ID {
				foundHost = true
			}
		}

		if !foundHost {
			t.Error("Expected for find testHost in event:", e)
		}

	}

	if !found {
		t.Error("Expected test event to show")
	}
}

func TestGetByGuest(t *testing.T) {
	es, err := GetByGuest(testGuest, 0, testdb)
	if err != nil {
		t.Error("Unexpected error:", err)
		return
	}

	if len(es) == 0 {
		t.Error("Expected at least one event")
		return
	}

	found := false
	for _, e := range es {
		if e.ID == testID {
			found = true
		}

		if e.GoingCount != 2 {
			t.Error("Expected testHost and testGuest to be going")
		}
	}

	if !found {
		t.Error("Expected test event to show")
	}

}
