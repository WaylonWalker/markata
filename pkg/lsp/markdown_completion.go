package lsp

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

type markdownSuggestion struct{ label, insert, detail string }

// Keep container vocabulary aligned with the editor command catalog and
// spec/spec/MARKDOWN_CONTAINERS.md. The palette remains the source for full snippets.
var containerClasses = []markdownSuggestion{
	{"admonition", "admonition", "Callout container"},
	{"custom-class", "custom-class", "Custom site container class"},
	{"road", "road", "Road diagram"},
	{"md-box", "md-box", "Bounded box"}, {"md-card", "md-card", "Article card"},
	{"md-checklist", "md-checklist", "Interactive checklist"}, {"md-columns", "md-columns", "Card columns"},
	{"md-connector", "md-connector", "Diagram connector"}, {"md-delta", "md-delta", "Before and after metric"},
	{"md-divider", "md-divider", "Layout divider"}, {"md-feature", "md-feature", "Icon feature row"},
	{"md-feature-body", "md-feature-body", "Feature text"}, {"md-feature-icon", "md-feature-icon", "Feature icon"},
	{"md-grid", "md-grid", "Responsive card grid"}, {"md-icon", "md-icon", "Icon tile"},
	{"md-layout", "md-layout", "Layout group"}, {"md-meter", "md-meter", "Measurement meter"},
	{"md-progress", "md-progress", "Progress gauge"}, {"md-section", "md-section", "Story section"},
	{"md-source", "md-source", "Markdown source disclosure"}, {"md-sparkline", "md-sparkline", "Trend sparkline"},
	{"md-stack", "md-stack", "Vertical stack"},
	{"md-tab", "md-tab", "Native tab panel"}, {"md-tabs", "md-tabs", "Native accessible tabs"},
	{"md-stat", "md-stat", "Metric statistic"}, {"md-status", "md-status", "Status badge"},
	{"wa-animated-image", "wa-animated-image", "Web Awesome animated image"}, {"wa-badge", "wa-badge", "Web Awesome badge"},
	{"wa-carousel", "wa-carousel", "Web Awesome carousel"}, {"wa-comparison", "wa-comparison", "Web Awesome image comparison"},
	{"wa-copy", "wa-copy", "Web Awesome copy button"}, {"wa-details", "wa-details", "Web Awesome details"},
	{"wa-qr", "wa-qr", "Web Awesome QR code"}, {"wa-tab", "wa-tab", "Web Awesome tab"},
	{"wa-tabs", "wa-tabs", "Web Awesome tabs"}, {"wa-tag", "wa-tag", "Web Awesome tag"},
	{"wa-tooltip", "wa-tooltip", "Web Awesome tooltip"},
}

var roadStyles = []markdownSuggestion{
	{"constellation", "constellation", "Road style"}, {"cycle", "cycle", "Road style"},
	{"editorial", "editorial", "Prose-led road"}, {"git", "git", "Road style"},
	{"reveal", "reveal", "Editorial entrance modifier"}, {"snake", "snake", "Road style"},
	{"steps", "steps", "Road style"}, {"story", "story", "Editorial typography modifier"},
	{"timeline", "timeline", "Road style"}, {"winding", "winding", "Road style"},
}

var containerAttributes = []markdownSuggestion{
	{".class", ".class", "Add a CSS class"},
	{"#id", "#id", "Set an element ID"},
	{"data-align", "data-align", "Center or align card content"},
	{"data-columns", "data-columns", "Grid columns: 2, 3, or 4"},
	{"data-direction", "data-direction", "Connector direction"},
	{"data-gap", "data-gap", "Grid or stack spacing"},
	{"data-label", "data-label", "Accessible label"},
	{"data-max", "data-max", "Meter maximum"},
	{"data-size", "data-size", "Icon size"},
	{"data-tone", "data-tone", "Status color"},
	{"data-value", "data-value", "Metric or meter value"},
	{"caption", "caption", "Image comparison caption"},
	{"content", "content", "Tooltip content"},
	{"label", "label", "Tab label"},
	{"navigation", "navigation", "Carousel navigation"},
	{"pagination", "pagination", "Carousel pagination"},
	{"position", "position", "Image comparison position"},
	{"summary", "summary", "Details summary"},
	{"variant", "variant", "Web Awesome component variant"},
}

var attributeValues = map[string][]markdownSuggestion{
	"data-align":     {{"center", "center", "Align content"}},
	"data-columns":   {{"2", "2", "Two columns"}, {"3", "3", "Three columns"}, {"4", "4", "Four columns"}},
	"data-direction": {{"both", "both", "Two-way connector"}, {"forward", "forward", "Forward connector"}, {"reverse", "reverse", "Reverse connector"}},
	"data-gap":       {{"lg", "lg", "Large spacing"}, {"sm", "sm", "Small spacing"}},
	"data-size":      {{"lg", "lg", "Large icon"}, {"sm", "sm", "Small icon"}},
	"data-tone":      {{"danger", "danger", "Danger"}, {"info", "info", "Information"}, {"success", "success", "Success"}, {"warning", "warning", "Warning"}},
	"navigation":     {{"true", "true", "Show carousel navigation"}},
	"pagination":     {{"true", "true", "Show carousel pagination"}},
	"position":       {{"50", "50", "Center comparison position"}},
	"variant":        {{"brand", "brand", "Brand style"}, {"success", "success", "Success style"}},
}

