package cmd

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/WaylonWalker/markata/pkg/suggest"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type inputError struct {
	problem string
	matches []string
	prefix  string
	next    string
}

// Diagnostic carries context needed to render actionable, consistently styled
// command failures. Empty fields are omitted by the renderer.
type Diagnostic struct {
	Severity   string
	Code       string
	Problem    string
	Cause      error
	Source     string
	Fix        string
	Command    string
	Tips       []string
	Details    []string
	StatusCode int
}

func (d *Diagnostic) Error() string {
	if d == nil {
		return ""
	}
	return d.Problem
}
func (d *Diagnostic) Unwrap() error {
	if d == nil {
		return nil
	}
	return d.Cause
}
func (d *Diagnostic) ExitCode() int {
	if d == nil || d.StatusCode == 0 {
		return 1
	}
	return d.StatusCode
}

func (e *inputError) Error() string {
	var b strings.Builder
	b.WriteString(e.problem)
	if len(e.matches) > 0 {
		b.WriteString("\nDid you mean:")
		for _, match := range e.matches {
			fmt.Fprintf(&b, "\n  %s%s", e.prefix, match)
		}
	}
	if e.next != "" {
		fmt.Fprintf(&b, "\nNext: %s", e.next)
	}
	return b.String()
}

func clarifyFlagError(cmd *cobra.Command, err error) error {
	message := err.Error()
	if strings.HasPrefix(message, "invalid argument ") && strings.Contains(message, `for "--port" flag`) {
		return &inputError{problem: `invalid value for --port; enter a whole number from 1 to 65535`, next: "Example: 'markata-go serve --port 8001'."}
	}
	name, isUnknown := strings.CutPrefix(message, "unknown flag: --")
	if isUnknown {
		name = strings.TrimSpace(name)
		candidates := make([]string, 0)
		addFlags := func(flags *pflag.FlagSet) {
			flags.VisitAll(func(flag *pflag.Flag) {
				if !flag.Hidden && flag.Deprecated == "" {
					candidates = append(candidates, flag.Name)
				}
			})
		}
		addFlags(cmd.Flags())
		addFlags(cmd.InheritedFlags())
		return &inputError{
			problem: message,
			matches: suggest.Closest(name, candidates, 3),
			prefix:  "--",
			next:    "Run '" + cmd.CommandPath() + " --help' to see valid flags.",
		}
	}
	return &inputError{problem: message, next: "Run '" + cmd.CommandPath() + " --help' to check flag syntax."}
}

// FormatError renders a single fatal diagnostic for the command entry point.
// It preserves plain text when stderr is redirected or coloring is disabled.
func FormatError(err error) string {
	if err == nil {
		return ""
	}
	message := strings.TrimSpace(err.Error())
	if message == "" {
		return ""
	}
	var input *inputError
	if errors.As(err, &input) {
		return styleErrorLines(message)
	}
	var diagnostic *Diagnostic
	if errors.As(err, &diagnostic) {
		message = diagnostic.Problem
		if diagnostic.Source != "" {
			message += "\nFrom: " + diagnostic.Source
		}
		if diagnostic.Fix != "" {
			message += "\nFix: " + diagnostic.Fix
		}
		if diagnostic.Command != "" {
			message += "\nTry:\n  " + diagnostic.Command
		}
		for _, tip := range diagnostic.Tips {
			message += "\nTip: " + tip
		}
		for _, detail := range diagnostic.Details {
			message += "\n" + detail
		}
	}
	return styleErrorLines(message)
}

func styleErrorLines(message string) string {
	parts := strings.Split(message, "\n")
	theme := currentLogTheme
	enabled := colorEnabledFor(errorOutputIsTerminal())
	parts[0] = colorize("Error", theme.Error, enabled) + ": " + parts[0]
	for i := 1; i < len(parts); i++ {
		line := parts[i]
		switch {
		case line == "Did you mean:":
			parts[i] = colorize(line, theme.Component, enabled)
		case strings.HasPrefix(line, "  --") || strings.HasPrefix(line, "  markata-go "):
			parts[i] = colorize(line, theme.Success, enabled)
		case strings.HasPrefix(line, "Next:"):
			parts[i] = colorize("Next", theme.Component, enabled) + strings.TrimPrefix(line, "Next")
		}
	}
	return strings.Join(parts, "\n")
}

func configKeyError(key string) error {
	lookupKey := strings.TrimPrefix(key, "markata-go.")
	prefix := strings.TrimSuffix(key, lookupKey)
	candidates := make([]string, 0)
	for _, field := range configSettingsForSuggestions() {
		if field == lookupKey {
			return &inputError{
				problem: fmt.Sprintf("configuration key %q is not set in this file", key),
				next:    "Run 'markata-go config show' to see its resolved value, including defaults.",
			}
		}
		candidates = append(candidates, field)
	}
	sort.Strings(candidates)
	return &inputError{
		problem: fmt.Sprintf("unknown configuration key %q", key),
		matches: suggest.Closest(lookupKey, candidates, 3),
		prefix:  prefix,
		next:    "Run 'markata-go config show' to inspect the active configuration.",
	}
}

func invalidChoiceError(flag, value string, choices []string, command string) error {
	choice := suggest.Closest(value, choices, 1)
	d := &Diagnostic{Severity: "error", Code: "input.invalid_choice", Problem: fmt.Sprintf("invalid value %q for %s", value, flag), Details: []string{"Valid values: " + strings.Join(choices, ", ")}, StatusCode: exitCodeUsage}
	if len(choice) > 0 {
		d.Fix = "use " + choice[0]
		if command != "" {
			d.Command = command + " " + flag + " " + choice[0]
		}
	}
	return newUsageError(d)
}
