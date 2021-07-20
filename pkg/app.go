package pkg

import (
	"bufio"
	"errors"
	"github.com/blackbass1988/access_logs_stats/pkg/template"
	"log"
	"os"
	"time"

	"github.com/blackbass1988/access_logs_stats/pkg/input"
	"github.com/blackbass1988/access_logs_stats/pkg/re"
)

var (
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
	config *Config
	buffer []byte

	senderCollection *SenderCollection
	ir               input.BufferedReader

	fileReader *bufio.Reader
}

//NewApp creates new parser
func NewApp(config *Config) (app *App, err error) {
	app = new(App)

	err, tmpl := template.NewTempate(config.InputDsn)

	if err != nil {
		return nil, err
	}

	err, config.InputDsn = tmpl.ProcessTemplate(config.TemplateVars)

	if err != nil {
		return nil, err
	}

	app.config = config
	return app, err
}

//NewRow create new rowEntry
func NewRow(rawString string, rex re.RegExp) (row *RowEntry, err error) {
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
	checkOrFail(err)

	defer func() {
		a.ir.Close()
	}()

	lineChannel := make(chan string)
	go a.ir.ReadToChannel(lineChannel)

	if a.config.ExitAfterOneTick {
		a.appendLine(lineChannel)
		a.senderCollection.sendStats()
	} else {
		go a.appendLine(lineChannel)
		//read to buffer in background

		for {
			select {
			case <-tick:
				go a.senderCollection.sendStats()
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
	a.buffer = []byte{}
	a.senderCollection = NewSenderCollection(a.config)
}

func (a *App) appendLine(linesChannel <-chan string) {
	var (
		rawString string
		err       error
		logRow    *RowEntry
		more      bool
	)
	for {
		rawString, more = <-linesChannel

		if !more {
			break
		}

		logRow, err = NewRow(rawString, a.config.Rex)

		if err != nil && err == errEmptyResult {
			continue
		}
		checkOrFail(err)

		a.senderCollection.appendData(logRow)
	}
}
