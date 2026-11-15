package logging

import (
	"bytes"
	"strings"
	"testing"
)

// Guardrail: a panic value that isn't a runtime error is logged by type
// only, because it could carry the user's data.
func TestLogPanicHidesTheValue(t *testing.T) {
	var buf bytes.Buffer
	logger := NewTo(&buf)
	LogPanic(logger, "tool", "Subject: private")
	if strings.Contains(buf.String(), "private") || !strings.Contains(buf.String(), "panic=string") {
		t.Errorf("log: %s", buf.String())
	}

	buf.Reset()
	func() {
		defer func() { LogPanic(logger, "tool", recover()) }()
		var s []int
		_ = s[len(buf.String())] // a runtime error: its text describes the bug
	}()
	if !strings.Contains(buf.String(), "index out of range") {
		t.Errorf("runtime error not logged: %s", buf.String())
	}
}
