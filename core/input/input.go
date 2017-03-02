package input

import (
	"errors"
	"log"
	"strings"
)

type InputBufferedReader interface {
	ReadToBuffer()
	FlushBuffer() []byte
	Close()
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func GetFileReader(inputDsn string) (InputBufferedReader, error) {

	var err error
	var r InputBufferedReader

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
