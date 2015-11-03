package dao

import (
	m "partisan/models"
	"testing"
	"time"
)

func TestGetFeedByUserIDs(t *testing.T) {
	var userIDs []uint64
	var items []m.FeedItem
	for i := 0; i <= 30; i++ {
		userIDs = append(userIDs, uint64(i))
		item := m.FeedItem{UserID: uint64(i), Action: "post", RecordID: uint64(i)}
		db.Create(&item)
		items = append(items, item)
	}
	defer db.Delete(&items)

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
}

func TestGetFeedByUserIDsPage2(t *testing.T) {
	var userIDs []uint64
	var items []m.FeedItem
	for i := 0; i <= 30; i++ {
		userIDs = append(userIDs, uint64(i))
		item := m.FeedItem{UserID: uint64(i), Action: "post", RecordID: uint64(i)}
		db.Create(&item)
		items = append(items, item)
	}
	defer db.Delete(&items)

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
}

func TestGetFeedByUserIDsAfter(t *testing.T) {
	var userIDs []uint64
	var items []m.FeedItem
	for i := 0; i < 10; i++ {
		userIDs = append(userIDs, uint64(i))
		item := m.FeedItem{UserID: uint64(i), Action: "post", RecordID: uint64(i)}
		db.Create(&item)
		if i%2 == 0 {
			item.CreatedAt = time.Unix(int64(0), int64(0))
			db.Save(&item)
		}
		items = append(items, item)
	}
	defer db.Delete(&items)

	after := time.Now().AddDate(0, 0, -1)
	fs, err := GetFeedByUserIDsAfter(uint64(1), userIDs, after, &db)
	if err != nil {
		t.Error(err)
	}

	if len(fs) != 5 {
		t.Error("Expected 5 items, got:", len(fs))
	}

	for _, f := range fs {
		if f.CreatedAt == time.Unix(int64(0), int64(0)) {
			t.Error("Expected no items to be zero time")
		}
	}
}
