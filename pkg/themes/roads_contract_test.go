package themes

import (
	"strings"
	"testing"
)

func TestWindingRoadPointsStayWithinRoadFlow(t *testing.T) {
	js, err := ReadStatic("js/roads.js")
	if err != nil {
		t.Fatalf("read roads.js: %v", err)
	}
	source := string(js)

	if !strings.Contains(source, `p.el.style.position = "absolute";`) {
		t.Fatal("road point anchors must be removed from normal document flow")
	}
	if !strings.Contains(source, `{ x: w / 2, y: h }`) {
		t.Fatal("winding path must end at the road container boundary")
	}
	if strings.Contains(source, `{ x: w / 2, y: h * 1.02 }`) {
		t.Fatal("winding path must not extend into following article content")
	}
}
