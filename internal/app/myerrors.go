package app

import "errors"

// AUTH___

var (
	ErrForbidden       = errors.New("auth forbidden")
	ErrInvalidCode     = errors.New("auth invalid code")
	ErrCodeExpired     = errors.New("auth code expired")
	ErrCodeNotFound    = errors.New("auth code not found")
	ErrSessionNotFound = errors.New("auth session not found")
	ErrSessionExpired  = errors.New("auth session expired")
	ErrSessionRevoked  = errors.New("auth session revoked")
	ErrSessionRequired = errors.New("auth session required")
)

// USER___

var (
	ErrUserExists    = errors.New("user already exists")
	ErrUserNotFound  = errors.New("user not found")
	ErrLoginRequired = errors.New("login required")

	ErrInvalidName      = errors.New("user invalid name")
	ErrInvalidPhone     = errors.New("user invalid phone")
	ErrInvalidLastLogin = errors.New("user invalid last login")
)

// USER_ADDR__

var (
	ErrUserAddrNotFound = errors.New("user address not found")
)

var (
	ErrInvalidMaskPath = errors.New("invalid mask path")
)

// GOOGLE_MAPS__

var (
	ErrMapxUnavailable        = errors.New("mapx unavailable")
	ErrMapxQueryRequired      = errors.New("mpax query required")
	ErrMapxPlaceRefRequired   = errors.New("mpax place ref required")
	ErrMapxPlaceTokenInvalid  = errors.New("mpax place token invalid")
	ErrMapxPlaceTokenRequired = errors.New("mpax place token required")
	ErrMapxCoordinatesInvalid = errors.New("mapx invalid coordinates")
)
