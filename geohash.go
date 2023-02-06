package geohash

import (
	"github.com/comfforts/errors"
	"github.com/comfforts/geocode"

	"github.com/comfforts/geohash/pkg/talwar"
)

type HashStrategy string

const (
	VENESS HashStrategy = "VENESS"
	TALWAR HashStrategy = "TALWAR"
)

type GeoHash interface {
	Encode(lat float64, lon float64, percision int) (string, error)
	Decode(hash string) (*geocode.RangeBounds, error)
}

func NewGeoHasher(strategy HashStrategy) (GeoHash, error) {
	switch strategy {
	case VENESS:
		return talwar.NewGeoCoder(), nil
	case TALWAR:
		return talwar.NewGeoCoder(), nil
	default:
		return nil, errors.NewAppError("undefined strategy")
	}
}
