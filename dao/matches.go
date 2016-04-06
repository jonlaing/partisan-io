package dao

import (
	"errors"
	m "partisan/models"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
)

const (
	centerBounds int = 50
)

func GetMatches(user m.User, gender string, minAge, maxAge int, radius float64, page int, db *gorm.DB) (users []m.User, err error) {
	offset := (page - 1) * 24

	err = db.Where("id != ?", user.ID).Scopes(
		inBounds(user),
		inGeoRadius(user.Latitude, user.Longitude, radius),
		withGender(gender),
		withAgeRange(minAge, maxAge),
		noFriends(user),
	).Limit(24).Offset(offset).Find(&users).Error
	if err != nil {
		return users, &ErrNoMatches{err}
	}

	return
}

func inBounds(user m.User) func(db *gorm.DB) *gorm.DB {
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
		minLat, maxLat, minLong, maxLong, err := geoBounds(lat, long, rad)
		if err != nil {
			return db
		}

		return db.Where("latitude > ? AND latitude < ?", minLat, maxLat).
			Where("longitude > ? AND longitude < ?", minLong, maxLong)
	}
}

func withGender(gender string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// make sure you can't overload on gender
		if len(gender) > 0 && len(gender) < 256 {
			return db.Where("gender ILIKE ?", "%"+gender+"%")
		}

		return db
	}
}

func withAgeRange(minAge, maxAge int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if minAge > -1 {
			year := strconv.Itoa(time.Now().Year() - minAge)
			yearParse, err := time.Parse("2006", year)
			if err == nil {
				db = db.Where("birthdate < ?", yearParse.Format("2006-01-02"))
			}
		}

		if maxAge > -1 {
			year := strconv.Itoa(time.Now().Year() - maxAge)
			yearParse, err := time.Parse("2006", year)
			if err == nil {
				db = db.Where("birthdate > ?", yearParse.Format("2006-01-02"))
			}
		}

		return db
	}
}

func noFriends(user m.User) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		dbClone := db.New() // Since we're in a scope, we need a fresh db so as not to apply the scope to ConfirmedFriends
		friendIDs, _ := ConfirmedFriendIDs(user, dbClone)
		if len(friendIDs) > 0 {
			return db.Where("id NOT IN (?)", friendIDs)
		}

		return db
	}
}

func geoBounds(lat, long, rad float64) (minLat, maxLat, minLong, maxLong float64, err error) {
	if rad > 90 {
		err = errors.New("Radius cannot be above 90")
	}

	// LATITUDE
	minLat = lat - rad
	maxLat = lat + rad
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
	minLong = long - rad
	maxLong = long + rad
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

	return
}
