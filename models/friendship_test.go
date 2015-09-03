package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	// "log"
	"os"
	"testing"
)

func TestFindingFriendshipsForUser(t *testing.T) {
	dbURL := os.Getenv("DB_URL")

	dbSession, err := mgo.Dial(dbURL)
	if err != nil {
		panic(err)
	}
	defer dbSession.Close()

	mIO := InitModels(dbSession, "partisan_test")
	db := mIO.DB

	f1Id := bson.NewObjectId()
	f2Id := bson.NewObjectId()

	f1 := Friendship{
		ID:        bson.NewObjectId(),
		UserID:    f1Id,
		FriendID:  f2Id,
		Confirmed: false,
	}
	f1.SetCreateTime()

	f2 := Friendship{
		ID:        bson.NewObjectId(),
		UserID:    f2Id,
		FriendID:  f1Id,
		Confirmed: false,
	}
	f2.SetCreateTime()

	err = db.C("friendships").Insert(f1, f2)
	if err != nil {
		t.Errorf("Couldn't Insert Friendships:", err)
	}

	fs := []Friendship{}
	err = FindFriendshipsByUser(mIO, f1Id, &fs)
	if err != nil {
		t.Errorf("Couldn't Find Friendships:", err)
	}

	if fs[0].UserID != f1Id || fs[0].FriendID != f2Id {
		t.Log(fs[0].UserID)
		t.Log(fs[0].FriendID)
		t.Fail()
	}
}

func TestFindFriendshipsForFriends(t *testing.T) {
	dbURL := os.Getenv("DB_URL")

	dbSession, err := mgo.Dial(dbURL)
	if err != nil {
		panic(err)
	}

	mIO := InitModels(dbSession, "partisan_test")
	db := mIO.DB

	f1Id := bson.NewObjectId()
	f2Id := bson.NewObjectId()

	f1 := Friendship{
		ID:        bson.NewObjectId(),
		UserID:    f1Id,
		FriendID:  f2Id,
		Confirmed: false,
	}
	f1.SetCreateTime()

	f2 := Friendship{
		ID:        bson.NewObjectId(),
		UserID:    f2Id,
		FriendID:  f1Id,
		Confirmed: false,
	}
	f2.SetCreateTime()

	err = db.C("friendships").Insert(f1, f2)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	fs := [2]Friendship{}
	err = FindFriendship(mIO, f1Id, f2Id, &fs)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if fs[0].UserID != f1Id || fs[0].FriendID != f2Id {
		t.Log(fs[0].UserID)
		t.Log(fs[0].FriendID)
		t.Fail()
	}

	if fs[1].UserID != f2Id || fs[1].FriendID != f1Id {
		t.Log(fs[1].UserID)
		t.Log(fs[1].FriendID)
		t.Fail()
	}
}

func TestFindingFriendIds(t *testing.T) {
	dbURL := os.Getenv("DB_URL")

	dbSession, err := mgo.Dial(dbURL)
	if err != nil {
		panic(err)
	}

	mIO := InitModels(dbSession, "partisan_test")
	db := mIO.DB

	f1Id := bson.NewObjectId()
	f2Id := bson.NewObjectId()

	f1 := Friendship{
		ID:        bson.NewObjectId(),
		UserID:    f1Id,
		FriendID:  f2Id,
		Confirmed: false,
	}
	f1.SetCreateTime()

	f2 := Friendship{
		ID:        bson.NewObjectId(),
		UserID:    f2Id,
		FriendID:  f1Id,
		Confirmed: false,
	}
	f2.SetCreateTime()

	err = db.C("friendships").Insert(f1, f2)

	ids, err := FindFriendIds(mIO, f1Id)
	if err != nil {
		t.Errorf("Couldn't find any friendships", err)
	}

	if len(ids) < 1 {
		t.Log(len(ids))
		panic("not enough ids")
	}

	if ids[0] != f2Id {
		t.Log(ids[0])
		t.Log(f2Id)
		t.Fail()
	}
}

func TestInsertingFriendship(t *testing.T) {
	dbURL := os.Getenv("DB_URL")

	dbSession, err := mgo.Dial(dbURL)
	if err != nil {
		panic(err)
	}
	defer dbSession.Close()

	mIO := InitModels(dbSession, "partisan_test")
	db := mIO.DB

	f1Id := bson.NewObjectId()
	f2Id := bson.NewObjectId()

	err = InsertFriendship(mIO, f1Id, f2Id)
	if err != nil {
		t.Errorf("Couldn't insert friendships:", err)
	}

	f1 := Friendship{}
	f2 := Friendship{}
	err = db.C("friendships").Find(bson.M{"user_id": f1Id, "friend_id": f2Id}).One(&f1)
	if err != nil {
		t.Error("Couldn't find friendships:", err)
	}
	err = db.C("friendships").Find(bson.M{"user_id": f2Id, "friend_id": f1Id}).One(&f2)
	if err != nil {
		t.Error("Couldn't find friendships:", err)
	}
}

func TestConfirmingFriendship(t *testing.T) {
	dbURL := os.Getenv("DB_URL")

	dbSession, err := mgo.Dial(dbURL)
	if err != nil {
		panic(err)
	}
	defer dbSession.Close()

	mIO := InitModels(dbSession, "partisan_test")
	db := mIO.DB

	f1Id := bson.NewObjectId()
	f2Id := bson.NewObjectId()

	f1 := Friendship{
		ID:        bson.NewObjectId(),
		UserID:    f1Id,
		FriendID:  f2Id,
		Confirmed: false,
	}
	f1.SetCreateTime()

	f2 := Friendship{
		ID:        bson.NewObjectId(),
		UserID:    f2Id,
		FriendID:  f1Id,
		Confirmed: false,
	}
	f2.SetCreateTime()

	err = db.C("friendships").Insert(f1, f2)
	if err != nil {
		t.Errorf("Couldn't Insert Friendships:", err)
	}

	err = ConfirmFriendship(mIO, f1Id, f2Id)
	if err != nil {
		t.Error(err)
	}

	f1a := Friendship{}
	f2a := Friendship{}
	err = db.C("friendships").Find(bson.M{"user_id": f1Id, "friend_id": f2Id}).One(&f1a)
	if err != nil {
		t.Error("Couldn't find friendships:", err)
	}
	err = db.C("friendships").Find(bson.M{"user_id": f2Id, "friend_id": f1Id}).One(&f2a)
	if err != nil {
		t.Error("Couldn't find friendships:", err)
	}

	if !f1a.Confirmed {
		t.Fail()
	}

	if !f2a.Confirmed {
		t.Fail()
	}
}
