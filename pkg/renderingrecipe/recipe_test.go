package renderingrecipe

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/WaylonWalker/markata/pkg/renderingcontract"
)

func TestCompile_CanonicalBundleHasRealHashes(t *testing.T) {
	theme, err := LoadCanonicalTheme()
	if err != nil {
		t.Fatal(err)
	}
	one, err := Compile(theme)
	if err != nil {
		t.Fatal(err)
	}
	two, err := Compile(theme)
	if err != nil {
		t.Fatal(err)
	}
	if one.Manifest.SemanticHash == "" || one.Manifest.RecipeHash == "" {
		t.Fatal("missing hashes")
	}
	if one.Manifest.RecipeHash != two.Manifest.RecipeHash {
		t.Fatal("recipe hash is not deterministic")
	}
	for _, asset := range one.Manifest.Assets {
		if len(asset.SHA256) != 64 || !bytes.Equal(one.Assets[asset.Path], two.Assets[asset.Path]) {
			t.Fatalf("asset %s is not deterministic", asset.Path)
		}
	}
}

func TestCanonicalJSON_RejectsDuplicateKeys(t *testing.T) {
	var value map[string]any
	if err := decodeStrict([]byte(`{"a":1,"a":2}`), &value); err == nil {
		t.Fatal("duplicate key accepted")
	}
}

func TestNormalize_RejectsMalformedTypes(t *testing.T) {
	if _, err := normalize(map[string]any{"theme": map[string]any{"palette": nil}}); err == nil {
		t.Fatal("null palette accepted")
	}
	if _, err := normalize(map[string]any{"theme": map[string]any{"motif": map[string]any{"size": 71}}}); err == nil {
		t.Fatal("numeric motif size accepted")
	}
	if _, err := normalize(map[string]any{"theme": map[string]any{"motif": map[string]any{"size": "39px"}}}); err == nil {
		t.Fatal("undersized motif accepted")
	}
	if _, err := normalize(map[string]any{"theme": map[string]any{"motif": map[string]any{"gap": "33px"}}}); err == nil {
		t.Fatal("oversized gap accepted")
	}
}

func TestMixColor_InterpolatesInSRGB(t *testing.T) {
	if got := mixColor("#000000", "#ffffff", .5); got != "#808080" {
		t.Fatalf("got %s", got)
	}
}

func TestRecipeSources_MatchCanonicalSpec(t *testing.T) {
	for _, id := range []string{"texture-screenprint-v1", "heading-splatter-v1", "motif-block-w-v1"} {
		canonical, err := os.ReadFile(filepath.Join("..", "..", "spec", "rendering-contract", "recipes", "source", id+".json"))
		if err != nil {
			t.Fatal(err)
		}
		embedded, err := testdata.ReadFile("testdata/recipe-sources/" + id + ".json")
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(canonical, embedded) {
			t.Fatalf("embedded source %s is stale", id)
		}
	}
}

func TestRecipeSources_SeedAffectsGeometry(t *testing.T) {
	sources, err := loadRecipeSources()
	if err != nil {
		t.Fatal(err)
	}
	m := renderingcontract.MotifState{ColorMix: 1, Color: "ink", Size: "78px", Gap: "10px"}
	a := motifSVG(m, "#000000", "#ffffff", sources["motif-block-w-v1"])
	source := sources["motif-block-w-v1"]
	source.Seed++
	b := motifSVG(m, "#000000", "#ffffff", source)
	if bytes.Equal(a, b) {
		t.Fatal("source seed does not affect motif geometry")
	}
}

func TestCompile_HeadingPassCarriesResolvedPaint(t *testing.T) {
	theme, err := LoadCanonicalTheme()
	if err != nil {
		t.Fatal(err)
	}
	bundle, err := Compile(theme)
	if err != nil {
		t.Fatal(err)
	}
	for _, pass := range bundle.Manifest.Passes {
		if pass.ID == "heading-wear" && pass.Paint == "" {
			t.Fatal("heading paint missing")
		}
	}
}

func TestHeadingPaint_TracksMixEndpoints(t *testing.T) {
	if got := mixColor("#101010", "#e0e0e0", 0); got != "#101010" {
		t.Fatalf("zero mix = %s", got)
	}
	if got := mixColor("#101010", "#e0e0e0", 1); got != "#e0e0e0" {
		t.Fatalf("full mix = %s", got)
	}
}

