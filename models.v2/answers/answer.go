package answers

import "errors"

// Answer is an answer to a question, included are the coordinates
// of the question, and whether or not the user agreed
// The Map should be in the form of [0,4,8,12] (which would be the entire far-left).
// Check Matcher for more details on the map
type Answer struct {
	Map   []int `json:"map" form:"map"` // defined in matcher.go
	Mask  []int `json:"mask" form:"mask"`
	Agree bool  `json:"agree" form:"agree"`
}

var ErrMap = errors.New("Answer doesn't have map. Probably an error in binding")
