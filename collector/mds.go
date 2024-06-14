package collector

import (
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"os/exec"
)

const (
	GetMDTUsageCommand = "lctl get_param osd-*.*MDT*.kbytestotal"
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
	output, err := getMDTInfo()
	if err != nil {
		mc.logger.Log("error", err.Error())
	}
	metric := parseMDTInfo(output)

	for mdt, mdtMetric := range metric.MDTList {
		ch <- prometheus.MustNewConstMetric(mc.KBFree, prometheus.GaugeValue, float64(mdtMetric.KBFree), mdt)
		ch <- prometheus.MustNewConstMetric(mc.KBTotal, prometheus.GaugeValue, float64(mdtMetric.KBTotal), mdt)
		ch <- prometheus.MustNewConstMetric(mc.InodesFree, prometheus.GaugeValue, float64(mdtMetric.InodesFree), mdt)
		ch <- prometheus.MustNewConstMetric(mc.InodesTotal, prometheus.GaugeValue, float64(mdtMetric.InodesTotal), mdt)
	}
}

func getMDTInfo() (string, error) {
	output, err := execCommand(GetMDTUsageCommand)
	if err != nil {

	}
}

func parseMDTInfo(output string) MDSSpaceMetric {

}

func execCommand(command string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", command)

	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	} else {
		return string(stdout), nil
	}
}
