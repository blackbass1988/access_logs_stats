package input

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type FileInputReader struct {
	InputBufferedReader

	file *os.File
	fi   os.FileInfo

	fileReader *bufio.Reader

	mutex sync.Mutex

	buffer []byte
}

func CreateFileReader(dsn string) (r *FileInputReader, err error) {
	filename := strings.Replace(dsn, "file:", "", 1)

	r = &FileInputReader{}
	r.buffer = []byte{}
	r.openFile(filename)
	r.file.Seek(0, 2)

	go r.checkFile()

	return r, err
}

func (r *FileInputReader) openFile(filename string) {
	var err error
	r.file, err = os.Open(filename)
	r.fi, err = r.file.Stat()
	r.fileReader = bufio.NewReader(r.file)
	check(err)
}

func (r *FileInputReader) checkFile() {
	tick1s := time.Tick(1 * time.Second)
	for {
		select {
		case <-tick1s:
			fi, err := os.Stat(r.file.Name())
			check(err)
			if !os.SameFile(fi, r.fi) {
				log.Println("reopen input file")
				r.mutex.Lock()
				r.file.Close()
				r.openFile(r.file.Name())
				r.mutex.Unlock()

			}
		}
	}
}

func (r *FileInputReader) Close() {
	r.file.Close()
}

func (r *FileInputReader) ReadToBuffer() {
	log.Println("reading...")
	for {
		r.mutex.Lock()
		bytesBuf, err := r.fileReader.ReadBytes('\n')
		r.mutex.Unlock()
		if err == io.EOF {
			time.Sleep(10 * time.Millisecond)
			continue
		} else if err != nil {
			check(err)
		}
		r.mutex.Lock()
		r.buffer = append(r.buffer, bytesBuf...)
		r.mutex.Unlock()
	}
}

func (r *FileInputReader) FlushBuffer() []byte {
	r.mutex.Lock()
	buffer := r.buffer
	r.buffer = []byte{}
	r.mutex.Unlock()
	return buffer
}
