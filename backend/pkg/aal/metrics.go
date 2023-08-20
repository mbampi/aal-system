package aal

import "github.com/prometheus/client_golang/prometheus"

var (
	eventToObservationDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "event_to_observation_duration_seconds",
			Help:    "Duration of receiving an event to sending the observation in seconds",
			Buckets: prometheus.DefBuckets, // Default histogram buckets
		},
		[]string{"property"},
	)
	eventToFindingDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "event_to_finding_duration_seconds",
			Help:    "Duration of receiving an event to sending the finding in seconds",
			Buckets: prometheus.DefBuckets, // Default histogram buckets
		},
		[]string{"property"},
	)
	eventsCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "events_count",
			Help: "Number of events received",
		},
		[]string{"event"},
	)
	observationsCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "observations_count",
			Help: "Number of observations sent",
		},
		[]string{"observation"},
	)
	findingsCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "findings_count",
			Help: "Number of findings sent",
		},
		[]string{"finding"},
	)
)

func initMetrics() {
	prometheus.MustRegister(eventToObservationDuration)
	prometheus.MustRegister(eventToFindingDuration)
	prometheus.MustRegister(eventsCount)
	prometheus.MustRegister(observationsCount)
	prometheus.MustRegister(findingsCount)
}
