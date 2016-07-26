package matches

import (
	"math"
	"strings"

	"partisan/models.v2/users"
)

type Match struct {
	User     users.User `json:"user"`
	Match    float64    `json:"match"`
	Distance float64    `json:"-"`
}

type Matches []Match

// Len satisfies sort.Interface
func (ms Matches) Len() int {
	return len(ms)
}

// Less satisfies sort.Interface
func (ms Matches) Less(a, b int) bool {
	if ms[a].Distance-ms[b].Distance < 10 {
		return ms[a].Match < ms[b].Match
	}

	return ms[a].Distance < ms[b].Distance
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

func (s SearchBinding) GenderGroup() ([]string, error) {
	genderGuesses := [][]string{
		// cisfemme
		{
			"female",
			"woman",
			"girl",
			"lady",
			"ciswoman", "ciswomyn",
			"cis woman", "cis womyn",
		},
		// cismasc
		{
			"male",
			"man",
			"boy",
			"guy",
			"cisman",
			"cis man",
		},
		// transfemme
		{
			"femme", "feminine",
			"transwoman", "transwomyn",
			"trans woman", "trans womyn",
			"mtf", "m -> f",
		},
		// transmasc
		{
			"masc", "masculine",
			"transman",
			"trans man",
			"ftm", "f -> m",
		},
	}

	for _, gs := range genderGuesses {
		for _, g := range gs {
			if strings.ToLower(s.Gender) == g {
				return gs, nil
			}
		}
	}

	return []string{}, ErrGenderGroup
}
