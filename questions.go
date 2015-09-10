package main

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
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
			Prompt: "Capitalism and free-market principles are most compatible with human nature",
			Map:    []int{2, 3, 6, 7, 10, 11, 14, 15},
		},
		Question{
			Prompt: "Capitalism exacerbates negative aspects of humanity like greed and short-sighted planning.",
			Map:    []int{0, 1, 4, 5, 8, 9, 12, 13},
		},
		Question{
			Prompt: "Society requires hierarchical organization to remain coherent.",
			Map:    []int{0, 1, 2, 3, 4, 5, 6, 7},
		},
		Question{
			Prompt: "One should always remain highly skeptical of the government.",
			Map:    []int{7, 8, 9, 10, 11, 12, 13, 14, 15},
		},
		Question{
			Prompt: "The free-market will regulate itself.",
			Map:    []int{3, 7, 11, 15},
		},
		Question{
			Prompt: "Capitalism cannot be reformed, only abolished.",
			Map:    []int{0, 4, 8, 12},
		},
		Question{
			Prompt: "Both capitalism and the state should be abolished.",
			Map:    []int{8, 12},
		},
		Question{
			Prompt: "Capitalism can be made more efficient and ethical through well thought out reform.",
			Map:    []int{5, 6},
		},
		Question{
			Prompt: "Loyalty to one's country is of utmost importance.",
			Map:    []int{2, 3, 6, 7},
		},
		Question{
			Prompt: "The police cannot be trusted.",
			Map:    []int{0, 4, 8, 12, 13, 14, 15},
		},
		Question{
			Prompt: "Voting is not just a right; it's a duty.",
			Map:    []int{2, 3, 6, 7, 10, 11},
		},
		Question{
			Prompt: "It is the duty of the enlightened world to spread democracy, by force if necessary.",
			Map:    []int{2, 3, 6, 7, 10, 11},
		},
		Question{
			Prompt: "Feminism: yes or no?",
			Map:    []int{0, 1, 4, 5, 8, 9, 12, 13},
		},
	}
}

// QuestionShow shows a random question
func QuestionShow(c *gin.Context) {
	index := rand.Intn(len(questions) - 1)

	quest := questions[index]

	c.JSON(http.StatusOK, quest)
}
