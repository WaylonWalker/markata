package plugins

import (
	"fmt"
	"html"
	"math"
	"strings"
)

// renderNativeChartSVG intentionally accepts a conservative Chart.js subset.
// Auto mode falls back to Chart.js for styled, interactive, and complex charts.
// Strict Go mode surfaces an actionable error for unsupported configurations.
func renderNativeChartSVG(config map[string]interface{}, chartID, title string) (string, error) {
	chartType, ok := config["type"].(string)
	if !ok {
		return "", fmt.Errorf("chart type is required")
	}
	switch chartType {
	case "bar", "line", "scatter", "pie", "doughnut":
	default:
		return "", fmt.Errorf("chart type %q requires Chart.js; set chartjs.mode = \"client\"", chartType)
	}
	if opts, exists := config["options"]; exists && opts != nil {
		options, ok := opts.(map[string]interface{})
		if !ok || len(options) != 0 {
			return "", fmt.Errorf("custom Chart.js options require client mode")
		}
	}
	data, ok := config["data"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("chart data must be an object")
	}
	rawSets, ok := data["datasets"].([]interface{})
	if !ok || len(rawSets) != 1 {
		return "", fmt.Errorf("native charts currently require exactly one dataset")
	}
	dataset, ok := rawSets[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("chart dataset must be an object")
	}
	for key := range dataset {
		switch key {
		case "data", "label", "fill":
		default:
			return "", fmt.Errorf("custom dataset option %q requires client mode", key)
		}
	}
	if fill, exists := dataset["fill"]; exists {
		if chartType != "line" {
			return "", fmt.Errorf("fill is only supported for native line charts")
		}
		if _, ok := fill.(bool); !ok {
			return "", fmt.Errorf("fill must be a boolean for native charts")
		}
	}

	accessible := strings.TrimSpace(title)
	if accessible == "" {
		accessible = chartType + " chart"
	}
	var out strings.Builder
	const width, height = 720, 340
	fmt.Fprintf(&out, "<svg xmlns=\"http://www.w3.org/2000/svg\" class=\"chartjs-native-svg\" role=\"img\" aria-labelledby=\"%s-title\" viewBox=\"0 0 %d %d\" style=\"width:100%%;max-width:%dpx;height:auto\">", chartID, width, height, width)
	fmt.Fprintf(&out, "<title id=\"%s-title\">%s</title>", chartID, html.EscapeString(accessible))

	var err error
	switch chartType {
	case "scatter":
		err = renderNativeScatter(&out, dataset)
	case "pie", "doughnut":
		var labels []string
		labels, err = nativeChartLabels(data)
		if err == nil {
			var values []float64
			values, err = nativeChartNumbers(dataset, len(labels))
			if err == nil {
				err = renderNativePie(&out, chartType == "doughnut", labels, values)
			}
		}
	default:
		var labels []string
		labels, err = nativeChartLabels(data)
		if err == nil {
			var values []float64
			values, err = nativeChartNumbers(dataset, len(labels))
			if err == nil {
				filled, _ := dataset["fill"].(bool)
				err = renderNativeCartesian(&out, chartType, labels, values, filled)
			}
		}
	}
	if err != nil {
		return "", err
	}
	out.WriteString("</svg>")
	return out.String(), nil
}

func nativeChartLabels(data map[string]interface{}) ([]string, error) {
	raw, ok := data["labels"].([]interface{})
	if !ok || len(raw) == 0 || len(raw) > 60 {
		return nil, fmt.Errorf("native charts need 1 to 60 category labels")
	}
	labels := make([]string, 0, len(raw))
	for _, item := range raw {
		switch value := item.(type) {
		case string:
			labels = append(labels, value)
		case float64:
			labels = append(labels, fmt.Sprint(value))
		default:
			return nil, fmt.Errorf("native chart labels must be strings or numbers")
		}
	}
	return labels, nil
}

func nativeChartNumbers(dataset map[string]interface{}, count int) ([]float64, error) {
	raw, ok := dataset["data"].([]interface{})
	if !ok || len(raw) != count {
		return nil, fmt.Errorf("native dataset length must match category labels")
	}
	values := make([]float64, 0, len(raw))
	for _, item := range raw {
		value, ok := item.(float64)
		if !ok || math.IsNaN(value) || math.IsInf(value, 0) {
			return nil, fmt.Errorf("native chart values must be finite numbers")
		}
		values = append(values, value)
	}
	return values, nil
}

func nativeChartColor(index int) string {
	colors := []string{
		"var(--color-primary, #7c3aed)",
		"var(--color-link, #2563eb)",
		"var(--color-success, #059669)",
		"var(--color-warning, #d97706)",
		"var(--color-error, #dc2626)",
		"var(--color-text-muted, #6b7280)",
	}
	return colors[index%len(colors)]
}

