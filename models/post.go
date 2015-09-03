package models

import (
	"fmt"
	"github.com/kennygrant/sanitize"
	"gopkg.in/mgo.v2/bson"
	"net/url"
	c "partisan-api/models/concerns"
	v "partisan-api/models/validators"
)

// Post is the primary user created content. It can be just about anything.
type Post struct {
	ID       bson.ObjectId `json:"id" bson:"_id,omitempty"`          // Primary key
	UserID   bson.ObjectId `json:"user_id" bson:"user_id,omitempty"` // ID of user that created Post
	Body     string        `json:"body" bson:"body"`                 // The text based body
	Likes    c.Likes       `json:"likes" bson:"likes"`               // The likes from other users
	Dislikes c.Dislikes    `json:"dislikes" bson:"dislikes"`         // The dislikes from other users
	Stars    c.Stars       `json:"stars" bson:"stars"`               // Stars are a special like type
	datedRecord
}

// Validate satisfies the Model interface
func (p *Post) Validate() (verrs map[string]error) {
	if err := v.NonBlank(p.Body, "body"); err != nil {
		verrs["body"] = err
	}

	if err := v.NonEpoch(p.Created, "created"); err != nil {
		verrs["created"] = err
	}

	if err := v.NonEpoch(p.Updated, "updated"); err != nil {
		verrs["updated"] = err
	}

	return
}

// Validate satisfies the Model interface
func (p *Posts) Validate() map[string]error {
	return map[string]error{"all": fmt.Errorf("Can't batch insert Posts")}
}

// FindPostsByUser finds all the posts a user has made
func FindPostsByUser(userID bson.ObjectId, ps *[]Post) error {
	return ModelIO.DB.C("posts").Find(bson.M{"user_id": userID}).All(ps)
}

// Like increments Post.Likes and decrements Post.Dislikes if
// the same user has already disliked
func (p *Post) Like(userID bson.ObjectId) error {
	id := userID.String()

	if _, ok := p.Likes[id]; !ok {
		p.Likes[id] = true
	}

	if _, ok := p.Dislikes[id]; ok {
		delete(p.Dislikes, id)
	}

	return ModelIO.Update(p)
}

// Dislike increments Post.Dislikes and decrements Post.Likes if
// the same user had already liked.
func (p *Post) Dislike(userID bson.ObjectId) error {
	id := userID.String()

	if _, ok := p.Dislikes[id]; !ok {
		p.Dislikes[id] = true
	}

	if _, ok := p.Likes[id]; ok {
		delete(p.Likes, id)
	}

	return ModelIO.Update(p)
}
