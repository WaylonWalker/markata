package logging

import (
	"bytes"
	"fmt"
	"io"
	stdlog "log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/WaylonWalker/markata/pkg/palettes"
)

const (
	ansiBlue   = "\033[34m"
	ansiCyan   = "\033[36m"
	ansiDim    = "\033[2m"
	ansiGreen  = "\033[32m"
	ansiRed    = "\033[31m"
	ansiReset  = "\033[0m"
	ansiYellow = "\033[33m"

	metaPrefix = "\x1emarkata|"
	metaSep    = "\x1f"
)

type Format string

const (
	FormatAuto  Format = "auto"
	FormatPlain Format = "plain"
	FormatRich  Format = "rich"
)

type Theme struct {
	Timestamp  string
	Component  string
	Success    string
	Warning    string
	Error      string
	PhaseColor map[string]string
}

type Options struct {
	Writer io.Writer
	// Observer receives decoded entries before presentation. It must not log
	// through the standard logger, which would recursively call the writer.
	Observer   func(Entry, string)
	Format     Format
	ForceColor bool
	NoColor    bool
	IsTTY      bool
	Theme      Theme
}

type Writer struct {
	mu       sync.Mutex
	out      io.Writer
	format   Format
	color    bool
	theme    Theme
	buf      bytes.Buffer
	observer func(Entry, string)
}

type Entry struct {
	Component string
	Phase     string
	Level     string
}

type Logger struct {
	entry Entry
}

var (
	bracketPrefixPattern = regexp.MustCompile(`^\[([^\]]+)\]\s*(.*)$`)
	colonPrefixPattern   = regexp.MustCompile(`^([a-z][a-z0-9_/-]*):\s*(.*)$`)
	excludedPrefixes     = map[string]struct{}{
		"warning": {},
		"error":   {},
		"info":    {},
		"debug":   {},
	}
)

func ParseFormat(raw string) (Format, error) {
	switch Format(strings.TrimSpace(strings.ToLower(raw))) {
	case "", FormatAuto:
		return FormatAuto, nil
	case FormatPlain:
		return FormatPlain, nil
	case FormatRich:
		return FormatRich, nil
	default:
		return "", fmt.Errorf("invalid log format %q (expected auto, plain, or rich)", raw)
	}
}

func ConfigureStandardLogger(opts Options) {
	writer := NewWriter(opts)
	stdlog.SetFlags(0)
	stdlog.SetPrefix("")
	stdlog.SetOutput(writer)
}

func NewWriter(opts Options) *Writer {
	format := opts.Format
	if format == "" {
		format = FormatAuto
	}

	writer := opts.Writer
	if writer == nil {
		writer = os.Stderr
	}

	theme := opts.Theme
	if theme.Component == "" {
		theme = DefaultTheme()
	}
	color := allowColor(opts)

	return &Writer{
		out:      writer,
		observer: opts.Observer,
		format:   resolveFormat(format, opts.IsTTY || color),
		color:    color,
		theme:    theme,
	}
}

func DefaultTheme() Theme {
	return Theme{
		Timestamp: ansiDim,
		Component: ansiCyan,
		Success:   ansiGreen,
		Warning:   ansiYellow,
		Error:     ansiRed,
		PhaseColor: map[string]string{
			"configure": ansiCyan,
			"validate":  ansiYellow,
			"glob":      ansiBlue,
			"load":      ansiCyan,
			"transform": ansiBlue,
			"render":    ansiCyan,
			"collect":   ansiYellow,
			"write":     ansiBlue,
			"cleanup":   ansiDim,
		},
	}
}

