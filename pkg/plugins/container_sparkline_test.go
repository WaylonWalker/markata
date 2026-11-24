package plugins

import (
	"strings"
	"testing"
)

func TestContainerSparklineRendersStaticSVG(t *testing.T) {
	// Container attributes come from a map, so check them independently of order.
	got := renderContainerMarkdown(t, `::: md-sparkline {data-values="10, 20 40 30" data-label="Jobs per year" data-tone="success"}

Up 3x since year one

:::`)
	for _, want := range []string{
		`<div class="md-sparkline"`,
		`data-values="10, 20 40 30"`,
		"\">\n<svg class=\"md-sparkline-chart\"",
		`<svg class="md-sparkline-chart" viewBox="0 0 100 24" preserveAspectRatio="none" role="img" aria-label="Jobs per year: 4 points from 10 to 30, low 10, high 40" focusable="false">`,
		`<polyline class="md-sparkline-line" points="0,22 33.33,15.33 66.67,2 100,8.67" vector-effect="non-scaling-stroke"/>`,
		`Up 3x since year one`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("want %q\ngot:\n%s", want, got)
		}
	}
	if strings.Contains(got, "md-sparkline-area") {
		t.Fatalf("area must be opt-in: %s", got)
	}
}

func TestContainerSparklineAreaModifierAddsFill(t *testing.T) {
	got := renderContainerMarkdown(t, "::: md-sparkline area {data-values=\"10 20 40 30\" data-label=\"Jobs\"}\n\nCaption\n\n:::")
	want := `<polygon class="md-sparkline-area" points="0,24 0,22 33.33,15.33 66.67,2 100,8.67 100,24"/><polyline class="md-sparkline-line"`
	if !strings.Contains(got, want) {
		t.Fatalf("want area before line %q\ngot:\n%s", want, got)
	}
}

func TestContainerSparklineFlatSeriesDrawsMidline(t *testing.T) {
	got := renderContainerMarkdown(t, "::: md-sparkline {data-values=\"5 5 5\" data-label=\"Steady\"}\n\nNo change\n\n:::")
	if !strings.Contains(got, `points="0,12 50,12 100,12"`) {
		t.Fatalf("flat series should draw through the middle: %s", got)
	}
}

func TestContainerSparklineEscapesLabel(t *testing.T) {
	got := renderContainerMarkdown(t, "::: md-sparkline {data-values=\"1 2\" data-label=\"A & B\"}\n\nCaption\n\n:::")
	if !strings.Contains(got, `aria-label="A &amp; B: 2 points`) {
		t.Fatalf("label must be escaped: %s", got)
	}
}

func TestContainerSparklineRejectsInvalidOrUnlabeledData(t *testing.T) {
	for _, attributes := range []string{
		`data-values="1 2 3"`,
		`data-values="" data-label="Empty"`,
		`data-values="7" data-label="One point"`,
		`data-values="1 NaN 3" data-label="NaN"`,
		`data-values="1 inf" data-label="Infinite"`,
		`data-values="1 two 3" data-label="Not a number"`,
		`data-values="1 2" data-label="   "`,
		`data-values="` + strings.Repeat("1 ", maxSparklinePoints+1) + `" data-label="Too many"`,
	} {
		t.Run(attributes[:min(len(attributes), 40)], func(t *testing.T) {
			got := renderContainerMarkdown(t, "::: md-sparkline {"+attributes+"}\n\nNo trend is available.\n\n:::")
			if strings.Contains(got, "<svg") {
				t.Fatalf("invalid data should not render a sparkline: %s", got)
			}
			if !strings.Contains(got, "No trend is available.") {
				t.Fatalf("fallback text must remain visible: %s", got)
			}
		})
	}
}

func TestContainerSparklineOnlyRendersWhenOptedIn(t *testing.T) {
	got := renderContainerMarkdown(t, "::: md-card {data-values=\"1 2 3\" data-label=\"Not a sparkline\"}\n\nRegular card.\n\n:::")
	if strings.Contains(got, "md-sparkline-chart") {
		t.Fatalf("non-sparkline container rendered an SVG: %s", got)
	}
}

func TestContainerSparklineFiniteExtremeValues(t *testing.T) {
	got := renderContainerMarkdown(t, `::: md-sparkline {data-values="-1e308 0 1e308" data-label="Extremes"}

Caption

:::`)
	if !strings.Contains(got, `points="0,22 50,12 100,2"`) {
		t.Fatalf("finite extremes should render a linear trend: %s", got)
	}
}