func TestHeadingGeometry_IndependentOfMix(t *testing.T) {
	sources, err := loadRecipeSources()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(headingSVG(0, sources["heading-splatter-v1"]), headingSVG(1, sources["heading-splatter-v1"])) {
		t.Fatal("heading geometry varies with mix")
	}
}

func TestHeadingSVG_HasMostlyOpaqueFieldAndRealWearAlpha(t *testing.T) {
	sources, err := loadRecipeSources()
	if err != nil {
		t.Fatal(err)
	}
	data := string(headingSVG(0, sources["heading-splatter-v1"]))
	if !bytes.Contains([]byte(data), []byte(`fill-rule="evenodd"`)) || !bytes.Contains([]byte(data), []byte(`opacity=".42"`)) {
		t.Fatalf("heading asset does not contain direct transparent and partial-alpha wear: %s", data)
	}
	if bytes.Contains([]byte(data), []byte(`<mask`)) {
		t.Fatal("heading asset contains nested SVG mask")
	}
}

func TestMotifFieldDimensions_DeriveFromState(t *testing.T) {
	sources, err := loadRecipeSources()
	if err != nil {
		t.Fatal(err)
	}
	canonical := renderingcontract.MotifState{Size: "71px", Gap: "11px"}
	wantW, wantH, err := motifFieldDimensions(canonical, sources["motif-block-w-v1"])
	if err != nil || wantW != 1328 || wantH != 507 {
		t.Fatalf("canonical field = %dx%d (%v)", wantW, wantH, err)
	}
	changed := renderingcontract.MotifState{Size: "60px", Gap: "18px"}
	gotW, gotH, err := motifFieldDimensions(changed, sources["motif-block-w-v1"])
	if err != nil || gotW == wantW || gotH == wantH {
		t.Fatalf("field did not respond to size/gap: %dx%d (%v)", gotW, gotH, err)
	}
}

func TestMotifDimensions_AllowsZeroGap(t *testing.T) {
	size, gap, err := motifDimensions(renderingcontract.MotifState{Size: "40px", Gap: "0px"})
	if err != nil || size != 40 || gap != 0 {
		t.Fatalf("size=%v gap=%v err=%v", size, gap, err)
	}
}

func TestLoadCanonicalTheme_RejectsObsoleteFields(t *testing.T) {
	if _, err := normalize(map[string]any{"theme": map[string]any{"motif": map[string]any{"scale": 1}}}); err == nil {
		t.Fatal("obsolete motif.scale accepted")
	}
}

func TestCanonicalJSON_PreservesZeroAndUnicode(t *testing.T) {
	a, err := canonicalJSON(map[string]any{"zero": 0, "text": "café"})
	if err != nil {
		t.Fatal(err)
	}
	b, err := canonicalJSON(map[string]any{"text": "café", "zero": 0})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(a, b) || !bytes.Contains(a, []byte(`"zero":0`)) || !bytes.Contains(a, []byte("café")) {
		t.Fatalf("canonical output = %s", a)
	}
}

func TestHashVectors_AssetAndFontChangesChangeRecipe(t *testing.T) {
	theme, err := LoadCanonicalTheme()
	if err != nil {
		t.Fatal(err)
	}
	one, err := Compile(theme)
	if err != nil {
		t.Fatal(err)
	}
	if AssetHash([]byte("a")) == AssetHash([]byte("b")) {
		t.Fatal("asset bytes share hash")
	}
	identity := one.Manifest
	identity.RecipeHash = ""
	identity.Fonts[0].SHA256 = "0" + identity.Fonts[0].SHA256[1:]
	encoded, err := canonicalJSON(identity)
	if err != nil {
		t.Fatal(err)
	}
	if digest([]byte(recipeDomain), encoded) == one.Manifest.RecipeHash {
		t.Fatal("font digest did not change recipe identity")
	}
}

