package visualizations

import (
	"strings"
	"testing"
)

func TestExcalidrawRendererVersion_IncludesBundledFonts(t *testing.T) {
	first, err := ExcalidrawRendererVersion()
	if err != nil {
		t.Fatal(err)
	}
	second, err := ExcalidrawRendererVersion()
	if err != nil {
		t.Fatal(err)
	}
	if first != second || !strings.HasPrefix(first, "excalidraw-0.18.1-static-svg-v2:") {
		t.Fatalf("renderer versions = %q, %q; want a stable version with the bundled font hash", first, second)
	}
}
