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
	buffer []byte
	nowait bool
}

//CreateStdinReader creates new StdInputReader
func CreateStdinReader(dsn string) (r *StdInputReader, err error) {
	r = &StdInputReader{}
	r.buffer = []byte{}

	if option := strings.Replace(dsn, "stdin:", "", 1); option == "nowait" || option == "" {
		r.nowait = true
	} else {
		return nil, errors.New("unknown or not implemented option: " + option)
	}
	return r, err
}

//ReadToBuffer implements ReadToBuffer method of BufferedReader for StdInputReader
func (r *StdInputReader) ReadToBuffer() {
	var (
		b   byte
		err error
	)
	r.reader = bufio.NewReader(os.Stdin)

	for {
		b, err = r.reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
		} else {
			r.buffer = append(r.buffer, b)
		}
	}
}

//FlushBuffer implements FlushBuffer method of BufferedReader for StdInputReader
func (r *StdInputReader) FlushBuffer() []byte {
	b := r.buffer
	r.buffer = []byte{}
	return b
}

//Close implements Close method of BufferedReader for StdInputReader
func (r *StdInputReader) Close() {
}
