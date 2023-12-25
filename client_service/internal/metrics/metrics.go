package metrics

import "github.com/prometheus/client_golang/prometheus"

func init() {
	prometheus.MustRegister(CreatedOrdersCounter)
	prometheus.MustRegister(ClientCanceledCounter)
}

var (
	CreatedOrdersCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "created_orders_counter",
			Help: "Сounts of created orders",
		},
	)
)

var (
	ClientCanceledCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "client_canceled_counter",
			Help: "Сounts of orders canceled by client",
		},
	)
)

var (
	LogoutCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "logout_counter",
			Help: "Сounts of logout operations",
		},
	)
)

func IncreaseSuccessfulLoginsCounter() {
	SuccessfulLoginsCounter.Inc()
}

func IncreaseUnsuccessfulLoginsCounter() {
	UnsuccessfulLoginsCounter.Inc()
}

func IncreaseLogoutCounter() {
	LogoutCounter.Inc()
}
