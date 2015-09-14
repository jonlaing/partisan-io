package v1

import (
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"partisan/auth"
	"partisan/db"
	"partisan/matcher"
	m "partisan/models"
	"sort"
)

const (
	earthRadius  float64 = 3959 // in miles
	searchBounds float64 = float64(5) / earthRadius * float64(180) / math.Pi
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
	db, err := db.InitDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	user, err := auth.CurrentUser(c, &db)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	// LATITUDE
	minLat := user.Latitude - searchBounds
	maxLat := user.Latitude + searchBounds
	if minLat > 90 {
		minLat = 90 - minLat
	} else if minLat < -90 {
		minLat = 90 + minLat
	}
	if maxLat > 90 {
		maxLat = 90 - maxLat
	} else if maxLat < -90 {
		maxLat = 90 + maxLat
	}

	// LONGITUDE
	minLong := user.Longitude - searchBounds
	maxLong := user.Longitude + searchBounds
	if minLong > 180 {
		minLong = 180 - minLong
	} else if minLong < -180 {
		minLong = 180 + minLong
	}
	if maxLong > 180 {
		maxLong = 180 - maxLong
	} else if maxLong < -180 {
		maxLong = 180 + maxLong
	}

	// MATCH BOUNDS
	// minX := user.CenterX - 10
	// maxX := user.CenterX + 10
	// minY := user.CenterY - 10
	// maxY := user.CenterY + 10

	var users []m.User
	// if err := db.Where("latitude > ? AND latitude < ? AND longitude > ? AND longitude < ? AND center_x > ? AND center_x < ? AND center_y > ? AND center_y < ?",
	// 	minLat, maxLat, minLong, maxLong, minX, maxX, minY, maxY).Limit(50).Find(&users).Error; err != nil {
	// 	c.AbortWithError(http.StatusNotAcceptable, err)
	// 	return
	// }
	if err := db.Where("latitude > ? AND latitude < ? AND longitude > ? AND longitude < ? AND id != ?",
		minLat, maxLat, minLong, maxLong, user.ID).Limit(50).Find(&users).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}
	// if err := db.Limit(50).Find(&users).Error; err != nil {
	// 	c.AbortWithError(http.StatusNotAcceptable, err)
	// 	return
	// }

	var matches MatchCollectionResp
	for _, u := range users {
		match, _ := matcher.Match(user.PoliticalMap, u.PoliticalMap)
		match = float64(int(match*1000)) / 10
		matches = append(matches, MatchResp{User: u, Match: match})
	}

	sort.Sort(sort.Reverse(matches))

	c.JSON(http.StatusOK, matches)
}
