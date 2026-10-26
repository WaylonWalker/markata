package cmd

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/WaylonWalker/markata/pkg/config"
)

const (
	// diffContext is the number of unchanged lines shown around each change.
	diffContext = 2
	// maxDiffCells bounds the line-matching table; larger edits fall back to
	// showing the changed region as one replaced block.
	maxDiffCells = 4_000_000
)

// sensitiveLinePattern captures the key of a TOML "key = value" or YAML
// "key: value" line.
var sensitiveLinePattern = regexp.MustCompile(`^(\s*"?([A-Za-z0-9_-]+)"?\s*[=:]\s*)(.+)$`)

type diffOp struct {
	kind byte // ' ', '-', '+'
	line string
}

// unifiedDiff renders a unified diff (without file headers) from before to
// after. Config edits are small and local, so the common prefix and suffix
// are trimmed before matching the middle.
func unifiedDiff(before, after string) string {
	a, b := diffLines(before), diffLines(after)
	prefix := 0
	for prefix < len(a) && prefix < len(b) && a[prefix] == b[prefix] {
		prefix++
	}
	suffix := 0
	for suffix < len(a)-prefix && suffix < len(b)-prefix && a[len(a)-1-suffix] == b[len(b)-1-suffix] {
		suffix++
	}
	ops := make([]diffOp, 0, len(a)+len(b))
	for _, line := range a[:prefix] {
		ops = append(ops, diffOp{' ', line})
	}
	ops = append(ops, diffMiddle(a[prefix:len(a)-suffix], b[prefix:len(b)-suffix])...)
	for _, line := range a[len(a)-suffix:] {
		ops = append(ops, diffOp{' ', line})
	}
	return formatHunks(ops)
}

func diffLines(text string) []string {
	if text == "" {
		return nil
	}
	return strings.Split(strings.TrimSuffix(strings.ReplaceAll(text, "\r\n", "\n"), "\n"), "\n")
}

// diffMiddle matches lines with a longest-common-subsequence table.
func diffMiddle(a, b []string) []diffOp {
	var ops []diffOp
	if len(a)*len(b) > maxDiffCells {
		for _, line := range a {
			ops = append(ops, diffOp{'-', line})
		}
		for _, line := range b {
			ops = append(ops, diffOp{'+', line})
		}
		return ops
	}
	lcs := make([][]int, len(a)+1)
	for i := range lcs {
		lcs[i] = make([]int, len(b)+1)
	}
	for i := len(a) - 1; i >= 0; i-- {
		for j := len(b) - 1; j >= 0; j-- {
			if a[i] == b[j] {
				lcs[i][j] = lcs[i+1][j+1] + 1
			} else {
				lcs[i][j] = max(lcs[i+1][j], lcs[i][j+1])
			}
		}
	}
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		switch {
		case a[i] == b[j]:
			ops = append(ops, diffOp{' ', a[i]})
			i++
			j++
		case lcs[i+1][j] >= lcs[i][j+1]:
			ops = append(ops, diffOp{'-', a[i]})
			i++
		default:
			ops = append(ops, diffOp{'+', b[j]})
			j++
		}
	}
	for ; i < len(a); i++ {
		ops = append(ops, diffOp{'-', a[i]})
	}
	for ; j < len(b); j++ {
		ops = append(ops, diffOp{'+', b[j]})
	}
	return ops
}

// formatHunks groups changed lines with diffContext lines of context.
func formatHunks(ops []diffOp) string {
	var out strings.Builder
	oldLine, newLine := 1, 1
	for start := 0; start < len(ops); {
		if ops[start].kind == ' ' {
			start++
			oldLine++
			newLine++
			continue
		}
		// Extend the hunk while changes are within 2*diffContext lines.
		end := start
		for k := start; k < len(ops); k++ {
			if ops[k].kind != ' ' {
				end = k
			} else if k-end > 2*diffContext {
				break
			}
		}
		from := max(start-diffContext, 0)
		to := min(end+diffContext, len(ops)-1)
		hunkOld, hunkNew := oldLine-(start-from), newLine-(start-from)
		oldCount, newCount := 0, 0
		var body strings.Builder
		for k := from; k <= to; k++ {
			op := ops[k]
			body.WriteByte(op.kind)
			body.WriteString(op.line)
			body.WriteByte('\n')
			if op.kind != '+' {
				oldCount++
			}
			if op.kind != '-' {
				newCount++
			}
		}
		out.WriteString("@@ -" + hunkRange(hunkOld, oldCount) + " +" + hunkRange(hunkNew, newCount) + " @@\n")
		out.WriteString(body.String())
		for k := start; k <= to; k++ {
			if ops[k].kind != '+' {
				oldLine++
			}
			if ops[k].kind != '-' {
				newLine++
			}
		}
		start = to + 1
	}
	return out.String()
}

func hunkRange(start, count int) string {
	if count == 0 {
		start--
	}
	return strconv.Itoa(start) + "," + strconv.Itoa(count)
}

// redactSensitiveLines hides the values of secret keys that appear as diff
// context, matching the sidebar, which never shows sensitive settings.
func redactSensitiveLines(diff string) string {
	lines := strings.Split(diff, "\n")
	for i, line := range lines {
		if line == "" || line[0] == '@' {
			continue
		}
		m := sensitiveLinePattern.FindStringSubmatch(line[1:])
		if len(m) > 2 && config.SensitiveSettingName(m[2]) {
			lines[i] = line[:1] + m[1] + "\"…\""
		}
	}
	return strings.Join(lines, "\n")
}
