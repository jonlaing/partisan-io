package notifications

import (
	"fmt"
	"math/rand"

	"github.com/jinzhu/gorm"
	"github.com/timehop/apns"
	"partisan/models.v2/users"
)

type PushNotification struct {
	DeviceToken    string
	Message        string
	Action         string
	Meta           map[string]interface{}
	NotificationID string
}

func (n Notification) NewPushNotification(db *gorm.DB) (PushNotification, error) {
	to, err := users.GetByID(n.ToID, db)
	if err != nil {
		return PushNotification{}, err
	}

	from, err := users.GetByID(n.UserID, db)
	if err != nil {
		return PushNotification{}, err
	}

	if len(to.DeviceToken) == 0 {
		return PushNotification{}, ErrDeviceToken
	}

	pn := PushNotification{
		DeviceToken:    to.DeviceToken,
		Message:        n.pnMessage(from.Username),
		Action:         string(n.Action),
		Meta:           map[string]interface{}{"record_id": n.RecordID},
		NotificationID: n.ID,
	}

	if pn.Message == "" {
		return pn, ErrNotifMessage
	}

	return pn, nil
}

func (n PushNotification) Prepare() apns.Notification {
	payload := apns.NewPayload()
	payload.APS.Alert.Body = n.Message
	payload.APS.Sound = "bingbong.aiff"
	payload.SetCustomValue("action", n.Action)
	payload.SetCustomValue("meta", n.Meta)

	pn := apns.NewNotification()
	pn.DeviceToken = n.DeviceToken
	pn.Payload = payload
	pn.Priority = apns.PriorityImmediate
	pn.Identifier = uint32(rand.Intn(100000))
	pn.ID = n.NotificationID

	return pn
}

func (n Notification) pnMessage(username string) string {
	switch n.Action {
	case AFriendRequest:
		return fmt.Sprintf("@%s requested to be your friend.", username)
	case AFriendAccept:
		return fmt.Sprintf("@%s accepted your friend request.", username)
	case AUserTag:
		return fmt.Sprintf("@%s tagged you in something.", username)
	case ALike:
		return fmt.Sprintf("@%s liked something you wrote.", username)
	case AComment:
		return fmt.Sprintf("@%s commented on your post.", username)
	case AEventUpdate:
		return fmt.Sprintf("@%s updated their event.", username)
	}

	return ""
}
