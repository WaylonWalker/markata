package config

//go:generate go run ../../internal/configdocs/gen

import (
	"encoding"
	"errors"
	"fmt"
	"math"
	"path/filepath"
	"reflect"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/WaylonWalker/markata/pkg/fontpacks"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/palettes"
	"github.com/WaylonWalker/markata/pkg/renderingcontract"

	"github.com/alecthomas/chroma/v2/styles"
)

// Setting kinds describe how a config value is edited.
const (
	SettingString  = "string"
	SettingBool    = "bool"
	SettingInt     = "int"
	SettingFloat   = "float"
	SettingList    = "list"
	SettingComplex = "complex"
)

// maxSettingStringLen bounds a single edited string value.
const maxSettingStringLen = 4096

// SettingField describes one config key below the markata-go wrapper.
type SettingField struct {
	// Key is the dotted key path, e.g. "theme.switcher.enabled".
	Key string `json:"key"`
	// Section is the first key segment for nested settings, or "" for
	// top-level site settings.
	Section string `json:"section"`
	Label   string `json:"label"`
	Kind    string `json:"kind"`
	// Value is the effective value. It is nil for unset optional settings,
	// complex settings, and sensitive settings.
	Value any `json:"value"`
	// Summary describes complex values, e.g. "3 items".
	Summary string `json:"summary,omitempty"`
	Doc     string `json:"doc,omitempty"`
	// Options are the allowed values for a string setting. When Closed is
	// true any other non-empty value is rejected; otherwise they are
	// suggestions.
	Options []string `json:"options,omitempty"`
	Closed  bool     `json:"closed,omitempty"`
	// Min and Max bound numeric settings when set.
	Min       *float64 `json:"min,omitempty"`
	Max       *float64 `json:"max,omitempty"`
	Optional  bool     `json:"optional,omitempty"`
	Sensitive bool     `json:"sensitive,omitempty"`
	// Unsupported marks model fields that the config loader does not read
	// from config files, so editing them would have no effect.
	Unsupported bool `json:"unsupported,omitempty"`
	Editable    bool `json:"editable"`

	path    []string
	intBits int
}

// SensitiveSettingName reports whether a config key name (the last path
// segment) holds a secret that must not be shown or edited in the dev UI.
func SensitiveSettingName(name string) bool {
	return sensitiveSettingPattern.MatchString(name)
}

// ErrUnknownSetting is returned for keys that are not editable settings.
var ErrUnknownSetting = errors.New("unknown or read-only setting")

var (
	sensitiveSettingPattern = regexp.MustCompile(`(^|_)(secret|token|password|passphrase|api_key|private_key|credentials?)$`)
	validValuesPattern      = regexp.MustCompile(`(?i)(?:valid values(?: are)?:?|options:)([^.]*)`)
	quotedValuePattern      = regexp.MustCompile(`"([^"]+)"`)
	textUnmarshalerType     = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()

	settingTemplatesOnce sync.Once
	settingTemplates     []SettingField
)

// Settings describes every config setting with its effective value in cfg.
// Settings are returned in declaration order of models.Config.
func Settings(cfg *models.Config) []SettingField {
	templates := settingFieldTemplates()
	out := make([]SettingField, len(templates))
	root := reflect.ValueOf(cfg)
	memo := &settingOptionsMemo{}
	for i := range templates {
		field := templates[i]
		field.Options, field.Closed = memo.options(field, cfg)
		field.Min, field.Max = settingRange(field.Key)
		if v, ok := settingReflectValue(root, field.path); ok && !field.Sensitive {
			field.Value, field.Summary = settingValueOf(v, field.Kind)
		}
		out[i] = field
	}
	return out
}

// LookupSettingIn returns the editable setting for key with the options and
// ranges that apply to cfg, so CoerceSettingValue can enforce them.
func LookupSettingIn(cfg *models.Config, key string) (SettingField, bool) {
	field, ok := LookupSetting(key)
	if !ok {
		return field, false
	}
	field.Options, field.Closed = settingOptions(field, cfg)
	field.Min, field.Max = settingRange(field.Key)
	return field, true
}

// LookupSetting returns the editable setting for key.
func LookupSetting(key string) (SettingField, bool) {
	templates := settingFieldTemplates()
	for i := range templates {
		field := &templates[i]
		if field.Key == key {
			return *field, field.Editable
		}
	}
	return SettingField{}, false
}

// SettingPath returns the key path segments of a setting.
func (f SettingField) SettingPath() []string {
	return append([]string{}, f.path...)
}

