// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	assetsPkg "github.com/WaylonWalker/markata/pkg/assets"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/runtimeenv"
	"github.com/yuin/goldmark"
)

func TestWebAwesomePlugin_ProcessComparisonContainer(t *testing.T) {
	plugin := NewWebAwesomePlugin()
	post := &models.Post{
		ArticleHTML: `<div position="35" caption="Homepage redesign" class="webawesome comparison">
<p><img src="/before.webp" alt="Before image"> <img src="/after.webp" alt="After image"></p>
</div>`,
	}

	if err := plugin.processPost(post); err != nil {
		t.Fatalf("processPost() error = %v", err)
	}

	want := []string{
		`<figure class="markata-webawesome-figure">`,
		`<wa-comparison class="markata-webawesome-comparison" position="35">`,
		`<img slot="after" src="/before.webp" alt="Before image" loading="lazy">`,
		`<img slot="before" src="/after.webp" alt="After image" loading="lazy">`,
		`<figcaption>Homepage redesign</figcaption>`,
	}
	for _, expected := range want {
		if !strings.Contains(post.ArticleHTML, expected) {
			t.Fatalf("ArticleHTML missing %q\nGot: %s", expected, post.ArticleHTML)
		}
	}
}

func TestWebAwesomePlugin_ProcessComparisonContainerWithFigure(t *testing.T) {
	plugin := NewWebAwesomePlugin()
	post := &models.Post{
		ArticleHTML: `<div class="webawesome comparison" position="42" caption="Compare the same generated graphic before and after a color treatment.">
<figure>
<img src="/before.webp" alt="Before image">
<img src="/after.webp" alt="After image">
</figure>
</div>`,
	}

	if err := plugin.processPost(post); err != nil {
		t.Fatalf("processPost() error = %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `<wa-comparison class="markata-webawesome-comparison" position="42">`) {
		t.Fatalf("ArticleHTML missing wa-comparison\nGot: %s", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, `<figcaption>Compare the same generated graphic before and after a color treatment.</figcaption>`) {
		t.Fatalf("ArticleHTML missing full caption\nGot: %s", post.ArticleHTML)
	}
}

func TestWebAwesomePlugin_ProcessComparisonContainerWithLinkedImages(t *testing.T) {
	plugin := NewWebAwesomePlugin()
	post := &models.Post{
		ArticleHTML: `<div class="wa-comparison">
<figure>
<a href="/before.webp" class="glightbox-link"><img class="glightbox" src="/before.webp" alt="Before image" data-glightbox="description: Before image"></a>
<a href="/after.webp" class="glightbox-link"><img class="glightbox" src="/after.webp" alt="After image" data-glightbox="description: After image"></a>
</figure>
</div>`,
	}

	if err := plugin.processPost(post); err != nil {
		t.Fatalf("processPost() error = %v", err)
	}

	want := []string{
		`<wa-comparison class="markata-webawesome-comparison">`,
		`<img slot="after" src="/before.webp" alt="Before image" loading="lazy">`,
		`<img slot="before" src="/after.webp" alt="After image" loading="lazy">`,
	}
	for _, expected := range want {
		if !strings.Contains(post.ArticleHTML, expected) {
			t.Fatalf("ArticleHTML missing %q\nGot: %s", expected, post.ArticleHTML)
		}
	}

	if strings.Contains(post.ArticleHTML, `glightbox-link`) {
		t.Fatalf("ArticleHTML should not retain glightbox wrappers\nGot: %s", post.ArticleHTML)
	}
}

func TestWebAwesomePlugin_ProcessUsefulContentContainers(t *testing.T) {
	plugin := NewWebAwesomePlugin()
	post := &models.Post{ArticleHTML: `<div class="webawesome details" summary="Install notes">
<p>Use the binary for your platform.</p>
</div>
<div class="wa-tabs">
<div class="wa-tab" label="macOS">
<pre><code>brew install markata-go</code></pre>
</div>
<div class="wa-tab" label="Linux">
<pre><code>curl -fsSL example.com</code></pre>
</div>
</div>
<div class="webawesome copy"><p>go test ./...</p></div>
<div class="webawesome qr"><p>https://example.com/post</p></div>
<div class="webawesome badge" variant="brand"><p>New</p></div>
<div class="webawesome tag" variant="success"><p>Stable</p></div>
<div class="webawesome tooltip" content="Static Site Generator"><p>SSG</p></div>
<div class="webawesome carousel" navigation="true">
<figure><img src="/one.webp" alt="One"><img src="/two.webp" alt="Two"></figure>
</div>
<div class="webawesome animated-image">
<p><img src="/demo.webp" alt="Animation demo"></p>
</div>`}

	if err := plugin.processPost(post); err != nil {
		t.Fatalf("processPost() error = %v", err)
	}

	want := []string{
		`<wa-details summary="Install notes">`,
		`<wa-tab-group>`,
		`<wa-tab slot="nav" panel="macos">macOS</wa-tab>`,
		`<wa-tab-panel name="linux">`,
		`<wa-copy-button value="go test ./..."></wa-copy-button>`,
		`<wa-qr-code value="https://example.com/post"></wa-qr-code>`,
		`<wa-badge variant="brand">New</wa-badge>`,
		`<wa-tag variant="success">Stable</wa-tag>`,
		`<span class="markata-wa-tooltip-anchor"`,
		`>SSG</span><wa-tooltip for="`,
		`>Static Site Generator</wa-tooltip>`,
		`<wa-carousel navigation="true">`,
		`<wa-carousel-item><img src="/one.webp" alt="One" loading="lazy" decoding="async" style="display: block; max-width: 100%; max-height: 100%; width: auto !important; height: auto !important; object-fit: contain !important; border-radius: var(--radius-lg);"></wa-carousel-item>`,
		`<wa-animated-image src="/demo.webp" alt="Animation demo"></wa-animated-image>`,
	}
	for _, expected := range want {
		if !strings.Contains(post.ArticleHTML, expected) {
			t.Fatalf("ArticleHTML missing %q\nGot: %s", expected, post.ArticleHTML)
		}
	}
}

func TestWebAwesomePlugin_RenderProcessesWaCarouselContainerWithLeadingClass(t *testing.T) {
	plugin := NewWebAwesomePlugin()
	plugin.config.Source = "cdn"
	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{Extra: map[string]interface{}{}})
	m.SetPosts([]*models.Post{{ArticleHTML: `<div class="not-prose wa-carousel" navigation="true" pagination="true">
<figure><img src="/one.webp" alt="One"><img src="/two.webp" alt="Two"></figure>
</div>`}})

	if err := plugin.Render(m); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	got := m.Posts()[0].ArticleHTML
	want := []string{
		`<wa-carousel navigation="true" pagination="true">`,
		`<wa-carousel-item><img src="/one.webp" alt="One" loading="lazy" decoding="async" style="display: block; max-width: 100%; max-height: 100%; width: auto !important; height: auto !important; object-fit: contain !important; border-radius: var(--radius-lg);"></wa-carousel-item>`,
		`<wa-carousel-item><img src="/two.webp" alt="Two" loading="lazy" decoding="async" style="display: block; max-width: 100%; max-height: 100%; width: auto !important; height: auto !important; object-fit: contain !important; border-radius: var(--radius-lg);"></wa-carousel-item>`,
	}
	for _, expected := range want {
		if !strings.Contains(got, expected) {
			t.Fatalf("ArticleHTML missing %q\nGot: %s", expected, got)
		}
	}
	if enabled, ok := m.Config().Extra["webawesome_enabled"].(bool); !ok || !enabled {
		t.Fatalf("webawesome_enabled = %v, want true", m.Config().Extra["webawesome_enabled"])
	}
}

func TestWebAwesomePlugin_ProcessCarouselCaptions(t *testing.T) {
	plugin := NewWebAwesomePlugin()
	post := &models.Post{ArticleHTML: `<div class="wa-carousel" navigation="true" pagination="true">
<p><img src="/family.webp" alt="Family">
family
<img src="/ayla.webp" alt="Ayla">
<figcaption>Ayla 11</figcaption>
<img src="/wyatt.webp" alt="Wyatt">
Wyatt 13</p>
</div>`}

	if err := plugin.processPost(post); err != nil {
		t.Fatalf("processPost() error = %v", err)
	}

	want := []string{
		`<wa-carousel navigation="true" pagination="true">`,
		`<figure class="markata-wa-carousel-figure"><img src="/family.webp" alt="Family" loading="lazy" decoding="async" style="display: block; max-width: 100%; max-height: 100%; width: auto !important; height: auto !important; object-fit: contain !important; border-radius: var(--radius-lg);"><figcaption>family</figcaption></figure>`,
		`<figure class="markata-wa-carousel-figure"><img src="/ayla.webp" alt="Ayla" loading="lazy" decoding="async" style="display: block; max-width: 100%; max-height: 100%; width: auto !important; height: auto !important; object-fit: contain !important; border-radius: var(--radius-lg);"><figcaption>Ayla 11</figcaption></figure>`,
		`<figure class="markata-wa-carousel-figure"><img src="/wyatt.webp" alt="Wyatt" loading="lazy" decoding="async" style="display: block; max-width: 100%; max-height: 100%; width: auto !important; height: auto !important; object-fit: contain !important; border-radius: var(--radius-lg);"><figcaption>Wyatt 13</figcaption></figure>`,
	}
	for _, expected := range want {
		if !strings.Contains(post.ArticleHTML, expected) {
			t.Fatalf("ArticleHTML missing %q\nGot: %s", expected, post.ArticleHTML)
		}
	}
}

func TestWebAwesomePlugin_ProcessNestedTabsPreservesNestedDivs(t *testing.T) {
	plugin := NewWebAwesomePlugin()
	post := &models.Post{ArticleHTML: `<div class="wa-tabs">
<div class="wa-tab" label="macOS">
<div class="inner"><p>brew install markata-go</p></div>
</div>
<div class="wa-tab" label="Linux">
<div class="inner"><p>curl -fsSL example.com</p></div>
</div>
</div>`}

	if err := plugin.processPost(post); err != nil {
		t.Fatalf("processPost() error = %v", err)
	}

	want := []string{
		`<wa-tab-group>`,
		`<wa-tab slot="nav" panel="macos">macOS</wa-tab>`,
		`<wa-tab-panel name="macos"><div class="inner"><p>brew install markata-go</p></div></wa-tab-panel>`,
		`<wa-tab-panel name="linux"><div class="inner"><p>curl -fsSL example.com</p></div></wa-tab-panel>`,
	}
	for _, expected := range want {
		if !strings.Contains(post.ArticleHTML, expected) {
			t.Fatalf("ArticleHTML missing %q\nGot: %s", expected, post.ArticleHTML)
		}
	}
}

func TestWebAwesomePlugin_TooltipsUseUniqueAnchorIDs(t *testing.T) {
	plugin := NewWebAwesomePlugin()
	post := &models.Post{ArticleHTML: `<div class="webawesome tooltip" content="Static Site Generator"><p>SSG</p></div>
<div class="webawesome tooltip" content="Static Site Generator"><p>SSG</p></div>`}

	if err := plugin.processPost(post); err != nil {
		t.Fatalf("processPost() error = %v", err)
	}

	re := regexp.MustCompile(`<span class="markata-wa-tooltip-anchor" id="([^"]+)" tabindex="0">`)
	matches := re.FindAllStringSubmatch(post.ArticleHTML, -1)
	if len(matches) != 2 {
		t.Fatalf("expected 2 tooltip anchors, got %d in %s", len(matches), post.ArticleHTML)
	}
	if matches[0][1] == matches[1][1] {
		t.Fatalf("tooltip ids collided: %q", matches[0][1])
	}
	for _, match := range matches {
		if !strings.Contains(post.ArticleHTML, fmt.Sprintf(`<wa-tooltip for=%q>Static Site Generator</wa-tooltip>`, match[1])) {
			t.Fatalf("missing tooltip target for anchor %q in %s", match[1], post.ArticleHTML)
		}
	}
}

func TestWebAwesomePlugin_RenderEnablesAssetsForRawComponent(t *testing.T) {
	plugin := NewWebAwesomePlugin()
	plugin.config.Source = "cdn"
	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{Extra: map[string]interface{}{}})
	m.SetPosts([]*models.Post{{ArticleHTML: `<wa-button>Click</wa-button>`}})

	if err := plugin.Render(m); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if got, ok := m.Config().Extra["webawesome_enabled"].(bool); !ok || !got {
		t.Fatalf("webawesome_enabled = %v, want true", m.Config().Extra["webawesome_enabled"])
	}
	if got := m.Config().Extra["webawesome_css_url"]; got != webAwesomeDefaultCDNBase+"/styles/themes/default.css" {
		t.Fatalf("webawesome_css_url = %v", got)
	}
	if got := m.Config().Extra["webawesome_loader_url"]; got != webAwesomeDefaultCDNBase+"/webawesome.loader.js" {
		t.Fatalf("webawesome_loader_url = %v", got)
	}
	if needs, ok := m.Posts()[0].Extra["needs_webawesome"].(bool); !ok || !needs {
		t.Fatalf("needs_webawesome = %v, want true", m.Posts()[0].Extra["needs_webawesome"])
	}
}

