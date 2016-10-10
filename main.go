package main

import (
	"flag"
	"github.com/blackbass1988/access_logs_stats/core"
	"log"
	"os"

	"io"
	"runtime"
	"runtime/pprof"
	"time"
)

const (
	PROG_NAME = "AccessLogsStats"
	VERSION   = "0.5.0"
)

func PrintHello() {
	log.Printf("%s ver.%s", PROG_NAME, VERSION)
}

func main() {
	var (
		fileconfig string
		profile    bool
	)

	PrintHello()

	flag.StringVar(&fileconfig, "c", "", "config path")
	flag.BoolVar(&profile, "p", false, "enable profiling")
	flag.Parse()

	if fileconfig == "" {
		log.Print("ERROR config not set")
		flag.PrintDefaults()
		os.Exit(2)
	}

	config, err := core.NewConfig(fileconfig)
	if err != nil {
		log.Fatal(err)
	}

	app, err := core.NewApp(config)

	if profile {
		go codeProfile()
	}

	app.Start()
}

func codeProfile() {
	m := &runtime.MemStats{}
	tick1m := time.Tick(1 * time.Minute)
	tick5s := time.Tick(5 * time.Second)

	f_cpu_profiling, err := os.Create("profile.prof")
	if err != nil {
		panic(err)
	}
	pprof.StartCPUProfile(f_cpu_profiling)
	defer func() {
		pprof.StopCPUProfile()
	}()

	for {
		select {
		case <-tick5s:
			runtime.ReadMemStats(m)
			log.Println("")
			log.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
			log.Printf("~ Goroutines count %d\n", runtime.NumGoroutine())
			log.Printf("~ Alloc %dKB\n", m.Alloc/1024)
			log.Printf("~ TotalAlloc %dKB\n", m.TotalAlloc/1024)
			log.Printf("~ Sys (sum of XxxSys below) %dKB\n", m.Sys/1024)
			log.Printf("~ Lookups (number of pointer lookups) %d\n", m.Lookups)
			log.Printf("~ Mallocs %d\n", m.Mallocs)
			log.Printf("~ Frees %d\n", m.Frees)

			log.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")

			log.Printf("~ HeapAlloc %dKB\n", m.HeapAlloc/1024)
			log.Printf("~ HeapSys %dKB\n", m.HeapSys/1024)
			log.Printf("~ HeapIdle %dKB\n", m.HeapIdle/1024)
			log.Printf("~ HeapInuse %dKB\n", m.HeapInuse/1024)
			log.Printf("~ HeapReleased %dKB\n", m.HeapReleased/1024)
			log.Printf("~ HeapObjects %d\n", m.HeapObjects)

			log.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")

			log.Printf("~ NextGC %d\n", m.NextGC)
			log.Printf("~ LastGC %v\n", time.Unix(0, int64(m.LastGC)))
			log.Printf("~ PauseTotalNs %d\n", m.PauseTotalNs)
			log.Printf("~ NumGC %d\n", m.PauseTotalNs)
			log.Printf("~ GCCPUFraction %f\n", m.GCCPUFraction)
			log.Printf("~ EnableGC %v\n", m.EnableGC)
			log.Printf("~ DebugGC %v\n", m.DebugGC)
			log.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
			log.Println("")

		case <-tick1m:
			var f_heap_profiling io.Writer
			f_heap_profiling, _ = os.Create("profile_heap.prof")
			pprof.WriteHeapProfile(f_heap_profiling)
			log.Println("~ head saved")
		}
	}
}
