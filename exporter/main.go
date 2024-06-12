package main

import (
	"fmt"
	"github.com/Uuq114/lustre_exporter/collector"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type HelloHandler struct{}

func (h *HelloHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte(`<html>
             <head><title>Lustre Exporter</title></head>
             <body>
             <h1>Lustre Metrics Exporter</h1>
             <p><a href='/metrics'>Metrics</a></p>
             </body>
             </html>`))
}

func newHelloHandler() *HelloHandler {
	return &HelloHandler{}
}

func main() {
	le := collector.NewLustreExporter()
	le.LoadCollector()

	serverUrl := fmt.Sprintf(":%v", le.Config.Port)
	http.Handle("/", newHelloHandler())
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(serverUrl, nil); err != nil {
		le.Logger.Log("msg", "start http server fail.", "error", err.Error())
	}
}
