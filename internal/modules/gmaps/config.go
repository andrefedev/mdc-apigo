package gmaps

import "googlemaps.github.io/maps"

type Config struct {
	Language           string
	Region             string
	CountryCode        string
	SearchCenter       maps.LatLng
	SearchBounds       maps.LatLngBounds
	SearchRadiusMeters uint
	StrictBounds       bool
	AutocompleteLimit  int
}

func DefaultConfig() Config {
	return Config{
		Region:      "co",
		Language:    "es",
		CountryCode: "CO",
		SearchCenter: maps.LatLng{
			Lat: 6.244203,
			Lng: -75.581211,
		},
		SearchBounds: maps.LatLngBounds{
			SouthWest: maps.LatLng{
				Lat: 6.105700,
				Lng: -75.729200,
			},
			NorthEast: maps.LatLng{
				Lat: 6.406100,
				Lng: -75.453800,
			},
		},
		StrictBounds:       false,
		AutocompleteLimit:  5,
		SearchRadiusMeters: 35000,
	}
}
