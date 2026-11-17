//go:build excalidraw_spike

package plugins

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/visualizations"
)

const excalidrawSpikeRuns = 6

func TestExcalidrawRendererSpike(t *testing.T) {
	peakRSS := startProcessTreeRSS(os.Getpid())
	defer peakRSS.close()
	fixtureDir := filepath.Join("testdata", "excalidraw", "renderer-corpus")
	outputDir := filepath.Join(os.TempDir(), "markata-excalidraw-spike", "chromium")
	if configuredOutput := os.Getenv("EXCALIDRAW_SPIKE_OUTPUT"); configuredOutput != "" {
		outputDir = configuredOutput
	}
	fixtureFilter := os.Getenv("EXCALIDRAW_SPIKE_FIXTURE")
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatal(err)
	}

	start := time.Now()
	config := models.NewExcalidrawConfig()
	config.Mode = excalidrawModeChromium
	renderer, err := newExcalidrawRenderer(config, true)
	if err != nil {
		t.Fatal(err)
	}
	initTime := time.Since(start)
	defer func() {
		if err := renderer.Close(); err != nil {
			t.Errorf("close renderer: %v", err)
		}
	}()

	entries, err := os.ReadDir(fixtureDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".excalidraw") && entry.Name() != "malformed.json" {
			continue
		}
		if fixtureFilter != "" && entry.Name() != fixtureFilter {
			continue
		}
		source, err := os.ReadFile(filepath.Join(fixtureDir, entry.Name()))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := validateExcalidrawScene(string(source)); err != nil {
			fmt.Printf("fixture=%s scene_error=%q\n", entry.Name(), err)
			continue
		}

		timings := make([]time.Duration, 0, excalidrawSpikeRuns)
		var output []byte
		for range excalidrawSpikeRuns {
			renderStart := time.Now()
			output, err = renderer.RenderSVG(context.Background(), visualizations.Request{Source: string(source)})
			timings = append(timings, time.Since(renderStart))
			if err != nil {
				t.Fatalf("render %s: %v", entry.Name(), err)
			}
		}

		asset, validationErr := visualizations.NewSVGAsset(output)
		if err := os.WriteFile(filepath.Join(outputDir, entry.Name()+".raw.svg"), output, 0o600); err != nil {
			t.Fatal(err)
		}
		validation := "accepted"
		if validationErr != nil {
			validation = validationErr.Error()
		} else if err := os.WriteFile(filepath.Join(outputDir, entry.Name()+".svg"), output, 0o600); err != nil {
			t.Fatal(err)
		}
		hash := sha256.Sum256(output)
		fmt.Printf("fixture=%s first_ms=%.3f warm_median_ms=%.3f svg_bytes=%d sha256=%s validator=%q external_font=%t\n",
			entry.Name(), float64(timings[0].Microseconds())/1000, durationMedianMillis(timings[1:]),
			len(output), hex.EncodeToString(hash[:]), validation,
			strings.Contains(string(output), "https://esm.sh/"))
		_ = asset
	}
	fmt.Printf("renderer=%s init_ms=%.3f process_tree_peak_rss_kb=%d output_dir=%s\n", renderer.Version(), float64(initTime.Microseconds())/1000, peakRSS.max.Load(), outputDir)
}

func durationMedianMillis(values []time.Duration) float64 {
	sorted := append([]time.Duration(nil), values...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	return float64(sorted[len(sorted)/2].Microseconds()) / 1000
}

type processTreeRSS struct {
	root int
	max  atomic.Int64
	done chan struct{}
}

func startProcessTreeRSS(root int) *processTreeRSS {
	sampler := &processTreeRSS{root: root, done: make(chan struct{})}
	go func() {
		ticker := time.NewTicker(5 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if rss := sampler.sample(); rss > sampler.max.Load() {
					sampler.max.Store(rss)
				}
			case <-sampler.done:
				return
			}
		}
	}()
	return sampler
}

func (sampler *processTreeRSS) close() {
	close(sampler.done)
}

func (sampler *processTreeRSS) sample() int64 {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return 0
	}
	parents := make(map[int]int, len(entries))
	rssByPID := make(map[int]int64, len(entries))
	for _, entry := range entries {
		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}
		status, err := os.ReadFile(filepath.Join("/proc", entry.Name(), "status"))
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(status), "\n") {
			if value, ok := strings.CutPrefix(line, "PPid:"); ok {
				parents[pid], _ = strconv.Atoi(strings.TrimSpace(value))
			}
			if value, ok := strings.CutPrefix(line, "VmRSS:"); ok {
				fields := strings.Fields(value)
				if len(fields) > 0 {
					rssByPID[pid], _ = strconv.ParseInt(fields[0], 10, 64)
				}
			}
		}
	}
	children := map[int]bool{sampler.root: true}
	for changed := true; changed; {
		changed = false
		for pid, parent := range parents {
			if children[parent] && !children[pid] {
				children[pid] = true
				changed = true
			}
		}
	}
	var total int64
	for pid := range children {
		total += rssByPID[pid]
	}
	return total
}
