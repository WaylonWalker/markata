package plugins

import (
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
)

func TestPublishHTMLPluginCriticalStageErrors(t *testing.T) {
	plugin := NewPublishHTMLPlugin()
	if !plugin.CriticalStageErrors(lifecycle.StageWrite) {
		t.Fatal("publish_html write failures must be critical")
	}
	if plugin.CriticalStageErrors(lifecycle.StageRender) {
		t.Fatal("publish_html render failures should keep lifecycle defaults")
	}
}