func TestWebAwesomePlugin_RenderUsesSharedVendorAssetURL(t *testing.T) {
	plugin := NewWebAwesomePlugin()
	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{Extra: map[string]interface{}{
		"asset_urls": map[string]string{
			"webawesome": "/assets/vendor/webawesome",
		},
	}})
	m.SetPosts([]*models.Post{{ArticleHTML: `<wa-button>Click</wa-button>`}})

	if err := plugin.Render(m); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if got := m.Config().Extra["webawesome_css_url"]; got != "/assets/vendor/webawesome/styles/themes/default.css" {
		t.Fatalf("webawesome_css_url = %v", got)
	}
	if got := m.Config().Extra["webawesome_loader_url"]; got != "/assets/vendor/webawesome/webawesome.loader.js" {
		t.Fatalf("webawesome_loader_url = %v", got)
	}
}

func TestWebAwesomePlugin_ConfigureRegistersDefaultVendorAssetWithoutConfig(t *testing.T) {
	plugin := NewWebAwesomePlugin()
	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{})

	if err := plugin.Configure(m); err != nil {
		t.Fatalf("Configure() error = %v", err)
	}

	extraAssets, ok := m.Config().Extra["cdn_assets_extra"].([]interface{})
	if ok {
		t.Fatalf("cdn_assets_extra has unexpected type []interface{}: %#v", extraAssets)
	}
	assets, ok := m.Config().Extra["cdn_assets_extra"].([]assetsPkg.Asset)
	if !ok {
		t.Fatalf("cdn_assets_extra type = %T, want []assets.Asset", m.Config().Extra["cdn_assets_extra"])
	}
	if len(assets) != 1 {
		t.Fatalf("len(cdn_assets_extra) = %d, want 1", len(assets))
	}
	if assets[0].Name != webAwesomeAssetName {
		t.Fatalf("asset name = %q, want %q", assets[0].Name, webAwesomeAssetName)
	}
	if assets[0].Integrity != assetsPkg.NewWebAwesomeAsset(webAwesomeDefaultVersion, "").Integrity {
		t.Fatalf("integrity = %q, want default SRI", assets[0].Integrity)
	}
}

