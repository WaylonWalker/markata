package fontpacks

import (
	"errors"
	"io"
	"io/fs"
	"reflect"
	"slices"
	"strings"
	"testing"
	"testing/fstest"
	"unicode/utf8"
)

func visibleSignature(source string) string {
	seen := map[rune]bool{}
	for _, r := range VisibleText(source) {
		if r != utf8.RuneError {
			if r == '\t' || r == '\n' || r == '\f' || r == '\r' {
				r = ' '
			}
			seen[r] = true
		}
	}
	runes := make([]rune, 0, len(seen))
	for r := range seen {
		runes = append(runes, r)
	}
	slices.Sort(runes)
	return string(runes)
}

func TestCoverageMatchesVisibleText(t *testing.T) {
	for _, source := range []string{
		"", "<!-- empty --><br>", "fragment without tags",
		"<p>Hello &amp; &#x100; &#1046; &nbsp; 中文 😀</p>",
		`<SCRIPT>hidden Ж</SCRIPT><style>hidden 中</style><p title="invisible">visible</p>`,
		"<script/>still hidden</script><style>unfinished",
		"<p>malformed <b>nest</p>tail &bogus; <unfinished",
		"\xff � &#xfffd; <p>\t\n</p>",
		"<textarea>&amp; Ā</textarea><title>Title</title>",
	} {
		t.Run(source, func(t *testing.T) {
			coverage, err := CollectCoverage(strings.NewReader(source))
			if err != nil {
				t.Fatal(err)
			}
			if got, want := coverage.Signature(), visibleSignature(source); got != want {
				t.Fatalf("signature = %q, want %q", got, want)
			}
			// Every byte boundary, including UTF-8, entities and raw text.
			for i := 0; i <= len(source); i++ {
				reader := io.MultiReader(strings.NewReader(source[:i]), strings.NewReader(source[i:]))
				split, err := CollectCoverage(reader)
				if err != nil || split.Signature() != coverage.Signature() {
					t.Fatalf("split %d = %q, %v; want %q", i, split.Signature(), err, coverage.Signature())
				}
			}
		})
	}
	articles := []string{"<article>One &am", "p;</article><scr", "ipt>hidden Ж", "</script><p>Ā</p>"}
	readers := make([]io.Reader, 0, 2*len(articles))
	for _, article := range articles {
		readers = append(readers, strings.NewReader(article), strings.NewReader("\n"))
	}
	coverage, err := CollectCoverage(io.MultiReader(readers...))
	if err != nil || coverage.Signature() != visibleSignature(strings.Join(articles, "\n")+"\n") {
		t.Fatalf("article coverage = %q, %v", coverage.Signature(), err)
	}
	if (Coverage{}).Signature() != "" {
		t.Fatal("zero coverage must be empty")
	}
}

type failingCoverageReader struct{ err error }

func (r failingCoverageReader) Read([]byte) (int, error) { return 0, r.err }

func TestCoverageReaderError(t *testing.T) {
	want := errors.New("read failed")
	for _, prefix := range []string{"", "<p>visible</p>", "<script>hidden"} {
		_, err := CollectCoverage(io.MultiReader(strings.NewReader(prefix), failingCoverageReader{want}))
		if !errors.Is(err, want) {
			t.Fatalf("error = %v, want %v", err, want)
		}
	}
}

func TestCoverageNormalizesCollapsedHTMLWhitespace(t *testing.T) {
	coverage, err := CollectCoverage(strings.NewReader("<p>A\tB\nC\fD\rE&nbsp;中</p>"))
	if err != nil {
		t.Fatal(err)
	}
	if got, want := coverage.Signature(), " ABCDE\u00a0中"; got != want {
		t.Fatalf("coverage signature = %q, want %q", got, want)
	}
}

func TestCoverageFullAndFallback(t *testing.T) {
	c := testCatalog(t)
	profiles := compileSubsetProfiles(c.SubsetProfiles)
	for _, test := range []struct {
		name     string
		base     string
		html     string
		tiers    map[string]Tier
		expected map[string]bool
	}{
		{"base full", "full", "Ā Ж", map[string]Tier{"full": {}}, map[string]bool{"full": true}},
		{"unsupported uses system fallback", "prose-core", "Hello Ā Ж 中", map[string]Tier{"prose-core": {}, "latin-ext": {}, "full": {}}, map[string]bool{"prose-core": true, "latin-ext": true}},
		{"missing extended tier uses system fallback", "prose-core", "Hello Ā Ȁ", map[string]Tier{"prose-core": {}, "full": {}}, map[string]bool{"prose-core": true}},
		{"base fallback", "missing", "Hello Ā", map[string]Tier{"full": {}}, map[string]bool{"full": true}},
		{"missing manifest tier uses system fallback", "prose-core", "Ā", map[string]Tier{"prose-core": {}}, map[string]bool{"prose-core": true}},
	} {
		t.Run(test.name, func(t *testing.T) {
			pack := FontPack{Roles: map[string]Role{"body": {Source: "demo", Tier: test.base}}}
			manifests := map[string]Manifest{"demo": {Tiers: test.tiers}}
			got := requiredTiers(pack, coverageFromString(test.html), profiles, manifests)["demo"]
			if !reflect.DeepEqual(got, test.expected) {
				t.Fatalf("tiers = %v, want %v", got, test.expected)
			}
		})
	}
}

