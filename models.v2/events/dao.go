package events

import (
	"time"

	"partisan/location"

	"github.com/jinzhu/gorm"
	"partisan/models.v2/users"
)

func GetByID(id string, guest Subscriber, db *gorm.DB) (e Event, err error) {
	err = db.Where("id = ?", id).Find(&e).Error
	if err != nil {
		return
	}

	e.CountGuests(db)
	e.GetHosts(db)
	e.GetRSVP(guest, db)

	return
}

func SearchForUser(u users.User, db *gorm.DB) (es Events, err error) {
	minX := u.CenterX - 10
	maxX := u.CenterX + 10
	minY := u.CenterY - 10
	maxY := u.CenterY + 10

	minLat, maxLat, minLong, maxLong, err := location.Bounds(u.Latitude, u.Longitude, 25)
	if err != nil {
		return
	}

	err = db.Where("centerx > ? AND centerx < ?", minX, maxX).
		Where("centery > ? AND centery < ?", minY, maxY).
		Where("latitude > ? AND latitude < ?", minLat, maxLat).
		Where("longitude > ? AND longitude < ?", minLong, maxLong).
		Find(&es).Error
	if err != nil {
		return
	}

	if len(es) == 0 {
		return es, ErrEventsNotFound
	}

	es.CollectHosts(db)
	es.CountGuests(db)
	es.CollectRSVPs(u, db)

	return
}

func GetByHost(host Subscriber, db *gorm.DB) (es Events, err error) {
	err = db.Joins("LEFT JOIN event_subscriptions ON event_subscriptions.event_id = events.id").
		Where("event_subscriptions.subscriber_type = ?", host.GetSubscriberType()).
		Where("event_subscriptions.subscriber_id = ?", host.GetID()).
		Where("event_subscriptions.rsvp = ?", RTHost).
		Where("events.start_date > ?::timestamp", time.Now()).
		Find(&es).Error
	if err != nil {
		return
	}

	if len(es) == 0 {
		return es, ErrEventsNotFound
	}

	es.CollectHosts(db)
	es.CountGuests(db)
	es.CollectRSVPs(host, db)

	return
}

func GetByGuest(guest Subscriber, db *gorm.DB) (es Events, err error) {
	err = db.Joins("LEFT JOIN event_subscriptions ON event_subscriptions.event_id = events.id").
		Where("event_subscriptions.subscriber_type = ?", guest.GetSubscriberType()).
		Where("event_subscriptions.subscriber_id = ?", guest.GetID()).
		Where("event_subscriptions.rsvp IN (?)", []RSVPType{RTGoing, RTMaybe}).
		Where("events.start_date > ?::timestamp", time.Now()).
		Find(&es).Error
	if err != nil {
		return
	}

	if len(es) == 0 {
		return es, ErrEventsNotFound
	}

	es.CollectHosts(db)
	es.CountGuests(db)
	es.CollectRSVPs(guest, db)

	return
}

func GetPastByGuest(guest Subscriber, db *gorm.DB) (es Events, err error) {
	err = db.Joins("event_subscriptions ON event_subscriptions.event_id = events.id").
		Where("event_subscriptions.subscriber_type = ?", guest.GetSubscriberType()).
		Where("event_subscriptions.subscriber_id = ?", guest.GetID()).
		Where("event_subscriptions.rsvp IN (?)", []RSVPType{RTGoing, RTMaybe}).
		Where("events.start_date < ?::timestamp", time.Now()).
		Find(&es).Error
	if err != nil {
		return
	}

	es.CollectHosts(db)
	es.CountGuests(db)

	return
}

func (e *Event) GetHosts(db *gorm.DB) error {
	var hosts EventSubscriptions
	if err := db.Where("event_id = ? AND rsvp = ?", e.ID, RTHost).
		Find(&hosts).Error; err != nil {
		return err
	}

	if err := hosts.CollectSubscribers(db); err != nil {
		return err
	}

	for _, host := range hosts {
		e.Hosts = append(e.Hosts, host.Subscriber)
	}

	return nil
}

