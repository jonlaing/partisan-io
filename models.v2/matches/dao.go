package matches

import (
	"sort"
	"strconv"
	"time"

	"partisan/location"
	"partisan/matcher"

	"partisan/models.v2/friendships"
	"partisan/models.v2/users"

	"github.com/jinzhu/gorm"
)

const centerBounds = 50

func List(user users.User, search SearchBinding, db *gorm.DB) (matches Matches, err error) {
	var users []users.User
	offset := (search.Page - 1) * 24

	err = db.Where("id != ?", user.ID).Scopes(
		inBounds(user),
		inGeoRadius(user.Latitude, user.Longitude, search.Degrees()),
		withGender(search),
		withAgeRange(search.MinAge, search.MaxAge),
		lookingFor(search.LookingFor),
		noFriends(user),
	).Limit(24).Offset(offset).Find(&users).Error

	if err == nil {
		for _, u := range users {
			match, _ := matcher.Match(user.PoliticalMap, u.PoliticalMap)
			matches = append(matches, Match{
				User:  u,
				Match: matcher.ToHuman(match),
			})
		}
	}

	sort.Sort(sort.Reverse(matches))

	return
}

func inBounds(user users.User) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// MATCH BOUNDS
		minX := user.CenterX - centerBounds
		maxX := user.CenterX + centerBounds
		minY := user.CenterY - centerBounds
		maxY := user.CenterY + centerBounds

		return db.Where("centerx > ? AND centerx < ?", minX, maxX).
			Where("centery > ? AND centery < ?", minY, maxY)
	}
}

func inGeoRadius(lat, long, rad float64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		minLat, maxLat, minLong, maxLong, err := location.Bounds(lat, long, rad)
		if err != nil {
			return db
		}

		return db.Where("latitude > ? AND latitude < ?", minLat, maxLat).
			Where("longitude > ? AND longitude < ?", minLong, maxLong)
	}
}

func withGender(search SearchBinding) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// make sure you can't overload on gender
		group, err := search.GenderGroup()
		if err != nil {
			if len(search.Gender) > 0 && len(search.Gender) < 256 {
				return db.Where("gender ILIKE ?", "%"+search.Gender+"%")
			}

			return db
		}

		return db.Where("gender IN (?)", group)
	}
}

func withAgeRange(minAge, maxAge int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if minAge > 0 {
			year := strconv.Itoa(time.Now().Year() - minAge)
			yearParse, err := time.Parse("2006", year)
			if err == nil {
				db = db.Where("birthdate < ?", yearParse.Format("2006-01-02"))
			}
		}

		if maxAge > 0 {
			year := strconv.Itoa(time.Now().Year() - maxAge)
			yearParse, err := time.Parse("2006", year)
			if err == nil {
				db = db.Where("birthdate > ?", yearParse.Format("2006-01-02"))
			}
		}

		return db
	}
}

func noFriends(user users.User) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		dbClone := db.New() // Since we're in a scope, we need a fresh db so as not to apply the scope to ListConfirmedIDsByUserID
		friendIDs, _ := friendships.ListConfirmedIDsByUserID(user.ID, dbClone)
		if len(friendIDs) > 0 {
			return db.Where("id NOT IN (?)", friendIDs)
		}

		return db
	}
}

func lookingFor(looking int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if looking == 0 {
			return db
		}

		var out []int
		n := uint(looking)
		for i := uint(1); i <= n; i++ {
			if n|i == n {
				out = append(out, int(i))
			}
		}

		if len(out) > 0 {
			return db.Where("looking_for IN (?)", out)
		}

		return db
	}
}
