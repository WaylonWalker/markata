package workbench

import (
	"strings"
	"testing"
)

func TestMarkSourceContainers(t *testing.T) {
	source := strings.Join([]string{
		"Intro",
		"::: md-section",
		"body",
		"::: ",
		"```markdown",
		"::: ignored",
		"```",
		":::: md-grid {data-columns=\"2\"}",
		"inside",
		"::::",
		"~~~markdown",
		"::: also-ignored",
		"~~~",
	}, "\n")

	got, anchors := markSourceContainers(source)
	for _, want := range []string{
		`::: md-section {data-markata-source-line="1" data-markata-source-end-line="3"}`,
		"```markdown\n::: ignored\n```",
		`:::: md-grid {data-columns="2" data-markata-source-line="7" data-markata-source-end-line="9"}`,
		"~~~markdown\n::: also-ignored\n~~~",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("marked source does not contain %q:\n%s", want, got)
		}
	}
	if gotLines, sourceLines := strings.Count(got, "\n"), strings.Count(source, "\n"); gotLines != sourceLines {
		t.Errorf("source line count changed: got %d newlines, want %d", gotLines, sourceLines)
	}
	if want := []int{1, 3, 7, 9}; len(anchors) != len(want) || anchors[0] != want[0] || anchors[1] != want[1] || anchors[2] != want[2] || anchors[3] != want[3] {
		t.Errorf("container anchor lines = %v, want %v", anchors, want)
	}
	if twice, _ := markSourceContainers(got); twice != got {
		t.Errorf("marking source twice changed it again:\nfirst:\n%s\nsecond:\n%s", got, twice)
	}
}

func TestMarkSourceContainers_UnclosedContainerHasNoEndAnchor(t *testing.T) {
	marked, anchors := markSourceContainers("::: note\nbody")
	if !strings.Contains(marked, `data-markata-source-line="0"`) {
		t.Fatalf("unclosed container has no start anchor: %s", marked)
	}
	if strings.Contains(marked, "data-markata-source-end-line") {
		t.Fatalf("unclosed container has an end anchor: %s", marked)
	}
	if len(anchors) != 1 || anchors[0] != 0 {
		t.Fatalf("unclosed container anchors = %v, want [0]", anchors)
	}
}
