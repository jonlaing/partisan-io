package posts

import "database/sql/driver"

type Action string

const (
	APost    = "post"
	AComment = "comment"
	ALike    = "like"
)

type ParentType string

const (
	// PTNoType is a post to a feed
	PTNoType ParentType = ""
	// PTPost is a comment or like to a post
	PTPost = "post"
	// PTComment is a like of a comment
	PTComment = "comment"
	// PTEvent is a post to an event
	PTEvent = "event"
	// PTGroup is a post to a group
	PTGroup = "group"
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

func (p *ParentType) Scan(src interface{}) error {
	pstring, ok := src.([]byte)
	if !ok {
		*p = ParentType("")
		return nil
	}

	*p = ParentType(pstring)
	return nil
}

func (p ParentType) Value() (driver.Value, error) {
	return string(p), nil
}

func validAction(s Action) bool {
	switch s {
	case APost:
		fallthrough
	case AComment:
		fallthrough
	case ALike:
		return true
	default:
		return false
	}
}

func validParentType(s ParentType) bool {
	switch s {
	case PTNoType:
		fallthrough
	case PTPost:
		fallthrough
	case PTComment:
		fallthrough
	case PTEvent:
		fallthrough
	case PTGroup:
		return true
	default:
		return false
	}
}

func validPostParentType(s ParentType) bool {
	switch s {
	case PTEvent:
		fallthrough
	case PTGroup:
		return true
	default:
		return false
	}
}
