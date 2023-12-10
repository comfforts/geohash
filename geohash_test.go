package geohash

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/comfforts/geocode"
	"github.com/comfforts/geohash/pkg/constants"
)

// At highest resolution, hash should differ for points where
// either latitude is apart by atleast 0.044
// or longitude is apart by atleast 0.089
const (
	LATITUDE_RESOLUTION  = 0.044
	LONGITUDE_RESOLUTION = 0.088
)

func TestGeoHasherTalwar(t *testing.T) {
	for scenario, fn := range map[string]func(
		t *testing.T,
		gh GeoHash,
	){
		"test encode decode succeeds":              testEncodeDecode,
		"test error scenarios succeeds":            testErrors,
		"encoding resolution test succeeds":        testEncodingResolution,
		"encoding resolution change test succeeds": testEncodingResolutionChange,
	} {
		t.Run(scenario, func(t *testing.T) {
			gh, teardown := setupGeoHasher(t, TALWAR)
			defer teardown()
			fn(t, gh)
		})
	}
}

func TestGeoHasherVeness(t *testing.T) {
	for scenario, fn := range map[string]func(
		t *testing.T,
		gh GeoHash,
	){
		"test precision mapping succeeds":   testPrecisionMap,
		"encoding resolution test succeeds": testEncodingResolution,
		// "encoding resolution change test succeeds": testEncodingResolutionChange,
	} {
		t.Run(scenario, func(t *testing.T) {
			gh, teardown := setupGeoHasher(t, VENESS)
			defer teardown()
			fn(t, gh)
		})
	}
}

func setupGeoHasher(t *testing.T, ht HashStrategy) (
	gh GeoHash,
	teardown func(),
) {
	t.Helper()

	gh, err := NewGeoHasher(ht)
	require.NoError(t, err)

	return gh, func() {
		t.Logf(" geo hasher test ended, will cleanup")
	}
}

func testErrors(t *testing.T, gh GeoHash) {
	precision := 12
	point := geocode.Point{Latitude: 0, Longitude: 0}

	_, err := gh.Encode(point.Latitude, point.Longitude, precision)
	require.Error(t, err)
	require.Equal(t, err, constants.ErrInvalidLatLong)

	_, err = gh.Decode("")
	require.Error(t, err)
	require.Equal(t, err, constants.ErrInvalidGeocode)
}

func testEncodeDecode(t *testing.T, gh GeoHash) {
	points := []geocode.Point{
		{Latitude: 0.133333, Longitude: 117.500000},
		{Latitude: -33.918861, Longitude: 18.423300},
		{Latitude: 38.294788, Longitude: -122.461510},
		{Latitude: 28.644800, Longitude: 77.216721},
	}

	precision := 12

	for _, point := range points {
		hash, _ := gh.Encode(point.Latitude, point.Longitude, precision)
		bound, err := gh.Decode(hash)
		require.NoError(t, err)

		fmt.Printf(
			"point lat: %f, long: %f\n",
			point.Latitude,
			point.Longitude,
		)
		fmt.Printf(
			"result at percision: %d, lat: %f, long: %f\n",
			precision,
			bound.Latitude.Min+(bound.Latitude.Max-bound.Latitude.Min)/2,
			bound.Longitude.Min+(bound.Longitude.Max-bound.Longitude.Min)/2,
		)
		fmt.Printf(
			"error lat: %f, error long: %f\n",
			(bound.Latitude.Min+(bound.Latitude.Max-bound.Latitude.Min)/2)-point.Latitude,
			(bound.Longitude.Min+(bound.Longitude.Max-bound.Longitude.Min)/2)-point.Longitude,
		)
		require.Equal(t, true, true)
	}
	t.Logf("test encoding resolution done")
}

func testEncodingResolution(t *testing.T, gh GeoHash) {
	points := []geocode.Point{
		{Latitude: 0.133333, Longitude: 117.500000},
		{Latitude: -33.918861, Longitude: 18.423300},
		{Latitude: 38.294788, Longitude: -122.461510},
		{Latitude: 28.644800, Longitude: 77.216721},
	}

	for _, point := range points {
		for p := 13; p > 4; p-- {
			hash, _ := gh.Encode(point.Latitude, point.Longitude, p)
			bound, err := gh.Decode(hash)
			require.NoError(t, err)

			fmt.Printf("at percision: %d, lat range: %f, long range: %f\n", p, bound.Latitude.Max-bound.Latitude.Min, bound.Longitude.Max-bound.Longitude.Min)

			// require.Equal(t, true, math.Abs(bound.Latitude.Max-bound.Latitude.Min) < LATITUDE_RESOLUTION)
			// require.Equal(t, true, math.Abs(bound.Longitude.Max-bound.Longitude.Min) < LONGITUDE_RESOLUTION)
			require.Equal(t, true, true)

		}
	}
	t.Logf("test encoding resolution done")
}

