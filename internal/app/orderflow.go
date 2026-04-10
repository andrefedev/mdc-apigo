package app

import "strings"

const (
	orderStatusPending      = "pending"
	orderStatusAcepted      = "acepted"
	orderStatusCanceled     = "canceled"
	orderStatusDispatched   = "dispatched"
	orderStatusSuccessfully = "successfully"

	orderPaymentPending    = "pending"
	orderPaymentAuthorized = "authorized"
	orderPaymentRefunded   = "refunded"

	orderPaymentMethodBank   = "bank"
	orderPaymentMethodCash   = "cash"
	orderPaymentMethodBreb   = "breb"
	orderPaymentMethodQrcode = "qrcode"
)

var orderStatusTransitions = map[string]map[string]struct{}{
	orderStatusPending: {
		orderStatusAcepted:  {},
		orderStatusCanceled: {},
	},
	orderStatusAcepted: {
		orderStatusCanceled:   {},
		orderStatusDispatched: {},
	},
	orderStatusDispatched: {
		orderStatusSuccessfully: {},
	},
	orderStatusCanceled:     {},
	orderStatusSuccessfully: {},
}

func normalizeOrderStatus(value string) (string, error) {
	status := strings.TrimSpace(value)
	switch status {
	case orderStatusPending,
		orderStatusAcepted,
		orderStatusCanceled,
		orderStatusDispatched,
		orderStatusSuccessfully:
		return status, nil
	default:
		return "", WrapInvalidOrderStatus(nil)
	}
}

func normalizeOrderPaymentStatus(value string) (string, error) {
	status := strings.TrimSpace(value)
	switch status {
	case orderPaymentPending,
		orderPaymentAuthorized,
		orderPaymentRefunded:
		return status, nil
	default:
		return "", WrapInvalidOrderPaymentStatus(nil)
	}
}

func normalizeOrderPaymentMethod(value string) (string, error) {
	method := strings.TrimSpace(value)
	switch method {
	case orderPaymentMethodBank,
		orderPaymentMethodCash,
		orderPaymentMethodBreb,
		orderPaymentMethodQrcode:
		return method, nil
	default:
		return "", WrapInvalidOrderPaymentMethod(nil)
	}
}

func canTransitionOrderStatus(current, next string) bool {
	if current == next {
		return true
	}

	nexts, ok := orderStatusTransitions[current]
	if !ok {
		return false
	}

	_, ok = nexts[next]
	return ok
}

func canChangeOrderPaymentMethod(orderStatus, paymentStatus string) bool {
	if paymentStatus != orderPaymentPending {
		return false
	}

	return orderStatus == orderStatusPending || orderStatus == orderStatusAcepted
}

func canTransitionOrderPaymentStatus(orderStatus, current, next string) bool {
	if current == next {
		return true
	}

	switch next {
	case orderPaymentAuthorized:
		if current != orderPaymentPending {
			return false
		}
		return orderStatus == orderStatusPending ||
			orderStatus == orderStatusAcepted ||
			orderStatus == orderStatusDispatched

	case orderPaymentRefunded:
		if current != orderPaymentPending && current != orderPaymentAuthorized {
			return false
		}
		return orderStatus != orderStatusSuccessfully

	case orderPaymentPending:
		return current == orderPaymentPending

	default:
		return false
	}
}
