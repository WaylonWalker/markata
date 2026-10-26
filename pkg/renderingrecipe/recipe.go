// Package renderingrecipe compiles normalized rendering semantics into portable,
// deterministic assets. Consumers receive the bundle; they do not draw it.
package renderingrecipe

import (
	"bytes"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/WaylonWalker/markata/pkg/renderingcontract"
	jsoncanonicalizer "github.com/cyberphone/json-canonicalization/go/src/webpki.org/jsoncanonicalizer"
)

const motifSourceID = "motif-block-w-v1"
const (
	recipeNone          = "none"
	recipeSVGMediaType  = "image/svg+xml"
	recipeLayerSandwich = "sandwich"
)

//go:embed testdata/*.json testdata/recipe-sources/*.json
var testdata embed.FS

const (
	SchemaVersion   = "rendering-recipe-v1"
	ContractVersion = 1
	semanticDomain  = "rendering-semantic-v1\x00"
	recipeDomain    = "rendering-recipe-v1\x00"
)

type Theme struct {
	ContractVersion int                          `json:"contract_version"`
	Palette         string                       `json:"palette"`
	Aesthetic       string                       `json:"aesthetic"`
	Fontpack        string                       `json:"fontpack"`
	Texture         renderingcontract.Dial       `json:"texture"`
	HeadingTexture  renderingcontract.Dial       `json:"heading_texture"`
	Motif           renderingcontract.MotifState `json:"motif"`
	Variables       map[string]string            `json:"variables,omitempty"`
}

