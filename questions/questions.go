package questions

import (
	"fmt"
	"math/rand"
	"partisan/matcher"
)

// The radius for an acceptable question
const radius = 17

var (
	MProState    = []int{0, 1, 2, 3, 4, 5, 6, 7}
	MProCapital  = []int{2, 3, 6, 7, 10, 11, 14, 15}
	MAntiState   = []int{8, 9, 10, 11, 12, 13, 14, 15}
	MAntiCapital = []int{0, 1, 4, 5, 8, 9, 12, 13}

	MMiddleLeft              = []int{1, 5, 9, 13}
	MMiddleRight             = []int{2, 6, 10, 14}
	MMiddleAuthoritarian     = []int{4, 5, 6, 7}
	MMiddleAntiAuthoritarian = []int{8, 9, 10, 11}

	MAuthoritarian     = []int{0, 1, 2, 3}
	MFarRight          = []int{3, 7, 11, 15}
	MAntiAuthoritarian = []int{12, 13, 14, 15}
	MFarLeft           = []int{0, 4, 8, 12}

	MAuthSocialist  = []int{0, 1, 4, 5}
	MAuthCapitalist = []int{2, 3, 6, 7}
	MLibSocialist   = []int{8, 9, 12, 13}
	MLibCapitalist  = []int{10, 11, 14, 15}
)

// Question holds questions to guage a user's political leanings
type Question struct {
	Prompt string `json:"prompt"`
	Map    []int  `json:"map"` // defined in matcher.go
}

// Questions are in groups of four to contain pro-state, anti-state, pro-capital, and anti-capital sentiments
type Questions [4]Question

// QuestionSet is a collection of questions. Using ValidSet, the algorithm can determine whether this is an appropriate question given
// a User's current Center.
type QuestionSet struct {
	// In case you forget, as you have 100 times, yes, Mask is necessary! If, for instance, you disagree with something mapped to
	// MMiddleLeft, without the Mask, both the Far Left AND the whole Right Wing will get points added! With the Mask, only the
	// Far Left will get highlighted.
	Mask      []int     `json:"mask"`
	Questions Questions `json:"questions"`
}

func (qs QuestionSet) ValidSet(x, y int) bool {
	var pMap matcher.PoliticalMap
	for _, q := range qs.Questions {
		pMap.Add(q.Map, qs.Mask, true)
	}

	centerX, centerY := pMap.Center()

	fmt.Println("CENTER:", centerX, centerY)
	fmt.Println("RADIUS:", centerX+radius, centerY+radius, centerX-radius, centerY-radius)

	return x <= centerX+radius &&
		x >= centerX-radius ||
		y <= centerY+radius &&
			y >= centerY-radius
}

func (qs QuestionSet) HasMask() bool {
	return len(qs.Mask) > 0
}

func (qs QuestionSet) InMask(n int) bool {
	for _, v := range qs.Mask {
		if v == n {
			return true
		}
	}

	return false
}

func (qs QuestionSet) Shuffle() QuestionSet {
	for i := range qs.Questions {
		j := rand.Intn(i + 1)
		qs.Questions[i], qs.Questions[j] = qs.Questions[j], qs.Questions[i]
	}

	return qs
}

type QuestionSets []QuestionSet

func (qss QuestionSets) NextSet(x, y int) (QuestionSet, error) {
	var validSets []QuestionSet

	for _, qs := range qss {
		if qs.ValidSet(x, y) {
			validSets = append(validSets, qs)
		}
	}

	// once you get a list of valid sets, choose one at random
	if l := len(validSets); l > 0 {
		return validSets[rand.Intn(l)].Shuffle(), nil // Intn chooses from [0, n-1] (i.e. Intn(1) is always 0)
	}

	return QuestionSet{}, &ErrNoneValid{x, y}
}

type ErrNoneValid struct {
	X int
	Y int
}

func (e *ErrNoneValid) Error() string {
	return fmt.Sprintf("partisan/questions: Cannot find valid question set: %d, %d", e.X, e.Y)
}
