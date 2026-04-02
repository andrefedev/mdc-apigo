package okgrpc

import (
	"context"

	"apigo/internal/features/auth"
)

type contextKey string

const sessionCtx contextKey = "session"

func withSession(ctx context.Context, session *auth.Session) context.Context {
	return context.WithValue(ctx, sessionCtx, session)
}

func sessionFromContext(ctx context.Context) (*auth.Session, bool) {
	session, ok := ctx.Value(sessionCtx).(*auth.Session)
	return session, ok
}
