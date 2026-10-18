package plugins

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"maps"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func writeMinifyFixture(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
}

func readMinifyFixture(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// statMinifyFixture captures identity from an open handle: os.Stat on Windows
// defers loading file IDs until SameFile, which may reopen an already replaced
// path. Close the handle before returning so it cannot block replacement.
func statMinifyFixture(t *testing.T, path string) os.FileInfo {
	t.Helper()
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	info, statErr := file.Stat()
	closeErr := file.Close()
	if statErr != nil {
		t.Fatal(statErr)
	}
	if closeErr != nil {
		t.Fatal(closeErr)
	}
	return info
}

func newMinifyFixture(t *testing.T, kind string) (*lifecycle.Config, *minifyCache, string) {
	t.Helper()
	site := t.TempDir()
	config := lifecycle.NewConfig()
	config.ContentDir = site
	config.OutputDir = filepath.Join(site, "output")
	if err := os.MkdirAll(config.OutputDir, 0o700); err != nil {
		t.Fatal(err)
	}
	cache, err := newMinifyCache(config, kind)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(cache.close)
	return config, cache, filepath.Join(config.OutputDir, "site."+strings.TrimSuffix(kind, "_minify"))
}

func captureMinifyLog(t *testing.T) *bytes.Buffer {
	t.Helper()
	var output bytes.Buffer
	old := log.Writer()
	log.SetOutput(&output)
	t.Cleanup(func() { log.SetOutput(old) })
	return &output
}

func requireMinifyStatus(t *testing.T, path, recipe string, cache *minifyCache, transform minifyFunc, status string) {
	t.Helper()
	result, err := minifyAsset(path, recipe, cache, transform)
	if err != nil {
		t.Fatal(err)
	}
	if result.status != status {
		t.Fatalf("status = %q, want %q", result.status, status)
	}
}

func TestMinifyCache_ExactReuse(t *testing.T) {
	for _, kind := range []string{"js_minify", "css_minify"} {
		t.Run(kind, func(t *testing.T) {
			config, cache, path := newMinifyFixture(t, kind)
			raw := []byte("function hello ( ) { return  42; }\n")
			transform := NewJSMinifyPlugin().minifyBytes
			uncached := NewJSMinifyPlugin().minifyFile
			if kind == "css_minify" {
				raw = []byte("/*! Copyright */\nbody { color: #ffffff; margin: 0px; }\n")
				css := NewCSSMinifyPlugin()
				css.config.PreserveComments = []string{"Copyright"}
				transform = css.minifyBytes
				uncached = css.minifyFile
			}
			expected, err := transform(raw)
			if err != nil {
				t.Fatal(err)
			}
			referencePath := filepath.Join(config.ContentDir, "reference"+filepath.Ext(path))
			writeMinifyFixture(t, referencePath, raw)
			if _, _, err := uncached(referencePath); err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(readMinifyFixture(t, referencePath), expected) {
				t.Fatal("byte transform differs from uncached minifyFile wrapper")
			}
			calls := 0
			counted := func(data []byte) ([]byte, error) {
				calls++
				return transform(data)
			}
			recipe := minifyRecipe(kind, nil)
			writeMinifyFixture(t, path, raw)
			requireMinifyStatus(t, path, recipe, cache, counted, "transformed")
			if !bytes.Equal(readMinifyFixture(t, path), expected) {
				t.Fatal("cold output differs from uncached transform")
			}
			recordPath := filepath.Join(cache.dir, minifyRecordPath(filepath.Base(path)))
			recordBefore := statMinifyFixture(t, recordPath)
			// Reopen to prove reuse survives process/plugin lifetimes.
			cache.close()
			reopened, err := newMinifyCache(config, kind)
			if err != nil {
				t.Fatal(err)
			}
			defer reopened.close()
			for i := 0; i < 3; i++ {
				writeMinifyFixture(t, path, raw)
				targetBefore := statMinifyFixture(t, path)
				requireMinifyStatus(t, path, recipe, reopened, counted, "restored")
				targetAfter := statMinifyFixture(t, path)
				if os.SameFile(targetBefore, targetAfter) {
					t.Fatal("exact-input restore failed to replace target")
				}
				if !bytes.Equal(readMinifyFixture(t, path), expected) {
					t.Fatal("raw recopy did not restore exact bytes")
				}
			}
			recordAfter := statMinifyFixture(t, recordPath)
			if !os.SameFile(recordBefore, recordAfter) ||
				!recordBefore.ModTime().Equal(recordAfter.ModTime()) {
				t.Fatal("exact-input restore rewrote the record")
			}
			if calls != 1 {
				t.Fatalf("minifier calls = %d, want 1 (zero for repeated raw copies)", calls)
			}
		})
	}
}

func TestMinifyCache_ContentNotStatAndStableIdentity(t *testing.T) {
	_, cache, path := newMinifyFixture(t, "js_minify")
	deep := filepath.Join(cache.output, "assets", "vendor", "deep", "site.js")
	files := []string{path, deep}
	calls := 0
	transform := func(data []byte) ([]byte, error) {
		calls++
		return bytes.TrimSpace(data), nil
	}
	recipe := minifyRecipe("js_minify", nil)
	for _, file := range files {
		writeMinifyFixture(t, file, []byte(" alpha "))
		requireMinifyStatus(t, file, recipe, cache, transform, "transformed")
	}
	before, err := os.Stat(deep)
	if err != nil {
		t.Fatal(err)
	}
	writeMinifyFixture(t, deep, []byte("bravo"))
	if err := os.Chtimes(deep, before.ModTime(), before.ModTime()); err != nil {
		t.Fatal(err)
	}
	requireMinifyStatus(t, deep, recipe, cache, transform, "transformed")
	requireMinifyStatus(t, path, recipe, cache, transform, "transformed")
	if calls != 4 || string(readMinifyFixture(t, deep)) != "bravo" {
		t.Fatal("same-size/mtime edit missed or current input was not transformed")
	}
	rootRecord, err := cache.readRecord("site.js")
	if err != nil {
		t.Fatal(err)
	}
	deepRecord, err := cache.readRecord("assets/vendor/deep/site.js")
	if err != nil {
		t.Fatal(err)
	}
	if rootRecord == nil || deepRecord == nil {
		t.Fatal("root/deep assets did not retain output-relative identities")
	}
}

func TestMinifyCache_RecipeChanges(t *testing.T) {
	for _, recopy := range []bool{false, true} {
		t.Run(map[bool]string{false: "processed-only", true: "raw-recopy"}[recopy], func(t *testing.T) {
			_, cache, path := newMinifyFixture(t, "css_minify")
			raw := []byte("/* License A */\n/* License B */\nbody { color: red; }\n")
			p := NewCSSMinifyPlugin()
			p.config.PreserveComments = []string{"License A"}
			writeMinifyFixture(t, path, raw)
			requireMinifyStatus(t, path, minifyRecipe(p.Name(), p.config.PreserveComments), cache, p.minifyBytes, "transformed")
			if recopy {
				writeMinifyFixture(t, path, raw)
			}
			current := readMinifyFixture(t, path)
			p.config.PreserveComments = []string{"License B"}
			expected, err := p.minifyBytes(current)
			if err != nil {
				t.Fatal(err)
			}
			requireMinifyStatus(t, path, minifyRecipe(p.Name(), p.config.PreserveComments), cache, p.minifyBytes, "transformed")
			if !bytes.Equal(readMinifyFixture(t, path), expected) {
				t.Fatal("option change did not transform current stage input")
			}
			seen := []byte(nil)
			revised := func(data []byte) ([]byte, error) {
				seen = bytes.Clone(data)
				return p.minifyBytes(data)
			}
			current = readMinifyFixture(t, path)
			requireMinifyStatus(t, path, minifyHash([]byte("next-wrapper-revision")), cache, revised, "transformed")
			if !bytes.Equal(seen, current) {
				t.Fatal("wrapper revision change did not receive current stage input")
			}
		})
	}
	if minifyRecipe("css_minify", []string{"B", "A", "A"}) != minifyRecipe("css_minify", []string{"A", "B"}) {
		t.Fatal("equivalent comment predicates have different recipes")
	}
	if minifyRecipe("css_minify", nil) == minifyRecipe("js_minify", nil) {
		t.Fatal("kind absent from recipe")
	}
}

func TestMinifyCache_MissingAndCorruptPayloads(t *testing.T) {
	for _, payload := range []string{"result-missing", "result-corrupt", "source-missing", "source-corrupt", "source-unavailable"} {
		t.Run(payload, func(t *testing.T) {
			logs := captureMinifyLog(t)
			_, cache, path := newMinifyFixture(t, "js_minify")
			raw := []byte(" raw ")
			recipe := minifyRecipe("js_minify", nil)
			calls := 0
			transform := func(data []byte) ([]byte, error) {
				calls++
				return bytes.TrimSpace(data), nil
			}
			writeMinifyFixture(t, path, raw)
			requireMinifyStatus(t, path, recipe, cache, transform, "transformed")
			record, err := cache.readRecord("site.js")
			if err != nil {
				t.Fatal(err)
			}
			hash := record.SourceHash
			if strings.HasPrefix(payload, "result") {
				hash = record.OutputHash
			}
			blobPath := filepath.Join(cache.dir, "blobs", hash)
			if strings.HasPrefix(payload, "source") {
				// Historical snapshots are ignored; new transforms store results only.
				writeMinifyFixture(t, blobPath, raw)
			}
			switch {
			case strings.HasSuffix(payload, "missing"):
				if err := os.Remove(blobPath); err != nil {
					t.Fatal(err)
				}
			case strings.HasSuffix(payload, "corrupt"):
				writeMinifyFixture(t, blobPath, []byte("corrupt"))
			case payload == "source-unavailable":
				record.SourceHash = ""
				data, err := json.Marshal(record)
				if err != nil {
					t.Fatal(err)
				}
				writeMinifyFixture(t, filepath.Join(cache.dir, minifyRecordPath("site.js")), data)
			}
			if strings.HasPrefix(payload, "result") {
				writeMinifyFixture(t, path, raw)
			} else {
				recipe = minifyHash([]byte("new recipe"))
			}
			current := readMinifyFixture(t, path)
			requireMinifyStatus(t, path, recipe, cache, transform, "transformed")
			if calls != 2 {
				t.Fatalf("calls/logs = %d / %s", calls, logs.String())
			}
			repaired, err := cache.readRecord("site.js")
			if err != nil {
				t.Fatal(err)
			}
			if strings.HasPrefix(payload, "source") {
				if repaired.SourceHash != minifyHash(current) || strings.Contains(logs.String(), "degraded provenance") {
					t.Fatal("historical snapshot affected the current-input relation")
				}
			} else {
				if !strings.Contains(logs.String(), "Warning") {
					t.Fatal("unavailable result was not warned about")
				}
				if _, err := cache.readBlob(repaired.OutputHash); err != nil {
					t.Fatal("result payload not repaired")
				}
			}
		})
	}
}

func TestMinifyCache_InvalidRecordsAndLegacy(t *testing.T) {
	for _, invalid := range []string{"null", "{", `{"site.js":"legacyhash"}`, `{"version":1,"output_hash":"../../outside"}`} {
		t.Run(invalid, func(t *testing.T) {
			logs := captureMinifyLog(t)
			_, cache, path := newMinifyFixture(t, "js_minify")
			writeMinifyFixture(t, path, []byte(" raw "))
			recordPath := filepath.Join(cache.dir, minifyRecordPath("site.js"))
			writeMinifyFixture(t, recordPath, []byte(invalid))
			requireMinifyStatus(t, path, minifyRecipe("js_minify", nil), cache,
				func(data []byte) ([]byte, error) { return bytes.TrimSpace(data), nil }, "transformed")
			if !strings.Contains(logs.String(), "invalid asset record") {
				t.Fatal("invalid record was not warned about")
			}
			record, err := cache.readRecord("site.js")
			if err != nil || record == nil {
				t.Fatalf("record not repaired: %v", err)
			}
		})
	}
	_, cache, path := newMinifyFixture(t, "js_minify")
	writeMinifyFixture(t, path, []byte(" raw "))
	legacy := filepath.Join(cache.output, "assets", "vendor", ".markata-js_minify-cache")
	legacyBytes := []byte(`{"site.js":"obsolete"}`)
	writeMinifyFixture(t, legacy, legacyBytes)
	requireMinifyStatus(t, path, minifyRecipe("js_minify", nil), cache,
		func(data []byte) ([]byte, error) { return bytes.TrimSpace(data), nil }, "transformed")
	if !bytes.Equal(readMinifyFixture(t, legacy), legacyBytes) {
		t.Fatal("legacy sidecar was traversed or changed")
	}
}

func TestMinifyCache_TargetFailureRetries(t *testing.T) {
	_, cache, path := newMinifyFixture(t, "js_minify")
	raw := []byte(" raw ")
	writeMinifyFixture(t, path, raw)
	recipe := minifyRecipe("js_minify", nil)
	calls := 0
	failWrite := func(data []byte) ([]byte, error) {
		calls++
		if err := os.Remove(path); err != nil {
			t.Fatal(err)
		}
		if err := os.Mkdir(path, 0o700); err != nil {
			t.Fatal(err)
		}
		return bytes.TrimSpace(data), nil
	}
	if _, err := minifyAsset(path, recipe, cache, failWrite); err == nil {
		t.Fatal("expected target rename failure")
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	writeMinifyFixture(t, path, raw)
	requireMinifyStatus(t, path, recipe, cache, failWrite, "restored")
	if calls != 1 || string(readMinifyFixture(t, path)) != "raw" {
		t.Fatal("write retry transformed source instead of restoring verified result")
	}
}

func TestMinifyCache_RecipeRefreshTargetFailureRetainsCurrentInput(t *testing.T) {
	for _, retry := range []string{"restore", "result-missing", "another-recipe"} {
		t.Run(retry, func(t *testing.T) {
			_, cache, path := newMinifyFixture(t, "css_minify")
			raw := []byte("/* License A */\n/* License B */\nbody { color: red; }\n")
			p := NewCSSMinifyPlugin()
			p.config.PreserveComments = []string{"License A"}
			writeMinifyFixture(t, path, raw)
			requireMinifyStatus(t, path, minifyRecipe(p.Name(), p.config.PreserveComments), cache, p.minifyBytes, "transformed")
			prior := readMinifyFixture(t, path)
			p.config.PreserveComments = []string{"License B"}
			recipe := minifyRecipe(p.Name(), p.config.PreserveComments)
			failWrite := func(data []byte) ([]byte, error) {
				if !bytes.Equal(data, prior) {
					t.Fatal("recipe refresh did not receive current input")
				}
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
				if err := os.Mkdir(path, 0o700); err != nil {
					t.Fatal(err)
				}
				return p.minifyBytes(data)
			}
			if _, err := minifyAsset(path, recipe, cache, failWrite); err == nil {
				t.Fatal("expected target failure")
			}
			if err := os.Remove(path); err != nil {
				t.Fatal(err)
			}
			writeMinifyFixture(t, path, prior)
			status := "restored"
			if retry == "result-missing" {
				record, err := cache.readRecord("site.css")
				if err != nil {
					t.Fatal(err)
				}
				if err := os.Remove(filepath.Join(cache.dir, "blobs", record.OutputHash)); err != nil {
					t.Fatal(err)
				}
				status = "transformed"
			} else if retry == "another-recipe" {
				p.config.PreserveComments = []string{"License A", "License B"}
				recipe = minifyRecipe(p.Name(), p.config.PreserveComments)
				status = "transformed"
			}
			calls := 0
			retryTransform := func(data []byte) ([]byte, error) {
				calls++
				if !bytes.Equal(data, prior) {
					t.Fatal("target failure retry did not use current input")
				}
				return p.minifyBytes(data)
			}
			requireMinifyStatus(t, path, recipe, cache, retryTransform, status)
			expected, err := p.minifyBytes(prior)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(readMinifyFixture(t, path), expected) || (retry == "restore" && calls != 0) {
				t.Fatal("recipe-refresh retry did not recover exact result")
			}
		})
	}
}

func TestMinifyCache_UnavailableStorageFallback(t *testing.T) {
	logs := captureMinifyLog(t)
	config, old, path := newMinifyFixture(t, "js_minify")
	old.close()
	blocker := filepath.Join(config.ContentDir, "cache-is-file")
	config.Extra["cache_dir"] = blocker
	writeMinifyFixture(t, blocker, []byte("not a directory"))
	cache, err := newMinifyCache(config, "js_minify")
	if err == nil {
		t.Fatal("expected storage initialization error")
	}
	t.Cleanup(cache.close)
	writeMinifyFixture(t, path, []byte(" raw "))
	requireMinifyStatus(t, path, minifyRecipe("js_minify", nil), cache,
		func(data []byte) ([]byte, error) { return bytes.TrimSpace(data), nil }, "transformed")
	if string(readMinifyFixture(t, path)) != "raw" || !strings.Contains(logs.String(), "persisting cache") {
		t.Fatal("unavailable fresh storage did not produce explicit uncached fallback")
	}
}

func TestMinifyCache_FreshRecordPersistenceFailureFallsBack(t *testing.T) {
	logs := captureMinifyLog(t)
	_, cache, path := newMinifyFixture(t, "js_minify")
	writeMinifyFixture(t, path, []byte(" raw "))
	transform := func(data []byte) ([]byte, error) {
		// The current-byte check has already established no prior record.
		cache.close()
		return bytes.TrimSpace(data), nil
	}
	requireMinifyStatus(t, path, minifyRecipe("js_minify", nil), cache, transform, "transformed")
	if string(readMinifyFixture(t, path)) != "raw" || !strings.Contains(logs.String(), "no previous cache record") {
		t.Fatal("fresh cache failure unnecessarily blocked uncached transformation")
	}
}

func TestMinifyCache_PersistenceFailureInvalidatesOldProvenance(t *testing.T) {
	for _, blockInvalidation := range []bool{false, true} {
		t.Run(map[bool]string{false: "invalidate-succeeds", true: "invalidate-fails"}[blockInvalidation], func(t *testing.T) {
			logs := captureMinifyLog(t)
			config, cache, path := newMinifyFixture(t, "js_minify")
			sourceA, sourceB := []byte("source A"), []byte("source B")
			shared := []byte("same output")
			recipe := minifyRecipe("js_minify", nil)
			constant := func([]byte) ([]byte, error) { return bytes.Clone(shared), nil }
			writeMinifyFixture(t, path, sourceA)
			requireMinifyStatus(t, path, recipe, cache, constant, "transformed")
			writeMinifyFixture(t, path, sourceB)
			blocker := filepath.Join(cache.dir, "blobs", minifyHash(shared))
			if blockInvalidation {
				cache.close()
			} else {
				if err := os.Remove(blocker); err != nil {
					t.Fatal(err)
				}
				if err := os.Mkdir(blocker, 0o700); err != nil {
					t.Fatal(err)
				}
			}
			_, err := minifyAsset(path, recipe, cache, constant)
			if blockInvalidation {
				if err == nil || !strings.Contains(err.Error(), "leaving target unchanged") {
					t.Fatalf("error = %v, want failed safe invalidation", err)
				}
				if !bytes.Equal(readMinifyFixture(t, path), sourceB) {
					t.Fatal("mutated target while stale many-to-one provenance remained")
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}
				if record, err := cache.readRecord("site.js"); err != nil || record != nil {
					t.Fatal("obsolete source A relation survived uncached fallback")
				}
				if err := os.Remove(blocker); err != nil {
					t.Fatal(err)
				}
			}
			if !strings.Contains(logs.String(), "persisting cache") {
				t.Fatal("cache failure was silent")
			}
			reopened, err := newMinifyCache(config, "js_minify")
			if err != nil {
				t.Fatal(err)
			}
			defer reopened.close()
			seen := []byte(nil)
			revised := func(data []byte) ([]byte, error) {
				seen = bytes.Clone(data)
				return []byte("revised"), nil
			}
			requireMinifyStatus(t, path, minifyHash([]byte("revised recipe")), reopened, revised, "transformed")
			if bytes.Equal(seen, sourceA) {
				t.Fatal("recipe change resurrected source A after source B produced identical output")
			}
		})
	}
}

func TestMinifyCache_IsolationPermissionsAndHardlinks(t *testing.T) {
	shared := t.TempDir()
	var scopes []string
	for i := 0; i < 2; i++ {
		config, old, path := newMinifyFixture(t, "js_minify")
		old.close()
		config.Extra["cache_dir"] = shared
		cache, err := newMinifyCache(config, "js_minify")
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(cache.close)
		scopes = append(scopes, cache.scope)
		raw := []byte(strings.Repeat(" ", i+1) + "raw")
		writeMinifyFixture(t, path, raw)
		originalLink := filepath.Join(config.ContentDir, "original-link")
		if err := os.Link(path, originalLink); err != nil {
			t.Skipf("hardlinks unavailable: %v", err)
		}
		requireMinifyStatus(t, path, minifyRecipe("js_minify", nil), cache,
			func(data []byte) ([]byte, error) { return bytes.TrimSpace(data), nil }, "transformed")
		if !bytes.Equal(readMinifyFixture(t, originalLink), raw) {
			t.Fatal("target replacement mutated hardlink")
		}
		record, err := cache.readRecord("site.js")
		if err != nil {
			t.Fatal(err)
		}
		blobPath := filepath.Join(cache.dir, "blobs", record.OutputHash)
		blobLink := filepath.Join(config.ContentDir, "blob-link")
		if err := os.Link(blobPath, blobLink); err != nil {
			t.Fatal(err)
		}
		// Repair by atomic replacement must not mutate the linked prior inode.
		if err := cache.writePrivate(filepath.Join("blobs", record.OutputHash), []byte("replacement")); err != nil {
			t.Fatal(err)
		}
		if string(readMinifyFixture(t, blobLink)) != "raw" {
			t.Fatal("blob repair truncated hardlink")
		}
		if runtime.GOOS != "windows" {
			for _, entry := range []struct {
				path string
				mode os.FileMode
			}{{cache.dir, 0o700}, {filepath.Join(cache.dir, "records"), 0o700},
				{filepath.Join(cache.dir, minifyRecordPath("site.js")), 0o600}, {blobPath, 0o600}} {
				info, err := os.Stat(entry.path)
				if err != nil || info.Mode().Perm() != entry.mode {
					t.Fatalf("private mode %s = %v, error %v", entry.path, info, err)
				}
			}
		}
	}
	if scopes[0] == scopes[1] {
		t.Fatal("two sites sharing a cache with identical relative asset names collided")
	}
}

func TestMinifyCache_ConfiguredOutputScopes(t *testing.T) {
	config, first, _ := newMinifyFixture(t, "js_minify")
	config.OutputDir = filepath.Join(config.ContentDir, "other-output")
	second, err := newMinifyCache(config, "js_minify")
	if err != nil {
		t.Fatal(err)
	}
	defer second.close()
	if first.scope == second.scope {
		t.Fatal("configured output roots share scope")
	}
	config.Extra["build_cache"] = map[string]interface{}{"enabled": false}
	cache, err := newMinifyCache(config, "js_minify")
	if err != nil {
		t.Fatal(err)
	}
	defer cache.close()
	if cache.scope != second.scope {
		t.Fatal("build-cache enabled gate changed independent minification scope")
	}
}

func TestMinifyCache_UnsafePaths(t *testing.T) {
	for _, location := range []string{"direct", "cache-symlink", "output-symlink", "internal-symlink"} {
		t.Run(location, func(t *testing.T) {
			config, cache, _ := newMinifyFixture(t, "js_minify")
			cache.close()
			switch location {
			case "direct":
				config.Extra["cache_dir"] = filepath.Join(config.OutputDir, "private")
			case "cache-symlink":
				link := filepath.Join(config.ContentDir, "link")
				if err := os.Symlink(config.OutputDir, link); err != nil {
					t.Skipf("symlinks unavailable: %v", err)
				}
				config.Extra["cache_dir"] = filepath.Join(link, "missing")
			case "output-symlink":
				link := filepath.Join(config.ContentDir, "published-link")
				if err := os.Symlink(config.OutputDir, link); err != nil {
					t.Skipf("symlinks unavailable: %v", err)
				}
				config.OutputDir = link
				config.Extra["cache_dir"] = filepath.Join(cache.output, "private")
			case "internal-symlink":
				base := filepath.Join(config.ContentDir, "external-cache")
				if err := os.Mkdir(base, 0o700); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(config.OutputDir, filepath.Join(base, "asset-minify")); err != nil {
					t.Skipf("symlinks unavailable: %v", err)
				}
				config.Extra["cache_dir"] = base
			}
			unsafe, err := newMinifyCache(config, "js_minify")
			if unsafe != nil {
				defer unsafe.close()
			}
			if err == nil {
				t.Fatal("accepted cache location leading into output")
			}
			entries, err := os.ReadDir(cache.output)
			if err != nil || len(entries) != 0 {
				t.Fatal("private cache data created in output")
			}
		})
	}
	_, cache, path := newMinifyFixture(t, "js_minify")
	outside := filepath.Join(t.TempDir(), "private.js")
	writeMinifyFixture(t, outside, []byte("private"))
	if _, err := cache.identity(outside); err == nil {
		t.Fatal("outside target gained authority")
	}
	if _, err := cache.readBlob("../private"); err == nil {
		t.Fatal("invalid blob digest accepted")
	}
	if err := os.Symlink(outside, path); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	if _, err := minifyAsset(path, minifyRecipe("js_minify", nil), cache,
		func([]byte) ([]byte, error) { t.Fatal("followed target symlink"); return nil, nil }); err == nil {
		t.Fatal("symlink target accepted")
	}
	if string(readMinifyFixture(t, outside)) != "private" {
		t.Fatal("outside path mutated")
	}
}

func TestMinifyCache_DependencyIdentity(t *testing.T) {
	module := readMinifyFixture(t, filepath.Join("..", "..", "go.mod"))
	found := make(map[string]string)
	for _, line := range strings.Split(string(module), "\n") {
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			found[fields[0]] = fields[1]
		}
	}
	for dependency, version := range map[string]string{
		"github.com/tdewolff/minify/v2": assetMinifyVersion,
		"github.com/tdewolff/parse/v2":  assetParseVersion,
	} {
		if found[dependency] != version {
			t.Fatalf("recipe dependency identity %s = %s, go.mod = %s", dependency, version, found[dependency])
		}
	}
}

func TestMinifyCache_DerivedStorageOutputCollision(t *testing.T) {
	for _, collision := range []string{"asset-minify", "v1", "private-parent", "cache-symlink", "internal-symlink", "records-symlink"} {
		t.Run(collision, func(t *testing.T) {
			logs := captureMinifyLog(t)
			site := t.TempDir()
			base := filepath.Join(site, "cache")
			output := filepath.Join(base, "asset-minify")
			switch collision {
			case "v1":
				output = filepath.Join(output, "v1")
			case "private-parent":
				output = filepath.Join(output, "published")
			case "internal-symlink", "records-symlink":
				output = filepath.Join(site, "published")
			}
			if err := os.MkdirAll(output, 0o755); err != nil {
				t.Fatal(err)
			}
			config := lifecycle.NewConfig()
			config.ContentDir = site
			config.OutputDir = output
			config.Extra["cache_dir"] = base
			switch collision {
			case "cache-symlink":
				link := filepath.Join(site, "cache-link")
				if err := os.Symlink(base, link); err != nil {
					t.Skipf("symlinks unavailable: %v", err)
				}
				config.Extra["cache_dir"] = link
			case "internal-symlink":
				if err := os.MkdirAll(base, 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(output, filepath.Join(base, "asset-minify")); err != nil {
					t.Skipf("symlinks unavailable: %v", err)
				}
			case "records-symlink":
				scope := minifyHash([]byte(site + "\x00" + output + "\x00" + cssMinifyPluginName))
				dir := filepath.Join(base, "asset-minify", "v1", scope)
				if err := os.MkdirAll(dir, 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(output, filepath.Join(dir, "records")); err != nil {
					t.Skipf("symlinks unavailable: %v", err)
				}
			}
			path := filepath.Join(output, "site.css")
			writeMinifyFixture(t, path, []byte("/* private source comment */\n.CLASS { margin: 1em; }\n"))
			before := minifyDirectoryModes(t, site)
			cache, err := newMinifyCache(config, cssMinifyPluginName)
			if cache != nil {
				cache.close()
				t.Errorf("unsafe derived storage must disable caching, not return an active cache")
			}
			if err == nil || !strings.Contains(err.Error(), "unsafe asset cache") {
				t.Errorf("error = %v, want unsafe derived-cache warning", err)
			}
			if after := minifyDirectoryModes(t, site); !maps.Equal(before, after) {
				t.Fatal("cache validation created directories or changed permissions before rejecting unsafe storage")
			}
			manager := lifecycle.NewManager()
			manager.SetConfig(config)
			plugin := NewCSSMinifyPlugin()
			if err := plugin.Configure(manager); err != nil {
				t.Fatal(err)
			}
			if err := plugin.Write(manager); err != nil {
				t.Fatal(err)
			}
			if after := minifyDirectoryModes(t, site); !maps.Equal(before, after) {
				t.Fatal("uncached fallback created private directories or changed directory permissions")
			}
			if string(readMinifyFixture(t, path)) != ".CLASS{margin:1em}" ||
				!strings.Contains(logs.String(), "unsafe asset cache") {
				t.Fatal("unsafe derived storage did not produce warned, normal uncached output")
			}
			if entries, err := os.ReadDir(output); err != nil || len(entries) != 1 || entries[0].Name() != "site.css" {
				t.Fatal("private cache files appeared in published output")
			}
		})
	}
}

func minifyDirectoryModes(t *testing.T, root string) map[string]os.FileMode {
	t.Helper()
	modes := make(map[string]os.FileMode)
	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			modes[path] = info.Mode().Perm()
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	return modes
}

func TestMinifyCache_BaseAncestorCreationGuard(t *testing.T) {
	t.Run("missing-shared-parent", func(t *testing.T) {
		config := lifecycle.NewConfig()
		config.ContentDir = t.TempDir()
		parent := filepath.Join(config.ContentDir, "missing")
		config.OutputDir = filepath.Join(parent, "published")
		config.Extra["cache_dir"] = filepath.Join(parent, "cache")
		cache, err := newMinifyCache(config, cssMinifyPluginName)
		if cache != nil {
			cache.close()
			t.Fatal("unsafe missing cache ancestor did not disable caching")
		}
		if err == nil || !strings.Contains(err.Error(), "unsafe asset cache ancestor") {
			t.Fatalf("error = %v, want unsafe creation ancestor", err)
		}
		if _, err := os.Lstat(parent); !os.IsNotExist(err) {
			t.Fatal("validation created a private ancestor of published output")
		}
	})
	t.Run("missing-private-parent", func(t *testing.T) {
		config := lifecycle.NewConfig()
		config.ContentDir = t.TempDir()
		config.OutputDir = filepath.Join(config.ContentDir, "published")
		parent := filepath.Join(config.ContentDir, "missing", "private")
		config.Extra["cache_dir"] = filepath.Join(parent, "cache")
		cache, err := newMinifyCache(config, cssMinifyPluginName)
		if cache != nil {
			defer cache.close()
		}
		if err != nil {
			t.Fatalf("missing ancestors outside published output should be safe: %v", err)
		}
		if cache == nil {
			t.Fatal("safe missing cache ancestors disabled caching")
		}
		if _, err := os.Stat(filepath.Join(cache.dir, "records")); err != nil {
			t.Fatalf("safe private storage was not created: %v", err)
		}
		if _, err := os.Lstat(config.OutputDir); !os.IsNotExist(err) {
			t.Fatal("cache setup created published output")
		}
	})
	t.Run("existing-shared-parent", func(t *testing.T) {
		config := lifecycle.NewConfig()
		config.ContentDir = t.TempDir()
		config.OutputDir = filepath.Join(config.ContentDir, "published")
		config.Extra["cache_dir"] = config.ContentDir
		if err := os.Mkdir(config.OutputDir, 0o755); err != nil {
			t.Fatal(err)
		}
		before, err := os.Stat(config.ContentDir)
		if err != nil {
			t.Fatal(err)
		}
		cache, err := newMinifyCache(config, cssMinifyPluginName)
		if err != nil {
			t.Fatalf("existing shared base should be safe when private storage is a sibling: %v", err)
		}
		defer cache.close()
		after, err := os.Stat(config.ContentDir)
		if err != nil {
			t.Fatal(err)
		}
		if before.Mode() != after.Mode() {
			t.Fatal("shared existing base permissions changed")
		}
		if entries, err := os.ReadDir(config.OutputDir); err != nil || len(entries) != 0 {
			t.Fatal("safe sibling cache published private storage")
		}
	})
}

func TestMinifyCache_IdenticalFreshTransformPublication(t *testing.T) {
	for _, kind := range []string{"css_minify", "js_minify"} {
		for _, empty := range []bool{false, true} {
			for _, enabled := range []bool{false, true} {
				name := kind
				if empty {
					name += "/empty"
				}
				if enabled {
					name += "/cached"
				} else {
					name += "/uncached"
				}
				t.Run(name, func(t *testing.T) {
					config, cache, path := newMinifyFixture(t, kind)
					if !enabled {
						cache.close()
						cache = nil
					}
					transform := NewCSSMinifyPlugin().minifyBytes
					raw := []byte(".CLASS{margin:1em}")
					if kind == "js_minify" {
						transform = NewJSMinifyPlugin().minifyBytes
						raw = []byte("var value=1")
					}
					if empty {
						raw = nil
					}
					result, err := transform(raw)
					if err != nil {
						t.Fatal(err)
					}
					if !bytes.Equal(raw, result) {
						t.Fatalf("fixture does not reproduce a byte-identical fresh transformation: %q -> %q", raw, result)
					}
					writeMinifyFixture(t, path, raw)
					before := statMinifyFixture(t, path)
					release := filepath.Join(config.ContentDir, "old-release")
					if err := os.Link(path, release); err != nil {
						t.Skipf("hardlinks unavailable: %v", err)
					}
					recipe := minifyRecipe(kind, nil)
					requireMinifyStatus(t, path, recipe, cache, transform, "transformed")
					after := statMinifyFixture(t, path)
					if kind == "js_minify" && empty {
						if !os.SameFile(before, after) || before.Mode() != after.Mode() {
							t.Fatal("empty JavaScript no-write exception changed")
						}
					} else {
						if os.SameFile(before, after) {
							t.Error("byte-identical fresh transform skipped baseline atomic replacement")
						}
						if runtime.GOOS != "windows" && after.Mode().Perm() != 0o644 {
							t.Errorf("fresh target mode = %o, want 644", after.Mode().Perm())
						}
					}
					oldRelease := statMinifyFixture(t, release)
					if !os.SameFile(before, oldRelease) || before.Mode() != oldRelease.Mode() ||
						!bytes.Equal(readMinifyFixture(t, release), raw) {
						t.Fatal("atomic publication changed old hard-linked release")
					}
					if enabled {
						requireMinifyStatus(t, path, recipe, cache, transform, "restored")
						warm := statMinifyFixture(t, path)
						if after.Mode() != warm.Mode() ||
							(os.SameFile(after, warm) != (kind == "js_minify" && empty)) {
							t.Fatal("exact-input restore did not follow atomic publication/empty-JS rules")
						}
					}
				})
			}
		}
	}
}

func TestMinifyCache_LocationValidationIOFailureProtectsProvenance(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix directory search permissions required")
	}
	logs := captureMinifyLog(t)
	config, cache, path := newMinifyFixture(t, jsMinifyPluginName)
	sourceA, sourceB := []byte("source A"), []byte("source B")
	shared := []byte("same output")
	recipe := minifyRecipe(jsMinifyPluginName, nil)
	transform := func([]byte) ([]byte, error) { return bytes.Clone(shared), nil }
	writeMinifyFixture(t, path, sourceA)
	requireMinifyStatus(t, path, recipe, cache, transform, "transformed")
	writeMinifyFixture(t, path, sourceB)
	cache.close()
	if err := os.Chmod(cache.dir, 0); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chmod(cache.dir, 0o700); err != nil {
			t.Error(err)
		}
	})
	blocked, err := newMinifyCache(config, jsMinifyPluginName)
	if err == nil {
		blocked.close()
		t.Skip("directory permission failures unavailable")
	}
	if !errors.Is(err, os.ErrPermission) || blocked == nil {
		t.Fatalf("I/O failure must retain unavailable cache state, got cache %v, error %v", blocked, err)
	}
	defer blocked.close()
	runMinification(jsMinifyPluginName, []string{path}, func(string) bool { return false }, transform, 1, blocked, recipe)
	if !bytes.Equal(readMinifyFixture(t, path), sourceB) || !strings.Contains(logs.String(), "leaving target unchanged") {
		t.Fatal("location validation I/O failure silently allowed obsolete many-to-one source resurrection")
	}
}

