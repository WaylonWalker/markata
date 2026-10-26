package cmd

import (
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestEvaluateEncryptionKeyPolicy_MissingKey(t *testing.T) {
	cfg := models.NewConfig()
	cfg.Encryption.Enabled = true
	cfg.Encryption.DefaultKey = "default"
	cfg.Encryption.MinPasswordLength = 14
	cfg.Encryption.MinEstimatedCrackTime = "10y"

	t.Setenv("MARKATA_GO_ENCRYPTION_KEY_DEFAULT", "")

	results, _, _, err := evaluateEncryptionKeyPolicy(cfg, "")
	if err != nil {
		t.Fatalf("evaluateEncryptionKeyPolicy() error = %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("results length = %d, want 1", len(results))
	}
	if results[0].Err == nil {
		t.Fatal("expected missing key error")
	}
}

func TestEvaluateEncryptionKeyPolicy_StrongKey(t *testing.T) {
	cfg := models.NewConfig()
	cfg.Encryption.Enabled = true
	cfg.Encryption.DefaultKey = "default"
	cfg.Encryption.MinPasswordLength = 14
	cfg.Encryption.MinEstimatedCrackTime = "10y"

	t.Setenv("MARKATA_GO_ENCRYPTION_KEY_DEFAULT", "h7Qm!2Vx9#Lp4@Td")

	results, _, _, err := evaluateEncryptionKeyPolicy(cfg, "")
	if err != nil {
		t.Fatalf("evaluateEncryptionKeyPolicy() error = %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("results length = %d, want 1", len(results))
	}
	if results[0].Err != nil {
		t.Fatalf("expected success, got error: %v", results[0].Err)
	}
}
