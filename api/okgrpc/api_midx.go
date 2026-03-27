package okgrpc

import (
	"context"
	"slices"
	"strings"

	v1 "apigo/protobuf/gen/v1"

	"apigo/internal/features/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	userContext     = "userContext"
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
		v1.AuthService_Code_FullMethodName,
		v1.AuthService_CodeDetail_FullMethodName,
		v1.AuthService_CodeVerify_FullMethodName,
		"/server.reflection.v1alpha.ServerReflection/ServerReflectionInfo",
	}

	return slices.Contains(m, method)
}

// ############
// # METHOD's #
// ############

func requireLogin(ctx context.Context) (*auth.Session, error) {
	session, ok := SessionFromContext(ctx)
	if !ok || session == nil || !session.IsAuthenticated() {
		return nil, auth.WrapSessionRequired(nil)
	}

	return session, nil
}

func requireStaff(ctx context.Context) (*auth.Session, error) {
	identity, err := RequireLogin(ctx)
	if err != nil {
		return nil, err
	}
	if !identity.CanAccessBackoffice() {
		return nil, WrapForbidden(nil)
	}

	return identity, nil
}

func requireSuperUser(ctx context.Context) (*auth.Session, error) {
	session, err := RequireLogin(ctx)
	if err != nil {
		return nil, err
	}

	// require user...

	if !session.CanManageUsers() {
		return nil, WrapForbidden(nil)
	}

	return identity, nil
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

		session, err := srv.AuthService.SessionByIdToken(ctx, value)
		if err != nil {
			return nil, err
		}

		return handler(WithSession(ctx, session), req)
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
			return nil, auth.WrapSessionRequired(nil)
		}

		session, err := srv.AuthService.SessionByIdToken(ctx, idk)
		if err != nil {
			return nil, err
		}

		return handler(WithSession(ctx, session), req)
	}
}
