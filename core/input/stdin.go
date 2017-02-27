package input

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
)

type StdInputReader struct {
	InputBufferedReader

	reader *bufio.Reader
	buffer []byte
	nowait bool
}

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

func (r *StdInputReader) FlushBuffer() []byte {
	b := r.buffer
	r.buffer = []byte{}
	return b
}

func (r *StdInputReader) Close() {
}
