package plugins

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"html"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	htmlparser "golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const (
	pageStyleHref = "href"
	cssStylesheetRel = "stylesheet"
)

var pageCSSURL = regexp.MustCompile(`(?i)url\(\s*["']?([^"')\s]+)`)

type pageStyle struct {
	start, end int
	css        []byte
	source     string
}

// Cleanup bundles final stylesheet bytes, preserving the cascade and script
// boundaries. It only touches explicitly selected pages in production builds.
func (p *CSSBundlePlugin) Cleanup(m *lifecycle.Manager) error {
	if !p.config.Enabled || len(p.config.Pages) == 0 || m.Config().Extra["fast_mode"] == true {
		return nil
	}
	updatedHTMLHashes := make(map[string]string)
	for _, page := range p.config.Pages {
		if !filepath.IsLocal(page) || strings.ToLower(filepath.Ext(page)) != ".html" {
			return fmt.Errorf("CSS bundle page must be an output-relative HTML path: %q", page)
		}
		path := filepath.Join(m.Config().OutputDir, page)
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading CSS bundle page %q: %w", page, err)
		}
		bundled, err := p.bundlePageStyles(content, m.Config().OutputDir)
		if err != nil {
			return fmt.Errorf("bundling page %q: %w", page, err)
		}
		if !bytes.Equal(content, bundled) {
			//nolint:gosec // Generated public HTML must be readable by the web server.
			if err := os.WriteFile(path, bundled, 0o644); err != nil {
				return fmt.Errorf("writing CSS bundle page %q: %w", page, err)
			}
			updatedHTMLHashes[filepath.Clean(page)] = buildcache.ContentHash(string(bundled))
		}
	}
	// This runs after static-assets CSS reconciliation so bundled pages read
	// real versioned stylesheets. Update ownership again for any pages changed
	// by bundling rather than leaving the next warm build with stale hashes.
	if err := refreshImageLibraryCSSOwnership(GetBuildCache(m), m.Config().OutputDir, updatedHTMLHashes); err != nil {
		return fmt.Errorf("updating bundled page ownership: %w", err)
	}
	return nil
}

type pageBundleWriter struct {
	content []byte
	output  string
	result  bytes.Buffer
	group   []pageStyle
	cursor  int
}

func (w *pageBundleWriter) flush() error {
	if len(w.group) < 2 {
		w.group = nil
		return nil
	}
	var css bytes.Buffer
	var sources []string
	for _, style := range w.group {
		css.Write(style.css)
		css.WriteByte('\n')
		if style.source != "" {
			sources = append(sources, style.source)
		}
	}
	name := fmt.Sprintf("page-%x.css", sha256.Sum256(css.Bytes()))
	dir := filepath.Join(w.output, "css")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	//nolint:gosec // Generated public stylesheet.
	if err := os.WriteFile(filepath.Join(dir, name), css.Bytes(), 0o644); err != nil {
		return err
	}
	for i, style := range w.group {
		w.result.Write(w.content[w.cursor:style.start])
		if i == 0 {
			fmt.Fprintf(&w.result, `<link rel="stylesheet" href="/css/%s" data-css-sources="%s">`, name, html.EscapeString(strings.Join(sources, " ")))
		}
		w.cursor = style.end
	}
	w.group = nil
	return nil
}

