package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Flag struct {
	ID         bson.ObjectId `json:"id" bson:"_id,omitempty"`
	UserID     bson.ObjectId `json:"user_id" bson:"user_id"`
	RecordID   bson.ObjectId `json:"record_id" bson:"record_id"`
	RecordType string        `json:"record_type" bson:"record_type"`
}