func nativeChartShortLabel(label string) string {
	letters := []rune(label)
	if len(letters) > 13 {
		return string(letters[:12]) + "…"
	}
	return label
}

func nativeChartAxes(out *strings.Builder, minY, maxY float64) {
	const left, top, plotH, plotW = 65.0, 25.0, 245.0, 625.0
	if maxY <= minY {
		maxY = minY + 1
	}
	for tick := 0; tick <= 4; tick++ {
		value := minY + float64(tick)*(maxY-minY)/4
		y := top + plotH - float64(tick)*plotH/4
		fmt.Fprintf(out, "<line x1=\"%.1f\" y1=\"%.1f\" x2=\"%.1f\" y2=\"%.1f\" stroke=\"var(--color-border, #d1d5db)\" stroke-width=\"0.8\"/>", left, y, left+plotW, y)
		fmt.Fprintf(out, "<text x=\"%.1f\" y=\"%.1f\" text-anchor=\"end\" font-size=\"11\" fill=\"var(--color-text-muted, #6b7280)\">%g</text>", left-8, y+4, value)
	}
}

func renderNativeCartesian(out *strings.Builder, chartType string, labels []string, values []float64, filled bool) error {
	const left, top, plotH, plotW = 65.0, 25.0, 245.0, 625.0
	minY, maxY := math.Min(0, values[0]), math.Max(0, values[0])
	for _, value := range values {
		minY = math.Min(minY, value)
		maxY = math.Max(maxY, value)
	}
	if minY == maxY {
		maxY = minY + 1
	}
	nativeChartAxes(out, minY, maxY)
	valueY := func(value float64) float64 {
		return top + (maxY-value)*plotH/(maxY-minY)
	}
	baseY := valueY(0)
	n := len(values)
	column := plotW / float64(n)
	barWidth := math.Min(48, column*0.65)
	var path strings.Builder
	for i, value := range values {
		x := left + (float64(i)+0.5)*column
		if chartType == "line" {
			if n > 1 {
				x = left + float64(i)*plotW/float64(n-1)
			} else {
				x = left + plotW/2
			}
		}
		y := valueY(value)
		if chartType == "bar" {
			height := math.Abs(y-baseY)
			fmt.Fprintf(out, "<rect x=\"%.2f\" y=\"%.2f\" width=\"%.2f\" height=\"%.2f\" rx=\"3\" fill=\"%s\"><title>%s: %g</title></rect>", x-barWidth/2, math.Min(y, baseY), barWidth, height, nativeChartColor(0), html.EscapeString(labels[i]), value)
		} else {
			if i == 0 {
				fmt.Fprintf(&path, "M %.2f %.2f", x, y)
			} else {
				fmt.Fprintf(&path, " L %.2f %.2f", x, y)
			}
		}
		step := int(math.Ceil(float64(n) / 12))
		if i%step == 0 {
			fmt.Fprintf(out, "<text x=\"%.2f\" y=\"%d\" text-anchor=\"middle\" font-size=\"10\" fill=\"var(--color-text-muted, #6b7280)\">%s</text>", x, 291, html.EscapeString(nativeChartShortLabel(labels[i])))
		}
	}
	if chartType == "line" {
		if filled {
			fmt.Fprintf(out, "<path d=\"%s L %.2f %.2f L %.2f %.2f Z\" fill=\"%s\" fill-opacity=\"0.16\"/>", path.String(), left+plotW, baseY, left, baseY, nativeChartColor(0))
		}
		fmt.Fprintf(out, "<path d=\"%s\" fill=\"none\" stroke=\"%s\" stroke-width=\"2.5\" stroke-linecap=\"round\" stroke-linejoin=\"round\"/>", path.String(), nativeChartColor(0))
		for i, value := range values {
			x := left + plotW/2
			if n > 1 {
				x = left + float64(i)*plotW/float64(n-1)
			}
			fmt.Fprintf(out, "<circle cx=\"%.2f\" cy=\"%.2f\" r=\"4\" fill=\"%s\"><title>%s: %g</title></circle>", x, valueY(value), nativeChartColor(0), html.EscapeString(labels[i]), value)
		}
	}
	return nil
}

