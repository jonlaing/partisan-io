package main

import (
	"partisan/matcher"
	"testing"
)

func TestNoBias(t *testing.T) {
	// testing to make sure questions collectively don't have too much of a bias
	tolerance := 3

	var agreer matcher.PoliticalMap
	for _, q := range questions {
		agreer.Add(q.Map, true)
	}

	x, y := agreer.Center()

	if x > tolerance {
		t.Error("Right wing bias!")
	}

	if x < -tolerance {
		t.Error("Left-wing bias!")
	}

	if y > tolerance {
		t.Error("Statist bias!")
	}

	if y < -tolerance {
		t.Error("Anti-state bias!")
	}

	if t.Failed() {
		t.Log("Center:", x, ",", y, "\n")
		t.Log("\n",
			agreer[0:4], "\n",
			agreer[4:8], "\n",
			agreer[8:12], "\n",
			agreer[12:], "\n")
	}
}
