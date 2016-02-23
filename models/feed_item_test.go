package models

import "testing"

func TestGetPostIDs(t *testing.T) {
	fs := FeedItems{
		FeedItem{RecordType: "post", RecordID: 1},
		FeedItem{RecordType: "post", RecordID: 2},
		FeedItem{RecordType: "comment", RecordID: 3},
		FeedItem{RecordType: "post", RecordID: 4},
	}

	ids := fs.GetPostIDs()

	if len(ids) != 3 {
		t.Error("Wrong number of ids. Should be 3, was:", len(ids))
	}

	for _, v := range ids {
		if v == 3 {
			t.Error("FeedItem with RecordID 3 should not have been in list")
		}
	}
}
