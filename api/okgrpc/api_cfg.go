package okgrpc

import (
	"context"

	"apigo/internal/app"
)

type contextKey string

const sessionCtx contextKey = "session"

func withSession(ctx context.Context, session *app.Session) context.Context {
	return context.WithValue(ctx, sessionCtx, session)
}

func sessionFromContext(ctx context.Context) (*app.Session, bool) {
	session, ok := ctx.Value(sessionCtx).(*app.Session)
	return session, ok
}
