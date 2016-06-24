package notifications

import "errors"

var (
	ErrScanAction   = errors.New("Could not scan action")
	ErrNotifySelf   = errors.New("Cannot notify yourself")
	ErrDeviceToken  = errors.New("This user does not have a valid device token")
	ErrNotifMessage = errors.New("This notification does not have a valid message")
)
