package main

import (
	"bytes"

	"github.com/alexcesaro/log"
	"github.com/alexcesaro/log/golog"
)

var logger *golog.Logger

// SetLogger overrides the default logger for the package.
func setLogger(l *golog.Logger) {
	logger = l
}

func init() {
	// Set a default null logger
	var b bytes.Buffer
	setLogger(golog.New(&b, log.Debug))
}
