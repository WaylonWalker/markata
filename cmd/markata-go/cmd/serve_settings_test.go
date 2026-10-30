package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func doSettings(t *testing.T, method, body string, headers map[string]string) (*httptest.ResponseRecorder, settingsResponse) {
	t.Helper()
	req := httptest.NewRequest(method, "http://localhost:8000"+serveSettingsEndpoint, strings.NewReader(body))
	req.RemoteAddr = "127.0.0.1:52100"
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	handleSettings(rec, req)
	var resp settingsResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response %q: %v", rec.Body.String(), err)
	}
	return rec, resp
}

func settingsField(t *testing.T, resp settingsResponse, key string) settingsFieldJSON {
	t.Helper()
	for i := range resp.Fields {
		f := &resp.Fields[i]
		if f.Key == key {
			return *f
		}
	}
	t.Fatalf("setting %q not listed", key)
	return settingsFieldJSON{}
}

func readSiteFile(t *testing.T, site, name string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(site, name))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

var settingsTestFiles = map[string]string{
	"markata-go.toml":   "[markata-go]\ntitle = \"Site\" # name\nlicense = false\ninclude = [\"config/*.toml\"]\n",
	"config/theme.toml": "# theme\n[markata-go.theme]\npalette = \"nord-dark\"\n",
}

func TestHandleSettings_DescribesSourcesAndTargets(t *testing.T) {
	setupThemeBakeSite(t, settingsTestFiles)
	t.Setenv("MARKATA_GO_AUTHOR", "Env Author")

	rec, resp := doSettings(t, http.MethodGet, "", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET = %d %+v", rec.Code, resp)
	}
	if len(resp.Files) != 2 {
		t.Errorf("files = %v, want root and include", resp.Files)
	}
	title := settingsField(t, resp, "title")
	if title.Value != "Site" || title.Source != "markata-go.toml" || title.Target != "markata-go.toml" || !title.Editable {
		t.Errorf("title = %+v", title)
	}
	switcher := settingsField(t, resp, "theme.switcher.enabled")
	if switcher.Source != "" || switcher.Target != filepath.Join("config", "theme.toml") {
		t.Errorf("theme.switcher.enabled = %+v, want default written to config/theme.toml", switcher)
	}
	if nav := settingsField(t, resp, "nav"); nav.Editable || nav.Kind != "complex" {
		t.Errorf("nav = %+v, want read-only complex", nav)
	}
	if secret := settingsField(t, resp, "builder_admin.webhook.secret"); secret.Editable || !secret.Sensitive || secret.Value != nil {
		t.Errorf("secret = %+v, want hidden", secret)
	}
	if author := settingsField(t, resp, "author"); author.Env != "MARKATA_GO_AUTHOR" {
		t.Errorf("author env = %q", author.Env)
	}
	if fontpack := settingsField(t, resp, "theme.fontpack"); len(fontpack.Options) == 0 {
		t.Error("theme.fontpack has no options")
	}
}

func TestHandleSettings_BakesIntoOwningFiles(t *testing.T) {
	site := setupThemeBakeSite(t, settingsTestFiles)
	rebuilt := 0
	serveRequestFullRebuild = func() { rebuilt++ }

	body := `{"changes":[{"key":"title","value":"New"},{"key":"concurrency","value":4},` +
		`{"key":"theme.switcher.enabled","value":false},{"key":"glob.patterns","value":["posts/**/*.md"]}]}`
	rec, resp := doSettings(t, http.MethodPost, body, map[string]string{"Origin": "http://localhost:8000"})
	if rec.Code != http.StatusOK {
		t.Fatalf("POST = %d %+v", rec.Code, resp)
	}
	if rebuilt != 1 || len(resp.Changed) != 4 {
		t.Errorf("rebuilt = %d, changed = %+v", rebuilt, resp.Changed)
	}
	wantRoot := "[markata-go]\ntitle = \"New\" # name\nlicense = false\ninclude = [\"config/*.toml\"]\nconcurrency = 4\n\n[markata-go.glob]\npatterns = [\"posts/**/*.md\"]\n"
	if got := readSiteFile(t, site, "markata-go.toml"); got != wantRoot {
		t.Errorf("markata-go.toml:\n%s\nwant:\n%s", got, wantRoot)
	}
	wantTheme := "# theme\n[markata-go.theme]\npalette = \"nord-dark\"\n\n[markata-go.theme.switcher]\nenabled = false\n"
	if got := readSiteFile(t, site, "config/theme.toml"); got != wantTheme {
		t.Errorf("config/theme.toml:\n%s\nwant:\n%s", got, wantTheme)
	}
}

