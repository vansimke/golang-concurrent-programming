package alog

import "testing"

// 01-01
func TestNewInitializesShutdownChannels(t *testing.T) {
	alog := New(nil)
	if alog.shutdownCh == nil {
		t.Error("shutdownCh field not initialized")
	}

	if alog.shutdownCompleteCh == nil {
		t.Error("shutdownCompleteCh field not initialized")
	}
}

// 02-01
