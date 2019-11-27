package middelware

import "net/http"

/*
 * Copyright (c) 2019 Norwegian University of Science and Technology
 */

func NoCacheAndSniff(h http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "private, no-cache, no-store, must-revalidate, max-age=0")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(f)
}
