package okgrpcx

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func UnaryErrorInterceptor(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	resp, err := handler(ctx, req)
	if err != nil {
		return nil, StatusError(err)
	}

	return resp, nil
}

func UnaryLoggingInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()

	logger := slog.With(
		"transport", "grpc",
		"full_method", info.FullMethod,
	)

	if p, ok := peer.FromContext(ctx); ok && p.Addr != nil {
		logger = logger.With("peer_addr", p.Addr.String())
	}

	resp, err := handler(ctx, req)

	st := status.Convert(err)
	attrs := []any{
		"grpc_code", st.Code().String(),
		"duration_ms", time.Since(start).Milliseconds(),
	}

	if st.Code() != codes.OK {
		logger.WarnContext(ctx, "grpc request completed", attrs...)
		return resp, err
	}

	logger.InfoContext(ctx, "grpc request completed", attrs...)
	return resp, err
}
