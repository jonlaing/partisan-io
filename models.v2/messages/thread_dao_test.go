package messages

import (
	"testing"

	"partisan/models.v2/users"
)

func TestGetMessageThread(t *testing.T) {
	id := spoofUUID()
	thread := Thread{ID: id}
	testdb.Create(&thread)
	defer testdb.Delete(&thread)

	_, err := GetThread(id, testdb)
	if err != nil {
		t.Error(err)
	}

	_, err = GetThread(spoofUUID(), testdb)
	if err == nil {
		t.Error("Expected an error")
	}
}

func TestGetMessageThreads(t *testing.T) {
	u := users.User{ID: spoofUUID()}
	otherUser := spoofUUID()
	var mtus []ThreadUser
	var threads []Thread

	for i := 0; i < 10; i++ {
		threadID := spoofUUID()
		userID := u.ID
		if i%2 == 0 {
			userID = otherUser
		}

		mtu := ThreadUser{UserID: userID, ThreadID: threadID}
		thread := Thread{ID: threadID}
		testdb.Create(&mtu)
		testdb.Create(&thread)
		mtus = append(mtus, mtu)
		threads = append(threads, thread)
	}
	defer testdb.Delete(&mtus)
	defer testdb.Delete(&threads)

	ts, err := ListThreads(u.ID, testdb)
	if err != nil {
		t.Error(err)
	}

	if len(ts) != 5 {
		t.Error("Expected 5 threads, got:", len(ts))
	}
}

func TestGetMessageThreadIDs(t *testing.T) {
	u := users.User{ID: spoofUUID()}
	otherUser := spoofUUID()
	var mtus []ThreadUser
	var threads []Thread

	for i := 0; i < 10; i++ {
		threadID := spoofUUID()
		userID := u.ID
		if i%2 == 0 {
			userID = otherUser
		}

		mtu := ThreadUser{UserID: userID, ThreadID: threadID}
		thread := Thread{ID: threadID}
		testdb.Create(&mtu)
		testdb.Create(&thread)
		mtus = append(mtus, mtu)
		threads = append(threads, thread)
	}
	defer testdb.Delete(&mtus)
	defer testdb.Delete(&threads)

	ts, err := GetThreadIDs(u.ID, testdb)
	if err != nil {
		t.Error(err)
	}

	if len(ts) != 5 {
		t.Error("Expected 5 threads, got:", len(ts))
	}
}

func TestMessageThreadHasUnread(t *testing.T) {
	user1 := spoofUUID()
	user2 := spoofUUID()
	thread1 := spoofUUID()
	thread2 := spoofUUID()

	m1 := Message{UserID: user1, Read: false, ThreadID: thread1}
	mtu1 := ThreadUser{UserID: user1, ThreadID: thread1}
	m2 := Message{UserID: user2, Read: true, ThreadID: thread2}
	mtu2 := ThreadUser{UserID: user2, ThreadID: thread2}
	testdb.Create(&m1)
	testdb.Create(&mtu1)
	testdb.Create(&m2)
	testdb.Create(&mtu2)
	defer testdb.Delete(&m1)
	defer testdb.Delete(&m2)
	defer testdb.Delete(&mtu1)
	defer testdb.Delete(&mtu2)

	unread, err := HasUnread(user2, thread1, testdb)
	if err != nil {
		t.Error(err)
	}

	if !unread {
		t.Error("Shouldn't have any unread messages")
	}

	unread, err = HasUnread(user2, thread2, testdb)
	if err != nil {
		t.Error(err)
	}

	if unread {
		t.Error("Should have unread messages")
	}
}

func TestMessageThreadHasUser(t *testing.T) {
	user1 := spoofUUID()
	user2 := spoofUUID()
	thread1 := spoofUUID()
	thread2 := spoofUUID()

	m1 := Message{UserID: user1, ThreadID: thread1}
	mtu1 := ThreadUser{ThreadID: thread1, UserID: user1}
	m2 := Message{UserID: user2, ThreadID: thread2}
	mtu2 := ThreadUser{ThreadID: thread2, UserID: user2}
	testdb.Create(&m1)
	testdb.Create(&mtu1)
	testdb.Create(&m2)
	testdb.Create(&mtu2)
	defer testdb.Delete(&m1)
	defer testdb.Delete(&m2)
	defer testdb.Delete(&mtu1)
	defer testdb.Delete(&mtu2)

	hasUser, err := HasUser(user1, thread1, testdb)
	if err != nil {
		t.Error(err)
	}

	if !hasUser {
		t.Error("Expected user in thread 1")
	}

	hasUser, err = HasUser(user1, thread2, testdb)
	if err != nil {
		t.Error(err)
	}

	if hasUser {
		t.Error("Didn't expect user in thread 2")
	}
}

func TestMessageThreadByUsers(t *testing.T) {
	// NOTE: This test seems to be very finicky, and will fail seemingly at random sometimes.
	// Not sure what the issue is, but will need to work it out eventually
	user1 := spoofUUID()
	user2 := spoofUUID()
	thread1 := spoofUUID()
	thread2 := spoofUUID()

	t1 := Thread{ID: thread1}
	t2 := Thread{ID: thread2}
	mtu1 := ThreadUser{UserID: user1, ThreadID: thread1}
	mtu2 := ThreadUser{UserID: user2, ThreadID: thread1}
	mtu3 := ThreadUser{UserID: user2, ThreadID: thread2}
	if err := testdb.Create(&t1).Error; err != nil {
		panic(err)
	}
	if err := testdb.Create(&t2).Error; err != nil {
		panic(err)
	}
	if err := testdb.Create(&mtu1).Error; err != nil {
		panic(err)
	}
	if err := testdb.Create(&mtu2).Error; err != nil {
		panic(err)
	}
	if err := testdb.Create(&mtu3).Error; err != nil {
		panic(err)
	}
	defer testdb.Delete(&t1)
	defer testdb.Delete(&t2)
	defer testdb.Delete(&mtu1)
	defer testdb.Delete(&mtu2)
	defer testdb.Delete(&mtu3)

	_, err := GetByUsers(testdb, user1, user2)
	if err != nil {
		t.Error("Uexpected error:", err)
	}

	th, err := GetByUsers(testdb, spoofUUID(), spoofUUID())
	if err == nil {
		t.Error("Expected an error:", th)
	}

	th, err = GetByUsers(testdb, user2, spoofUUID())
	if err == nil {
		t.Error("Expected an error:", th)
	}

	if err != ErrThreadUnreciprocated {
		t.Error("Should have been ErrThreadUnreciprocated error, was:", err)
	}
}
