package plugins

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

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
