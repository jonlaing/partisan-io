package main

import (
	"fmt"
	"math/rand"
	"net/http"

	"partisan/matcher"

	"partisan/Godeps/_workspace/src/github.com/gin-gonic/gin"
)

// Question holds questions to guage a user's political leanings
type Question struct {
	Prompt string `json:"prompt"`
	Map    []int  `json:"map"` // defined in matcher.go
}

var questions []Question

func init() {
	questions = []Question{
		Question{
			Prompt: "Of all economic systems, Capitalism is most compatible with human nature.",
			Map:    []int{2, 3, 6, 7, 10, 11, 14, 15},
		},
		Question{
			Prompt: "Despite corruption, governments are, at their core, for the good of society.",
			Map:    []int{0, 1, 2, 3, 4, 5, 6, 7},
		},
		Question{
			Prompt: "Loyalty to one's country is of utmost importance.",
			Map:    []int{2, 3, 6, 7},
		},
		Question{
			Prompt: "Voting is not just a right; it's a duty.",
			Map:    []int{2, 3, 6, 7, 10, 11},
		},
		Question{
			Prompt: "Men are more suited for positions of leadership, while women are more suited for positions of nurturing.",
			Map:    []int{2, 3, 6, 7, 10, 11, 14, 15},
		},
		Question{
			Prompt: "Violence is never the answer.",
			Map:    []int{5, 6, 9, 10},
		},
		Question{
			Prompt: "The police, in general, are good.",
			Map:    []int{1, 2, 3, 5, 6, 7, 10, 11},
		},
		Question{
			Prompt: "You should always support the troops.",
			Map:    []int{2, 3, 6, 7, 10, 11, 15},
		},
		Question{
			Prompt: "People over profits.",
			Map:    []int{0, 1, 4, 5, 8, 9, 12, 13},
		},
		Question{
			Prompt: "Healthcare should be free of charge.",
			Map:    []int{0, 1, 4, 5, 8, 9, 12, 13},
		},
		Question{
			Prompt: "The traditional institution of marriage is sacred.",
			Map:    []int{2, 3, 7, 11},
		},
		Question{
			Prompt: "Prisons make society safer.",
			Map:    []int{1, 2, 3, 5, 6, 7, 10, 11},
		},
		Question{
			Prompt: "The best way to institute change in the system is to work within it.",
			Map:    []int{1, 3, 5, 6, 7, 10, 11, 14, 15},
		},
		Question{
			Prompt: "The borders should be patrolled to prevent illegal immigration.",
			Map:    []int{2, 3, 6, 7, 10, 11, 14, 15},
		},
		Question{
			Prompt: "When someone is poor, it is mostly their fault.",
			Map:    []int{2, 3, 6, 7, 10, 11, 14, 15},
		},
		Question{
			Prompt: "Government/State is the best way to complete large-scale projects such as building roads.",
			Map:    []int{0, 1, 2, 3, 4, 5, 6, 7, 10, 11},
		},
		Question{
			Prompt: "Management/CEOs deserve a higher salary than their employees.",
			Map:    []int{2, 3, 6, 7, 10, 11, 14, 15},
		},
		Question{
			Prompt: "Social welfare programs, such as food stamps, just make people lazy.",
			Map:    []int{2, 3, 6, 7, 10, 11, 14, 15},
		},
		Question{
			Prompt: "Housing should be considered a human right.",
			Map:    []int{0, 1, 4, 5, 8, 9, 12, 13},
		},
		Question{
			Prompt: "All education, including higher education, should be free",
			Map:    []int{0, 1, 4, 5, 8, 9, 12, 13},
		},
		Question{
			Prompt: "Some races/ethnicities are superior to others.",
			Map:    []int{2, 3, 7, 11, 14, 15},
		},
		Question{
			Prompt: "The best way to end a recession is to remove restrictions on wealth creators.",
			Map:    []int{3, 7, 11, 15},
		},
		Question{
			Prompt: "Those with higher income should be allowed to pay for better healthcare than those with lower income.",
			Map:    []int{2, 3, 6, 7, 10, 11, 14, 15},
		},
		Question{
			Prompt: "Authority should always be questioned.",
			Map:    []int{8, 9, 10, 11, 12, 13, 14, 15},
		},
		Question{
			Prompt: "The primary function of education should be to prepare students for their careers.",
			Map:    []int{2, 3, 6, 7, 10, 11, 14, 15},
		},
		Question{
			Prompt: "A country shouldn't concern itself with the problems of refugees.",
			Map:    []int{3, 7, 11, 15},
		},
		Question{
			Prompt: "Markets are a poor way to distribute resources.",
			Map:    []int{0, 4, 5, 8, 12},
		},
		Question{
			Prompt: "Free Markets = Free People",
			Map:    []int{3, 7, 11, 15},
		},
		Question{
			Prompt: "People should be able to lead their own lives, free of government intervention.",
			Map:    []int{8, 9, 10, 11, 12, 13, 14, 15},
		},
	}

	var total matcher.PoliticalMap
	for _, q := range questions {
		total.Add(q.Map, true)
	}

	fmt.Println(total[0:4])
	fmt.Println(total[4:8])
	fmt.Println(total[8:12])
	fmt.Println(total[12:])

}

// QuestionShow shows a random question
func QuestionShow(c *gin.Context) {
	index := rand.Intn(len(questions))

	quest := questions[index]

	c.JSON(http.StatusOK, quest)
}
