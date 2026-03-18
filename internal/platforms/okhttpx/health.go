package okhttpx

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Healthz(w http.ResponseWriter, _ *http.Request) {
	Json(w, http.StatusOK, map[string]bool{"ok": true})
}

func ReadyZ(pool *pgxpool.Pool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := pool.Ping(r.Context()); err != nil {
			Fail(w, r, err)
			return
		}
		Json(w, http.StatusOK, map[string]bool{"ok": true})
	}
}

func Readyz(pool *pgxpool.Pool) http.HandlerFunc {
	return ReadyZ(pool)
}
