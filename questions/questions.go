package questions

// Question holds questions to guage a user's political leanings
type Question struct {
	Prompt string `json:"prompt"`
	Map    []int  `json:"map"` // defined in matcher.go
}

// Questions are in groups of four to contain pro-state, anti-state, pro-capital, and anti-capital sentiments
type Questions [4]Question

// QuestionSet is a collection of questions. Using ValidSet, the algorithm can determine whether this is an appropriate question given
// a User's current Center.
type QuestionSet struct {
	Mask      []int               `json:"-"`
	Questions Questions           `json:"questions"`
	ValidSet  func(int, int) bool `json:"-"` // ValidSet returns whether the provided coordinates match with the target of the QuestionSet
}

func (qs QuestionSet) HasMask() bool {
	return len(qs.Mask) > 0
}

func (qs QuestionSet) InMask(n int) bool {
	for _, v := range qs.Mask {
		if v == n {
			return true
		}
	}

	return false
}

type QuestionSets []QuestionSet

func (qss QuestionSets) NextSet(x, y int) (QuestionSet, error) {
	for _, qs := range qss {
		if qs.ValidSet(x, y) {
			return qs, nil
		}
	}

	return QuestionSet{}, &ErrNoneValid{}
}

type ErrNoneValid struct{}

func (e *ErrNoneValid) Error() string {
	return "partisan/questions: Cannot find valid question set"
}
