package models

import (
	"gopkg.in/mgo.v2/bson"
)

// Friendship is a joining table between two users who are friends. For each
// friendship there are two Friendship records: one for each user
type Friendship struct {
	ID        bson.ObjectId `json:"id" bson:"_id,omitempty"`              // Primary key
	UserID    bson.ObjectId `json:"user_id" bson:"user_id,omitempty"`     // ID of first user
	FriendID  bson.ObjectId `json:"friend_id" bson:"friend_id,omitempty"` // ID of second user
	Confirmed bool          `json:"confirmed" bson:"confirmed"`           // Whether the second user has confirmed
	datedRecord
}

// FindFriendshipsByUser finds all friendships for a given user
func FindFriendshipsByUser(db Database, userID bson.ObjectId, fs *[]Friendship) error {
	err := db.Collection(fs).Find(bson.M{"user_id": userID}).All(fs)
	return err
}

// FindFriendship finds the two Friendships for a griven friendship
func FindFriendship(db Database, userID, friendID bson.ObjectId, fs *[2]Friendship) error {
	f1 := Friendship{}
	f2 := Friendship{}
	err := db.Collection(fs[0]).Find(bson.M{"user_id": userID, "friend_id": friendID}).One(&f1)
	if err != nil {
		return err
	}
	err = db.Collection(fs[0]).Find(bson.M{"user_id": friendID, "friend_id": userID}).One(&f2)
	if err != nil {
		return err
	}

	fs[0] = f1
	fs[1] = f2

	return nil
}

// FindFriendIds finds all the ids of a user's friends
func FindFriendIds(db Database, userID bson.ObjectId) (ids []bson.ObjectId, err error) {
	fs := []Friendship{}
	err = FindFriendshipsByUser(db, userID, &fs)
	for _, v := range fs {
		ids = append(ids, v.FriendID)
	}

	return ids, err
}

// InsertFriendship creates a new friendships by inserting two Friendship records
func InsertFriendship(db Database, userID, friendID bson.ObjectId) error {
	f1 := Friendship{
		ID:        bson.NewObjectId(),
		UserID:    userID,
		FriendID:  friendID,
		Confirmed: false,
	}
	f2 := Friendship{
		ID:        bson.NewObjectId(),
		UserID:    friendID,
		FriendID:  userID,
		Confirmed: false,
	}

	return db.Collection(f1).Insert(f1, f2)
}

// ConfirmFriendship confirms the friendship for two users
func ConfirmFriendship(db Database, userID, friendID bson.ObjectId) error {
	fs := [2]Friendship{}
	FindFriendship(db, userID, friendID, &fs)
	fs[0].Confirmed = true
	fs[1].Confirmed = true

	err := db.Collection(fs[0]).UpdateId(fs[0].ID, fs[0])
	if err != nil {
		return err
	}

	return db.Collection(fs[0]).UpdateId(fs[1].ID, fs[1])
}

// RemoveFriendship unfriends two users
func RemoveFriendship(db Database, userID, friendID bson.ObjectId) error {
	fs := [2]Friendship{}
	FindFriendship(db, userID, friendID, &fs)

	err := db.Collection(fs[0]).RemoveId(fs[0].ID)
	if err != nil {
		return nil
	}

	return db.Collection(fs[0]).RemoveId(fs[1].ID)
}
