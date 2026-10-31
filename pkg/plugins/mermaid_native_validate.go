//go:build !js

package plugins

import (
    "encoding/xml"
    "fmt"
    "io"
    "net/url"
    "regexp"
    "strings"
)

// Native SVG is incorporated into trusted site HTML. Keep the trust boundary
// here even if a diagram library currently escapes content itself.
const maxNativeMermaidSVGBytes = 16 << 20

var svgLocalURL = regexp.MustCompile(`(?i)url\(\s*#[a-z0-9_.:-]+\s*\)`)

func validateNativeMermaidSVG(data []byte) error {
    if len(data) == 0 || len(data) > maxNativeMermaidSVGBytes {
        return fmt.Errorf("native Mermaid SVG is empty or exceeds %d bytes", maxNativeMermaidSVGBytes)
    }
    dec := xml.NewDecoder(strings.NewReader(string(data)))
    foundRoot := false
    depth := 0
    for {
        tok, err := dec.Token()
        if err == io.EOF {break}
        if err != nil {return fmt.Errorf("native Mermaid SVG XML: %w",err)}
        switch t := tok.(type) {
        case xml.StartElement:
            name := strings.ToLower(t.Name.Local)
            if depth == 0 {
                if foundRoot || name != "svg" {
                    return fmt.Errorf("native Mermaid returned an unexpected SVG root")
                }
                foundRoot = true
            }
            depth++
            switch name {
            case "script", "foreignobject", "iframe", "object", "embed",
                "image", "use", "feimage", "animate", "animatetransform",
                "animatemotion", "set", "style", "audio", "video":
                return fmt.Errorf("native Mermaid SVG includes forbidden element %q",name)
            }
            if t.Name.Space != "" && t.Name.Space != "http://www.w3.org/2000/svg" {
                return fmt.Errorf("native Mermaid SVG includes foreign namespace")
            }
            for _,a := range t.Attr {
                key := strings.ToLower(a.Name.Local)
                value := strings.TrimSpace(a.Value)
                if len(key) > 2 && strings.HasPrefix(key,"on") {
                    return fmt.Errorf("native Mermaid SVG includes event attribute %q",key)
                }
                if key == "href" || key == "src" {
                    if !safeNativeMermaidURL(value) {
                        return fmt.Errorf("native Mermaid SVG includes unsafe URL")
                    }
                }
                // Only local SVG marker references may appear in CSS url().
                remaining := svgLocalURL.ReplaceAllString(value,"")
                if strings.Contains(strings.ToLower(remaining),"url(") {
                    return fmt.Errorf("native Mermaid SVG includes external CSS URL")
                }
            }
        case xml.EndElement:
            depth--
            if depth < 0 {return fmt.Errorf("native Mermaid SVG has an unbalanced element")}
        case xml.Directive:
            return fmt.Errorf("native Mermaid SVG contains an XML directive")
        case xml.ProcInst:
            // The renderer should not need XML processing instructions.
            return fmt.Errorf("native Mermaid SVG contains a processing instruction")
        case xml.CharData:
            if depth == 0 && strings.TrimSpace(string(t)) != "" {
                return fmt.Errorf("native Mermaid SVG has content outside its root")
            }
        }
    }
    if !foundRoot || depth != 0 {
        return fmt.Errorf("native Mermaid SVG has no valid root")
    }
    return nil
}

func safeNativeMermaidURL(value string) bool {
    if value == "" || strings.ContainsAny(value, "\r\n\t\\") ||
        strings.HasPrefix(value, "//") {
        return false
    }
    parsed,err := url.Parse(value)
    if err != nil {return false}
    switch strings.ToLower(parsed.Scheme) {
    case "https","http":
        return parsed.Host != "" && parsed.User == nil
    case "mailto":
        return parsed.Opaque != "" || parsed.Path != ""
    case "":
        return !strings.HasPrefix(value, "://")
    default:
        return false
    }
}
