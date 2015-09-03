package models

import (
	"github.com/kennygrant/sanitize"
	omniauth "github.com/stretchr/gomniauth/common"
	"gopkg.in/mgo.v2/bson"
	"net/url"
	"partisan-api/match"
	// v "partisan-api/models/validators"
)

// User is a record of an app user in the DB. We don't store passwords because we use OAuth
type User struct {
	ID        bson.ObjectId  `json:"id" bson:"_id,omitempty"`
	Username  string         `json:"username" bson:"username"`
	FullName  string         `json:"full_name" bson:"full_name"`
	Email     string         `json:"email" bson:"email"`
	AvatarURL string         `json:"avatar_url" bson:"avatar_url"`
	StarCount int            `json:"star_count,omitempty" bson:"star_count,omitempty"` // how many stars the user has to give
	Answers   []match.Answer `json:"-" bson:"answers"`
	Xcoord    int            `json:"xcoord" bson:"xcoord"`
	Ycoord    int            `json:"ycoord" bson:"ycoord"`
	Zcoord    int            `json:"zcoord" bson:"zcoord"`
	datedRecord
}

// NewUserFromAuth creates a new user based on a successful OAuth authenitcation
func NewUserFromAuth(ou omniauth.User) User {
	oauthIds := make(map[string]string)
	oauthIds["github"] = ou.IDForProvider("github")

	return User{
		ID:        bson.NewObjectId(),
		Username:  ou.Nickname(),
		Email:     ou.Email(),
		AvatarURL: ou.AvatarURL(),
	}
}

// GetID satisfies Model interface
func (u *User) GetID() bson.ObjectId {
	return u.ID
}

// Scan implements Model interface
func (u *User) Scan(params url.Values) error {
	username, err := sanitize.HTMLAllowing(params.Get("username"))
	if err != nil {
		return err
	}

	fullname, err := sanitize.HTMLAllowing(params.Get("full_name"))
	if err != nil {
		return err
	}

	email, err := sanitize.HTMLAllowing(params.Get("email"))
	if err != nil {
		return err
	}

	avatarURL, err := sanitize.HTMLAllowing(params.Get("avatar_url"))
	if err != nil {
		return err
	}

	u.Username = username
	u.FullName = fullname
	u.Email = email
	u.AvatarURL = avatarURL
	u.SetUpdateTime()

	return nil
}

// Validate implements Model interface
func (u *User) Validate() map[string]error {
	return map[string]error{}
}

// FindUserByEmail finds a user based on their email address
func FindUserByEmail(email string, u *User) error {
	return ModelIO.Collection(u).Find(bson.M{"email": email}).One(u)
}

// FindFriends finds all a users' friends
func (u User) FindFriends(db Database) ([]User, error) {
	fids, _ := FindFriendIds(db, u.ID)
	fs := []User{}
	err := db.Collection(u).Find(bson.M{"_id": bson.M{"$in": fids}}).All(&fs)
	return fs, err
}
