package input

import (
	"errors"
	"log"
	"strings"
)

//BufferedReader describes interface of implementations
type BufferedReader interface {
	ReadToBuffer()
	FlushBuffer() []byte
	Close()
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//GetFileReader is a "factory method"
func GetFileReader(inputDsn string) (BufferedReader, error) {

	var err error
	var r BufferedReader

	if strings.HasPrefix(inputDsn, "file:") {
		r, err = CreateFileReader(inputDsn)
	} else if strings.HasPrefix(inputDsn, "syslog:") {
		r, err = CreateSyslogInputReader(inputDsn)
	} else if strings.HasPrefix(inputDsn, "stdin:") {
		r, err = CreateStdinReader(inputDsn)
	} else {
		err = errors.New("unknown input type: " + inputDsn)
	}

	return r, err

}
