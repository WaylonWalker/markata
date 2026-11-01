//go:build !js

package plugins

import (
	"strings"
	"testing"
)

func TestValidateNativeMermaidSVG_RejectsActiveContent(t *testing.T) {
	const prefix = `<svg xmlns="http://www.w3.org/2000/svg">`
	for name, body := range map[string]string{
		"script":         `<script>alert(1)</script>`,
		"event":          `<rect onload="alert(1)" />`,
		"foreign_object": `<foreignObject><p>HTML</p></foreignObject>`,
		"javascript":     `<a href="javascript:alert(1)"><text>Click</text></a>`,
		"data":           `<a href="data:text/html;base64,PHNjcmlwdD4="><text>Click</text></a>`,
		"remote_image":   `<image href="https://evil.example/pixel.png" />`,
		"remote_fill":    `<rect fill="url(https://evil.example/pixel.svg)" />`,
		"style_tag":      `<style>@import url(https://evil.example/style.css);</style>`,
		"doctype":        `<!DOCTYPE svg [<!ENTITY external SYSTEM "file:///etc/passwd">]>`,
		"extra_root":     `</svg><svg>`,
	} {
		t.Run(name, func(t *testing.T) {
			payload := []byte(prefix + body + `</svg>`)
			if err := validateNativeMermaidSVG(payload); err == nil {
				t.Errorf("native SVG security gate accepted %q", name)
			}
		})
	}
}

func TestValidateNativeMermaidSVG_AllowsSafeMarkersAndLinks(t *testing.T) {
	for _, body := range []string{
		`<defs><marker id="m"><path d="M0,0 L10,10"/></marker></defs><path marker-end="url(#m)" d="M1,2 L2,3"/>`,
		`<a href="https://example.com/docs"><text>Docs</text></a>`,
		`<a href="/docs"><text>Internal</text></a>`,
		`<text>Escaped &amp; characters &lt;not markup&gt;</text>`,
	} {
		svg := []byte(`<svg xmlns="http://www.w3.org/2000/svg">` + body + `</svg>`)
		if err := validateNativeMermaidSVG(svg); err != nil {
			t.Errorf("rejected safe SVG %q: %v", body, err)
		}
	}
}

func TestNativeMermaid_ClickDirectiveDoesNotInjectActiveSVG(t *testing.T) {
	cfg := nativeMermaidRenderer{}
	cfg.config.Mode = mermaidModeGo
	for _, link := range []string{"javascript:alert(1)", "data:text/html,hello", "https://example.com"} {
		src := "flowchart LR\n A[Clickable] --> B[Done]\n click A href \"" + link + "\""
		svg, err := cfg.render(src)
		if err != nil {
			t.Fatalf("render %q: %v", link, err)
		}
		if strings.Contains(strings.ToLower(svg), "javascript:") || strings.Contains(strings.ToLower(svg), "data:text/html") {
			t.Errorf("unsafe link appeared in SVG output: %s", svg)
		}
	}
}
