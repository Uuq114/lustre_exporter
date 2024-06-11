package exporter

import (
	"fmt"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"net/http"
	"os"
)

var (
	version       = "0.0.1"
	configFile    = "config/config.yaml"
	collectorFile = "config/collector.yaml"
)

//var (
//	collectors         []string
//	collectorFactories map[string]func(log.Logger) collector.Collector
//)

type LustreExporter struct {
	config Config
	logger log.Logger
}

type Config struct {
	port string
}

func (le *LustreExporter) init() {
	le.logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	le.logger = log.With(le.logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
}

func (le *LustreExporter) InitConfig() {
	le.init()
	le.logger.Log("version", version, "msg", "lustre exporter init now.")
	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		le.logger.Log("msg", "read config file fail", "error", err.Error())
	}
	le.config = Config{
		port: viper.GetString("port"),
	}
	le.logger.Log("config", fmt.Sprintf("%+v", le.config))
}

//func RegisterCollector(name string, factory func(logger log.Logger) collector.Collector) {
//	collectors = append(collectors, name)
//	collectorFactories[name] = factory
//}

//func (le *LustreExporter) loadCollector() {
//	for _, name := range collector. {
//		logger := promlog.New(&promlog.Config{})
//		prometheus.MustRegister(collectorFactories[name](log.With(logger, "collector", name)))
//	}
//}

type HelloHandler struct{}

func (h *HelloHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte(`<html>
             <head><title>GPFS Exporter</title></head>
             <body>
             <h1>GPFS Metrics Exporter</h1>
             <p><a href='/metrics'>Metrics</a></p>
             </body>
             </html>`))
}

func newHelloHandler() *HelloHandler {
	return &HelloHandler{}
}

func main() {
	le := &LustreExporter{}
	le.InitConfig()

	//exampleCollector := collector.NewExampleCollector()
	//collector.RegisterCollector(exampleCollector)

	le.loadCollector()

	serverUrl := fmt.Sprintf(":%v", le.config.port)
	http.Handle("/", newHelloHandler())
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(serverUrl, nil); err != nil {
		le.logger.Log("msg", "start http server fail.", "error", err.Error())
	}
}
