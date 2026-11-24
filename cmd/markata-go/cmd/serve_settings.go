package cmd

import (
	"crypto/rand"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/WaylonWalker/markata/pkg/config"
	"github.com/WaylonWalker/markata/pkg/models"
)

// Settings sidebar endpoints. They exist only under `markata-go serve`;
// production builds never reference them.
const (
	serveSettingsEndpoint   = "/__markata/settings"
	serveSettingsScriptPath = "/__markata/settings.js"
	serveSettingsPreview    = "/__markata/settings/preview"
	maxSettingsChanges      = 200
)

//go:embed devui/settings.js
var serveSettingsScript []byte

// Target kinds flag config files a bake would not normally write to.
const (
	// targetKindOverride marks a --merge-config file: builds without that
	// flag do not see settings written there.
	targetKindOverride = "override"
	// targetKindGlobal marks the user-level config discovered from
	// ~/.config/markata-go; it applies to every site, so bakes refuse it.
	targetKindGlobal = "global"
)

// serveSessionID identifies this dev server process so the sidebar can tell
// a restart (which drops the in-memory preview) from a reset.
var serveSessionID = newServeSessionID()

func newServeSessionID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "session"
	}
	return hex.EncodeToString(b)
}

type settingsFieldJSON struct {
	config.SettingField
	// Default is the value the setting takes when no config file sets it.
	Default any `json:"default"`
	// Source is the config file that currently defines the setting.
	Source string `json:"source,omitempty"`
	// Sources lists every config file that defines the setting; resetting
	// it to the default removes it from all of them.
	Sources []string `json:"sources,omitempty"`
	// Target is the config file a change would be written to.
	Target string `json:"target"`
	// TargetKind flags an unusual target (override or global).
	TargetKind string `json:"target_kind,omitempty"`
	// Env names an environment variable that overrides the setting.
	Env string `json:"env,omitempty"`
}

type settingsResponse struct {
	Fields []settingsFieldJSON `json:"fields,omitempty"`
	// Preview lists unsaved settings the dev server is currently applying.
	Preview  []settingsChange  `json:"preview"`
	Files    []string          `json:"files,omitempty"`
	Changed  []settingsChanged `json:"changed,omitempty"`
	Diffs    []settingsDiff    `json:"diffs,omitempty"`
	Warnings []string          `json:"warnings,omitempty"`
	Error    string            `json:"error,omitempty"`
	// Session changes whenever the dev server restarts.
	Session string `json:"session,omitempty"`
	// SinglePage is true for `markata-go serve <file>`.
	SinglePage bool `json:"single_page,omitempty"`
}

type settingsChanged struct {
	Key    string `json:"key"`
	Target string `json:"target"`
	Unset  bool   `json:"unset,omitempty"`
}

// settingsDiff previews the edit a bake makes to one file.
type settingsDiff struct {
	Target  string `json:"target"`
	Kind    string `json:"kind,omitempty"`
	Created bool   `json:"created,omitempty"`
	Diff    string `json:"diff"`
}

type settingsChange struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
	// Unset resets the setting to its default by removing it from every
	// config file that defines it.
	Unset bool `json:"unset,omitempty"`
}

type settingsRequest struct {
	Changes []settingsChange `json:"changes"`
	// DryRun returns the diffs a bake would make without writing.
	DryRun bool `json:"dry_run,omitempty"`
}

// Preview state: unsaved settings merged over the config files for every
// rebuild while `markata-go serve` runs. They live only in memory.
var (
	servePreviewMu       sync.RWMutex
	servePreviewSettings []plannedSetting
)

// servePreviewConfig returns the current preview for config loading.
func servePreviewConfig() configPreview {
	servePreviewMu.RLock()
	defer servePreviewMu.RUnlock()
	return previewConfig(servePreviewSettings)
}

func previewConfig(planned []plannedSetting) configPreview {
	var preview configPreview
	settings := make([]config.BakeSetting, 0, len(planned))
	for _, p := range planned {
		if p.unset {
			preview.remove = append(preview.remove, p.path)
			continue
		}
		settings = append(settings, config.BakeSetting{Path: p.path, Value: p.value})
	}
	preview.overlay = config.SettingsOverlay(settings)
	return preview
}

func servePreviewChanges() []settingsChange {
	servePreviewMu.RLock()
	defer servePreviewMu.RUnlock()
	out := make([]settingsChange, len(servePreviewSettings))
	for i, p := range servePreviewSettings {
		out[i] = settingsChange{Key: p.key, Value: p.value, Unset: p.unset}
	}
	return out
}

