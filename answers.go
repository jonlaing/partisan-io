package main

import (
	// "fmt"
	"github.com/gin-gonic/gin"
	// "net/http"
)


// Answer is an answer to a question, included are the coordinates
// of the question, and whether or not the user agreed
// The Map should be in the form of [1,5,9,13] (which would be the entire far-left).
// Check Matcher for more details on the map
type Answer struct {
	Map   []int  `json:"map"` // defined in matcher.go
	Agree bool `json:"agree"`
}

// AnswersUpdate updates the coordinates of user based on question answers
func AnswersUpdate(c *gin.Context) {
	// db, err := initDB()
	// if err != nil {
	// 	c.AbortWithError(http.StatusInternalServerError, err)
	// 	return
	// }
	// defer db.Close()

	// as := []Answer{}
	// user := User{}

	// if err := c.BindJSON(&as); err != nil {
	// 	c.AbortWithError(http.StatusBadRequest, err)
	// 	return
	// }

	// // Get the User
	// userID, ok := c.Get("user_id")
	// if !ok {
	// 	c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("User ID not set"))
	// 	return
	// }

	// if err := db.First(&user, userID).Error; err != nil {
	// 	c.AbortWithError(http.StatusNotFound, err)
	// 	return
	// }

	// var mapErrors []map[string]interface{}
	// for _, a := range as {
	// 	err := user.PoliticalMap.Add(a)
	// 	if err != nil {
	// 		mapErrors = append(mapErrors, map[string]interface{}{"index": a.Index, "error": err.Error()})
	// 	}
	// }

	// if err := db.Save(&user).Error; err != nil {
	// 	c.AbortWithError(http.StatusNotAcceptable, err)
	// 	return
	// }

	// if len(mapErrors) > 0 {
	// 	c.JSON(http.StatusOK, mapErrors)
	// } else {
	// 	c.JSON(http.StatusOK, gin.H{"message": "updated"})
	// }
}
