package builderadmin

import (
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/buildstats"
)

func TestBuildSpanViewOrdersDepthAndCriticalChain(t *testing.T) {
	t.Parallel()

	spans := []buildstats.SpanTiming{
		{ID: "child-b", ParentID: "root", Name: "late child", StartOffset: 3 * time.Second, Duration: 8 * time.Second, Status: "ok"},
		{ID: "root", Name: "root", StartOffset: 0, Duration: 12 * time.Second, Status: "ok"},
		{ID: "child-a", ParentID: "root", Name: "early child", StartOffset: time.Second, Duration: 2 * time.Second, Status: "ok"},
		{ID: "grandchild", ParentID: "child-b", Name: "latest leaf", StartOffset: 8 * time.Second, Duration: 2 * time.Second, Status: "ok", Attributes: map[string]string{"feed": "archive"}},
	}
	rows := buildSpanView(spans)
	if len(rows) != 4 {
		t.Fatalf("row count = %d, want 4", len(rows))
	}
	for i, want := range []string{"root", "child-a", "child-b", "grandchild"} {
		if rows[i].ID != want {
			t.Fatalf("row %d = %q, want %q", i, rows[i].ID, want)
		}
	}
	if rows[0].Depth != 0 || rows[1].Depth != 1 || rows[2].Depth != 1 || rows[3].Depth != 2 {
		t.Fatalf("depths = %d,%d,%d,%d", rows[0].Depth, rows[1].Depth, rows[2].Depth, rows[3].Depth)
	}
	if !rows[0].Critical || rows[1].Critical || !rows[2].Critical || !rows[3].Critical {
		t.Fatalf("critical flags = root:%v child-a:%v child-b:%v grandchild:%v", rows[0].Critical, rows[1].Critical, rows[2].Critical, rows[3].Critical)
	}
	if rows[3].StartMS != 8000 || rows[3].DurationMS != 2000 {
		t.Fatalf("grandchild timing = start:%d duration:%d", rows[3].StartMS, rows[3].DurationMS)
	}
	rows[3].Attributes["feed"] = "changed"
	if got := spans[3].Attributes["feed"]; got != "archive" {
		t.Fatalf("view mutation changed input attribute: %q", got)
	}
}

func TestBuildSpanViewHandlesMissingParentAndCycle(t *testing.T) {
	t.Parallel()

	rows := buildSpanView([]buildstats.SpanTiming{
		{ID: "orphan", ParentID: "missing", Name: "orphan", Duration: time.Second},
		{ID: "a", ParentID: "b", Name: "a", StartOffset: time.Second, Duration: time.Second},
		{ID: "b", ParentID: "a", Name: "b", StartOffset: 2 * time.Second, Duration: time.Second},
	})
	if len(rows) != 3 {
		t.Fatalf("row count = %d, want 3", len(rows))
	}
	for _, row := range rows {
		if row.Depth < 0 || row.Depth > 2 {
			t.Fatalf("unexpected depth for %s: %d", row.ID, row.Depth)
		}
	}
}

func TestBuildSpanViewCycleHasNoCompletionChain(t *testing.T) {
	t.Parallel()
	rows := buildSpanView([]buildstats.SpanTiming{
		{ID: "a", ParentID: "b", Duration: time.Second},
		{ID: "b", ParentID: "a", Duration: 2 * time.Second},
	})
	for _, row := range rows {
		if row.Critical {
			t.Fatalf("cycle row %q marked critical", row.ID)
		}
	}
}
