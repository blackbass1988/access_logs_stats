package input

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
)

//StdInputReader implements reading from stdin
type StdInputReader struct {
	BufferedReader

	reader *bufio.Reader
	nowait bool
}

//CreateStdinReader creates new StdInputReader
func CreateStdinReader(dsn string) (r *StdInputReader, err error) {
	r = &StdInputReader{}

	if option := strings.Replace(dsn, "stdin:", "", 1); option == "nowait" || option == "" {
		r.nowait = true
	} else {
		return nil, errors.New("unknown or not implemented option: " + option)
	}
	return r, err
}

//ReadToChannel implements ReadToChannel
func (r *StdInputReader) ReadToChannel(lineChannel chan<- string) {
	var (
		b   []byte
		err error
	)

	r.reader = bufio.NewReader(os.Stdin)

	for {
		b, err = r.reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				close(lineChannel)
				break
			}
		} else {
			lineChannel <- string(b)
		}
	}
}

//Close implements Close method of BufferedReader for StdInputReader
func (r *StdInputReader) Close() {
}
