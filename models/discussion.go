package models

import (
	"github.com/kennygrant/sanitize"
	"gopkg.in/mgo.v2/bson"
	"net/url"
	// v "partisan-api/models/validators"
)

// Discussion is a sort of minimum viable product content type for some more
// advanced content types. Not sure if this is going to stay. Basically it's
// just meant to inspire discussion (aka arguments) between a wider group than Post
type Discussion struct {
	ID      bson.ObjectId `json:"id" bson:"_id,omitempty"`          // Primary key
	UserID  bson.ObjectId `json:"user_id" bson:"user_id,omitempty"` // ID of user that created Discussion
	Subject string        `json:"subject" bson:"subject"`           // Prompt for discussion
	Body    string        `json:"body" bson:"body"`                 // Further detail for discussion
	voteableRecord
}

// GetID satisfies Model interface
func (d *Discussion) GetID() bson.ObjectId {
	return d.ID
}

// Scan satisfies Model interface
func (d *Discussion) Scan(params url.Values) error {
	subject, err := sanitize.HTMLAllowing(params.Get("sanitize"))
	if err != nil {
		return err
	}

	body, err := sanitize.HTMLAllowing(params.Get("body"))
	if err != nil {
		return err
	}

	d.Subject = subject
	d.Body = body
	d.SetUpdateTime()

	return nil
}

// Validate implements Model interface
func (d Discussion) Validate() map[string]error {
	return map[string]error{}
}