func setServePreview(planned []plannedSetting) {
	sort.Slice(planned, func(i, j int) bool { return planned[i].key < planned[j].key })
	servePreviewMu.Lock()
	servePreviewSettings = planned
	servePreviewMu.Unlock()
}

// dropServePreviewKeys removes baked keys from the preview.
func dropServePreviewKeys(keys map[string]bool) {
	servePreviewMu.Lock()
	defer servePreviewMu.Unlock()
	kept := servePreviewSettings[:0:0]
	for _, p := range servePreviewSettings {
		if !keys[p.key] {
			kept = append(kept, p)
		}
	}
	servePreviewSettings = kept
}

// settingsState is the effective config plus the files that compose it.
type settingsState struct {
	cfg     *models.Config
	sources []string
	raws    []map[string]any
	kinds   map[string]string
}

func loadSettingsState() (*settingsState, error) {
	cfg, sources, err := resolveServeConfigSources(cfgFile)
	if err != nil {
		return nil, err
	}
	state := &settingsState{cfg: cfg, sources: sources, kinds: configSourceKinds(sources)}
	for _, path := range sources {
		raw, err := config.LoadRawConfigFile(path)
		if err != nil {
			return nil, err
		}
		state.raws = append(state.raws, raw)
	}
	return state, nil
}

// configSourceKinds flags --merge-config files and, when the root config was
// discovered in the user's config directory, that file and its includes.
func configSourceKinds(sources []string) map[string]string {
	kinds := map[string]string{}
	for _, path := range mergeConfigFiles {
		if abs, err := filepath.Abs(path); err == nil {
			kinds[abs] = targetKindOverride
		}
	}
	if len(sources) == 0 || cfgFile != "" {
		return kinds
	}
	home, err := os.UserHomeDir()
	if err != nil || sources[0] != filepath.Join(home, ".config", "markata-go", "config.toml") {
		return kinds
	}
	for _, path := range sources {
		if kinds[path] == "" {
			kinds[path] = targetKindGlobal
		}
	}
	return kinds
}

// source returns the highest-precedence file that defines path, or "".
func (s *settingsState) source(path []string) string {
	found := ""
	for _, src := range s.definedIn(path) {
		found = src
	}
	return found
}

// definedIn lists the files that define path, lowest precedence first.
func (s *settingsState) definedIn(path []string) []string {
	var out []string
	for i, raw := range s.raws {
		if config.SettingDefinedDepth(raw, path) == len(path) {
			out = append(out, s.sources[i])
		}
	}
	return out
}

// target picks the file that receives a change to path: the file defining
// the most specific part of the key path (later files win ties because they
// override earlier ones), else the root config. A config-less site gets a
// new markata-go.toml in the working directory.
func (s *settingsState) target(path []string) (string, error) {
	if len(s.sources) == 0 {
		return filepath.Abs("markata-go.toml")
	}
	best, bestDepth := s.sources[0], 0
	for i, raw := range s.raws {
		if depth := config.SettingDefinedDepth(raw, path); depth > 0 && depth >= bestDepth {
			best, bestDepth = s.sources[i], depth
		}
	}
	return best, nil
}

func settingEnvName(key string) string {
	return "MARKATA_GO_" + strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
}

func settingEnvOverride(key string) string {
	name := settingEnvName(key)
	if _, ok := os.LookupEnv(name); ok {
		return name
	}
	return ""
}

// settingDefaults returns the value of every setting when no config file
// sets it and no environment variable overrides it.
func settingDefaults() map[string]any {
	defaults := map[string]any{}
	cfg, err := config.LoadWithMergeOptions(config.LoadOptions{DisableDotEnv: true, DisableEnvOverrides: true}, "")
	if err != nil {
		return defaults
	}
	fields := config.Settings(cfg)
	for i := range fields {
		field := &fields[i]
		defaults[field.Key] = field.Value
	}
	return defaults
}

func (s *settingsState) describe() (settingsResponse, error) {
	resp := settingsResponse{Preview: servePreviewChanges(), Session: serveSessionID, SinglePage: serveSourceFile != ""}
	for _, path := range s.sources {
		resp.Files = append(resp.Files, displayConfigPath(path))
	}
	defaults := settingDefaults()
	fields := config.Settings(s.cfg)
	for i := range fields {
		field := &fields[i]
		path := field.SettingPath()
		target, err := s.target(path)
		if err != nil {
			return resp, err
		}
		entry := settingsFieldJSON{
			SettingField: *field,
			Target:       displayConfigPath(target),
			TargetKind:   s.kinds[target],
			Env:          settingEnvOverride(field.Key),
		}
		if !field.Sensitive {
			entry.Default = defaults[field.Key]
		}
		for _, src := range s.definedIn(path) {
			entry.Sources = append(entry.Sources, displayConfigPath(src))
			entry.Source = displayConfigPath(src)
		}
		resp.Fields = append(resp.Fields, entry)
	}
	return resp, nil
}

