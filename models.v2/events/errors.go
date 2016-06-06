package events

import "errors"

var (
	ErrScanRSVPType       = errors.New("Error scanning RSVP Type")
	ErrScanSubscriberType = errors.New("Error scanning Subscriber Type")
	ErrStartDate          = errors.New("Invalid start date")
	ErrEndDate            = errors.New("Invalid end date")
	ErrEventsNotFound     = errors.New("Couldn't find any events for this search")
)