func ThemeFromPalette(palette *palettes.Palette) Theme {
	if palette == nil {
		return DefaultTheme()
	}

	resolve := func(names ...string) string {
		for _, name := range names {
			if hex := palette.Resolve(name); hex != "" {
				return hex
			}
		}
		return ""
	}

	theme := DefaultTheme()
	if value := resolve("text-muted", "border", "text-secondary", "text"); value != "" {
		theme.Timestamp = value
	}
	if value := resolve("primary", "info", "text-primary", "text"); value != "" {
		theme.Component = value
	}
	if value := resolve("success", "primary", "secondary"); value != "" {
		theme.Success = value
	}
	if value := resolve("warning", "secondary", "primary"); value != "" {
		theme.Warning = value
	}
	if value := resolve("error", "primary", "secondary"); value != "" {
		theme.Error = value
	}

	theme.PhaseColor = map[string]string{
		"configure": resolve("primary", "info", "text-primary"),
		"validate":  resolve("warning", "secondary", "primary"),
		"glob":      resolve("secondary", "info", "primary"),
		"load":      resolve("info", "primary", "secondary"),
		"transform": resolve("tertiary", "primary", "secondary"),
		"render":    resolve("primary", "tertiary", "info"),
		"collect":   resolve("success", "secondary", "primary"),
		"write":     resolve("warning", "primary", "secondary"),
		"cleanup":   resolve("text-muted", "border", "text-secondary"),
	}

	defaults := DefaultTheme().PhaseColor
	for phase, value := range theme.PhaseColor {
		if value == "" {
			theme.PhaseColor[phase] = defaults[phase]
		}
	}

	return theme
}

func Component(name string) Logger {
	return Logger{entry: Entry{Component: name}}
}

func (l Logger) Phase(phase string) Logger {
	l.entry.Phase = phase
	return l
}

func (l Logger) Level(level string) Logger {
	l.entry.Level = level
	return l
}

func (l Logger) Printf(format string, args ...any) {
	stdlog.Print(encodeEntry(l.entry, fmt.Sprintf(format, args...)))
}

func (l Logger) Infof(format string, args ...any) {
	l.Level("info").Printf(format, args...)
}

func (l Logger) Warnf(format string, args ...any) {
	l.Level("warning").Printf(format, args...)
}

func (l Logger) Errorf(format string, args ...any) {
	l.Level("error").Printf(format, args...)
}

func (w *Writer) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	written, err := w.buf.Write(p)
	if err != nil {
		return 0, err
	}

	for {
		line, readErr := w.buf.ReadString('\n')
		if readErr != nil {
			remaining := line
			w.buf.Reset()
			_, _ = w.buf.WriteString(remaining)
			break
		}
		line = strings.TrimSuffix(line, "\n")
		if w.observer != nil {
			entry, body := decodeEntry(line)
			w.observer(entry, body)
		}
		if _, err := io.WriteString(w.out, w.render(line)+"\n"); err != nil {
			return written, err
		}
	}

	return written, nil
}

func (w *Writer) render(msg string) string {
	entry, body := decodeEntry(msg)
	if entry.Component == "" {
		entry.Component, body = splitComponent(body)
	}
	timestamp := time.Now().Format("2006/01/02 15:04:05")

	if w.format == FormatPlain {
		if entry.Component == "" {
			return timestamp + " " + body
		}
		return fmt.Sprintf("%s [%s] %s", timestamp, entry.Component, body)
	}

	styledTimestamp := style(timestamp, w.theme.Timestamp, w.color)
	if entry.Component == "" {
		return styledTimestamp + " " + w.styleMessage(body, entry.Level)
	}

	styledComponent := style("["+entry.Component+"]", w.componentColor(entry), w.color)
	return styledTimestamp + " " + styledComponent + " " + w.styleMessage(body, entry.Level)
}

func splitComponent(msg string) (component, body string) {
	trimmed := strings.TrimSpace(msg)
	if trimmed == "" {
		return "", ""
	}

	if matches := bracketPrefixPattern.FindStringSubmatch(trimmed); len(matches) == 3 {
		return strings.TrimSpace(matches[1]), strings.TrimSpace(matches[2])
	}

	if matches := colonPrefixPattern.FindStringSubmatch(trimmed); len(matches) == 3 {
		prefix := strings.TrimSpace(matches[1])
		if _, excluded := excludedPrefixes[prefix]; !excluded {
			return prefix, strings.TrimSpace(matches[2])
		}
	}

	return "", trimmed
}

