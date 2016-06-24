package notifications

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/timehop/apns"
	"partisan/models.v2/users"
)

type PushNotification struct {
	DeviceToken string
	Message     string
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
		DeviceToken: to.DeviceToken,
		Message:     n.pnMessage(from.Username),
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

	pn := apns.NewNotification()
	pn.DeviceToken = n.DeviceToken
	pn.Payload = payload
	pn.Priority = apns.PriorityImmediate

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