// SettingValue returns the effective value of key in cfg.
func SettingValue(cfg *models.Config, key string) (any, bool) {
	templates := settingFieldTemplates()
	for i := range templates {
		field := &templates[i]
		if field.Key != key {
			continue
		}
		v, ok := settingReflectValue(reflect.ValueOf(cfg), field.path)
		if !ok {
			return nil, true
		}
		value, _ := settingValueOf(v, field.Kind)
		return value, true
	}
	return nil, false
}

// CoerceSettingValue converts a JSON-decoded value to the Go type stored for
// field, rejecting values of the wrong shape.
//
//nolint:gocyclo // Each setting kind has distinct coercion and validation rules.
func CoerceSettingValue(field SettingField, raw any) (any, error) {
	switch field.Kind {
	case SettingString:
		s, ok := raw.(string)
		if !ok {
			return nil, fmt.Errorf("%s must be a string", field.Key)
		}
		if err := checkSettingString(field.Key, s); err != nil {
			return nil, err
		}
		return s, checkSettingOption(field, s)
	case SettingBool:
		b, ok := raw.(bool)
		if !ok {
			return nil, fmt.Errorf("%s must be true or false", field.Key)
		}
		return b, nil
	case SettingInt:
		n, ok := raw.(float64)
		if !ok || n != math.Trunc(n) || math.IsInf(n, 0) {
			return nil, fmt.Errorf("%s must be a whole number", field.Key)
		}
		limit := math.Ldexp(1, field.intBits-1)
		if field.intBits >= 53 {
			limit = 1 << 53
		}
		if n >= limit || n < -limit {
			return nil, fmt.Errorf("%s is out of range", field.Key)
		}
		return int64(n), checkSettingRange(field, n)
	case SettingFloat:
		n, ok := raw.(float64)
		if !ok || math.IsNaN(n) || math.IsInf(n, 0) {
			return nil, fmt.Errorf("%s must be a number", field.Key)
		}
		return n, checkSettingRange(field, n)
	case SettingList:
		items, ok := raw.([]any)
		if !ok {
			return nil, fmt.Errorf("%s must be a list of strings", field.Key)
		}
		if len(items) > 256 {
			return nil, fmt.Errorf("%s has too many items", field.Key)
		}
		list := make([]string, 0, len(items))
		for _, item := range items {
			s, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("%s must be a list of strings", field.Key)
			}
			if err := checkSettingString(field.Key, s); err != nil {
				return nil, err
			}
			list = append(list, s)
		}
		return list, nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknownSetting, field.Key)
	}
}

// checkSettingOption rejects values outside a closed option list. The empty
// string is allowed and means "use the default".
func checkSettingOption(field SettingField, s string) error {
	if !field.Closed || s == "" || slices.Contains(field.Options, s) {
		return nil
	}
	if accept := settingAccepts[field.Key]; accept != nil && accept(s) {
		return nil
	}
	shown := field.Options
	if len(shown) > 12 {
		return fmt.Errorf("%s: unsupported value %q", field.Key, s)
	}
	quoted := make([]string, len(shown))
	for i, o := range shown {
		quoted[i] = strconv.Quote(o)
	}
	return fmt.Errorf("%s must be one of: %s", field.Key, strings.Join(quoted, ", "))
}

func checkSettingRange(field SettingField, n float64) error {
	if field.Min != nil && n < *field.Min {
		return fmt.Errorf("%s must be at least %v", field.Key, *field.Min)
	}
	if field.Max != nil && n > *field.Max {
		return fmt.Errorf("%s must be at most %v", field.Key, *field.Max)
	}
	return nil
}

func checkSettingString(key, s string) error {
	if len(s) > maxSettingStringLen {
		return fmt.Errorf("%s is too long", key)
	}
	if !utf8.ValidString(s) {
		return fmt.Errorf("%s is not valid UTF-8", key)
	}
	for _, r := range s {
		if r < 0x20 && r != '\n' && r != '\t' {
			return fmt.Errorf("%s contains control characters", key)
		}
	}
	return nil
}

// SettingsEqual reports whether two setting values are equal after
// normalizing numbers and lists.
func SettingsEqual(a, b any) bool {
	if la, ok := a.([]string); ok && len(la) == 0 {
		a = nil
	}
	if lb, ok := b.([]string); ok && len(lb) == 0 {
		b = nil
	}
	return reflect.DeepEqual(bakeComparable(a), bakeComparable(b))
}

