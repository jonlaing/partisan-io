package models

import (
	"gopkg.in/mgo.v2/bson"
)

// Trying to keep track of user personalities to gameify the app
const (
	Controversial = iota // Posts things with likes and dislikes being close
	Prolific             // Posts a lot
	Participant          // Comments a lot
	Debater              // Creates a lot of discussion topics
	Friendly             // Someone who has a lot of friends
)

// Profile is a user profile
type Profile struct {
	ID       bson.ObjectId `json:"id" bson:"_id,omitempty"`
	UserID   bson.ObjectId `json:"user_id" bson:"user_id"`
	About    string        `json:"about" bson:"about"`
	Badges   uint          `json:"badges" bson:"badges"`
	Location string        `json:"location,omitempty" bson:"location,omitempty"`
	datedRecord
}
