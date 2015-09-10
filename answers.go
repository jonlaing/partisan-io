package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Answer is an answer to a question, included are the coordinates
// of the question, and whether or not the user agreed
// The Map should be in the form of [1,5,9,13] (which would be the entire far-left).
// Check Matcher for more details on the map
type Answer struct {
	Map   []int `json:"map" form:"map"` // defined in matcher.go
	Agree bool  `json:"agree" form:"agree"`
}

// AnswersUpdate updates the coordinates of user based on question answers
func AnswersUpdate(c *gin.Context) {
	db, err := initDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	user, err := CurrentUser(c, &db)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	a := Answer{}

	if err := c.BindJSON(&a); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

        fmt.Println(c.Request)
	if len(a.Map) == 0 {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Answer doesn't have map. Probably an error in binding"))
		return
	}

	fmt.Println(a)

	err = user.PoliticalMap.Add(a)
	if err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
	}

	// fmt.Println(user.PoliticalMap)

	if err := db.Save(&user).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}
