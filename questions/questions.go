package questions

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"partisan/matcher"
)

var Sets QuestionSets

const radius = 17

func init() {
	Sets = QuestionSets{}
	Sets = append(Sets, InitialQuestions...)
	Sets = append(Sets, LeftWingQuestions...)
	Sets = append(Sets, RightWingQuestions...)
	Sets = append(Sets, AuthoritarianQuestions...)
	Sets = append(Sets, LibertarianQuestions...)
	Sets = append(Sets, AuthoritarianSocialistQuestions...)
	Sets = append(Sets, LibertarianSocialistQuestions...)
	Sets = append(Sets, AuthoritarianCapitalistQuestions...)
	Sets = append(Sets, LibertarianCapitalistQuestions...)
}

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

func (qs QuestionSet) center() (int, int) {
	var x, y, t float64 // t is the total points
	var p [16]int

	for _, q := range qs.Questions {
		for _, i := range q.Map {
			p[i]++
		}
	}

	// distance from the "origin"
	// moving up by two places at a time (since there's no 0 on the grid)
	xCoef := []int{
		-2, -1, 1, 2,
		-2, -1, 1, 2,
		-2, -1, 1, 2,
		-2, -1, 1, 2,
	}

	yCoef := []int{
		2, 2, 2, 2,
		1, 1, 1, 1,
		-1, -1, -1, -1,
		-2, -2, -2, -2,
	}

	for k, v := range p {
		x += float64(v * xCoef[k])
		y += float64(v * yCoef[k])
		t += float64(v)
	}

	if t > 0 {
		return int(math.Ceil(x * 100 / t)), int(math.Ceil(y * 100 / t))
	}

	return 0, 0
}

type QuestionSets []QuestionSet

func (qss QuestionSets) NextSet(x, y, dx, dy int) (QuestionSet, error) {
	if x != 0 && y != 0 && dx == 0 && dy == 0 {
		return QuestionSet{}, errors.New("No deltas")
	}

	// var validSets QuestionSets

	// for _, qs := range qss {
	// 	if qs.ValidSet(x, y) {
	// 		validSets = append(validSets, qs)
	// 	}
	// }

	return qss.match(x, y, dx, dy)
}

// find a question set that matches the detals
func (qss QuestionSets) match(x, y, dx, dy int) (QuestionSet, error) {
	if len(qss) < 1 {
		return QuestionSet{}, errors.New("No question sets to favor, empty slice")
	}

	if len(qss) == 1 {
		return qss[0], nil
	}

	if dx == 0 && dy == 0 {
		return qss[0], nil
	}

	var validSets []QuestionSet
	bestPoints := 0

	for _, set := range qss {
		xS, yS := set.center()
		points := 0

		if dx < 0 && xS < x {
			points++
		}

		if dx > 0 && xS > x {
			points++
		}

		if dy < 0 && yS < y {
			points++
		}

		if dy > 0 && yS > y {
			points++
		}

		if points == bestPoints {
			validSets = append(validSets, set)
		}

		if points > bestPoints {
			bestPoints = points
			validSets = []QuestionSet{set}
		}
	}

	if l := len(validSets); l > 0 {
		return validSets[rand.Intn(l)].Shuffle(), nil
	}

	return QuestionSet{}, errors.New("Couldn't find any set that satisfied deltas")
}

type ErrNoneValid struct {
	X int
	Y int
}

func (e *ErrNoneValid) Error() string {
	return fmt.Sprintf("partisan/questions: Cannot find valid question set: %d, %d", e.X, e.Y)
}
