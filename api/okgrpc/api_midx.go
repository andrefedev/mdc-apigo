package okgrpc

import (
	"context"
	"log/slog"
	"slices"
	"strings"
	"time"

	"apigo/internal/app"
	v1 "apigo/protobuf/gen/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

const (
	headerAuthorize = "authorization"
)

func bearerToken(values []string) string {
	if len(values) == 0 {
		return ""
	}

	value := strings.TrimSpace(values[0])
	if value == "" {
		return ""
	}

	return strings.TrimSpace(value)
}

func isPublicMethod(method string) bool {
	m := []string{
		v1.ApiService_Code_FullMethodName,
		v1.ApiService_CodeDetail_FullMethodName,
		v1.ApiService_CodeVerify_FullMethodName,
		"/server.reflection.v1alpha.ServerReflection/ServerReflectionInfo",
	}

	return slices.Contains(m, method)
}

// ############
// # METHOD's #
// ############

func requireLogin(ctx context.Context) (*app.Session, error) {
	session, ok := sessionFromContext(ctx)
	if !ok || session == nil {
		return nil, app.WrapSessionRequired(nil)
	}

	return session, nil
}

func requireRootUser(ctx context.Context) (*app.Session, error) {
	session, err := requireLogin(ctx)
	if err != nil {
		return nil, err
	}

	// require user...
	//if !session.CanManageUsers() {
	//	return nil, WrapForbidden(nil)
	//}

	return session, nil
}

func requireStaffUser(ctx context.Context) (*app.Session, error) {
	session, err := requireLogin(ctx)
	if err != nil {
		return nil, err
	}

	if !session.IsEmployee() {
		return nil, app.WrapForbidden(nil)
	}

	return session, nil
}

// #################
// # INTERCEPTOR's #
// #################

func SessionUnaryInterceptor(srv *Server) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return handler(ctx, req)
		}

		value := bearerToken(md.Get(headerAuthorize))
		if value == "" {
			return handler(ctx, req)
		}

		session, err := srv.useservice.SessionByIdToken(ctx, value)
		if err != nil {
			return nil, err
		}

		return handler(withSession(ctx, session), req)
	}
}

func AuthorizeUnaryInterceptor(srv *Server) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, _ := metadata.FromIncomingContext(ctx)
		idk := bearerToken(md.Get(headerAuthorize))
		if idk == "" {
			if isPublicMethod(info.FullMethod) {
				return handler(ctx, req)
			}
			return nil, app.WrapSessionRequired(nil)
		}

		session, err := srv.useservice.SessionByIdToken(ctx, idk)
		if err != nil {
			return nil, err
		}

		return handler(withSession(ctx, session), req)
	}
}

// #################
// # INTERCEPTOR's #
// #################

func UnaryErrorInterceptor(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	resp, err := handler(ctx, req)
	if err != nil {
		return nil, statusError(err)
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

	grpcCodeOf := func(err error) codes.Code {
		if err == nil {
			return codes.OK
		}
		return status.Code(statusError(err))
	}

	resp, err := handler(ctx, req)
	attrs := []any{
		"grpc_code", grpcCodeOf(err).String(),
		"duration_ms", time.Since(start).Milliseconds(),
	}
	if err != nil {
		attrs = append(attrs, "err", err)
	}

	if err != nil {
		logger.WarnContext(ctx, "grpc request completed", attrs...)
		return resp, err
	}

	logger.InfoContext(ctx, "grpc request completed", attrs...)
	return resp, err
}
