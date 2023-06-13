package main

import (
	"net/http"
	"testing"
	"time"

	prom "github.com/bnb-chain/zkbnb-smt/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var Monitor *prom.Collector

func init() {
	Monitor = prom.NewCollector()
}

func TestResource(t *testing.T) {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			panic(err)
		}
	}()
	for _, env := range prepareEnv() {
		smt, _ := initSMTWithMetrics(env, Monitor)
		// set operations
		opts, _ := generateSetOperations(12000)
		opts = append(opts, TestOperation{
			method: Commit,
		})
		execOperationsWithSleep(smt, opts, time.Duration(10*time.Microsecond))
		time.Sleep(1200 * time.Second)
	}
}