func TestWebAwesomePlugin_OfflineDefaultUsesCDNWithoutRequiredVendorAsset(t *testing.T) {
	t.Setenv(runtimeenv.EnvOffline, "true")

	plugin := NewWebAwesomePlugin()
	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{})

	if err := plugin.Configure(m); err != nil {
		t.Fatalf("Configure() error = %v", err)
	}

	if plugin.config.Source != "cdn" {
		t.Fatalf("source = %q, want cdn", plugin.config.Source)
	}
	if _, ok := m.Config().Extra["cdn_assets_extra"]; ok {
		t.Fatalf("cdn_assets_extra should not be registered by default while offline: %#v", m.Config().Extra["cdn_assets_extra"])
	}
}

func TestWebAwesomePlugin_OfflineExplicitVendorStillRequiresVendorAsset(t *testing.T) {
	t.Setenv(runtimeenv.EnvOffline, "true")

	plugin := NewWebAwesomePlugin()
	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{Extra: map[string]interface{}{
		"webawesome": map[string]interface{}{
			"source": "vendor",
		},
	}})

	if err := plugin.Configure(m); err != nil {
		t.Fatalf("Configure() error = %v", err)
	}

	assets, ok := m.Config().Extra["cdn_assets_extra"].([]assetsPkg.Asset)
	if !ok || len(assets) != 1 {
		t.Fatalf("cdn_assets_extra = %#v, want one vendor asset", m.Config().Extra["cdn_assets_extra"])
	}
	if assets[0].Name != webAwesomeAssetName {
		t.Fatalf("asset name = %q, want %q", assets[0].Name, webAwesomeAssetName)
	}
}

