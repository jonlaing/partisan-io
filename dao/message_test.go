package dao

import (
	m "partisan/models"
	"testing"
	"time"
)

func TestMessageUnreadCount(t *testing.T) {
	u := m.User{ID: 0}
	var msgs []m.Message
	var threads []m.MessageThread

	for i := 1; i <= 10; i++ {
		threadID := i
		if i%2 == 0 {
			threadID = 0
		}

		msg1 := m.Message{UserID: uint64(i), ThreadID: uint64(threadID), Read: i%2 == 0}
		msg2 := m.Message{UserID: u.ID, ThreadID: uint64(threadID), Read: true}
		thread := m.MessageThread{ID: uint64(i)}
		db.Create(&msg1)
		db.Create(&msg2)
		db.Create(&thread)
		msgs = append(msgs, msg1, msg2)
		threads = append(threads, thread)
	}
	defer db.Delete(&msgs)
	defer db.Delete(&threads)

	c, err := MessageUnreadCount(u.ID, &db)
	if err != nil {
		t.Error(err)
	}

	if c != 5 {
		t.Error("Expected 5, got:", c)
	}
}

func TestGetMessages(t *testing.T) {
	thread := m.MessageThread{ID: 1}

	var msgs []m.Message

	for i := 1; i <= 10; i++ {
		// UserID will be either 1 or 0
		msg := m.Message{ID: uint64(i), ThreadID: uint64(i % 2)}
		db.Create(&msg)
		msgs = append(msgs, msg)
	}
	defer db.Delete(&msgs)

	ms, err := GetMessages(thread.ID, &db)
	if err != nil {
		t.Error(err)
	}

	if len(ms) != 5 {
		t.Error("Expected 5, got:", len(ms))
	}

	ms, err = GetMessages(2, &db)
	if err != nil {
		t.Error(err)
	}

	if len(ms) != 0 {
		t.Error("Didn't expect to find any messages")
	}
}

func TestGetMessagesAfter(t *testing.T) {
	var msgs []m.Message

	for i := 1; i <= 10; i++ {
		msg := m.Message{ID: uint64(i), ThreadID: 1}
		db.Create(&msg)
		if i%2 == 0 {
			msg.CreatedAt = time.Time{}
			db.Save(&msg)
		}
		msgs = append(msgs, msg)
	}

	after := time.Now().AddDate(0, 0, -1)
	ms, err := GetMessagesAfter(1, after, &db)
	if err != nil {
		t.Error(err)
	}

	if len(ms) != 5 {
		t.Error("Expected 5, got:", len(ms))
	}

	ms, err = GetMessages(2, &db)
	if err != nil {
		t.Error(err)
	}

	if len(ms) != 0 {
		t.Error("Didn't expect to find any messages")
	}
}

func TestCollectMessageUsers(t *testing.T) {
	var ms []m.Message
	var us []m.User

	for i := 0; i < 10; i++ {
		ms = append(ms, m.Message{UserID: uint64(i)})
		us = append(us, m.User{ID: uint64(i)})
	}

	collectMessageUsers(ms, us)

	for _, msg := range ms {
		found := false
		for _, u := range us {
			if u.ID == msg.User.ID {
				found = true
			}
		}

		if !found {
			t.Error("Expected User to match up:", msg)
		}
	}
}
