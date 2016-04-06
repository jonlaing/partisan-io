package v1

import (
	"fmt"
	"math"
	"net/http"
	"partisan/auth"
	"partisan/dao"
	"partisan/db"
	"partisan/matcher"
	m "partisan/models"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	earthRadius float64 = 3959 // in miles
)

// MatchResp is the JSON schema we respond with
type MatchResp struct {
	User  m.User  `json:"user"`
	Match float64 `json:"match"`
}

// MatchCollectionResp is the JSON collection schema we respond with
type MatchCollectionResp []MatchResp

// Len satisfies sort.Interface
func (ms MatchCollectionResp) Len() int {
	return len(ms)
}

// Less satisfies sort.Interface
func (ms MatchCollectionResp) Less(a, b int) bool {
	return ms[a].Match < ms[b].Match
}

// Swap satisfies sort.Interface
func (ms MatchCollectionResp) Swap(a, b int) {
	ms[b], ms[a] = ms[a], ms[b]
}

// MatchesIndex returns a list of matches orderd by location and match percentage
func MatchesIndex(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		handleError(err, c)
		return
	}

	page := getPage(c)

	// search radius
	distance := c.DefaultQuery("distance", "25")
	radius, err := convertMilesToDegrees(distance)
	if err != nil {
		handleError(err, c)
		return
	}

	// Gender
	gender := c.Query("gender")

	// Age
	minA := c.DefaultQuery("minAge", "-1")
	maxA := c.DefaultQuery("maxAge", "-1")
	minAge, err := strconv.Atoi(minA)
	if err != nil {
		minAge = -1
	}
	maxAge, err := strconv.Atoi(maxA)
	if err != nil {
		maxAge = -1
	}

	db.LogMode(true)
	users, err := dao.GetMatches(user, gender, minAge, maxAge, radius, page, db)
	if err != nil {
		handleError(err, c)
		return
	}
	db.LogMode(false)

	fmt.Println(users)

	var matches MatchCollectionResp
	for _, u := range users {
		match, err := matcher.Match(user.PoliticalMap, u.PoliticalMap)
		// The only error that can happen here is if neither user has a map, just ignore
		if err == nil {
			match = float64(int(match*1000)) / 10
			matches = append(matches, MatchResp{User: u, Match: match})
		}
	}

	sort.Sort(sort.Reverse(matches))

	c.JSON(http.StatusOK, matches)
}

// UserHasMap will return a 200 if a user has a valid map,
// otherwise will return a 404
func UserHasMap(c *gin.Context) {
	user, err := auth.CurrentUser(c)
	if err != nil {
		handleError(err, c)
		return
	}

	if user.PoliticalMap.IsEmpty() {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// converts miles from the context params (string) into coordinate degrees
func convertMilesToDegrees(m string) (float64, error) {
	miles, err := strconv.Atoi(m)
	if err != nil {
		return 0, err
	}

	geoBounds := float64(miles) / earthRadius * float64(180) / math.Pi
	return geoBounds, nil
}
