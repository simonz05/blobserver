package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"

	"git.tideland.biz/goas/monitoring"
	"github.com/simonz05/util/log"
	"github.com/simonz05/blobserver"
	"github.com/simonz05/blobserver/server"
	"github.com/simonz05/blobserver/config"
)

var (
	help       = flag.Bool("h", false, "show help text")
	laddr      = flag.String("http", ":6064", "set bind address for the HTTP server")
	version    = flag.Bool("version", false, "show version number and exit")
	configFilename = flag.String("config", "config.toml", "config file path")
	cpuprofile = flag.String("debug.cpuprofile", "", "write cpu profile to file")
)

var Version = "0.1.0"

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()
	log.Println("start blobserver service …")

	if *version {
		fmt.Fprintln(os.Stdout, Version)
		return
	}

	if *help {
		flag.Usage()
		os.Exit(1)
	}

	conf, err := config.FromFile(*configFilename)

	if err != nil {
		log.Fatal(err)
	}

	if conf.Listen == "" && *laddr == ""{
		log.Fatal("Listen address required")
	} else if conf.Listen == "" {
		conf.Listen = *laddr
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	storage, err := blobserver.CreateStorage(conf)

	if err != nil {
		log.Fatalf("error instantiating storage for type %s: %v",
			conf.StorageType(), err)
	}

	err = server.ListenAndServe(*laddr, storage)

	if err != nil {
		log.Errorln(err)
	}

	monitoring.MeasuringPointsPrintAll()
}