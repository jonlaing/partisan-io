package models

import (
	"gopkg.in/mgo.v2/bson"
)

// FeedAction indicates the type of action that
// created the FeedItem
type FeedAction uint

const (
	// FeedCreate indicates a create action
	FeedCreate FeedAction = iota

	// FeedLike indicates a like action
	FeedLike

	// FeedDislike indicates a dislike action
	FeedDislike

	// FeedStar indicates a star action
	FeedStar

	// FeedComment indicates a comment action
	FeedComment

	// FeedEdit indicates a edit action. Not sure if
	// I'm going to use this.
	FeedEdit
)

// FeedItem is a record that shows up in a user's feed and is created by their
// friends. Think of it like a Twitter or Facebook feed item. FeedItem is polymorphic.
type FeedItem struct {
	ID       bson.ObjectId `json:"id" bson:"_id,omitempty"`              // Primary key
	UserID   bson.ObjectId `json:"user_id" bson:"user_id,omitempty"`     // ID of user that initiated the action
	RecordID bson.ObjectId `json:"record_id" bson:"record_id,omitempty"` // ID of record for polymorphism
	Action   FeedAction    `json:"action" bson:"action"`                 // The action that created the FeedItem
	datedRecord
}

// // GetID implements Model interface
// func (f *FeedItem) GetID() bson.ObjectId {
// 	return f.ID
// }

// // SetAttribute implements Model interface
// // This is not how FeedItems are intended to be used
// //
// func (f *FeedItem) SetAttribute() error {
// 	return nil
// }

// // Validate implements Model interface
// func (f *FeedItem) Validate() []*v.ValidationError {
// 	return []*v.ValidationError{nil}
// }

// NewFeedItem creates a new FeedItem from Model
func NewFeedItem(db Database, item Model, userID bson.ObjectId, action FeedAction) error {
	f := FeedItem{
		ID:       bson.NewObjectId(),
		UserID:   userID,
		RecordID: item.GetID(),
		Action:   action,
	}
	f.SetCreateTime()

	return db.Collection(f).Insert(f)
}

// FindFeedItems finds all the feed items that will go into a user's feed. For a given
// user, it finds friends and the FeedItems they've created.
// TODO: Need to put time constraint
func FindFeedItems(db Database, userID bson.ObjectId, fs []FeedItem) error {
	fids, _ := FindFriendIds(db, userID)
	err := db.Collection(fs).Find(bson.M{"user_id": bson.M{"$in": fids}}).All(&fs)
	return err
}
