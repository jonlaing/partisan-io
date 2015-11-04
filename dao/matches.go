package dao

import (
	"errors"
	m "partisan/models"
	"strconv"
	"time"

	"partisan/Godeps/_workspace/src/github.com/jinzhu/gorm"
)

const (
	centerBounds int = 50
)

func GetMatches(user m.User, gender string, minAge, maxAge int, radius float64, page int, db *gorm.DB) (users []m.User, err error) {
	offset := (page - 1) * 24

	minLat, maxLat, minLong, maxLong, err := geoBounds(user.Latitude, user.Longitude, radius)
	if err != nil {
		return users, err
	}

	// MATCH BOUNDS
	minX := user.CenterX - centerBounds
	maxX := user.CenterX + centerBounds
	minY := user.CenterY - centerBounds
	maxY := user.CenterY + centerBounds

	query := db.Where("id != ?", user.ID)
	query = query.Where("latitude > ? AND latitude < ?", minLat, maxLat)
	query = query.Where("longitude > ? AND longitude < ?", minLong, maxLong)
	query = query.Where("center_x > ? AND center_x < ?", minX, maxX)
	query = query.Where("center_y > ? AND center_y < ?", minY, maxY)

	// make sure you can't overload on gender
	if len(gender) > 0 && len(gender) < 256 {
		query = query.Where("gender ILIKE ?", "%"+gender+"%")
	}

	if minAge > -1 {
		year := strconv.Itoa(time.Now().Year() - minAge)
		yearParse, err := time.Parse("2006", year)
		if err == nil {
			query = query.Where("birthdate < ?::date", yearParse)
		}
	}

	if maxAge > -1 {
		year := strconv.Itoa(time.Now().Year() - maxAge)
		yearParse, err := time.Parse("2006", year)
		if err == nil {
			query = query.Where("birthdate > ?::date", yearParse)
		}
	}

	friendIDs, _ := ConfirmedFriendIDs(user, db)
	if len(friendIDs) > 0 {
		query = query.Where("id NOT IN (?)", friendIDs)
	}

	err = query.Limit(24).Offset(offset).Find(&users).Error

	return
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