func (p *CSSBundlePlugin) bundlePageStyles(content []byte, output string) ([]byte, error) {
	tokens := htmlparser.NewTokenizer(bytes.NewReader(content))
	w := pageBundleWriter{content: content, output: output}
	offset := 0
	inHead := false
	for {
		kind := tokens.Next()
		start := offset
		offset += len(tokens.Raw())
		if kind == htmlparser.ErrorToken {
			if !errors.Is(tokens.Err(), io.EOF) {
				return nil, tokens.Err()
			}
			break
		}
		if kind != htmlparser.StartTagToken && kind != htmlparser.SelfClosingTagToken && kind != htmlparser.EndTagToken {
			continue
		}
		token := tokens.Token()
		if token.DataAtom == atom.Head {
			inHead = kind != htmlparser.EndTagToken
			if !inHead {
				if err := w.flush(); err != nil {
					return nil, err
				}
			}
			continue
		}
		if !inHead || kind == htmlparser.EndTagToken {
			continue
		}
		switch token.DataAtom {
		case atom.Script:
			if err := w.flush(); err != nil {
				return nil, err
			}
		case atom.Noscript, atom.Template:
			if err := w.flush(); err != nil {
				return nil, err
			}
			skipPageStyleContainer(tokens, &offset, token.DataAtom)
		case atom.Link, atom.Style:
			style, ok := p.readPageStyle(token, tokens, start, &offset, output)
			if ok {
				w.group = append(w.group, style)
			} else if err := w.flush(); err != nil {
				return nil, err
			}
		default:
		}
	}
	if err := w.flush(); err != nil {
		return nil, err
	}
	w.result.Write(content[w.cursor:])
	return w.result.Bytes(), nil
}

func skipPageStyleContainer(tokens *htmlparser.Tokenizer, offset *int, container atom.Atom) {
	depth := 1
	for depth > 0 {
		kind := tokens.Next()
		*offset += len(tokens.Raw())
		if kind == htmlparser.ErrorToken {
			return
		}
		if kind != htmlparser.StartTagToken && kind != htmlparser.EndTagToken {
			continue
		}
		if tokens.Token().DataAtom == container {
			if kind == htmlparser.StartTagToken {
				depth++
			} else {
				depth--
			}
		}
	}
}

func (p *CSSBundlePlugin) readPageStyle(token htmlparser.Token, tokens *htmlparser.Tokenizer, start int, offset *int, output string) (pageStyle, bool) {
	style := pageStyle{start: start, end: *offset}
	if token.DataAtom == atom.Link {
		css, ok := p.localPageStylesheet(token.Attr, output)
		style.css = css
		for _, attr := range token.Attr {
			if attr.Key == pageStyleHref {
				style.source = attr.Val
			}
		}
		return style, ok
	}
	var css bytes.Buffer
	for tokens.Next() == htmlparser.TextToken {
		css.Write(tokens.Raw())
		*offset += len(tokens.Raw())
	}
	*offset += len(tokens.Raw())
	style.end = *offset
	style.css = css.Bytes()
	return style, len(token.Attr) == 0 && safePageCSS(style.css)
}

func (p *CSSBundlePlugin) localPageStylesheet(attrs []htmlparser.Attribute, output string) ([]byte, bool) {
	var href, rel string
	for _, attr := range attrs {
		switch attr.Key {
		case pageStyleHref:
			href = attr.Val
		case "rel":
			rel = attr.Val
		case "type":
			if attr.Val != "text/css" {
				return nil, false
			}
		default:
			return nil, false
		}
	}
	u, err := url.Parse(href)
	if err != nil || rel != cssStylesheetRel || u.IsAbs() || u.Host != "" || !strings.HasPrefix(u.Path, "/") || u.RawQuery != "" || u.Fragment != "" {
		return nil, false
	}
	local := filepath.FromSlash(strings.TrimPrefix(u.Path, "/"))
	if !filepath.IsLocal(local) || p.isExcluded(local) {
		return nil, false
	}
	css, err := os.ReadFile(filepath.Join(output, local))
	return css, err == nil && safePageCSS(css)
}

// Relative URLs change meaning in a bundle. Leave those files untouched rather
// than guessing at CSS escapes, import semantics, or namespace placement.
func safePageCSS(css []byte) bool {
	lower := bytes.ToLower(css)
	for _, rule := range []string{"@import", "@charset", "@namespace"} {
		if bytes.Contains(lower, []byte(rule)) {
			return false
		}
	}
	for _, match := range pageCSSURL.FindAllSubmatch(css, -1) {
		value := string(match[1])
		u, err := url.Parse(value)
		if err != nil || (!u.IsAbs() && !strings.HasPrefix(value, "/") && !strings.HasPrefix(value, "#")) {
			return false
		}
	}
	return true
}