func TestHandleSettings_RollsBackInvalidChanges(t *testing.T) {
	site := setupThemeBakeSite(t, settingsTestFiles)
	serveRequestFullRebuild = func() { t.Error("rebuild requested for a rejected change") }

	tests := []struct {
		name   string
		body   string
		status int
	}{
		{"validation error", `{"changes":[{"key":"title","value":"x"},{"key":"url","value":"example.com"}]}`, http.StatusUnprocessableEntity},
		{"unknown option", `{"changes":[{"key":"title","value":"x"},{"key":"theme.fontpack","value":"not-a-pack"}]}`, http.StatusBadRequest},
		{"out of range", `{"changes":[{"key":"theme.motif.color_mix","value":2}]}`, http.StatusBadRequest},
		{"unknown key", `{"changes":[{"key":"no_such_key","value":"x"}]}`, http.StatusBadRequest},
		{"read-only key", `{"changes":[{"key":"nav","value":"x"}]}`, http.StatusBadRequest},
		{"sensitive key", `{"changes":[{"key":"builder_admin.webhook.secret","value":"x"}]}`, http.StatusBadRequest},
		{"wrong type", `{"changes":[{"key":"concurrency","value":"four"}]}`, http.StatusBadRequest},
		{"fractional int", `{"changes":[{"key":"concurrency","value":1.5}]}`, http.StatusBadRequest},
		{"duplicate", `{"changes":[{"key":"title","value":"a"},{"key":"title","value":"b"}]}`, http.StatusBadRequest},
		{"unknown field", `{"changes":[],"path":"/etc"}`, http.StatusBadRequest},
		{"empty", `{"changes":[]}`, http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec, resp := doSettings(t, http.MethodPost, tt.body, nil)
			if rec.Code != tt.status {
				t.Errorf("status = %d (%+v), want %d", rec.Code, resp, tt.status)
			}
		})
	}
	for name, want := range settingsTestFiles {
		if got := readSiteFile(t, site, name); got != want {
			t.Errorf("%s changed:\n%s", name, got)
		}
	}
}

func TestHandleSettings_RejectsCrossOriginAndRebinding(t *testing.T) {
	setupThemeBakeSite(t, settingsTestFiles)
	body := `{"changes":[{"key":"title","value":"x"}]}`
	if rec, _ := doSettings(t, http.MethodPost, body, map[string]string{"Origin": "https://evil.example"}); rec.Code != http.StatusForbidden {
		t.Errorf("cross-origin status = %d", rec.Code)
	}

	req := httptest.NewRequest(http.MethodGet, "http://rebind.evil.example:8000"+serveSettingsEndpoint, http.NoBody)
	req.RemoteAddr = "127.0.0.1:52100"
	req.Header.Set("Origin", "http://rebind.evil.example:8000")
	rec := httptest.NewRecorder()
	handleSettings(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Errorf("rebinding host status = %d, want 403", rec.Code)
	}

	lan := httptest.NewRequest(http.MethodPost, "http://192.168.1.5:8000"+serveSettingsEndpoint, strings.NewReader(body))
	lan.RemoteAddr = "192.168.1.20:40000"
	lan.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	handleSettings(rec, lan)
	if rec.Code != http.StatusForbidden {
		t.Errorf("non-loopback client status = %d, want 403", rec.Code)
	}
	for addr, want := range map[string]bool{"127.0.0.1:1": true, "[::1]:1": true, "192.168.1.20:1": false, "10.0.0.2:1": false, "": false} {
		if got := loopbackRemote(addr); got != want {
			t.Errorf("loopbackRemote(%q) = %v, want %v", addr, got, want)
		}
	}

	for _, host := range []string{"localhost:8000", "127.0.0.1:8000", "[::1]:8000", "192.168.1.5:8000", "site.localhost"} {
		if !localDevHost(host) {
			t.Errorf("localDevHost(%q) = false", host)
		}
	}
}

func TestInjectDevScripts_LoadsSettingsSidebar(t *testing.T) {
	html := injectDevScripts("<html><head></head><body></body></html>", BuildStatus{})
	if !strings.Contains(html, `window.__markataSettingsEndpoint = "`+serveSettingsEndpoint+`"`) ||
		!strings.Contains(html, `<script src="`+serveSettingsScriptPath+`" defer></script></body>`) {
		t.Errorf("dev scripts missing settings sidebar:\n%s", html)
	}
	if len(serveSettingsScript) == 0 || !strings.Contains(string(serveSettingsScript), "__markataSettingsEndpoint") {
		t.Error("settings script not embedded")
	}
}

func doSettingsPreview(t *testing.T, body string) (*httptest.ResponseRecorder, settingsResponse) {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "http://localhost:8000"+serveSettingsPreview, strings.NewReader(body))
	req.RemoteAddr = "127.0.0.1:52100"
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handleSettingsPreview(rec, req)
	var resp settingsResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response %q: %v", rec.Body.String(), err)
	}
	return rec, resp
}

