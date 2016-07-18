package matcher

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// PoliticalMap is a 4x4 map of poltical opinions.
// Everytime a question is answered, a point is added
// to a sub quadrant.
//
//      ---------------------
//      |  0 |  1 |  2 |  3 |
//      ---------------------
//      |  4 |  5 |  6 |  7 |
//      ---------------------
//      |  8 |  9 | 10 | 11 |
//      ---------------------
//      | 12 | 13 | 14 | 15 |
//      ---------------------
//
type PoliticalMap [16]int

// Add records an answer and places it in the map.
// NOTE: In case you forget, as you have 100 times, yes, the mask is necessary! If, for instance, you disagree with something mapped to
// MMiddleLeft, without the mask, both the Far Left AND the whole Right Wing will get points added! With the mask, only the
// Far Left will get highlighted.
func (p *PoliticalMap) Add(aMap []int, mask []int, agree bool) error {
	var sign int

	if agree {
		sign = 1
	} else {
		sign = -1
	}

	for k, v := range aMap {
		if v > 16 {
			return &ErrOutOfRange{k, v}
		}

		// if there's no mask, or the index is within the mask, add it
		if len(mask) == 0 || contains(mask, v) {
			p[v] += sign
		}
	}

	p.normalize() // if any term is < 0, shift the whole thing up

	return nil
}

// if any part of the PoliticalMap is < 0, then shift everything up,
// so that the smallest value is zero
func (p *PoliticalMap) normalize() {
	smallest := 0
	for _, v := range p {
		if v < smallest {
			smallest = v
		}
	}

	if smallest < 0 {
		for k := range p {
			p[k] -= smallest
		}
	}
}

// Center finds the center of gravity of the map for faster lookups in SQL
// Based on actual center of gravity math for two dimenions.
// There are some people that will come up with a weird coordinate, but it should
// be the minority. Most people should have fairly large clusters in one spot,
// so it'll make the initial SQL look up easier to match them.
// TODO: Actually test this with real people and real questions
func (p *PoliticalMap) Center() (int, int) {
	var x, y, t float64 // t is the total points

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
		return int(math.Ceil(x * 100 / t / 3)), int(math.Ceil(y * 100 / t / 3))
	}

	return 0, 0
}

// IsEmpty will check for empty poltical maps. These cause big problems
// when trying to match
func (p *PoliticalMap) IsEmpty() bool {
	for _, v := range p {
		if v > 0 {
			return false
		}
	}

	return true
}

// Scan satisfies sql.Scanner interface
func (p *PoliticalMap) Scan(src interface{}) error {
	str, ok := src.([]byte)
	if !ok {
		return ErrScan{src}
	}

	strs := strings.Split(string(str), ",")
	for i, val := range strs {
		var err error
		p[i], err = strconv.Atoi(val)
		if err != nil {
			return ErrScan{val}
		}
	}

	return nil
}

// Value satisfies driver.Valuer interface
func (p PoliticalMap) Value() (driver.Value, error) {
	str := ""

	for i, val := range p {
		if i == 0 {
			str += fmt.Sprintf("%d", val)
		} else {
			str += fmt.Sprintf(",%d", val)
		}
	}

	return str, nil
}

// Match returns the % match between two PoliticalMaps
// Only match subquadrants will be compared. If one or both maps has
// 0 points at a subquadrant, it will be ignored.
func Match(p1, p2 PoliticalMap) (float64, error) {
	// if p1.IsEmpty() || p2.IsEmpty() {
	// 	return -1.0, errors.New("One or both maps are empty")
	// }

	matchPoints := 0 // Points among matching coordinates
	totalPoints := 0 // total points of all subquadrants in both maps

	for i := range p1 {
		totalPoints += p1[i] + p2[i] // increase total points

		// if both maps have points at this subquadrant
		if p1[i] != 0 && p2[i] != 0 {
			matchPoints += p1[i] + p2[i] // add points of intersecting subquadrants
		}
	}

	if totalPoints < 1 {
		return 0.0, errors.New("No maps between users")
	}

	// just the simple ratio is too generous so...
	naive := float64(matchPoints) / float64(totalPoints)

	// we put the naive calculation on a new curve (a circle) to
	// exaggerate disagreement. It makes the matcher a little
	// more interesting (otherwise you match with freakin' everyone)
	return -math.Sqrt(1-naive) + 1, nil
}

// Merge merges two or more maps by averaging their values at each subquadrant
func Merge(maps ...PoliticalMap) (p PoliticalMap, err error) {
	numMaps := len(maps)
	if numMaps < 2 {
		return p, errors.New("Cannot merge less than 2 maps")
	}

	for i := 0; i < 16; i++ {
		sum := 0
		for _, m := range maps {
			sum += m[i]
		}

		p[i] = sum / numMaps
	}

	return
}

func contains(a []int, i int) bool {
	for _, v := range a {
		if i == v {
			return true
		}
	}
	return false
}

type ErrScan struct {
	Val interface{}
}

func (e ErrScan) Error() string {
	return fmt.Sprintf("partisan/matcher: Can't scan values: %v", e.Val)
}

type ErrOutOfRange struct {
	Index int
	Val   int
}

func (e *ErrOutOfRange) Error() string {
	return fmt.Sprintf("partisan/matcher: Answer map coordinate out of range at %d! Must be 16 or less. Was %d", e.Index, e.Val)
}

func ToHuman(match float64) float64 {
	return float64(int(match*1000)) / 10
}
