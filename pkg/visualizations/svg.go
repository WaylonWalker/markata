package visualizations

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"unicode"
)

const MaxSVGBytes = 4 << 20 // 4 MiB per diagram until real-engine output budgets are established.

const svgNamespace = "http://www.w3.org/2000/svg"

// NewSVGAsset validates bounded SVG from a trusted backend renderer and
// supplies intrinsic dimensions plus an immutable content-addressed path.
// This is NOT a general-purpose sanitizer for untrusted SVG or a complete
// offline browser-network audit.
func NewSVGAsset(svg []byte) (Asset, error) {
	if len(svg) == 0 || len(svg) > MaxSVGBytes {
		return Asset{}, fmt.Errorf("SVG must contain 1..%d bytes", MaxSVGBytes)
	}
	width, height, err := inspectSVG(svg)
	if err != nil {
		return Asset{}, err
	}
	sum := sha256.Sum256(svg)
	sha := hex.EncodeToString(sum[:])
	return Asset{
		Path:   "assets/visualizations/" + sha + ".svg",
		SHA256: sha,
		SVG:    append([]byte(nil), svg...),
		Width:  width,
		Height: height,
	}, nil
}

type svgInspection struct {
	rootSeen bool
	depth    int
	inStyle  bool
	width    int
	height   int
}

func inspectSVG(data []byte) (width, height int, err error) {
	decoder := xml.NewDecoder(bytes.NewReader(data))
	var state svgInspection
	for {
		token, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return 0, 0, fmt.Errorf("invalid SVG XML: %w", err)
		}
		if err := state.accept(token); err != nil {
			return 0, 0, err
		}
	}
	if !state.rootSeen || state.depth != 0 {
		return 0, 0, errors.New("missing or incomplete SVG root")
	}
	return state.width, state.height, nil
}

func (state *svgInspection) accept(token xml.Token) error {
	switch item := token.(type) {
	case xml.Directive:
		return errors.New("SVG directives and DOCTYPE are not allowed")
	case xml.ProcInst:
		if !strings.EqualFold(item.Target, "xml") || state.rootSeen {
			return errors.New("SVG processing instructions are not allowed")
		}
	case xml.StartElement:
		return state.start(item)
	case xml.EndElement:
		if strings.EqualFold(item.Name.Local, "style") {
			state.inStyle = false
		}
		state.depth--
	case xml.CharData:
		if state.depth == 0 && strings.TrimSpace(string(item)) != "" {
			return errors.New("text outside SVG root")
		}
		if state.inStyle && unsafeCSS(string(item)) {
			return errors.New("SVG contains external or active CSS")
		}
	}
	return nil
}

func (state *svgInspection) start(item xml.StartElement) error {
	if state.depth == 0 {
		if state.rootSeen || item.Name.Local != "svg" {
			return errors.New("SVG must have one svg root")
		}
		if item.Name.Space != "" && item.Name.Space != svgNamespace {
			return errors.New("SVG root has an unsupported namespace")
		}
		state.rootSeen = true
		var err error
		state.width, state.height, err = rootDimensions(item.Attr)
		if err != nil {
			return err
		}
	}
	if err := validateSVGElement(item); err != nil {
		return err
	}
	if strings.EqualFold(item.Name.Local, "style") {
		state.inStyle = true
	}
	state.depth++
	return nil
}

func validateSVGElement(item xml.StartElement) error {
	if forbiddenElement(item.Name.Local) {
		return fmt.Errorf("SVG element %q is not allowed in static output", item.Name.Local)
	}
	for _, attr := range item.Attr {
		name := strings.ToLower(attr.Name.Local)
		value := strings.TrimSpace(attr.Value)
		if strings.HasPrefix(name, "on") || name == "base" || name == "src" {
			return fmt.Errorf("SVG attribute %q is not allowed", name)
		}
		if name == "href" && !strings.HasPrefix(value, "#") {
			return errors.New("SVG href must point to a local fragment")
		}
		if cssAttribute(name) && unsafeCSS(value) {
			return fmt.Errorf("SVG attribute %q contains external or active CSS", name)
		}
	}
	return nil
}

