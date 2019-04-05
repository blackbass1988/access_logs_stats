package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/blackbass1988/access_logs_stats/pkg"
	_ "github.com/blackbass1988/access_logs_stats/pkg/output/console"
	_ "github.com/blackbass1988/access_logs_stats/pkg/output/zabbix"
	prof "github.com/blackbass1988/yet_another_pprof_wrapper"
)

var (
	version   = "0.10.0"
	buildTime = "unknown"
	commit    = "unknown"
	branch    = "unknown"
)

func init() {
	if version == "" {
		version = "unknown"
	}
	if commit == "" {
		commit = "unknown"
	}
	if branch == "" {
		branch = "unknown"
	}
}

func printHello() {
	fmt.Printf("AccessLogsStats ver.%s@%s (git %s %s)\n", version, buildTime, branch, commit)
}

func main() {
	var (
		fileconfig       string
		heapProfile      string
		cpuProfile       string
		exitAfterOneTick bool
		showVersion      bool
		templateVars     templateVarsArray
		templateVarsMap  map[string]string
	)

	printHello()

	flag.BoolVar(&showVersion, "version", false, "show current version")
	flag.StringVar(&fileconfig, "c", "", "config path")
	flag.StringVar(&heapProfile, "heapprofile", "", "enable heap profiling")
	flag.StringVar(&cpuProfile, "cpuprofile", "", "Write the cpu heapProfile to `filename`")
	flag.BoolVar(&exitAfterOneTick, "one", false, "make one tick end exit")
	flag.Var(&templateVars,
		"template-var",
		`Extra variables to set into output template.
You can pass many variables.
For example: -template-var key=value -template-var foo=bar`)
	flag.Parse()

	if showVersion {
		os.Exit(0)
	}

	if len(templateVars) > 0 {
		templateVarsMap = arrayToMap(templateVars)
	}

	if cpuProfile != "" {
		cWriter, err := os.Create(cpuProfile)
		if err != nil {
			panic(err)
		}
		go prof.ProfileCpu(cWriter)
	}

	if heapProfile != "" {
		mWriter, err := os.Create(heapProfile)
		if err != nil {
			panic(err)
		}
		go prof.ProfileMemory(mWriter, 10*time.Second, true)
	}

	if fileconfig == "" {
		log.Print("ERROR config not set")
		flag.PrintDefaults()
		os.Exit(2)
	}

	config, err := pkg.NewConfig(fileconfig, templateVarsMap)
	if err != nil {
		log.Fatal(err)
	}
	config.ExitAfterOneTick = exitAfterOneTick

	app, err := pkg.NewApp(config)
	if err != nil {
		log.Fatal(err)
	}

	app.Start()
}

type templateVarsArray []string

func (t *templateVarsArray) String() string {
	return "[" + strings.Join(*t, ",") + "]"
}

func (t *templateVarsArray) Set(value string) error {
	*t = append(*t, value)
	return nil
}

func arrayToMap(arr []string) map[string]string {
	m := make(map[string]string)

	for _, i := range arr {
		splited := strings.Split(i, "=")
		m[splited[0]] = splited[1]
	}

	return m
}
