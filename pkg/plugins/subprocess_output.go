package plugins

import (
	"bytes"
	"strings"
	"sync"

	"github.com/WaylonWalker/markata/pkg/logging"
)

const maxSubprocessLineBytes = 16 * 1024

type subprocessLineWriter struct {
	mu       sync.Mutex
	logger   logging.Logger
	line     []byte
	dropping bool
}

func newSubprocessLineWriter(logger logging.Logger) *subprocessLineWriter {
	return &subprocessLineWriter{logger: logger, line: make([]byte, 0, 256)}
}

func (w *subprocessLineWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	for _, b := range p {
		if b == '\n' {
			w.flushLineLocked()
			continue
		}
		if w.dropping {
			continue
		}
		if len(w.line) == maxSubprocessLineBytes {
			w.line = append(w.line, []byte(" … [line truncated]")...)
			w.dropping = true
			continue
		}
		w.line = append(w.line, b)
	}
	return len(p), nil
}

func (w *subprocessLineWriter) Flush() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.flushLineLocked()
}

func (w *subprocessLineWriter) flushLineLocked() {
	if len(w.line) > 0 {
		w.logger.Printf("%s", strings.TrimSpace(string(w.line)))
	}
	w.line = w.line[:0]
	w.dropping = false
}

type subprocessTailWriter struct {
	mu        sync.Mutex
	buf       []byte
	max       int
	truncated bool
}

func newSubprocessTailWriter(maxBytes int) *subprocessTailWriter {
	if maxBytes < 1 {
		maxBytes = 1
	}
	return &subprocessTailWriter{max: maxBytes}
}

func (w *subprocessTailWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if len(p) >= w.max {
		w.truncated = w.truncated || len(w.buf) > 0 || len(p) > w.max
		w.buf = append(w.buf[:0], p[len(p)-w.max:]...)
		return len(p), nil
	}
	if excess := len(w.buf) + len(p) - w.max; excess > 0 {
		copy(w.buf, w.buf[excess:])
		w.buf = w.buf[:len(w.buf)-excess]
		w.truncated = true
	}
	w.buf = append(w.buf, p...)
	return len(p), nil
}

func (w *subprocessTailWriter) String() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	output := string(bytes.Clone(w.buf))
	if w.truncated {
		return "[earlier subprocess output truncated]\n" + output
	}
	return output
}

// logSubprocessOutput forwards captured child-process output through the
// configured logger so interactive Serve can keep ownership of the terminal.
func logSubprocessOutput(logger logging.Logger, output string) {
	for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
		if line = strings.TrimSpace(line); line != "" {
			logger.Printf("%s", line)
		}
	}
}
