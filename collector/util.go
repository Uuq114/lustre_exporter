package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

/*
注册collector需要实现Collector接口，要求要实现两个方法：
- registerCollector，要求知道collector提供名字和指标名称
- exportMetric，将metric以map[string]any的形式返回
*/

type Collector interface {
	Describe(ch chan<- *prometheus.Desc)
	Collect(ch chan<- prometheus.Metric)
}

func RegisterCollector(collector Collector) {
	prometheus.MustRegister(collector)
}
