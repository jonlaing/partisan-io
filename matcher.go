package main

import (
	"fmt"
)

// Map is a 4x4 map of poltical opinions.
// Everytime a question is answered, a point is added
// to a sub quadrant.
//
//      ---------------------
//      |  1 |  2 |  3 |  4 |
//      ---------------------
//      |  5 |  6 |  7 |  8 |
//      ---------------------
//      |  9 | 10 | 11 | 12 |
//      ---------------------
//      | 13 | 14 | 15 | 16 |
//      ---------------------
//
type Map [16]uint

// PoliticalMap holds the map, because postgresql doesn't like storing arrays, but is fine with structs (apparently?)
type PoliticalMap struct {
	Map Map
}

// Add records an answer and places it in the map
func (p *PoliticalMap) Add(a Answer) error {
	if !a.Agree {
		return nil
	}

	for k, v := range a.Map {
		if v > 16 {
			return fmt.Errorf("Answer map coordinate out of range at %d! Must be 16 or less. Was %d", k, v)
		}

		// add 1 to the index, this score will be used in the index
		p.Map[v]++
	}

	return nil
}

// Match returns the % match between two PoliticalMaps
// Only match subquadrants will be compared. If one or both maps has
// 0 points at a subquadrant, it will be ignored.
func Match(p1, p2 PoliticalMap) (float64, error) {
	var matchPoints, totalPoints uint
	matchPoints = 0 // Points among matching coordinates
	totalPoints = 0 // total points of all subquadrants in both maps

	for i := range p1.Map {
		totalPoints += p1.Map[i] + p2.Map[i] // increase total points

		// if both maps have points at this subquadrant
		if p1.Map[i] != 0 && p2.Map[i] != 0 {
			matchPoints += p1.Map[i] + p2.Map[i] // add points of intersecting subquadrants
		}
	}

	return float64(matchPoints) / float64(totalPoints), nil
}

// Enemy returns the % enemy between two PoliticalMaps
// This calculates the heat of quadrants where there
// is no overlap in the two political maps. If both
// maps have points on a quadrant, it is ignored.
func Enemy(p1, p2 PoliticalMap) (float64, error) {
	var enemyPoints, totalPoints uint
	enemyPoints = 0 // Points among matching coordinates
	totalPoints = 0 // total points of all subquadrants in both maps

	for i := range p1.Map {
		totalPoints += p1.Map[i] + p2.Map[i] // increase total points

		// Logical XOR: Only add if only one map has points on the quadrant
		if (p1.Map[i] != 0) != (p2.Map[i] != 0) {
			enemyPoints += p1.Map[i] + p2.Map[i] // add points of intersecting subquadrants
		}
	}

	return float64(enemyPoints) / float64(totalPoints), nil
}
