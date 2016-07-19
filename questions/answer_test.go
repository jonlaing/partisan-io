package questions

import (
	"partisan/matcher"
	"testing"
)

type answerCase struct {
	name       string
	acceptable [][]int
}

var answerCases = []answerCase{
	{"Anarchist", [][]int{MAntiState, MAntiCapital, MAntiAuthoritarian, MFarLeft, MLibSocialist}},
	{"Authoritarian Right", [][]int{MProState, MProCapital, MAuthoritarian, MFarRight}},
	{"Authoritarian Communist", [][]int{MProState, MAntiCapital, MAuthoritarian, MFarLeft}},
	{"Liberal", [][]int{MProState, MProCapital, MMiddleAuthoritarian, MMiddleRight}},
	{"Confused", [][]int{
		MProState,
		MAntiState,
		MProCapital,
		MAntiCapital,
		MAuthoritarian,
		MMiddleAuthoritarian,
		MMiddleAntiAuthoritarian,
		MAntiAuthoritarian,
		MFarLeft,
		MMiddleLeft,
		MFarRight,
		MMiddleRight,
	}},
}

func TestWeights(t *testing.T) {
	for _, set := range Sets {
		x, y := set.center()
		t.Log("Question Set Center: (", x, ",", y, ")")
	}
}

func TestFullWeight(t *testing.T) {
	pMap := matcher.PoliticalMap{}
	for _, set := range Sets {
		for _, q := range set.Questions {
			pMap.Add(q.Map, set.Mask, true)
		}
	}

	x, y := pMap.Center()

	t.Log("Total Partisan Lean:", x, y)
}

func TestAnswerCases(t *testing.T) {
	for _, aCase := range answerCases {
		answerCount := 0
		pMap := matcher.PoliticalMap{}
		dx := 0
		dy := 0

		// 4 sets of 4 is 16 questions altogether
		for answerCount < 4 {
			answerCount++
			x, y := pMap.Center()
			set, err := Sets.NextSet(x, y, dx, dy)
			if err != nil {
				t.Error("Unexpected error when getting next set for", aCase.name, ":", err)
				continue
			}

			for _, q := range set.Questions {
				dx, dy, err = pMap.Add(q.Map, set.Mask, hasMap(q.Map, aCase.acceptable...))
				if err != nil {
					t.Error("Unexpected error adding to map:", err)
				}
			}
		}

		x, y := pMap.Center()
		t.Log(aCase.name, "has center of (", x, ",", y, ")")
		t.Logf("%s has PoliticalMap:\n%s", aCase.name, pMap.ToHuman())
	}
}

func hasMap(m []int, mTypes ...[]int) bool {
	for _, mType := range mTypes {
		matched := true
		for i := range m {
			if m[i] != mType[i] {
				matched = false
				break
			}
		}

		if matched {
			return true
		}
	}

	return false
}
