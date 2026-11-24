package plugins

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
)

func TestCSSBundlePages_CascadeAndIdempotence(t *testing.T) {
	output := t.TempDir()
	writePageCSS(t, output, "first.css", `@layer tokens { :root { --ink: black } }`)
	writePageCSS(t, output, "last.css", `.hero { color: blue; background: url('/image.svg') }`)
	page := []byte(`<!doctype html><html><head><link rel="stylesheet" href="/first.css"><meta name="x" content="keep"><style>.hero { color: red }</style><link rel="stylesheet" href="/last.css"></head><body><p>Keep body</p></body></html>`)
	p := NewCSSBundlePlugin()
	got, err := p.bundlePageStyles(page, output)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(got), `rel="stylesheet"`) != 1 || !strings.Contains(string(got), `<meta name="x" content="keep">`) || !strings.Contains(string(got), `<p>Keep body</p>`) {
		t.Fatalf("incorrect bundled markup: %s", got)
	}
	if !strings.Contains(string(got), `data-css-sources="/first.css /last.css"`) {
		t.Fatal("bundle must advertise its sources to the conditional CSS loader")
	}
	files, err := filepath.Glob(filepath.Join(output, "css", "page-*.css"))
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 {
		t.Fatalf("got %d bundles, want one", len(files))
	}
	css, err := os.ReadFile(files[0])
	if err != nil {
		t.Fatal(err)
	}
	if string(css) != "@layer tokens { :root { --ink: black } }\n.hero { color: red }\n.hero { color: blue; background: url('/image.svg') }\n" {
		t.Fatalf("cascade order changed: %s", css)
	}
	again, err := p.bundlePageStyles(got, output)
	if err != nil || !bytes.Equal(got, again) {
		t.Fatal("bundling must be idempotent")
	}
}

func TestCSSBundlePages_PreservesBoundaries(t *testing.T) {
	for _, boundary := range []string{
		`<script>window.color = getComputedStyle(document.documentElement).color</script>`,
		`<link rel="stylesheet" href="https://example.com/styles.css">`,
		`<link rel="stylesheet" href="/missing.css">`,
		`<link rel="stylesheet" href="/relative.css">`,
		`<link rel="stylesheet" href="/imports.css">`,
		`<link rel="stylesheet" href="/a.css" media="print">`,
		`<link rel="stylesheet" href="/a.css" integrity="sha256-test">`,
		`<style nonce="keep">p { color: green }</style>`,
		`<noscript><link rel="stylesheet" href="/a.css"></noscript>`,
		`<template><style>p { color: green }</style></template>`,
	} {
		t.Run(boundary, func(t *testing.T) {
			output := t.TempDir()
			writePageCSS(t, output, "a.css", `p { color: red }`)
			writePageCSS(t, output, "relative.css", `p { background: url(../image.png) }`)
			writePageCSS(t, output, "imports.css", `@import "/other.css";`)
			pair := `<link rel="stylesheet" href="/a.css"><style>p { color: blue }</style>`
			page := []byte(`<html><head>` + pair + boundary + pair + `</head><body></body></html>`)
			got, err := NewCSSBundlePlugin().bundlePageStyles(page, output)
			if err != nil {
				t.Fatal(err)
			}
			before, after, ok := strings.Cut(string(got), boundary)
			if !ok || strings.Count(before, "data-css-sources=") != 1 || strings.Count(after, "data-css-sources=") != 1 {
				t.Fatalf("must retain the boundary between bundles: %s", got)
			}
		})
	}
}

func TestCSSBundlePages_ExplicitSelectionAndValidation(t *testing.T) {
	for _, test := range []struct {
		name, page    string
		enabled, fast bool
		wantError     bool
	}{
		{name: "selected", page: "index.html", enabled: true},
		{name: "disabled", page: "index.html"},
		{name: "fast", page: "index.html", enabled: true, fast: true},
		{name: "escape", page: "../index.html", enabled: true, wantError: true},
		{name: "not HTML", page: "css/main.css", enabled: true, wantError: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			output := t.TempDir()
			writePageCSS(t, output, "a.css", `p { color: red }`)
			page := `<html><head><link rel="stylesheet" href="/a.css"><style>p { color: blue }</style></head></html>`
			writePageCSS(t, output, "index.html", page)
			writePageCSS(t, output, "untouched.html", page)
			m := lifecycle.NewManager()
			m.SetConfig(&lifecycle.Config{OutputDir: output, Extra: map[string]interface{}{
				"fast_mode":  test.fast,
				"css_bundle": map[string]interface{}{"enabled": test.enabled, "pages": []interface{}{test.page}},
			}})
			p := NewCSSBundlePlugin()
			if err := p.Configure(m); err != nil {
				t.Fatal(err)
			}
			if err := p.Cleanup(m); (err != nil) != test.wantError {
				t.Fatalf("Cleanup() = %v, want error %v", err, test.wantError)
			}
			got, err := os.ReadFile(filepath.Join(output, "index.html"))
			if err != nil {
				t.Fatal(err)
			}
			changed := string(got) != page
			if changed != (test.enabled && !test.fast && !test.wantError) {
				t.Fatalf("unexpected page change: %s", got)
			}
			untouched, err := os.ReadFile(filepath.Join(output, "untouched.html"))
			if err != nil {
				t.Fatal(err)
			}
			if string(untouched) != page {
				t.Fatal("unselected pages must remain unchanged")
			}
		})
	}
}

