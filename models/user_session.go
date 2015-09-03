package models

import (
	"gopkg.in/mgo.v2/bson"
	"net/url"
	"time"
)

// UserSession is the session of an authenticated user. It is meant to expire in 48 hours.
type UserSession struct {
	ID        bson.ObjectId `json:"id" bson:"_id,omitempty"` // ID will be used as X-Access-Token for authenticated requests
	UserID    bson.ObjectId `json:"user_id" bson:"user_id"`
	LastLogin time.Time     `json:"last_login" bson:"last_login"`
	LastSeen  time.Time     `json:"last_seen" bson:"last_seen"`
}

// GetID satisfies Model interface
func (us *UserSession) GetID() bson.ObjectId {
	return us.ID
}

// SetAttributes implements Model interface
// This isn't really intended to be used, but UserSession needs to implement this interface
func (us *UserSession) SetAttributes(params url.Values) error {
	return nil
}

// Validate implements Model interface
func (us *UserSession) Validate() map[string]error {
	return map[string]error{}
}

// IsExpired returns true if the session is older than 48 hours without activity
func (us UserSession) IsExpired() bool {
	return time.Now().Sub(us.LastSeen).Hours() > float64(48)
}
