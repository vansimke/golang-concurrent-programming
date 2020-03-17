package alog

import (
	"bytes"
	"errors"
	"reflect"
	"regexp"
	"testing"
	"time"
)

const messageTimestampPattern = `\[\d{4}-\d{2}-\d{2}\ \d{2}:\d{2}:\d{2}] - `

// 01-01 (Task 01, Test 01)
func TestMessageChannelModule1(t *testing.T) {
	alog := New(nil)
	if alog.msgCh == nil {
		t.Fatal("msgCh field not initialized. Should have type 'chan string' but it is currently nil")
	}
}

// 02-01
func TestErrorChannelModule1(t *testing.T) {
	alog := New(nil)
	if alog.errorCh == nil {
		t.Fatal("errorCh field not initialized. Should have type 'chan string' but it is currently nil")
	}
}

// 03-01
func TestMessageChannelMethodModule1(t *testing.T) {
	alog := New(nil)
	if alog.MessageChannel() != alog.msgCh {
		t.Fatal("MessageChannel method does not return the msgCh field")
	}
	messageChannelDir := reflect.ValueOf(alog.MessageChannel()).Type().ChanDir()
	if messageChannelDir != reflect.SendDir {
		t.Fatal("MessageChannel does not return send-only channel")
	}
}

// 04-01
func TestErrorChannelMethodModule1(t *testing.T) {
	alog := New(nil)
	if alog.ErrorChannel() != alog.errorCh {
		t.Fatal("ErrorChannel method does not return the errorCh field")
	}
	errorChannelDir := reflect.ValueOf(alog.ErrorChannel()).Type().ChanDir()
	if errorChannelDir != reflect.RecvDir {
		t.Fatal("ErrorChannel does not return receive-only channel")
	}
}

// 05-01
func TestWritesToWriterModule1(t *testing.T) {
	b := bytes.NewBuffer([]byte{})
	alog := New(b)
	alog.write("test", nil)

	written := b.String()
	if written == "" {
		t.Fatal("Nothing written to log")
	}
	if !regexp.MustCompile(messageTimestampPattern + "test\n$").Match([]byte(written)) {
		t.Error("Properly formatted string not written to log. Did you pass the message to 'formatMessage'?")
	}
}

// 06-01

type errorWriter struct {
	b *bytes.Buffer
}

func (ew errorWriter) Write(data []byte) (int, error) {
	ew.b.Write(data)
	return 0, errors.New("error")
}
func TestWriteSendsErrorsToErrorChannelModule1(t *testing.T) {
	alog := New(&errorWriter{bytes.NewBuffer([]byte{})})
	alog.errorCh = make(chan error, 1)
	alog.write("test", nil)
	if (<-alog.errorCh).Error() != "error" {
		t.Fatal("Did not receive destination writer's error on errorCh")
	}
}

// 07-01
type sleepingWriter struct {
	b *bytes.Buffer
}

func (sw sleepingWriter) Write(data []byte) (int, error) {
	sw.b.Write(data)
	time.Sleep(1 * time.Second)
	sw.b.WriteString("write complete")
	return 0, nil
}

func TestStartHandlesMessagesModule1(t *testing.T) {
	b := bytes.NewBuffer([]byte{})
	alog := New(sleepingWriter{b})
	go alog.Start()
	alog.msgCh <- "test message"
	time.Sleep(100 * time.Millisecond)
	written := b.Bytes()
	if !regexp.MustCompile(messageTimestampPattern + "test message\n$").Match(written) {
		t.Error("Message not written to logger's destination")
	}
	if alog.m != nil {
		alog.m.Unlock()
	}
	alog.msgCh <- "second message"
	time.Sleep(100 * time.Millisecond)
	written = b.Bytes()
	if !regexp.MustCompile(messageTimestampPattern + "test message\n" + messageTimestampPattern + "second message\n").Match(written) {
		t.Error("write method not called as a goroutine")
	}
}

// 08-01

func TestMutexIsIntializedModule1(t *testing.T) {
	alog := New(nil)
	if alog.m == nil {
		t.Fatal("Alog's mutex field 'm' not initialized")
	}
}

// 08-02
func TestWriteSendsWriteRequestsSequentiallyModule1(t *testing.T) {
	b := bytes.NewBuffer([]byte{})
	alog := New(sleepingWriter{b})
	go alog.write("test message", nil)
	time.Sleep(100 * time.Millisecond)
	go alog.write("second message", nil)
	time.Sleep(1500 * time.Millisecond)
	written := b.Bytes()
	if !regexp.MustCompile(messageTimestampPattern + "test message\nwrite complete" + messageTimestampPattern + "second message\n").Match(written) {
		t.Error("Mutex not protecting Alog.dest#Write from concurrent calls")
	}
}

// 09-01
type panickingWriter struct {
	b *bytes.Buffer
}

func (pw panickingWriter) Write(data []byte) (int, error) {
	pw.b.Write(data)
	panic("panicking!")
}
func TestWriteSendsWriteRequestsSequentiallyWhenPanickingModule1(t *testing.T) {
	b := bytes.NewBuffer([]byte{})
	alog := New(panickingWriter{b})
	go func() {
		defer func() {
			recover()
		}()
		alog.write("test message", nil)
	}()
	time.Sleep(100 * time.Millisecond)
	go func() {
		defer func() {
			recover()
		}()
		alog.write("second message", nil)
	}()
	time.Sleep(1500 * time.Millisecond)
	written := b.Bytes()
	if !regexp.MustCompile(messageTimestampPattern + "test message\n" + messageTimestampPattern + "second message\n").Match(written) {
		t.Error("Mutex not unlocked when panicking", string(written))
	}
}

// 10-01

func TestWriteSendsErrorsAsynchronouslyModule1(t *testing.T) {
	b := bytes.NewBuffer([]byte{})
	alog := New(&errorWriter{b})
	go alog.write("first", nil)
	time.Sleep(100 * time.Millisecond)
	go alog.write("second", nil)
	time.Sleep(100 * time.Millisecond)
	written := b.Bytes()
	if !regexp.MustCompile(`.*first.*\n.*second.*`).Match(written) {
		t.Fatal("Error messages not sent to error channel asynchronously", string(written))
	}
}
