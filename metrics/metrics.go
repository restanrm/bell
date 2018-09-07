package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	// HTTPRequestsCount count http requests with status code and endpoint
	HTTPRequestsCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Count http requests",
		},
		[]string{"handler", "code", "method"},
	)

	// HTTPRequestDuration represent the duration of the requests
	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Measure the duration of the requests",
			Buckets: []float64{0.01, 0.05, 0.1, 0.2, 0.4, 0.8, 1, 2, 10},
		},
		[]string{"handler", "method"},
	)
)

func init() {
	prometheus.MustRegister(
		HTTPRequestDuration,
		HTTPRequestsCount,
	)
}
