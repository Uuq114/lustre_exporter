package main

import (
	"fmt"
	"github.com/go-kit/log"
	"github.com/spf13/viper"
	"net/http"
	"os"
)

var (
	version    = "0.0.1"
	configFile = "config/config.yaml"
)

type LustreExporter struct {
	config Config
	logger log.Logger
}

type Config struct {
	port string
}

func (le *LustreExporter) init() {
	le.logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
}

func (le *LustreExporter) InitConfig() {
	le.init()
	le.logger.Log("msg", "lustre exporter init now")
	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		le.logger.Log("msg", "read config fail", "err", err.Error())
	}
	le.config = Config{
		port: viper.GetString("port"),
	}
	le.logger.Log("config", fmt.Sprintf("%+v", le.config))
}

func landingHandler(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte(`<html>
             <head><title>GPFS Exporter</title></head>
             <body>
             <h1>Lustre Metrics Exporter</h1>
             <p><a href='/metrics'>Metrics</a></p>
             </body>
             </html>`))
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<html>
             <head><title>GPFS Exporter</title></head>
             <body>
             <h1>Lustre Metrics Exporter</h1>
             <p>Metrics Page</a></p>
             </body>
             </html>`))
}

func main() {
	le := &LustreExporter{}
	le.InitConfig()

	http.HandleFunc("/", landingHandler)
	http.HandleFunc("/metrics", metricsHandler)

	le.logger.Log("msg", "lustre exporter is working")
	serverUrl := fmt.Sprintf(":%v", le.config.port)
	if err := http.ListenAndServe(serverUrl, nil); err != nil {
		le.logger.Log("msg", "start http server fail", "err", err.Error())
	}
}
