package events

import (
	"time"

	"partisan/matcher"

	models "partisan/models.v2"
)

type EventSubscription struct {
	ID             string    `json:"id" gorm:"primary_key" sql:"type:uuid;default:uuid_generate_v4()"`
	SubscriberType string    `json:"-"`
	SubscriberID   string    `json:"-" sql:"type:uuid"`
	EventID        string    `json:"event_id" sql:"type:uuid"`
	RSVP           RSVPType  `json:"rsvp"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	Subscriber Subscriber `json:"subscriber" sql:"-"`
}

type EventSubscriptions []EventSubscription

func (subs EventSubscriptions) GetUserIDs() (ids []string) {
	for _, sub := range subs {
		if sub.SubscriberType == STUser {
			ids = append(ids, sub.SubscriberID)
		}
	}

	return
}

func (subs EventSubscriptions) GetOrgIDs() (ids []string) {
	for _, sub := range subs {
		if sub.SubscriberType == STOrg {
			ids = append(ids, sub.SubscriberID)
		}
	}

	return
}

func (e *Event) NewHost(host Subscriber) (s EventSubscription, errs models.ValidationErrors) {
	errs = make(models.ValidationErrors)

	s.SubscriberType = host.GetSubscriberType()
	s.SubscriberID = host.GetID()
	s.EventID = e.ID
	s.RSVP = RTHost
	s.CreatedAt = time.Now()
	s.UpdatedAt = s.CreatedAt

	e.Hosts = append(e.Hosts, host)

	if len(e.Hosts) > 1 {
		var maps []matcher.PoliticalMap
		for _, host := range e.Hosts {
			maps = append(maps, host.GetPoliticalMap())
		}

		pMap, err := matcher.Merge(maps...)
		if err != nil {
			errs["political_map"] = err
		}

		e.PoliticalMap = pMap
	} else {
		e.PoliticalMap = host.GetPoliticalMap()
	}

	x, y := e.PoliticalMap.Center()
	e.CenterX = x
	e.CenterY = y
	e.UpdatedAt = time.Now()

	return
}

func (e *Event) NewGuest(guest Subscriber, rsvp RSVPType) (s EventSubscription, errs models.ValidationErrors) {
	s.SubscriberType = guest.GetSubscriberType()
	s.SubscriberID = guest.GetID()
	s.EventID = e.ID
	s.RSVP = rsvp
	s.CreatedAt = time.Now()
	s.UpdatedAt = s.CreatedAt

	if rsvp == RTGoing {
		e.GoingCount++
	}

	if rsvp == RTMaybe {
		e.MaybeCount++
	}

	return
}
