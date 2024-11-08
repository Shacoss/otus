package metric

import "github.com/prometheus/client_golang/prometheus"

var (
	RequestLatency = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:       "http_request_latency",
		Help:       "Latency of HTTP requests.",
		Objectives: map[float64]float64{0.5: 0.05, 0.95: 0.005, 0.99: 0.001},
	}, []string{"method", "path"})

	RequestCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of API requests",
	}, []string{"method", "endpoint", "status"})
)

func RegisterMetrics() {
	prometheus.MustRegister(RequestLatency)
	prometheus.MustRegister(RequestCounter)
}
