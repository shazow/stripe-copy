package main

import (
	"fmt"
	"os"

	"github.com/alexcesaro/log"
	"github.com/alexcesaro/log/golog"
	"github.com/jessevdk/go-flags"
)

// version is overridden during the build process
var version string = "dev"

// cmdOptions contains the flag options
type cmdOptions struct {
	Verbose   []bool `short:"v" long:"verbose" description:"Show verbose logging."`
	Pretend   bool   `short:"p" long:"pretend" description:"Do everything read-only, skip writes."`
	StopAfter int    `long:"stopAfter" description:"Stop after this many write operations."`
	Version   bool   `long:"version"`
}

// logLevels corresponds to the number of Verbose flags set
var logLevels = []log.Level{
	log.Warning,
	log.Info,
	log.Debug,
}

func fail(code int, format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(code)
}

func main() {
	options := cmdOptions{}
	parser := flags.NewParser(&options, flags.Default)
	p, err := parser.Parse()
	if err != nil {
		if p == nil {
			fmt.Print(err)
		}
		os.Exit(1)
		return
	}

	if options.Version {
		fmt.Printf(version)
		os.Exit(0)
	}

	// Figure out the log level
	numVerbose := len(options.Verbose)
	if numVerbose > len(logLevels) {
		numVerbose = len(logLevels) - 1
	}

	logLevel := logLevels[numVerbose]
	setLogger(golog.New(os.Stderr, logLevel))

	api, err := newStripeAPI(os.Getenv(envStripeSource), os.Getenv(envStripeTarget))
	if err != nil {
		fail(1, "Failed to initialize API: %s\n", err)
	}

	if options.Pretend {
		logger.Info("Running in pretend mode. Write operations will be skipped.")
		api.pretend = true
	}

	if options.StopAfter > 0 {
		logger.Debugf("Will stop after %d write operations.", options.StopAfter)
		api.stopAfter = options.StopAfter
	}

	err = api.SyncPlans()
	if err != nil {
		fail(2, "Failed to sync plans: %s\n", err)

	}

	fmt.Fprintln(os.Stderr, "Done.")
	os.Exit(0)
}
