package okgrpcx

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
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
	attrs = append(attrs, errorLogAttrs(err)...)

	if st.Code() != codes.OK {
		logger.WarnContext(ctx, "grpc request completed", attrs...)
		return resp, err
	}

	logger.InfoContext(ctx, "grpc request completed", attrs...)
	return resp, err
}

func errorLogAttrs(err error) []any {
	if err == nil {
		return nil
	}

	attrs := []any{"err", err}

	if appAttrs := appErrorAttrs(err); len(appAttrs) > 0 {
		return append(attrs, appAttrs...)
	}

	st, ok := status.FromError(err)
	if !ok {
		return attrs
	}

	for _, detail := range st.Details() {
		info, ok := detail.(*errdetails.ErrorInfo)
		if !ok || info == nil || info.Domain != errorDomain {
			continue
		}

		if kind := info.Metadata["app_kind"]; kind != "" {
			attrs = append(attrs, "app_kind", kind)
		}
		if code := info.Metadata["app_code"]; code != "" {
			attrs = append(attrs, "app_code", code)
		}
		if op := info.Metadata["app_op"]; op != "" {
			attrs = append(attrs, "app_op", op)
		}
		break
	}

	return attrs
}

func appErrorAttrs(err error) []any {
	spec, ok := appErrorSpecOf(err)
	if !ok {
		return nil
	}

	attrs := []any{
		"app_kind", spec.kind,
		"app_code", spec.code,
	}

	return attrs
}