type Asset struct {
	Path      string `json:"path"`
	SHA256    string `json:"sha256"`
	MediaType string `json:"media_type"`
	ViewBox   string `json:"view_box"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

type Pass struct {
	ID         string `json:"id"`
	Role       string `json:"role"`
	Mode       string `json:"mode"`
	Asset      string `json:"asset"`
	MediaType  string `json:"media_type"`
	ViewBox    string `json:"view_box"`
	Repeat     string `json:"repeat"`
	ScaleMilli int    `json:"scale_milli"`
	// ScaleMilli is the repeat tile size multiplier in thousandths. Consumers
	// apply it to the asset ViewBox dimensions in CSS pixels.
	Scope        string `json:"scope,omitempty"`
	OpacityMilli *int   `json:"opacity_milli,omitempty"`
	Paint        string `json:"paint,omitempty"`
	SizeMilli    int    `json:"size_milli,omitempty"`
	GapMilli     int    `json:"gap_milli,omitempty"`
}

type FontAsset struct {
	Role          string   `json:"role"`
	Family        string   `json:"family"`
	AssetPath     string   `json:"asset_path"`
	SHA256        string   `json:"sha256"`
	Weight        int      `json:"weight"`
	Style         string   `json:"style"`
	FallbackStack []string `json:"fallback_stack"`
}

type Manifest struct {
	SchemaVersion   string      `json:"schema_version"`
	ContractVersion int         `json:"contract_version"`
	SemanticHash    string      `json:"semantic_hash"`
	RecipeHash      string      `json:"recipe_hash,omitempty"`
	Assets          []Asset     `json:"assets"`
	Passes          []Pass      `json:"passes"`
	Fonts           []FontAsset `json:"fonts"`
}

type Bundle struct {
	Manifest Manifest
	Assets   map[string][]byte
}

type recipeSource struct {
	SchemaVersion string `json:"schema_version"`
	ID            string `json:"id"`
	Kind          string `json:"kind"`
	Primitive     string `json:"primitive"`
	ViewBox       string `json:"view_box"`
	Path          string `json:"path"`
	Columns       int    `json:"columns"`
	Rows          int    `json:"rows"`
	Seed          int    `json:"seed"`
}

func loadRecipeSources() (map[string]recipeSource, error) {
	result := make(map[string]recipeSource)
	for _, id := range []string{"texture-screenprint-v1", "heading-splatter-v1", "motif-block-w-v1"} {
		data, err := testdata.ReadFile("testdata/recipe-sources/" + id + ".json")
		if err != nil {
			return nil, err
		}
		var raw map[string]any
		if err := decodeStrict(data, &raw); err != nil {
			return nil, fmt.Errorf("%s: %w", id, err)
		}
		allowed := []string{"schema_version", "id", "kind", "primitive", "view_box", "seed"}
		if id == motifSourceID {
			allowed = append(allowed, "columns", "rows", "path")
		}
		if err := rejectKeys(raw, allowed...); err != nil {
			return nil, fmt.Errorf("%s: %w", id, err)
		}
		var source recipeSource
		if err := decodeMap(raw, &source); err != nil {
			return nil, fmt.Errorf("%s: %w", id, err)
		}
		if source.SchemaVersion != "rendering-recipe-source-v1" || source.ID != id {
			return nil, fmt.Errorf("invalid recipe source %s", id)
		}
		if source.Seed < 0 || source.Primitive != "explicit-circles" && id != motifSourceID || source.Primitive != "explicit-path" && id == motifSourceID {
			return nil, fmt.Errorf("invalid recipe source geometry %s", id)
		}
		if id == motifSourceID && (source.Columns <= 0 || source.Rows <= 0 || source.Columns > 64 || source.Rows > 64) {
			return nil, fmt.Errorf("invalid motif source grid")
		}
		if source.ViewBox == "" {
			return nil, fmt.Errorf("empty recipe source view_box %s", id)
		}
		result[id] = source
	}
	return result, nil
}

func LoadCanonicalTheme() (Theme, error) {
	data, err := testdata.ReadFile("testdata/canonical-theme.json")
	if err != nil {
		return Theme{}, err
	}
	var raw map[string]any
	if err := decodeStrict(data, &raw); err != nil {
		return Theme{}, err
	}
	return normalize(raw)
}

//nolint:gocyclo // Validation order is part of the canonical recipe contract.
func normalize(raw map[string]any) (Theme, error) {
	if raw == nil {
		return Theme{}, errors.New("theme must be an object")
	}
	if nested, ok := raw["theme"].(map[string]any); ok {
		if err := rejectKeys(raw, "theme"); err != nil {
			return Theme{}, err
		}
		raw = nested
	}
	if err := rejectKeys(raw, "contract_version", "palette", "aesthetic", "fontpack", "texture", "heading_texture", "motif", "variables"); err != nil {
		return Theme{}, err
	}
	contract, err := renderingcontract.Load()
	if err != nil {
		return Theme{}, err
	}
	if value, ok := raw["contract_version"]; ok {
		if !isNumber(value) {
			return Theme{}, errors.New("contract_version must be an integer")
		}
	}
	palette, err := requiredString(raw, "palette", contract.Defaults.Palette)
	if err != nil {
		return Theme{}, err
	}
	aesthetic, err := requiredString(raw, "aesthetic", contract.Defaults.Aesthetic)
	if err != nil {
		return Theme{}, err
	}
	fontpack, err := requiredString(raw, "fontpack", contract.Defaults.Fontpack)
	if err != nil {
		return Theme{}, err
	}
	if raw["contract_version"] != nil && contractValue(raw["contract_version"]) != ContractVersion {
		return Theme{}, fmt.Errorf("unsupported contract_version")
	}
	if !enum(contract.Palettes, palette) {
		return Theme{}, fmt.Errorf("unknown palette %q", palette)
	}
	if !contains(contract.Enums["aesthetics"], aesthetic) {
		return Theme{}, fmt.Errorf("unknown aesthetic %q", aesthetic)
	}
	if !contains(contract.Enums["fontpacks"], fontpack) {
		return Theme{}, fmt.Errorf("unknown fontpack %q", fontpack)
	}
	texture, err := dialValue(raw, "texture", contract.Defaults.Texture)
	if err != nil {
		return Theme{}, err
	}
	heading, err := dialValue(raw, "heading_texture", contract.Defaults.HeadingTexture)
	if err != nil {
		return Theme{}, err
	}
	motif, err := motifValue(raw, contract.Defaults.Motif)
	if err != nil {
		return Theme{}, err
	}
	variables := map[string]string{}
	if values, ok := raw["variables"].(map[string]any); ok {
		for key, value := range values {
			stringValue, ok := value.(string)
			if !ok {
				return Theme{}, fmt.Errorf("variables.%s must be a string", key)
			}
			variables[key] = stringValue
		}
	} else if _, exists := raw["variables"]; exists {
		return Theme{}, errors.New("variables must be an object")
	}
	theme := Theme{ContractVersion: 1, Palette: palette, Aesthetic: aesthetic, Fontpack: fontpack, Texture: texture, HeadingTexture: heading, Motif: motif, Variables: variables}
	if err := validateTheme(theme); err != nil {
		return Theme{}, err
	}
	return theme, nil
}

func requiredString(raw map[string]any, key, fallback string) (string, error) {
	value, ok := raw[key]
	if !ok {
		return fallback, nil
	}
	result, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("%s must be a string", key)
	}
	return result, nil
}

func contractValue(v any) int {
	n, err := strconv.Atoi(fmt.Sprint(v))
	if err != nil {
		return 0
	}
	return n
}
func isNumber(value any) bool {
	switch value.(type) {
	case json.Number, float64, int, int64:
		return true
	default:
		return false
	}
}
func contains(values []string, value string) bool {
	for _, item := range values {
		if item == value {
			return true
		}
	}
	return false
}
func enum(values []renderingcontract.Palette, value string) bool {
	for _, item := range values {
		if item.ID == value {
			return true
		}
	}
	return false
}

func dialValue(raw map[string]any, key string, fallback renderingcontract.Dial) (renderingcontract.Dial, error) {
	dial := fallback
	if rawValue, exists := raw[key]; exists {
		value, ok := rawValue.(map[string]any)
		if !ok {
			return dial, fmt.Errorf("%s must be an object", key)
		}
		if err := rejectKeys(value, "kind", "color_mix", "scale", "scope", "glyph", "size", "gap", "row_offset", "wobble", "scatter", "layer", "color", "url"); err != nil {
			return dial, fmt.Errorf("%s: %w", key, err)
		}
		if err := decodeMap(value, &dial); err != nil {
			return dial, fmt.Errorf("%s: %w", key, err)
		}
		for _, field := range []string{"kind", "scope", "glyph", "size", "gap", "layer", "color", "url"} {
			if v, exists := value[field]; exists {
				if _, ok := v.(string); !ok {
					return dial, fmt.Errorf("%s.%s must be a string", key, field)
				}
			}
		}
		for _, field := range []string{"scale", "row_offset", "wobble", "scatter"} {
			if v, exists := value[field]; exists {
				if !isNumber(v) {
					return dial, fmt.Errorf("%s.%s must be a number", key, field)
				}
			}
		}
		if _, ok := value["color_mix"]; ok {
			if !isNumber(value["color_mix"]) {
				return dial, fmt.Errorf("%s.color_mix must be a number", key)
			}
		}
	}
	dial.ColorMix = fixed(dial.ColorMix, 0, 1)
	dial.Scale = fixed(dial.Scale, .25, 3)
	return dial, nil
}

func rejectKeys(value map[string]any, allowed ...string) error {
	set := make(map[string]struct{}, len(allowed))
	for _, key := range allowed {
		set[key] = struct{}{}
	}
	for key := range value {
		if _, ok := set[key]; !ok {
			return fmt.Errorf("obsolete or unknown recipe field %q", key)
		}
	}
	return nil
}

func motifValue(raw map[string]any, fallback renderingcontract.Dial) (renderingcontract.MotifState, error) {
	if rawValue, exists := raw["motif"]; exists {
		if _, ok := rawValue.(map[string]any); !ok {
			return renderingcontract.MotifState{}, errors.New("motif must be an object")
		}
	}
	if value, ok := raw["motif"].(map[string]any); ok {
		if err := rejectKeys(value, "kind", "color_mix", "glyph", "size", "gap", "row_offset", "wobble", "scatter", "layer", "color", "url"); err != nil {
			return renderingcontract.MotifState{}, fmt.Errorf("motif: %w", err)
		}
	}
	dial, err := dialValue(raw, "motif", fallback)
	if err != nil {
		return renderingcontract.MotifState{}, err
	}
	return renderingcontract.MotifState{Kind: dial.Kind, ColorMix: dial.ColorMix, Glyph: dial.Glyph, Size: dial.Size, Gap: dial.Gap, RowOffset: fixed(dial.RowOffset, 0, 1), Wobble: fixed(dial.Wobble, 0, 1), Scatter: fixed(dial.Scatter, 0, 1), Layer: dial.Layer, Color: dial.Color, URL: dial.URL}, nil
}

func fixed(value, minimum, maximum float64) float64 {
	value = math.Max(minimum, math.Min(maximum, value))
	return math.Round(value*1000) / 1000
}

func SemanticHash(theme Theme) (string, error) {
	b, err := canonicalJSON(theme)
	if err != nil {
		return "", err
	}
	return digest([]byte(semanticDomain), b), nil
}

// AssetHash is the digest of the exact standalone bytes. Domain separation is
// reserved for semantic and manifest identities; asset digests are portable
// content addresses and must match ordinary SHA-256 tooling.
func AssetHash(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func Compile(theme Theme) (Bundle, error) {
	if err := validateTheme(theme); err != nil {
		return Bundle{}, err
	}
	sources, err := loadRecipeSources()
	if err != nil {
		return Bundle{}, err
	}
	semantic, err := SemanticHash(theme)
	if err != nil {
		return Bundle{}, err
	}
	assets := map[string][]byte{}
	add := func(path string, data []byte, media, view string, width, height int) Asset {
		assets[path] = append([]byte(nil), data...)
		return Asset{Path: path, SHA256: AssetHash(data), MediaType: media, ViewBox: view, Width: width, Height: height}
	}
	colors := paletteColors(theme.Palette)
	surface := add("assets/surface-screenprint-v1.svg", surfaceSVG(theme.Texture.ColorMix, colors.background, colors.ink, sources["texture-screenprint-v1"]), "image/svg+xml", "0 0 180 180", 180, 180)
	heading := add("assets/heading-splatter-v1.svg", headingSVG(theme.HeadingTexture.ColorMix, sources["heading-splatter-v1"]), "image/svg+xml", "0 0 180 180", 180, 180)
	motifSource := sources["motif-block-w-v1"]
	motifView := motifSource.ViewBox
	motifWidth, motifHeight, err := motifFieldDimensions(theme.Motif, motifSource)
	if err != nil {
		return Bundle{}, err
	}
	motif := add("assets/motif-block-w-v1.svg", motifSVG(theme.Motif, motifRoleColor(theme.Motif.Color, colors), colors.background, motifSource), "image/svg+xml", motifView, motifWidth, motifHeight)
	fonts := canonicalFonts()
	passes := []Pass{}
	if theme.Motif.Layer == "under" || theme.Motif.Layer == recipeLayerSandwich {
		passes = append(passes, Pass{ID: "motif-under", Role: "motif-under", Mode: "image", Asset: motif.Path, MediaType: motif.MediaType, ViewBox: motif.ViewBox, Repeat: "repeat", ScaleMilli: 1000})
	}
	if theme.Texture.Kind != recipeNone {
		passes = append(passes, Pass{ID: "surface-texture", Role: "surface texture", Mode: "image", Asset: surface.Path, MediaType: surface.MediaType, ViewBox: surface.ViewBox, Repeat: "repeat", ScaleMilli: int(fixed(theme.Texture.Scale, .25, 3) * 1000), Scope: theme.Texture.Scope})
	}
	if theme.Motif.Layer == "over" || theme.Motif.Layer == recipeLayerSandwich {
		passes = append(passes, Pass{ID: "motif-over", Role: "motif-over", Mode: "image", Asset: motif.Path, MediaType: motif.MediaType, ViewBox: motif.ViewBox, Repeat: "repeat", ScaleMilli: 1000})
	}
	if theme.HeadingTexture.Kind != "none" {
		passes = append(passes,
			Pass{ID: "heading-wear", Role: "heading wear", Mode: "alpha-mask", Asset: heading.Path, MediaType: heading.MediaType, ViewBox: heading.ViewBox, Repeat: "repeat", ScaleMilli: int(fixed(theme.HeadingTexture.Scale, .25, 3) * 1000), Paint: mixColor(colors.background, colors.ink, theme.HeadingTexture.ColorMix)},
		)
	}
	for i := range passes {
		if strings.HasPrefix(passes[i].ID, "motif-") {
			size, gap, err := motifDimensions(theme.Motif)
			if err != nil {
				return Bundle{}, err
			}
			passes[i].SizeMilli, passes[i].GapMilli = int(size*1000), int(gap*1000)
		}
	}
	manifest := Manifest{SchemaVersion: SchemaVersion, ContractVersion: ContractVersion, SemanticHash: semantic, Assets: []Asset{surface, heading, motif}, Passes: passes, Fonts: fonts}
	identity := manifest
	identity.RecipeHash = ""
	encoded, err := canonicalJSON(identity)
	if err != nil {
		return Bundle{}, err
	}
	manifest.RecipeHash = digest([]byte(recipeDomain), encoded)
	if err := validateManifest(manifest); err != nil {
		return Bundle{}, err
	}
	return Bundle{Manifest: manifest, Assets: assets}, nil
}

//nolint:gocyclo // Manifest validation checks independent schema invariants.
func validateManifest(manifest Manifest) error {
	if manifest.SchemaVersion != SchemaVersion || manifest.ContractVersion != ContractVersion || len(manifest.SemanticHash) != 64 || len(manifest.RecipeHash) != 64 {
		return errors.New("invalid rendering recipe manifest identity")
	}
	if len(manifest.Assets) != 3 || len(manifest.Passes) == 0 || len(manifest.Fonts) != 3 {
		return errors.New("invalid rendering recipe manifest cardinality")
	}
	for i := range manifest.Assets {
		asset := &manifest.Assets[i]
		if asset.MediaType != recipeSVGMediaType || !isHexDigest(asset.SHA256) || asset.ViewBox == "" || asset.Width <= 0 || asset.Height <= 0 {
			return fmt.Errorf("invalid manifest asset %q", asset.Path)
		}
	}
	for i := range manifest.Passes {
		pass := &manifest.Passes[i]
		if pass.ID == "" || pass.Asset == "" || pass.MediaType != recipeSVGMediaType || pass.ViewBox == "" || pass.ScaleMilli <= 0 {
			return fmt.Errorf("invalid manifest pass %q", pass.ID)
		}
		if pass.Mode != "image" && pass.Mode != "alpha-mask" {
			return fmt.Errorf("invalid manifest pass mode %q", pass.Mode)
		}
		if pass.Paint != "" && !isHexColor(pass.Paint) {
			return fmt.Errorf("invalid manifest paint %q", pass.Paint)
		}
	}
	return nil
}
func isHexDigest(value string) bool {
	if len(value) != 64 {
		return false
	}
	for _, r := range value {
		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f')) {
			return false
		}
	}
	return true
}
func isHexColor(value string) bool {
	if len(value) != 7 || value[0] != '#' {
		return false
	}
	for _, r := range value[1:] {
		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f')) {
			return false
		}
	}
	return true
}

func canonicalFonts() []FontAsset {
	return []FontAsset{
		{Role: "body", Family: "Space Grotesk", AssetPath: "internal/fontcatalog/space-grotesk/space-grotesk-prose-core.woff2", SHA256: "5010c6e79c7e65a058dfd39a71d8d60f7b56070a3c32cc911d9dddbaec7f0851", Weight: 400, Style: "normal", FallbackStack: []string{"system-ui", "sans-serif"}},
		{Role: "heading", Family: "Knewave", AssetPath: "internal/fontcatalog/knewave/knewave-display-core.woff2", SHA256: "56fc586096f8e158fc214e3de31ee47c15ff80fc0abd5e0646ddff47a6af27ba", Weight: 400, Style: "normal", FallbackStack: []string{"cursive", "sans-serif"}},
		{Role: "mono", Family: "DM Mono", AssetPath: "internal/fontcatalog/dm-mono/dm-mono-code-core.woff2", SHA256: "90c42888e73fa88ebf3d0e27e0461535acd661e3892bbcb6eb21047f33b67c32", Weight: 400, Style: "normal", FallbackStack: []string{"ui-monospace", "monospace"}},
	}
}

func digest(domain, data []byte) string {
	sum := sha256.Sum256(append(append([]byte{}, domain...), data...))
	return hex.EncodeToString(sum[:])
}

func canonicalJSON(value any) ([]byte, error) {
	_ = writeCanonical
	b, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	return jsoncanonicalizer.Transform(b)
}

//nolint:gocyclo // Compatibility encoder mirrors JSON's type dispatch.
func writeCanonical(out *bytes.Buffer, value any) error { // retained for compatibility with package tests
	switch value := value.(type) {
	case nil:
		out.WriteString("null")
	case string:
		if !utf8.ValidString(value) {
			return errors.New("invalid UTF-8 string")
		}
		b, err := json.Marshal(value)
		if err != nil {
			return err
		}
		out.Write(b)
	case bool:
		if value {
			out.WriteString("true")
		} else {
			out.WriteString("false")
		}
	case int:
		out.WriteString(strconv.Itoa(value))
	case int64:
		out.WriteString(strconv.FormatInt(value, 10))
	case float64:
		if !math.IsNaN(value) && !math.IsInf(value, 0) {
			out.WriteString(strconv.FormatFloat(value, 'f', -1, 64))
		} else {
			return errors.New("non-finite number")
		}
	case []string:
		out.WriteByte('[')
		for i, item := range value {
			if i > 0 {
				out.WriteByte(',')
			}
			if err := writeCanonical(out, item); err != nil {
				return err
			}
		}
		out.WriteByte(']')
	case []Asset:
		return writeStruct(out, value)
	case []Pass:
		return writeStruct(out, value)
	case []FontAsset:
		return writeStruct(out, value)
	case []any:
		out.WriteByte('[')
		for i, item := range value {
			if i > 0 {
				out.WriteByte(',')
			}
			if err := writeCanonical(out, item); err != nil {
				return err
			}
		}
		out.WriteByte(']')
	case map[string]any:
		keys := make([]string, 0, len(value))
		for key := range value {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		out.WriteByte('{')
		for i, key := range keys {
			if i > 0 {
				out.WriteByte(',')
			}
			b, err := json.Marshal(key)
			if err != nil {
				return err
			}
			out.Write(b)
			out.WriteByte(':')
			if err := writeCanonical(out, value[key]); err != nil {
				return err
			}
		}
		out.WriteByte('}')
	default:
		b, err := json.Marshal(value)
		if err != nil {
			return err
		}
		var generic any
		if err := json.Unmarshal(b, &generic); err != nil {
			return err
		}
		return writeCanonical(out, generic)
	}
	return nil
}

func writeStruct(out *bytes.Buffer, value any) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	var generic any
	if err := decodeStrict(b, &generic); err != nil {
		return err
	}
	return writeCanonical(out, generic)
}

func decodeStrict(data []byte, target any) error {
	if !utf8.Valid(data) {
		return errors.New("input is not valid UTF-8")
	}
	if hasLoneSurrogateEscape(data) {
		return errors.New("lone surrogate in JSON string")
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	var value any
	if err := decodeValue(dec, &value); err != nil {
		return err
	}
	if err := rejectSurrogates(value); err != nil {
		return err
	}
	var trailing any
	if err := dec.Decode(&trailing); err != io.EOF {
		if err == nil {
			return errors.New("trailing JSON")
		}
		return fmt.Errorf("trailing JSON: %w", err)
	}
	encoded, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return json.Unmarshal(encoded, target)
}
func hasLoneSurrogateEscape(data []byte) bool {
	for i := 0; i+5 < len(data); i++ {
		if data[i] != '\\' || data[i+1] != 'u' {
			continue
		}
		n, err := strconv.ParseUint(string(data[i+2:i+6]), 16, 16)
		if err == nil && n >= 0xD800 && n <= 0xDFFF {
			return true
		}
	}
	return false
}
func rejectSurrogates(value any) error {
	switch v := value.(type) {
	case string:
		for _, r := range v {
			if r >= 0xD800 && r <= 0xDFFF {
				return errors.New("lone surrogate in JSON string")
			}
		}
	case map[string]any:
		for k, item := range v {
			if err := rejectSurrogates(k); err != nil {
				return err
			}
			if err := rejectSurrogates(item); err != nil {
				return err
			}
		}
	case []any:
		for _, item := range v {
			if err := rejectSurrogates(item); err != nil {
				return err
			}
		}
	}
	return nil
}
func decodeMap(value map[string]any, target any) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, target)
}
func decodeValue(dec *json.Decoder, out *any) error {
	token, err := dec.Token()
	if err != nil {
		return err
	}
	switch token := token.(type) {
	case json.Delim:
		if token == '{' {
			m := map[string]any{}
			for dec.More() {
				key, err := dec.Token()
				if err != nil {
					return err
				}
				ks, ok := key.(string)
				if !ok {
					return errors.New("invalid JSON object key")
				}
				if _, exists := m[ks]; exists {
					return fmt.Errorf("duplicate JSON key %q", ks)
				}
				var v any
				if err := decodeValue(dec, &v); err != nil {
					return err
				}
				m[ks] = v
			}
			if _, err := dec.Token(); err != nil {
				return err
			}
			*out = m
			return nil
		}
		if token == '[' {
			a := []any{}
			for dec.More() {
				var v any
				if err := decodeValue(dec, &v); err != nil {
					return err
				}
				a = append(a, v)
			}
			if _, err := dec.Token(); err != nil {
				return err
			}
			*out = a
			return nil
		}
		return errors.New("invalid JSON delimiter")
	case json.Number:
		if _, err := strconv.ParseFloat(string(token), 64); err != nil {
			return errors.New("invalid or non-finite number")
		}
		*out = token
		return nil
	default:
		*out = token
		return nil
	}
}

type paletteColorSet struct{ background, ink, accent, muted, shadow string }

func paletteColors(id string) paletteColorSet {
	contract, err := renderingcontract.Load()
	if err != nil {
		return paletteColorSet{"#fff", "#000", "#000", "#000", "#fff"}
	}
	for _, palette := range contract.Palettes {
		if palette.ID == id {
			final := renderingcontract.FinalRenderPalette(palette)
			return paletteColorSet{final["background"], final["text"], final["link"], final["text"], final["background"]}
		}
	}
	return paletteColorSet{"#fff", "#000", "#000", "#000", "#fff"}
}
func validateTheme(theme Theme) error {
	if theme.ContractVersion != ContractVersion {
		return fmt.Errorf("unsupported contract_version %d", theme.ContractVersion)
	}
	contract, err := renderingcontract.Load()
	if err != nil {
		return err
	}
	if !enum(contract.Palettes, theme.Palette) {
		return fmt.Errorf("unknown palette %q", theme.Palette)
	}
	if !contains(contract.Enums["aesthetics"], theme.Aesthetic) {
		return fmt.Errorf("unknown aesthetic %q", theme.Aesthetic)
	}
	if theme.Fontpack != "brush" {
		return fmt.Errorf("fontpack %q is not compiled in recipe v1", theme.Fontpack)
	}
	if theme.Texture.Kind != "none" && theme.Texture.Kind != "screenprint" {
		return fmt.Errorf("recipe v1 does not support texture kind %q", theme.Texture.Kind)
	}
	if theme.HeadingTexture.Kind != recipeNone && theme.HeadingTexture.Kind != "splatter" && theme.HeadingTexture.Kind != "inherit" {
		return fmt.Errorf("recipe v1 does not support heading texture kind %q", theme.HeadingTexture.Kind)
	}
	if theme.Motif.Kind != "block-w" {
		return fmt.Errorf("recipe v1 supports block-w motif representative primitive only")
	}
	if theme.Motif.Layer != "under" && theme.Motif.Layer != "over" && theme.Motif.Layer != "sandwich" {
		return fmt.Errorf("invalid motif layer %q", theme.Motif.Layer)
	}
	if !contains(contract.Enums["motif_colors"], theme.Motif.Color) {
		return fmt.Errorf("invalid motif color %q", theme.Motif.Color)
	}
	if theme.Texture.Scope != "all" && theme.Texture.Scope != "quiet" {
		return fmt.Errorf("invalid texture scope %q", theme.Texture.Scope)
	}
	if _, _, err := motifDimensions(theme.Motif); err != nil {
		return err
	}
	return nil
}
func surfaceSVG(mix float64, background, ink string, source recipeSource) []byte {
	color := mixColor(background, ink, mix)
	offset := float64(source.Seed % 17)
	// The texture is a transparent decoration. Its owning surface supplies the
	// background; baking that color into the tile makes quiet scope opaque and
	// leaks a second, incorrect surface through transparent chrome.
	return []byte(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox=%q><circle cx="%.3f" cy="%.3f" r="11" fill=%q/><circle cx="%.3f" cy="%.3f" r="8" fill=%q/></svg>`, source.ViewBox, 32+offset, 48+offset, color, 121-offset, 137-offset, color))
}
func headingSVG(_ float64, source recipeSource) []byte {
	offset := float64(source.Seed % 13)
	// This is a direct alpha image. The even-odd path cuts genuine transparent
	// holes from the opaque field; the second hole is restored only partially.
	return []byte(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox=%q><path fill="white" fill-rule="evenodd" d="M0 0H180V180H0ZM%.3f 47m-3 0a3 3 0 1 0 6 0a3 3 0 1 0-6 0ZM%.3f 113m-4 0a4 4 0 1 0 8 0a4 4 0 1 0-8 0Z"/><circle cx="%.3f" cy="113" r="4" fill="white" opacity=".42"/></svg>`, source.ViewBox, 31+offset, 129-offset, 129-offset))
}
func motifSVG(m renderingcontract.MotifState, color, background string, source recipeSource) []byte {
	path := source.Path
	if path == "" {
		return nil
	}
	var svg strings.Builder
	width, height, err := motifFieldDimensions(m, source)
	if err != nil {
		return nil
	}
	svg.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" data-field="%dx%d">`, width, height, source.Columns, source.Rows))
	markColor := mixColor(background, colorForMotif(m.Color, color), m.ColorMix)
	markSize, gap, err := motifDimensions(m)
	if err != nil {
		return nil
	}
	// The path's normalized bounds are 1688.4 by 905.76. Size is the visible
	// mark width; gap is the space between those normalized bounds.
	markScale := markSize / 1688.4
	for index := 0; index < source.Columns*source.Rows; index++ {
		row, column := index/source.Columns, index%source.Columns
		x := motifPadding + float64(column)*(markSize+gap) + float64(row%2)*m.RowOffset*(markSize+gap)
		y := motifPadding + float64(row)*(markSize*905.76/1688.4+gap)
		seedIndex := index + source.Seed
		rotation := (float64(seedIndex%7) - 3) * m.Wobble * 2
		scale := 1 + (float64((seedIndex*13)%11)-5)*m.Scatter*.01
		fmt.Fprintf(&svg, `<path data-index="%d" fill="%s" transform="translate(%.3f %.3f) rotate(%.3f 844.2 452.9) scale(%.5f)" d="%s"/>`, index, markColor, x, y, rotation, scale*markScale, path)
	}
	svg.WriteString(`</svg>`)
	return []byte(svg.String())
}