func writeSettingsJSON(w http.ResponseWriter, status int, resp settingsResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	//nolint:errcheck // Best effort write to HTTP response
	json.NewEncoder(w).Encode(resp)
}

func handleSettingsScript(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	//nolint:errcheck // Best effort write to HTTP response
	w.Write(serveSettingsScript)
}

// handleSettings serves GET (describe every setting with its value, source,
// and target file) and POST (bake changes into the config files, verify that
// the reloaded config reflects them, and queue a full rebuild).
func handleSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		w.Header().Set("Allow", "GET, POST")
		writeSettingsJSON(w, http.StatusMethodNotAllowed, settingsResponse{Error: "method not allowed"})
		return
	}
	if !sameOriginRequest(r) {
		writeSettingsJSON(w, http.StatusForbidden, settingsResponse{Error: "cross-origin requests are not allowed"})
		return
	}

	serveConfigMu.Lock()
	defer serveConfigMu.Unlock()

	state, err := loadSettingsState()
	if err != nil {
		writeSettingsJSON(w, http.StatusInternalServerError, settingsResponse{Error: err.Error()})
		return
	}
	if r.Method == http.MethodGet {
		resp, err := state.describe()
		if err != nil {
			writeSettingsJSON(w, http.StatusInternalServerError, settingsResponse{Error: err.Error()})
			return
		}
		writeSettingsJSON(w, http.StatusOK, resp)
		return
	}

	req, ok := decodeSettingsRequest(w, r)
	if !ok {
		return
	}
	resp, status := applySettingsChanges(state, req.Changes, req.DryRun)
	writeSettingsJSON(w, status, resp)
}

func decodeSettingsRequest(w http.ResponseWriter, r *http.Request) (settingsRequest, bool) {
	var req settingsRequest
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		writeSettingsJSON(w, http.StatusUnsupportedMediaType, settingsResponse{Error: "Content-Type must be application/json"})
		return req, false
	}
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 256<<10))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		writeSettingsJSON(w, http.StatusBadRequest, settingsResponse{Error: "invalid request: " + err.Error()})
		return req, false
	}
	return req, true
}

// handleSettingsPreview replaces the in-memory preview with the posted
// changes (an empty list resets it). The config files are not touched: the
// preview is checked like a bake (it must load, add no validation errors, and
// take effect), then a full rebuild applies it and live reload shows it.
func handleSettingsPreview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		writeSettingsJSON(w, http.StatusMethodNotAllowed, settingsResponse{Error: "method not allowed"})
		return
	}
	if !sameOriginRequest(r) {
		writeSettingsJSON(w, http.StatusForbidden, settingsResponse{Error: "cross-origin requests are not allowed"})
		return
	}
	req, ok := decodeSettingsRequest(w, r)
	if !ok {
		return
	}

	serveConfigMu.Lock()
	defer serveConfigMu.Unlock()

	state, err := loadSettingsState()
	if err != nil {
		writeSettingsJSON(w, http.StatusInternalServerError, settingsResponse{Error: err.Error()})
		return
	}
	if req.DryRun {
		writeSettingsJSON(w, http.StatusBadRequest, settingsResponse{Error: "dry_run applies only to bakes"})
		return
	}
	resp, status := previewSettingsChanges(state, req.Changes)
	writeSettingsJSON(w, status, resp)
}

func previewSettingsChanges(state *settingsState, changes []settingsChange) (resp settingsResponse, status int) {
	planned, err := planSettingsChanges(state, changes, true)
	if err != nil {
		resp.Error = err.Error()
		resp.Preview = servePreviewChanges()
		return resp, http.StatusBadRequest
	}
	if len(planned) > 0 {
		previewed, _, _, err := loadManagerConfigWith(cfgFile, previewConfig(planned))
		if err == nil {
			err = checkSettingsTookEffect(validationErrorSet(state.cfg), previewed, previewed, planned, &resp)
		}
		if err != nil {
			resp.Error = err.Error()
			resp.Preview = servePreviewChanges()
			return resp, http.StatusUnprocessableEntity
		}
	}
	setServePreview(planned)
	resp.Preview = servePreviewChanges()
	if len(planned) == 0 {
		infof("Settings preview reset")
	} else {
		infof("Previewing unsaved settings (%s)", strings.Join(plannedKeys(planned), ", "))
	}
	if serveRequestFullRebuild != nil {
		serveRequestFullRebuild()
	}
	return resp, http.StatusOK
}

