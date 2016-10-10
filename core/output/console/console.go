package console

import (
	"github.com/blackbass1988/access_logs_stats/core/output"
	"log"
)

var c *console

type console struct {
}

func (c *console) send(key string, value string) {
	log.Printf("%s = %s\n", key, value)
}

func Send(messages []*output.Message) {

	for _, message := range messages {
		c.send(message.Key, message.Value)
	}

}

func Init(params map[string]string) {
	//nothing to do?
}

func init() {
	c = new(console)
	output.RegisterOutput("console", Send, Init)
}