func TestBuiltinHandwrittenTiersWithArticleSeparators(t *testing.T) {
	source, err := BuiltinSource()
	if err != nil {
		t.Fatal(err)
	}
	article, err := CollectCoverage(strings.NewReader("<p>Hello world</p>"))
	if err != nil {
		t.Fatal(err)
	}
	withSeparator, err := CollectCoverage(io.MultiReader(strings.NewReader("<p>Hello world</p>"), strings.NewReader("\n")))
	if err != nil {
		t.Fatal(err)
	}
	for _, picker := range []bool{false, true} {
		names := []string{"handwritten"}
		if picker {
			names = SortedKeys(source.Catalog.FontPacks)
		}
		base, err := source.Catalog.ResolveManyFSWithCoverage(names, source.FS, source.Root, article, ResolveOptions{})
		if err != nil {
			t.Fatal(err)
		}
		separated, err := source.Catalog.ResolveManyFSWithCoverage(names, source.FS, source.Root, withSeparator, ResolveOptions{})
		if err != nil {
			t.Fatal(err)
		}
		if len(base.Assets) == 0 {
			t.Fatal("expected built-in font assets")
		}
		if len(separated.Assets) == 0 {
			t.Fatal("expected assets with article separator")
		}
		if !reflect.DeepEqual(separated.Assets, base.Assets) {
			t.Errorf("picker=%t: newline changed tiers; plain=%v separated=%v", picker, summarizeAssetTiers(base.Assets), summarizeAssetTiers(separated.Assets))
		}
	}
}

func summarizeAssetTiers(assets []Asset) map[string]string {
	result := make(map[string]string, len(assets))
	for _, asset := range assets {
		result[asset.Source] = asset.Tier
	}
	return result
}

type countingCoverageFS struct {
	fs.FS
	opens map[string]int
}

func (f countingCoverageFS) Open(name string) (fs.File, error) {
	f.opens[name]++
	return f.FS.Open(name)
}

func coverageFixture() fstest.MapFS {
	return fstest.MapFS{
		"demo/manifest.yaml": {Data: []byte("id: demo\nfamily: Demo\ntiers:\n  prose-core:\n    file: demo.woff2\n    profile: prose-core\n")},
		"demo/demo.woff2":    {Data: []byte("woff2")},
	}
}

func TestCoverageResolutionSharesValidationAndPreservesOutput(t *testing.T) {
	c := testCatalog(t)
	c.FontPacks["other"] = c.FontPacks["bundled"]
	data := coverageFixture()
	// Include checksum reads when verifying per-call validation reuse.
	hash, _, err := AssetSHA256FS(data, "demo/demo.woff2")
	if err != nil {
		t.Fatal(err)
	}
	data["demo/manifest.yaml"].Data = append(data["demo/manifest.yaml"].Data, []byte("    sha256: "+hash+"\n")...)
	counted := countingCoverageFS{FS: data, opens: map[string]int{}}
	coverage := coverageFromString("<p>Hello</p>")
	got, err := c.ResolveManyFSWithCoverage([]string{"bundled", "other", "bundled"}, counted, ".", coverage, ResolveOptions{ValidateChecksums: true})
	if err != nil {
		t.Fatal(err)
	}
	if counted.opens["demo/manifest.yaml"] != 1 || counted.opens["demo/demo.woff2"] != 2 {
		t.Fatalf("opens = %v, want one manifest, one stat and one checksum read", counted.opens)
	}
	single, err := c.ResolveFS("bundled", data, ".", "<p>Hello</p>")
	if err != nil {
		t.Fatal(err)
	}
	if got.CSS != c.cssForPacks(got.Packs, single.Assets, nil) || !reflect.DeepEqual(got.Assets, single.Assets) || got.Bytes != single.Bytes {
		t.Fatal("multi-pack output differs from existing serialization")
	}
	reversed, err := c.ResolveManyFS([]string{"other", "bundled"}, data, ".", "<p>Hello</p>")
	if err != nil {
		t.Fatal(err)
	}
	if reversed.CSS != got.CSS || !reflect.DeepEqual(reversed.Assets, got.Assets) {
		t.Fatal("reversed/duplicate selection changed CSS or assets")
	}
	// No validation cache may escape one resolve call.
	data["demo/demo.woff2"].Data = []byte("corrupted")
	if _, err := c.ResolveManyFS([]string{"bundled", "other"}, data, ".", "<p>Hello</p>"); err == nil || !strings.Contains(err.Error(), "checksum") {
		t.Fatalf("corruption error = %v", err)
	}
}

