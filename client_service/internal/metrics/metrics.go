package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type Metrics struct {
	CreatedOrdersCounter prometheus.Counter
	CanceledTripCounter  prometheus.Counter
	EndedTripCounter     prometheus.Counter
	InTheQueueCounter    prometheus.Gauge
}

func NewMetrics() *Metrics {
	return &Metrics{
		CreatedOrdersCounter: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: "client", Name: "created_order_counter", Help: "Counts the number of created orders",
		}),
		CanceledTripCounter: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: "client", Name: "cancel_trip_counter", Help: "Counts the number of cancelled trips",
		}),
		EndedTripCounter: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: "client", Name: "ended_trip_counter", Help: "Counts the number of ended trips",
		}),
		InTheQueueCounter: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: "client", Name: "number_waiting_users", Help: "Counts the number of users in queue",
		}),
	}
}

func (metrics *Metrics) Serve() {
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":9000", nil)
}
