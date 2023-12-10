package talwar

import (
	"fmt"

	"github.com/comfforts/geocode"

	"github.com/comfforts/geohash/pkg/constants"
)

type geoCoder struct {
	charMap         map[string]uint
	encodingQuadMap map[string]string
	decodingQuadMap map[string]string
}

func NewGeoCoder() *geoCoder {
	return &geoCoder{
		charMap:         map[string]uint{"a": 1, "b": 2, "c": 3, "d": 4},
		encodingQuadMap: map[string]string{"00": "a", "01": "b", "11": "c", "10": "d"},
		decodingQuadMap: map[string]string{"a": "00", "b": "01", "c": "11", "d": "10"},
	}
}

func (gc *geoCoder) Encode(lat float64, lon float64, percision int) (string, error) {

	if lat == 0 || lon == 0 {
		return "", constants.ErrInvalidLatLong
	}

	if lat < constants.LAT_MIN || lat > constants.LAT_MAX || lon < constants.LON_MIN || lon > constants.LON_MAX {
		return "", constants.ErrInvalidLatLong
	}

	if percision < 1 || percision > 12 {
		percision = 12
	}

	var hash string
	var latMin float64 = constants.LAT_MIN
	var latMax float64 = constants.LAT_MAX
	var lonMin float64 = constants.LON_MIN
	var lonMax float64 = constants.LON_MAX

	for len(hash) < percision {
		var quad string
		var latMid float64 = (latMin + latMax) / 2
		if lat >= latMid {
			quad = fmt.Sprintf("%s%d", quad, 1)
			latMin = latMid
		} else {
			quad = fmt.Sprintf("%s%d", quad, 0)
			latMax = latMid
		}

		var lonMid float64 = (lonMin + lonMax) / 2
		if lon >= lonMid {
			quad = fmt.Sprintf("%s%d", quad, 1)
			lonMin = lonMid
		} else {
			quad = fmt.Sprintf("%s%d", quad, 0)
			lonMax = lonMid
		}
		hash = fmt.Sprintf("%s%s", hash, gc.encodingQuadMap[quad])
	}

	return hash, nil
}

func (gc *geoCoder) Decode(hash string) (*geocode.RangeBounds, error) {
	return gc.bounds(hash)
}

func (gc *geoCoder) bounds(hash string) (*geocode.RangeBounds, error) {
	if len(hash) == 0 {
		return nil, constants.ErrInvalidGeocode
	}
	var latMin float64 = constants.LAT_MIN
	var latMax float64 = constants.LAT_MAX
	var lonMin float64 = constants.LON_MIN
	var lonMax float64 = constants.LON_MAX
	for i := 0; i < len(hash); i++ {
		char := string(hash[i])
		quad := gc.decodingQuadMap[char]
		latHash := quad[0]
		lonHash := quad[1]

		latMid := (latMin + latMax) / 2
		if latHash == '0' {
			latMax = latMid
		} else {
			latMin = latMid
		}

		lonMid := (lonMin + lonMax) / 2
		if lonHash == '0' {
			lonMax = lonMid
		} else {
			lonMin = lonMid
		}
	}

	var latRange, lonRange geocode.Range
	latRange.Min = latMin
	latRange.Max = latMax
	lonRange.Min = lonMin
	lonRange.Max = lonMax

	var bounds geocode.RangeBounds
	bounds.Latitude = latRange
	bounds.Longitude = lonRange

	return &bounds, nil
}
