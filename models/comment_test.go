package models

import "testing"

func TestCommentGetID(t *testing.T) {
	c := Comment{ID: 5}

	if c.GetID() != 5 {
		t.Error("GetID should have returned 5, returned:", c.GetID())
	}
}

func TestCommentGetUserID(t *testing.T) {
	c := Comment{UserID: 5}

	if c.GetUserID() != 5 {
		t.Error("GetUserID should have returned 5, returned:", c.GetUserID())
	}
}

func TestCommentGetType(t *testing.T) {
	c := Comment{}

	if c.GetType() != "comment" {
		t.Error("GetType should have returned \"comment\", but returned:", c.GetType())
	}
}

func TestCommentGetContent(t *testing.T) {
	body := "blah blah"
	c := Comment{Body: body}

	if c.GetContent() != body {
		t.Error("GetContent should have returned \"", body, "\" but returned:", c.GetContent())
	}
}

func TestCommentGetRecordUserID(t *testing.T) {
	post := Post{ID: 5, UserID: 7}
	testDB.Create(&post)
	defer testDB.Delete(&post)

	c := Comment{PostID: 5}

	if id, err := c.GetRecordUserID(&testDB); err != nil {
		t.Error(err)
	} else {
		if id != 7 {
			t.Error("ID should have been 7, instead was:", id)
		}
	}
}
