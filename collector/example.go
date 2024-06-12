package collector

import (
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"math/rand"
)

type ExampleCollector struct {
	randomFloatMetric *prometheus.Desc
	logger            log.Logger
}

func NewExampleCollector(logger log.Logger) Collector {
	return &ExampleCollector{
		randomFloatMetric: prometheus.NewDesc("random_float_metric", "example metric", nil, nil),
		logger:            logger,
	}
}

func init() {
	registerCollector("example", NewExampleCollector)
}

func (c *ExampleCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.randomFloatMetric
}

func (c *ExampleCollector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(c.randomFloatMetric, prometheus.GaugeValue, rand.Float64())
}
