package main

import (
	"github.com/gin-gonic/gin"
	// "net/http"
)

// Question holds questions to guage a user's political leanings
type Question struct {
	Prompt string `json:"prompt"`
	Map    Map    `json:"map"` // defined in matcher.go
}

func init() {
}

// QuestionsTest is a blah blah
func QuestionsTest(c *gin.Context) {
}
