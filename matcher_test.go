package main

import (
	"testing"
)

func TestAddingPointsToMap(t *testing.T) {
	a1 := Answer{
		Map:   []int{0, 4, 8, 12},
		Agree: true,
	}

	a2 := Answer{
		Map:   []int{3, 7, 11, 15},
		Agree: true,
	}

	p1 := PoliticalMap{}
	p1.Add(a1)
	for i, val := range p1 {
		if i%4 == 0 && val != 1 {
			t.Error("PoliticalMap didn't mark right fields:", p1)
		} else if i%4 != 0 && val != 0 {
			t.Error("PoliticalMap didn't mark right fields:", p1)
		}
	}

	p1.Add(a1)
	for i, val := range p1 {
		if i%4 == 0 && val != 2 {
			t.Error("PoliticalMap didn't mark right fields:", p1)
		} else if i%4 != 0 && val != 0 {
			t.Error("PoliticalMap didn't mark right fields:", p1)
		}
	}

	a1.Agree = false
	p1.Add(a1)
	for i, val := range p1 {
		if i%4 == 0 && val != 1 {
			t.Error("PoliticalMap didn't mark right fields:", p1)
		} else if i%4 != 0 && val != 0 {
			t.Error("PoliticalMap didn't mark right fields:", p1)
		}
	}

	p1.Add(a2)
	for i, val := range p1 {
		if (i+1)%4 == 0 && val != 1 {
			t.Error("PoliticalMap didn't mark right fields:", p1)
		}
	}
}

func TestMatching(t *testing.T) {
	p1 := PoliticalMap{
		1, 0, 0, 0,
		1, 0, 0, 0,
		2, 2, 0, 0,
		5, 3, 1, 0,
	}

	p2 := PoliticalMap{
		0, 0, 0, 0,
		0, 0, 2, 0,
		0, 0, 3, 0,
		0, 4, 6, 0,
	}

	p3 := PoliticalMap{
		0, 1, 1, 1,
		0, 1, 1, 1,
		0, 0, 1, 1,
		0, 0, 0, 1,
	}

	m, err := Match(p1, p1)
	if err != nil {
		t.Error(err)
	}

	if m != 1.0 {
		t.Error("Match incorrect:", m)
	}

	m, err = Match(p1, p2)
	if err != nil {
		t.Error(err)
	}

	if m != 14.0/30.0 {
		t.Error("Match incorrect:", m)
	}

	m, err = Match(p1, p3)
	if err != nil {
		t.Error(err)
	}

	if m != 0.0 {
		t.Error("Match incorrect:", m)
	}
}
