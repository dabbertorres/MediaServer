package websrv

import (
	"context"
	"database/sql"
	"net/http"
)

type contextKey int

const contextKeyDB contextKey = 0

func MiddlewareDB(db *sql.DB, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), contextKeyDB, db)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
