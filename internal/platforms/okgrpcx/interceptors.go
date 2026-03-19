package okgrpcx

import (
	"context"

	"google.golang.org/grpc"
)

func UnaryErrorInterceptor(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	resp, err := handler(ctx, req)
	if err != nil {
		return nil, StatusError(err)
	}

	return resp, nil
}
