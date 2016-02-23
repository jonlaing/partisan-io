package models

import "testing"

func TestFriendshipGetID(t *testing.T) {
	f := Friendship{ID: 5}

	if f.GetID() != 5 {
		t.Error("GetID should have returned 5, but returned:", f.GetID())
	}
}

func TestFriendshipGetRecordUserID(t *testing.T) {
	f := Friendship{
		UserID:    4,
		FriendID:  5,
		Confirmed: false,
	}

	// ignoring error, since it's impossible to return one in this implementation
	if id, _ := f.GetRecordUserID(&testDB); id != 5 {
		t.Error("GetRecordUserID should have returned 5, but returned:", id)
	}

	f.Confirmed = true

	// ignoring error, since it's impossible to return one in this implementation
	if id, _ := f.GetRecordUserID(&testDB); id != 4 {
		t.Error("GetRecordUserID should have returned 4, but returned:", id)
	}
}

func TestFriendshipGetType(t *testing.T) {
	f := Friendship{}

	if f.GetType() != "friendship" {
		t.Fail()
	}
}
