package plugins

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

func TestStandaloneCapabilities_Fixture(t *testing.T) {
	fixture, err := os.ReadFile("testdata/renderer_contract/capabilities.json")
	if err != nil {
		t.Fatal(err)
	}
	var want map[string]any
	if err := json.Unmarshal(fixture, &want); err != nil {
		t.Fatal(err)
	}
	actual, err := json.Marshal(GetStandaloneCapabilities())
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err := json.Unmarshal(actual, &got); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("standalone capabilities drifted from public fixture\ngot: %s\nwant: %s", actual, fixture)
	}
}

func TestStandaloneCapabilities_DoNotClaimUnimplementedModes(t *testing.T) {
	caps := GetStandaloneCapabilities()
	if caps.SchemaVersion != 1 {
		t.Fatalf("breaking schema version: %d", caps.SchemaVersion)
	}
	if caps.LoadsSiteConfig || caps.ExternalURLAudit {
		t.Fatal("the standalone renderer neither loads site configuration nor audits authored external URLs")
	}
	if !reflect.DeepEqual(caps.DiagramModes["chartjs"], []string{"client"}) {
		t.Fatalf("chartjs advertises unimplemented modes: %v", caps.DiagramModes["chartjs"])
	}
	if !reflect.DeepEqual(caps.DiagramModes["mermaid"], availableMermaidModes()) {
		t.Fatalf("mermaid capabilities differ from supported build modes: got %v, want %v",
			caps.DiagramModes["mermaid"], availableMermaidModes())
	}
	caps.Features[0] = "modified"
	if GetStandaloneCapabilities().Features[0] == "modified" {
		t.Fatal("capabilities must not share mutable backing state")
	}
}
