package friendships

import "errors"

var (
	ErrFriendSelf = errors.New("You cannot friend yourself. Being friends with yourself is weird. Go be friends with other people.")
)
