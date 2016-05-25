package messages

import (
	"testing"
	"time"

	"github.com/nu7hatch/gouuid"
	"partisan/models.v2/users"
)

func spoofUUID() string {
	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	return id.String()
}

func TestMessageUnreadCount(t *testing.T) {
	u := users.User{ID: spoofUUID()}
	var msgs []Message
	var mtus []ThreadUser
	var threads []Thread

	for i := 1; i <= 10; i++ {
		threadID := spoofUUID()
		if i%2 == 0 {
			threadID = u.ID // just as a baseline
		}

		id := spoofUUID()

		msg1 := Message{ThreadID: threadID, UserID: id, Read: i%2 == 0}
		mtu1 := ThreadUser{UserID: id, ThreadID: threadID}
		msg2 := Message{ThreadID: threadID, UserID: u.ID, Read: true}
		mtu2 := ThreadUser{UserID: u.ID, ThreadID: threadID}
		thread := Thread{ID: threadID}
		testdb.Create(&msg1)
		testdb.Create(&mtu1)
		testdb.Create(&msg2)
		testdb.Create(&mtu2)
		if threadID != u.ID {
			testdb.Create(&thread)
			threads = append(threads, thread)
		}
		msgs = append(msgs, msg1, msg2)
		mtus = append(mtus, mtu1, mtu2)
	}
	defer testdb.Delete(&msgs)
	defer testdb.Delete(&mtus)
	defer testdb.Delete(&threads)

	c, err := UnreadCount(u.ID, testdb)
	if err != nil {
		t.Error(err)
	}

	if c != 5 {
		t.Error("Expected 5, got:", c)
	}
}

func TestGetMessages(t *testing.T) {
	thread := Thread{ID: spoofUUID()}
	otherThreadID := spoofUUID()

	var msgs []Message
	var mtus []ThreadUser

	for i := 1; i <= 10; i++ {
		id := spoofUUID()

		var tID string
		if i%2 == 0 {
			tID = thread.ID
		} else {
			tID = otherThreadID
		}

		msg := Message{ID: id, UserID: id, ThreadID: tID}
		mtu := ThreadUser{UserID: id, ThreadID: tID}
		testdb.Create(&msg)
		testdb.Create(&mtu)
		msgs = append(msgs, msg)
		mtus = append(mtus, mtu)
	}
	defer testdb.Delete(&msgs)
	defer testdb.Delete(&mtus)

	ms, err := GetMessages(thread.ID, testdb)
	if err != nil {
		t.Error(err)
	}

	if len(ms) != 5 {
		t.Error("Expected 5, got:", len(ms))
	}

	ms, err = GetMessages(spoofUUID(), testdb)
	if err != nil {
		t.Error(err)
	}

	if len(ms) != 0 {
		t.Error("Didn't expect to find any messages")
	}
}

func TestGetMessagesAfter(t *testing.T) {
	var msgs []Message
	var mtus []ThreadUser

	threadID := spoofUUID()

	for i := 1; i <= 10; i++ {
		id := spoofUUID()

		msg := Message{ID: id, UserID: id, ThreadID: threadID}
		mtu := ThreadUser{UserID: id, ThreadID: threadID}
		testdb.Create(&msg)
		testdb.Create(&mtu)
		if i%2 == 0 {
			msg.CreatedAt = time.Time{}
			testdb.Save(&msg)
		}
		msgs = append(msgs, msg)
		mtus = append(mtus, mtu)
	}
	defer testdb.Delete(&msgs)
	defer testdb.Delete(&mtus)

	after := time.Now().AddDate(0, 0, -1)
	ms, err := GetMessagesAfter(threadID, after, testdb)
	if err != nil {
		t.Error(err)
	}

	if len(ms) != 5 {
		t.Error("Expected 5, got:", len(ms))
	}

	ms, err = GetMessages(spoofUUID(), testdb)
	if err != nil {
		t.Error(err)
	}

	if len(ms) != 0 {
		t.Error("Didn't expect to find any messages")
	}
}

// This changed, so I'm going to have to rewrite the test
// func TestCollectMessageUsers(t *testing.T) {
// 	var ms []Message
// 	var us []users.User

// 	for i := 0; i < 10; i++ {
// 		id := spoofUUID()

// 		ms = append(ms, m.Message{UserID: id})
// 		us = append(us, m.User{ID: id})
// 	}

// 	collectMessageUsers(ms, us)

// 	for _, msg := range ms {
// 		found := false
// 		for _, u := range us {
// 			if u.ID == msg.User.ID {
// 				found = true
// 			}
// 		}

// 		if !found {
// 			t.Error("Expected User to match up:", msg)
// 		}
// 	}
// }
