package middelware

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

/*
 * Copyright (c) 2019 Norwegian University of Science and Technology
 */

var (
	ignore = []string{"/healthz", "/metrics", "/swagger"}
)

func isIgnore(path string) bool {
	for _, p := range ignore {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}

// LogEntry logs every request
func LogEntry(logger *zap.SugaredLogger, cutPath bool) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		f := func(w http.ResponseWriter, r *http.Request) {
			if isIgnore(r.URL.Path) {
				h.ServeHTTP(w, r)
				return
			}

			path := r.URL.Path
			if cutPath {
				if idx := strings.LastIndexByte(r.URL.Path, '/'); idx > 0 {
					path = r.URL.Path[:idx]
				}
			}
			spanId := r.Header.Get(SpanID)
			ra := remoteAddr(r, logger)
			logger.Infow(fmt.Sprintf("%s %s %s", ra, r.Method, path), "span", spanId)
			defer func(s time.Time) {
				logger.Infow(fmt.Sprintf("%s %s %s took %s", ra, r.Method, path, time.Since(s)), "span", spanId)
			}(time.Now())

			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(f)
	}
}

func remoteAddr(r *http.Request, logger *zap.SugaredLogger) string {
	adr, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		logger.Errorw("split host and port failed", "error", err)
		return ""
	}
	return adr
}
