package models

import "testing"

func TestLikeGetID(t *testing.T) {
	l := Like{ID: 5}

	if l.GetID() != 5 {
		t.Error("Expected ID to be 5, got:", l.GetID())
	}
}

func TestLikeGetType(t *testing.T) {
	l := Like{}

	if l.GetType() != "like" {
		t.Fail()
	}
}

func TestLikeGetRecordUserID(t *testing.T) {
	post := Post{ID: 5, UserID: 6}
	if err := testDB.Create(&post).Error; err != nil {
		t.Error(err)
	}
	defer testDB.Delete(&post)

	l := Like{RecordType: "post", RecordID: 5}

	if id, err := l.GetRecordUserID(testDB); err != nil {
		t.Error(err)
	} else {
		if id != 6 {
			t.Error("Expected ID to be 6, got:", id)
		}
	}
}
