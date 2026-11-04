package plugins

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"math"
	"strings"
	"time"
)

var errNativeContributionGraphUnsupported = errors.New("unsupported native contribution graph layout")


const contributionGraphYearDomain = "year"

func nativeContributionGraphYear(options map[string]interface{}) (int, error) {
	for key, value := range options {
		switch key {
		case contributionGraphYearDomain, "domain", "subDomain", "maxValue", "maxPercentile":
		case "range":
			if n, ok := value.(float64); !ok || n != 1 {
				return 0, fmt.Errorf("%w: range must be 1; set contribution_graph.mode = \"client\"", errNativeContributionGraphUnsupported)
			}
		default:
			return 0, fmt.Errorf("%w: option %q; set contribution_graph.mode = \"client\"", errNativeContributionGraphUnsupported, key)
		}
	}
	if domain, ok := options["domain"]; ok && domain != contributionGraphYearDomain {
		return 0, fmt.Errorf("%w: native contribution graph supports domain = \"year\" only; set contribution_graph.mode = \"client\"", errNativeContributionGraphUnsupported)
	}
	if subDomain, ok := options["subDomain"]; ok && subDomain != "day" {
		return 0, fmt.Errorf("%w: native contribution graph supports subDomain = \"day\" only; set contribution_graph.mode = \"client\"", errNativeContributionGraphUnsupported)
	}

	year := time.Now().UTC().Year()
	if value, ok := options[contributionGraphYearDomain].(float64); ok {
		year = int(value)
		if float64(year) != value {
			return 0, fmt.Errorf("invalid contribution graph year %v", value)
		}
	}
	if year < 1 || year > 9998 {
		return 0, fmt.Errorf("contribution graph year must be between 1 and 9998")
	}
	return year, nil
}

func nativeContributionGraphValues(data interface{}) (map[string]int, error) {
	var points []contributionGraphDataPoint
	encoded, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("encode contribution data: %w", err)
	}
	if err := json.Unmarshal(encoded, &points); err != nil {
		return nil, fmt.Errorf("parse contribution data: %w", err)
	}
	values := make(map[string]int, len(points))
	for _, point := range points {
		day, err := time.Parse("2006-01-02", point.Date)
		if err != nil || day.Format("2006-01-02") != point.Date {
			return nil, fmt.Errorf("invalid contribution date %q", point.Date)
		}
		if point.Value < 0 {
			return nil, fmt.Errorf("negative contribution value for %s", point.Date)
		}
		values[point.Date] = point.Value
	}
	return values, nil
}

// renderNativeContributionGraph creates a responsive, palette-aware SVG calendar
// using only Go. It intentionally supports the GitHub-style year/day layout.
// Authors requiring other Cal-Heatmap domains can set mode = "client".
func renderNativeContributionGraph(data interface{}, options map[string]interface{}, maxValue int, graphID string) (string, error) {
	year, err := nativeContributionGraphYear(options)
	if err != nil {
		return "", err
	}
	values, err := nativeContributionGraphValues(data)
	if err != nil {
		return "", err
	}
	if maxValue < 1 {
		maxValue = 1
	}

	first := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := first.AddDate(1, 0, 0)
	gridStart := first.AddDate(0, 0, -int(first.Weekday()))
	dayCount := int(end.Sub(gridStart).Hours() / 24)
	weeks := (dayCount + 6) / 7
	const cell, step, left, top = 11, 14, 38, 26
	width := left + step*weeks + 4
	height := top + step*7 + 8

	var out strings.Builder
	fmt.Fprintf(&out, "<svg xmlns=\"http://www.w3.org/2000/svg\" class=\"contribution-graph-svg\" viewBox=\"0 0 %d %d\" role=\"img\" aria-labelledby=\"%s-title\" style=\"display:block;width:100%%;max-width:%dpx;height:auto;margin:auto\">", width, height, graphID, width)
	fmt.Fprintf(&out, "<title id=\"%s-title\">%d contribution calendar</title>", graphID, year)
	for month := time.January; month <= time.December; month++ {
		start := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
		week := int(start.Sub(gridStart).Hours()/24) / 7
		fmt.Fprintf(&out, "<text x=\"%d\" y=\"13\" font-size=\"10\" fill=\"var(--color-text-muted, #657076)\">%s</text>", left+step*week, month.String()[:3])
	}
	for _, day := range []struct {
		label string
		row   int
	}{
		{"Mon", 1}, {"Wed", 3}, {"Fri", 5},
	} {
		fmt.Fprintf(&out, "<text x=\"0\" y=\"%d\" font-size=\"10\" fill=\"var(--color-text-muted, #657076)\">%s</text>", top+step*day.row+9, day.label)
	}

	for date := gridStart; date.Before(end); date = date.AddDate(0, 0, 1) {
		if date.Before(first) {
			continue
		}
		key := date.Format("2006-01-02")
		value := values[key]
		week := int(date.Sub(gridStart).Hours()/24) / 7
		x, y := left+step*week, top+step*int(date.Weekday())
		fill := "var(--color-surface, #ebedf0)"
		opacity := 1.0
		if value > 0 {
			fill = "var(--color-primary, #216e39)"
			opacity = math.Min(1, 0.2+0.8*float64(value)/float64(maxValue))
		}
		fmt.Fprintf(&out, "<rect x=\"%d\" y=\"%d\" width=\"%d\" height=\"%d\" rx=\"2\" fill=\"%s\" fill-opacity=\"%.3f\"><title>%d posts on %s</title></rect>", x, y, cell, cell, fill, opacity, value, html.EscapeString(date.Format("Jan 2, 2006")))
	}
	out.WriteString("</svg>")
	return out.String(), nil
}
