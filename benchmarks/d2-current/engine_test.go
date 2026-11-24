package d2native

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
	"testing"
)

func graph(nodes int) string {
	var source strings.Builder
	source.WriteString("direction: right\n")
	for i := 0; i < nodes-1; i++ {
		fmt.Fprintf(&source, "node_%d -> node_%d\n", i, i+1)
	}
	return source.String()
}

func compressedBytes(svg []byte) (int, error) {
	var compressed bytes.Buffer
	writer := gzip.NewWriter(&compressed)
	if _, err := writer.Write(svg); err != nil {
		return 0, err
	}
	if err := writer.Close(); err != nil {
		return 0, err
	}
	return compressed.Len(), nil
}

func TestNativeD2SVG(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}
	for _, fixture := range []struct {
		name   string
		source string
	}{
		{name: "small", source: graph(3)},
		{name: "medium", source: graph(18)},
		{name: "large", source: graph(72)},
	} {
		t.Run(fixture.name, func(t *testing.T) {
			svg, err := engine.Render(context.Background(), fixture.source)
			if err != nil {
				t.Fatal(err)
			}
			if len(svg) == 0 || len(svg) > 4<<20 {
				t.Fatalf("SVG size %d outside Markata's initial asset budget", len(svg))
			}
			if !bytes.Contains(svg, []byte("<svg")) {
				t.Fatal("D2 did not produce an SVG element")
			}
			decoder := xml.NewDecoder(bytes.NewReader(svg))
			for {
				if _, err := decoder.Token(); err != nil {
					if err == io.EOF {
						break
					}
					t.Fatalf("invalid SVG XML: %v", err)
				}
			}
			compressed, err := compressedBytes(svg)
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("nodes=%d svg_bytes=%d gzip_bytes=%d embedded_fonts=%t",
				strings.Count(fixture.source, " -> ")+1, len(svg), compressed,
				bytes.Contains(svg, []byte("data:font/")))
		})
	}
}

func BenchmarkNativeD2Warm(b *testing.B) {
	engine, err := NewEngine()
	if err != nil {
		b.Fatal(err)
	}
	for _, fixture := range []struct {
		name  string
		nodes int
	}{
		{name: "small", nodes: 3},
		{name: "medium", nodes: 18},
		{name: "large", nodes: 72},
	} {
		b.Run(fixture.name, func(b *testing.B) {
			source := graph(fixture.nodes)
			var svg []byte
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				svg, err = engine.Render(context.Background(), source)
				if err != nil {
					b.Fatal(err)
				}
			}
			b.StopTimer()
			compressed, err := compressedBytes(svg)
			if err != nil {
				b.Fatal(err)
			}
			b.ReportMetric(float64(len(svg)), "svg-bytes")
			b.ReportMetric(float64(compressed), "gzip-bytes")
		})
	}
}

func BenchmarkNativeD2EngineInit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := NewEngine(); err != nil {
			b.Fatal(err)
		}
	}
}
