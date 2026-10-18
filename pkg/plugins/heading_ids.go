package plugins

import (
	"bytes"
	"html"
	"strconv"

	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// HeadingIDTransformer regenerates auto heading IDs from the heading's
// visible text. goldmark's auto ID generator slugs the raw source line, so
// `## Recent [TIL](/til/)` becomes `recent-tiltil` because the link
// destination leaks into the ID. Headings with an explicit `{#id}` attribute
// are left untouched.
type HeadingIDTransformer struct{}

// Transform implements parser.ASTTransformer.
func (t *HeadingIDTransformer) Transform(doc *ast.Document, reader text.Reader, _ parser.Context) {
	source := reader.Source()
	idCounts := make(map[string]int)
	if err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		h, ok := n.(*ast.Heading)
		if !ok {
			return ast.WalkContinue, nil
		}
		idVal, ok := h.AttributeString("id")
		if !ok {
			return ast.WalkContinue, nil
		}
		id, ok := idVal.([]byte)
		if !ok {
			return ast.WalkContinue, nil
		}

		// Only rewrite IDs goldmark derived from the raw line; anything else
		// was set explicitly by the author.
		if hasExplicitAttributes(h, source) || !isGoldmarkAutoID(id, goldmarkHeadingSlug(headingRawLine(h, source))) {
			return ast.WalkContinue, nil
		}

		h.SetAttribute([]byte("id"), []byte(generateHeadingID(string(headingPlainText(h, source)), idCounts)))
		return ast.WalkContinue, nil
	}); err != nil {
		// The transformer callbacks do not return errors, but preserve the
		// parser's contract if a future callback does.
		return
	}
}

// generateHeadingID assigns IDs in source order using the same collision
// behavior as the table of contents.
func generateHeadingID(text string, idCounts map[string]int) string {
	id := models.Slugify(text)
	if id == "" {
		id = "heading"
	}

	count := idCounts[id]
	idCounts[id] = count + 1
	if count > 0 {
		// Match goldmark/GitHub numbering: intro, intro-1, intro-2, ...
		return id + "-" + strconv.Itoa(count)
	}
	return id
}

// isGoldmarkAutoID reports whether id is goldmark's auto ID for a heading
// whose raw-line slug is slug, including its "-N" collision suffixes.
func isGoldmarkAutoID(id, slug []byte) bool {
	if bytes.Equal(id, slug) {
		return true
	}
	if !bytes.HasPrefix(id, slug) || len(id) < len(slug)+2 || id[len(slug)] != '-' {
		return false
	}
	for _, c := range id[len(slug)+1:] {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// hasExplicitAttributes reports whether the heading's source line ends with
// an attribute block such as {#id}. goldmark trims that block from the
// heading's text segment, so it sits between the segment end and the newline.
func hasExplicitAttributes(h *ast.Heading, source []byte) bool {
	lines := h.Lines()
	if lines.Len() == 0 {
		return false
	}
	stop := lines.At(lines.Len() - 1).Stop
	if stop < 0 || stop > len(source) {
		return false
	}
	rest := source[stop:]
	if nl := bytes.IndexByte(rest, '\n'); nl >= 0 {
		rest = rest[:nl]
	}
	rest = bytes.TrimSpace(rest)
	return bytes.HasPrefix(rest, []byte("{")) && bytes.HasSuffix(rest, []byte("}"))
}

func headingRawLine(h *ast.Heading, source []byte) []byte {
	lines := h.Lines()
	if lines.Len() == 0 {
		return nil
	}
	seg := lines.At(lines.Len() - 1)
	return seg.Value(source)
}

// headingPlainText concatenates the text of all inline descendants.
func headingPlainText(h *ast.Heading, source []byte) []byte {
	var buf bytes.Buffer
	if err := ast.Walk(h, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch v := n.(type) {
		case *ast.Text:
			buf.Write(v.Segment.Value(source))
			if v.SoftLineBreak() || v.HardLineBreak() {
				buf.WriteByte(' ')
			}
		case *ast.String:
			// The typographer extension emits entities such as &ldquo;;
			// decode them so they do not leak into the slug as "ldquo".
			buf.WriteString(html.UnescapeString(string(v.Value)))
		}
		return ast.WalkContinue, nil
	}); err != nil {
		return buf.Bytes()
	}
	return buf.Bytes()
}

// goldmarkHeadingSlug mirrors goldmark's default IDs.Generate for headings,
// minus collision suffixes.
func goldmarkHeadingSlug(value []byte) []byte {
	value = util.TrimLeftSpace(value)
	value = util.TrimRightSpace(value)
	result := make([]byte, 0, len(value))
	for i := 0; i < len(value); {
		v := value[i]
		l := util.UTF8Len(v)
		i += int(l)
		if l != 1 {
			continue
		}
		switch {
		case util.IsAlphaNumeric(v):
			if 'A' <= v && v <= 'Z' {
				v += 'a' - 'A'
			}
			result = append(result, v)
		case util.IsSpace(v) || v == '-' || v == '_':
			result = append(result, '-')
		}
	}
	if len(result) == 0 {
		return []byte("heading")
	}
	return result
}
