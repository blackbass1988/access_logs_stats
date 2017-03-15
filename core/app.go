package core

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/blackbass1988/access_logs_stats/core/input"
	"io"
	"log"
	"os"
	"regexp"
	"time"
)

var (

	//ProgName is a just name
	ProgName = "AccessLogsStats"
	//Version it is a version of application will be overridden on build
	Version = "dev"
	//BuildTime it is a build time of application will be overridden on build
	BuildTime = "dev"

	errEmptyResult   = errors.New("bad string or regular expression")
	errFiltersNotSet = errors.New("filters not set")
	errOutputNotSet  = errors.New("there are least one output must be specified. 0 found")
)

//RowEntry contains raw input string and parsed fields of it
type RowEntry struct {
	Fields map[string]string
	Raw    string
}

//App is a main struct of application
type App struct {
	fi     os.FileInfo
	file   *os.File
	config Config
	buffer []byte

	senderCollection *SenderCollection
	ir               input.BufferedReader

	fileReader *bufio.Reader

	processBufferSync chan bool
}

//NewApp creates new parser
func NewApp(config Config) (app *App, err error) {
	app = new(App)
	app.config = config
	return app, err
}

//NewRow create new rowEntry
func NewRow(rawString string, rex *regexp.Regexp) (row *RowEntry, err error) {
	row = new(RowEntry)
	row.Fields = make(map[string]string)
	row.Raw = rawString

	matches := rex.FindStringSubmatch(rawString)

	if len(matches) == 0 {
		return nil, errEmptyResult
	}

	for i, name := range rex.SubexpNames() {
		if len(name) > 0 {
			row.Fields[name] = matches[i]
		}
	}
	return row, err
}

//Start starts an app
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
		a.processBufferSync <- true
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

func (a *App) openReader() (err error) {
	a.ir, err = input.GetFileReader(a.config.InputDsn)
	return err
}

func (a *App) stop() {
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
	log.Println("[processBuffer] start")

	buffer := a.ir.FlushBuffer()

	<-a.processBufferSync

	log.Println("[processBuffer] buffer read done")

	byteReader := bytes.NewReader(buffer)
	bufReader := bufio.NewReader(byteReader)

	a.senderCollection.resetData()

	log.Println("[processBuffer] resetData")

	for {
		rawString, err = bufReader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			check(err)
		}

		logRow, err := NewRow(rawString, a.config.Rex)
		if err != nil && err == errEmptyResult {
			log.Println(err, rawString)
			continue
		}
		check(err)

		a.senderCollection.appendData(logRow)
		lastString = rawString
	}
	log.Println("[processBuffer] buffer append done")
	log.Println(lastString)

	go a.senderCollection.sendStats()

}
