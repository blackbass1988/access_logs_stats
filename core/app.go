package core

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/blackbass1988/access_logs_stats/core/input"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

var (
	PROG_NAME  string = "AccessLogsStats"
	VERSION    string = "dev"
	BUILD_TIME string = "dev"

	ERR_EMPTY_RESULT    error = errors.New("bad string or regular expression")
	ERR_FILTERS_NOT_SET error = errors.New("filters not set")
	ERR_OUTPUT_NOT_SET  error = errors.New("there are least one output must be specified. 0 found")
)

type Row struct {
	Fields map[string]string
	Raw    string
}

type App struct {
	fi     os.FileInfo
	file   *os.File
	config Config
	buffer []byte

	senderCollection *SenderCollection
	ir               input.InputBufferedReader

	fileReader *bufio.Reader

	processBufferSync chan bool
	m                 sync.Mutex
}

func (a *App) openReader() (err error) {
	a.ir, err = input.GetFileReader(a.config.InputDsn)
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

	if a.config.ExitAfterOneTick {
		a.ir.ReadToBuffer()
		a.processBuffer()

	} else {
		//read to buffer in background
		go a.ir.ReadToBuffer()

		for {
			select {
			case <-tick:
				a.processBufferSync <- true
				go a.processBuffer()
			}
		}
	}
}

func (a *App) Stop() {
	os.Exit(0)
}

func (a *App) init() {
	a.processBufferSync = make(chan bool, 1)
	a.buffer = []byte{}
	a.senderCollection = NewSenderCollection(&a.config)
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

	a.senderCollection.resetData()

	for {
		rawString, err = bufReader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			check(err)
		}

		logRow, err := NewRow(rawString, a.config.Rex)
		if err != nil && err == ERR_EMPTY_RESULT {
			log.Println(err, rawString)
			continue
		}
		check(err)

		a.senderCollection.appendData(logRow)
		lastString = rawString
	}
	log.Println(lastString)

	if a.config.ExitAfterOneTick {
		a.senderCollection.sendStats()
	} else {
		go a.senderCollection.sendStats()
		<-a.processBufferSync
	}
}

func NewApp(config Config) (app *App, err error) {
	app = new(App)
	app.config = config
	return app, err
}
