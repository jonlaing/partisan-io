package dao

import (
	m "partisan/models"
	"strconv"
	"testing"
	"time"
)

func TestGetMatchesCenter(t *testing.T) {
	user := m.User{CenterX: 0, CenterY: 0}
	var users []m.User
	for i := 0; i < 10; i++ {
		u := m.User{ID: uint64(i), CenterX: 100 / (10 - i), CenterY: 0}
		users = append(users, u)
		db.Create(&u)
	}
	defer db.Delete(&users)

	matches, err := GetMatches(user, "", -1, -1, float64(90), 1, &db)
	if err != nil {
		t.Error(err)
	}

	if len(matches) != 8 {
		t.Error("Expected 8 matches, got:", len(matches))
	}
}

func TestGetMatchesLocation(t *testing.T) {
	user := m.User{Latitude: 0, Longitude: 0, CenterX: 0, CenterY: 0}
	var users []m.User
	for i := 0; i < 10; i++ {
		u := m.User{ID: uint64(i), Latitude: float64(45 / (10 - i)), Longitude: 0, CenterX: 0, CenterY: 0}
		users = append(users, u)
		db.Create(&u)
	}
	defer db.Delete(&users)

	matches, err := GetMatches(user, "", -1, -1, float64(10), 1, &db)
	if err != nil {
		t.Error(err)
	}

	if len(matches) != 6 {
		t.Error("Expected 6 matches, got:", len(matches))
	}
}

func TestGetMatchesGender(t *testing.T) {
	user := m.User{}
	var users []m.User
	for i := 0; i < 10; i++ {
		var u m.User
		if i%2 == 0 {
			u = m.User{ID: uint64(i), Gender: "Foo"}
		} else {
			u = m.User{ID: uint64(i), Gender: "Bar"}
		}
		users = append(users, u)
		db.Create(&u)
	}
	defer db.Delete(&users)

	matches, err := GetMatches(user, "Foo", -1, -1, float64(10), 1, &db)
	if err != nil {
		t.Error(err)
	}

	if len(matches) != 5 {
		t.Error("Expected 5 matches, got:", len(matches))
	}

	for _, m := range matches {
		if m.Gender != "Foo" {
			t.Error("Expected gender to be \"Foo\", got:", m.Gender)
		}
	}
}

func TestGetMatchesAge(t *testing.T) {
	user := m.User{}
	var users []m.User
	for i := 0; i < 10; i++ {
		var u m.User
		birthdate, _ := time.Parse("2006", strconv.Itoa(time.Now().Year()-i-20))
		u = m.User{ID: uint64(i), Birthdate: birthdate}
		users = append(users, u)
		db.Create(&u)
	}
	defer db.Delete(&users)

	matches, err := GetMatches(user, "", 25, 40, float64(10), 1, &db)
	if err != nil {
		t.Error(err)
	}

	if len(matches) != 5 {
		t.Error("Expected 5 matches, got:", len(matches))
	}
}
