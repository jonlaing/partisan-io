package dao

import (
	m "partisan/models"
	"testing"
)

func TestGetMessageThreads(t *testing.T) {
	u := m.User{ID: 1}
	var msgs []m.Message
	var threads []m.MessageThread

	for i := 0; i < 10; i++ {
		// UserID will be either 1 or 0
		msg := m.Message{UserID: uint64(i % 2), ThreadID: uint64(i)}
		thread := m.MessageThread{ID: uint64(i)}
		db.Create(&msg)
		db.Create(&thread)
		msgs = append(msgs, msg)
		threads = append(threads, thread)
	}
	defer db.Delete(&msgs)
	defer db.Delete(&threads)

	ts, err := GetMessageThreads(u.ID, &db)
	if err != nil {
		t.Error(err)
	}

	if len(ts) != 5 {
		t.Error("Expected 5 threads, got:", len(ts))
	}
}

func TestGetMessageThreadIDs(t *testing.T) {
	u := m.User{ID: 1}
	var msgs []m.Message
	var threads []m.MessageThread

	for i := 0; i < 10; i++ {
		// UserID will be either 1 or 0
		msg := m.Message{UserID: uint64(i % 2), ThreadID: uint64(i)}
		thread := m.MessageThread{ID: uint64(i)}
		db.Create(&msg)
		db.Create(&thread)
		msgs = append(msgs, msg)
		threads = append(threads, thread)
	}
	defer db.Delete(&msgs)
	defer db.Delete(&threads)

	ts, err := GetMessageThreadIDs(u.ID, &db)
	if err != nil {
		t.Error(err)
	}

	if len(ts) != 5 {
		t.Error("Expected 5 threads, got:", len(ts))
	}
}

func TestMessageThreadHasUnread(t *testing.T) {
	m1 := m.Message{ThreadID: 1, Read: false}
	m2 := m.Message{ThreadID: 2, Read: true}
	db.Create(&m1)
	db.Create(&m2)
	defer db.Delete(&m1)
	defer db.Delete(&m2)

	unread, err := MessageThreadHasUnread(1, &db)
	if err != nil {
		t.Error(err)
	}

	if !unread {
		t.Error("Shouldn't have any unread messages")
	}

	unread, err = MessageThreadHasUnread(2, &db)
	if err != nil {
		t.Error(err)
	}

	if unread {
		t.Error("Should have unread messages")
	}
}

func TestMessageThreadHasUser(t *testing.T) {
	m1 := m.Message{ThreadID: 1, UserID: 1}
	m2 := m.Message{ThreadID: 2, UserID: 2}
	db.Create(&m1)
	db.Create(&m2)
	defer db.Delete(&m1)
	defer db.Delete(&m2)

	hasUser, err := MessageThreadHasUser(1, 1, &db)
	if err != nil {
		t.Error(err)
	}

	if !hasUser {
		t.Error("Expected user in thread 1")
	}

	hasUser, err = MessageThreadHasUser(1, 2, &db)
	if err != nil {
		t.Error(err)
	}

	if hasUser {
		t.Error("Didn't expect user in thread 2")
	}
}
