package okgrpc

import (
	"apigo/internal/features/auth"
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func bearerToken(values []string) string {
	if len(values) == 0 {
		return ""
	}

	const prefix = "Bearer "
	value := strings.TrimSpace(values[0])
	if value == "" || !strings.HasPrefix(value, prefix) {
		return ""
	}

	return strings.TrimSpace(strings.TrimPrefix(value, prefix))
}

func (m *Server) isPublicMethod(method string) bool {
	_, ok := m.publicMethods[method]
	return ok
}

func (m *Server) AttachIdentityUnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	idToken := bearerToken(md.Get(headerAuthorize))
	if idToken == "" {
		return handler(ctx, req)
	}

	identity, err := m.Service.IdentityByIdToken(ctx, idToken)
	if err != nil {
		return nil, err
	}

	return handler(WithIdentity(ctx, identity), req)
}

func (m *Server) AuthorizeUnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	idToken := bearerToken(md.Get(headerAuthorize))
	if idToken == "" {
		if m.isPublicMethod(info.FullMethod) {
			return handler(ctx, req)
		}
		return nil, WrapSessionRequired(nil)
	}

	identity, err := m.Service.IdentityByIdToken(ctx, idToken)
	if err != nil {
		return nil, err
	}

	return handler(WithIdentity(ctx, identity), req)
}

// ############
// # METHOD's #
// ############

func RequireLogin(ctx context.Context) (*Identity, error) {
	identity, ok := IdentityFromContext(ctx)
	if !ok || identity == nil || !identity.IsAuthenticated() {
		return nil, auth.WrapSessionRequired(nil)
	}

	return identity, nil
}

func RequireStaff(ctx context.Context) (*Identity, error) {
	identity, err := RequireLogin(ctx)
	if err != nil {
		return nil, err
	}
	if !identity.CanAccessBackoffice() {
		return nil, WrapForbidden(nil)
	}

	return identity, nil
}

func RequireSuperUser(ctx context.Context) (*Identity, error) {
	session, err := RequireLogin(ctx)
	if err != nil {
		return nil, err
	}
	if !session.CanManageUsers() {
		return nil, WrapForbidden(nil)
	}

	return identity, nil
}
