package plugins

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"

	xhtml "golang.org/x/net/html"
)

// Render profile names are shared across native Go, CLI and Python consumers.
// A profile controls renderer behavior, not a site build's entire lifecycle.
const (
	RenderProfileWeb                = "web"
	RenderProfileOfflineInteractive = "offline-interactive"
	RenderProfileOfflineStatic      = "offline-static"
)

// ResolveStandaloneProfile returns effective options and rejects a requested
// guarantee that this standalone renderer cannot currently fulfill.
//
// No profile silently downgrades a requested static rendering policy to client
// JavaScript. Existing empty-profile calls retain their previous behavior.
func ResolveStandaloneProfile(options StandaloneOptions) (StandaloneOptions, error) {
	switch options.Profile {
	case "", RenderProfileWeb:
		return options, nil
	case RenderProfileOfflineInteractive:
		if options.Unsafe {
			return options, fmt.Errorf("%s does not allow unsafe authored HTML", options.Profile)
		}
		// Bundled dependencies must not need a separately served asset endpoint.
		options.InlineAssets = true
		return options, nil
	case RenderProfileOfflineStatic:
		return options, fmt.Errorf("%s is not supported by the standalone renderer: Mermaid and Chart.js require client-side JavaScript; no static fallback is implemented", options.Profile)
	default:
		return options, fmt.Errorf("unknown render profile %q (expected web, offline-interactive, or offline-static)", options.Profile)
	}
}

const (
	standaloneTagStyle = "style"
	standaloneAttrSrc  = "src"
)

var remoteCSSURL = regexp.MustCompile(`(?i)(?:url\s*\(\s*['"]?\s*|@import\s+(?:url\s*\(\s*)?['"]?\s*)(?:https?:)?//`)

func isRemoteResourceURL(value string) bool {
	value = strings.TrimSpace(value)
	if strings.HasPrefix(value, "//") {
		return true
	}
	u, err := url.Parse(value)
	if err != nil {
		return true // Invalid resource URLs cannot be certified offline.
	}
	return u.Scheme == schemeHTTP || u.Scheme == schemeHTTPS
}

// checkStandaloneOfflineResources catches browser resource loads reachable from
// generated HTML and declared assets. It deliberately ignores ordinary <a href>
// links, since following a link is an explicit navigation rather than a load.
//
// This is a conservative first-pass resource audit, NOT yet a certification that
// vendor JavaScript cannot initiate arbitrary network requests. That broader
// browser verification remains tracked in #1854.
func checkStandaloneOfflineResources(result *StandaloneResult) error {
	for _, asset := range result.Assets {
		if isRemoteResourceURL(asset.URL) ||
			(asset.Kind == standaloneTagStyle && remoteCSSURL.MatchString(asset.Content)) {
			return fmt.Errorf("offline-interactive: asset %q has a remote dependency", asset.ID)
		}
	}
	if len(result.ExternalURLs) > 0 {
		return fmt.Errorf("offline-interactive: renderer declared external dependencies: %s",
			strings.Join(result.ExternalURLs, ", "))
	}
	return checkStandaloneHTMLResources(result.HTML)
}

func checkStandaloneHTMLResources(markup string) error {
	tokenizer := xhtml.NewTokenizer(strings.NewReader(markup))
	for {
		kind := tokenizer.Next()
		switch kind {
		case xhtml.ErrorToken:
			if err := tokenizer.Err(); err != nil && !errors.Is(err, io.EOF) {
				return fmt.Errorf("offline-interactive: parse generated HTML: %w", err)
			}
			return nil
		case xhtml.StartTagToken, xhtml.SelfClosingTagToken:
			token := tokenizer.Token()
			for _, attr := range token.Attr {
				if err := checkStandaloneHTMLResourceAttr(token.Data, attr); err != nil {
					return err
				}
			}
		}
	}
}

func checkStandaloneHTMLResourceAttr(tag string, attr xhtml.Attribute) error {
	if attr.Key == standaloneTagStyle && remoteCSSURL.MatchString(attr.Val) {
		return fmt.Errorf("offline-interactive: inline style has a remote resource")
	}
	if attr.Key == "srcset" {
		for _, image := range strings.Split(attr.Val, ",") {
			item := strings.Fields(strings.TrimSpace(image))
			if len(item) > 0 && isRemoteResourceURL(item[0]) {
				return fmt.Errorf("offline-interactive: %s srcset has a remote resource", tag)
			}
		}
		return nil
	}
	if attr.Key == standaloneAttrSrc || attr.Key == "poster" || attr.Key == "data" ||
		(attr.Key == "href" && tag != "a") {
		if isRemoteResourceURL(attr.Val) {
			return fmt.Errorf("offline-interactive: %s %s loads a remote resource", tag, attr.Key)
		}
	}
	return nil
}
