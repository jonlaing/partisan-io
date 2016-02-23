package questions

import "testing"

func TestHasMask(t *testing.T) {
	qs1 := QuestionSet{}
	qs2 := QuestionSet{Mask: []int{0}}

	if qs1.HasMask() {
		t.Error("QuestionSet #1 should not have a mask", len(qs1.Mask))
	}

	if !qs2.HasMask() {
		t.Error("QuestionSet #2 should have a mask:", len(qs2.Mask))
	}
}

func TestInMask(t *testing.T) {
	nums := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	qs := QuestionSet{Mask: []int{0, 2, 4, 6, 8, 10, 12, 14}}

	for _, n := range nums {
		if n%2 == 0 && !qs.InMask(n) {
			t.Error("Subquadrant", n, "should have been in mask")
		}

		if n%2 != 0 && qs.InMask(n) {
			t.Error("Subquadrant", n, "should not have been in mask")
		}
	}
}

func TestNextSet(t *testing.T) {
	qss1 := QuestionSets{
		QuestionSet{Mask: []int{0}, ValidSet: func(x, y int) bool { return x < 0 && y < 0 }},
		QuestionSet{Mask: []int{1}, ValidSet: func(x, y int) bool { return x > 0 && y < 0 }},
		QuestionSet{Mask: []int{1}, ValidSet: func(x, y int) bool { return x > 0 && y > 0 }},
		QuestionSet{Mask: []int{1}, ValidSet: func(x, y int) bool { return x < 0 && y > 0 }},
	}

	qs1, err := qss1.NextSet(-1, -1)

	if err != nil {
		t.Error(err)
	}

	if qs1.Mask[0] != 0 {
		t.Error("Got wrong question set:", qs1)
	}

	qss2 := QuestionSets{
		QuestionSet{Mask: []int{0}, ValidSet: func(x, y int) bool { return x > 0 && y < 0 }},
		QuestionSet{Mask: []int{1}, ValidSet: func(x, y int) bool { return x > 0 && y < 0 }},
		QuestionSet{Mask: []int{1}, ValidSet: func(x, y int) bool { return x > 0 && y > 0 }},
		QuestionSet{Mask: []int{1}, ValidSet: func(x, y int) bool { return x < 0 && y > 0 }},
	}

	qs2, err := qss2.NextSet(-1, -1)

	if err == nil {
		t.Error("Should have thrown an error, got:", qs2)
	}
}
