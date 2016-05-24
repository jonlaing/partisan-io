package matches

import (
	"math"

	"partisan/models.v2/users"
)

type Match struct {
	User  users.User
	Match float64
}

type Matches []Match

// Len satisfies sort.Interface
func (ms Matches) Len() int {
	return len(ms)
}

// Less satisfies sort.Interface
func (ms Matches) Less(a, b int) bool {
	return ms[a].Match < ms[b].Match
}

// Swap satisfies sort.Interface
func (ms Matches) Swap(a, b int) {
	ms[b], ms[a] = ms[a], ms[b]
}

type SearchBinding struct {
	Gender     string  `json:"gender"`
	MinAge     int     `json:"min_age"`
	MaxAge     int     `json:"max_age"`
	Radius     float64 `json:"radius"` // in miles
	LookingFor int     `json:"looking_for"`
	Page       int     `json:"page"`
}

// Degrees converts miles into coordinate degrees
func (s SearchBinding) Degrees() float64 {
	earthRadius := float64(3959) // in miles
	return float64(s.Radius) / earthRadius * float64(180) / math.Pi
}
