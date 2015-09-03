package models

import (
	"github.com/kennygrant/sanitize"
	"gopkg.in/mgo.v2/bson"
	"net/url"
	c "partisan-api/models/concerns"
	v "partisan-api/models/validators"
)

// Comment is a polymorphic user generated comment.
type Comment struct {
	ID       bson.ObjectId `json:"id" bson:"_id,omitempty"`              // Primary key
	UserID   bson.ObjectId `json:"user_id" bson:"user_id,omitempty"`     // ID of user
	RecordID bson.ObjectId `json:"record_id" bson:"record_id,omitempty"` // The ID of the record for polymorphism
	Body     string        `json:"body" bson:"body"`                     // The body of the comment
	Likes    c.Likes       `json:"likes" bson:"likes"`                   // Likes for comment
	voteableRecord
}

// GetID satisfies the Model interface
func (c *Comment) GetID() bson.ObjectId {
	return c.ID
}

// Scan implements Model interface
func (c *Comment) Scan(params url.Values) error {
	body, err := sanitize.HTMLAllowing(params.Get("body"))
	if err != nil {
		return err
	}

	c.Body = body
	c.SetUpdateTime()
	return nil
}

// Validate implements Model interface
func (c *Comment) Validate() (verrs map[string]error) {
	if err := v.NonBlank(c.Body, "body"); err != nil {
		verrs["body"] = err
	}

	if err := v.NonEpoch(c.Created, "created"); err != nil {
		verrs["created"] = err
	}

	if err := v.NonEpoch(c.Updated, "updated"); err != nil {
		verrs["updated"] = err
	}

	return
}
