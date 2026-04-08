package app

import "errors"

func WrapInvalidCode(cause error) error {
	if cause == nil {
		return ErrInvalidCode
	}
	return errors.Join(ErrInvalidCode, cause)
}

func WrapCodeExpired(cause error) error {
	if cause == nil {
		return ErrCodeExpired
	}
	return errors.Join(ErrCodeExpired, cause)
}

func WrapCodeNotFound(cause error) error {
	if cause == nil {
		return ErrCodeNotFound
	}
	return errors.Join(ErrCodeNotFound, cause)
}

func WrapUserNotFound(cause error) error {
	if cause == nil {
		return ErrUserNotFound
	}
	return errors.Join(ErrUserNotFound, cause)
}

func WrapSessionNotFound(cause error) error {
	if cause == nil {
		return ErrSessionNotFound
	}
	return errors.Join(ErrSessionNotFound, cause)
}

func WrapSessionExpired(cause error) error {
	if cause == nil {
		return ErrSessionExpired
	}
	return errors.Join(ErrSessionExpired, cause)
}

func WrapSessionRevoked(cause error) error {
	if cause == nil {
		return ErrSessionRevoked
	}
	return errors.Join(ErrSessionRevoked, cause)
}

func WrapSessionRequired(cause error) error {
	if cause == nil {
		return ErrSessionRequired
	}
	return errors.Join(ErrSessionRequired, cause)
}

func WrapForbidden(cause error) error {
	if cause == nil {
		return ErrForbidden
	}
	return errors.Join(ErrForbidden, cause)
}

func WrapLoginRequired(cause error) error {
	if cause == nil {
		return ErrLoginRequired
	}
	return errors.Join(ErrLoginRequired, cause)
}

func WrapInvalidName(cause error) error {
	if cause == nil {
		return ErrInvalidName
	}
	return errors.Join(ErrInvalidName, cause)
}

func WrapInvalidPhone(cause error) error {
	if cause == nil {
		return ErrInvalidPhone
	}

	return errors.Join(ErrInvalidPhone, cause)
}

func WrapUserExists(cause error) error {
	if cause == nil {
		return ErrUserExists
	}
	return errors.Join(ErrUserExists, cause)
}

// USER_ADDR__

func WrapUserAddrNotFound(cause error) error {
	if cause == nil {
		return ErrUserAddrNotFound
	}
	return errors.Join(ErrUserAddrNotFound, cause)
}

func WrapInvalidFlatQuery(cause error) error {
	if cause == nil {
		return ErrInvalidFlatQuery
	}
	return errors.Join(ErrInvalidFlatQuery, cause)
}

// SALES__

func WrapOrderNotFound(cause error) error {
	if cause == nil {
		return ErrOrderNotFound
	}
	return errors.Join(ErrOrderNotFound, cause)
}

func WrapInvalidOrderLineItem(cause error) error {
	if cause == nil {
		return ErrInvalidOrderLineItem
	}
	return errors.Join(ErrInvalidOrderLineItem, cause)
}

func WrapInvalidOrderLineQuantity(cause error) error {
	if cause == nil {
		return ErrInvalidOrderLineQuantity
	}
	return errors.Join(ErrInvalidOrderLineQuantity, cause)
}

func WrapInvalidOrderLineBasePrice(cause error) error {
	if cause == nil {
		return ErrInvalidOrderLineBasePrice
	}
	return errors.Join(ErrInvalidOrderLineBasePrice, cause)
}

func WrapInvalidOrderLineOfferPrice(cause error) error {
	if cause == nil {
		return ErrInvalidOrderLineOfferPrice
	}
	return errors.Join(ErrInvalidOrderLineOfferPrice, cause)
}
