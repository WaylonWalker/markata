package themes

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestPaletteSwitcherDefersHiddenGridsUntilPickerOpens(t *testing.T) {
	source, err := fs.ReadFile(DefaultStatic(), "js/palette-switcher.js")
	if err != nil {
		t.Fatal(err)
	}
	text := string(source)

	mountStart := strings.Index(text, "function mountPicker(container)")
	openStart := strings.Index(text, "function openPicker(options)")
	closeStart := strings.Index(text, "function closePicker(returnFocus)")
	if mountStart < 0 || openStart < 0 || closeStart < 0 || openStart >= closeStart {
		t.Fatal("palette picker open lifecycle functions were not found in the expected order")
	}
	openPicker := text[openStart:closeStart]
	mountRemainder := text[mountStart:]
	nextFunction := strings.Index(mountRemainder[len("function mountPicker(container)"):], "\n  function ")
	if nextFunction < 0 {
		t.Fatal("could not find the end of mountPicker")
	}
	mountPicker := mountRemainder[:len("function mountPicker(container)")+nextFunction]
	for _, renderCall := range []string{"renderStyleGrid();", "renderFontGrid();"} {
		if strings.Contains(mountPicker, renderCall) {
			t.Errorf("mountPicker must not call %s before the picker is opened", renderCall)
		}
		if !strings.Contains(openPicker, renderCall) {
			t.Errorf("openPicker must call %s to populate the picker on demand", renderCall)
		}
	}
}

func TestThemePicker_ResetButtonIsLabeledInBothTemplates(t *testing.T) {
	resetButton := regexp.MustCompile(`(?s)<button[^>]*data-picker-reset[^>]*>\s*<svg.*?</svg>\s*<span>Reset</span>\s*</button>`)

	for _, tc := range []struct {
		name string
		read func() ([]byte, error)
	}{
		{"embedded theme", func() ([]byte, error) {
			return fs.ReadFile(DefaultTemplates(), "partials/palette-switcher.html")
		}},
		{"root templates", func() ([]byte, error) {
			return os.ReadFile(filepath.Join("..", "..", "templates", "partials", "palette-switcher.html"))
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			content, err := tc.read()
			if err != nil {
				t.Fatal(err)
			}
			if !resetButton.Match(content) {
				t.Error("theme picker reset control must show a visible Reset label")
			}
		})
	}
}
