package output

type output struct {
	name    string
	send    func([]*Message)
	init    func(map[string]string)
	enabled bool
}

var outputs = []output{}

type Message struct {
	Key   string
	Value string
}

func RegisterOutput(name string, send func(messages []*Message), init func(params map[string]string)) error {
	outputs = append(outputs, output{name, send, init, false})
	return nil
}

type Output struct {
	prefix   string
	messages []*Message
}

/*
 * set common prefix for all keys
 */
func (s *Output) SetPrefix(prefix string) {
	s.prefix = prefix
}

/*
 * add message to message pack
 */
func (s *Output) AddMessage(key string, value string) {

	if len(s.prefix) > 0 {
		key = s.prefix + key
	}
	m := new(Message)
	m.Key = key
	m.Value = value
	s.messages = append(s.messages, m)
}

/*
 * sends message pack by output
 */
func (s *Output) Send() {

	for _, aOutput := range outputs {
		if aOutput.enabled {
			aOutput.send(s.messages)
		}
	}
	s.messages = []*Message{}
}

func (s *Output) Init(senderName string, params map[string]string) {
	for i, aOutput := range outputs {
		if aOutput.name == senderName {
			outputs[i].enabled = true
			aOutput.init(params)
			break
		}
	}
}
