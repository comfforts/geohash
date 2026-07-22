# Geohash

Geohash is a small Go package for encoding latitude/longitude coordinates into
geohash strings and decoding geohash strings back into latitude/longitude
bounds.

The package exposes a simple strategy-based API:

- `VENESS`: standard Base32 geohash encoding.
- `TALWAR`: quadrant-based geohash encoding using the characters `a`, `b`, `c`,
  and `d`.

## Installation

```sh
go get github.com/comfforts/geohash
```

## Quick Start

```go
package main

import (
	"fmt"
	"log"

	"github.com/comfforts/geohash"
)

func main() {
	hasher, err := geohash.NewGeoHasher(geohash.TALWAR)
	if err != nil {
		log.Fatal(err)
	}

	hash, err := hasher.Encode(37.4438238, -122.0869895, 15)
	if err != nil {
		log.Fatal(err)
	}

	bounds, err := hasher.Decode(hash)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("hash:", hash)
	fmt.Printf("latitude:  %f to %f\n", bounds.Latitude.Min, bounds.Latitude.Max)
	fmt.Printf("longitude: %f to %f\n", bounds.Longitude.Min, bounds.Longitude.Max)
}
```

## API

### NewGeoHasher

```go
hasher, err := geohash.NewGeoHasher(geohash.TALWAR)
```

`NewGeoHasher` returns a value that satisfies the `GeoHash` interface.

Supported strategies:

| Strategy | Description | Precision |
| --- | --- | --- |
| `geohash.TALWAR` | Encodes each latitude/longitude bisection as one of `a`, `b`, `c`, or `d`. | `1` to `25`; out-of-range values default to `12`. |
| `geohash.VENESS` | Encodes coordinates with the standard geohash Base32 alphabet. | `1` to `12`; out-of-range values default to `12`. |

Unknown strategies return an error.

### Encode

```go
hash, err := hasher.Encode(latitude, longitude, precision)
```

`Encode` converts a coordinate into a geohash string.

Current validation behavior:

- Latitude must be between `-90` and `90`.
- Longitude must be between `-180` and `180`.
- Latitude and longitude must both be non-zero.
- Invalid coordinates return `constants.ErrInvalidLatLong`.

### Decode

```go
bounds, err := hasher.Decode(hash)
```

`Decode` converts a geohash string into a `*geocode.RangeBounds` value from
[`github.com/comfforts/geocode`](https://github.com/comfforts/geocode). The
decoded bounds contain minimum and maximum latitude and longitude values for the
cell represented by the hash.

Empty or invalid hash strings return `constants.ErrInvalidGeocode`.

## Development

Run the test suite:

```sh
go test ./...
```

The repository layout is intentionally small:

```text
.
├── geohash.go                # Public package API and strategy selection
├── geohash_test.go           # Package tests
├── pkg/constants/constants.go
├── pkg/talwar/geohash.go     # TALWAR strategy implementation
└── pkg/veness/geohash.go     # VENESS strategy implementation
```

## Notes

- The public method parameter is currently named `percision` in code; it
  represents geohash precision.
- `VENESS` follows the common Base32 geohash alphabet:
  `0123456789bcdefghjkmnpqrstuvwxyz`.
- `TALWAR` uses a custom alphabet and is not interchangeable with standard
  Base32 geohashes.