func (e *Event) CountGuests(db *gorm.DB) error {
	if err := db.Model(EventSubscription{}).
		Where("event_id = ? AND rsvp = ?", e.ID, RTMaybe).
		Count(&e.MaybeCount).Error; err != nil {
		return err
	}

	if err := db.Model(EventSubscription{}).
		Where("event_id = ? AND rsvp IN (?)", e.ID, []RSVPType{RTGoing, RTHost}).
		Count(&e.GoingCount).Error; err != nil {
		return err
	}

	return nil
}

func (e *Event) GetRSVP(guest Subscriber, db *gorm.DB) error {
	var sub EventSubscription
	if err := db.Where("event_id = ?", e.ID).
		Where("subscriber_type = ?", guest.GetSubscriberType()).
		Where("subscriber_id = ?", guest.GetID()).
		Find(&sub).Error; err != nil {
		return err
	}

	e.RSVP = sub.RSVP

	return nil
}

func (es *Events) CollectHosts(db *gorm.DB) error {
	eIDs := es.collectIDs()

	var subs EventSubscriptions
	if err := db.Where("event_id IN (?) AND rsvp = ?", eIDs, RTHost).Find(&subs).Error; err != nil {
		return err
	}

	if err := subs.CollectSubscribers(db); err != nil {
		return err
	}

	events := []Event(*es)
	for i := range events {
		for _, sub := range subs {
			if sub.EventID == events[i].ID {
				events[i].Hosts = append(events[i].Hosts, sub.Subscriber)
			}
		}
	}

	*es = Events(events)

	return nil
}

func (es *Events) CountGuests(db *gorm.DB) error {
	eIDs := es.collectIDs()

	var subs EventSubscriptions
	if err := db.Where("event_id IN (?)", eIDs).Find(&subs).Error; err != nil {
		return err
	}

	events := []Event(*es)
	for i := range events {
		for _, sub := range subs {
			if sub.EventID == events[i].ID {
				if sub.RSVP == RTMaybe {
					events[i].MaybeCount++
				}

				if sub.RSVP == RTGoing || sub.RSVP == RTHost {
					events[i].GoingCount++
				}
			}
		}
	}

	*es = Events(events)

	return nil

}

func (es *Events) CollectRSVPs(guest Subscriber, db *gorm.DB) error {
	eIDs := es.collectIDs()

	var subs EventSubscriptions
	if err := db.Where("event_id IN (?)", eIDs).
		Where("subscriber_type = ?", guest.GetSubscriberType()).
		Where("subscriber_id = ?", guest.GetID()).
		Find(&subs).Error; err != nil {
		return err
	}

	events := []Event(*es)
	for i := range events {
		for _, sub := range subs {
			if sub.EventID == events[i].ID {
				events[i].RSVP = sub.RSVP
			}
		}
	}

	*es = Events(events)

	return nil
}

func (subs *EventSubscriptions) CollectSubscribers(db *gorm.DB) error {
	users, err := users.ListRelated(subs, db)
	if err != nil {
		return err
	}

	subscriptions := []EventSubscription(*subs)
	for i, s := range subscriptions {
		for _, u := range users {
			if s.SubscriberType == STUser && s.SubscriberID == u.ID {
				subscriptions[i].Subscriber = u
			}
		}
	}

	// TODO: Actually implement organizations
	// orgs, err := organizations.ListRelated(subs, db)
	// if err != nil {
	// 	return err
	// }

	// for i, s := range subscriptions {
	// 	for _, o := range orgs {
	// 		if s.SubscriberType == STOrg && s.SubscriberID == org.ID {
	// 			subscriptions[i].Subscriber == o
	// 		}
	// 	}
	// }

	*subs = EventSubscriptions(subscriptions)

	return nil
}

func (es Events) collectIDs() (ids []string) {
	for _, e := range es {
		ids = append(ids, e.ID)
	}

	return
}
