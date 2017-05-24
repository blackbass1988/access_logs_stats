package output

type output struct {
	name    string
	send    func([]*Message)
	init    func(map[string]string)
	enabled bool
}

var outputs = []output{}

//Message is key=value presentation of calculation
type Message struct {
	Key   string
	Value string
}

//RegisterOutput registers new output
func RegisterOutput(name string, send func(messages []*Message), init func(params map[string]string)) error {
	outputs = append(outputs, output{name, send, init, false})
	return nil
}

//Output is base struct of log target
type Output struct {
	prefix   string
	messages []*Message
}

//SetPrefix sets common prefix for all keys
func (s *Output) SetPrefix(prefix string) {
	s.prefix = prefix
}

//AddMessage adds message to message pack
func (s *Output) AddMessage(key string, value string) {

	if len(s.prefix) > 0 {
		key = s.prefix + key
	}
	m := new(Message)
	m.Key = key
	m.Value = value
	s.messages = append(s.messages, m)
}

//Send sends message pack by output
func (s *Output) Send() {

	currentMessages := s.messages
	for _, aOutput := range outputs {
		if aOutput.enabled {
			go aOutput.send(currentMessages)
		}
	}
	s.messages = []*Message{}
}

//Init initializes parent outputs
func (s *Output) Init(senderName string, params map[string]string) {
	for i, aOutput := range outputs {
		if aOutput.name == senderName {
			outputs[i].enabled = true
			aOutput.init(params)
			break
		}
	}
}
