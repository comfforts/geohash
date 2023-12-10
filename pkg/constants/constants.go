package constants

import "github.com/comfforts/errors"

const (
	ERROR_DECODING_BOUNDS string = "error decoding bounds"
	INVALID_LAT_LONG      string = "error: invalid latitude or longitude"
	INVALID_HASH          string = "error: invalid geocode hash string"
)

const (
	LAT_MIN float64 = -90
	LAT_MAX float64 = 90
	LON_MIN float64 = -180
	LON_MAX float64 = 180
)

var (
	ErrInvalidLatLong = errors.NewAppError(INVALID_LAT_LONG)
	ErrDecodingBounds = errors.NewAppError(ERROR_DECODING_BOUNDS)
	ErrInvalidGeocode = errors.NewAppError(INVALID_HASH)
)
