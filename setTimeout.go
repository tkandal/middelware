package middelware

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

/*
 * Copyright (c) 2019 Norwegian University of Science and Technology
 */

// SetTimeout sets the timeout for every request
func SetTimeout(to time.Duration) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		f := func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithDeadline(r.Context(), time.Now().Add(to))
			defer cancel()
			h.ServeHTTP(w, r.Clone(ctx))
		}
		return http.HandlerFunc(f)
	}
}