func TestHashVectors_CanonicalFixture(t *testing.T) {
	vectorPath := filepath.Join("..", "..", "spec", "rendering-contract", "recipes", "hash-vectors-v1.json")
	data, err := os.ReadFile(vectorPath)
	if err != nil {
		t.Fatal(err)
	}
	var vectors struct {
		CanonicalFixture struct {
			SemanticHash string `json:"semantic_hash"`
			RecipeHash   string `json:"recipe_hash"`
			AssetXSHA256 string `json:"asset_x_sha256"`
		} `json:"canonical_fixture"`
		Cases []json.RawMessage `json:"cases"`
	}
	if err := json.Unmarshal(data, &vectors); err != nil {
		t.Fatal(err)
	}
	theme, err := LoadCanonicalTheme()
	if err != nil {
		t.Fatal(err)
	}
	semantic, err := SemanticHash(theme)
	if err != nil {
		t.Fatal(err)
	}
	bundle, err := Compile(theme)
	if err != nil {
		t.Fatal(err)
	}
	if semantic != vectors.CanonicalFixture.SemanticHash {
		t.Fatalf("semantic hash = %s", semantic)
	}
	if bundle.Manifest.RecipeHash != vectors.CanonicalFixture.RecipeHash {
		t.Fatalf("recipe hash = %s", bundle.Manifest.RecipeHash)
	}
	if AssetHash([]byte(`x`)) != vectors.CanonicalFixture.AssetXSHA256 {
		t.Fatalf("asset vector changed")
	}
	if len(vectors.Cases) != 11 {
		t.Fatalf("cases = %d, want 11", len(vectors.Cases))
	}
	for _, raw := range vectors.Cases {
		var c struct {
			ID                string  `json:"id"`
			Value             any     `json:"value"`
			ExpectedCanonical string  `json:"expected_canonical"`
			Background        string  `json:"background"`
			Foreground        string  `json:"foreground"`
			Expected          string  `json:"expected"`
			Mix               float64 `json:"mix"`
			A                 string  `json:"a"`
			B                 string  `json:"b"`
			ASHA              string  `json:"a_sha256"`
			BSHA              string  `json:"b_sha256"`
			Size              string  `json:"size"`
			Gap               string  `json:"gap"`
			ExpectedSize      int     `json:"expected_size_milli"`
			ExpectedGap       int     `json:"expected_gap_milli"`
			ExpectedDifferent bool    `json:"expected_different"`
		}
		if err := json.Unmarshal(raw, &c); err != nil {
			t.Fatal(err)
		}
		if c.ID == "" {
			t.Fatal("vector case missing id")
		}
		t.Run(c.ID, func(t *testing.T) {
			switch c.ID {
			case "explicit-zero":
				if got, err := json.Marshal(c.Value); err != nil || string(got) != c.ExpectedCanonical {
					t.Fatalf("canonical = %s err=%v", got, err)
				}
			case "unicode":
				if got, err := json.Marshal(c.Value); err != nil || string(got) != c.ExpectedCanonical {
					t.Fatalf("canonical = %s err=%v", got, err)
				}
			case "object-key-order", "dial-precision-bounds":
				got, err := canonicalJSON(c.Value)
				if err != nil {
					t.Fatal(err)
				}
				if string(got) != c.ExpectedCanonical {
					t.Fatalf("canonical = %s", got)
				}
			case "asset-byte-change":
				if AssetHash([]byte(c.A)) != c.ASHA || AssetHash([]byte(c.B)) != c.BSHA || c.ASHA == c.BSHA {
					t.Fatal("asset vector mismatch")
				}
			case "color-mix-0", "color-mix-intermediate", "color-mix-1":
				if got := mixColor(c.Background, c.Foreground, c.Mix); got != c.Expected {
					t.Fatalf("mix = %s", got)
				}
			case "motif-size-gap":
				size, gap, err := motifDimensions(renderingcontract.MotifState{Size: c.Size, Gap: c.Gap})
				if err != nil || int(size*1000) != c.ExpectedSize || int(gap*1000) != c.ExpectedGap {
					t.Fatalf("size=%v gap=%v err=%v", size, gap, err)
				}
			case "primitive-source-change":
				sources, err := loadRecipeSources()
				if err != nil {
					t.Fatal(err)
				}
				source := sources["motif-block-w-v1"]
				a := motifSVG(renderingcontract.MotifState{Size: "40px", Gap: "0px"}, "#000", "#fff", source)
				source.Seed++
				b := motifSVG(renderingcontract.MotifState{Size: "40px", Gap: "0px"}, "#000", "#fff", source)
				if bytes.Equal(a, b) != !c.ExpectedDifferent {
					t.Fatal("source mutation identity mismatch")
				}
			case "font-digest-change":
				identity := bundle.Manifest
				identity.RecipeHash = ""
				identity.Fonts[0].SHA256 = "0" + identity.Fonts[0].SHA256[1:]
				encoded, err := canonicalJSON(identity)
				if err != nil {
					t.Fatal(err)
				}
				if (digest([]byte(recipeDomain), encoded) == bundle.Manifest.RecipeHash) != !c.ExpectedDifferent {
					t.Fatal("font mutation identity mismatch")
				}
			default:
				t.Fatalf("unknown hash vector case %q", c.ID)
			}
		})
	}
}
