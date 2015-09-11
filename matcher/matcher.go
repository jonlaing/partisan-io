package matcher

import (
	"database/sql/driver"
	"fmt"
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

// Add records an answer and places it in the map
func (p *PoliticalMap) Add(aMap []int, agree bool) error {
	var sign int

	if agree {
		sign = 1
	} else {
		sign = -1
	}

	for k, v := range aMap {
		if v > 16 {
			return fmt.Errorf("Answer map coordinate out of range at %d! Must be 16 or less. Was %d", k, v)
		}

		if sign > 0 || p[v] > 0 {
			p[v] += sign
		}
	}

        // not sure this is what I want so commenting it out
	// for i := 0; i < 16; i++ {
	// 	if !contains(aMap, i) {
	// 		if sign < 0 || p[i] > 0 {
	// 			p[i] -= sign
	// 		}
	// 	}
	// }

	return nil
}

// Scan satisfies sql.Scanner interface
func (p *PoliticalMap) Scan(src interface{}) error {
	str, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("Cannot scan values for PoliticalMap: %v", src)
	}

	strs := strings.Split(string(str), ",")
	for i, val := range strs {
		var err error
		p[i], err = strconv.Atoi(val)
		if err != nil {
			return err
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

	fmt.Println(str)
	return str, nil
}

// Match returns the % match between two PoliticalMaps
// Only match subquadrants will be compared. If one or both maps has
// 0 points at a subquadrant, it will be ignored.
func Match(p1, p2 PoliticalMap) (float64, error) {
	matchPoints := 0 // Points among matching coordinates
	totalPoints := 0 // total points of all subquadrants in both maps

	for i := range p1 {
		totalPoints += p1[i] + p2[i] // increase total points

		// if both maps have points at this subquadrant
		if p1[i] != 0 && p2[i] != 0 {
			matchPoints += p1[i] + p2[i] // add points of intersecting subquadrants
		}
	}

	return float64(matchPoints) / float64(totalPoints), nil
}

// Enemy returns the % enemy between two PoliticalMaps
// This calculates the heat of quadrants where there
// is no overlap in the two political maps. If both
// maps have points on a quadrant, it is ignored.
func Enemy(p1, p2 PoliticalMap) (float64, error) {
	enemyPoints := 0 // Points among matching coordinates
	totalPoints := 0 // total points of all subquadrants in both maps

	for i := range p1 {
		totalPoints += p1[i] + p2[i] // increase total points

		// Logical XOR: Only add if only one map has points on the quadrant
		if (p1[i] != 0) != (p2[i] != 0) {
			enemyPoints += p1[i] + p2[i] // add points of intersecting subquadrants
		}
	}

	return float64(enemyPoints) / float64(totalPoints), nil
}

func contains(a []int, i int) bool {
	for _, v := range a {
		if i == v {
			return true
		}
	}
	return false
}