func TestWebAwesomePlugin_OfflineExplicitSelfHostedAssetsStillRequiresVendorAsset(t *testing.T) {
	t.Setenv(runtimeenv.EnvOffline, "true")

	plugin := NewWebAwesomePlugin()
	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{Extra: map[string]interface{}{
		"assets": map[string]interface{}{
			"mode": "self-hosted",
		},
	}})

	if err := plugin.Configure(m); err != nil {
		t.Fatalf("Configure() error = %v", err)
	}

	assets, ok := m.Config().Extra["cdn_assets_extra"].([]assetsPkg.Asset)
	if !ok || len(assets) != 1 {
		t.Fatalf("cdn_assets_extra = %#v, want one vendor asset", m.Config().Extra["cdn_assets_extra"])
	}
	if assets[0].Name != webAwesomeAssetName {
		t.Fatalf("asset name = %q, want %q", assets[0].Name, webAwesomeAssetName)
	}
}

func TestWebAwesomePlugin_ConfigureOmitsIntegrityForVersionOverride(t *testing.T) {
	plugin := NewWebAwesomePlugin()
	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{Extra: map[string]interface{}{
		"webawesome": map[string]interface{}{
			"version": "3.6.0",
		},
	}})

	if err := plugin.Configure(m); err != nil {
		t.Fatalf("Configure() error = %v", err)
	}

	assets, ok := m.Config().Extra["cdn_assets_extra"].([]assetsPkg.Asset)
	if !ok || len(assets) != 1 {
		t.Fatalf("cdn_assets_extra = %#v, want one asset", m.Config().Extra["cdn_assets_extra"])
	}
	if assets[0].Version != "3.6.0" {
		t.Fatalf("asset version = %q, want %q", assets[0].Version, "3.6.0")
	}
	if assets[0].Integrity != "" {
		t.Fatalf("integrity = %q, want empty for version override", assets[0].Integrity)
	}
}

