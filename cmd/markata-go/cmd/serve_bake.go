package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"sync"

	"github.com/WaylonWalker/markata/pkg/config"
	"github.com/WaylonWalker/markata/pkg/models"
)

// serveThemeBakeEndpoint is the dev-server endpoint the theme picker uses to
// write the visitor's current choices into the site's config.
const serveThemeBakeEndpoint = "/__markata/theme/bake"

const themeBakeGroup = "theme"

var (
	// themeBakeNamePattern allows palette, aesthetic, and fontpack identifiers
	// such as "st.-patrick's-day" while rejecting control characters.
	themeBakeNamePattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9 ._'&+()-]{0,79}$`)

	// serveRequestFullRebuild queues a full rebuild; set by runServe.
	serveRequestFullRebuild func()

	// serveConfigMu serializes config edits, previews, and the config loads of
	// serve rebuilds, so a rebuild never reads a half-applied bake or one that
	// is about to be rolled back.
	serveConfigMu sync.Mutex
)

type themeBakeRequest struct {
	Palette      string `json:"palette"`
	PaletteLight string `json:"palette_light"`
	PaletteDark  string `json:"palette_dark"`
	Seasonal     bool   `json:"seasonal"`
	FallbackMode string `json:"fallback_mode"`
	Aesthetic    string `json:"aesthetic"`
	Fontpack     string `json:"fontpack"`
	TextSize     string `json:"text_size"`
}

type themeBakeResponse struct {
	Target string `json:"target"`
	Path   string `json:"path"`
	Exists bool   `json:"exists"`
	// TargetKind flags an unusual target (override or global).
	TargetKind string   `json:"target_kind,omitempty"`
	Keys       []string `json:"keys,omitempty"`
	Warnings   []string `json:"warnings,omitempty"`
	Error      string   `json:"error,omitempty"`
}

// resolveServeConfigSources loads the effective config and returns every
// config file that contributes to it as absolute paths, ordered from lowest
// to highest precedence: root config, its includes, then --merge-config files.
func resolveServeConfigSources(cfgPath string) (*models.Config, []string, error) {
	cfg, _, configPaths, err := loadManagerConfigWith(cfgPath, configPreview{})
	if err != nil {
		return nil, nil, err
	}
	absPaths := make([]string, 0, len(configPaths))
	for _, p := range configPaths {
		abs, err := filepath.Abs(p)
		if err != nil {
			return nil, nil, fmt.Errorf("resolve config path %s: %w", p, err)
		}
		absPaths = append(absPaths, abs)
	}
	if len(absPaths) == 0 {
		return cfg, nil, nil
	}
	sources, err := discoverConfigSourcePaths(absPaths)
	if err != nil {
		return nil, nil, err
	}
	return cfg, sources, nil
}

func validThemeBakeName(field, value string) error {
	if value == "" || themeBakeNamePattern.MatchString(value) {
		return nil
	}
	return fmt.Errorf("invalid %s %q", field, value)
}

// themeBakeValues converts a picker request into theme config keys, in the
// order they are written. seasonalSet reports whether the site currently
// turns seasonal palettes on, which a palette pick must switch off to show.
func themeBakeValues(req themeBakeRequest, cfg *models.Config, seasonalSet bool) ([]config.BakeValue, error) {
	for field, value := range map[string]string{
		"palette":       req.Palette,
		"palette_light": req.PaletteLight,
		"palette_dark":  req.PaletteDark,
		"aesthetic":     req.Aesthetic,
		"fontpack":      req.Fontpack,
	} {
		if err := validThemeBakeName(field, value); err != nil {
			return nil, err
		}
	}
	fontpacksFile := ""
	if cfg != nil {
		fontpacksFile = cfg.FontpacksFile
	}
	if !config.KnownFontpack(req.Fontpack, fontpacksFile) {
		return nil, fmt.Errorf("unknown fontpack %q", req.Fontpack)
	}
	switch req.FallbackMode {
	case "", "light", "dark":
	default:
		return nil, fmt.Errorf("invalid fallback_mode %q", req.FallbackMode)
	}
	switch req.TextSize {
	case "", models.TextSizeSmall, models.TextSizeMedium, models.TextSizeLarge, models.TextSizeXLarge:
	default:
		return nil, fmt.Errorf("invalid text_size %q", req.TextSize)
	}

	var values []config.BakeValue
	add := func(key, value string) {
		if value != "" {
			values = append(values, config.BakeValue{Key: key, Value: value})
		}
	}
	if req.Seasonal {
		// Seasonal picks change by date; the configured palettes remain the fallback.
		values = append(values, config.BakeValue{Key: "seasonal", Value: true})
	} else {
		add("palette", req.Palette)
		add("palette_light", req.PaletteLight)
		add("palette_dark", req.PaletteDark)
		if len(values) > 0 && (seasonalSet || (cfg != nil && cfg.Theme.Seasonal)) {
			values = append(values, config.BakeValue{Key: "seasonal", Value: false})
		}
	}
	add("fallback_mode", req.FallbackMode)
	add("aesthetic", req.Aesthetic)
	add("fontpack", req.Fontpack)
	add("text_size", req.TextSize)
	if len(values) == 0 {
		return nil, errors.New("no theme settings to bake")
	}
	return values, nil
}

