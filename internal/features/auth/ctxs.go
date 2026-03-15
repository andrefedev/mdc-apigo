package auth

import "context"

type identityKey struct{}

func WithIdentity(ctx context.Context, v *Identity) context.Context {
	return context.WithValue(ctx, identityKey{}, v)
}

func IdentityFromContext(ctx context.Context) (*Identity, bool) {
	v, ok := ctx.Value(identityKey{}).(*Identity)
	return v, ok
}