func TestWebAwesomePlugin_RenderDoesNotEnableAssetsWithoutComponents(t *testing.T) {
	plugin := NewWebAwesomePlugin()
	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{Extra: map[string]interface{}{}})
	m.SetPosts([]*models.Post{{ArticleHTML: `<p>No Web Awesome here.</p>`}})

	if err := plugin.Render(m); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if _, ok := m.Config().Extra["webawesome_enabled"]; ok {
		t.Fatalf("webawesome_enabled was set for a page without Web Awesome components")
	}
	if m.Posts()[0].Extra != nil {
		if _, ok := m.Posts()[0].Extra["needs_webawesome"]; ok {
			t.Fatalf("needs_webawesome was set for a page without Web Awesome components")
		}
	}
}

func TestWebAwesomePlugin_ComponentModulesDeduplicatesDetectedComponents(t *testing.T) {
	plugin := NewWebAwesomePlugin()
	plugin.config.Source = "cdn"
	modules := plugin.componentModules(`<wa-comparison></wa-comparison><wa-button></wa-button><wa-button></wa-button>`)
	want := []string{
		webAwesomeDefaultCDNBase + "/components/comparison/comparison.js",
		webAwesomeDefaultCDNBase + "/components/button/button.js",
	}
	if len(modules) != len(want) {
		t.Fatalf("modules = %#v, want %#v", modules, want)
	}
	for i := range want {
		if modules[i] != want[i] {
			t.Fatalf("modules = %#v, want %#v", modules, want)
		}
	}
}

func TestWebAwesomePlugin_DefaultsToVendorSource(t *testing.T) {
	plugin := NewWebAwesomePlugin()
	if plugin.config.Source != "vendor" {
		t.Fatalf("default source = %q, want %q", plugin.config.Source, "vendor")
	}
	if plugin.config.OutputDir != "assets/vendor/webawesome" {
		t.Fatalf("default output dir = %q, want %q", plugin.config.OutputDir, "assets/vendor/webawesome")
	}
}

