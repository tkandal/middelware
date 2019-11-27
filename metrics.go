package middelware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

/*
 * Copyright (c) 2019 Norwegian University of Science and Technology
 */

var (
	requestsDuration *prometheus.HistogramVec
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// PrometheusMetrics collect request metrics
func PrometheusMetrics(progName string) mux.MiddlewareFunc {
	requestsDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: fmt.Sprintf("%s_requests_duration_ms", progName),
			Help: fmt.Sprintf("%s requests duration ms distribution", strings.Title(progName)),
			Buckets: []float64{
				1, 2, 5, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 125, 150, 175, 200, 250, 300, 400, 500,
				600, 700, 750, 1000, 1250, 1500, 2000, 2500, 3000, 4000, 5000, 7500, 10000, 20000, 30000,
			},
		},
		[]string{"route", "method", "status"})

	prometheus.MustRegister(requestsDuration)

	return func(h http.Handler) http.Handler {
		f := func(w http.ResponseWriter, r *http.Request) {
			if isIgnore(r.URL.Path) {
				h.ServeHTTP(w, r)
				return
			}

			rw := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			defer func(s time.Time) {
				requestsDuration.With(prometheus.Labels{
					"route":  r.URL.Path,
					"method": r.Method,
					"status": strconv.FormatInt(int64(rw.statusCode), 10),
				}).Observe(float64(time.Now().Sub(s).Nanoseconds() / int64(time.Millisecond)))
			}(time.Now())

			h.ServeHTTP(rw, r)
		}
		return http.HandlerFunc(f)
	}
}
