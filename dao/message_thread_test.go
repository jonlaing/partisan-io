package dao

import (
	m "partisan/models"
	"testing"
)

func TestGetMessageThread(t *testing.T) {
	thread := m.MessageThread{ID: 1}
	db.Create(&thread)
	defer db.Delete(&thread)

	mtDAO := MessageThreadDAO{db}

	_, err := mtDAO.Get(1)
	if err != nil {
		t.Error(err)
	}

	_, err = mtDAO.Get(2)
	if err == nil {
		t.Error("Expected an error")
	}

	if _, ok := err.(*ErrThreadNotFound); !ok {
		t.Error("Expected ErrThreadNotFound, but got:", err)
	}
}

func TestGetMessageThreads(t *testing.T) {
	u := m.User{ID: 1}
	var mtus []m.MessageThreadUser
	var threads []m.MessageThread

	for i := 0; i < 10; i++ {
		// UserID will be either 1 or 0
		mtu := m.MessageThreadUser{UserID: uint64(i % 2), ThreadID: uint64(i)}
		thread := m.MessageThread{ID: uint64(i)}
		db.Create(&mtu)
		db.Create(&thread)
		mtus = append(mtus, mtu)
		threads = append(threads, thread)
	}
	defer db.Delete(&mtus)
	defer db.Delete(&threads)

	mtDAO := MessageThreadDAO{db}

	ts, err := mtDAO.List(u.ID)
	if err != nil {
		t.Error(err)
	}

	if len(ts) != 5 {
		t.Error("Expected 5 threads, got:", len(ts))
	}
}

func TestGetMessageThreadIDs(t *testing.T) {
	u := m.User{ID: 1}
	var mtus []m.MessageThreadUser
	var threads []m.MessageThread

	for i := 0; i < 10; i++ {
		// UserID will be either 1 or 0
		mtu := m.MessageThreadUser{UserID: uint64(i % 2), ThreadID: uint64(i)}
		thread := m.MessageThread{ID: uint64(i)}
		db.Create(&mtu)
		db.Create(&thread)
		mtus = append(mtus, mtu)
		threads = append(threads, thread)
	}
	defer db.Delete(&mtus)
	defer db.Delete(&threads)

	mtDAO := MessageThreadDAO{db}

	ts, err := mtDAO.GetIDs(u.ID)
	if err != nil {
		t.Error(err)
	}

	if len(ts) != 5 {
		t.Error("Expected 5 threads, got:", len(ts))
	}
}

func TestMessageThreadHasUnread(t *testing.T) {
	m1 := m.Message{UserID: 1, Read: false, ThreadID: 1}
	mtu1 := m.MessageThreadUser{UserID: 1, ThreadID: 1}
	m2 := m.Message{UserID: 2, Read: true, ThreadID: 2}
	mtu2 := m.MessageThreadUser{UserID: 2, ThreadID: 2}
	db.Create(&m1)
	db.Create(&mtu1)
	db.Create(&m2)
	db.Create(&mtu2)
	defer db.Delete(&m1)
	defer db.Delete(&m2)
	defer db.Delete(&mtu1)
	defer db.Delete(&mtu2)

	mtDAO := MessageThreadDAO{db}

	unread, err := mtDAO.HasUnread(2, 1)
	if err != nil {
		t.Error(err)
	}

	if !unread {
		t.Error("Shouldn't have any unread messages")
	}

	unread, err = mtDAO.HasUnread(2, 2)
	if err != nil {
		t.Error(err)
	}

	if unread {
		t.Error("Should have unread messages")
	}
}

func TestMessageThreadHasUser(t *testing.T) {
	m1 := m.Message{UserID: 1}
	mtu1 := m.MessageThreadUser{ThreadID: 1, UserID: 1}
	m2 := m.Message{UserID: 2}
	mtu2 := m.MessageThreadUser{ThreadID: 2, UserID: 2}
	db.Create(&m1)
	db.Create(&mtu1)
	db.Create(&m2)
	db.Create(&mtu2)
	defer db.Delete(&m1)
	defer db.Delete(&m2)
	defer db.Delete(&mtu1)
	defer db.Delete(&mtu2)

	mtDAO := MessageThreadDAO{db}

	hasUser, err := mtDAO.HasUser(1, 1)
	if err != nil {
		t.Error(err)
	}

	if !hasUser {
		t.Error("Expected user in thread 1")
	}

	hasUser, err = mtDAO.HasUser(1, 2)
	if err != nil {
		t.Error(err)
	}

	if hasUser {
		t.Error("Didn't expect user in thread 2")
	}
}

func TestMessageThreadByUsers(t *testing.T) {
	t1 := m.MessageThread{ID: 1}
	t2 := m.MessageThread{ID: 2}
	mtu1 := m.MessageThreadUser{UserID: 1, ThreadID: 1}
	mtu2 := m.MessageThreadUser{UserID: 2, ThreadID: 1}
	mtu3 := m.MessageThreadUser{UserID: 2, ThreadID: 2}
	db.Create(&t1)
	db.Create(&t2)
	db.Create(&mtu1)
	db.Create(&mtu2)
	db.Create(&mtu3)
	defer db.Delete(&t1)
	defer db.Delete(&t2)
	defer db.Delete(&mtu1)
	defer db.Delete(&mtu2)
	defer db.Delete(&mtu3)

	mtDAO := MessageThreadDAO{db}

	_, err := mtDAO.GetByUsers(1, 2)
	if err != nil {
		t.Error(err)
	}

	th, err := mtDAO.GetByUsers(3, 4)
	if err == nil {
		t.Error("Expected an error:", th)
	}

	th, err = mtDAO.GetByUsers(2, 3)
	if err == nil {
		t.Error("Expected an error:", th)
	}

	if _, ok := err.(*MessageThreadUnreciprocated); !ok {
		t.Error("Should have been MessageThreadUnreciprocated error")
	}
}
