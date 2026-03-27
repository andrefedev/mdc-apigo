package okgrpc

import (
	"apigo/internal/features/auth"
	"context"
)

type contextKey string

const sessionCtx contextKey = "session"

func WithSession(ctx context.Context, session *auth.Session) context.Context {
	return context.WithValue(ctx, sessionCtx, session)
}

func SessionFromContext(ctx context.Context) (*auth.Session, bool) {
	session, ok := ctx.Value(sessionCtx).(*auth.Session)
	return session, ok
}