func encodeEntry(entry Entry, message string) string {
	return metaPrefix + sanitizeMeta(entry.Component) + "|" + sanitizeMeta(entry.Phase) + "|" + sanitizeMeta(entry.Level) + metaSep + message
}

func decodeEntry(message string) (entry Entry, body string) {
	if !strings.HasPrefix(message, metaPrefix) {
		return Entry{}, message
	}

	payload := strings.TrimPrefix(message, metaPrefix)
	idx := strings.Index(payload, metaSep)
	if idx == -1 {
		return Entry{}, message
	}

	meta := strings.SplitN(payload[:idx], "|", 3)
	if len(meta) != 3 {
		return Entry{}, message
	}

	return Entry{
		Component: meta[0],
		Phase:     meta[1],
		Level:     meta[2],
	}, payload[idx+len(metaSep):]
}

func sanitizeMeta(value string) string {
	value = strings.ReplaceAll(value, "|", "/")
	value = strings.ReplaceAll(value, metaSep, " ")
	return strings.TrimSpace(value)
}

func resolveFormat(format Format, isTTY bool) Format {
	if format == FormatAuto {
		if isTTY {
			return FormatRich
		}
		return FormatPlain
	}
	return format
}

func allowColor(opts Options) bool {
	return ColorEnabled(opts.IsTTY, opts.Format, opts.ForceColor, opts.NoColor)
}

// ColorEnabled applies the shared color policy for CLI output and logging.
// Explicit disable options take precedence, then explicit force options, then
// automatic TTY and environment detection.
func ColorEnabled(isTTY bool, format Format, forceColor, noColor bool) bool {
	if noColor || format == FormatPlain {
		return false
	}
	if forceColor || (os.Getenv("FORCE_COLOR") != "" && os.Getenv("FORCE_COLOR") != "0") {
		return true
	}
	if os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" || !isTTY {
		return false
	}
	return true
}

func style(text, color string, enabled bool) string {
	if !enabled || text == "" {
		return text
	}
	if strings.HasPrefix(color, "#") {
		if rgb := ansiTrueColor(color); rgb != "" {
			return rgb + text + ansiReset
		}
	}
	return color + text + ansiReset
}

func ansiTrueColor(hex string) string {
	parsed, err := palettes.ParseHexColor(hex)
	if err != nil {
		return ""
	}
	r, g, b, _ := parsed.RGBA()
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", r>>8, g>>8, b>>8)
}

func (w *Writer) componentColor(entry Entry) string {
	if entry.Phase != "" {
		if color := w.theme.PhaseColor[entry.Phase]; color != "" {
			return color
		}
	}
	return w.theme.Component
}

func (w *Writer) styleMessage(message, level string) string {
	switch strings.ToLower(level) {
	case "warning", "warn":
		return stylePrefixedMessage(message, "Warning:", w.theme.Warning, w.color)
	case "error":
		return stylePrefixedMessage(message, "Error:", w.theme.Error, w.color)
	}

	lower := strings.ToLower(message)
	switch {
	case strings.HasPrefix(lower, "warning:"):
		return stylePrefixedMessage(message, "Warning:", w.theme.Warning, w.color)
	case strings.HasPrefix(lower, "error:"):
		return stylePrefixedMessage(message, "Error:", w.theme.Error, w.color)
	default:
		return message
	}
}

// stylePrefixedMessage colors prefix. If message already starts with prefix
// (case-insensitive), that prefix is replaced; otherwise prefix is prepended.
func stylePrefixedMessage(message, prefix, color string, enabled bool) string {
	if len(message) >= len(prefix) && strings.EqualFold(message[:len(prefix)], prefix) {
		return style(prefix, color, enabled) + message[len(prefix):]
	}
	if message == "" {
		return style(prefix, color, enabled)
	}
	return style(prefix, color, enabled) + " " + message
}
