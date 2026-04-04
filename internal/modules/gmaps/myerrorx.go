package gmaps

import "errors"

func WrapPlaceNotFound(cause error) error {
	if cause == nil {
		return ErrPlaceNotFound
	}
	return errors.Join(ErrPlaceNotFound, cause)
}

func WrapQueryRequired(cause error) error {
	if cause == nil {
		return ErrQueryRequired
	}
	return errors.Join(ErrQueryRequired, cause)
}

func WrapApiKeyRequired(cause error) error {
	if cause == nil {
		return ErrApiKeyRequired
	}
	return errors.Join(ErrApiKeyRequired, cause)
}

func WrapPlaceRefRequired(cause error) error {
	if cause == nil {
		return ErrPlaceRefRequired
	}
	return errors.Join(ErrPlaceRefRequired, cause)
}

func WrapCoordinatesInvalid(cause error) error {
	if cause == nil {
		return ErrCoordinatesInvalid
	}
	return errors.Join(ErrCoordinatesInvalid, cause)
}

func WrapPlaceOutOfCoverage(cause error) error {
	if cause == nil {
		return ErrPlaceOutOfCoverage
	}
	return errors.Join(ErrPlaceOutOfCoverage, cause)
}

func WrapUnavailable(cause error) error {
	if cause == nil {
		return ErrUnavailable
	}
	return errors.Join(ErrUnavailable, cause)
}

func WrapPlaceTokenInvalid(cause error) error {
	if cause == nil {
		return ErrPlaceTokenInvalid
	}
	return errors.Join(ErrPlaceTokenInvalid, cause)
}

func WrapPlaceTokenRequired(cause error) error {
	if cause == nil {
		return ErrPlaceTokenRequired
	}
	return errors.Join(ErrPlaceTokenRequired, cause)
}
