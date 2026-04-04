package gmaps

import "errors"

var (
	ErrApiKeyRequired     = errors.New("gmaps api key required")
	ErrQueryRequired      = errors.New("gmaps query required")
	ErrPlaceRefRequired   = errors.New("gmaps place id required")
	ErrCoordinatesInvalid = errors.New("gmaps invalid coordinates")
	ErrPlaceNotFound      = errors.New("gmaps place not found")
	ErrPlaceOutOfCoverage = errors.New("gmaps place outside medellin coverage")
)
