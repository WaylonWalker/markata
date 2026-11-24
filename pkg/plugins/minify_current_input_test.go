package plugins

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// An author can intentionally replace static source with a previous result.
// Neither an immediate nor a later option change may resurrect older comments.
func TestMinifyCache_StaticEditToPreviousOutput(t *testing.T) {
	for _, delayed := range []bool{false, true} {
		name := "immediate"
		if delayed {
			name = "delayed"
		}
		t.Run(name, func(t *testing.T) {
			config, cache, path := newMinifyFixture(t, cssMinifyPluginName)
			cache.close()
			config.Extra["theme"] = "no-such-minify-test-theme"
			config.Extra[cssMinifyPluginName] = models.CSSMinifyConfig{Enabled: true}
			source := filepath.Join(config.ContentDir, "static", filepath.Base(path))
			writeMinifyFixture(t, source, []byte("/* License A */\nbody { color: red; }\n"))
			manager := lifecycle.NewManager()
			manager.SetConfig(config)
			static, css := NewStaticAssetsPlugin(), NewCSSMinifyPlugin()
			build := func() {
				t.Helper()
				if err := static.Configure(manager); err != nil {
					t.Fatal(err)
				}
				if err := css.Configure(manager); err != nil {
					t.Fatal(err)
				}
				if err := static.Write(manager); err != nil {
					t.Fatal(err)
				}
				if err := css.Write(manager); err != nil {
					t.Fatal(err)
				}
				if err := static.Cleanup(manager); err != nil {
					t.Fatal(err)
				}
			}
			build()
			edited := readMinifyFixture(t, path)
			if bytes.Contains(edited, []byte("License A")) {
				t.Fatal("initial build did not strip the comment")
			}
			writeMinifyFixture(t, source, edited)
			if delayed {
				build()
				if !bytes.Equal(readMinifyFixture(t, path), edited) {
					t.Fatal("same-recipe build changed the edited source")
				}
			}
			config.Extra[cssMinifyPluginName] = models.CSSMinifyConfig{
				Enabled: true, PreserveComments: []string{"License A"},
			}
			build()
			// The oracle sees the SAME edited static source as the cached build.
			expected, err := css.minifyBytes(readMinifyFixture(t, source))
			if err != nil {
				t.Fatal(err)
			}
			if actual := readMinifyFixture(t, path); !bytes.Equal(actual, expected) || bytes.Contains(actual, []byte("License A")) {
				t.Fatalf("stale comment resurrected: got %q, current-input oracle %q", actual, expected)
			}
		})
	}
}

// Legacy retry records can associate an observed digest with a result derived
// from a different historical snapshot. That digest must never authorize reuse.
func TestMinifyCache_InputHashFailureRetryCollision(t *testing.T) {
	for _, changedRecipe := range []bool{false, true} {
		name := "same-recipe"
		if changedRecipe {
			name = "changed-recipe"
		}
		t.Run(name, func(t *testing.T) {
			_, cache, path := newMinifyFixture(t, cssMinifyPluginName)
			raw := []byte("/* License A */\nbody { color: red; }\n")
			edited, err := NewCSSMinifyPlugin().minifyBytes(raw)
			if err != nil {
				t.Fatal(err)
			}
			css := NewCSSMinifyPlugin()
			css.config.PreserveComments = []string{"License A"}
			recipe := minifyRecipe(css.Name(), css.config.PreserveComments)
			writeMinifyFixture(t, path, raw)
			failWrite := func(data []byte) ([]byte, error) {
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
				if err := os.Mkdir(path, 0o700); err != nil {
					t.Fatal(err)
				}
				return css.minifyBytes(data)
			}
			if _, err := minifyAsset(path, recipe, cache, failWrite); err == nil {
				t.Fatal("expected target replacement failure")
			}
			record, err := cache.readRecord(filepath.Base(path))
			if err != nil || record == nil {
				t.Fatalf("failed write did not retain a retryable record: %v", err)
			}
			// Emulate the earlier snapshot-substitution record format.
			data, err := json.Marshal(struct {
				*minifyRecord
				InputHash string `json:"input_hash"`
			}{record, minifyHash(edited)})
			if err != nil {
				t.Fatal(err)
			}
			if err := cache.writePrivate(minifyRecordPath(record.Asset), data); err != nil {
				t.Fatal(err)
			}
			if err := os.Remove(path); err != nil {
				t.Fatal(err)
			}
			writeMinifyFixture(t, path, edited)
			if changedRecipe {
				recipe = minifyHash([]byte("next-recipe"))
			}
			calls := 0
			transform := func(current []byte) ([]byte, error) {
				calls++
				if !bytes.Equal(current, edited) {
					t.Fatalf("retry selected historical source: %q", current)
				}
				return css.minifyBytes(current)
			}
			requireMinifyStatus(t, path, recipe, cache, transform, "transformed")
			if calls != 1 || !bytes.Equal(readMinifyFixture(t, path), edited) {
				t.Fatal("retry digest restored an obsolete comment")
			}
			updated, err := cache.readRecord(record.Asset)
			if err != nil || updated == nil || updated.SourceHash != minifyHash(edited) {
				t.Fatalf("new relation does not describe exact current input: %#v, %v", updated, err)
			}
			if bytes.Contains(readMinifyFixture(t, filepath.Join(cache.dir, minifyRecordPath(record.Asset))), []byte(`"input_hash"`)) {
				t.Fatal("new record produced a legacy retry digest")
			}
		})
	}
}

