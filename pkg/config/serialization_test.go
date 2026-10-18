package config

import (
	"bytes"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestCanonicalJSON_UsesNestedTheme(t *testing.T) {
	config := models.NewConfig()
	data, err := CanonicalJSON(config)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(data, []byte(`"contract_version"`)) || bytes.Contains(data, []byte(`"texture_strength"`)) {
		t.Fatalf("not canonical: %s", data)
	}
}

func TestCanonicalJSON_PreservesZeroDials(t *testing.T) {
	config := models.NewConfig()
	config.Theme.Texture.ColorMix = 0
	config.Theme.Motif.ColorMix = 0
	config.Theme.Motif.RowOffset = 0
	data, err := CanonicalJSON(config)
	if err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{`"color_mix": 0`, `"row_offset": 0`} {
		if !bytes.Contains(data, []byte(field)) {
			t.Fatalf("canonical JSON omitted zero field %s: %s", field, data)
		}
	}
}
