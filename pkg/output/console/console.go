package console

import (
	"log"

	"github.com/blackbass1988/access_logs_stats/pkg/output"
)

var c *console

type console struct {
}

func (c *console) send(key string, value string) {
	log.Printf("%s = %s\n", key, value)
}

//Send sends messages to console
func Send(messages []*output.Message) {

	for _, message := range messages {
		c.send(message.Key, message.Value)
	}

}

//Init initializes console sender
func Init(params map[string]string) {
	//nothing to do?
}

func init() {
	c = new(console)
	output.RegisterOutput("console", Send, Init)
}
