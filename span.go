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
	SpanID = "SpanID"
	idLen  = 16
)

func SetSpanID() mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		f := func(w http.ResponseWriter, r *http.Request) {
			r.Header.Set(SpanID, randstr.Hex(idLen))
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(f)
	}
}
