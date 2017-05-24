package input

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

//FileInputReader is a BufferedReader for a file reading
type FileInputReader struct {
	BufferedReader

	file *os.File
	fi   os.FileInfo

	fileReader *bufio.Reader

	m sync.Mutex

	buffer []byte
}

//CreateFileReader create new FileInputReader
func CreateFileReader(dsn string) (r *FileInputReader, err error) {
	filename := strings.Replace(dsn, "file:", "", 1)

	r = &FileInputReader{}
	r.buffer = []byte{}
	r.openFile(filename)
	r.file.Seek(0, 2)

	go r.checkFile()

	return r, err
}

//Close implements Close method of BufferedReader for FileInputReader
func (r *FileInputReader) Close() {
	r.file.Close()
}

//ReadToBuffer implements ReadToBuffer method of BufferedReader for FileInputReader
func (r *FileInputReader) ReadToBuffer() {
	log.Println("reading...")
	for {
		r.m.Lock()
		bytesBuf, err := r.fileReader.ReadBytes('\n')
		r.m.Unlock()
		if err == io.EOF {
			time.Sleep(10 * time.Millisecond)
			continue
		} else if err != nil {
			check(err)
		}
		r.m.Lock()
		r.buffer = append(r.buffer, bytesBuf...)
		r.m.Unlock()
	}
}

//FlushBuffer implements FlushBuffer method of BufferedReader for FileInputReader
func (r *FileInputReader) FlushBuffer() []byte {
	r.m.Lock()
	buffer := r.buffer
	r.buffer = []byte{}
	r.m.Unlock()
	return buffer
}

func (r *FileInputReader) openFile(filename string) {
	var err error

	r.file, err = os.Open(filename)
	r.fi, err = r.file.Stat()
	if err != nil && !os.IsExist(err) {
		err = fmt.Errorf("file \"%s\" not exists", filename)
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
				r.m.Lock()
				r.file.Close()
				r.openFile(r.file.Name())
				r.m.Unlock()
			}
			prevSize = fi.Size()
		}
	}
}
