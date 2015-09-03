package models

import (
	"errors"
	"time"
)

type datedRecord struct {
	Created time.Time `json:"created" bson:"created,inline"`
	Updated time.Time `json:"updated" bson:"updated,inline"`
}

func (d *datedRecord) SetCreateTime() error {
	if d.Created.IsZero() {
		d.Created = time.Now()
		d.Updated = time.Now()
		return nil
	}

	return errors.New("Created time is not zero")
}

func (d *datedRecord) SetUpdateTime() {
	if d.Created.IsZero() {
		d.SetCreateTime()
		return
	}

	d.Updated = time.Now()
}