func plannedKeys(planned []plannedSetting) []string {
	keys := make([]string, len(planned))
	for i, p := range planned {
		keys[i] = p.key
		if p.unset {
			keys[i] += " (reset)"
		}
	}
	return keys
}

// checkSettingsTookEffect fails when cfg has validation errors not in before,
// or a planned value did not load back (env overrides only warn). A reset
// setting must match its value in expected, the config loaded with the
// planned changes applied as a preview.
func checkSettingsTookEffect(before map[string]bool, cfg, expected *models.Config, planned []plannedSetting, resp *settingsResponse) error {
	for msg := range validationErrorSet(cfg) {
		if !before[msg] {
			return errors.New(msg)
		}
	}
	for _, p := range planned {
		if p.env != "" {
			resp.Warnings = append(resp.Warnings, p.env+" is set and overrides "+p.key)
			continue
		}
		got, _ := config.SettingValue(cfg, p.key)
		want := p.value
		if p.unset {
			want, _ = config.SettingValue(expected, p.key)
		}
		if !config.SettingsEqual(got, want) {
			return fmt.Errorf("%s did not take effect (the loaded config reports %v)", p.key, got)
		}
	}
	return nil
}

type plannedSetting struct {
	key   string
	path  []string
	value any
	unset bool
	// targets are the files a bake edits: the owning file for a change, or
	// every file that defines the key for a reset.
	targets []string
	env     string
}

func planSettingsChanges(state *settingsState, changes []settingsChange, allowEmpty bool) ([]plannedSetting, error) {
	if len(changes) == 0 && !allowEmpty {
		return nil, errors.New("no changes to bake")
	}
	if len(changes) > maxSettingsChanges {
		return nil, fmt.Errorf("too many changes (max %d)", maxSettingsChanges)
	}
	seen := map[string]bool{}
	planned := make([]plannedSetting, 0, len(changes))
	for _, change := range changes {
		field, ok := config.LookupSettingIn(state.cfg, change.Key)
		if !ok {
			return nil, fmt.Errorf("%w: %q", config.ErrUnknownSetting, change.Key)
		}
		if seen[change.Key] {
			return nil, fmt.Errorf("duplicate change for %s", change.Key)
		}
		seen[change.Key] = true
		path := field.SettingPath()
		p := plannedSetting{key: change.Key, path: path, env: settingEnvOverride(change.Key)}
		if change.Unset {
			if change.Value != nil {
				return nil, fmt.Errorf("%s: a reset takes no value", change.Key)
			}
			p.unset = true
			p.targets = state.definedIn(path)
			if len(p.targets) == 0 {
				return nil, fmt.Errorf("%s is not set in any config file; it already uses the default", change.Key)
			}
		} else {
			value, err := config.CoerceSettingValue(field, change.Value)
			if err != nil {
				return nil, err
			}
			target, err := state.target(path)
			if err != nil {
				return nil, err
			}
			p.value, p.targets = value, []string{target}
		}
		planned = append(planned, p)
	}
	return planned, nil
}

type fileBackup struct {
	path    string
	data    []byte
	existed bool
	mode    os.FileMode
}

func backupFile(path string) (fileBackup, error) {
	b := fileBackup{path: path}
	info, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return b, nil
	}
	if err != nil {
		return b, err
	}
	b.data, err = os.ReadFile(path)
	if err != nil {
		return b, err
	}
	b.existed, b.mode = true, info.Mode().Perm()
	return b, nil
}

func (b fileBackup) restore() error {
	if !b.existed {
		if err := os.Remove(b.path); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
		return nil
	}
	return config.WriteFileAtomic(b.path, b.data, b.mode)
}

func validationErrorSet(cfg *models.Config) map[string]bool {
	set := map[string]bool{}
	actual, _ := config.SplitErrorsAndWarnings(config.ValidateConfig(cfg))
	for _, err := range actual {
		set[err.Error()] = true
	}
	return set
}

// globalTargetError explains why a bake refuses the user-level config.
func globalTargetError(path string) error {
	return fmt.Errorf("this change would be written to your user-level config %s, which applies to every site; "+
		"create markata-go.toml in the site directory to bake settings for this site", path)
}

