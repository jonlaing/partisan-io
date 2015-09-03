package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	// "log"
	"os"
	"testing"
)

func TestNewFeedItem(t *testing.T) {
	dbURL := os.Getenv("DB_URL")

	dbSession, err := mgo.Dial(dbURL)
	if err != nil {
		panic(err)
	}

	mIO := InitModels(dbSession, "partisan_test")
	db := mIO.DB

	uid := bson.NewObjectId()
	p := Post{
		ID:     bson.NewObjectId(),
		UserID: uid,
		Body:   "Test Post",
	}
	err = NewFeedItem(mIO, &p, uid, FeedCreate)
	if err != nil {
		t.Errorf("Couldn't insert New Feed Item:", err)
	}

	f := FeedItem{}
	err = db.C("feed_items").Find(nil).Sort("-_id").One(&f)
	if err != nil {
		t.Errorf("Couldn't find feed item")
	}

	if f.RecordID != p.ID {
		t.Log(f)
		t.Fail()
	}

	if f.Action != FeedCreate {
		t.Log(f)
		t.Fail()
	}
}
