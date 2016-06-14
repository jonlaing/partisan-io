package notifications

import "database/sql/driver"

type Action string

const (
	AFriendRequest Action = "friendrequest"
	AFriendAccept         = "friendaccept"
	AMention              = "mention"
	ALike                 = "like"
	AComment              = "comment"
	AEventUpdate          = "eventupdate"
)

func (a *Action) Scan(src interface{}) error {
	astring, ok := src.([]byte)
	if !ok {
		return ErrScanAction
	}

	*a = Action(astring)
	return nil
}

func (a Action) Value() (driver.Value, error) {
	return string(a), nil
}
