package plugins

import (
	"strings"
	"testing"
)

func TestContainerGaugesRenderSemanticElements(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want string
	}{
		{
			name: "meter for utilization",
			src: `::: md-meter {data-value="18" data-max="100" data-label="Before cleanup: 18% free" data-tone="caution"}

18% free

:::`,
			want: `<meter value="18" max="100" aria-label="Before cleanup: 18% free">18 / 100</meter>`,
		},
		{
			name: "progress for tasks",
			src: `::: md-progress {data-value="0" data-max="10" data-label="Build completed units"}

0 of 10 complete

:::`,
			want: `<progress value="0" max="10" aria-label="Build completed units">0 / 10</progress>`,
		},
		{
			name: "decimal meter and escaped label",
			src: `::: md-meter {data-value="7.3" data-max="100" data-label="Storage A & B"}

7.3 GiB used

:::`,
			want: `<meter value="7.3" max="100" aria-label="Storage A &amp; B">7.3 / 100</meter>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderContainerMarkdown(t, tt.src)
			if !strings.Contains(got, tt.want) {
				t.Fatalf("want native gauge %q\ngot:\n%s", tt.want, got)
			}
		})
	}
}

func TestContainerGaugesRejectInvalidOrUnlabeledData(t *testing.T) {
	for _, attributes := range []string{
		`data-value="18" data-max="100"`,
		`data-value="18" data-max="0" data-label="Invalid max"`,
		`data-value="200" data-max="100" data-label="Over max"`,
		`data-value="-3" data-max="100" data-label="Negative"`,
		`data-value="NaN" data-max="100" data-label="NaN"`,
		`data-value="inf" data-max="100" data-label="Infinite"`,
		`data-value="oops" data-max="100" data-label="Not a number"`,
	} {
		t.Run(attributes, func(t *testing.T) {
			got := renderContainerMarkdown(t, "::: md-meter {"+attributes+"}\n\nNo numeric chart is available.\n\n:::")
			if strings.Contains(got, "<meter") {
				t.Fatalf("invalid numeric data should not render meter: %s", got)
			}
			if !strings.Contains(got, "No numeric chart is available.") {
				t.Fatalf("fallback text must remain visible: %s", got)
			}
		})
	}
}

func TestContainerGaugeOnlyAddsNativeHTMLWhenOptedIn(t *testing.T) {
	got := renderContainerMarkdown(t, `::: md-card {data-value="18" data-max="100" data-label="Not a gauge"}

Regular card.

:::`)
	if strings.Contains(got, "<meter") || strings.Contains(got, "<progress") {
		t.Fatalf("plain cards must not inject gauge elements: %s", got)
	}
}
