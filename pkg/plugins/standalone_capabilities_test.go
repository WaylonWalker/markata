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
	for name, modes := range caps.DiagramModes {
		if len(modes) != 1 || modes[0] != "client" {
			t.Fatalf("%s incorrectly advertises prerendering: %v", name, modes)
		}
	}
	caps.Features[0] = "modified"
	if GetStandaloneCapabilities().Features[0] == "modified" {
		t.Fatal("capabilities must not share mutable backing state")
	}
}
