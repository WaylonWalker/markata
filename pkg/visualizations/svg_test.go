package visualizations

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewSVGAssetDimensionsAndHash(t *testing.T) {
	a, err := NewSVGAsset([]byte(sampleSVG))
	if err != nil {
		t.Fatal(err)
	}
	if a.Width != 960 || a.Height != 540 || !strings.HasSuffix(a.Path, a.SHA256+".svg") {
		t.Fatalf("asset: %+v", a)
	}
	a.SVG[0] = 'X'
	b, err := NewSVGAsset([]byte(sampleSVG))
	if err != nil || b.SVG[0] != '<' {
		t.Fatal("asset bytes not owned by caller")
	}
	scaled, err := NewSVGAsset([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="480px" viewBox="0 0 960 540"><path d="M0 0"/></svg>`))
	if err != nil || scaled.Width != 480 || scaled.Height != 270 {
		t.Fatalf("derived size: %+v %v", scaled, err)
	}
}

func TestInvalidOrExternalSVG(t *testing.T) {
	invalid := []string{
		`<svg xmlns="http://www.w3.org/2000/svg"><rect/></svg>`,
		`<svg viewBox="0 0 -1 2"/>`,
		`<svg xmlns="http://example.com/fake" viewBox="0 0 1 2"/>`,
		`<svg width="40" height="40" viewBox="0 0 -1 2"/>`,
		`<svg width="40" height="40" viewBox="NaN 0 10 10"/>`,
		`<svg width="40" height="40" viewBox="0 0 Infinity 10"/>`,
		`<svg width="50%" height="40" viewBox="0 0 10 20"/>`,
		`<!DOCTYPE svg [<!ENTITY x SYSTEM "file:///etc/passwd">]><svg viewBox="0 0 1 1"/>`,
		`<svg viewBox="0 0 1 1"><script>alert(1)</script></svg>`,
		`<?xml-stylesheet href="https://example.com/evil.css"?><svg viewBox="0 0 1 1"/>`,
		`<svg viewBox="0 0 1 1"><foreignObject><p>oops</p></foreignObject></svg>`,
		`<svg viewBox="0 0 1 1"><rect onclick="alert(1)"/></svg>`,
		`<svg viewBox="0 0 1 1"><use href="https://example.com/a.svg#symbol"/></svg>`,
		`<svg viewBox="0 0 1 1"><style>@import url(https://example.com/x.css)</style></svg>`,
		`<svg viewBox="0 0 1 1"><rect style="fill:url(https://example.com/x)"/></svg>`,
		`<svg viewBox="0 0 1 1"><image href="file:///etc/passwd"/></svg>`,
		`<svg viewBox="0 0 1 1"/>trailing text`,
	}
	for _, source := range invalid {
		if _, err := NewSVGAsset([]byte(source)); err == nil {
			t.Errorf("unexpectedly accepted %q", source)
		}
	}
	if _, err := NewSVGAsset(bytes.Repeat([]byte("x"), MaxSVGBytes+1)); err == nil {
		t.Fatal("oversize SVG accepted")
	}
	valid := `<svg viewBox="0 0 10 5"><style>.a{font-family:system-ui;fill:url(#paint)}</style><defs><linearGradient id="paint"/></defs><rect fill="url(#paint)"/><text>url(https://example.com) is just a label</text></svg>`
	if _, err := NewSVGAsset([]byte(valid)); err != nil {
		t.Fatalf("internal CSS and ordinary labels rejected: %v", err)
	}
}

func TestPublishAndFigure(t *testing.T) {
	asset, err := NewSVGAsset([]byte(sampleSVG))
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	path, err := Publish(root, asset)
	if err != nil {
		t.Fatal(err)
	}
	if path != asset.Path {
		t.Fatalf("path=%s", path)
	}
	stored, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
	if err != nil || !bytes.Equal(stored, asset.SVG) {
		t.Fatal("published data mismatch", err)
	}
	if _, err := Publish(root, asset); err != nil {
		t.Fatal("identical publish must be idempotent", err)
	}
	asset.Path = "assets/visualizations/forged.svg"
	if _, err := Publish(root, asset); err == nil {
		t.Fatal("forged asset path accepted")
	}
	asset.Path = path
	asset.Width = 999
	if _, err := FigureHTML(asset, "description", "", false); err == nil {
		t.Fatal("forged dimensions accepted")
	}
	asset.Width = 960
	markup, err := FigureHTML(asset, `A <B> & C`, `A "caption"`, true)
	if err != nil || !strings.Contains(markup, `width="960" height="540" loading="eager"`) || !strings.Contains(markup, `A &lt;B&gt; &amp; C`) || !strings.Contains(markup, `A &#34;caption&#34;`) {
		t.Fatalf("markup=%s error=%v", markup, err)
	}
	if _, err := FigureHTML(asset, "", "", false); err == nil {
		t.Fatal("missing alt accepted")
	}
}