func TestRunMinification_SeparatedCounts(t *testing.T) {
	logs := captureMinifyLog(t)
	_, cache, path := newMinifyFixture(t, "js_minify")
	writeMinifyFixture(t, path, []byte("raw"))
	failed := filepath.Join(cache.output, "failed.js")
	excluded := filepath.Join(cache.output, "excluded.js")
	runMinification("js_minify", []string{path, failed, excluded}, func(path string) bool { return filepath.Base(path) == "excluded.js" },
		func(data []byte) ([]byte, error) { return data, nil }, 2, cache, minifyRecipe("js_minify", nil))
	if !strings.Contains(logs.String(), "1 transformed, 0 restored, 1 excluded, 1 failed") {
		t.Fatal(logs.String())
	}
}

func TestMinifyPlugins_ReconfigureAndSelection(t *testing.T) {
	for _, kind := range []string{"js_minify", "css_minify"} {
		t.Run(kind, func(t *testing.T) {
			config, cache, path := newMinifyFixture(t, kind)
			cache.close()
			manager := lifecycle.NewManager()
			manager.SetConfig(config)
			js, css := NewJSMinifyPlugin(), NewCSSMinifyPlugin()
			configure, write, excluded := js.Configure, js.Write, js.isExcluded
			if kind == "css_minify" {
				configure, write, excluded = css.Configure, css.Write, css.isExcluded
			}
			config.Extra[kind] = map[string]interface{}{"exclude": []interface{}{"site.*"}}
			if err := configure(manager); err != nil {
				t.Fatal(err)
			}
			if !excluded(path) {
				t.Fatal("configured exclusion missing")
			}
			config.Extra[kind] = map[string]interface{}{"exclude": []interface{}{"other.*"}}
			if err := configure(manager); err != nil {
				t.Fatal(err)
			}
			other := filepath.Join(config.OutputDir, "other"+filepath.Ext(path))
			if excluded(path) || !excluded(other) {
				t.Fatal("reconfiguration accumulated prior exclusions")
			}
			config.Extra = nil
			if err := configure(manager); err != nil {
				t.Fatal(err)
			}
			if excluded(path) || excluded(other) {
				t.Fatal("nil Extra retained obsolete exclusions")
			}
			for _, gate := range []string{"fast", "disabled", "excluded", "pagefind", "min-js"} {
				if gate == "min-js" && kind == "css_minify" {
					continue
				}
				testPath := path
				config.Extra = make(map[string]interface{})
				switch gate {
				case "fast":
					config.Extra["fast_mode"] = true
				case "disabled":
					config.Extra[kind] = map[string]interface{}{"enabled": false}
				case "excluded":
					config.Extra[kind] = map[string]interface{}{"exclude": []interface{}{"site.*"}}
				case "pagefind":
					testPath = filepath.Join(config.OutputDir, "_pagefind", filepath.Base(path))
				case "min-js":
					testPath = filepath.Join(config.OutputDir, "site.min.js")
				}
				raw := []byte("/* untouched */")
				writeMinifyFixture(t, testPath, raw)
				if err := configure(manager); err != nil {
					t.Fatal(err)
				}
				if err := write(manager); err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(readMinifyFixture(t, testPath), raw) {
					t.Fatalf("%s changed gated asset", gate)
				}
				if err := os.Remove(testPath); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestMinifyCache_StaticAssetsFinalCycle(t *testing.T) {
	logs := captureMinifyLog(t)
	config, cache, _ := newMinifyFixture(t, "css_minify")
	cache.close()
	config.Extra["theme"] = "no-such-minify-test-theme"
	config.Extra["css_minify"] = models.CSSMinifyConfig{Enabled: true, PreserveComments: []string{"Copyright"}}
	cssRaw := []byte("/*! Copyright */\nbody { color: red; padding: 0px; }\n")
	jsRaw := []byte("function greet ( ) { return  42; }\n")
	writeMinifyFixture(t, filepath.Join(config.ContentDir, "static", "deep", "site.css"), cssRaw)
	writeMinifyFixture(t, filepath.Join(config.ContentDir, "static", "root.js"), jsRaw)
	manager := lifecycle.NewManager()
	manager.SetConfig(config)
	static, css, js := NewStaticAssetsPlugin(), NewCSSMinifyPlugin(), NewJSMinifyPlugin()
	if err := css.Configure(manager); err != nil {
		t.Fatal(err)
	}
	if err := js.Configure(manager); err != nil {
		t.Fatal(err)
	}
	cssExpected, err := css.minifyBytes(cssRaw)
	if err != nil {
		t.Fatal(err)
	}
	jsExpected, err := js.minifyBytes(jsRaw)
	if err != nil {
		t.Fatal(err)
	}
	for cycle := 0; cycle < 4; cycle++ {
		if err := static.Configure(manager); err != nil {
			t.Fatal(err)
		}
		if err := static.Write(manager); err != nil {
			t.Fatal(err)
		}
		if err := css.Write(manager); err != nil {
			t.Fatal(err)
		}
		if err := js.Write(manager); err != nil {
			t.Fatal(err)
		}
		if err := static.Cleanup(manager); err != nil {
			t.Fatal(err)
		}
		for asset, expected := range map[string][]byte{"deep/site.css": cssExpected, "root.js": jsExpected} {
			asset = filepath.FromSlash(asset)
			hash, ok := manager.AssetHashes()[asset]
			if !ok || hash == "" {
				t.Fatalf("cycle %d asset %s has no hash", cycle, asset)
			}
			ext := filepath.Ext(asset)
			alias := strings.TrimSuffix(asset, ext) + "." + hash + ext
			for _, relative := range []string{asset, alias} {
				if !bytes.Equal(readMinifyFixture(t, filepath.Join(config.OutputDir, relative)), expected) {
					t.Fatalf("cycle %d asset %s differs from uncached bytes", cycle, relative)
				}
			}
		}
	}
	if !strings.Contains(logs.String(), "restored") {
		t.Fatal("static raw recopies did not exercise result restoration")
	}
	if err := filepath.Walk(config.OutputDir, func(path string, _ os.FileInfo, err error) error {
		if err == nil && (strings.Contains(path, "asset-minify") || strings.HasSuffix(path, ".json")) {
			t.Fatalf("private data published: %s", path)
		}
		return err
	}); err != nil {
		t.Fatal(err)
	}
}

func TestMinifyCache_TransformFailureDoesNotPublish(t *testing.T) {
	_, cache, path := newMinifyFixture(t, "js_minify")
	raw := []byte("raw")
	writeMinifyFixture(t, path, raw)
	if _, err := minifyAsset(path, minifyRecipe("js_minify", nil), cache,
		func([]byte) ([]byte, error) { return nil, errors.New("minifier failed") }); err == nil {
		t.Fatal("transform error swallowed")
	}
	if !bytes.Equal(readMinifyFixture(t, path), raw) {
		t.Fatal("failed transform changed target")
	}
}

func TestMinifyCache_NoCacheAndEmptySource(t *testing.T) {
	for _, raw := range [][]byte{nil, []byte(" raw ")} {
		path := filepath.Join(t.TempDir(), "site.js")
		writeMinifyFixture(t, path, raw)
		requireMinifyStatus(t, path, minifyRecipe("js_minify", nil), nil,
			func(data []byte) ([]byte, error) { return bytes.TrimSpace(data), nil }, "transformed")
	}
	_, cache, path := newMinifyFixture(t, "js_minify")
	writeMinifyFixture(t, path, nil)
	requireMinifyStatus(t, path, minifyRecipe("js_minify", nil), cache, NewJSMinifyPlugin().minifyBytes, "transformed")
	record, err := cache.readRecord("site.js")
	if err != nil || record.SourceHash != minifyHash(nil) {
		t.Fatal("empty source confused with unavailable source")
	}
}

func TestMinifyCache_RecordSymlinkCannotEscape(t *testing.T) {
	_, cache, path := newMinifyFixture(t, "js_minify")
	outside := filepath.Join(t.TempDir(), "record.json")
	writeMinifyFixture(t, outside, []byte("private"))
	recordPath := filepath.Join(cache.dir, minifyRecordPath("site.js"))
	if err := os.Symlink(outside, recordPath); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	writeMinifyFixture(t, path, []byte(" raw "))
	logs := captureMinifyLog(t)
	requireMinifyStatus(t, path, minifyRecipe("js_minify", nil), cache,
		func(data []byte) ([]byte, error) { return bytes.TrimSpace(data), nil }, "transformed")
	if string(readMinifyFixture(t, outside)) != "private" || !strings.Contains(logs.String(), "Warning") {
		t.Fatal("record symlink escaped cache root or cache failure was silent")
	}
}

func TestMinifyCache_WorkerConcurrency(t *testing.T) {
	old := log.Writer()
	log.SetOutput(io.Discard)
	defer log.SetOutput(old)
	_, cache, _ := newMinifyFixture(t, "js_minify")
	var files []string
	for i := 0; i < 12; i++ {
		path := filepath.Join(cache.output, strings.Repeat("a", i+1)+".js")
		writeMinifyFixture(t, path, []byte(" raw "))
		files = append(files, path)
	}
	runMinification("js_minify", files, func(string) bool { return false },
		func(data []byte) ([]byte, error) { time.Sleep(time.Millisecond); return bytes.TrimSpace(data), nil },
		4, cache, minifyRecipe("js_minify", nil))
	for _, path := range files {
		if string(readMinifyFixture(t, path)) != "raw" {
			t.Fatal("worker result lost or borrowed pooled memory")
		}
	}
}
