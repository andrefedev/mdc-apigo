package gmaps

import "errors"

var (
	ErrPlaceNotFound      = errors.New("gmaps place not found")
	ErrQueryRequired      = errors.New("gmaps query required")
	ErrApiKeyRequired     = errors.New("gmaps api key required")
	ErrPlaceRefRequired   = errors.New("gmaps place id required")
	ErrCoordinatesInvalid = errors.New("gmaps invalid coordinates")
	ErrPlaceOutOfCoverage = errors.New("gmaps place outside medellin coverage")
)

var (
	ErrUnavailable        = errors.New("gmaps unavailable")
	ErrPlaceTokenInvalid  = errors.New("gmaps place token invalid")
	ErrPlaceTokenRequired = errors.New("gmaps place token required")
)