func TestMinifyCache_ExactInputRestoresAtomically(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix file permissions required")
	}
	for _, kind := range []string{cssMinifyPluginName, jsMinifyPluginName} {
		t.Run(kind, func(t *testing.T) {
			config, cache, path := newMinifyFixture(t, kind)
			raw := []byte("raw")
			result := []byte("raw")
			recipe := minifyRecipe(kind, nil)
			writeMinifyFixture(t, path, raw)
			requireMinifyStatus(t, path, recipe, cache,
				func(data []byte) ([]byte, error) { return bytes.TrimSpace(data), nil }, "transformed")
			// Replace with fresh identical exact input, not just a chmod on the
			// previously published target. Link it to an existing release.
			if err := os.Remove(path); err != nil {
				t.Fatal(err)
			}
			writeMinifyFixture(t, path, result)
			before := statMinifyFixture(t, path)
			release := filepath.Join(config.ContentDir, "release")
			if err := os.Link(path, release); err != nil {
				t.Skipf("hardlinks unavailable: %v", err)
			}
			noTransform := func([]byte) ([]byte, error) {
				t.Fatal("exact-input hit invoked the engine")
				return nil, nil
			}
			requireMinifyStatus(t, path, recipe, cache, noTransform, "restored")
			after := statMinifyFixture(t, path)
			old := statMinifyFixture(t, release)
			if os.SameFile(before, after) || after.Mode().Perm() != 0o644 ||
				!os.SameFile(before, old) || old.Mode().Perm() != 0o600 {
				t.Fatal("exact-input restore was not hard-link-safe atomic replacement")
			}
			if !bytes.Equal(readMinifyFixture(t, path), result) || !bytes.Equal(readMinifyFixture(t, release), result) {
				t.Fatal("normalization changed output or release bytes")
			}
			requireMinifyStatus(t, path, recipe, cache, noTransform, "restored")
			warm := statMinifyFixture(t, path)
			if os.SameFile(after, warm) || warm.Mode().Perm() != 0o644 {
				t.Fatal("exact-input hit did not atomically republish normally permissioned target")
			}
		})
	}
}

func TestMinifyCache_NonIdempotentImportantComments(t *testing.T) {
	_, cache, path := newMinifyFixture(t, cssMinifyPluginName)
	css := NewCSSMinifyPlugin()
	css.config.PreserveComments = []string{"Copyright"}
	recipe := minifyRecipe(css.Name(), css.config.PreserveComments)
	raw := []byte("/*! Copyright author */\nbody { color: red; }\n")
	writeMinifyFixture(t, path, raw)
	calls := 0
	transform := func(current []byte) ([]byte, error) {
		calls++
		return css.minifyBytes(current)
	}
	// Retain each output as the next invocation's exact input. The wrapper
	// prepends the important comment and the engine also preserves it, so
	// skipping based on a previous output digest would diverge from main.
	for i := 0; i < 4; i++ {
		current := readMinifyFixture(t, path)
		expected, err := css.minifyBytes(current)
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Equal(current, expected) {
			t.Fatal("fixture must demonstrate a non-idempotent transform")
		}
		requireMinifyStatus(t, path, recipe, cache, transform, "transformed")
		if actual := readMinifyFixture(t, path); !bytes.Equal(actual, expected) {
			t.Fatalf("invocation %d: got %q, direct current-input transform %q", i, actual, expected)
		}
		record, err := cache.readRecord(filepath.Base(path))
		if err != nil || record == nil || record.SourceHash != minifyHash(current) {
			t.Fatalf("record does not identify exact invocation input: %#v, %v", record, err)
		}
	}
	if calls != 4 {
		t.Fatalf("retained non-idempotent output needs 4 engine calls, got %d", calls)
	}
	expected, err := css.minifyBytes(raw)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 4; i++ {
		writeMinifyFixture(t, path, raw)
		status := "restored"
		if i == 0 {
			status = "transformed" // Prime the exact raw-input relation again.
		}
		requireMinifyStatus(t, path, recipe, cache, transform, status)
		if !bytes.Equal(readMinifyFixture(t, path), expected) {
			t.Fatal("raw recopy diverged from direct transform")
		}
	}
	if calls != 5 {
		t.Fatalf("primed repeated raw copies invoked the engine: total calls %d, want 5", calls)
	}
	if _, err := os.Stat(filepath.Join(cache.dir, "blobs", minifyHash(raw))); !os.IsNotExist(err) {
		t.Fatalf("raw source snapshot was produced: %v", err)
	}
}
