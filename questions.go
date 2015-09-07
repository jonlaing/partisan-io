package main

import (
	"github.com/gin-gonic/gin"
	// "net/http"
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
			Map:    []int{3, 4, 7, 8, 11, 12, 15, 16},
		},
		Question{
			Prompt: "Capitalism exacerbates negative aspects of humanity like greed and short-sighted planning.",
			Map:    []int{1, 2, 5, 6, 9, 10, 13, 14},
		},
		Question{
			Prompt: "Society requires hierarchical organization to remain coherent.",
			Map:    []int{1, 2, 3, 4, 5, 6, 7, 8},
		},
		Question{
			Prompt: "One should always remain highly skeptical of the government.",
			Map:    []int{8, 9, 10, 11, 12, 13, 14, 15, 16},
		},
		Question{
			Prompt: "The free-market will regulate itself.",
			Map:    []int{4, 8, 12, 16},
		},
		Question{
			Prompt: "Capitalism cannot be reformed, only abolished.",
			Map:    []int{1, 5, 9, 13},
		},
		Question{
			Prompt: "Both capitalism and the state should be abolished.",
			Map:    []int{9, 13},
		},
		Question{
			Prompt: "Capitalism can be made more efficient and ethical through well thought out reform.",
			Map:    []int{6, 7},
		},
		Question{
			Prompt: "Loyalty to one's country is of utmost importance.",
			Map:    []int{3, 4, 7, 8},
		},
		Question{
			Prompt: "The police cannot be trusted.",
			Map:    []int{1, 5, 9, 13, 14, 15, 16},
		},
		Question{
			Prompt: "Voting is not just a right; it's a duty.",
			Map:    []int{3, 4, 7, 8, 11, 12},
		},
		Question{
			Prompt: "It is the duty of the enlightened world to spread democracy, by force if necessary.",
			Map:    []int{3, 4, 7, 8, 11, 12},
		},
		Question{
                  Prompt: "Feminism: yes or no?",
			Map:    []int{1,2,5,6,9,10,13,14},
		},
	}
}

// QuestionsTest is a blah blah
func QuestionsTest(c *gin.Context) {
}