var mermaidDiagramKinds = []markdownSuggestion{
	{"classDiagram", "classDiagram", "Class diagram"}, {"erDiagram", "erDiagram", "Entity relationship diagram"},
	{"flowchart", "flowchart", "Flowchart"}, {"gantt", "gantt", "Project timeline"},
	{"gitGraph", "gitGraph", "Git history"}, {"graph", "graph", "Flowchart"},
	{"mindmap", "mindmap", "Mind map"}, {"pie", "pie", "Pie chart"},
	{"sequenceDiagram", "sequenceDiagram", "Sequence diagram"}, {"stateDiagram-v2", "stateDiagram-v2", "State diagram"},
	{"timeline", "timeline", "Timeline"},
}

var mermaidDirections = []markdownSuggestion{
	{"BT", "BT", "Bottom to top"}, {"LR", "LR", "Left to right"},
	{"RL", "RL", "Right to left"}, {"TB", "TB", "Top to bottom"}, {"TD", "TD", "Top down"},
}

var containerOpenLine = regexp.MustCompile(`^\s*:{3,}\s+(.*)$`)
var mermaidFence = regexp.MustCompile(`^\s*` + "```" + `\s*mermaid\s*$`)
var markdownFenceClose = regexp.MustCompile(`^\s*` + "```" + `\s*$`)
var mermaidDirectionLine = regexp.MustCompile(`^(?:flowchart|graph)\s+(\w*)$`)
var mermaidNodeShape = regexp.MustCompile(`\w+\s*([\[{(])([^\]})]*)$`)
var mermaidClassDef = regexp.MustCompile(`^\s*classDef\s+\S+\s+(.+)$`)

func (s *Server) handleMarkdownCompletion(msg *Message, params CompletionParams, lines []string, line string, col int) (bool, error) {
	if items, ok := containerCompletions(line, col, params); ok {
		return true, s.sendCompletion(msg.ID, params, &CompletionList{Items: items})
	}
	if items, ok := attributeCompletions(line, col, params); ok {
		return true, s.sendCompletion(msg.ID, params, &CompletionList{Items: items})
	}
	if items, ok := mermaidCompletions(lines, params.Position.Line, line, col, params); ok {
		return true, s.sendCompletion(msg.ID, params, &CompletionList{Items: items})
	}
	return false, nil
}

func containerCompletions(line string, col int, params CompletionParams) ([]CompletionItem, bool) {
	if col < 0 || col > len(line) {
		return nil, false
	}
	match := containerOpenLine.FindStringSubmatchIndex(line[:col])
	if match == nil {
		return nil, false
	}
	opener := line[match[2]:col]
	if strings.Contains(opener, "{") {
		return nil, false
	}
	fields := strings.Fields(opener)
	start := match[2]
	prefix := ""
	trailingSpace := strings.HasSuffix(opener, " ") || strings.HasSuffix(opener, "\t")
	if len(fields) > 0 && !trailingSpace {
		prefix = fields[len(fields)-1]
		start = col - len(prefix)
	}
	suggestions := containerClasses
	// Suggest modifiers only after an explicit road or grid class.
	if len(fields) > 0 && fields[0] == "road" && (len(fields) > 1 || trailingSpace) {
		suggestions = roadStyles
	}
	if len(fields) > 0 && fields[0] == "md-grid" && (len(fields) > 1 || trailingSpace) {
		suggestions = []markdownSuggestion{{"pending", "pending", "Cards that are not ready"}, {"rail", "rail", "Overview card rail"}}
	}
	return filterSuggestions(suggestions, prefix, params, start, col), true
}

func attributeCompletions(line string, col int, params CompletionParams) ([]CompletionItem, bool) {
	if col < 0 || col > len(line) {
		return nil, false
	}
	open := strings.LastIndex(line[:col], "{")
	if open < 0 || strings.Contains(line[open+1:col], "}") {
		return nil, false
	}
	before := line[open+1 : col]
	tokenStart := 0
	var quote rune
	for i, char := range before {
		if quote != 0 {
			if char == quote {
				quote = 0
			}
			continue
		}
		if char == '"' || char == '\'' {
			quote = char
		} else if char == ' ' || char == '\t' {
			tokenStart = i + 1
		}
	}
	current := before[tokenStart:]
	prefix := ""
	if eq := strings.LastIndex(current, "="); eq >= 0 {
		key := strings.Trim(strings.TrimSpace(current[:eq]), "\"'")
		valuePrefix := strings.TrimLeft(current[eq+1:], " \t\"'")
		valuePrefix = strings.TrimRight(valuePrefix, "\"'")
		values, ok := attributeValues[key]
		if !ok {
			return nil, false
		}
		start := col - len(valuePrefix)
		return filterSuggestions(values, valuePrefix, params, start, col), true
	}
	for current != "" && current[len(current)-1] != ' ' && current[len(current)-1] != '\t' {
		prefix = current[len(current)-1:] + prefix
		current = current[:len(current)-1]
	}
	// Include familiar authoring attributes from the shared command catalog.
	// Unknown container attributes remain valid; this is a suggestion list.
	return filterSuggestions(containerAttributes, prefix, params, col-len(prefix), col), true
}

