package okgrpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func statusError(err error) error {
	if err == nil {
		return nil
	}

	if _, ok := status.FromError(err); ok {
		return err
	}

	if st := grpcStatusError(err); st != nil {
		return st
	}

	return status.Error(codes.Internal, "Ha ocurrido un error interno del servidor")
}