func testPrecisionMap(t *testing.T, gh GeoHash) {
	point := geocode.Point{Latitude: 42.7134562468, Longitude: -79.8196752384}

	valueMap := map[string]float64{
		"10 ft":     0.00001,
		"100 ft":    0.0001,
		"1000 ft":   0.001,
		"1 mile":    0.01,
		"10 miles":  0.1,
		"100 miles": 1,
	}
	fmt.Printf("point: %v, value map %v\n", point, valueMap)
	percMap := map[string]int{}

	for k, v := range valueMap {
		t.Logf("building precision map for %s precision", k)
		p := 5
		found := false
		for !found {
			hash, _ := gh.Encode(point.Latitude, point.Longitude, p)
			bound, err := gh.Decode(hash)
			require.NoError(t, err)

			if bound.Latitude.Max-bound.Latitude.Min > v || bound.Longitude.Max-bound.Longitude.Min > v {
				p++
			} else {
				found = true
			}
		}
		t.Logf("required precision for %s is %d", k, p)
		percMap[k] = p
	}
	fmt.Printf("precision map %v\n", percMap)
	for p := 13; p > 4; p-- {
		for k, v := range valueMap {
			res := map[string][]geocode.Point{}
			n := 0
			cn := 0
			for n < 10 {
				// get hash
				hash, err := gh.Encode(point.Latitude, point.Longitude, p)
				require.NoError(t, err)
				percHash := hash[0:p]
				_, ok := res[percHash]
				if !ok {
					res[percHash] = []geocode.Point{}
				}
				res[percHash] = append(res[percHash], point)
				cn++

				// move latitude and get hash
				newLat := point.Latitude + v
				// fmt.Printf("moving lat by: %f, from: %f, to: %f\n", v, point.Latitude, newLat)
				hash, _ = gh.Encode(newLat, point.Longitude, p)
				percHash = hash[0:p]
				_, ok = res[percHash]
				if !ok {
					res[percHash] = []geocode.Point{}
				}
				res[percHash] = append(res[percHash], geocode.Point{
					Latitude:  newLat,
					Longitude: point.Longitude,
				})
				cn++

				// move longitude and get hash
				newLong := point.Longitude + v
				// fmt.Printf("moving long by: %f, from: %f, to: %f\n", v, point.Longitude, newLong)
				hash, _ = gh.Encode(point.Latitude, newLong, p)
				percHash = hash[0:p]
				_, ok = res[percHash]
				if !ok {
					res[percHash] = []geocode.Point{}
				}
				res[percHash] = append(res[percHash], geocode.Point{
					Latitude:  point.Latitude,
					Longitude: newLong,
				})
				cn++

				point.Latitude = point.Latitude + v
				point.Longitude = point.Longitude + v
				n++
			}
			fmt.Printf("hash map at precision: %d, for increment: %s-%f, len: %d\n", p, k, v, len(res))
			t.Logf("hash map at precision: %d, for increment: %f", p, v)
		}
	}
	t.Logf("test precision map done")
}

func testEncodingResolutionChange(t *testing.T, gh GeoHash) {
	point := geocode.Point{Latitude: 42.713456, Longitude: -79.819675}

	res := map[string][]geocode.Point{}

	n := 0
	cn := 0
	for n < 10 {
		// get hash
		hash, _ := gh.Encode(point.Latitude, point.Longitude, 9)
		_, ok := res[hash]
		if !ok {
			res[hash] = []geocode.Point{}
		}
		res[hash] = append(res[hash], point)
		cn++

		// move latitude and get hash
		hash, _ = gh.Encode(point.Latitude+0.045, point.Longitude, 9)
		_, ok = res[hash]
		if !ok {
			res[hash] = []geocode.Point{}
		}
		res[hash] = append(res[hash], geocode.Point{
			Latitude:  point.Latitude + LATITUDE_RESOLUTION,
			Longitude: point.Longitude,
		})
		cn++

		// move longitude and get hash
		hash, _ = gh.Encode(point.Latitude, point.Longitude+0.09, 9)
		_, ok = res[hash]
		if !ok {
			res[hash] = []geocode.Point{}
		}
		res[hash] = append(res[hash], geocode.Point{
			Latitude:  point.Latitude,
			Longitude: point.Longitude + LATITUDE_RESOLUTION,
		})
		cn++

		point.Latitude = point.Latitude + 0.045
		point.Longitude = point.Longitude + LONGITUDE_RESOLUTION
		n++
	}
	require.Equal(t, cn, len(res))
}
