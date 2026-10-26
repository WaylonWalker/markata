package plugins

import (
	"io"
	stdlog "log"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/logging"
)

func TestSubprocessLineWriter_TruncatesOversizedLine(t *testing.T) {
	previousWriter, previousFlags, previousPrefix := stdlog.Writer(), stdlog.Flags(), stdlog.Prefix()
	t.Cleanup(func() {
		stdlog.SetOutput(previousWriter)
		stdlog.SetFlags(previousFlags)
		stdlog.SetPrefix(previousPrefix)
	})

	var logged []string
	logging.ConfigureStandardLogger(logging.Options{
		Writer: io.Discard,
		Observer: func(_ logging.Entry, message string) {
			logged = append(logged, message)
		},
	})
	writer := newSubprocessLineWriter(logging.Component("child"))
	if _, err := writer.Write([]byte(strings.Repeat("x", maxSubprocessLineBytes+100) + "\n")); err != nil {
		t.Fatal(err)
	}
	writer.Flush()
	if len(logged) != 1 {
		t.Fatalf("logged entries = %d, want 1", len(logged))
	}
	if len(logged[0]) > maxSubprocessLineBytes+len(" … [line truncated]") {
		t.Fatalf("logged line has %d bytes, want a bounded line", len(logged[0]))
	}
	if !strings.HasSuffix(logged[0], "… [line truncated]") {
		t.Fatalf("logged line does not report truncation: %q", logged[0][len(logged[0])-32:])
	}
}

func TestSubprocessTailWriterKeepsOnlyNewestBytes(t *testing.T) {
	writer := newSubprocessTailWriter(8)
	if _, err := writer.Write([]byte("first")); err != nil {
		t.Fatal(err)
	}
	if _, err := writer.Write([]byte("-last")); err != nil {
		t.Fatal(err)
	}
	if got := writer.String(); got != "[earlier subprocess output truncated]\nrst-last" {
		t.Fatalf("tail = %q, want truncated tail", got)
	}
}
