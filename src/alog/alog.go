// Package alog provides a simple asynchronous logger that will write to provided io.Writers without blocking calling
// goroutines.
package alog

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Alog is a type that defines a logger. It starts as a simple, synchronous pass-through to an underlying writer
// but will evolve into a fully asynchronous logging provider
type Alog struct {
	dest               io.Writer
	m                  *sync.Mutex
	msgCh              chan string
	errorCh            chan error
	shutdownCh         chan struct{}
	shutdownCompleteCh chan struct{}
}

// New creates a new Alog object that writes to the provided io.Writer.
// If nil is provided the output will be directed to os.Stdout
func New(w io.Writer) *Alog {
	if w == nil {
		w = os.Stdout
	}
	return &Alog{
		dest:               w,
		m:                  &sync.Mutex{},
		msgCh:              make(chan string),
		shutdownCh:         make(chan struct{}),
		shutdownCompleteCh: make(chan struct{}),
		errorCh:            make(chan error),
	}
}

// Start begins the message loop for the asychronous logger. It should be initiated as a goroutine to prevent
// the caller from being blocked.
func (al Alog) Start() {
	wg := &sync.WaitGroup{}
	for {
		select {
		case msg := <-al.msgCh:
			wg.Add(1)
			go al.write(msg, wg)
		case <-al.shutdownCh:
			wg.Wait()
			al.shutdown()
		}
	}
}

func (al Alog) write(msg string, wg *sync.WaitGroup) {
	defer wg.Done()
	logEntry := fmt.Sprintf("[%v] - %v", time.Now().Format("2006-01-02"), msg)
	al.m.Lock()
	defer al.m.Unlock()
	_, err := al.dest.Write([]byte(logEntry))
	if err != nil {
		go func(err error) {
			al.errorCh <- err
		}(err)
	}
}

func (al Alog) shutdown() {
	close(al.msgCh)
	close(al.errorCh)
	close(al.shutdownCh)
	al.shutdownCompleteCh <- struct{}{}
}

// MessageChannel returns a channel that accepts messages that should be written to the log.
func (al Alog) MessageChannel() chan<- string {
	return al.msgCh
}

// ErrorChannel returns a channel that will be populated when an error is raised during a write operation.
// This channel should always be monitored in some way to prevent deadlock goroutines from being generated
// when errors occur.
func (al Alog) ErrorChannel() <-chan error {
	return al.errorCh
}

// Stop shuts down the logger. It will wait for all pending messages to be written and then return.
// The logger will no longer function after this method has been called.
func (al Alog) Stop() {
	al.shutdownCh <- struct{}{}

	<-al.shutdownCompleteCh
}
