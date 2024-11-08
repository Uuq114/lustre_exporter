package collector

import (
	"errors"
	"os/exec"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	MDTTotalSpaceCommand = "lctl get_param osd-*.*MDT*.kbytestotal"
	MDTFreeSpaceCommand  = "lctl get_param osd-*.*MDT*.kbytesfree"
	MDTTotalInodeCommand = "lctl get_param osd-*.*MDT*.filestotal"
	MDTFreeInodeCommand  = "lctl get_param osd-*.*MDT*.filesfree"
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
	KBFree      int64
	KBTotal     int64
	InodesFree  int64
	InodesTotal int64
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
	commands := []string{MDTTotalSpaceCommand, MDTFreeSpaceCommand}
	for _, command := range commands {
		output, err := execCommand(command)
		if err != nil {
			mc.logger.Log("error", err.Error())
		}
		parseMDTInfo(output)
	}

	// export parsed data
	for mdt, mdtMetric := range mdsMetric.MDTList {
		ch <- prometheus.MustNewConstMetric(mc.KBFree, prometheus.GaugeValue, float64(mdtMetric.KBFree), mdt)
		ch <- prometheus.MustNewConstMetric(mc.KBTotal, prometheus.GaugeValue, float64(mdtMetric.KBTotal), mdt)
		ch <- prometheus.MustNewConstMetric(mc.InodesFree, prometheus.GaugeValue, float64(mdtMetric.InodesFree), mdt)
		ch <- prometheus.MustNewConstMetric(mc.InodesTotal, prometheus.GaugeValue, float64(mdtMetric.InodesTotal), mdt)
	}
}

// converse metric name to struct field, make parsing easier
func mapMetricNameToField(name string) string {
	mapping := map[string]string{
		"kbytestotal": "KBTotal",
		"kbytesfree":  "KBFree",
		"filesfree":   "InodesFree",
		"filestotal":  "InodesTotal",
	}
	if value, ok := mapping[name]; ok {
		return value
	} else {
		return ""
	}
}

// parse metrics into a global map
func parseMDTInfo(output string) error {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		words := strings.Split(line, "=")
		mdtName := strings.Split(words[0], ".")[1]
		metricName := strings.Split(words[0], ".")[2]
		value, _ := strconv.ParseInt(words[1], 10, 64)

		fieldName := mapMetricNameToField(metricName)
		if fieldName == "" {
			return errors.New("cannot map metric name to valid struct field")
		}
		if mdt, ok := mdsMetric.MDTList[mdtName]; ok {
			MDTReflectValue := reflect.ValueOf(mdt).Elem()
			field := MDTReflectValue.FieldByName(fieldName)
			if field.IsValid() && field.CanSet() {
				field.SetInt(value)
			}
		} else {
			var mdt MDTMetric
			MDTReflectValue := reflect.ValueOf(mdt).Elem()
			field := MDTReflectValue.FieldByName(fieldName)
			if field.IsValid() && field.CanSet() {
				field.SetInt(value)
			}
			mdsMetric.MDTList[mdtName] = mdt
		}
	}
	return nil
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
