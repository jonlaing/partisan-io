package location

import "errors"

func Bounds(lat, long, rad float64) (minLat, maxLat, minLong, maxLong float64, err error) {
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
