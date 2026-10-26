package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/WaylonWalker/markata/pkg/encryption"
	"github.com/WaylonWalker/markata/pkg/models"
)

const encryptionEnvPrefix = "MARKATA_GO_ENCRYPTION_KEY_"

type encryptionKeyPolicyResult struct {
	KeyName    string
	EnvName    string
	Estimated  time.Duration
	Required   bool
	Configured bool
	Err        error
}

func evaluateEncryptionKeyPolicy(cfg *models.Config, keyFilter string) ([]encryptionKeyPolicyResult, time.Duration, int, error) {
	if cfg == nil {
		return nil, 0, 0, fmt.Errorf("config is nil")
	}

	minLength := cfg.Encryption.MinPasswordLength
	if minLength == 0 {
		minLength = encryption.DefaultMinPasswordLength
	}

	minDurationValue := cfg.Encryption.MinEstimatedCrackTime
	if minDurationValue == "" {
		minDurationValue = encryption.DefaultMinEstimatedCrackTime
	}

	minDuration, err := encryption.ParseEstimatedCrackDuration(minDurationValue)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("invalid encryption.min_estimated_crack_time: %w", err)
	}

	requiredKeys := configuredEncryptionKeys(cfg)
	if keyFilter != "" {
		requiredKeys = map[string]struct{}{strings.ToLower(keyFilter): {}}
	}

	keys := make([]string, 0, len(requiredKeys))
	for key := range requiredKeys {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	results := make([]encryptionKeyPolicyResult, 0, len(keys))
	for _, key := range keys {
		envName := encryptionEnvPrefix + strings.ToUpper(key)
		password := os.Getenv(envName)
		result := encryptionKeyPolicyResult{
			KeyName:    key,
			EnvName:    envName,
			Required:   true,
			Configured: password != "",
		}
		if password == "" {
			result.Err = fmt.Errorf("missing key in environment")
			results = append(results, result)
			continue
		}

		result.Estimated = encryption.EstimateCrackTime(password)

		if err := encryption.ValidatePassword(password, minLength, minDuration); err != nil {
			result.Err = err
			results = append(results, result)
			continue
		}
		results = append(results, result)
	}

	return results, minDuration, minLength, nil
}

func configuredEncryptionKeys(cfg *models.Config) map[string]struct{} {
	keys := make(map[string]struct{})
	defaultKey := strings.ToLower(strings.TrimSpace(cfg.Encryption.DefaultKey))
	if defaultKey != "" {
		keys[defaultKey] = struct{}{}
	}
	for _, key := range cfg.Encryption.PrivateTags {
		normalized := strings.ToLower(strings.TrimSpace(key))
		if normalized == "" {
			continue
		}
		keys[normalized] = struct{}{}
	}
	return keys
}
