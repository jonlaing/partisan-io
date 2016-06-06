package events

import (
	"database/sql/driver"
	"partisan/matcher"
)

type Subscriber interface {
	GetID() string
	GetSubscriberType() string
	GetPoliticalMap() matcher.PoliticalMap
}

type RSVPType string

const (
	RTNone  RSVPType = ""
	RTGoing          = "going"
	RTMaybe          = "maybe"
	RTHost           = "host"
)

type SubscriberType string

const (
	STNone SubscriberType = ""
	STUser                = "user"
	STOrg                 = "organization"
)

func (r *RSVPType) Scan(src interface{}) error {
	rstring, ok := src.([]byte)
	if !ok {
		return ErrScanRSVPType
	}

	*r = RSVPType(rstring)
	return nil
}

func (s RSVPType) Value() (driver.Value, error) {
	return string(s), nil
}

func (s *SubscriberType) Scan(src interface{}) error {
	sstring, ok := src.([]byte)
	if !ok {
		return ErrScanSubscriberType
	}

	*s = SubscriberType(sstring)
	return nil
}

func (s SubscriberType) Value() (driver.Value, error) {
	return string(s), nil
}
