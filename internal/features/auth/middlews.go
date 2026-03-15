package auth

import (
	"apigo/internal/platforms/httpx"
	"log"
	"net/http"
	"strings"
)

const (
	headerAuthorize = "authorization"
)

type Middleware struct {
	Service *Service
}

func NewMiddleware(service *Service) *Middleware {
	return &Middleware{
		Service: service,
	}
}

func bearerToken(v string) string {
	const prefix = "Bearer "
	if !strings.HasPrefix(v, prefix) {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(v, prefix))
}

func (m *Middleware) AttachIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		header := r.Header.Get(headerAuthorize)
		idToken := bearerToken(header)
		if idToken == "" {
			next.ServeHTTP(w, r)
			return
		}

		log.Printf("idToken: %s", idToken)

		identity, err := m.Service.IdentityByIdToken(ctx, idToken)
		if err != nil {
			httpx.Fail(w, r, err)
			return
		}

		log.Printf("pass test")

		ctx = WithIdentity(ctx, identity)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		identity, ok := IdentityFromContext(ctx)
		if !ok || identity == nil || !identity.IsAuthenticated() {
			httpx.Fail(w, r, ErrAuthenticationRequired(nil))
			return
		}

		next.ServeHTTP(w, r)
	})
}
