package collector

import (
	"fmt"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/promlog"
	"github.com/spf13/viper"
	"os"
	"sync"
)

/*
自定义的collector需要实现Collector接口，要求要实现两个方法：
- Describe，metric的元数据？
- Collect，返回收集到的metric

每个自定义的collector必须有：
- 一个init函数，用于注册collector
*/

/* lustre exporter */

var (
	configFile = "config/config.yaml"
)

type LustreExporter struct {
	Config     LustreExporterConfig
	Collectors map[string]Collector
	sync.Mutex
	Logger log.Logger
}

type LustreExporterConfig struct {
	Port string
}

func NewLustreExporter() LustreExporter {
	return LustreExporter{
		Config:     parseConfig(configFile, "yaml"),
		Collectors: make(map[string]Collector),
		Logger:     log.With(log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr)), "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller),
	}
}

func parseConfig(configFile, configType string) LustreExporterConfig {
	viper.SetConfigFile(configFile)
	viper.SetConfigType(configType)
	viper.ReadInConfig()
	return LustreExporterConfig{
		Port: viper.GetString("port"),
	}
}

func (le *LustreExporter) LoadCollector() {
	le.Lock()
	for _, name := range collectors {
		logger := promlog.New(&promlog.Config{})
		collector := factories[name](logger)
		le.Collectors[name] = collector
		prometheus.MustRegister(collector)
	}
	defer le.Unlock()
}

/* collector */

var (
	collectors    []string
	factories     = make(map[string]func(logger log.Logger) Collector)
	factoriesLock sync.Mutex
)

type Collector interface {
	Describe(ch chan<- *prometheus.Desc)
	Collect(ch chan<- prometheus.Metric)
}

/* utils */

const (
	namespace = "Lustre"
)

func registerCollector(name string, factory func(logger log.Logger) Collector) {
	collectors = append(collectors, name)
	factoriesLock.Lock()
	defer factoriesLock.Unlock()
	factories[name] = factory
}

type ExecTimeoutError struct {
	Message string
}

func (e *ExecTimeoutError) Error() string {
	return fmt.Sprintf("Message: %s", e.Message)
}