func writePageCSS(t *testing.T, output, path, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(output, path), []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestCSSBundlePages_RefreshesImageLibraryOwnershipAfterFinalCleanup(t *testing.T) {
	output := t.TempDir()
	if err := os.MkdirAll(filepath.Join(output, "images"), 0o755); err != nil {
		t.Fatal(err)
	}
	writePageCSS(t, output, "first.css", "p { color: red }")
	writePageCSS(t, output, "second.css", "p { color: blue }")
	page := filepath.Join(output, "images", "index.html")
	initial := []byte(`<html><head><link rel="stylesheet" href="/first.css"><link rel="stylesheet" href="/second.css"></head><body>Images</body></html>`)
	if err := os.WriteFile(page, initial, 0o600); err != nil {
		t.Fatal(err)
	}
	cacheDir := t.TempDir()
	cache := buildcache.New(cacheDir)
	cache.SetImageLibraryOutputs(output, output, map[string]string{
		page: buildcache.ContentHash(string(initial)),
	})
	if err := cache.Save(); err != nil {
		t.Fatal(err)
	}
	m := lifecycle.NewManager()
	m.Cache().Set("build_cache", cache)
	m.SetConfig(&lifecycle.Config{OutputDir: output, Extra: map[string]interface{}{
		"css_bundle": map[string]interface{}{"enabled": true, "pages": []interface{}{"images/index.html"}},
	}})
	p := NewCSSBundlePlugin()
	if err := p.Configure(m); err != nil {
		t.Fatal(err)
	}
	if NewStaticAssetsPlugin().Priority(lifecycle.StageCleanup) >= p.Priority(lifecycle.StageCleanup) {
		t.Fatal("page bundling needs published versioned CSS before it runs")
	}
	if err := p.Cleanup(m); err != nil {
		t.Fatal(err)
	}
	bundled, err := os.ReadFile(page)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(bundled, initial) {
		t.Fatal("expected page bundling to rewrite image library HTML")
	}
	want := buildcache.ContentHash(string(bundled))
	if got, ok := cache.ImageLibraryOutputHash(output, output, page); !ok || got != want {
		t.Fatalf("bundled output ownership hash %q, owned=%v, want %q", got, ok, want)
	}
	loaded, err := buildcache.Load(cacheDir)
	if err != nil {
		t.Fatal(err)
	}
	if got, ok := loaded.ImageLibraryOutputHash(output, output, page); !ok || got != want {
		t.Fatalf("persisted bundled page hash %q, owned=%v, want %q", got, ok, want)
	}
}

func TestCSSBundlePages_FinalFingerprintThenBundleColdWarm(t *testing.T) {
	out := t.TempDir()
	for _, dir := range []string{"css", "images"} {
		if err := os.MkdirAll(filepath.Join(out, dir), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	css := map[string]string{
		"css/first.css":  ".first { color: red }",
		"css/second.css": ".second { color: blue }",
	}
	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{OutputDir: out, Extra: map[string]interface{}{
		"css_bundle": map[string]interface{}{"enabled": true, "pages": []interface{}{"images/index.html"}},
	}})
	var links strings.Builder
	for asset, content := range css {
		if err := os.WriteFile(filepath.Join(out, asset), []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
		// Rendered HTML may contain old Configure hashes prior to purge.
		m.SetAssetHash(asset, "deadbeef")
		fmt.Fprintf(&links, `<link rel="stylesheet" href="/%s">`, strings.TrimSuffix(asset, ".css")+".deadbeef.css")
	}
	initial := []byte("<html><head>" + links.String() + "</head><body>Images</body></html>")
	page := filepath.Join(out, "images", "index.html")
	if err := os.WriteFile(page, initial, 0o600); err != nil {
		t.Fatal(err)
	}
	cacheDir := t.TempDir()
	cache := buildcache.New(cacheDir)
	cache.SetImageLibraryOutputs(out, out, map[string]string{page: buildcache.ContentHash(string(initial))})
	m.Cache().Set("build_cache", cache)
	static := NewStaticAssetsPlugin()
	bundle := NewCSSBundlePlugin()
	if err := bundle.Configure(m); err != nil {
		t.Fatal(err)
	}
	if err := static.Cleanup(m); err != nil {
		t.Fatal(err)
	}
	for asset, content := range css {
		hash := fmt.Sprintf("%x", sha256.Sum256([]byte(content)))[:8]
		published := strings.TrimSuffix(asset, ".css") + "." + hash + ".css"
		data, err := os.ReadFile(filepath.Join(out, published))
		if err != nil || string(data) != content {
			t.Fatalf("incorrect published versioned CSS %s: %v", published, err)
		}
	}
	if err := bundle.Cleanup(m); err != nil {
		t.Fatal(err)
	}
	final, err := os.ReadFile(page)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(final, []byte(`data-css-sources=`)) || strings.Contains(string(final), "deadbeef.css") {
		t.Fatalf("bundle did not use final CSS hashes: %s", final)
	}
	bundles, err := filepath.Glob(filepath.Join(out, "css", "page-*.css"))
	if err != nil || len(bundles) != 1 {
		t.Fatalf("expected one final CSS bundle: %v, %v", bundles, err)
	}
	bundleBytes, err := os.ReadFile(bundles[0])
	if err != nil {
		t.Fatal(err)
	}
	for _, content := range css {
		if !bytes.Contains(bundleBytes, []byte(content)) {
			t.Fatalf("missing stylesheet in final bundle: %s", bundleBytes)
		}
	}
	if got, ok := cache.ImageLibraryOutputHash(out, out, page); !ok || got != buildcache.ContentHash(string(final)) {
		t.Fatal("final bundled HTML ownership does not match published bytes")
	}
	if err := static.Cleanup(m); err != nil {
		t.Fatal(err)
	}
	if err := bundle.Cleanup(m); err != nil {
		t.Fatal(err)
	}
	warm, err := os.ReadFile(page)
	if err != nil || !bytes.Equal(warm, final) {
		t.Fatalf("unchanged warm output differs from cold: %v", err)
	}
}
