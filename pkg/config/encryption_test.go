package config

import (
	"testing"

	"github.com/WaylonWalker/markata/pkg/encryption"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestParseTOML_EncryptionDefaults(t *testing.T) {
	data := []byte(`
[markata-go]
title = "Test"
`)

	cfg, err := ParseTOML(data)
	if err != nil {
		t.Fatalf("ParseTOML() error = %v", err)
	}

	if !cfg.Encryption.EnforceStrength {
		t.Error("EnforceStrength should default to true")
	}
	if cfg.Encryption.MinEstimatedCrackTime != encryption.DefaultMinEstimatedCrackTime {
		t.Errorf("MinEstimatedCrackTime = %q, want %q", cfg.Encryption.MinEstimatedCrackTime, encryption.DefaultMinEstimatedCrackTime)
	}
	if cfg.Encryption.MinPasswordLength != encryption.DefaultMinPasswordLength {
		t.Errorf("MinPasswordLength = %d, want %d", cfg.Encryption.MinPasswordLength, encryption.DefaultMinPasswordLength)
	}
}

func TestParseTOML_EncryptionOverrides(t *testing.T) {
	data := []byte(`
[markata-go.encryption]
enforce_strength = false
min_estimated_crack_time = "5d"
min_password_length = 20
`)

	cfg, err := ParseTOML(data)
	if err != nil {
		t.Fatalf("ParseTOML() error = %v", err)
	}

	if cfg.Encryption.EnforceStrength {
		t.Error("EnforceStrength should reflect override value")
	}
	if cfg.Encryption.MinEstimatedCrackTime != "5d" {
		t.Errorf("MinEstimatedCrackTime = %q, want %q", cfg.Encryption.MinEstimatedCrackTime, "5d")
	}
	if cfg.Encryption.MinPasswordLength != 20 {
		t.Errorf("MinPasswordLength = %d, want 20", cfg.Encryption.MinPasswordLength)
	}
}

func TestMergeEncryptionConfig(t *testing.T) {
	base := models.NewEncryptionConfig()
	override := models.EncryptionConfig{
		Enabled:               false,
		EnforceStrength:       false,
		MinEstimatedCrackTime: "1d",
		MinPasswordLength:     32,
	}

	result := mergeEncryptionConfig(base, override)
	if result.EnforceStrength {
		t.Error("EnforceStrength should follow override")
	}
	if result.MinEstimatedCrackTime != "1d" {
		t.Errorf("MinEstimatedCrackTime = %q, want %q", result.MinEstimatedCrackTime, "1d")
	}
	if result.MinPasswordLength != 32 {
		t.Errorf("MinPasswordLength = %d, want 32", result.MinPasswordLength)
	}
}