func settingFieldTemplates() []SettingField {
	settingTemplatesOnce.Do(func() {
		settingTemplates = collectSettingFields(reflect.TypeOf(models.Config{}), nil, map[reflect.Type]bool{})
		for i := range settingTemplates {
			if settingTemplates[i].Editable && !settingParses(settingTemplates[i]) {
				settingTemplates[i].Editable = false
				settingTemplates[i].Unsupported = true
			}
		}
	})
	return settingTemplates
}

// settingParses reports whether the config loader reads field from a config
// file, by parsing a minimal document that sets it to a probe value.
func settingParses(field SettingField) bool {
	var probe any
	switch field.Kind {
	case SettingString:
		probe = "probe"
	case SettingBool:
		probe = true
	case SettingInt:
		probe = int64(7)
	case SettingFloat:
		probe = 0.375
	case SettingList:
		probe = []string{"probe"}
	default:
		return false
	}
	table, leaf := field.path[:len(field.path)-1], field.path[len(field.path)-1]
	data, err := bakeTOML(nil, table, []BakeValue{{Key: leaf, Value: probe}})
	if err != nil {
		return false
	}
	cfg, err := ParseTOML(data)
	if err != nil {
		return false
	}
	v, ok := settingReflectValue(reflect.ValueOf(cfg), field.path)
	if !ok {
		return false
	}
	got, _ := settingValueOf(v, field.Kind)
	return SettingsEqual(got, probe)
}

func settingTOMLName(f reflect.StructField) string {
	tag, ok := f.Tag.Lookup("toml")
	if !ok {
		return ""
	}
	name, _, _ := strings.Cut(tag, ",")
	if name == "-" {
		return ""
	}
	return name
}

func collectSettingFields(t reflect.Type, prefix []string, seen map[reflect.Type]bool) []SettingField {
	if seen[t] || len(prefix) > 6 {
		return nil
	}
	seen[t] = true
	defer delete(seen, t)

	var out []SettingField
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		name := settingTOMLName(f)
		if name == "" {
			continue
		}
		path := append(append([]string{}, prefix...), name)
		ft := f.Type
		optional := false
		if ft.Kind() == reflect.Pointer {
			ft = ft.Elem()
			optional = true
		}
		field := SettingField{
			Key:      strings.Join(path, "."),
			Label:    settingLabel(name),
			Doc:      settingDocs[t.Name()+"."+f.Name],
			Optional: optional,
			path:     path,
		}
		if len(path) > 1 {
			field.Section = path[0]
		}
		switch kind := settingKind(ft); kind {
		case "":
			out = append(out, collectSettingFields(ft, path, seen)...)
			continue
		default:
			field.Kind = kind
		}
		if field.Kind == SettingInt {
			field.intBits = ft.Bits()
		}
		field.Sensitive = sensitiveSettingPattern.MatchString(name)
		field.Editable = field.Kind != SettingComplex && !field.Sensitive
		out = append(out, field)
	}
	return out
}

// settingKind classifies t; "" means a nested settings struct.
func settingKind(t reflect.Type) string {
	if t.Kind() == reflect.Struct && (t.PkgPath() != "" && !strings.HasSuffix(t.PkgPath(), "/pkg/models") ||
		reflect.PointerTo(t).Implements(textUnmarshalerType)) {
		return SettingComplex
	}
	switch t.Kind() { //nolint:exhaustive // Remaining kinds are complex values.
	case reflect.String:
		return SettingString
	case reflect.Bool:
		return SettingBool
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		if t.PkgPath() == "time" {
			return SettingComplex
		}
		return SettingInt
	case reflect.Float32, reflect.Float64:
		return SettingFloat
	case reflect.Slice:
		if t.Elem().Kind() == reflect.String {
			return SettingList
		}
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			if t.Field(i).IsExported() && settingTOMLName(t.Field(i)) != "" {
				return ""
			}
		}
	}
	return SettingComplex
}

func settingLabel(name string) string {
	label := strings.ReplaceAll(name, "_", " ")
	if label == "" {
		return label
	}
	return strings.ToUpper(label[:1]) + label[1:]
}

// settingReflectValue walks root along the TOML key path.
func settingReflectValue(root reflect.Value, path []string) (reflect.Value, bool) {
	v := root
	for _, key := range path {
		for v.Kind() == reflect.Pointer {
			if v.IsNil() {
				return reflect.Value{}, false
			}
			v = v.Elem()
		}
		if v.Kind() != reflect.Struct {
			return reflect.Value{}, false
		}
		found := false
		for i := 0; i < v.NumField(); i++ {
			if v.Type().Field(i).IsExported() && settingTOMLName(v.Type().Field(i)) == key {
				v = v.Field(i)
				found = true
				break
			}
		}
		if !found {
			return reflect.Value{}, false
		}
	}
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return reflect.Value{}, false
		}
		v = v.Elem()
	}
	return v, true
}

