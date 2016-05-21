package notifications

import "errors"

var (
	ErrScanAction = errors.New("Could not scan action")
	ErrNotifySelf = errors.New("Cannot notify yourself")
)
