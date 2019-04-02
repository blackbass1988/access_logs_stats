package zabbix

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"net"
	"os"

	"github.com/blackbass1988/access_logs_stats/pkg/output"
	"github.com/blackbass1988/access_logs_stats/pkg/template"
)

var (
	z      *zabbix
	header = []byte("ZBXD\x01")
)

//@link https://www.zabbix.org/wiki/Docs/protocols/zabbix_sender/2.0
type message struct {
	Request string `json:"request"`
	Data    []data `json:"data"`
}

type data struct {
	Host  string `json:"host"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type zabbix struct {
	zabbixHost string
	zabbixPort string
	host       string

	template *template.Template

	templateVars map[string]string
}

func (z *zabbix) getData(messages []*output.Message) []data {
	var els []data

	for _, message := range messages {
		key := message.Field + "_" + message.Metric

		err, key := z.template.Process(message.Field, message.Metric, z.templateVars)

		if err != nil {
			//ok for dev version
			log.Panic(err)
		}

		el := data{Host: z.host, Key: key, Value: message.Value}
		els = append(els, el)
	}
	return els
}

func (z *zabbix) send(messages []*output.Message) {
	//todo refact
	//todo persist connect?

	//generate json
	d := z.getData(messages)

	m := message{"sender data", d}
	jsonBytes, err := json.Marshal(m)
	if err != nil {
		log.Println("json marshal error:", err)
		return
	}

	//send to server
	conn, err := net.Dial("tcp4", z.zabbixHost+":"+z.zabbixPort)
	if err != nil {
		log.Println("zabbix connect error:", err)
		return
	}
	defer conn.Close()

	length := len(jsonBytes)

	buf := bytes.NewBuffer(header)
	binary.Write(buf, binary.LittleEndian, uint64(length))
	buf.Write(jsonBytes)
	_, err = conn.Write(buf.Bytes())
	if err != nil {
		log.Println("zabbix socket write error:", err)
	}

	//read response
	response := []byte{}
	tmp := make([]byte, 64)
	for {
		_, err := conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				log.Print("zabbix socket read error:", err)
			}
			break
		}
		response = append(response, tmp...)
	}

	//if failed - > log it!
	//log.Println(string(buf.Bytes()))
	//log.Println(string(response))

}

//Send sends messages to zabbix
func Send(messages []*output.Message) {
	z.send(messages)
}

//Init initializes zabbix sender
func Init(params map[string]string, templateVars map[string]string) {
	var err error

	z.templateVars = templateVars
	templateString := output.DefaultTemplate

	for k, v := range params {
		if v == "${hostname}" {
			v, _ = os.Hostname()
		}

		switch k {
		case "zabbix_host":
			z.zabbixHost = v
		case "zabbix_port":
			z.zabbixPort = v
		case "host":
			z.host = v
		case "template":
			templateString = v
		}
	}

	err, z.template = template.NewTempate(templateString)
	if err != nil {
		log.Fatalln("invalid template", templateString, "error was:", err)
	}

	if z.zabbixHost == "" || z.zabbixPort == "" || z.host == "" {
		log.Fatal("zabbix settings is incorrect. You must specify ",
			"zabbix_host, zabbix_port and host")
	}
}

func init() {
	z = new(zabbix)
	output.RegisterOutput("zabbix", Send, Init)
}
