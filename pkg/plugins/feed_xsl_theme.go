package plugins

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/templates"
)

const (
	feedXSLThemeStart   = "<!-- markata:theme-head -->"
	feedXSLThemeEnd     = "<!-- /markata:theme-head -->"
	feedXSLHeadTemplate = "partials/feed-xsl-head.html"
)

var (
	xslScriptPattern = regexp.MustCompile(`(?s)(<script\b[^>]*>)(.*?)(</script>)`)
	xslVoidPattern   = regexp.MustCompile(`(?i)<(link|meta)\b([^>]*?)\s*/?>`)
)

// themeFeedXSL replaces the markata:theme-head region of an XSL stylesheet
// with the rendered feed-xsl-head partial so RSS/Atom pages use the same
// palette, color mode, fontpack, text size and aesthetic as the site. The
// content is returned unchanged when the markers are missing or rendering
// fails, which keeps custom stylesheets working.
func (p *PublishFeedsPlugin) themeFeedXSL(content []byte, config *lifecycle.Config, templatesDir string) []byte {
	start := bytes.Index(content, []byte(feedXSLThemeStart))
	if start < 0 {
		return content
	}
	end := bytes.Index(content[start:], []byte(feedXSLThemeEnd))
	if end < 0 {
		return content
	}
	end += start + len(feedXSLThemeEnd)

	engine, err := p.getOrCreateEngine(templatesDir, feedThemeName(config))
	if err != nil || !engine.TemplateExists(feedXSLHeadTemplate) {
		return content
	}
	ctx := templates.Context{Config: ToModelsConfig(config), Extra: map[string]interface{}{}}
	head, err := engine.Render(feedXSLHeadTemplate, ctx)
	if err != nil {
		publishFeedsLog.Warnf("rendering %s for feed stylesheets failed: %v", feedXSLHeadTemplate, err)
		return content
	}

	out := make([]byte, 0, len(content)+len(head))
	out = append(out, content[:start]...)
	out = append(out, xmlSafeHead(head)...)
	out = append(out, content[end:]...)
	return out
}

// xmlSafeHead turns rendered HTML head markup into well-formed XML for
// embedding in an XSL stylesheet: script bodies are wrapped in CDATA and
// void elements are self-closed. The XSL html output method writes both back
// out as ordinary HTML.
func xmlSafeHead(head string) string {
	head = xslScriptPattern.ReplaceAllStringFunc(head, func(m string) string {
		parts := xslScriptPattern.FindStringSubmatch(m)
		body := parts[2]
		if strings.TrimSpace(body) == "" || strings.Contains(body, "<![CDATA[") {
			return m
		}
		body = strings.ReplaceAll(body, "]]>", "]]]]><![CDATA[>")
		return parts[1] + "<![CDATA[" + body + "]]>" + parts[3]
	})
	return xslVoidPattern.ReplaceAllString(head, "<$1$2 />")
}

// feedThemeName returns the configured theme name, defaulting to "default".
func feedThemeName(config *lifecycle.Config) string {
	themeName := ThemeDefault
	if config == nil || config.Extra == nil {
		return themeName
	}
	switch theme := config.Extra["theme"].(type) {
	case models.ThemeConfig:
		if theme.Name != "" {
			themeName = theme.Name
		}
	case map[string]interface{}:
		if name, ok := theme["name"].(string); ok && name != "" {
			themeName = name
		}
	case string:
		if theme != "" {
			themeName = theme
		}
	}
	return themeName
}
