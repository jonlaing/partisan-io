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
}