func settingValueOf(v reflect.Value, kind string) (value any, summary string) {
	switch kind {
	case SettingString:
		return v.String(), ""
	case SettingBool:
		return v.Bool(), ""
	case SettingInt:
		if v.CanInt() {
			return v.Int(), ""
		}
		return int64(v.Uint()), "" //nolint:gosec // Only uint32 and smaller are editable.
	case SettingFloat:
		return v.Float(), ""
	case SettingList:
		list := make([]string, v.Len())
		for i := range list {
			list[i] = v.Index(i).String()
		}
		return list, ""
	}
	switch v.Kind() { //nolint:exhaustive // Other complex values have no summary.
	case reflect.Slice, reflect.Array:
		return nil, pluralCount(v.Len(), "item")
	case reflect.Map:
		return nil, pluralCount(v.Len(), "entry")
	}
	return nil, ""
}

func pluralCount(n int, noun string) string {
	if n == 1 {
		return "1 " + noun
	}
	if noun == "entry" {
		return fmt.Sprintf("%d entries", n)
	}
	return fmt.Sprintf("%d %ss", n, noun)
}

// settingEnums lists the allowed values of string settings whose doc comment
// does not carry a "Valid values:" line. Each list matches the templates or
// plugins that read the setting.
var settingEnums = map[string][]string{
	"theme.switcher.position":             {"header", "footer"},
	"components.nav.position":             {"header", "sidebar"},
	"components.nav.style":                {"horizontal", "vertical"},
	"components.doc_sidebar.position":     {"left", "right"},
	"components.feed_sidebar.position":    {"left", "right"},
	"components.content_sidebar.position": {"left", "right"},
	"search.backend":                      {"pagefind", "bleve"},
	"search.position":                     {"navbar", "sidebar", "footer", "custom"},
	"layout.docs.sidebar_position":        {"left", "right"},
	"layout.docs.toc_position":            {"left", "right"},
	"layout.blog.toc_position":            {"left", "right"},
	"layout.docs.header_style":            headerStyles,
	"layout.blog.header_style":            headerStyles,
	"layout.landing.header_style":         headerStyles,
	"layout.docs.footer_style":            footerStyles,
	"layout.blog.footer_style":            footerStyles,
	"layout.landing.footer_style":         footerStyles,
	"header.style":                        headerStyles,
	"sidebar.position":                    {"left", "right"},
	"sidebar.auto_generate.order_by":      {"title", "date", "nav_order", "filename"},
	"toc.position":                        {"left", "right"},
	"assets.mode":                         {"cdn", "self-hosted", "auto"},
	"markdown.render_profile":            {"web", "offline-interactive", "offline-static"},
	"feed_defaults.pagination_type":       paginationTypes,
	"blogroll.pagination_type":            paginationTypes,
}

var (
	headerStyles    = []string{"full", "minimal", "transparent", "none"}
	footerStyles    = []string{"full", "minimal", "none"}
	paginationTypes = []string{
		string(models.PaginationManual), string(models.PaginationHTMX),
		string(models.PaginationHTMXInfinite), string(models.PaginationJS),
	}
)

// settingAccepts admits values of closed settings that the option list
// cannot enumerate, such as palette aliases.
var settingAccepts = map[string]func(string) bool{
	"theme.palette":                  KnownPalette,
	"theme.palette_light":            KnownPalette,
	"theme.palette_dark":             KnownPalette,
	"theme_calendar.default_palette": KnownPalette,
}

func unitRange() (lo, hi *float64) {
	zero, one := 0.0, 1.0
	return &zero, &one
}

// settingRange returns the bounds enforced by ValidateConfig for numeric
// settings, so the sidebar can reject out-of-range values before a rebuild.
func settingRange(key string) (lo, hi *float64) {
	switch key {
	case "theme.texture.color_mix", "theme.heading_texture.color_mix", "theme.motif.color_mix",
		"theme.motif.row_offset", "theme.motif.wobble", "theme.motif.scatter":
		return unitRange()
	case "theme.texture.scale", "theme.heading_texture.scale":
		zero, most := 0.0, 3.0
		return &zero, &most
	}
	return nil, nil
}

// settingOptions returns the allowed values of a string setting and whether
// other values are rejected.
func settingOptions(field SettingField, cfg *models.Config) (options []string, closed bool) {
	return (&settingOptionsMemo{}).options(field, cfg)
}

