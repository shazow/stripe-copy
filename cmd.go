package main

import (
	"fmt"
	"os"

	"github.com/alexcesaro/log"
	"github.com/alexcesaro/log/golog"
	"github.com/jessevdk/go-flags"
)

// Options contains the flag options
type Options struct {
	Verbose []bool `short:"v" long:"verbose" description:"Show verbose logging."`
	Pretend bool   `short:"p" long:"pretend" description:"Do everything read-only, skip writes."`
}

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
	options := Options{}
	parser := flags.NewParser(&options, flags.Default)
	p, err := parser.Parse()
	if err != nil {
		if p == nil {
			fmt.Print(err)
		}
		os.Exit(1)
		return
	}

	// Figure out the log level
	numVerbose := len(options.Verbose)
	if numVerbose > len(logLevels) {
		numVerbose = len(logLevels) - 1
	}

	logLevel := logLevels[numVerbose]
	setLogger(golog.New(os.Stderr, logLevel))

	api := newStripeAPI(os.Getenv(EnvStripeFrom), os.Getenv(EnvStripeTo))
	err = api.SyncPlans()
	if err != nil {
		fail(1, "Failed to sync plans:", err)

	}

	fmt.Fprintln(os.Stderr, "Done.")
	os.Exit(0)
}
