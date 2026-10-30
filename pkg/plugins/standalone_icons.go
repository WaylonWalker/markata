package plugins

import (
	"embed"
	"strings"

	xhtml "golang.org/x/net/html"
)

const standaloneScriptKind = "script"

// Pinned lucide-static 1.48.0; see standalone_icons/LICENSE.
//
//go:embed standalone_icons/*.svg standalone_icons/LICENSE
var standaloneIcons embed.FS

func expandStandaloneIcons(source string) string {
	if !strings.Contains(source, ":lucide") {
		return source
	}
	p := NewIconsPlugin()
	p.autoVendor = false
	p.icons = make(map[string]iconAsset)
	for _, match := range iconShortcodeRegex.FindAllStringSubmatch(source, -1) {
		name, ok := lucideShortcodeIconName(match[1])
		if !ok {
			continue
		}
		if _, exists := p.icons["lucide/"+name]; exists {
			continue
		}
		raw, err := standaloneIcons.ReadFile("standalone_icons/" + name + ".svg")
		if err != nil {
			continue
		}
		svg, ok := prepareIconSVG(string(raw))
		if !ok {
			continue
		}
		asset := iconAsset{name: "lucide/" + name, svg: svg}
		p.icons["lucide/"+name], p.icons["lucide-"+name] = asset, asset
	}
	var out strings.Builder
	z := xhtml.NewTokenizer(strings.NewReader(source))
	protected := 0
	for {
		kind := z.Next()
		if kind == xhtml.ErrorToken {
			break
		}
		raw := string(z.Raw())
		if kind == xhtml.StartTagToken || kind == xhtml.EndTagToken {
			name, _ := z.TagName()
			if standaloneProtectedTag(string(name)) {
				if kind == xhtml.StartTagToken {
					protected++
				} else if protected > 0 {
					protected--
				}
			}
		}
		if kind == xhtml.TextToken && protected == 0 {
			raw = p.replaceShortcodes(raw)
		}
		out.WriteString(raw)
	}
	return out.String()
}

func standaloneProtectedTag(name string) bool {
	switch name {
	case "code", "pre", standaloneScriptKind, standaloneTagStyle, "textarea", "svg":
		return true
	default:
		return false
	}
}
