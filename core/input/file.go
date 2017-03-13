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
	tick1s := time.Tick(1 * time.Second)
	for {
		select {
		case <-tick1s:
			fi, err := os.Stat(r.file.Name())
			check(err)
			if !os.SameFile(fi, r.fi) {
				log.Println("reopen input file")
				r.chanLock <- true
				r.file.Close()
				r.openFile(r.file.Name())
				<-r.chanLock
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
		r.chanLock <- true
		bytesBuf, err := r.fileReader.ReadBytes('\n')
		<-r.chanLock
		if err == io.EOF {
			time.Sleep(10 * time.Millisecond)
			continue
		} else if err != nil {
			check(err)
		}
		r.chanLock <- true
		r.buffer = append(r.buffer, bytesBuf...)
		<-r.chanLock
	}
}

func (r *FileInputReader) FlushBuffer() []byte {
	r.chanLock <- true
	buffer := r.buffer
	r.buffer = []byte{}
	<-r.chanLock
	return buffer
}
