package dao

import (
	m "partisan/models"
	"testing"
)

func TestGetFeedByUserIDs(t *testing.T) {
	var userIDs []uint64
	for i := 0; i <= 30; i++ {
		userIDs = append(userIDs, uint64(i))
		db.Create(&m.FeedItem{UserID: uint64(i), Action: "post", RecordID: uint64(i)})
	}

	fs, err := GetFeedByUserIDs(uint64(1), userIDs, 1, &db)
	if err != nil {
		t.Error(err)
	}

	if len(fs) != 25 {
		t.Error("Wrong number of feed items:", len(fs))
	}

	for _, f := range fs {
		found := false
		for _, v := range userIDs {
			if f.UserID == v {
				found = true
			}
		}
		if !found {
			t.Error("Couldn't find id:", f.ID, "in", userIDs)
		}
	}

	db.Delete(&fs)
}

func TestGetFeedByUserIDsPage2(t *testing.T) {
	var userIDs []uint64
	for i := 0; i <= 30; i++ {
		userIDs = append(userIDs, uint64(i))
		db.Create(&m.FeedItem{UserID: uint64(i), Action: "post", RecordID: uint64(i)})
	}

	fs, err := GetFeedByUserIDs(uint64(1), userIDs, 2, &db)
	if err != nil {
		t.Error(err)
	}

	if len(fs) != 6 {
		t.Error("Wrong number of feed items:", len(fs))
	}

	for _, f := range fs {
		found := false
		for _, v := range userIDs {
			if f.UserID == v {
				found = true
			}
		}
		if !found {
			t.Error("Couldn't find id:", f.ID, "in", userIDs)
		}
	}

	db.Delete(&fs)
}
