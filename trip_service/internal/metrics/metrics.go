package metrics

import "github.com/prometheus/client_golang/prometheus"

func init() {
	prometheus.MustRegister(SuccessfulLoginsCounter)
	prometheus.MustRegister(UnsuccessfulLoginsCounter)
	prometheus.MustRegister(LogoutCounter)
}

var (
	SuccessfulLoginsCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "successful_logins_counter",
			Help: "Сounts of successful logins",
		},
	)
)

var (
	UnsuccessfulLoginsCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "unsuccessful_logins_counter",
			Help: "Сounts of unsuccessful logins",
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
