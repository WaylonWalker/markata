package plugins

import (
	"strings"
	"testing"
)

func TestTildeAttributionRendering(t *testing.T) {
	p := NewMentionsPlugin()
	known := map[string]*mentionEntry{
		"poteto": {Handle: "poteto", SiteURL: "https://github.com/poteto"},
	}
	tests := []struct {
		name, input, want string
	}{
		{"quote", "nothing else matters\n\n~poteto", "nothing else matters\n\n<a href=\"https://github.com/poteto\" class=\"mention mention--attribution\">~poteto</a>"},
		{"ordinary mention unchanged", "Thanks @poteto!", "Thanks <a href=\"https://github.com/poteto\" class=\"mention\">@poteto</a>!"},
		{"unresolved styled but unlinked", "~missing", "<span class=\"mention mention--attribution mention--unresolved\" data-handle=\"~missing\">~missing</span>"},
		{"unknown at-sign stays text", "@missing", "@missing"},
		{"fenced source unaffected", "```\n~poteto\n```", "```\n~poteto\n```"},
		{"escaped HTML safe", "~poteto", "<a href=\"https://github.com/poteto\" class=\"mention mention--attribution\">~poteto</a>"},
		{"double tilde unaffected", "~~poteto~~", "~~poteto~~"},
		{"tilde punctuated", "~poteto.", "<a href=\"https://github.com/poteto\" class=\"mention mention--attribution\">~poteto</a>."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := p.processMentionsWithMetadata(tt.input, known)
			if got != tt.want {
				t.Errorf("got %q; want %q", got, tt.want)
			}
		})
	}
	if got := p.processMentionsWithMetadata("~poteto", nil); !strings.Contains(got, "mention--unresolved") {
		t.Fatalf("empty map must still style attribution, got %q", got)
	}
}