// settingOptionsMemo shares palette discovery, the font catalog, and the
// rendering contract across the fields of one Settings call.
type settingOptionsMemo struct {
	palettes, fontpacks []string
	contract            *renderingcontract.Contract
	loaded              map[string]bool
}

func (m *settingOptionsMemo) once(name string) bool {
	if m.loaded == nil {
		m.loaded = map[string]bool{}
	}
	if m.loaded[name] {
		return false
	}
	m.loaded[name] = true
	return true
}

func (m *settingOptionsMemo) options(field SettingField, cfg *models.Config) (options []string, closed bool) {
	if field.Kind != SettingString {
		return nil, false
	}
	switch field.Key {
	case "theme.palette", "theme.palette_light", "theme.palette_dark", "theme_calendar.default_palette":
		if m.once("palettes") {
			m.palettes = paletteOptions()
		}
		return append([]string{}, m.palettes...), true
	case "fontpack", "theme.fontpack":
		if m.once("fontpacks") {
			m.fontpacks = fontpackOptions()
		}
		// A custom catalog can define any pack name.
		return append([]string{}, m.fontpacks...), cfg == nil || cfg.FontpacksFile == ""
	case "markdown.highlight.theme":
		return styles.Names(), true
	}
	if values, ok := settingEnums[field.Key]; ok {
		return append([]string{}, values...), true
	}
	contractGroups := map[string]string{
		"theme.aesthetic":            "aesthetics",
		"theme.texture.kind":         "textures",
		"theme.texture.scope":        "scopes",
		"theme.heading_texture.kind": "heading_textures",
		"theme.motif.kind":           "motifs",
		"theme.motif.layer":          "motif_layers",
		"theme.motif.color":          "motif_colors",
	}
	if group, ok := contractGroups[field.Key]; ok {
		if m.once("contract") {
			if c, err := renderingcontract.Load(); err == nil {
				m.contract = &c
			}
		}
		if m.contract != nil {
			return append([]string{}, m.contract.Enums[group]...), true
		}
	}
	return docOptions(field.Doc)
}

// docOptions parses a "Valid values: ..." or "Options: ..." doc line.
func docOptions(doc string) (options []string, closed bool) {
	m := validValuesPattern.FindStringSubmatch(doc)
	if m == nil {
		return nil, false
	}
	for _, q := range quotedValuePattern.FindAllStringSubmatch(m[1], -1) {
		if !slices.Contains(options, q[1]) {
			options = append(options, q[1])
		}
	}
	return options, len(options) > 1
}

func paletteOptions() []string {
	names := map[string]bool{}
	for _, id := range renderingcontract.PaletteIDs() {
		names[id] = true
	}
	if infos, err := palettes.NewLoader().Discover(); err == nil {
		for _, info := range infos {
			if info.Path != "" {
				names[strings.TrimSuffix(filepath.Base(info.Path), filepath.Ext(info.Path))] = true
			} else if info.Name != "" && !strings.Contains(info.Name, " ") {
				names[info.Name] = true
			}
		}
	}
	return sortedKeys(names)
}

func fontpackOptions() []string {
	source, err := fontpacks.BuiltinSource()
	if err != nil {
		return nil
	}
	names := map[string]bool{}
	for name := range source.Catalog.FontPacks {
		names[name] = true
	}
	for name := range source.Catalog.Aliases {
		names[name] = true
	}
	return sortedKeys(names)
}

func sortedKeys(set map[string]bool) []string {
	out := make([]string, 0, len(set))
	for key := range set {
		out = append(out, key)
	}
	sort.Strings(out)
	return out
}

// SettingsOverlay builds a raw config wrapper that sets each setting, for use
// as LoadOptions.Overlay. Lists become []any so they merge like parsed files.
func SettingsOverlay(settings []BakeSetting) map[string]any {
	if len(settings) == 0 {
		return nil
	}
	root := map[string]any{}
	for _, setting := range settings {
		if len(setting.Path) == 0 {
			continue
		}
		node := root
		for _, key := range setting.Path[:len(setting.Path)-1] {
			child, ok := node[key].(map[string]any)
			if !ok {
				child = map[string]any{}
				node[key] = child
			}
			node = child
		}
		value := setting.Value
		if list, ok := value.([]string); ok {
			items := make([]any, len(list))
			for i, item := range list {
				items[i] = item
			}
			value = items
		}
		node[setting.Path[len(setting.Path)-1]] = value
	}
	return map[string]any{"markata-go": root}
}
