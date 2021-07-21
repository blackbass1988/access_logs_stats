package console

import (
	"github.com/blackbass1988/access_logs_stats/pkg/output"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

type PrometheusServer struct {
	registry *prometheus.Registry
}

func (ps *PrometheusServer) Send(_ []*output.Message) {
}

func (ps *PrometheusServer) Init(params map[string]string, _ map[string]string) {
	if params["addr"] == "" {
		log.Fatalln("for prometheus output settings 'addr' must set")
	}
	ps.registry = prometheus.NewRegistry()
	http.Handle("/metrics", promhttp.HandlerFor(ps.registry, promhttp.HandlerOpts{}))
	go func() {
		err := http.ListenAndServe(params["addr"], nil)
		if err != nil {
			panic(err)
		}
	}()

}

func (ps *PrometheusServer) RegisterPrometheusCollector(collector prometheus.Collector) {
	ps.registry.MustRegister(collector)
}

func init() {
	prom := &PrometheusServer{}
	output.RegisterOutput("prometheus", prom.Send, prom.Init, prom.RegisterPrometheusCollector)
}
