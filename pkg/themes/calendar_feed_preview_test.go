package themes

import (
	"strings"
	"testing"
)

func TestCalendarFeedPreview_IsLazyAndUsesDropperThumbs(t *testing.T) {
	content, err := ReadStatic("js/calendar-feed.js")
	if err != nil {
		t.Fatalf("ReadStatic(calendar-feed.js) error = %v", err)
	}

	js := strings.ReplaceAll(string(content), "\r\n", "\n")
	for _, needle := range []string{
		"const dropperHosts = new Set([",
		"function dropperThumbURL(raw)",
		"if (!stem.endsWith('_thumb'))",
		"url.pathname = `${stem}_thumb${url.pathname.slice(dot)}`;",
		"for (const key of ['w', 'h', 'width', 'height']) url.searchParams.delete(key);",
		"function previewImageCandidates(raw)",
		"const previewData = new WeakMap();",
		"previewData.set(details, { year, month, day, posts });",
		"const ensurePreview = (details) => {",
		"image.src = imageCandidates[0];",
		"image.src = imageCandidates[1];",
		"preview.style.pointerEvents = 'none';",
		"preview.style.pointerEvents = 'auto';",
		"event.preventDefault();",
	} {
		if !strings.Contains(js, needle) {
			t.Fatalf("calendar-feed.js missing responsive preview behavior %q", needle)
		}
	}
}

func TestCalendarFeedPagesOneYearAtATimeAndDefersRendering(t *testing.T) {
	content, err := ReadStatic("js/calendar-feed.js")
	if err != nil {
		t.Fatalf("ReadStatic(calendar-feed.js) error = %v", err)
	}

	js := strings.ReplaceAll(string(content), "\r\n", "\n")
	for _, needle := range []string{
		"const sortedYears = Array.from(years).sort((a, b) => b - a);",
		"const initialYear = calendarYearFromURL(sortedYears) || sortedYears[0];",
		"const renderYear = (year, persistURL = true) => {",
		"calendar.replaceChildren();",
		"if (sortedYears.length > 1) calendar.append(yearNav);",
		"for (let month = 0; month < 12; month += 1)",
		"function calendarYearFromURL(years)",
		"url.searchParams.set('year', String(year));",
		"let calendarRendered = false;",
		"const ensureCalendar = async () => {",
		"fetch(calendarSource, { headers: { Accept: 'application/json' } })",
		"const ready = await ensureCalendar();",
		"if (request !== modeRequest || !root.isConnected) return;",
		"if (!ready) {",
	} {
		if !strings.Contains(js, needle) {
			t.Fatalf("calendar-feed.js missing year paging behavior %q", needle)
		}
	}

	if strings.Contains(js, "for (const year of Array.from(years).sort((a, b) => b - a))") {
		t.Fatalf("calendar-feed.js still eagerly renders every archive year")
	}
}
