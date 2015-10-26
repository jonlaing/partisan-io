package matcher

import (
	"testing"
)

type answer struct {
	Map   []int
	Agree bool
}

func TestAddingPointsToMap(t *testing.T) {
	a1 := answer{
		Map:   []int{0, 4, 8, 12},
		Agree: true,
	}

	a2 := answer{
		Map:   []int{3, 7, 11, 15},
		Agree: true,
	}

	p1 := PoliticalMap{}
	p1.Add(a1.Map, a1.Agree)
	for i, val := range p1 {
		if i%4 == 0 && val != 1 {
			t.Error("PoliticalMap didn't mark right fields:", p1)
		} else if i%4 != 0 && val != 0 {
			t.Error("PoliticalMap didn't mark right fields:", p1)
		}
	}

	p1.Add(a1.Map, a1.Agree)
	for i, val := range p1 {
		if i%4 == 0 && val != 2 {
			t.Error("PoliticalMap didn't mark right fields:", p1)
		} else if i%4 != 0 && val != 0 {
			t.Error("PoliticalMap didn't mark right fields:", p1)
		}
	}

	a1.Agree = false
	p1.Add(a1.Map, a1.Agree)
	for i, val := range p1 {
		if i%4 == 0 && val != 1 {
			t.Error("PoliticalMap didn't mark right fields:", p1)
		} else if i%4 != 0 && val != 0 {
			t.Error("PoliticalMap didn't mark right fields:", p1)
		}
	}

	p1.Add(a2.Map, a2.Agree)
	for i, val := range p1 {
		if (i+1)%4 == 0 && val != 1 {
			t.Error("PoliticalMap didn't mark right fields:", p1)
		}
	}
}

func TestNormalize(t *testing.T) {
	a := answer{
		Map:   []int{0, 1, 2, 3},
		Agree: false,
	}

	p := PoliticalMap{}
	p.Add(a.Map, a.Agree)

	for k, v := range p {
		if k >= 0 && k <= 3 && v != 0 {
			t.Error("p[", k, "] should have been zero, was:", v)
		} else if k > 3 && v <= 0 {
			t.Error("p[", k, "] should have been greater than zero, was:", v)
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

func TestCenter(t *testing.T) {
	p1 := PoliticalMap{
		1, 1, 1, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
	}

	if x, y := p1.Center(); x != 0 || y != 0 {
		t.Error("Incorrect center: ", x, y)
	}

	p2 := PoliticalMap{
		1, 1, 1, 1,
		1, 0, 0, 1,
		1, 0, 0, 1,
		1, 1, 1, 1,
	}

	if x, y := p2.Center(); x != 0 || y != 0 {
		t.Error("Incorrect center: ", x, y)
	}
}

func TestScan(t *testing.T) {
	s1 := []byte("0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0")
	p1 := PoliticalMap{}

	if err := p1.Scan(s1); err != nil {
		t.Error(err)
	}

	for k, v := range p1 {
		if v != 0 {
			t.Error("Value at", k, "should be 0. Was:", v)
		}
	}

	s2 := []byte("1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1")
	p2 := PoliticalMap{}

	if err := p2.Scan(s2); err != nil {
		t.Error(err)
	}

	for k, v := range p2 {
		if v != 1 {
			t.Error("Value at", k, "should be 1. Was:", v)
		}
	}

	s3 := []byte("bad string, it can't be parsed")
	p3 := PoliticalMap{}

	if err := p3.Scan(s3); err == nil {
		t.Error("Scanning bad string should have produced an error")
	}

	s4 := PoliticalMap{}
	p4 := PoliticalMap{}

	if err := p4.Scan(s4); err == nil {
		t.Error("Scanning non-byte slice should have produced an error")
	}
}

func TestValue(t *testing.T) {
	p1 := PoliticalMap{
		1, 1, 1, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
	}

	// skipping error, because there's no way to return it in this
	// implementation of Value
	if str, _ := p1.Value(); str != "1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1" {
		t.Error("Value returned", str)
	}
}