func TestCoverageResolutionErrors(t *testing.T) {
	for _, test := range []struct {
		name string
		edit func(*Catalog, fstest.MapFS)
		want string
	}{
		{"missing tier", func(c *Catalog, _ fstest.MapFS) {
			c.FontPacks["bundled"] = FontPack{Performance: Performance{Class: "bundled"}, Roles: map[string]Role{"body": {Source: "demo", Tier: "full"}}}
		}, `no required tier "full"`},
		{"capability", func(c *Catalog, _ fstest.MapFS) {
			c.FontPacks["other"] = FontPack{Performance: Performance{Class: "bundled"}, Roles: map[string]Role{"body": {Source: "demo", Tier: "prose-core", OpticalSize: 12}}}
		}, "cannot satisfy"},
		{"missing asset", func(_ *Catalog, f fstest.MapFS) { delete(f, "demo/demo.woff2") }, "is missing"},
		{"profile", func(_ *Catalog, f fstest.MapFS) {
			f["demo/manifest.yaml"].Data = []byte("tiers:\n  prose-core:\n    file: demo.woff2\n    profile: wrong\n")
		}, "declares profile"},
	} {
		t.Run(test.name, func(t *testing.T) {
			c, data := testCatalog(t), coverageFixture()
			c.FontPacks["other"] = c.FontPacks["bundled"]
			test.edit(c, data)
			text := "Hello"
			if test.name == "missing tier" {
				text = "Ж"
			}
			_, err := c.ResolveManyFSWithCoverage([]string{"bundled", "other"}, data, ".", coverageFromString(text), ResolveOptions{ValidateChecksums: true})
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("error = %v, want %q", err, test.want)
			}
		})
	}
}

func TestPartitionedProseRetainsFallbackAndExtendedCoverage(t *testing.T) {
	c := testCatalog(t)
	c.SubsetProfiles["prose-basic"] = SubsetProfile{Unicode: []string{"U+0020-007E"}}
	c.SubsetProfiles["prose-extra"] = SubsetProfile{Unicode: []string{"U+0100-017F"}}
	profiles := compileSubsetProfiles(c.SubsetProfiles)
	for _, tt := range []struct {
		name, base, text string
		tiers            map[string]Tier
		want             map[string]bool
	}{
		{"complete pair without legacy core", "prose-core", "Hello Ā", map[string]Tier{"prose-basic": {}, "prose-extra": {}, "latin-ext": {}, "full": {}}, map[string]bool{"prose-basic": true, "prose-extra": true}},
		{"display fallback uses pair", "display-core", "Hello", map[string]Tier{"prose-basic": {}, "prose-extra": {}, "full": {}}, map[string]bool{"prose-basic": true, "prose-extra": true}},
		{"partial pair retains core", "prose-core", "Hello", map[string]Tier{"prose-basic": {}, "prose-core": {}, "full": {}}, map[string]bool{"prose-core": true}},
		{"outside pair adds extended", "prose-core", "Hello Ȁ", map[string]Tier{"prose-basic": {}, "prose-extra": {}, "latin-ext": {}}, map[string]bool{"prose-basic": true, "prose-extra": true, "latin-ext": true}},
		{"explicit full remains full", "full", "Hello Ā", map[string]Tier{"prose-basic": {}, "prose-extra": {}, "full": {}}, map[string]bool{"full": true}},
		{"dedicated display remains display", "display-core", "Hello", map[string]Tier{"display-core": {}, "prose-basic": {}, "prose-extra": {}}, map[string]bool{"display-core": true}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			pack := FontPack{Roles: map[string]Role{"body": {Source: "demo", Tier: tt.base}}}
			got := requiredTiers(pack, coverageFromString(tt.text), profiles, map[string]Manifest{"demo": {Tiers: tt.tiers}})["demo"]
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("tiers=%v,want%v", got, tt.want)
			}
		})
	}
}
