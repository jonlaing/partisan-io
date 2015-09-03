package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	// "log"
	"os"
	c "partisan-api/models/concerns"
	"testing"
)

func TestModelName(t *testing.T) {
	p := Post{}
	f := FeedItem{}
	ps := []Post{}
	fs := []FeedItem{}

	if n := modelName(p); n != "post" {
		t.Log(n)
		t.Fail()
	}

	if n := modelName(&p); n != "post" {
		t.Log(n)
		t.Fail()
	}

	if n := modelName(ps); n != "post" {
		t.Log(n)
		t.Fail()
	}

	if n := modelName(&ps); n != "post" {
		t.Log(n)
		t.Fail()
	}

	if n := modelName(f); n != "feed_item" {
		t.Log(n)
		t.Fail()
	}

	if n := modelName(&f); n != "feed_item" {
		t.Log(n)
		t.Fail()
	}

	if n := modelName(fs); n != "feed_item" {
		t.Log(n)
		t.Fail()
	}

	if n := modelName(&fs); n != "feed_item" {
		t.Log(n)
		t.Fail()
	}
}

func TestCollectionName(t *testing.T) {
	p := Post{}
	ps := []Post{}
	if n := collectionName(p); n != "posts" {
		t.Log(n)
		t.Fail()
	}

	if n := collectionName(&p); n != "posts" {
		t.Log(n)
		t.Fail()
	}

	if n := collectionName(ps); n != "posts" {
		t.Log(n)
		t.Fail()
	}

	if n := collectionName(&ps); n != "posts" {
		t.Log(n)
		t.Fail()
	}
}

func TestFindingRecord(t *testing.T) {
	dbURL := os.Getenv("DB_URL")

	dbSession, err := mgo.Dial(dbURL)
	if err != nil {
		panic(err)
	}

	mIO := InitModels(dbSession, "partisan_test")

	p := Post{
		ID:       bson.NewObjectId(),
		UserID:   bson.NewObjectId(),
		Body:     "Test Post",
		Likes:    c.Likes{},
		Dislikes: c.Dislikes{},
		Stars:    c.Stars{},
	}
	p.SetCreateTime()

	err = mIO.DB.C("posts").Insert(p)
	if err != nil {
		t.Errorf("Couldn't put in test post: ", err)
	}

	p1 := Post{}
	err = mIO.Find(p.ID, &p1)
	if err != nil {
		t.Errorf("Couldn't find post (", p.ID, "):", err)
	}

	if p1.Body != "Test Post" {
		t.Fail()
	}
}

func TestInsertingRecord(t *testing.T) {
	dbURL := os.Getenv("DB_URL")

	dbSession, err := mgo.Dial(dbURL)
	if err != nil {
		panic(err)
	}

	mIO := InitModels(dbSession, "partisan_test")

	p := Post{
		ID:       bson.NewObjectId(),
		UserID:   bson.NewObjectId(),
		Body:     "Test Post",
		Likes:    c.Likes{},
		Dislikes: c.Dislikes{},
		Stars:    c.Stars{},
	}
	p.SetCreateTime()

	err = mIO.Insert(&p)
	if err != nil {
		t.Errorf("Couldn't insert test post: ", err)
	}

	p1 := Post{}
	err = mIO.DB.C("posts").FindId(p.ID).One(&p1)
	if err != nil {
		t.Errorf("Couldnt find post (", p.ID, "):", err)
	}

	if p1.Body != "Test Post" {
		t.Fail()
	}
}

func TestUpdateRecord(t *testing.T) {
	dbURL := os.Getenv("DB_URL")

	dbSession, err := mgo.Dial(dbURL)
	if err != nil {
		panic(err)
	}

	mIO := InitModels(dbSession, "partisan_test")

	p := Post{
		ID:       bson.NewObjectId(),
		UserID:   bson.NewObjectId(),
		Body:     "Test Post",
		Likes:    c.Likes{},
		Dislikes: c.Dislikes{},
		Stars:    c.Stars{},
	}
	p.SetCreateTime()

	err = mIO.DB.C("posts").Insert(p)
	if err != nil {
		t.Errorf("Couldn't put in test post: ", err)
	}

	p1 := Post{}
	err = mIO.DB.C("posts").FindId(p.ID).One(&p1)
	if err != nil {
		t.Errorf("Couldn't find test post:", err)
	}

	p1.Body = "Updated Post"
	err = mIO.Update(&p1)
	if err != nil {
		t.Errorf("Couldn't update test post:", err)
	}

	p2 := Post{}
	err = mIO.DB.C("posts").FindId(p.ID).One(&p2)
	if err != nil {
		t.Errorf("Couldn't find test post again:", err)
	}

	if p2.Body != "Updated Post" {
		t.Fail()
	}
}

func TestDestroyRecord(t *testing.T) {
	dbURL := os.Getenv("DB_URL")

	dbSession, err := mgo.Dial(dbURL)
	if err != nil {
		panic(err)
	}

	mIO := InitModels(dbSession, "partisan_test")

	p := Post{
		ID:       bson.NewObjectId(),
		UserID:   bson.NewObjectId(),
		Body:     "Test Post",
		Likes:    c.Likes{},
		Dislikes: c.Dislikes{},
		Stars:    c.Stars{},
	}
	p.SetCreateTime()

	err = mIO.DB.C("posts").Insert(p)
	if err != nil {
		t.Errorf("Couldn't put in test post: ", err)
	}

	p1 := Post{}
	err = mIO.DB.C("posts").FindId(p.ID).One(&p1)
	if err != nil {
		t.Errorf("Couldn't find test post:", err)
	}

	err = mIO.Destroy(&p1)
	if err != nil {
		t.Errorf("Couldn't destroy test post:", err)
	}

	p2 := Post{}
	err = mIO.DB.C("posts").FindId(p.ID).One(&p2)
	if err == nil {
		t.Fail()
	}
}
