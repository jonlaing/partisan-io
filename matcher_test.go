package main

import (
	"testing"
)

func TestAddingPointsToMap(t *testing.T) {
	p1 := PoliticalMap{}
	p2 := PoliticalMap{}

	a1 := Answer{
		Index: 4,
		Agree: true,
	}
	a2 := Answer{
		Index: 4,
		Agree: false,
	}

	a3 := Answer{
		Index: 16,
		Agree: true,
	}

	p1.Add(a1)
	p2.Add(a2)

	if p1[4] != 1 {
		t.Error("1: Map at index incorrect:", p1[4])
	}

	if p2[4] != 0 {
		t.Error("2: Map at index incorrect:", p2[4])
	}

	p1.Add(a2)
	if p1[4] != 1 {
		t.Error("3: Map at index incorrect:", p1[4])
	}

	p1.Add(a1)
	if p1[4] != 2 {
		t.Error("4: Map at index incorrect:", p1[4])
	}

	err := p1.Add(a3)
	if err == nil {
		t.Error("Should throw error for out of index:", p1)
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
