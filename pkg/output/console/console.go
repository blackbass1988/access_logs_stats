package console

import (
	"github.com/blackbass1988/access_logs_stats/pkg/output"
	"log"
)

var c *console

type console struct {
	template     *output.Template
	templateVars map[string]string
}

func (c *console) send(field string, metric string, value string) {

	err, key := c.template.Process(field, metric, c.templateVars)

	if err != nil {
		log.Println("ERROR:", err)
	} else {
		log.Printf("%s = %s\n", key, value)
	}

}

//Send sends messages to console
func Send(messages []*output.Message) {

	for _, message := range messages {
		c.send(message.Field, message.Metric, message.Value)
	}

}

//Init initializes console sender
func Init(params map[string]string, templateVars map[string]string) {
	var err error
	var templateString string
	var ok bool

	c.templateVars = templateVars
	if templateString, ok = params["template"]; !ok {
		templateString = output.DefaultTemplate
	}

	err, c.template = output.NewTempate(templateString)

	if err != nil {
		log.Fatalf("template init failed for template \"%s\". Error: \"%s\"", templateString, err.Error())
	}

}

func init() {
	c = new(console)
	output.RegisterOutput("console", Send, Init)
}
