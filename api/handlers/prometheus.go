package handlers

import (
	"fmt"
	"net/http"
	"otus/pkg/metric"
	"time"
)

func InstrumentedHandler(h http.HandlerFunc, method, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		h.ServeHTTP(rec, r)
		duration := time.Since(start).Seconds()
		metric.RequestLatency.WithLabelValues(method, endpoint).Observe(duration)
		metric.RequestCounter.WithLabelValues(method, endpoint, fmt.Sprint(rec.statusCode)).Inc()
	}
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}
