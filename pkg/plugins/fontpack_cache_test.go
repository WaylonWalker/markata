package plugins

import (
	"math"
	"reflect"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/fontpacks"
)

func fontpackTestCacheKey(t *testing.T, content, defaultName string, names []string, picker bool, source *fontpacks.CatalogSource) string {
	t.Helper()
	coverage, err := fontpacks.CollectCoverage(strings.NewReader(content))
	if err != nil {
		t.Fatal(err)
	}
	key, err := fontpackCacheKey(coverage.Signature(), defaultName, names, picker, source)
	if err != nil {
		t.Fatal(err)
	}
	return key
}

func TestFontpackCacheKeyCanonicalizesSelection(t *testing.T) {
	source, err := fontpacks.BuiltinSource()
	if err != nil {
		t.Fatal(err)
	}
	handwritten, _, err := source.Catalog.ResolvePack("handwritten")
	if err != nil {
		t.Fatal(err)
	}
	if source.Catalog.Aliases == nil {
		source.Catalog.Aliases = make(map[string]string)
	}
	source.Catalog.Aliases["reader-alias"] = handwritten
	first := fontpackTestCacheKey(t, "<p>text</p>", "system", []string{"system", "handwritten"}, false, source)
	names := []string{"reader-alias", "system", "handwritten"}
	original := append([]string(nil), names...)
	reordered := fontpackTestCacheKey(t, "<p>text</p>", "system", names, false, source)
	if first != reordered {
		t.Fatal("equivalent aliases, ordering and duplicate selections changed cache identity")
	}
	if !reflect.DeepEqual(names, original) {
		t.Fatal("computing cache identity mutated caller's names")
	}
}

func TestFontpackCacheKeyTracksEmissionInputs(t *testing.T) {
	source, err := fontpacks.BuiltinSource()
	if err != nil {
		t.Fatal(err)
	}
	names := []string{"system", "handwritten"}
	first := fontpackTestCacheKey(t, "<p>text</p>", "system", names, false, source)
	for _, test := range []struct {
		name    string
		content string
		def     string
		names   []string
		picker  bool
	}{
		{name: "default", content: "<p>text</p>", def: "handwritten", names: names},
		{name: "picker", content: "<p>text</p>", def: "system", names: names, picker: true},
		{name: "pack set", content: "<p>text</p>", def: "system", names: []string{"system"}},
		{name: "coverage", content: "<p>text Ж</p>", def: "system", names: names},
	} {
		t.Run(test.name, func(t *testing.T) {
			if fontpackTestCacheKey(t, test.content, test.def, test.names, test.picker, source) == first {
				t.Fatalf("%s did not invalidate cache identity", test.name)
			}
		})
	}
	source.ContentDigest = "different bundled metadata"
	if fontpackTestCacheKey(t, "<p>text</p>", "system", names, false, source) == first {
		t.Fatal("bundled manifest revision did not invalidate cache identity")
	}
}

func TestFontpackCacheKeyRejectsInvalidInputs(t *testing.T) {
	source, err := fontpacks.BuiltinSource()
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		name   string
		def    string
		names  []string
		source *fontpacks.CatalogSource
	}{
		{name: "missing source", def: "system", names: []string{"system"}},
		{name: "missing catalog", def: "system", names: []string{"system"}, source: &fontpacks.CatalogSource{}},
		{name: "unknown default", def: "unknown", names: []string{"system"}, source: source},
		{name: "unknown selected pack", def: "system", names: []string{"system", "unknown"}, source: source},
		{name: "default absent", def: "system", names: []string{"handwritten"}, source: source},
	} {
		t.Run(test.name, func(t *testing.T) {
			if _, err := fontpackCacheKey("", test.def, test.names, false, test.source); err == nil {
				t.Fatal("invalid cache inputs accepted")
			}
		})
	}
	source.ContentDigest = ""
	if _, err := fontpackCacheKey("", "system", []string{"system"}, false, source); err == nil {
		t.Fatal("bundled source without metadata digest accepted")
	}
	source.Builtin = false
	pack := source.Catalog.FontPacks["system"]
	role := pack.Roles["body"]
	role.Weight = math.NaN()
	pack.Roles["body"] = role
	if _, err := fontpackCacheKey("", "system", []string{"system"}, false, source); err == nil {
		t.Fatal("unserializable catalog accepted as an empty successful cache key")
	}
}
