package geohash

import (
	"errors"

	"github.com/comfforts/geocode"

	"github.com/comfforts/geohash/pkg/talwar"
	"github.com/comfforts/geohash/pkg/veness"
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
		return veness.NewGeoCoder(), nil
	case TALWAR:
		return talwar.NewGeoCoder(), nil
	default:
		return nil, errors.New("undefined strategy")
	}
}
