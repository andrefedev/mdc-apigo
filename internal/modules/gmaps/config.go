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

// 6.249581265075928, -75.58697583068572

func DefaultConfig() Config {
	return Config{
		Region:      "co",
		Language:    "es",
		CountryCode: "CO",
		SearchCenter: maps.LatLng{
			Lat: 6.2218,
			Lng: -75.5860,
		},
		SearchBounds: maps.LatLngBounds{
			SouthWest: maps.LatLng{
				Lat: 6.1450,
				Lng: -75.6200,
			},
			NorthEast: maps.LatLng{
				Lat: 6.2985,
				Lng: -75.5520,
			},
		},
		StrictBounds:       false,
		AutocompleteLimit:  10,
		SearchRadiusMeters: 35000,
	}
}
