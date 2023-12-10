package veness

import (
	"strings"

	"github.com/comfforts/geocode"
	"github.com/comfforts/geohash/pkg/constants"
)

type geoCoder struct {
	seed string
}

func NewGeoCoder() *geoCoder {
	return &geoCoder{
		seed: "0123456789bcdefghjkmnpqrstuvwxyz", // (geohash-specific) Base32 map
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

	hash := ""
	evBit := true
	idx, bit := 0, 0
	latMin, latMax := constants.LAT_MIN, constants.LAT_MAX
	lonMin, lonMax := constants.LON_MIN, constants.LON_MAX

	for len(hash) < percision {
		if evBit {
			// bisect E-W longitude
			lonMid := (lonMin + lonMax) / 2
			if lon >= lonMid {
				idx = idx*2 + 1
				lonMin = lonMid
			} else {
				idx = idx * 2
				lonMax = lonMid
			}
		} else {
			// bisect N-S latitude
			latMid := (latMin + latMax) / 2
			if lat >= latMid {
				idx = idx*2 + 1
				latMin = latMid
			} else {
				idx = idx * 2
				latMax = latMid
			}
		}
		evBit = !evBit

		bit++

		if bit > 4 {
			// 5 bits gives us a character: append it and start over
			hash += string(gc.seed[idx])
			bit = 0
			idx = 0
		}
	}

	return hash, nil
}

func (gc *geoCoder) Decode(hash string) (*geocode.RangeBounds, error) {
	bounds, err := gc.bounds(hash)
	if err != nil {
		return nil, err
	}

	return bounds, nil
}

func (gc *geoCoder) bounds(hash string) (*geocode.RangeBounds, error) {
	if len(hash) < 1 {
		return nil, constants.ErrInvalidGeocode
	}

	hash = strings.ToLower(hash)
	evBit := true
	latMin, latMax := constants.LAT_MIN, constants.LAT_MAX
	lonMin, lonMax := constants.LON_MIN, constants.LON_MAX

	for i := 0; i < len(hash); i++ {
		ch := string(hash[i])
		idx := strings.Index(gc.seed, ch)
		if idx < 0 {
			return nil, constants.ErrInvalidGeocode
		}

		for n := 4; n > 0; n-- {
			bitN := idx >> n & 1
			if evBit {
				// longitude
				lonMid := (lonMin + lonMax) / 2
				if bitN == 1 {
					lonMin = lonMid
				} else {
					lonMax = lonMid
				}
			} else {
				// latitude
				latMid := (latMin + latMax) / 2
				if bitN == 1 {
					latMin = latMid
				} else {
					latMax = latMid
				}
			}
			evBit = !evBit
		}
	}

	return &geocode.RangeBounds{
		Latitude: geocode.Range{
			Min: latMin,
			Max: latMax,
		},
		Longitude: geocode.Range{
			Min: lonMin,
			Max: lonMax,
		},
	}, nil
}

// /**
//  * Determines adjacent cell in given direction.
//  *
//  * @param   geohash - Cell to which adjacent cell is required.
//  * @param   direction - Direction from geohash (N/S/E/W).
//  * @returns {string} Geocode of adjacent cell.
//  * @throws  Invalid geohash.
//  */
//  static adjacent(geohash, direction) {
//     // based on github.com/davetroy/geohash-js

//     geohash = geohash.toLowerCase();
//     direction = direction.toLowerCase();

//     if (geohash.length == 0) throw new Error('Invalid geohash');
//     if ('nsew'.indexOf(direction) == -1) throw new Error('Invalid direction');

//     const neighbour = {
//         n: [ 'p0r21436x8zb9dcf5h7kjnmqesgutwvy', 'bc01fg45238967deuvhjyznpkmstqrwx' ],
//         s: [ '14365h7k9dcfesgujnmqp0r2twvyx8zb', '238967debc01fg45kmstqrwxuvhjyznp' ],
//         e: [ 'bc01fg45238967deuvhjyznpkmstqrwx', 'p0r21436x8zb9dcf5h7kjnmqesgutwvy' ],
//         w: [ '238967debc01fg45kmstqrwxuvhjyznp', '14365h7k9dcfesgujnmqp0r2twvyx8zb' ],
//     };
//     const border = {
//         n: [ 'prxz',     'bcfguvyz' ],
//         s: [ '028b',     '0145hjnp' ],
//         e: [ 'bcfguvyz', 'prxz'     ],
//         w: [ '0145hjnp', '028b'     ],
//     };

//     const lastCh = geohash.slice(-1);    // last character of hash
//     let parent = geohash.slice(0, -1); // hash without last character

//     const type = geohash.length % 2;

//     // check for edge-cases which don't share common prefix
//     if (border[direction][type].indexOf(lastCh) != -1 && parent != '') {
//         parent = Geohash.adjacent(parent, direction);
//     }

//     // append letter for direction to parent
//     return parent + base32.charAt(neighbour[direction][type].indexOf(lastCh));
// }
// /**
//      * Returns all 8 adjacent cells to specified geohash.
//      *
//      * @param   {string} geohash - Geohash neighbours are required of.
//      * @returns {{n,ne,e,se,s,sw,w,nw: string}}
//      * @throws  Invalid geohash.
//      */
// 	 static neighbours(geohash) {
//         return {
//             'n':  Geohash.adjacent(geohash, 'n'),
//             'ne': Geohash.adjacent(Geohash.adjacent(geohash, 'n'), 'e'),
//             'e':  Geohash.adjacent(geohash, 'e'),
//             'se': Geohash.adjacent(Geohash.adjacent(geohash, 's'), 'e'),
//             's':  Geohash.adjacent(geohash, 's'),
//             'sw': Geohash.adjacent(Geohash.adjacent(geohash, 's'), 'w'),
//             'w':  Geohash.adjacent(geohash, 'w'),
//             'nw': Geohash.adjacent(Geohash.adjacent(geohash, 'n'), 'w'),
//         };
//     }
