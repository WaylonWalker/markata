package diagnostics

import (
	"strings"
	"testing"
)

func TestBareDomainLinkDestinations(t *testing.T) {
	cases := []struct {
		name, content string
		want          int
	}{
		{"reported reproduction", "[pype.dev](pype.dev)", 1},
		{"path query fragment", "[site](blog.example.dev/posts?q=1#intro)", 1},
		{"percent-encoded query", "[site](example.dev/search?q=hello%20world)", 1},
		{"percent-encoded path", "[site](example.dev/a%2Fb)", 1},
		{"percent-encoded hostname", "[site](exa%6Dple.dev/path)", 0},
		{"percent-encoded local file", "[local](README.md?name=hello%20world)", 0},
		{"angle wrapped", "[site](<example.com/docs>)", 1},
		{"nested label", "[a [nested] label](example.org)", 1},
		{"multiple", "[a](example.org) and [b](example.net)", 2},
		{"quoted title", "[a](example.org \"site\")", 1},
		{"real external", "[site](https://example.com)", 0},
		{"local markdown", "[README](README.md) [notes](notes.html)", 0},
		{"relative path", "[a](./example.dev) [b](../example.com) [c](/example.com)", 0},
		{"anchor", "[a](#example.com)", 0},
		{"mailto", "[a](mailto:user@example.com)", 0},
		{"protocol relative", "[a](//example.com)", 0},
		{"inline code", "\x60[a](example.com)\x60 and [b](https://example.com)", 0},
		{"escaped link", "\\[a](example.com)", 0},
		{"fenced code", "\x60\x60\x60md\n[a](example.com)\n\x60\x60\x60", 0},
		{"tilde fence", "~~~\n[a](example.com)\n~~~", 0},
		{"indented code", "    [a](example.com)", 0},
		{"plain text", "example.com and [not a destination]", 0},
		{"image", "![alt](example.com)", 0},
		{"incomplete link", "[a](example.com", 0},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got := checkBareDomainLinks("post.md", tt.content, false, "")
			if len(got) != tt.want {
				t.Fatalf("got %d warnings, want %d: %+v", len(got), tt.want, got)
			}
			for _, issue := range got {
				if issue.Code != diagnosticCodeBareDomainLink || issue.Severity != SeverityWarning || issue.Fixable {
					t.Fatalf("unsafe diagnostic: %+v", issue)
				}
				if !strings.Contains(issue.Message, "https://") {
					t.Errorf("missing guidance: %+v", issue)
				}
			}
		})
	}
}

func TestBareDomainLinkSourceRangeWithFrontmatter(t *testing.T) {
	source := "---\ntitle: Example\n---\nA [link](pype.dev/path)."
	got := Check("post.md", source, nil)
	for _, issue := range got {
		if issue.Code != diagnosticCodeBareDomainLink {
			continue
		}
		if issue.Range.StartLine != 3 || issue.Range.EndLine != 3 || issue.Range.StartCol != 9 || issue.Range.EndCol != 22 {
			t.Fatalf("unexpected source range: %+v", issue.Range)
		}
		return
	}
	t.Fatal("missing bare-domain-link warning")
}