// themeBakeEnvWarnings reports MARKATA_GO_THEME_* variables that would
// override baked values.
func themeBakeEnvWarnings() []string {
	var warnings []string
	for _, entry := range os.Environ() {
		name, _, _ := strings.Cut(entry, "=")
		if strings.HasPrefix(name, "MARKATA_GO_THEME_") {
			warnings = append(warnings, name+" is set and overrides the config file")
		}
	}
	return warnings
}

func displayConfigPath(path string) string {
	if cwd, err := os.Getwd(); err == nil {
		if rel, err := filepath.Rel(cwd, path); err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return rel
		}
	}
	return path
}

// sameOriginRequest rejects writes that do not come from the developer's own
// machine. The client must connect over loopback, so serving on 0.0.0.0 or a
// LAN address never lets other hosts edit the config. The Host must be
// localhost or an IP literal so a DNS-rebinding page served from an attacker's
// domain cannot pass the same-origin check, and a page on another origin must
// not be able to edit the config through the local dev server.
func sameOriginRequest(r *http.Request) bool {
	if !loopbackRemote(r.RemoteAddr) || !localDevHost(r.Host) {
		return false
	}
	origin := r.Header.Get("Origin")
	if origin == "" {
		return r.Header.Get("Sec-Fetch-Site") == "" || r.Header.Get("Sec-Fetch-Site") == "same-origin"
	}
	parsed, err := url.Parse(origin)
	if err != nil {
		return false
	}
	return parsed.Host == r.Host
}

func writeThemeBakeJSON(w http.ResponseWriter, status int, resp themeBakeResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	//nolint:errcheck // Best effort write to HTTP response
	json.NewEncoder(w).Encode(resp)
}

// handleThemeBake serves GET (describe the target file) and POST (write the
// picker's choices). A POST is a settings bake of theme.* keys, so it shares
// the sidebar's target selection, verification, rollback, and preview
// handling: baked keys leave the unsaved preview so they take effect.
func handleThemeBake(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		w.Header().Set("Allow", "GET, POST")
		writeThemeBakeJSON(w, http.StatusMethodNotAllowed, themeBakeResponse{Error: "method not allowed"})
		return
	}
	if !sameOriginRequest(r) {
		writeThemeBakeJSON(w, http.StatusForbidden, themeBakeResponse{Error: "cross-origin requests are not allowed"})
		return
	}

	serveConfigMu.Lock()
	defer serveConfigMu.Unlock()

	state, err := loadSettingsState()
	if err != nil {
		writeThemeBakeJSON(w, http.StatusInternalServerError, themeBakeResponse{Error: err.Error()})
		return
	}
	target, err := state.target([]string{themeBakeGroup, "palette"})
	if err != nil {
		writeThemeBakeJSON(w, http.StatusInternalServerError, themeBakeResponse{Error: err.Error()})
		return
	}
	_, statErr := os.Stat(target)
	resp := themeBakeResponse{
		Target:     displayConfigPath(target),
		Path:       target,
		Exists:     statErr == nil,
		TargetKind: state.kinds[target],
		Warnings:   themeBakeEnvWarnings(),
	}
	if r.Method == http.MethodGet {
		writeThemeBakeJSON(w, http.StatusOK, resp)
		return
	}

	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		resp.Error = "Content-Type must be application/json"
		writeThemeBakeJSON(w, http.StatusUnsupportedMediaType, resp)
		return
	}
	var req themeBakeRequest
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		resp.Error = "invalid request: " + err.Error()
		writeThemeBakeJSON(w, http.StatusBadRequest, resp)
		return
	}
	seasonalSet := state.source([]string{themeBakeGroup, "seasonal"}) != ""
	values, err := themeBakeValues(req, state.cfg, seasonalSet)
	if err != nil {
		resp.Error = err.Error()
		writeThemeBakeJSON(w, http.StatusBadRequest, resp)
		return
	}
	changes := make([]settingsChange, len(values))
	for i, v := range values {
		changes[i] = settingsChange{Key: themeBakeGroup + "." + v.Key, Value: v.Value}
	}
	baked, status := applySettingsChanges(state, changes, false)
	resp.Warnings = append(resp.Warnings, baked.Warnings...)
	if baked.Error != "" {
		resp.Error = baked.Error
		writeThemeBakeJSON(w, status, resp)
		return
	}
	var targets []string
	for _, c := range baked.Changed {
		resp.Keys = append(resp.Keys, strings.TrimPrefix(c.Key, themeBakeGroup+"."))
		if !slices.Contains(targets, c.Target) {
			targets = append(targets, c.Target)
		}
	}
	resp.Target = strings.Join(targets, ", ")
	if len(targets) == 1 {
		if abs, err := filepath.Abs(targets[0]); err == nil {
			resp.Path = abs
		}
	}
	resp.Exists = true
	writeThemeBakeJSON(w, status, resp)
}

// loopbackRemote reports whether the client address is a loopback address.
func loopbackRemote(remoteAddr string) bool {
	host := remoteAddr
	if h, _, err := net.SplitHostPort(remoteAddr); err == nil {
		host = h
	}
	ip := net.ParseIP(strings.Trim(host, "[]"))
	return ip != nil && ip.IsLoopback()
}

func localDevHost(hostport string) bool {
	host := hostport
	if h, _, err := net.SplitHostPort(hostport); err == nil {
		host = h
	}
	host = strings.Trim(host, "[]")
	if host == "localhost" || strings.HasSuffix(host, ".localhost") {
		return true
	}
	return net.ParseIP(host) != nil
}
