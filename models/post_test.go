package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	// "log"
	"os"
	"testing"
)

func TestLikesIncrement(t *testing.T) {
	dbURL := os.Getenv("DB_URL")

	dbSession, err := mgo.Dial(dbURL)
	if err != nil {
		panic(err)
	}

	mIO := InitModels(dbSession, "partisan_test")
	db := mIO.DB

	p := Post{
		ID:     bson.NewObjectId(),
		UserID: bson.NewObjectId(),
		Body:   "Test Post",
	}
	p.SetCreateTime()

	err = db.C("posts").Insert(p)
	if err != nil {
		t.Errorf("Couldn't insert test post:", err)
	}

	p1 := Post{}
	err = db.C("posts").FindId(p.ID).One(&p1)

	uid := bson.NewObjectId()
	err = p1.Like(uid)
	if err != nil {
		t.Errorf("Couldn't like post:", err)
	}

	p2 := Post{}
	err = db.C("posts").FindId(p.ID).One(&p2)

	if len(p2.Likes) != 1 {
		t.Log(len(p2.Likes))
		t.Fail()
	}

	if _, ok := p2.Likes[uid.String()]; !ok {
		t.Log(uid)
		t.Log(p1.Likes)
		t.Fail()
	}
}

func TestDislikesIncrement(t *testing.T) {
	dbURL := os.Getenv("DB_URL")

	dbSession, err := mgo.Dial(dbURL)
	if err != nil {
		panic(err)
	}

	mIO := InitModels(dbSession, "partisan_test")
	db := mIO.DB

	p := Post{
		ID:     bson.NewObjectId(),
		UserID: bson.NewObjectId(),
		Body:   "Test Post",
	}
	p.SetCreateTime()

	err = db.C("posts").Insert(p)
	if err != nil {
		t.Errorf("Couldn't insert test post:", err)
	}

	p1 := Post{}
	err = db.C("posts").FindId(p.ID).One(&p1)

	uid := bson.NewObjectId()
	err = p1.Dislike(uid)
	if err != nil {
		t.Errorf("Couldn't like post:", err)
	}

	p2 := Post{}
	err = db.C("posts").FindId(p.ID).One(&p2)

	if len(p2.Dislikes) != 1 {
		t.Log(len(p2.Dislikes))
		t.Fail()
	}

	if _, ok := p2.Dislikes[uid.String()]; !ok {
		t.Log(uid)
		t.Log(p1.Dislikes)
		t.Fail()
	}
}

func TestLikesTwiceNotIncrement(t *testing.T) {
	dbURL := os.Getenv("DB_URL")

	dbSession, err := mgo.Dial(dbURL)
	if err != nil {
		panic(err)
	}

	mIO := InitModels(dbSession, "partisan_test")
	db := mIO.DB

	p := Post{
		ID:     bson.NewObjectId(),
		UserID: bson.NewObjectId(),
		Body:   "Test Post",
	}
	p.SetCreateTime()

	err = db.C("posts").Insert(p)
	if err != nil {
		t.Errorf("Couldn't insert test post:", err)
	}

	p1 := Post{}
	err = db.C("posts").FindId(p.ID).One(&p1)

	uid := bson.NewObjectId()
	err = p1.Like(uid)
	if err != nil {
		t.Errorf("Couldn't like post:", err)
	}
	err = p1.Like(uid)
	if err != nil {
		t.Errorf("Couldn't like post twice:", err)
	}

	p2 := Post{}
	err = db.C("posts").FindId(p.ID).One(&p2)

	if len(p2.Likes) != 1 {
		t.Log(len(p2.Likes))
		t.Fail()
	}

	if _, ok := p2.Likes[uid.String()]; !ok {
		t.Log(uid)
		t.Log(p1.Likes)
		t.Fail()
	}
}

