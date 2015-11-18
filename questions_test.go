package main

import (
	"math"
	"partisan/matcher"
	"testing"
)

func TestNoBias(t *testing.T) {
	// testing to make sure questions collectively don't have too much of a bias
	tolerance := 3

	for k, qs := range questions {
		var agreer matcher.PoliticalMap
		for _, q := range qs {
			agreer.Add(q.Map, true)
		}

		x, y := agreer.Center()

		if x > tolerance {
			t.Error("Right wing bias! Set:", k)
		}

		if x < -tolerance {
			t.Error("Left-wing bias! Set:", k)
		}

		if y > tolerance {
			t.Error("Statist bias! Set:", k)
		}

		if y < -tolerance {
			t.Error("Anti-state bias! Set", k)
		}

		if t.Failed() {
			t.Log("Center for Set", k, ":", x, ",", y, "\n")
			t.Log("\n",
				agreer[0:4], "\n",
				agreer[4:8], "\n",
				agreer[8:12], "\n",
				agreer[12:], "\n")
		}
	}
}

func TestHealthyQuiz(t *testing.T) {
	for k, qs := range questions {
		for _, i := range []int{0, 3, 12, 15} {
			var hardline matcher.PoliticalMap
			for _, q := range qs {
				agree := false
				for _, v := range q.Map {
					if v == i {
						agree = true
					}
				}

				hardline.Add(q.Map, agree)
			}

			x, y := hardline.Center()

			if math.Abs(float64(x)) < 50 || math.Abs(float64(y)) < 50 {
				t.Errorf("Questions aren't radical enough in Set %d for %d! Center: %d, %d", k, i, x, y)
			}
		}
	}
}
