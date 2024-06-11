package collector

import (
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

/*
注册collector需要实现Collector接口，要求要实现三个方法：
- registerCollector，要求知道collector提供名字和指标名称
- exportMetric，将metric以map[string]any的形式返回
- Constructor，将collector的构造函数返回

每个自定义的collector必须有：
- 一个init函数，用于注册collector
*/

/* lustre exporter */

/* collector */

var (
	collectors    []string
	factories     map[string]func(logger log.Logger) Collector
	factoriesLock sync.Mutex
)

func registerCollector(name string, factory func(logger log.Logger) Collector) {
	collectors = append(collectors, name)
	factoriesLock.Lock()
	defer factoriesLock.Unlock()
	factories[name] = factory
}

type Collector interface {
	Describe(ch chan<- *prometheus.Desc)
	Collect(ch chan<- prometheus.Metric)
}

//func RegisterCollector(collector Collector) {
//	prometheus.MustRegister(collector)
//}
