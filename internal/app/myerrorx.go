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

// ORDER__

func WrapOrderNotFound(cause error) error {
	if cause == nil {
		return ErrOrderNotFound
	}
	return errors.Join(ErrOrderNotFound, cause)
}

func WrapOrderDeleteNotAllowed(cause error) error {
	if cause == nil {
		return ErrOrderDeleteNotAllowed
	}
	return errors.Join(ErrOrderDeleteNotAllowed, cause)
}

func WrapInvalidOrderStatus(cause error) error {
	if cause == nil {
		return ErrInvalidOrderStatus
	}
	return errors.Join(ErrInvalidOrderStatus, cause)
}

func WrapInvalidOrderPaymentStatus(cause error) error {
	if cause == nil {
		return ErrInvalidOrderPaymentStatus
	}
	return errors.Join(ErrInvalidOrderPaymentStatus, cause)
}

func WrapInvalidOrderPaymentMethod(cause error) error {
	if cause == nil {
		return ErrInvalidOrderPaymentMethod
	}
	return errors.Join(ErrInvalidOrderPaymentMethod, cause)
}

func WrapOrderInvalidTransition(cause error) error {
	if cause == nil {
		return ErrOrderInvalidTransition
	}
	return errors.Join(ErrOrderInvalidTransition, cause)
}

func WrapOrderPaymentInvalidTransition(cause error) error {
	if cause == nil {
		return ErrOrderPaymentInvalidTransition
	}
	return errors.Join(ErrOrderPaymentInvalidTransition, cause)
}

func WrapOrderLineEmpty(cause error) error {
	if cause == nil {
		return ErrOrderLineEmpty
	}
	return errors.Join(ErrOrderLineEmpty, cause)
}

// __ORDER_LINE__

func WrapOrderLineNotFound(cause error) error {
	if cause == nil {
		return ErrOrderLineNotFound
	}
	return errors.Join(ErrOrderLineNotFound, cause)
}

func WrapInvalidOrderLinePid(cause error) error {
	if cause == nil {
		return ErrInvalidOrderLinePid
	}
	return errors.Join(ErrInvalidOrderLinePid, cause)
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

func WrapInvalidOrderLinePriceRange(cause error) error {
	if cause == nil {
		return ErrInvalidOrderLinePriceRange
	}
	return errors.Join(ErrInvalidOrderLinePriceRange, cause)
}

// DELIVERY_DAY__

func WrapDeliveryDayNotFound(cause error) error {
	if cause == nil {
		return ErrDeliveryDayNotFound
	}
	return errors.Join(ErrDeliveryDayNotFound, cause)
}

func WrapInvalidDeliveryDayDate(cause error) error {
	if cause == nil {
		return ErrInvalidDeliveryDayDate
	}
	return errors.Join(ErrInvalidDeliveryDayDate, cause)
}

func WrapInvalidDeliveryDayKind(cause error) error {
	if cause == nil {
		return ErrInvalidDeliveryDayKind
	}
	return errors.Join(ErrInvalidDeliveryDayKind, cause)
}

func WrapInvalidDeliveryDayRange(cause error) error {
	if cause == nil {
		return ErrInvalidDeliveryDayRange
	}
	return errors.Join(ErrInvalidDeliveryDayRange, cause)
}

func WrapInvalidDeliveryDayCutoff(cause error) error {
	if cause == nil {
		return ErrInvalidDeliveryDayCutoff
	}
	return errors.Join(ErrInvalidDeliveryDayCutoff, cause)
}

func WrapInvalidDeliveryDayCap(cause error) error {
	if cause == nil {
		return ErrInvalidDeliveryDayCap
	}
	return errors.Join(ErrInvalidDeliveryDayCap, cause)
}