// applySettingsChanges writes changes grouped by target file. If any write,
// reload, validation, or effect check fails, every touched file is restored.
// A dry run returns the diffs without writing.
//
//nolint:gocyclo // Batch writes must report distinct validation, rollback, and preview failures.
func applySettingsChanges(state *settingsState, changes []settingsChange, dryRun bool) (resp settingsResponse, status int) {
	planned, err := planSettingsChanges(state, changes, false)
	if err != nil {
		resp.Error = err.Error()
		return resp, http.StatusBadRequest
	}

	var order []string
	byFile := map[string][]config.BakeSetting{}
	for _, p := range planned {
		value := p.value
		if p.unset {
			value = config.BakeRemove
		}
		for _, target := range p.targets {
			if state.kinds[target] == targetKindGlobal {
				resp.Error = globalTargetError(displayConfigPath(target)).Error()
				return resp, http.StatusConflict
			}
			if _, ok := byFile[target]; !ok {
				order = append(order, target)
			}
			byFile[target] = append(byFile[target], config.BakeSetting{Path: p.path, Value: value})
		}
		if p.env != "" {
			resp.Warnings = append(resp.Warnings, p.env+" is set and overrides "+p.key)
		}
	}

	plans := make([]config.BakePlan, 0, len(order))
	for _, path := range order {
		plan, err := config.PlanBake(path, byFile[path])
		if err != nil {
			resp.Error = err.Error()
			if errors.Is(err, config.ErrBakeUnsupportedLayout) {
				return resp, http.StatusConflict
			}
			return resp, http.StatusInternalServerError
		}
		plans = append(plans, plan)
		resp.Diffs = append(resp.Diffs, settingsDiff{
			Target:  displayConfigPath(path),
			Kind:    state.kinds[path],
			Created: !plan.Exists,
			Diff:    redactSensitiveLines(unifiedDiff(string(plan.Before), string(plan.After))),
		})
	}
	if dryRun {
		resp.Preview = servePreviewChanges()
		return resp, http.StatusOK
	}
	resp.Warnings = nil

	// The expected config is what the preview of these changes loads; resets
	// must land on the same values once the files are edited.
	expected, _, _, err := loadManagerConfigWith(cfgFile, previewConfig(planned))
	if err != nil {
		resp.Error = "changes do not load: " + err.Error()
		return resp, http.StatusUnprocessableEntity
	}

	backups := make([]fileBackup, 0, len(plans))
	rollback := func() {
		for i := len(backups) - 1; i >= 0; i-- {
			if err := backups[i].restore(); err != nil {
				resp.Warnings = append(resp.Warnings, "could not restore "+displayConfigPath(backups[i].path)+": "+err.Error())
			}
		}
	}
	for _, plan := range plans {
		backup, err := backupFile(plan.Path)
		if err != nil {
			rollback()
			resp.Error = err.Error()
			return resp, http.StatusInternalServerError
		}
		backups = append(backups, backup)
		if err := plan.Write(); err != nil {
			rollback()
			resp.Error = err.Error()
			return resp, http.StatusInternalServerError
		}
	}

	before := validationErrorSet(state.cfg)
	reloaded, _, err := resolveServeConfigSources(cfgFile)
	if err != nil {
		rollback()
		resp.Error = "config no longer loads, changes rolled back: " + err.Error()
		return resp, http.StatusUnprocessableEntity
	}
	if err := checkSettingsTookEffect(before, reloaded, expected, planned, &resp); err != nil {
		rollback()
		resp.Error = err.Error() + "; changes rolled back"
		return resp, http.StatusUnprocessableEntity
	}

	keysByFile := map[string][]string{}
	baked := map[string]bool{}
	for _, p := range planned {
		baked[p.key] = true
		for _, path := range p.targets {
			target := displayConfigPath(path)
			resp.Changed = append(resp.Changed, settingsChanged{Key: p.key, Target: target, Unset: p.unset})
			label := p.key
			if p.unset {
				label += " (reset)"
			}
			keysByFile[target] = append(keysByFile[target], label)
		}
	}
	for _, path := range order {
		target := displayConfigPath(path)
		infof("Baked settings into %s (%s)", target, strings.Join(keysByFile[target], ", "))
	}
	dropServePreviewKeys(baked)
	resp.Preview = servePreviewChanges()
	if serveRequestFullRebuild != nil {
		serveRequestFullRebuild()
	}
	return resp, http.StatusOK
}
