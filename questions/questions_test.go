package questions

import "testing"

func TestQuestions(t *testing.T) {
	if num := len(InitialQuestions); num < 8 {
		t.Error("InitialQuestions should have 8 questions, got:", num)
	}

	if num := len(LeftWingQuestions); num < 8 {
		t.Error("LeftWingQuestions should have 8 questions, got:", num)
	}

	if num := len(RightWingQuestions); num < 8 {
		t.Error("RightWingQuestions should have 8 questions, got:", num)
	}

	if num := len(AuthoritarianSocialistQuestions); num < 8 {
		t.Error("AuthoritarianSocialistQuestions should have 8 questions, got:", num)
	}

	if num := len(LibertarianSocialistQuestions); num < 8 {
		t.Error("LibertarianSocialistQuestions should have 8 questions, got:", num)
	}

	if num := len(AuthoritarianQuestions); num < 8 {
		t.Error("AuthoritarianQuestions should have 8 questions, got:", num)
	}

	if num := len(LibertarianQuestions); num < 8 {
		t.Error("LibertarianQuestions should have 8 questions, got:", num)
	}

	if num := len(AuthoritarianCapitalistQuestions); num < 8 {
		t.Error("AuthoritarianCapitalistQuestions should have 8 questions, got:", num)
	}

	if num := len(LibertarianCapitalistQuestions); num < 8 {
		t.Error("LibertarianCapitalistQuestions should have 8 questions, got:", num)
	}
}

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

// func TestNextSet(t *testing.T) {
// 	qss1 := QuestionSets{
// 		QuestionSet{Mask: []int{0}},
// 		QuestionSet{Mask: []int{1}},
// 		QuestionSet{Mask: []int{1}},
// 		QuestionSet{Mask: []int{1}},
// 	}

// 	qs1, err := qss1.NextSet(-1, -1)

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if qs1.Mask[0] != 0 {
// 		t.Error("Got wrong question set:", qs1)
// 	}

// 	qss2 := QuestionSets{
// 		QuestionSet{Mask: []int{0}},
// 		QuestionSet{Mask: []int{1}},
// 		QuestionSet{Mask: []int{1}},
// 		QuestionSet{Mask: []int{1}},
// 	}

// 	qs2, err := qss2.NextSet(-1, -1)

// 	if err == nil {
// 		t.Error("Should have thrown an error, got:", qs2)
// 	}
// }
