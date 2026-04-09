package okgrpc

import (
	"apigo/internal/app"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func grpcStatusError(err error) error {
	// AUTH__

	if errors.Is(err, app.ErrInvalidCode) {
		return status.Error(codes.InvalidArgument, "El código ingresado no es válido")
	}
	if errors.Is(err, app.ErrCodeExpired) {
		return status.Error(codes.FailedPrecondition, "El código ingresado ya expiró")
	}
	if errors.Is(err, app.ErrCodeNotFound) {
		return status.Error(codes.NotFound, "El código solicitado no existe")
	}

	// USER__

	if errors.Is(err, app.ErrUserExists) {
		return status.Error(codes.AlreadyExists, "Ya existe un usuario con los datos suministrados")
	}
	if errors.Is(err, app.ErrLoginRequired) {
		return status.Error(codes.Unauthenticated, "Debes iniciar sesión para continuar")
	}
	if errors.Is(err, app.ErrInvalidName) {
		return status.Error(codes.InvalidArgument, "El nombre no es válido")
	}
	if errors.Is(err, app.ErrInvalidPhone) {
		return status.Error(codes.InvalidArgument, "El número de teléfono no es válido")
	}
	if errors.Is(err, app.ErrInvalidLastLogin) {
		return status.Error(codes.InvalidArgument, "La fecha de último ingreso no es válida")
	}

	if errors.Is(err, app.ErrUserNotFound) {
		return status.Error(codes.NotFound, "Usuario no encontrado")
	}

	// SESSION__

	if errors.Is(err, app.ErrSessionNotFound) {
		return status.Error(codes.Unauthenticated, "Debes iniciar sesión para continuar")
	}
	if errors.Is(err, app.ErrSessionRequired) {
		return status.Error(codes.Unauthenticated, "Debes iniciar sesión para continuar")
	}
	if errors.Is(err, app.ErrSessionRevoked) {
		return status.Error(codes.Unauthenticated, "Debes iniciar sesión para continuar")
	}
	if errors.Is(err, app.ErrSessionExpired) {
		return status.Error(codes.Unauthenticated, "Debes iniciar sesión para continuar")
	}
	if errors.Is(err, app.ErrForbidden) {
		return status.Error(codes.PermissionDenied, "No tienes permisos para realizar esta acción")
	}

	// USER_ADDR__

	if errors.Is(err, app.ErrUserAddrNotFound) {
		return status.Error(codes.NotFound, "La dirección del usuario no existe")
	}

	// ORDER__

	if errors.Is(err, app.ErrOrderNotFound) {
		return status.Error(codes.NotFound, "El pedido solicitado no existe")
	}
	if errors.Is(err, app.ErrOrderDeleteNotAllowed) {
		return status.Error(codes.FailedPrecondition, "El pedido no se puede eliminar en su estado actual")
	}
	if errors.Is(err, app.ErrInvalidOrderStatus) {
		return status.Error(codes.InvalidArgument, "El estado del pedido no es válido")
	}
	if errors.Is(err, app.ErrOrderInvalidTransition) {
		return status.Error(codes.FailedPrecondition, "La transición de estado del pedido no es válida")
	}
	if errors.Is(err, app.ErrOrderLineEmpty) {
		return status.Error(codes.FailedPrecondition, "El pedido no tiene líneas para cambiar a ese estado")
	}
	if errors.Is(err, app.ErrInvalidFlatQuery) {
		return status.Error(codes.InvalidArgument, "El filtro de búsqueda no es válido")
	}

	// ORDER_LINE__

	if errors.Is(err, app.ErrOrderLineNotFound) {
		return status.Error(codes.NotFound, "La línea del pedido no existe")
	}
	if errors.Is(err, app.ErrInvalidOrderLinePid) {
		return status.Error(codes.InvalidArgument, "La referencia del producto de la línea no es válida")
	}
	if errors.Is(err, app.ErrInvalidOrderLineQuantity) {
		return status.Error(codes.InvalidArgument, "La cantidad de la línea no es válida")
	}
	if errors.Is(err, app.ErrInvalidOrderLineBasePrice) {
		return status.Error(codes.InvalidArgument, "El precio base de la línea no es válido")
	}
	if errors.Is(err, app.ErrInvalidOrderLineOfferPrice) {
		return status.Error(codes.InvalidArgument, "El precio de oferta de la línea no es válido")
	}
	if errors.Is(err, app.ErrInvalidOrderLinePriceRange) {
		return status.Error(codes.InvalidArgument, "El rango de precios de la línea no es válido")
	}

	// DELIVERY_DAY__

	if errors.Is(err, app.ErrDeliveryDayNotFound) {
		return status.Error(codes.NotFound, "El día de entrega solicitado no existe")
	}
	if errors.Is(err, app.ErrInvalidDeliveryDayDate) {
		return status.Error(codes.InvalidArgument, "La fecha del día de entrega no es válida")
	}
	if errors.Is(err, app.ErrInvalidDeliveryDayKind) {
		return status.Error(codes.InvalidArgument, "El tipo del día de entrega no es válido")
	}
	if errors.Is(err, app.ErrInvalidDeliveryDayRange) {
		return status.Error(codes.InvalidArgument, "El rango horario del día de entrega no es válido")
	}
	if errors.Is(err, app.ErrInvalidDeliveryDayCutoff) {
		return status.Error(codes.InvalidArgument, "El cutoff del día de entrega no es válido")
	}
	if errors.Is(err, app.ErrInvalidDeliveryDayCap) {
		return status.Error(codes.InvalidArgument, "La capacidad del día de entrega no es válida")
	}

	// INPUT__

	if errors.Is(err, app.ErrInvalidMaskPath) {
		return status.Error(codes.InvalidArgument, "La máscara de actualización contiene rutas no válidas")
	}

	return nil
}
