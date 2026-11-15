// Package logging sets up logging for an MCP server on stdio. stdout
// carries the MCP protocol, so every log line goes to stderr, where MCP
// clients collect it.
package logging

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"runtime/debug"
)

// New returns a text logger to stderr and makes it the default.
func New() *slog.Logger {
	return NewTo(os.Stderr)
}

// NewTo returns a text logger to w and makes it the default: for tests.
func NewTo(w io.Writer) *slog.Logger {
	logger := slog.New(slog.NewTextHandler(w, nil))
	slog.SetDefault(logger)
	return logger
}

// LogPanic logs a recovered panic with its stack. A runtime error's value
// is logged (it describes the bug, e.g. a nil map write); any other value
// only by its type, because it could carry the user's data.
func LogPanic(logger *slog.Logger, where string, r any) {
	value := fmt.Sprintf("%T", r)
	if rerr, ok := r.(runtime.Error); ok {
		value = rerr.Error()
	}
	logger.Error("recovered from a panic", "in", where, "panic", value, "stack", string(debug.Stack()))
}
