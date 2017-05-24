package main

import (
	"flag"
	"log"
	"os"

	"github.com/blackbass1988/access_logs_stats/pkg"
	_ "github.com/blackbass1988/access_logs_stats/pkg/output/console"
	_ "github.com/blackbass1988/access_logs_stats/pkg/output/zabbix"
	prof "github.com/blackbass1988/yet_another_pprof_wrapper"
)

var (
	version   = "0.8.0"
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
	log.Printf("AccessLogsStats ver.%s@%s (git %s %s)", version, buildTime, branch, commit)
}

func main() {
	var (
		fileconfig       string
		heapProfile      string
		cpuProfile       string
		exitAfterOneTick bool
	)

	printHello()

	flag.StringVar(&fileconfig, "c", "", "config path")
	flag.StringVar(&heapProfile, "heapprofile", "", "enable heap profiling")
	flag.StringVar(&cpuProfile, "cpuprofile", "", "Write the cpu heapProfile to `filename`")
	flag.BoolVar(&exitAfterOneTick, "one", false, "make one tick end exit")
	flag.Parse()

	if cpuProfile != "" {
		go prof.ProfileCpu(cpuProfile)
	}

	if heapProfile != "" {
		go prof.ProfileMemory(heapProfile)
	}

	if fileconfig == "" {
		log.Print("ERROR config not set")
		flag.PrintDefaults()
		os.Exit(2)
	}

	config, err := pkg.NewConfig(fileconfig)
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
