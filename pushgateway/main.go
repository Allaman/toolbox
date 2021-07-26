package main

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

const (
	pushgatewayIngress = "localhost:9091"
	label              = "all"
)

var (
	counter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "iddh_pull_count",
			Help: "Docker Hub pull count",
		},
		[]string{"repo"},
	)
)

func main() {

	var pullCounterAll = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "count_all",
			Help: "Total count",
		},
		[]string{"label"},
	)
	pullCounterAll.WithLabelValues("Foo").Add(float64(1337))

	if err := push.New("localhost:9091", "count").
		Collector(pullCounterAll).
		Push(); err != nil {
		log.Fatalf("Failed to push metrics: %s", err)
	}
}
