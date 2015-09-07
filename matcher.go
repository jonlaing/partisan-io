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
type Map [16]int

// PoliticalMap holds the map, because postgresql doesn't like storing slices, but is fine with structs (apparently?)
type PoliticalMap struct {
	Map Map
}

// Add records an answer and places it in the map
func (p *PoliticalMap) Add(a Answer) error {
	var sign int

	if a.Agree {
		sign = 1
	} else {
		sign = -1
	}

	for k, v := range a.Map {
		if v > 16 {
			return fmt.Errorf("Answer map coordinate out of range at %d! Must be 16 or less. Was %d", k, v)
		}

		if sign > 0 || p.Map[v] > 0 {
			p.Map[v] += sign
		}
	}

	for i := 0; i < 16; i++ {
		if !contains(a.Map, i) {
			if sign < 0 || p.Map[i] > 0 {
				p.Map[i] -= sign
			}
		}
	}

	return nil
}

// Match returns the % match between two PoliticalMaps
// Only match subquadrants will be compared. If one or both maps has
// 0 points at a subquadrant, it will be ignored.
func Match(p1, p2 PoliticalMap) (float64, error) {
	matchPoints := 0 // Points among matching coordinates
	totalPoints := 0 // total points of all subquadrants in both maps

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
	enemyPoints := 0 // Points among matching coordinates
	totalPoints := 0 // total points of all subquadrants in both maps

	for i := range p1.Map {
		totalPoints += p1.Map[i] + p2.Map[i] // increase total points

		// Logical XOR: Only add if only one map has points on the quadrant
		if (p1.Map[i] != 0) != (p2.Map[i] != 0) {
			enemyPoints += p1.Map[i] + p2.Map[i] // add points of intersecting subquadrants
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
