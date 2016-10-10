package input

import "log"

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