func TestHandleSettingsPreview_AppliesWithoutWritingThenBakeAndReset(t *testing.T) {
	const root = "[markata-go]\ntitle = \"Disk\" # keep\nlicense = false\n"
	site := setupThemeBakeSite(t, map[string]string{"markata-go.toml": root})
	rebuilds := 0
	serveRequestFullRebuild = func() { rebuilds++ }

	rec, resp := doSettingsPreview(t, `{"changes":[{"key":"title","value":"Preview"},{"key":"concurrency","value":3}]}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("preview status %d: %s", rec.Code, resp.Error)
	}
	if len(resp.Preview) != 2 || rebuilds != 1 {
		t.Fatalf("preview = %v, rebuilds = %d", resp.Preview, rebuilds)
	}
	data, err := os.ReadFile(filepath.Join(site, "markata-go.toml"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != root {
		t.Fatalf("preview wrote the config file:\n%s", data)
	}
	cfg, _, _, err := loadManagerConfig("")
	if err != nil || cfg.Title != "Preview" || cfg.Concurrency != 3 {
		t.Fatalf("manager config does not include preview: %v %+v", err, cfg)
	}
	_, got := doSettings(t, http.MethodGet, "", nil)
	if f := settingsField(t, got, "title"); f.Value != "Disk" {
		t.Fatalf("GET should report saved value, got %v", f.Value)
	}
	if len(got.Preview) != 2 {
		t.Fatalf("GET preview = %v", got.Preview)
	}

	// An invalid preview is rejected and the previous preview is kept.
	rec, resp = doSettingsPreview(t, `{"changes":[{"key":"glob.slug_mode","value":"bogus"}]}`)
	if rec.Code != http.StatusUnprocessableEntity && rec.Code != http.StatusBadRequest {
		t.Fatalf("invalid preview status %d", rec.Code)
	}
	if len(resp.Preview) != 2 {
		t.Fatalf("invalid preview replaced state: %v", resp.Preview)
	}

	// Baking one key saves it and drops it from the preview.
	rec, resp = doSettings(t, http.MethodPost, `{"changes":[{"key":"title","value":"Preview"}]}`, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("bake status %d: %s", rec.Code, resp.Error)
	}
	if len(resp.Preview) != 1 || resp.Preview[0].Key != "concurrency" {
		t.Fatalf("preview after bake = %v", resp.Preview)
	}
	data, err = os.ReadFile(filepath.Join(site, "markata-go.toml"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "title = \"Preview\" # keep") {
		t.Fatalf("bake did not write title:\n%s", data)
	}

	// Reset clears the rest.
	rec, resp = doSettingsPreview(t, `{"changes":[]}`)
	if rec.Code != http.StatusOK || len(resp.Preview) != 0 {
		t.Fatalf("reset status %d preview %v", rec.Code, resp.Preview)
	}
	cfg, _, _, err = loadManagerConfig("")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Concurrency != 0 {
		t.Fatalf("reset left concurrency = %d", cfg.Concurrency)
	}
}

func TestHandleSettings_RenderProfilePreviewAndBake(t *testing.T) {
	site := setupThemeBakeSite(t, settingsTestFiles)
	serveRequestFullRebuild = func() {}

	rec, response := doSettings(t, http.MethodGet, "", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET settings = %d %+v", rec.Code, response)
	}
	field := settingsField(t, response, "markdown.render_profile")
	if !field.Editable {
		t.Fatalf("render profile is not editable: %+v", field)
	}

	change := `{"changes":[{"key":"markdown.render_profile","value":"offline-interactive"}]}`
	rec, response = doSettings(t, http.MethodPost, change, map[string]string{"Origin": "http://localhost:8000"})
	if rec.Code != http.StatusOK {
		t.Fatalf("bake profile = %d %+v", rec.Code, response)
	}
	if source := readSiteFile(t, site, "markata-go.toml"); !strings.Contains(source, `render_profile = "offline-interactive"`) {
		t.Fatalf("profile was not persisted to config: %s", source)
	}
	rec, response = doSettings(t, http.MethodGet, "", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("reload settings = %d %+v", rec.Code, response)
	}
	if field := settingsField(t, response, "markdown.render_profile"); field.Value != "offline-interactive" {
		t.Fatalf("baked render profile did not reload: %+v", field)
	}
}

func TestHandleSettings_RenderAssetDefaultsBake(t *testing.T) {
	site := setupThemeBakeSite(t, settingsTestFiles)
	serveRequestFullRebuild = func() {}
	payload := `{"changes":[{"key":"markdown.render_asset_mode","value":"inline"},{"key":"markdown.render_asset_base","value":"/local-md-assets"}]}`
	rec, resp := doSettings(t, http.MethodPost, payload, map[string]string{"Origin":"http://localhost:8000"})
	if rec.Code != http.StatusOK {
		t.Fatalf("bake renderer asset defaults = %d %+v", rec.Code, resp)
	}
	configText := readSiteFile(t, site, "markata-go.toml")
	for _, wanted := range []string{`render_asset_mode = "inline"`, `render_asset_base = "/local-md-assets"`} {
		if !strings.Contains(configText, wanted) { t.Fatalf("bake lost %s: %s", wanted, configText) }
	}
	rec, resp = doSettings(t, http.MethodGet, "", nil)
	if rec.Code != http.StatusOK { t.Fatal(rec.Body.String()) }
	if settingField(t, resp, "markdown.render_asset_mode").Value != "inline" {
		t.Fatal("asset delivery mode did not reload")
	}
	if settingField(t, resp, "markdown.render_asset_base").Value != "/local-md-assets" {
		t.Fatal("asset base did not reload")
	}
}
