package dao

import (
	m "partisan/models"
	"testing"
)

func TestFriends(t *testing.T) {
	u := m.User{ID: 1}
	var fships []m.Friendship
	var fs []m.User
	for i := 2; i < 12; i++ {
		f := m.User{ID: uint64(i)}
		ship := m.Friendship{UserID: u.ID, FriendID: f.ID, Confirmed: true}
		db.Create(&f)
		db.Create(&ship)
		fs = append(fs, f)
		fships = append(fships, ship)
	}
	defer db.Delete(&fships)
	defer db.Delete(&fs)

	friends, err := Friends(u, true, &db)
	if err != nil {
		t.Error(err)
	}

	if len(friends) != 10 {
		t.Error("Expected 10 friends, got:", len(friends))
	}

	for _, f := range friends {
		found := false
		for _, v := range fs {
			if f.ID == v.ID {
				found = true
			}
		}
		if !found {
			t.Error("Expected to find User ID:", f.ID, "in", fs)
		}
	}
}

func TestConfirmedFriends(t *testing.T) {
	u := m.User{ID: 1}
	var fships []m.Friendship
	var fs []m.User
	for i := 2; i < 12; i++ {
		confirmed := i%2 == 0
		f := m.User{ID: uint64(i)}
		ship := m.Friendship{UserID: u.ID, FriendID: f.ID, Confirmed: confirmed}
		db.Create(&f)
		db.Create(&ship)
		if confirmed {
			fs = append(fs, f)
			fships = append(fships, ship)
		}
	}
	defer db.Delete(&fships)
	defer db.Delete(&fs)

	friends, err := ConfirmedFriends(u, &db)
	if err != nil {
		t.Error(err)
	}

	if len(friends) != 5 {
		t.Error("Expected 5 friends, got:", len(friends))
	}

	for _, f := range friends {
		found := false
		for _, v := range fs {
			if f.ID == v.ID {
				found = true
			}
		}
		if !found {
			t.Error("Expected to find User ID:", f.ID, "in", fs)
		}
	}
}

func TestFriendIDs(t *testing.T) {
	u := m.User{ID: 1}
	var fships []m.Friendship
	var fs []m.User
	var idList []uint64
	for i := 2; i < 12; i++ {
		f := m.User{ID: uint64(i)}
		ship := m.Friendship{UserID: u.ID, FriendID: f.ID, Confirmed: true}
		db.Create(&f)
		db.Create(&ship)
		fs = append(fs, f)
		fships = append(fships, ship)
		idList = append(idList, f.ID)
	}
	defer db.Delete(&fships)
	defer db.Delete(&fs)

	friendIDs, err := FriendIDs(u, true, &db)
	if err != nil {
		t.Error(err)
	}

	if len(friendIDs) != 10 {
		t.Error("Expected 10 friends, got:", len(friendIDs))
	}

	for _, id := range friendIDs {
		found := false
		for _, v := range fs {
			if id == v.ID {
				found = true
			}
		}
		if !found {
			t.Error("Expected to find User ID:", id, "in", fs)
		}
	}
}

func TestConfirmedFriendIDs(t *testing.T) {
	u := m.User{ID: 1}
	var fships []m.Friendship
	var fs []m.User
	var idList []uint64
	for i := 2; i < 12; i++ {
		confirmed := i%2 == 0
		f := m.User{ID: uint64(i)}
		ship := m.Friendship{UserID: u.ID, FriendID: f.ID, Confirmed: confirmed}
		db.Create(&f)
		db.Create(&ship)
		if confirmed {
			fs = append(fs, f)
			fships = append(fships, ship)
			idList = append(idList, f.ID)
		}
	}
	defer db.Delete(&fships)
	defer db.Delete(&fs)

	friendIDs, err := ConfirmedFriendIDs(u, &db)
	if err != nil {
		t.Error(err)
	}

	if len(friendIDs) != 5 {
		t.Error("Expected 5 friends, got:", len(friendIDs))
	}

	for _, id := range friendIDs {
		found := false
		for _, v := range fs {
			if id == v.ID {
				found = true
			}
		}
		if !found {
			t.Error("Expected to find User ID:", id, "in", fs)
		}
	}
}

func TestGetFriendship(t *testing.T) {
	u1 := m.User{ID: uint64(1)}
	fs1 := m.Friendship{UserID: u1.ID, FriendID: uint64(2)}

	db.Create(&fs1)
	defer db.Delete(&fs1)

	fs, err := GetFriendship(u1, uint64(2), &db)
	if err != nil {
		t.Error(err)
	}

	if fs.UserID != u1.ID {
		t.Error("Expected UserID to be:", u1.ID, ", but was:", fs.UserID)
	}

	if fs.FriendID != uint64(2) {
		t.Error("Expected UserID to be: 2, but was:", fs.FriendID)
	}

	u2 := m.User{ID: uint64(3)}
	fs2 := m.Friendship{UserID: uint64(4), FriendID: u2.ID}

	db.Create(&fs2)
	defer db.Delete(&fs2)

	fs, err = GetFriendship(u2, uint64(4), &db)
	if err != nil {
		t.Error(err)
	}

	if fs.UserID != uint64(4) {
		t.Error("Expected UserID to be: 2, but was:", fs.UserID)
	}

	if fs.FriendID != u2.ID {
		t.Error("Expected FriendID to be:", u2.ID, ", but was:", fs.FriendID)
	}
}
