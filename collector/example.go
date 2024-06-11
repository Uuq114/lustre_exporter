package collector

import (
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"math/rand"
)

type ExampleCollector struct {
	randomFloatMetric  *prometheus.Desc
	randomStringMetric *prometheus.Desc
	logger             log.Logger
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func getRandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func NewExampleCollector(logger log.Logger) Collector {
	//return &ExampleCollector{
	//	randomFloatMetric:  rand.Float32(),
	//	randomStringMetric: getRandomString(rand.Intn(10)),
	//}
	return &ExampleCollector{
		randomFloatMetric:  prometheus.NewDesc("random_float_metric", "example metric", nil, nil),
		randomStringMetric: prometheus.NewDesc("random_string_metric", "example metric", nil, nil),
		logger:             logger,
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

//func (c *ExampleCollector) exportMetric() map[string]any {
//	rand.Seed(time.Now().UnixMilli())
//	c.randomFloatMetric = rand.Float32()
//	c.randomStringMetric = getRandomString(rand.Intn(10))
//
//	return map[string]any{
//		"RandomFloat":  c.randomFloatMetric,
//		"RandomString": c.randomStringMetric,
//	}
//}
