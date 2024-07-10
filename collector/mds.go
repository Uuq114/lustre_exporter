package collector

import (
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"os/exec"
	"strconv"
	"strings"
)

const (
	MDTTotalSpaceCommand = "lctl get_param osd-*.*MDT*.kbytestotal"
)

// All MDS metrics goes here
var (
	mdsMetric MDSSpaceMetric
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
	mdsMetric.MDTList = make(map[string]MDTMetric)
}

func (mc *MDSSpaceCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- mc.KBFree
	ch <- mc.KBTotal
	ch <- mc.InodesFree
	ch <- mc.InodesTotal
}

func (mc *MDSSpaceCollector) Collect(ch chan<- prometheus.Metric) {
	// execute some commands, parse as metric type
	output, err := execCommand(MDTTotalSpaceCommand)
	if err != nil {
		mc.logger.Log("error", err.Error())
	}
	parseMDTInfo(output)

	// export parsed data
	for mdt, mdtMetric := range mdsMetric.MDTList {
		ch <- prometheus.MustNewConstMetric(mc.KBFree, prometheus.GaugeValue, float64(mdtMetric.KBFree), mdt)
		ch <- prometheus.MustNewConstMetric(mc.KBTotal, prometheus.GaugeValue, float64(mdtMetric.KBTotal), mdt)
		ch <- prometheus.MustNewConstMetric(mc.InodesFree, prometheus.GaugeValue, float64(mdtMetric.InodesFree), mdt)
		ch <- prometheus.MustNewConstMetric(mc.InodesTotal, prometheus.GaugeValue, float64(mdtMetric.InodesTotal), mdt)
	}
}

// there may be multiple MDTs on one node, so return a map here
func parseMDTInfo(output string) {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		words := strings.Split(line, "=")
		mdtName := strings.Split(strings.Split(words[0], ".")[1], "-")[1]
		kbTotal, _ := strconv.Atoi(words[1])

		if mdt, ok := mdsMetric.MDTList[mdtName]; ok {
			mdt.KBTotal = kbTotal
		} else {
			mdsMetric.MDTList[mdtName] = MDTMetric{KBTotal: kbTotal}
		}
	}
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
