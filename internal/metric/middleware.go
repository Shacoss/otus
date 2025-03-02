package metric

import (
	"fmt"
	"net/http"
	"time"
)

func HttpMetricMiddleware(h http.HandlerFunc, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		h.ServeHTTP(rec, r)
		duration := time.Since(start).Seconds()
		RequestLatency.WithLabelValues(r.Method, endpoint).Observe(duration)
		RequestCounter.WithLabelValues(r.Method, endpoint, fmt.Sprint(rec.statusCode)).Inc()
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
