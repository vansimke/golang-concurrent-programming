// Package alog provides a simple asynchronous logger that will write to provided io.Writers without blocking calling
// goroutines.
package alog

import (
	"io"
	"os"
	"strings"
)

// Alog is a type that defines a logger. It starts as a simple, synchronous pass-through to an underlying writer
// but will evolve into a fully asynchronous logging provider
type Alog struct {
	dest io.Writer
}

// New creates a new Alog object that writes to the provided io.Writer.
// If nil is provided the output will be directed to os.Stdout
func New(w io.Writer) *Alog {
	if w == nil {
		w = os.Stdout
	}
	return &Alog{
		dest: w,
	}
}

func (a Alog) Write(msg string) (int, error) {
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	return a.dest.Write([]byte(msg))
}
