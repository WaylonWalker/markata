package builderadmin

import (
	"sort"
	"time"

	"github.com/WaylonWalker/markata/pkg/buildstats"
)

// SpanViewRow is the bounded presentation model used by Builder Admin trace
// views. It keeps layout concerns out of the buildstats package.
type SpanViewRow struct {
	ID         string
	ParentID   string
	Name       string
	Stage      string
	Plugin     string
	Depth      int
	StartMS    int64
	DurationMS int64
	Status     string
	Critical   bool
	Attributes map[string]string
}

func buildSpanView(spans []buildstats.SpanTiming) []SpanViewRow {
	if len(spans) == 0 {
		return nil
	}
	ordered := append([]buildstats.SpanTiming(nil), spans...)
	sort.SliceStable(ordered, func(i, j int) bool {
		if ordered[i].StartOffset == ordered[j].StartOffset {
			return ordered[i].ID < ordered[j].ID
		}
		return ordered[i].StartOffset < ordered[j].StartOffset
	})

	byID := make(map[string]buildstats.SpanTiming, len(ordered))
	for _, span := range ordered {
		if span.ID != "" {
			byID[span.ID] = span
		}
	}
	critical := criticalSpanChain(ordered, byID)
	depthMemo := make(map[string]int, len(ordered))
	rows := make([]SpanViewRow, 0, len(ordered))
	for _, span := range ordered {
		rows = append(rows, SpanViewRow{
			ID:         span.ID,
			ParentID:   span.ParentID,
			Name:       span.Name,
			Stage:      span.Stage,
			Plugin:     span.Plugin,
			Depth:      spanDepth(span.ID, byID, depthMemo, map[string]bool{}),
			StartMS:    durationMilliseconds(span.StartOffset),
			DurationMS: durationMilliseconds(span.Duration),
			Status:     span.Status,
			Critical:   critical[span.ID],
			Attributes: cloneStringMap(span.Attributes),
		})
	}
	return rows
}

func criticalSpanChain(spans []buildstats.SpanTiming, byID map[string]buildstats.SpanTiming) map[string]bool {
	critical := make(map[string]bool)
	if len(spans) == 0 {
		return critical
	}
	parents := make(map[string]bool, len(spans))
	for _, span := range spans {
		if span.ParentID != "" {
			parents[span.ParentID] = true
		}
	}

	var terminal buildstats.SpanTiming
	var terminalEnd time.Duration
	foundLeaf := false
	for _, span := range spans {
		if span.ID == "" || parents[span.ID] {
			continue
		}
		end := span.StartOffset + span.Duration
		if !foundLeaf || end > terminalEnd || (end == terminalEnd && span.Duration > terminal.Duration) {
			terminal = span
			terminalEnd = end
			foundLeaf = true
		}
	}
	if !foundLeaf {
		return critical
	}
	for terminal.ID != "" && !critical[terminal.ID] {
		critical[terminal.ID] = true
		if terminal.ParentID == "" {
			break
		}
		parent, ok := byID[terminal.ParentID]
		if !ok {
			break
		}
		terminal = parent
	}
	return critical
}

func spanDepth(id string, byID map[string]buildstats.SpanTiming, memo map[string]int, visiting map[string]bool) int {
	if id == "" {
		return 0
	}
	if depth, ok := memo[id]; ok {
		return depth
	}
	if visiting[id] {
		return 0
	}
	span, ok := byID[id]
	if !ok || span.ParentID == "" {
		memo[id] = 0
		return 0
	}
	visiting[id] = true
	depth := 0
	if _, ok := byID[span.ParentID]; ok {
		depth = spanDepth(span.ParentID, byID, memo, visiting) + 1
	}
	delete(visiting, id)
	memo[id] = depth
	return depth
}

func durationMilliseconds(duration time.Duration) int64 {
	if duration <= 0 {
		return 0
	}
	return duration.Milliseconds()
}

func cloneStringMap(input map[string]string) map[string]string {
	if len(input) == 0 {
		return nil
	}
	clone := make(map[string]string, len(input))
	for key, value := range input {
		clone[key] = value
	}
	return clone
}