func TestDislikesTwiceNotIncrement(t *testing.T) {
	dbURL := os.Getenv("DB_URL")

	dbSession, err := mgo.Dial(dbURL)
	if err != nil {
		panic(err)
	}

	mIO := InitModels(dbSession, "partisan_test")
	db := mIO.DB

	p := Post{
		ID:     bson.NewObjectId(),
		UserID: bson.NewObjectId(),
		Body:   "Test Post",
	}
	p.SetCreateTime()

	err = db.C("posts").Insert(p)
	if err != nil {
		t.Errorf("Couldn't insert test post:", err)
	}

	p1 := Post{}
	err = db.C("posts").FindId(p.ID).One(&p1)

	uid := bson.NewObjectId()
	err = p1.Dislike(uid)
	if err != nil {
		t.Errorf("Couldn't like post:", err)
	}
	err = p1.Dislike(uid)
	if err != nil {
		t.Errorf("Couldn't like post twice:", err)
	}

	p2 := Post{}
	err = db.C("posts").FindId(p.ID).One(&p2)

	if len(p2.Dislikes) != 1 {
		t.Log(len(p2.Dislikes))
		t.Fail()
	}

	if _, ok := p2.Dislikes[uid.String()]; !ok {
		t.Log(uid)
		t.Log(p1.Dislikes)
		t.Fail()
	}
}

func TestLikeUnlikes(t *testing.T) {
	dbURL := os.Getenv("DB_URL")

	dbSession, err := mgo.Dial(dbURL)
	if err != nil {
		panic(err)
	}

	mIO := InitModels(dbSession, "partisan_test")
	db := mIO.DB

	p := Post{
		ID:     bson.NewObjectId(),
		UserID: bson.NewObjectId(),
		Body:   "Test Post",
	}
	p.SetCreateTime()

	err = db.C("posts").Insert(p)
	if err != nil {
		t.Errorf("Couldn't insert test post:", err)
	}

	p1 := Post{}
	err = db.C("posts").FindId(p.ID).One(&p1)

	uid := bson.NewObjectId()
	err = p1.Dislike(uid)
	if err != nil {
		t.Errorf("Couldn't dislike post:", err)
	}
	err = p1.Like(uid)
	if err != nil {
		t.Errorf("Couldn't like post:", err)
	}

	p2 := Post{}
	err = db.C("posts").FindId(p.ID).One(&p2)

	if len(p2.Likes) != 1 {
		t.Log(len(p2.Likes))
		t.Fail()
	}

	if len(p2.Dislikes) != 0 {
		t.Log(len(p2.Dislikes))
		t.Fail()
	}

	if _, ok := p2.Likes[uid.String()]; !ok {
		t.Log(uid)
		t.Log(p1.Likes)
		t.Fail()
	}

	if _, ok := p2.Dislikes[uid.String()]; ok {
		t.Log(uid)
		t.Log(p1.Dislikes)
		t.Fail()
	}
}

func TestDislikeUnlikes(t *testing.T) {
	dbURL := os.Getenv("DB_URL")

	dbSession, err := mgo.Dial(dbURL)
	if err != nil {
		panic(err)
	}

	mIO := InitModels(dbSession, "partisan_test")
	db := mIO.DB

	p := Post{
		ID:     bson.NewObjectId(),
		UserID: bson.NewObjectId(),
		Body:   "Test Post",
	}
	p.SetCreateTime()

	err = db.C("posts").Insert(p)
	if err != nil {
		t.Errorf("Couldn't insert test post:", err)
	}

	p1 := Post{}
	err = db.C("posts").FindId(p.ID).One(&p1)

	uid := bson.NewObjectId()
	err = p1.Like(uid)
	if err != nil {
		t.Errorf("Couldn't like post:", err)
	}
	err = p1.Dislike(uid)
	if err != nil {
		t.Errorf("Couldn't dislike post:", err)
	}

	p2 := Post{}
	err = db.C("posts").FindId(p.ID).One(&p2)

	if len(p2.Dislikes) != 1 {
		t.Log(len(p2.Dislikes))
		t.Fail()
	}

	if len(p2.Likes) != 0 {
		t.Log(len(p2.Likes))
		t.Fail()
	}

	if _, ok := p2.Dislikes[uid.String()]; !ok {
		t.Log(uid)
		t.Log(p1.Dislikes)
		t.Fail()
	}

	if _, ok := p2.Likes[uid.String()]; ok {
		t.Log(uid)
		t.Log(p1.Likes)
		t.Fail()
	}
}
