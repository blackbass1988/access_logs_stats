package core

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/blackbass1988/access_logs_stats/core/input"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	ERR_EMPTY_RESULT    error = errors.New("bad string or regular expression")
	ERR_FILTERS_NOT_SET error = errors.New("filters not set")
	ERR_OUTPUT_NOT_SET  error = errors.New("there are least one output must be specified. 0 found")
)

type Row struct {
	fields map[string]string
	raw    string
}

type App struct {
	fi     os.FileInfo
	file   *os.File
	config Config
	buffer []byte

	subProcessCollection *SubProcessCollection
	ir                   input.InputBufferedReader

	fileReader *bufio.Reader

	processBufferSync chan bool
	m                 sync.Mutex
}

func (a *App) openReader() (err error) {
	if strings.HasPrefix(a.config.InputDsn, "file:") {
		a.ir, err = input.CreateFileReader(a.config.InputDsn)
	} else if strings.HasPrefix(a.config.InputDsn, "syslog:") {
		a.ir, err = input.CreateSyslogInputReader(a.config.InputDsn)
	} else {
		err = errors.New("unknown input type")
	}
	return err
}

func (a *App) Start() {
	var err error
	a.init()

	tick := time.Tick(a.config.Period)
	log.Println("start a reading...")
	err = a.openReader()
	check(err)

	defer func() {
		a.ir.Close()
	}()

	go a.ir.ReadToBuffer()

	for {
		select {
		case <-tick:
			a.processBufferSync <- true
			go a.processBuffer()
		}
	}
}

func (a *App) init() {
	a.processBufferSync = make(chan bool, 1)
	a.buffer = []byte{}
	a.subProcessCollection = NewSubProcessCollection(&a.config)
}

func (a *App) processBuffer() {

	var (
		rawString  string
		err        error
		lastString string
	)
	a.m.Lock()
	buffer := a.ir.FlushBuffer()
	a.m.Unlock()
	byteReader := bytes.NewReader(buffer)
	bufReader := bufio.NewReader(byteReader)

	a.subProcessCollection.resetData()

	for {
		rawString, err = bufReader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			check(err)
		}

		logRow, err := a.NewRow(rawString)
		if err != nil && err == ERR_EMPTY_RESULT {
			log.Println(err, rawString)
			continue
		}
		check(err)

		a.subProcessCollection.appendData(logRow)
		lastString = rawString
	}

	go a.subProcessCollection.sendStats()
	log.Println(lastString)
	<-a.processBufferSync
}

func (a *App) NewRow(rawString string) (row *Row, err error) {

	row = new(Row)
	row.fields = make(map[string]string)
	row.raw = rawString

	matches := a.config.Rex.FindStringSubmatch(rawString)

	if len(matches) == 0 {
		return nil, ERR_EMPTY_RESULT
	}

	for i, name := range a.config.Rex.SubexpNames() {
		row.fields[name] = matches[i]
	}

	return row, err
}

func NewApp(config Config) (app *App, err error) {
	app = new(App)
	app.config = config
	return app, err
}
