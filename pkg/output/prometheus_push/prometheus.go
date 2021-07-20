package console

import (
	"github.com/blackbass1988/access_logs_stats/pkg/output"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"log"
)

type PrometheusPushGw struct {
	registry *prometheus.Registry
	pusher   *push.Pusher
}

func (p *PrometheusPushGw) Send(_ []*output.Message) {
	err := p.pusher.Push()
	if err != nil {
		log.Printf("cant push to prometheus: %s\n", err.Error())
	}
}

func (p *PrometheusPushGw) Init(params map[string]string, _ map[string]string) {
	if params["url"] == "" || params["job"] == "" {
		log.Fatalln("for prometheus_push output settings 'job' and 'urls' must set")
	}
	p.pusher = push.New(params["url"], params["job"])
	p.registry = prometheus.NewRegistry()
}

func (p *PrometheusPushGw) RegisterPrometheusCollector(collector prometheus.Collector) {
	p.registry.MustRegister(collector)
}


func init() {
	p := &PrometheusPushGw{}
	output.RegisterOutput("prometheus_push", p.Send, p.Init, p.RegisterPrometheusCollector)
}