func TestWebAwesomePlugin_ProcessGenericContainerWithNestedDiv(t *testing.T) {
	plugin := NewWebAwesomePlugin()
	post := &models.Post{ArticleHTML: `<div class="webawesome details" summary="Install notes">
<div class="inner"><p>Use the binary for your platform.</p></div>
</div>`}

	if err := plugin.processPost(post); err != nil {
		t.Fatalf("processPost() error = %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `<wa-details summary="Install notes">`) {
		t.Fatalf("missing wa-details wrapper: %s", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, `<div class="inner"><p>Use the binary for your platform.</p></div>`) {
		t.Fatalf("nested div content was not preserved: %s", post.ArticleHTML)
	}
}

func TestContainerExtension_ParsesQuotedAttributes(t *testing.T) {
	md := goldmark.New(goldmark.WithExtensions(&ContainerExtension{}))
	input := `::: webawesome details {summary="Install notes" data-kind="quick start"}
Use the binary.
:::`

	var buf bytes.Buffer
	if err := md.Convert([]byte(input), &buf); err != nil {
		t.Fatalf("Convert() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, `class="webawesome details"`) {
		t.Fatalf("missing classes: %q", output)
	}
	if !strings.Contains(output, `summary="Install notes"`) {
		t.Fatalf("missing quoted summary attr: %q", output)
	}
	if !strings.Contains(output, `data-kind="quick start"`) {
		t.Fatalf("missing quoted data attr: %q", output)
	}
}

func TestWebAwesomeVendorIntegration_DefaultConfigDownloadsAndUsesSharedAssetURL(t *testing.T) {
	archiveData := buildWebAwesomeTestTarGz(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(archiveData); err != nil {
			t.Logf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	plugin := NewWebAwesomePlugin()
	cdn := NewCDNAssetsPlugin()
	m := lifecycle.NewManager()
	cacheDir := t.TempDir()
	outputDir := t.TempDir()
	assetsConfig := models.NewAssetsConfig()
	assetsConfig.Mode = "cdn"
	assetsConfig.CacheDir = cacheDir
	assetsConfig.OutputDir = "assets/vendor"
	m.SetConfig(&lifecycle.Config{
		OutputDir: outputDir,
		Extra: map[string]interface{}{
			"assets": assetsConfig,
		},
	})
	m.SetPosts([]*models.Post{{ArticleHTML: `<wa-button>Click</wa-button>`}})

	if err := plugin.Configure(m); err != nil {
		t.Fatalf("webawesome Configure() error = %v", err)
	}
	setRequestedAssetURL(t, m.Config(), server.URL+"/webawesome-3.5.0.tgz")

	if err := cdn.Configure(m); err != nil {
		t.Fatalf("cdn_assets Configure() error = %v", err)
	}

	assetURLs, ok := m.Config().Extra["asset_urls"].(map[string]string)
	if !ok {
		t.Fatalf("asset_urls type = %T, want map[string]string", m.Config().Extra["asset_urls"])
	}
	if got := assetURLs[webAwesomeAssetName]; got != "/assets/vendor/webawesome" {
		t.Fatalf("asset_urls[webawesome] = %q, want %q", got, "/assets/vendor/webawesome")
	}
	assertFileExists(t, filepath.Join(cacheDir, "webawesome", "webawesome.loader.js"))
	assertFileExists(t, filepath.Join(cacheDir, "webawesome", "styles", "webawesome.css"))

	if err := plugin.Render(m); err != nil {
		t.Fatalf("webawesome Render() error = %v", err)
	}
	if got := m.Config().Extra["webawesome_css_url"]; got != "/assets/vendor/webawesome/styles/themes/default.css" {
		t.Fatalf("webawesome_css_url = %v", got)
	}
	if got := m.Config().Extra["webawesome_loader_url"]; got != "/assets/vendor/webawesome/webawesome.loader.js" {
		t.Fatalf("webawesome_loader_url = %v", got)
	}

	if err := cdn.Write(m); err != nil {
		t.Fatalf("cdn_assets Write() error = %v", err)
	}
	assertFileExists(t, filepath.Join(outputDir, "assets", "vendor", "webawesome", "webawesome.loader.js"))
	assertFileExists(t, filepath.Join(outputDir, "assets", "vendor", "webawesome", "styles", "webawesome.css"))
}

func TestWebAwesomeVendorIntegration_CustomOutputDirIsPublishedAndUsed(t *testing.T) {
	archiveData := buildWebAwesomeTestTarGz(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(archiveData); err != nil {
			t.Logf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	plugin := NewWebAwesomePlugin()
	cdn := NewCDNAssetsPlugin()
	m := lifecycle.NewManager()
	cacheDir := t.TempDir()
	outputDir := t.TempDir()
	assetsConfig := models.NewAssetsConfig()
	assetsConfig.Mode = "cdn"
	assetsConfig.CacheDir = cacheDir
	assetsConfig.OutputDir = "assets/vendor"
	m.SetConfig(&lifecycle.Config{
		OutputDir: outputDir,
		Extra: map[string]interface{}{
			"assets": assetsConfig,
			"webawesome": map[string]interface{}{
				"output_dir": "assets/vendor/wa-kit",
			},
		},
	})
	m.SetPosts([]*models.Post{{ArticleHTML: `<wa-button>Click</wa-button>`}})

	if err := plugin.Configure(m); err != nil {
		t.Fatalf("webawesome Configure() error = %v", err)
	}
	setRequestedAssetURL(t, m.Config(), server.URL+"/webawesome-3.5.0.tgz")

	if err := cdn.Configure(m); err != nil {
		t.Fatalf("cdn_assets Configure() error = %v", err)
	}

	assetURLs, ok := m.Config().Extra["asset_urls"].(map[string]string)
	if !ok {
		t.Fatalf("asset_urls type = %T, want map[string]string", m.Config().Extra["asset_urls"])
	}
	if got := assetURLs[webAwesomeAssetName]; got != "/assets/vendor/wa-kit" {
		t.Fatalf("asset_urls[webawesome] = %q, want %q", got, "/assets/vendor/wa-kit")
	}

	if err := plugin.Render(m); err != nil {
		t.Fatalf("webawesome Render() error = %v", err)
	}
	if got := m.Config().Extra["webawesome_css_url"]; got != "/assets/vendor/wa-kit/styles/themes/default.css" {
		t.Fatalf("webawesome_css_url = %v", got)
	}
	if got := m.Config().Extra["webawesome_loader_url"]; got != "/assets/vendor/wa-kit/webawesome.loader.js" {
		t.Fatalf("webawesome_loader_url = %v", got)
	}

	if err := cdn.Write(m); err != nil {
		t.Fatalf("cdn_assets Write() error = %v", err)
	}
	assertFileExists(t, filepath.Join(outputDir, "assets", "vendor", "wa-kit", "webawesome.loader.js"))
	assertFileExists(t, filepath.Join(outputDir, "assets", "vendor", "wa-kit", "styles", "webawesome.css"))
	if _, err := os.Stat(filepath.Join(outputDir, "assets", "vendor", "webawesome", "webawesome.loader.js")); !os.IsNotExist(err) {
		t.Fatalf("webawesome asset unexpectedly published at default path")
	}
}

func TestWebAwesomeVendorIntegration_VersionOverrideDownloadsWithoutIntegrity(t *testing.T) {
	archiveData := buildWebAwesomeTestTarGz(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(archiveData); err != nil {
			t.Logf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	plugin := NewWebAwesomePlugin()
	cdn := NewCDNAssetsPlugin()
	m := lifecycle.NewManager()
	cacheDir := t.TempDir()
	assetsConfig := models.NewAssetsConfig()
	assetsConfig.Mode = "cdn"
	assetsConfig.CacheDir = cacheDir
	m.SetConfig(&lifecycle.Config{Extra: map[string]interface{}{
		"assets": assetsConfig,
		"webawesome": map[string]interface{}{
			"version": "9.9.9",
		},
	}})
	m.SetPosts([]*models.Post{{ArticleHTML: `<wa-button>Click</wa-button>`}})

	if err := plugin.Configure(m); err != nil {
		t.Fatalf("webawesome Configure() error = %v", err)
	}
	requested := setRequestedAssetURL(t, m.Config(), server.URL+"/webawesome-9.9.9.tgz")
	if requested.Integrity != "" {
		t.Fatalf("integrity = %q, want empty for version override", requested.Integrity)
	}

	if err := cdn.Configure(m); err != nil {
		t.Fatalf("cdn_assets Configure() error = %v", err)
	}

	assetURLs, ok := m.Config().Extra["asset_urls"].(map[string]string)
	if !ok {
		t.Fatalf("asset_urls type = %T, want map[string]string", m.Config().Extra["asset_urls"])
	}
	if got := assetURLs[webAwesomeAssetName]; got != "/assets/vendor/webawesome" {
		t.Fatalf("asset_urls[webawesome] = %q, want %q", got, "/assets/vendor/webawesome")
	}
	assertFileExists(t, filepath.Join(cacheDir, "webawesome", "webawesome.loader.js"))
}

func TestWebAwesomeVendorIntegration_BadIntegrityFallsBackToCDN(t *testing.T) {
	archiveData := buildWebAwesomeTestTarGz(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(archiveData); err != nil {
			t.Logf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	plugin := NewWebAwesomePlugin()
	cdn := NewCDNAssetsPlugin()
	m := lifecycle.NewManager()
	cacheDir := t.TempDir()
	assetsConfig := models.NewAssetsConfig()
	assetsConfig.Mode = "cdn"
	assetsConfig.CacheDir = cacheDir
	m.SetConfig(&lifecycle.Config{Extra: map[string]interface{}{
		"assets": assetsConfig,
	}})
	m.SetPosts([]*models.Post{{ArticleHTML: `<wa-button>Click</wa-button>`}})

	if err := plugin.Configure(m); err != nil {
		t.Fatalf("webawesome Configure() error = %v", err)
	}
	requested := setRequestedAssetURL(t, m.Config(), server.URL+"/webawesome-3.5.0.tgz")
	requested.Integrity = "sha512-invalid"
	setRequestedAsset(t, m.Config(), requested)

	if err := cdn.Configure(m); err != nil {
		t.Fatalf("cdn_assets Configure() error = %v", err)
	}

	assetURLs, ok := m.Config().Extra["asset_urls"].(map[string]string)
	if !ok {
		t.Fatalf("asset_urls type = %T, want map[string]string", m.Config().Extra["asset_urls"])
	}
	if _, ok := assetURLs[webAwesomeAssetName]; ok {
		t.Fatalf("asset_urls should not include webawesome when download fails: %#v", assetURLs)
	}
	if _, err := os.Stat(filepath.Join(cacheDir, "webawesome", "webawesome.loader.js")); !os.IsNotExist(err) {
		t.Fatalf("webawesome asset unexpectedly cached after integrity failure")
	}

	if err := plugin.Render(m); err != nil {
		t.Fatalf("webawesome Render() error = %v", err)
	}
	if got := m.Config().Extra["webawesome_css_url"]; got != webAwesomeDefaultCDNBase+"/styles/themes/default.css" {
		t.Fatalf("webawesome_css_url = %v", got)
	}
	if got := m.Config().Extra["webawesome_loader_url"]; got != webAwesomeDefaultCDNBase+"/webawesome.loader.js" {
		t.Fatalf("webawesome_loader_url = %v", got)
	}
}

func buildWebAwesomeTestTarGz(t *testing.T) []byte {
	t.Helper()
	return buildWebAwesomeArchive(t, map[string]string{
		"package/dist-cdn/webawesome.loader.js":      "console.log('loader')",
		"package/dist-cdn/styles/webawesome.css":     "body { color: red; }",
		"package/dist-cdn/styles/themes/default.css": "html { color: red; }",
	})
}

func buildWebAwesomeArchive(t *testing.T, files map[string]string) []byte {
	t.Helper()

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)

	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		content := files[name]
		header := &tar.Header{Name: name, Mode: 0o644, Size: int64(len(content))}
		if err := tw.WriteHeader(header); err != nil {
			t.Fatalf("write tar header: %v", err)
		}
		if _, err := tw.Write([]byte(content)); err != nil {
			t.Fatalf("write tar content: %v", err)
		}
	}

	if err := tw.Close(); err != nil {
		t.Fatalf("close tar writer: %v", err)
	}
	if err := gz.Close(); err != nil {
		t.Fatalf("close gzip writer: %v", err)
	}

	return buf.Bytes()
}

func setRequestedAssetURL(t *testing.T, config *lifecycle.Config, url string) assetsPkg.Asset {
	t.Helper()
	requestedAssets, ok := config.Extra["cdn_assets_extra"].([]assetsPkg.Asset)
	if !ok || len(requestedAssets) != 1 {
		t.Fatalf("cdn_assets_extra = %#v, want one requested asset", config.Extra["cdn_assets_extra"])
	}
	requestedAssets[0].URL = url
	if requestedAssets[0].Version == webAwesomeDefaultVersion {
		requestedAssets[0].Integrity = computeArchiveIntegrity(buildWebAwesomeTestTarGz(t))
	}
	config.Extra["cdn_assets_extra"] = requestedAssets
	return requestedAssets[0]
}

func setRequestedAsset(t *testing.T, config *lifecycle.Config, asset assetsPkg.Asset) {
	t.Helper()
	requestedAssets, ok := config.Extra["cdn_assets_extra"].([]assetsPkg.Asset)
	if !ok || len(requestedAssets) != 1 {
		t.Fatalf("cdn_assets_extra = %#v, want one requested asset", config.Extra["cdn_assets_extra"])
	}
	requestedAssets[0] = asset
	config.Extra["cdn_assets_extra"] = requestedAssets
}

func computeArchiveIntegrity(data []byte) string {
	h := sha512.Sum512(data)
	return "sha512-" + base64.StdEncoding.EncodeToString(h[:])
}

func assertFileExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file %s: %v", path, err)
	}
}
