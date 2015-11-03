package v1

import (
	"math"
	"net/http"
	"partisan/auth"
	"partisan/dao"
	"partisan/db"
	"partisan/matcher"
	m "partisan/models"
	"sort"
	"strconv"

	"partisan/Godeps/_workspace/src/github.com/gin-gonic/gin"
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
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	// Offset
	p := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(p)
	if err != nil {
		page = 1
	}

	// search radius
	distance := c.DefaultQuery("distance", "25")
	radius, err := convertMilesToDegrees(distance)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
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

	users, err := dao.GetMatches(user, gender, minAge, maxAge, radius, page, db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	var matches MatchCollectionResp
	for _, u := range users {
		match, _ := matcher.Match(user.PoliticalMap, u.PoliticalMap)
		match = float64(int(match*1000)) / 10
		matches = append(matches, MatchResp{User: u, Match: match})
	}

	sort.Sort(sort.Reverse(matches))

	c.JSON(http.StatusOK, matches)
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