const motifPadding = 8

func motifFieldDimensions(m renderingcontract.MotifState, source recipeSource) (width, height int, err error) {
	size, gap, err := motifDimensions(m)
	if err != nil {
		return 0, 0, err
	}
	markHeight := size * 905.76 / 1688.4
	return int(math.Ceil(float64(source.Columns)*(size+gap) + 2*motifPadding)), int(math.Ceil(float64(source.Rows)*(markHeight+gap) + 2*motifPadding)), nil
}

var pixelsRE = regexp.MustCompile(`^(40|[4-9]\d|1[0-3]\d|140)px$`)
var gapRE = regexp.MustCompile(`^(0|[1-9]|[12]\d|32)px$`)

func parseDimension(value string, pattern *regexp.Regexp, minimum, maximum float64) (float64, error) {
	if !pattern.MatchString(value) {
		return 0, fmt.Errorf("invalid dimension %q", value)
	}
	n, err := strconv.ParseFloat(strings.TrimSuffix(value, "px"), 64)
	if err != nil || n < minimum || n > maximum {
		return 0, fmt.Errorf("dimension out of bounds %q", value)
	}
	return n, nil
}
func motifDimensions(m renderingcontract.MotifState) (size, gap float64, err error) {
	size, err = parseDimension(m.Size, pixelsRE, 40, 140)
	if err != nil {
		return 0, 0, fmt.Errorf("motif.size: %w", err)
	}
	gap, err = parseDimension(m.Gap, gapRE, 0, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("motif.gap: %w", err)
	}
	return size, gap, nil
}
func colorForMotif(_, ink string) string { return ink }
func motifRoleColor(role string, colors paletteColorSet) string {
	switch role {
	case "accent":
		return colors.accent
	case "muted":
		return colors.muted
	case "shadow":
		return colors.shadow
	default:
		return colors.ink
	}
}
func mixColor(background, foreground string, mix float64) string {
	bg, err1 := parseColor(background)
	fg, err2 := parseColor(foreground)
	if err1 != nil || err2 != nil {
		return foreground
	}
	for i := range bg {
		bg[i] = bg[i]*(1-mix) + fg[i]*mix
	}
	return fmt.Sprintf("#%02x%02x%02x", uint8(math.Round(bg[0]*255)), uint8(math.Round(bg[1]*255)), uint8(math.Round(bg[2]*255)))
}
func parseColor(value string) ([3]float64, error) {
	c, err := strconv.ParseUint(strings.TrimPrefix(value, "#"), 16, 24)
	if err != nil {
		return [3]float64{}, err
	}
	return [3]float64{float64(c>>16) / 255, float64((c>>8)&255) / 255, float64(c&255) / 255}, nil
}
