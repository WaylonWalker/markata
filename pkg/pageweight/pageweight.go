// Package pageweight reports the uncompressed weight of generated HTML pages.
// It measures build artifacts, not network transfers or Core Web Vitals.
package pageweight

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const SchemaVersion = 1

var (
	scriptPattern = regexp.MustCompile(`(?is)<script\b([^>]*)>(.*?)</script\s*>`)
	stylePattern  = regexp.MustCompile(`(?is)<style\b[^>]*>(.*?)</style\s*>`)
	linkPattern   = regexp.MustCompile(`(?is)<link\b[^>]*>`)
	imagePattern  = regexp.MustCompile(`(?i)<img\b`)
	iframePattern = regexp.MustCompile(`(?i)<iframe\b`)
	graphPattern  = regexp.MustCompile(`(?is)\bdata-graph="([^"]*)"`)
	srcPattern    = regexp.MustCompile(`(?i)(?:^|\s)src\s*=`)
)

// Limits are uncompressed HTML-source budgets, in bytes.
type Limits struct {
	PostHTMLBytes int `json:"post_html_bytes"`
	OtherHTMLBytes int `json:"other_html_bytes"`
	EmbeddedJSONBytes int `json:"embedded_json_bytes"`
	InlineJSBytes int `json:"inline_js_bytes"`
	SidebarJSONBytes int `json:"sidebar_json_bytes"`
}

func DefaultLimits() Limits {
	return Limits{
		PostHTMLBytes: 160 * 1024,
		OtherHTMLBytes: 320 * 1024,
		EmbeddedJSONBytes: 32 * 1024,
		InlineJSBytes: 24 * 1024,
		SidebarJSONBytes: 16 * 1024,
	}
}

type Finding struct {
	Code string `json:"code"`
	Bytes int `json:"bytes"`
	BudgetBytes int `json:"budget_bytes"`
	Advice string `json:"advice"`
}

type Page struct {
	Path string `json:"path"`
	Kind string `json:"kind"`
	HTMLBytes int `json:"html_bytes"`
	EmbeddedJSONBytes int `json:"embedded_json_bytes"`
	InlineJSBytes int `json:"inline_js_bytes"`
	InlineCSSBytes int `json:"inline_css_bytes"`
	PostCopyJSONBytes int `json:"post_copy_json_bytes"`
	SidebarJSONBytes int `json:"sidebar_json_bytes"`
	GraphAttributeBytes int `json:"graph_attribute_bytes"`
	StylesheetLinks int `json:"stylesheet_links"`
	ExternalScripts int `json:"external_scripts"`
	ImageElements int `json:"image_elements"`
	IframeElements int `json:"iframe_elements"`
	Findings []Finding `json:"findings"`
}

type Summary struct {
	Pages int `json:"pages"`
	WarnedPages int `json:"warned_pages"`
	Findings int `json:"findings"`
	TotalHTMLBytes int64 `json:"total_html_bytes"`
	LargestHTMLBytes int `json:"largest_html_bytes"`
}

type Report struct {
	SchemaVersion int `json:"schema_version"`
	Measurement string `json:"measurement"`
	Limits Limits `json:"limits"`
	Summary Summary `json:"summary"`
	Pages []Page `json:"pages"`
}