func cssAttribute(name string) bool {
	switch name {
	case "style", "fill", "stroke", "filter", "clip-path", "mask":
		return true
	}
	return false
}

func forbiddenElement(name string) bool {
	switch strings.ToLower(name) {
	case "script", "foreignobject", "iframe", "object", "embed", "audio", "video", "image", "feimage", "a", "animate", "animatetransform", "animatemotion", "set":
		return true
	}
	return false
}

func unsafeCSS(css string) bool {
	css = strings.ToLower(css)
	if strings.Contains(css, "@import") || strings.Contains(css, "expression(") || strings.Contains(css, "javascript:") {
		return true
	}
	for {
		i := strings.Index(css, "url(")
		if i < 0 {
			return false
		}
		css = css[i+4:]
		end := strings.IndexByte(css, ')')
		if end < 0 {
			return true
		}
		value := strings.Trim(css[:end], " \t\n\r\"'")
		if !strings.HasPrefix(value, "#") && !strings.HasPrefix(value, "data:") {
			return true
		}
		css = css[end+1:]
	}
}

func rootDimensions(attrs []xml.Attr) (width, height int, err error) {
	var rawWidth, rawHeight, viewBox string
	for _, attr := range attrs {
		switch attr.Name.Local {
		case "width":
			rawWidth = attr.Value
		case "height":
			rawHeight = attr.Value
		case "viewBox":
			viewBox = attr.Value
		}
	}
	width, widthErr := pixelDimension(rawWidth)
	height, heightErr := pixelDimension(rawHeight)
	if widthErr != nil || heightErr != nil {
		return 0, 0, errors.New("SVG width and height must be positive CSS pixels")
	}
	if viewBox == "" {
		if width > 0 && height > 0 {
			return width, height, nil
		}
		return 0, 0, errors.New("SVG needs positive dimensions or a four-number viewBox")
	}
	vw, vh, err := parseViewBox(viewBox)
	if err != nil {
		return 0, 0, err
	}
	switch {
	case width == 0 && height == 0:
		width, height = roundDimension(vw), roundDimension(vh)
	case width == 0:
		width = roundDimension(float64(height) * vw / vh)
	case height == 0:
		height = roundDimension(float64(width) * vh / vw)
	}
	if width == 0 || height == 0 {
		return 0, 0, errors.New("SVG dimensions out of range")
	}
	return width, height, nil
}

func parseViewBox(viewBox string) (width, height float64, err error) {
	fields := strings.FieldsFunc(viewBox, func(r rune) bool { return unicode.IsSpace(r) || r == ',' })
	if len(fields) != 4 {
		return 0, 0, errors.New("SVG viewBox must contain four numbers")
	}
	values := [4]float64{}
	for i, value := range fields {
		n, err := strconv.ParseFloat(value, 64)
		if err != nil || math.IsNaN(n) || math.IsInf(n, 0) {
			return 0, 0, errors.New("SVG viewBox values must be finite numbers")
		}
		values[i] = n
	}
	if values[2] <= 0 || values[3] <= 0 {
		return 0, 0, errors.New("SVG viewBox size must be positive")
	}
	return values[2], values[3], nil
}

func pixelDimension(s string) (int, error) {
	if s == "" {
		return 0, nil
	}
	s = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(s), "px"))
	value, err := strconv.ParseFloat(s, 64)
	if err != nil || !(value > 0) || math.IsInf(value, 0) {
		return 0, errors.New("invalid SVG dimension")
	}
	dimension := roundDimension(value)
	if dimension == 0 {
		return 0, errors.New("SVG dimension out of range")
	}
	return dimension, nil
}

func roundDimension(value float64) int {
	if math.IsNaN(value) || math.IsInf(value, 0) || value <= 0 || value > 100000 {
		return 0
	}
	return max(1, int(math.Round(value)))
}