func mermaidCompletions(lines []string, lineNo int, line string, col int, params CompletionParams) ([]CompletionItem, bool) {
	if lineNo < 0 || lineNo >= len(lines) || col < 0 || col > len(line) {
		return nil, false
	}
	fenceStart, inDiagram := mermaidFenceStart(lines, lineNo)
	if !inDiagram {
		return nil, false
	}
	trimmed := strings.TrimLeft(line[:col], " \t")
	if items, ok := mermaidHeaderCompletions(lines, fenceStart, lineNo, trimmed, col, params); ok {
		return items, true
	}
	if match := mermaidDirectionLine.FindStringSubmatchIndex(trimmed); match != nil {
		prefix := trimmed[match[2]:]
		return filterSuggestions(mermaidDirections, prefix, params, col-len(prefix), col), true
	}
	if match := mermaidNodeShape.FindStringSubmatchIndex(trimmed); match != nil {
		opening := trimmed[match[2]:match[3]]
		shapes := mermaidNodeClosers(opening)
		return filterSuggestions(shapes, "", params, col, col), true
	}
	if match := mermaidClassDef.FindStringSubmatchIndex(trimmed); match != nil {
		prefix := trimmed[match[2]:]
		attrs := []markdownSuggestion{{"color", "color:", "Text color"}, {"fill", "fill:", "Fill color"}, {"stroke", "stroke:", "Border color"}, {"stroke-width", "stroke-width:", "Border width"}}
		return filterSuggestions(attrs, prefix, params, col-len(prefix), col), true
	}
	return nil, false
}

func mermaidNodeClosers(opening string) []markdownSuggestion {
	switch opening {
	case "(":
		return []markdownSuggestion{{")", ")", "Round node"}, {"))", "))", "Stadium node"}}
	case "[":
		return []markdownSuggestion{{"]", "]", "Rectangle node"}, {"])", "])", "Cylinder node"}}
	case "{":
		return []markdownSuggestion{{"}", "}", "Hexagon node"}}
	default:
		return nil
	}
}

func mermaidFenceStart(lines []string, lineNo int) (int, bool) {
	start := -1
	for i := 0; i <= lineNo; i++ {
		if start < 0 && mermaidFence.MatchString(lines[i]) {
			start = i
		} else if start >= 0 && markdownFenceClose.MatchString(lines[i]) {
			start = -1
		}
	}
	return start, start >= 0
}

func mermaidHeaderCompletions(lines []string, fenceStart, lineNo int, trimmed string, col int, params CompletionParams) ([]CompletionItem, bool) {
	prefix := strings.TrimSpace(trimmed)
	if lineNo == fenceStart+1 {
		if prefix == "" {
			return mermaidHeadItems(params, col), true
		}
		if !strings.Contains(prefix, " ") {
			items := filterSuggestions(mermaidDiagramKinds, prefix, params, col-len(prefix), col)
			if len(items) > 0 {
				return items, true
			}
		}
	}
	if prefix != "" {
		return nil, false
	}
	for i := fenceStart + 1; i < lineNo; i++ {
		if strings.TrimSpace(lines[i]) != "" {
			return nil, false
		}
	}
	return mermaidHeadItems(params, col), true
}

func mermaidHeadItems(params CompletionParams, col int) []CompletionItem {
	items := filterSuggestions(mermaidDiagramKinds, "", params, col, col)
	// Keep a helpful flowchart starter near the top.
	items = append([]CompletionItem{{Label: "flowchart TD", Kind: CompletionItemKindSnippet, Detail: "Top-down flowchart", InsertText: "flowchart TD\n  Start --> End", FilterText: "flowchart diagram starter", SortText: "00000"}}, items...)
	return items
}

func filterSuggestions(suggestions []markdownSuggestion, prefix string, params CompletionParams, start, end int) []CompletionItem {
	prefix = strings.ToLower(prefix)
	items := make([]CompletionItem, 0, len(suggestions))
	for i, suggestion := range suggestions {
		if prefix != "" && !strings.HasPrefix(strings.ToLower(suggestion.label), prefix) {
			continue
		}
		items = append(items, CompletionItem{Label: suggestion.label, Kind: CompletionItemKindKeyword, Detail: suggestion.detail, InsertText: suggestion.insert, FilterText: suggestion.label + " " + suggestion.detail, SortText: fmt.Sprintf("%05d", i), TextEdit: &TextEdit{Range: Range{Start: Position{Line: params.Position.Line, Character: start}, End: Position{Line: params.Position.Line, Character: end}}, NewText: suggestion.insert}})
	}
	sort.SliceStable(items, func(i, j int) bool { return items[i].Label < items[j].Label })
	return items
}
