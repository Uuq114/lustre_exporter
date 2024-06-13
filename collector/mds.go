package collector

import (
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

type MDSSpaceCollector struct {
	KBFree      *prometheus.Desc
	KBTotal     *prometheus.Desc
	InodesFree  *prometheus.Desc
	InodesTotal *prometheus.Desc
	logger      log.Logger
}

type MDSSpaceMetric struct {
	MDTList map[string]MDTMetric
}

type MDTMetric struct {
	KBFree      int
	KBTotal     int
	InodesFree  int
	InodesTotal int
}

func NewMDSSpaceCollector(logger log.Logger) Collector {
	return &MDSSpaceCollector{
		KBFree: prometheus.NewDesc(prometheus.BuildFQName(namespace, "mds", "kbfree"),
			"Lustre MDS free space in kilobytes", nil, nil),
		KBTotal: prometheus.NewDesc(prometheus.BuildFQName(namespace, "mds", "kbtotal"),
			"Lustre MFS total space in kilobytes", nil, nil),
		InodesFree: prometheus.NewDesc(prometheus.BuildFQName(namespace, "mds", "inodes_free"),
			"Lustre MFS free inodes", nil, nil),
		InodesTotal: prometheus.NewDesc(prometheus.BuildFQName(namespace, "mds", "inodes_total"),
			"Lustre MFS total inodes", nil, nil),
		logger: logger,
	}
}

func init() {
	registerCollector("mds", NewMDSSpaceCollector)
}

func (mc *MDSSpaceCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- mc.KBFree
	ch <- mc.KBTotal
	ch <- mc.InodesFree
	ch <- mc.InodesTotal
}

func (mc *MDSSpaceCollector) Collect(ch chan<- prometheus.Metric) {

}

func getMDTInfo() (string, error) {

}
