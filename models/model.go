package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	// "log"
	"net/url"
	"reflect"
	"regexp"
	"strings"
)

// ModelDatabase is used to communicate with the DB
type ModelDatabase struct {
	DB *mgo.Database
}

// Model is the interace to DRY up all Model CRUD actions
type Model interface {
	GetID() bson.ObjectId
	Scan(params url.Values) error
	Validate() map[string]error
}

// Find finds a Model record in the DB by ID
func (m ModelDatabase) Find(r Model) error {
	return m.Collection(r).FindId(r.GetID()).One(r)
}

// Where finds a Model record in the DB by Query
func (m ModelDatabase) Where(query interface{}, r Model) error {
	return m.Collection(r).Find(query).One(r)
}

// Insert inserts a Model record into the DB
func (m ModelDatabase) Insert(r Model) error {
	_, err := m.Collection(r).UpsertId(r.GetID(), r)
	return err
}

// Update updates a Model record in the DB
func (m ModelDatabase) Update(r Model) error {
	return m.Collection(r).UpdateId(r.GetID(), r)
}

// Destroy destroys Model record in the DB
func (m ModelDatabase) Destroy(r Model) error {
	return m.Collection(r).RemoveId(r.GetID())
}

func valueName(v reflect.Value) string {
	if v.Kind() == reflect.Ptr {
		if e := v.Elem(); e.Kind() == reflect.Slice {
			return valueName(e)
		}
		return v.Elem().Type().Name()
	}

	if v.Kind() == reflect.Slice {
		return v.Type().Elem().Name()
	}

	return v.Type().Name()
}

// Get the name of the model through reflection (i.e. Post returns "post")
func modelName(r interface{}) string {
	name := valueName(reflect.ValueOf(r))

	rx := regexp.MustCompile("([a-z])([A-Z])")
	name = rx.ReplaceAllString(name, "${1}_${2}")
	return strings.ToLower(name)
}

// Pluralizes modelName()
func collectionName(r interface{}) string {
	name := modelName(r)
	return strings.Join([]string{name, "s"}, "")
}

// Collection is to DRY and standardize mongo collections
func (m ModelDatabase) Collection(r interface{}) *mgo.Collection {
	return m.DB.C(collectionName(r))
}
