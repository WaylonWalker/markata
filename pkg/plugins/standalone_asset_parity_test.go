package plugins

import (
	"reflect"
	"testing"
)

func TestStandaloneProfile_AssetManifestParity(t *testing.T) {
	for _, source := range []string{
		"# Header\n\nPlain text",
		"\x60\x60\x60go\nfmt.Println(1)\n\x60\x60\x60",
		"\x60\x60\x60mermaid\ngraph TD\nA-->B\n\x60\x60\x60",
	} {
		linked, err := RenderStandalone(source, StandaloneOptions{Profile: RenderProfileWeb})
		if err != nil {
			t.Fatal(err)
		}
		inline, err := RenderStandalone(source, StandaloneOptions{Profile: RenderProfileWeb, InlineAssets: true})
		if err != nil {
			t.Fatal(err)
		}
		if linked.HTML != inline.HTML || !reflect.DeepEqual(linked.Features, inline.Features) {
			t.Fatalf("asset delivery policy changed Markdown or feature signals for %q", source)
		}
		if len(linked.Assets) != len(inline.Assets) {
			t.Fatalf("asset count changed for %q", source)
		}
		for i := range linked.Assets {
			left, right := linked.Assets[i], inline.Assets[i]
			if left.ID != right.ID || left.Kind != right.Kind || left.Placement != right.Placement || left.Module != right.Module {
				t.Fatalf("asset ordering or type changed at %d: %+v vs %+v", i, left, right)
			}
		}
	}
}
