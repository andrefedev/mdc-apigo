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

var (
	ErrOrderNotFound         = errors.New("order not found")
	ErrOrderDeleteNotAllowed = errors.New("order delete not allowed")
	ErrInvalidFlatQuery      = errors.New("invalid flat query")

	ErrOrderLineNotFound          = errors.New("order line not found")
	ErrInvalidOrderLinePid        = errors.New("invalid order line pid")
	ErrInvalidOrderLineQuantity   = errors.New("invalid order line quantity")
	ErrInvalidOrderLineBasePrice  = errors.New("invalid order line base price")
	ErrInvalidOrderLineOfferPrice = errors.New("invalid order line offer price")
	ErrInvalidOrderLinePriceRange = errors.New("invalid order line price range")
)
