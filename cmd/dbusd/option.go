package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
)

var (
	options struct {
		verbose            bool
		veryVerbose        bool
		configfile         string
		showversion        bool
		logfile            string
		debug              bool
		tick               int
		dryrun             bool
		lockfile           string
		diagnosticInterval int
	}
)

const (
	USAGE = `dbusd - Distributed Data Pipeline

Flags:
`
)

func parseFlags() {
	flag.BoolVar(&options.verbose, "v", false, "verbose")
	flag.BoolVar(&options.veryVerbose, "vv", false, "very verbose")
	flag.StringVar(&options.configfile, "conf", "etc/engine.cf", "main config file")
	flag.StringVar(&options.logfile, "log", "", "master log file path, default stdout")
	flag.StringVar(&options.lockfile, "lockfile", "var/dpiped.lock", "lockfile path")
	flag.BoolVar(&options.showversion, "version", false, "show version and exit")
	flag.BoolVar(&options.debug, "debug", false, "debug mode")
	flag.IntVar(&options.tick, "tick", 60*10, "tick interval in seconds to report sys stat")
	flag.BoolVar(&options.dryrun, "dryrun", false, "dry run")
	flag.Usage = showUsage
	flag.Parse()

	if options.veryVerbose {
		options.debug = true
	}
	if options.debug {
		options.verbose = true
	}

	if options.tick <= 0 {
		panic("tick must be possitive")
	}
}

func showUsage() {
	fmt.Fprint(os.Stderr, USAGE)
	flag.PrintDefaults()
}

func showVersionAndExit() {
	fmt.Fprintf(os.Stderr, "%s %s (build: %s)\n", os.Args[0], Version, BuildID)
	fmt.Fprintf(os.Stderr, "Built with %s %s for %s/%s\n",
		runtime.Compiler, runtime.Version(), runtime.GOOS, runtime.GOARCH)
	os.Exit(0)
}
