package notifications

import (
	"github.com/jinzhu/gorm"
	"partisan/models.v2/users"
)

func ListByUserID(userID string, db *gorm.DB) (ns Notifications, err error) {
	err = db.Where("to_id = ?", userID).Limit(25).Find(&ns).Error
	if err != nil {
		return
	}

	ns.CollectUsers(db)
	return
}

func (ns *Notifications) CollectUsers(db *gorm.DB) {
	us, err := users.ListByIDs(collectUserIDs(*ns), db)
	if err != nil {
		return
	}

	notifs := []Notification(*ns)

	for i := range notifs {
		for _, u := range us {
			if u.ID == notifs[i] {
				notifs[i].User = u
			}
		}
	}

	*ns = Notifications(notifs)
}

func (ns Notifications) MarkRead(db *gorm.DB) error {
	var ids []string

	for _, n := range ns {
		if !n.Read {
			ids = append(ids, n.ID)
		}
	}

	if len(ids) == 0 {
		return nil
	}

	return db.Table("notifications").Where("id IN (?)", ids).UpdateColumn("read", true).Error
}

func collectUserIDs(ns Notifications) (userIDs []string) {
	for _, n := range ns {
		userIDs = append(userIDs, n.UserID)
	}

	return
}
