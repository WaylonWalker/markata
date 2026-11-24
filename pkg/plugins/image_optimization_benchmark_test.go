package plugins

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

const (
	imageOptimizationProductionTargetOccurrences = 3689
	imageOptimizationProductionUniqueTargets     = 18
)

func TestImageOptimizationProductionShapeDeduplicatesDiskWork(t *testing.T) {
	plugin, manager, targets := newRepeatedImageOptimizationFixture(t)
	outputDir := manager.Config().OutputDir
	cacheDir := plugin.config.CacheDir
	for _, target := range targets {
		if err := plugin.optimizeTarget(outputDir, cacheDir, target); err != nil {
			t.Fatalf("baseline optimizeTarget() error = %v", err)
		}
	}
	if got := plugin.stats.variantsChecked; got != 2*imageOptimizationProductionTargetOccurrences {
		t.Fatalf("baseline variants checked = %d, want %d", got, 2*imageOptimizationProductionTargetOccurrences)
	}
	if plugin.stats.cacheHits != plugin.stats.variantsChecked || plugin.stats.encoderCalls != 0 {
		t.Fatalf("baseline cache-hit stats = %+v, want every repeated variant to hit without encoding", plugin.stats)
	}
	plugin.stats = imageOptimizationStats{}

	if err := plugin.Write(manager); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if got := plugin.stats.targetOccurrences; got != imageOptimizationProductionTargetOccurrences {
		t.Fatalf("target occurrences = %d, want %d", got, imageOptimizationProductionTargetOccurrences)
	}
	if got := plugin.stats.uniqueTargets; got != imageOptimizationProductionUniqueTargets {
		t.Fatalf("unique targets = %d, want %d", got, imageOptimizationProductionUniqueTargets)
	}
	if got := plugin.stats.variantsChecked; got != 2*imageOptimizationProductionUniqueTargets {
		t.Fatalf("variants checked = %d, want %d", got, 2*imageOptimizationProductionUniqueTargets)
	}
	if plugin.stats.cacheHits != plugin.stats.variantsChecked || plugin.stats.cacheMisses != 0 || plugin.stats.encoderCalls != 0 {
		t.Fatalf("cache-hit stats = %+v, want every unique variant to hit without encoding", plugin.stats)
	}
}

func BenchmarkImageOptimizationRepeatedTargetDiskWork(b *testing.B) {
	for _, mode := range []string{"before-deduplication", "after-deduplication"} {
		b.Run(mode, func(b *testing.B) {
			plugin, manager, targets := newRepeatedImageOptimizationFixture(b)
			outputDir := manager.Config().OutputDir
			cacheDir := plugin.config.CacheDir
			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				plugin.stats = imageOptimizationStats{}
				if mode == "before-deduplication" {
					for _, target := range targets {
						if err := plugin.optimizeTarget(outputDir, cacheDir, target); err != nil {
							b.Fatal(err)
						}
					}
					continue
				}
				if err := plugin.Write(manager); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func newRepeatedImageOptimizationFixture(tb testing.TB) (*ImageOptimizationPlugin, *lifecycle.Manager, []imageOptimizationTarget) {
	tb.Helper()
	root := tb.TempDir()
	outputDir := filepath.Join(root, "output")
	cacheDir := filepath.Join(root, "image-cache")
	if err := os.MkdirAll(cacheDir, 0o700); err != nil {
		tb.Fatal(err)
	}
	plugin := NewImageOptimizationPlugin()
	plugin.config = defaultImageOptimizationConfig()
	plugin.config.Widths = nil
	plugin.config.CacheDir = cacheDir
	plugin.availableFormats = []string{formatAVIF, formatWebP}
	plugin.avifencPath = "/fixture/bin/avifenc"
	plugin.cwebpPath = "/fixture/bin/cwebp"

	uniqueTargets := make([]imageOptimizationTarget, imageOptimizationProductionUniqueTargets)
	for index := range uniqueTargets {
		sourcePath := filepath.Join(outputDir, "images", fmt.Sprintf("photo-%02d.jpg", index))
		if err := os.MkdirAll(filepath.Dir(sourcePath), 0o700); err != nil {
			tb.Fatal(err)
		}
		if err := os.WriteFile(sourcePath, []byte("fixture image"), 0o600); err != nil {
			tb.Fatal(err)
		}
		info, err := os.Stat(sourcePath)
		if err != nil {
			tb.Fatal(err)
		}
		src := fmt.Sprintf("/images/photo-%02d.jpg", index)
		uniqueTargets[index] = imageOptimizationTarget{Src: src}
		for _, format := range plugin.availableFormats {
			encoder := plugin.encoderPathForFormat(format)
			variant := buildImageVariants(sourcePath, format, nil)[0]
			if err := os.WriteFile(variant.Path, []byte("optimized fixture"), 0o600); err != nil {
				tb.Fatal(err)
			}
			cachePath := imageOptimizationCachePath(cacheDir, sourcePath, format, variant.Width, plugin.qualityForFormat(format), encoder)
			if err := writeImageCache(cachePath, sourcePath, info, format, variant.Width, plugin.qualityForFormat(format), encoder); err != nil {
				tb.Fatal(err)
			}
		}
	}

	targets := make([]imageOptimizationTarget, imageOptimizationProductionTargetOccurrences)
	for index := range targets {
		targets[index] = uniqueTargets[index%len(uniqueTargets)]
	}
	manager := lifecycle.NewManager()
	manager.Config().OutputDir = outputDir
	post := models.NewPost("production-shape.md")
	post.Extra = map[string]any{"image_optimization": targets}
	manager.SetPosts([]*models.Post{post})
	return plugin, manager, targets
}
