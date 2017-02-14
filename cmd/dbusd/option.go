package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/funkygao/dbus"
)

var (
	options struct {
		debug  bool
		dryrun bool

		configfile    string
		showversion   bool
		visualizeFile string
		lockfile      string

		logfile  string
		loglevel string
	}
)

const (
	USAGE = `dbusd - Distributed DataBus

Flags:
`
)

func parseFlags() {
	flag.StringVar(&options.configfile, "conf", "", "main config file")
	flag.StringVar(&options.logfile, "logfile", "", "master log file path, default stdout")
	flag.StringVar(&options.loglevel, "loglevel", "trace", "log level")
	flag.BoolVar(&options.showversion, "version", false, "show version and exit")
	flag.BoolVar(&options.debug, "debug", false, "debug mode")
	flag.BoolVar(&options.dryrun, "dryrun", false, "dry run")
	flag.StringVar(&options.visualizeFile, "dump", "", "visualize the pipleline to a png file. graphviz must be installed")
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, USAGE)
		flag.PrintDefaults()
	}
	flag.Parse()
}

func showVersionAndExit() {
	fmt.Fprintf(os.Stderr, "%s %s (build: %s)\n", os.Args[0], dbus.Version, dbus.BuildID)
	fmt.Fprintf(os.Stderr, "Built with %s %s for %s/%s\n",
		runtime.Compiler, runtime.Version(), runtime.GOOS, runtime.GOARCH)
	os.Exit(0)
}
