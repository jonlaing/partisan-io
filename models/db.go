package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// ModelIO communicates with the database
var ModelIO ModelDatabase
var dbSession *mgo.Session

// Database is an interface to abstract writing and reading from DB.
// ModelDatabase implements Database. This is helpful for writing mocks
// and stubs for testing, and for abstracting handlers in the main app.
type Database interface {
	Collection(m interface{}) *mgo.Collection
	Find(id bson.ObjectId, m Model) error
	Insert(m Model) error
	Update(m Model) error
	Destroy(m Model) error
}

// InitModels sets up the database for the main app
func InitModels(appSession *mgo.Session, dbName string) ModelDatabase {
	dbSession = appSession.Copy()
	ModelIO.DB = dbSession.DB(dbName)
	return ModelIO
}
