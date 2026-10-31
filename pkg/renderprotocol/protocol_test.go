package renderprotocol

import (
	"fmt"
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/plugins"
)

func serveRequests(t *testing.T, input string) []Response {
	t.Helper()
	var output bytes.Buffer
	if err := Serve(strings.NewReader(input), &output); err != nil {
		t.Fatal(err)
	}
	var responses []Response
	for _, line := range strings.Split(strings.TrimSpace(output.String()), "\n") {
		if line == "" {
			continue
		}
		var response Response
		if err := json.Unmarshal([]byte(line), &response); err != nil {
			t.Fatalf("non-JSON protocol output %q: %v", line, err)
		}
		responses = append(responses, response)
	}
	return responses
}

func TestServe_MultipleRequestsAndCapabilities(t *testing.T) {
	input := `{"id":"a","op":"render","source":"# Hello"}
{"id":"b","op":"capabilities"}
{"id":"c","op":"render","source":"## Again","options":{"inline_assets":true}}
`
	responses := serveRequests(t, input)
	if len(responses) != 3 {
		t.Fatalf("expected three responses, got %d", len(responses))
	}
	if !responses[0].OK || responses[0].ID != "a" || !strings.Contains(responses[0].Result.HTML, "Hello") {
		t.Fatalf("bad first render: %+v", responses[0])
	}
	if !responses[1].OK || responses[1].Capabilities == nil ||
		responses[1].Capabilities.SchemaVersion != plugins.StandaloneSchemaVersion {
		t.Fatalf("missing capabilities: %+v", responses[1])
	}
	if !responses[2].OK || responses[2].ID != "c" || !strings.Contains(responses[2].Result.HeadHTML, "<style>") {
		t.Fatalf("bad inline render: %+v", responses[2])
	}
}

func TestServe_OneShotEquivalence(t *testing.T) {
	source := "# Exact output\n\n- one\n- two"
	oneShot, err := plugins.RenderStandalone(source, plugins.StandaloneOptions{InlineAssets: false})
	if err != nil {
		t.Fatal(err)
	}
	request, err := json.Marshal(Request{ID: "same", Op: "render", Source: source})
	if err != nil {
		t.Fatal(err)
	}
	responses := serveRequests(t, string(request)+"\n")
	if len(responses) != 1 || !responses[0].OK {
		t.Fatalf("expected successful response: %+v", responses)
	}
	if !reflect.DeepEqual(oneShot, responses[0].Result) {
		t.Fatal("stdio and one-shot renderer diverged")
	}
}

func TestServe_ErrorsDoNotKillStream(t *testing.T) {
	input := `{"id":"bad","op":"render","optionz":{}}
not json
{"id":"unknown","op":"frobnicate"}
{"id":"big","op":"render","source":"` + strings.Repeat("x", MaxSourceBytes+1) + `"}
{"id":"good","op":"render","source":"OK"}
`
	responses := serveRequests(t, input)
	if len(responses) != 5 {
		t.Fatalf("expected five responses, got %d", len(responses))
	}
	for i, code := range []string{"invalid_json", "invalid_json", "unsupported_operation", "source_too_large"} {
		if responses[i].OK || responses[i].Error == nil || responses[i].Error.Code != code {
			t.Fatalf("response %d: wanted %s, got %+v", i, code, responses[i])
		}
	}
	if !responses[4].OK || responses[4].ID != "good" {
		t.Fatalf("invalid line prevented next render: %+v", responses[4])
	}
}

func TestServe_OversizedFrameDrainsUntilNewline(t *testing.T) {
	input := strings.Repeat("x", MaxFrameBytes+1) + "\n" + `{"id":"after","op":"capabilities"}` + "\n"
	responses := serveRequests(t, input)
	if len(responses) != 2 || responses[0].Error == nil ||
		responses[0].Error.Code != "frame_too_large" || !responses[1].OK {
		t.Fatalf("frame boundary lost: %+v", responses)
	}
}

func TestServe_ShutdownStopsAtBoundary(t *testing.T) {
	responses := serveRequests(t, `{"id":"quit","op":"shutdown"}
{"id":"never","op":"render","source":"NOT REACHED"}
`)
	if len(responses) != 1 || !responses[0].OK || responses[0].ID != "quit" {
		t.Fatalf("shutdown did not terminate cleanly: %+v", responses)
	}
}

func TestDecodeRequest_RejectsExtraJSON(t *testing.T) {
	if _, err := DecodeRequest([]byte(`{"id":"a","op":"capabilities"} true`)); err == nil {
		t.Fatal("multiple JSON values accepted")
	}
	if _, err := DecodeRequest([]byte(`{"id":"a","op":"render","options":{"unknown":true}}`)); err == nil {
		t.Fatal("unknown nested option accepted")
	}
}

func TestHandle_RequiresID(t *testing.T) {
	r := Handle(Request{Op: "capabilities"})
	if r.OK || r.Error == nil || r.Error.Code != "invalid_id" {
		t.Fatalf("missing ID incorrectly accepted: %+v", r)
	}
}


func TestSessionCachesParserByOptions(t *testing.T) {
	s := &session{}
	first := s.handle(Request{ID: "1", Op: "render", Source: "# First"})
	second := s.handle(Request{ID: "2", Op: "render", Source: "# Second"})
	if !first.OK || !second.OK || !strings.Contains(second.Result.HTML, "Second") {
		t.Fatalf("cannot reuse parser: first=%+v second=%+v", first, second)
	}
	normal := Options{}
	cached := s.renderers[normal]
	if cached == nil || len(s.renderers) != 1 {
		t.Fatalf("missing reusable parser: %+v", s.renderers)
	}
	withHTML := s.handle(Request{ID: "3", Op: "render",
		Source: "<b>trusted</b>", Options: Options{Unsafe: true}})
	if !withHTML.OK || !strings.Contains(withHTML.Result.HTML, "<b>trusted</b>") {
		t.Fatalf("unsafe option lost: %+v", withHTML)
	}
	if s.renderers[normal] != cached || len(s.renderers) != 2 {
		t.Fatal("options shared state or original parser discarded")
	}
	for i := 0; i < 9; i++ {
		response := s.handle(Request{ID: "variant", Op: "render", Source: "plain",
			Options: Options{AssetBase: fmt.Sprintf("/assets-%d", i)}})
		if !response.OK {
			t.Fatalf("variant %d: %+v", i, response)
		}
	}
	if len(s.renderers) > 8 {
		t.Fatalf("unbounded parser cache: %d", len(s.renderers))
	}
}
