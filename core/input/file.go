package input

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type FileInputReader struct {
	BufferedReader

	file *os.File
	fi   os.FileInfo

	fileReader *bufio.Reader

	chanLock chan bool

	buffer []byte
}

func CreateFileReader(dsn string) (r *FileInputReader, err error) {
	filename := strings.Replace(dsn, "file:", "", 1)

	r = &FileInputReader{}
	r.chanLock = make(chan bool, 1)
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
	if err != nil && !os.IsExist(err) {
		err = errors.New(fmt.Sprintf("file \"%s\" not exists", filename))
	}
	r.fileReader = bufio.NewReader(r.file)
	check(err)
}

func (r *FileInputReader) checkFile() {
	var prevSize int64
	prevSize = -1
	tick1s := time.Tick(1 * time.Second)
	for {
		select {
		case <-tick1s:
			fi, err := os.Stat(r.file.Name())
			check(err)

			if prevSize == -1 {
				prevSize = fi.Size()
			}

			if !os.SameFile(fi, r.fi) || prevSize > fi.Size() {
				log.Println("reopen input file")
				r.chanLock <- true
				r.file.Close()
				r.openFile(r.file.Name())
				<-r.chanLock
			}
			prevSize = fi.Size()
		}
	}
}

func (r *FileInputReader) Close() {
	r.file.Close()
}

func (r *FileInputReader) ReadToBuffer() {
	log.Println("reading...")
	buffChan := make(chan []byte)
	go r.readToBuffer(buffChan)
	go r.writeToBuffer(buffChan)
}

func (r *FileInputReader) readToBuffer(buffChan chan<- []byte) {
	for {
		r.chanLock <- true
		bytesBuf, err := r.fileReader.ReadBytes('\n')
		<-r.chanLock
		if err == io.EOF {
			time.Sleep(10 * time.Millisecond)
			continue
		} else if err != nil {
			check(err)
		}
		buffChan <- bytesBuf
	}
}

func (r *FileInputReader) writeToBuffer(buffChan <-chan []byte) {
	for {
		r.buffer = append(r.buffer, <-buffChan...)
	}
}

func (r *FileInputReader) FlushBuffer() []byte {
	r.chanLock <- true
	buffer := r.buffer
	r.buffer = []byte{}
	<-r.chanLock
	return buffer
}
