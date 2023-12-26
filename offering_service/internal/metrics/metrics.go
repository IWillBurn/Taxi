package metrics

import "github.com/prometheus/client_golang/prometheus"

func init() {
	prometheus.MustRegister(OffersCounter)
	prometheus.MustRegister(DecodingRequestsCounter)
}

var (
	OffersCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "count_of_offers",
			Help: "Counts of created offers",
		},
	)
)

var (
	DecodingRequestsCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "decoding_requests_counter",
			Help: "Counts of unsuccessful logins",
		},
	)
)

func IncreaseOffersCounter() {
	OffersCounter.Inc()
}

func IncreaseDecodingRequestsCounter() {
	DecodingRequestsCounter.Inc()
}
