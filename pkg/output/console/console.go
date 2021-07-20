package console

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"

	"github.com/blackbass1988/access_logs_stats/pkg/output"
	"github.com/blackbass1988/access_logs_stats/pkg/template"
)

var c *console

type console struct {
	template     *template.Template
	templateVars map[string]string
}

func (c *console) send(field string, metric string, value string) {

	key := c.template.Process(field, metric, c.templateVars)

	log.Printf("%s = %s\n", key, value)

}

//Send sends messages to console
func Send(messages []*output.Message) {

	for _, message := range messages {
		c.send(message.Field, message.Metric, message.Value)
	}

}

//Init initializes console sender
func Init(params map[string]string, templateVars map[string]string) {
	var templateString string
	var ok bool

	c.templateVars = templateVars
	if templateString, ok = params["template"]; !ok {
		templateString = output.DefaultTemplate
	}

	c.template = template.NewTemplate(templateString)
}

func RegisterPrometheusCollector(collector prometheus.Collector) {

}

func init() {
	c = new(console)
	output.RegisterOutput("console", Send, Init, RegisterPrometheusCollector)
}