func renderNativeScatter(out *strings.Builder, dataset map[string]interface{}) error {
	raw, ok := dataset["data"].([]interface{})
	if !ok || len(raw) == 0 || len(raw) > 200 {
		return fmt.Errorf("native scatter charts need 1 to 200 points")
	}
	type point struct{ x, y float64 }
	points := make([]point, 0, len(raw))
	minX, maxX, minY, maxY := math.Inf(1), math.Inf(-1), math.Inf(1), math.Inf(-1)
	for _, rawPoint := range raw {
		pair, ok := rawPoint.(map[string]interface{})
		if !ok {
			return fmt.Errorf("scatter points must contain x and y")
		}
		x, xok := pair["x"].(float64)
		y, yok := pair["y"].(float64)
		if !xok || !yok || math.IsNaN(x) || math.IsNaN(y) || math.IsInf(x, 0) || math.IsInf(y, 0) {
			return fmt.Errorf("scatter points must contain finite x and y")
		}
		points = append(points, point{x, y})
		minX, maxX = math.Min(minX, x), math.Max(maxX, x)
		minY, maxY = math.Min(minY, y), math.Max(maxY, y)
	}
	if maxX == minX {
		maxX++
	}
	if maxY == minY {
		maxY++
	}
	nativeChartAxes(out, minY, maxY)
	const left, top, plotH, plotW = 65.0, 25.0, 245.0, 625.0
	// Show the numerical x scale: position alone is ambiguous for scatter data.
	for tick := 0; tick <= 4; tick++ {
		value := minX + float64(tick)*(maxX-minX)/4
		x := left + float64(tick)*plotW/4
		fmt.Fprintf(out, "<line x1=\"%.1f\" y1=\"%.1f\" x2=\"%.1f\" y2=\"%.1f\" stroke=\"var(--color-border, #d1d5db)\"/>", x, top+plotH, x, top+plotH+4)
		fmt.Fprintf(out, "<text x=\"%.1f\" y=\"%.1f\" text-anchor=\"middle\" font-size=\"11\" fill=\"var(--color-text-muted, #6b7280)\">%g</text>", x, top+plotH+21, value)
	}
	for _, p := range points {
		cx := left + (p.x-minX)*plotW/(maxX-minX)
		cy := top + (maxY-p.y)*plotH/(maxY-minY)
		fmt.Fprintf(out, "<circle cx=\"%.2f\" cy=\"%.2f\" r=\"5\" fill=\"%s\"><title>x: %g, y: %g</title></circle>", cx, cy, nativeChartColor(0), p.x, p.y)
	}
	return nil
}

func renderNativePie(out *strings.Builder, donut bool, labels []string, values []float64) error {
	if len(values) > 8 {
		return fmt.Errorf("native pie charts support up to 8 categories")
	}
	sum := 0.0
	for _, v := range values {
		if v < 0 {
			return fmt.Errorf("pie values must be nonnegative")
		}
		sum += v
	}
	if sum == 0 {
		return fmt.Errorf("pie chart needs a positive total")
	}
	const cx, cy, radius = 275.0, 165.0, 112.0
	start := -math.Pi/2
	for i, value := range values {
		if value == 0 {
			continue
		}
		delta := 2*math.Pi*value/sum
		if delta >= 2*math.Pi-1e-10 {
			fmt.Fprintf(out, "<circle cx=\"%.1f\" cy=\"%.1f\" r=\"%.1f\" fill=\"%s\"><title>%s: %g</title></circle>", cx, cy, radius, nativeChartColor(i), html.EscapeString(labels[i]), value)
		} else {
			end := start+delta
			x1, y1 := cx+radius*math.Cos(start), cy+radius*math.Sin(start)
			x2, y2 := cx+radius*math.Cos(end), cy+radius*math.Sin(end)
			large := 0
			if delta > math.Pi {
				large = 1
			}
			fmt.Fprintf(out, "<path d=\"M %.2f %.2f L %.2f %.2f A %.2f %.2f 0 %d 1 %.2f %.2f Z\" fill=\"%s\"><title>%s: %g</title></path>", cx, cy, x1, y1, radius, radius, large, x2, y2, nativeChartColor(i), html.EscapeString(labels[i]), value)
		}
		start += delta
		fmt.Fprintf(out, "<rect x=\"%d\" y=\"%d\" width=\"12\" height=\"12\" fill=\"%s\"/>", 420, 65+i*29, nativeChartColor(i))
		fmt.Fprintf(out, "<text x=\"%d\" y=\"%d\" font-size=\"12\" fill=\"var(--color-text, #111827)\">%s (%g)</text>", 439, 76+i*29, html.EscapeString(nativeChartShortLabel(labels[i])), value)
	}
	if donut {
		fmt.Fprintf(out, "<circle cx=\"%.1f\" cy=\"%.1f\" r=\"%.1f\" fill=\"var(--color-background, #ffffff)\"/>", cx, cy, radius*0.54)
	}
	return nil
}