func Analyze(path string, html []byte, limits Limits) Page {
	p := Page{
		Path: filepath.ToSlash(path),
		Kind: "page",
		HTMLBytes: len(html),
		Findings: []Finding{},
	}
	if bytes.Contains(html, []byte(`class="post h-entry`)) {
		p.Kind = "post"
	} else if bytes.Contains(html, []byte(`feed h-feed`)) {
		p.Kind = "feed"
	}
	for _, match := range scriptPattern.FindAllSubmatch(html, -1) {
		attrs := strings.ToLower(string(match[1]))
		size := len(match[2])
		if srcPattern.MatchString(attrs) {
			p.ExternalScripts++
			continue
		}
		if strings.Contains(attrs, "application/json") || strings.Contains(attrs, "ld+json") {
			p.EmbeddedJSONBytes += size
			if strings.Contains(attrs, "post-copy__payload") {
				p.PostCopyJSONBytes += size
			}
			if strings.Contains(attrs, "feed-sidebar-data") {
				p.SidebarJSONBytes += size
			}
		} else {
			p.InlineJSBytes += size
		}
	}
	for _, match := range stylePattern.FindAllSubmatch(html, -1) {
		p.InlineCSSBytes += len(match[1])
	}
	for _, match := range linkPattern.FindAll(html, -1) {
		if bytes.Contains(bytes.ToLower(match), []byte("stylesheet")) {
			p.StylesheetLinks++
		}
	}
	for _, match := range graphPattern.FindAllSubmatch(html, -1) {
		p.GraphAttributeBytes += len(match[1])
	}
	p.ImageElements = len(imagePattern.FindAllIndex(html, -1))
	p.IframeElements = len(iframePattern.FindAllIndex(html, -1))

	htmlBudget := limits.OtherHTMLBytes
	if p.Kind == "post" {
		htmlBudget = limits.PostHTMLBytes
	}
	if htmlBudget > 0 && p.HTMLBytes > htmlBudget {
		p.Findings = append(p.Findings, Finding{"HTML_OVER_BUDGET", p.HTMLBytes, htmlBudget,
			"Inspect repeated inline content, sidebars, embedded data, and large markup; move optional features to build-generated assets."})
	}
	if limits.EmbeddedJSONBytes > 0 && p.EmbeddedJSONBytes > limits.EmbeddedJSONBytes {
		p.Findings = append(p.Findings, Finding{"EMBEDDED_JSON_OVER_BUDGET", p.EmbeddedJSONBytes, limits.EmbeddedJSONBytes,
			"Split nonessential JSON into cacheable static data and load it only when its feature is opened."})
	}
	if limits.InlineJSBytes > 0 && p.InlineJSBytes > limits.InlineJSBytes {
		p.Findings = append(p.Findings, Finding{"INLINE_JS_OVER_BUDGET", p.InlineJSBytes, limits.InlineJSBytes,
			"Move reusable script logic into cacheable JS and initialize optional behavior on interaction."})
	}
	if limits.SidebarJSONBytes > 0 && p.SidebarJSONBytes > limits.SidebarJSONBytes {
		p.Findings = append(p.Findings, Finding{"SIDEBAR_JSON_OVER_BUDGET", p.SidebarJSONBytes, limits.SidebarJSONBytes,
			"Ship the current feed navigation in HTML, and fetch other feed membership lists only when selected."})
	}
	return p
}

// Scan checks generated HTML artifacts. It deliberately does not fetch assets or
// estimate compressed network transfer. Stable ordering makes agent diffs useful.
func Scan(outputDir string, limits Limits) (Report, error) {
	report := Report{
		SchemaVersion: SchemaVersion,
		Measurement: "uncompressed_generated_html_source",
		Limits: limits,
		Pages: []Page{},
	}
	err := filepath.WalkDir(outputDir, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			if entry.Name() == ".markata" {
				return filepath.SkipDir
			}
			return nil
		}
		if !entry.Type().IsRegular() || !strings.EqualFold(filepath.Ext(path), ".html") {
			return nil
		}
		relative, err := filepath.Rel(outputDir, path)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read generated page %s: %w", path, err)
		}
		p := Analyze(relative, data, limits)
		report.Pages = append(report.Pages, p)
		report.Summary.TotalHTMLBytes += int64(p.HTMLBytes)
		if p.HTMLBytes > report.Summary.LargestHTMLBytes {
			report.Summary.LargestHTMLBytes = p.HTMLBytes
		}
		if len(p.Findings) > 0 {
			report.Summary.WarnedPages++
			report.Summary.Findings += len(p.Findings)
		}
		return nil
	})
	if err != nil {
		return Report{}, fmt.Errorf("scan generated HTML in %s: %w", outputDir, err)
	}
	sort.Slice(report.Pages, func(i, j int) bool {
		if report.Pages[i].HTMLBytes != report.Pages[j].HTMLBytes {
			return report.Pages[i].HTMLBytes > report.Pages[j].HTMLBytes
		}
		return report.Pages[i].Path < report.Pages[j].Path
	})
	report.Summary.Pages = len(report.Pages)
	return report, nil
}

// WriteJSON publishes a stable, human-readable report at the caller's chosen path.
func WriteJSON(path string, report Report) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
