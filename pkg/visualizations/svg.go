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

const MaxSVGBytes = 4 << 20 // Bounded per-diagram SVG processing; configurable later.

// NewSVGAsset validates the static renderer output and supplies intrinsic
// dimensions plus an immutable, content-addressed asset path.
// It is deliberately restrictive: SVG from untrusted authors must be sanitized
// BEFORE reaching a renderer. This check is not a general-purpose SVG sanitizer.
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

func inspectSVG(data []byte) (int, int, error) {
	dec := xml.NewDecoder(bytes.NewReader(data))
	var rootSeen bool
	depth := 0
	var width, height int
	var inStyle bool
	for {
		token, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, 0, fmt.Errorf("invalid SVG XML: %w", err)
		}
		switch element := token.(type) {
		case xml.Directive:
			return 0, 0, errors.New("SVG directives and DOCTYPE are not allowed")
		case xml.ProcInst:
			if !strings.EqualFold(element.Target, "xml") {
				return 0, 0, errors.New("SVG processing instructions are not allowed")
			}
		case xml.StartElement:
			if depth == 0 {
				if rootSeen || element.Name.Local != "svg" {
					return 0, 0, errors.New("SVG must have one svg root")
				}
				rootSeen = true
				width, height, err = rootDimensions(element.Attr)
				if err != nil {
					return 0, 0, err
				}
			}
			depth++
			if strings.EqualFold(element.Name.Local, "style") {
				inStyle = true
			}
			if forbiddenElement(element.Name.Local) {
				return 0, 0, fmt.Errorf("SVG element %q is not allowed in static output", element.Name.Local)
			}
			for _, attr := range element.Attr {
				name := strings.ToLower(attr.Name.Local)
				value := strings.TrimSpace(attr.Value)
				if strings.HasPrefix(name, "on") || name == "base" || name == "src" {
					return 0, 0, fmt.Errorf("SVG attribute %q is not allowed", name)
				}
				if name == "href" && !strings.HasPrefix(value, "#") {
					return 0, 0, errors.New("SVG href must point to a local fragment")
				}
				if (name == "style" || name == "fill" || name == "stroke" || name == "filter" || name == "clip-path" || name == "mask") && unsafeCSS(value) {
					return 0, 0, fmt.Errorf("SVG attribute %q contains external or active CSS", name)
				}
			}
		case xml.EndElement:
			if strings.EqualFold(element.Name.Local, "style") {
				inStyle = false
			}
			depth--
		case xml.CharData:
			if depth == 0 && strings.TrimSpace(string(element)) != "" {
				return 0, 0, errors.New("text outside SVG root")
			}
			// CSS in <style> can otherwise trigger a network request even when
			// all href attributes are local. Scan CSS-bearing text conservatively.
			if inStyle && unsafeCSS(string(element)) {
				return 0, 0, errors.New("SVG contains external or active CSS")
			}
		}
	}
	if !rootSeen || depth != 0 {
		return 0, 0, errors.New("missing or incomplete SVG root")
	}
	return width, height, nil
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

func rootDimensions(attrs []xml.Attr) (int, int, error) {
	var rawWidth, rawHeight, viewBox string
	for _, a := range attrs {
		switch a.Name.Local {
		case "width":
			rawWidth = a.Value
		case "height":
			rawHeight = a.Value
		case "viewBox":
			viewBox = a.Value
		}
	}
	w, wErr := pixelDimension(rawWidth)
	h, hErr := pixelDimension(rawHeight)
	if wErr != nil || hErr != nil {
		return 0, 0, errors.New("SVG width and height must be positive CSS pixels")
	}
	if w > 0 && h > 0 {
		return w, h, nil
	}
	fields := strings.FieldsFunc(viewBox, func(r rune) bool { return unicode.IsSpace(r) || r == ',' })
	if len(fields) != 4 {
		return 0, 0, errors.New("SVG needs positive dimensions or a four-number viewBox")
	}
	vw, err1 := strconv.ParseFloat(fields[2], 64)
	vh, err2 := strconv.ParseFloat(fields[3], 64)
	if err1 != nil || err2 != nil || !(vw > 0 && vh > 0) || math.IsInf(vw, 0) || math.IsInf(vh, 0) {
		return 0, 0, errors.New("SVG viewBox size must be finite and positive")
	}
	if w == 0 && h == 0 {
		w, h = roundDimension(vw), roundDimension(vh)
	} else if w == 0 {
		w = roundDimension(float64(h) * vw / vh)
	} else {
		h = roundDimension(float64(w) * vh / vw)
	}
	if w == 0 || h == 0 {
		return 0, 0, errors.New("SVG dimensions out of range")
	}
	return w, h, nil
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
	d := roundDimension(value)
	if d == 0 {
		return 0, errors.New("SVG dimension out of range")
	}
	return d, nil
}

func roundDimension(value float64) int {
	if math.IsNaN(value) || math.IsInf(value, 0) || value <= 0 || value > 100000 {
		return 0
	}
	return max(1, int(math.Round(value)))
}
