package cmd

import (
	"os"
	"testing"
)

func TestColorEnabledFor_NoColorEnv(t *testing.T) {
	originalNoColor := noColor
	originalForceColor := forceColor
	originalLogFormat := logFormat
	defer func() { noColor = originalNoColor }()
	defer func() { forceColor = originalForceColor }()
	defer func() { logFormat = originalLogFormat }()

	t.Setenv("NO_COLOR", "1")
	noColor = false
	forceColor = false
	logFormat = "auto"

	if colorEnabledFor(true) {
		t.Fatal("expected color to be disabled when NO_COLOR is set")
	}
}

func TestColorEnabledFor_WhitespaceNoColorEnv(t *testing.T) {
	oldNoColor, oldForceColor, oldLogFormat := noColor, forceColor, logFormat
	t.Cleanup(func() { noColor, forceColor, logFormat = oldNoColor, oldForceColor, oldLogFormat })
	t.Setenv("NO_COLOR", " ")
	t.Setenv("FORCE_COLOR", "")
	t.Setenv("TERM", "xterm-256color")
	noColor, forceColor, logFormat = false, false, "auto"

	if colorEnabledFor(true) {
		t.Fatal("expected any non-empty NO_COLOR value to disable automatic color")
	}
}

func TestColorEnabledFor_TermDumb(t *testing.T) {
	originalNoColor := noColor
	originalForceColor := forceColor
	originalLogFormat := logFormat
	defer func() { noColor = originalNoColor }()
	defer func() { forceColor = originalForceColor }()
	defer func() { logFormat = originalLogFormat }()

	t.Setenv("NO_COLOR", "")
	t.Setenv("TERM", "dumb")
	noColor = false
	forceColor = false
	logFormat = "auto"

	if colorEnabledFor(true) {
		t.Fatal("expected color to be disabled when TERM=dumb")
	}
}

func TestColorEnabledFor_NoColorFlag(t *testing.T) {
	originalNoColor := noColor
	originalForceColor := forceColor
	originalLogFormat := logFormat
	defer func() { noColor = originalNoColor }()
	defer func() { forceColor = originalForceColor }()
	defer func() { logFormat = originalLogFormat }()

	t.Setenv("NO_COLOR", "")
	t.Setenv("TERM", "xterm-256color")
	noColor = true
	forceColor = false
	logFormat = "auto"

	if colorEnabledFor(true) {
		t.Fatal("expected color to be disabled when --no-color is set")
	}
}

func TestColorEnabledFor_NonTTY(t *testing.T) {
	originalNoColor := noColor
	originalForceColor := forceColor
	originalLogFormat := logFormat
	defer func() { noColor = originalNoColor }()
	defer func() { forceColor = originalForceColor }()
	defer func() { logFormat = originalLogFormat }()

	t.Setenv("NO_COLOR", "")
	t.Setenv("TERM", "xterm-256color")
	noColor = false
	forceColor = false
	logFormat = "auto"

	if colorEnabledFor(false) {
		t.Fatal("expected color to be disabled on non-TTY output")
	}
}

func TestColorEnabledFor_InteractiveTTY(t *testing.T) {
	originalNoColor := noColor
	originalForceColor := forceColor
	originalLogFormat := logFormat
	defer func() { noColor = originalNoColor }()
	defer func() { forceColor = originalForceColor }()
	defer func() { logFormat = originalLogFormat }()

	t.Setenv("NO_COLOR", "")
	t.Setenv("TERM", "xterm-256color")
	noColor = false
	forceColor = false
	logFormat = "auto"

	if !colorEnabledFor(true) {
		t.Fatal("expected color to be enabled for interactive TTY output")
	}
}

func TestColorEnabledFor_ForceColor(t *testing.T) {
	originalNoColor := noColor
	originalForceColor := forceColor
	originalLogFormat := logFormat
	defer func() { noColor = originalNoColor }()
	defer func() { forceColor = originalForceColor }()
	defer func() { logFormat = originalLogFormat }()

	t.Setenv("NO_COLOR", "1")
	t.Setenv("TERM", "xterm-256color")
	noColor = false
	forceColor = true
	logFormat = "rich"

	if !colorEnabledFor(false) {
		t.Fatal("expected --color to enable color for non-TTY output")
	}
}

func TestColorEnabledFor_ForceColorEnv(t *testing.T) {
	oldNoColor, oldForceColor, oldLogFormat := noColor, forceColor, logFormat
	t.Cleanup(func() { noColor, forceColor, logFormat = oldNoColor, oldForceColor, oldLogFormat })
	t.Setenv("NO_COLOR", "1")
	t.Setenv("FORCE_COLOR", "1")
	t.Setenv("TERM", "dumb")
	noColor, forceColor, logFormat = false, false, "auto"

	if !colorEnabledFor(false) {
		t.Fatal("expected FORCE_COLOR to override automatic detection")
	}
}

func TestColorEnabledFor_PlainLogFormat(t *testing.T) {
	originalNoColor := noColor
	originalForceColor := forceColor
	originalLogFormat := logFormat
	defer func() { noColor = originalNoColor }()
	defer func() { forceColor = originalForceColor }()
	defer func() { logFormat = originalLogFormat }()

	t.Setenv("NO_COLOR", "")
	t.Setenv("TERM", "xterm-256color")
	noColor = false
	forceColor = true
	logFormat = "plain"

	if colorEnabledFor(true) {
		t.Fatal("expected plain log format to disable color")
	}
}

// A character device is not necessarily an interactive terminal.
func TestFileLikeTerminalRejectsNullDevice(t *testing.T) {
	file, err := os.Open(os.DevNull)
	if err != nil {
		t.Skipf("null device unavailable: %v", err)
	}
	defer file.Close()
	if fileLikeTerminal(file) {
		t.Fatal("null device must not activate terminal UI")
	}
}
