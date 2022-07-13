package middelware

import (
	"github.com/gorilla/mux"
	"github.com/thanhpk/randstr"
	"net/http"
)

/*
 * Copyright (c) 2022 Norwegian University of Science and Technology
 */

const (
	// SpanID is the header name for a span identity.
	SpanID = "SpanID"
	idLen  = 16
)

// SetSpanID adds a span identity as a request header.
func SetSpanID() mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		f := func(w http.ResponseWriter, r *http.Request) {
			r.Header.Set(SpanID, randstr.Hex(idLen))
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(f)
	}
}
