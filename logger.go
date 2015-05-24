package main

import (
	"bytes"

	"github.com/alexcesaro/log"
	"github.com/alexcesaro/log/golog"
)

var logger *golog.Logger

// SetLogger overrides the default logger for the package.
func SetLogger(l *golog.Logger) {
	logger = l
}

func init() {
	// Set a default null logger
	var b bytes.Buffer
	SetLogger(golog.New(&b, log.Debug))
}
