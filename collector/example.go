package collector

import (
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"math/rand"
)

type ExampleCollector struct {
	randomFloatMetric1 *prometheus.Desc
	randomFloatMetric2 *prometheus.Desc
	logger             log.Logger
}

func NewExampleCollector(logger log.Logger) Collector {
	return &ExampleCollector{
		randomFloatMetric1: prometheus.NewDesc("random_float_metric1", "example metric", nil, nil),
		randomFloatMetric2: prometheus.NewDesc("random_float_metric2", "example metric", nil, nil),
		logger:             logger,
	}
}

func init() {
	registerCollector("example", NewExampleCollector)
}

func (c *ExampleCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.randomFloatMetric1
	ch <- c.randomFloatMetric2
}

func (c *ExampleCollector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(c.randomFloatMetric1, prometheus.GaugeValue, rand.Float64())
	ch <- prometheus.MustNewConstMetric(c.randomFloatMetric2, prometheus.GaugeValue, rand.Float64())
}
