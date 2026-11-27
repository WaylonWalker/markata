package prose

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// PreferencesFile is discovered by walking upward from a Markdown file.
const PreferencesFile = ".markata-prose.json"

// Preferences contains project-level rules shared by CLI and LSP.
// Personal dictionaries and per-occurrence ignores are a future extension.
type Preferences struct {
	IgnoredRules []string `json:"ignored_rules"`
}

func (p Preferences) Options() Options {
	ignored := make(map[string]bool, len(p.IgnoredRules))
	for _, id := range p.IgnoredRules {
		ignored[id] = true
	}
	return Options{IgnoredRules: ignored}
}

// FindPreferences locates the nearest project preferences file.
func FindPreferences(path string) (string, error) {
	dir, err := filepath.Abs(filepath.Dir(path))
	if err != nil {
		return "", err
	}
	for {
		candidate := filepath.Join(dir, PreferencesFile)
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		} else if !os.IsNotExist(err) {
			return "", fmt.Errorf("stat %s: %w", candidate, err)
		}
		next := filepath.Dir(dir)
		if next == dir {
			return "", nil
		}
		dir = next
	}
}

// LoadPreferences returns empty preferences if a project has not opted in.
func LoadPreferences(path string) (Preferences, error) {
	name, err := FindPreferences(path)
	if err != nil || name == "" {
		return Preferences{}, err
	}
	data, err := os.ReadFile(name)
	if err != nil {
		return Preferences{}, fmt.Errorf("read prose preferences %s: %w", name, err)
	}
	var p Preferences
	if err := json.Unmarshal(data, &p); err != nil {
		return Preferences{}, fmt.Errorf("parse prose preferences %s: %w", name, err)
	}
	return p, nil
}

// IgnoreRule persists a rule preference in the nearest configuration or
// creates a configuration alongside the supplied directory's Markdown files.
func IgnoreRule(path, id string) (string, error) {
	if id == "" {
		return "", fmt.Errorf("rule ID cannot be empty")
	}
	name, err := FindPreferences(path)
	if err != nil {
		return "", err
	}
	if name == "" {
		name = filepath.Join(filepath.Dir(path), PreferencesFile)
	}
	p, err := LoadPreferences(path)
	if err != nil {
		return "", err
	}
	for _, existing := range p.IgnoredRules {
		if existing == id {
			return name, nil
		}
	}
	p.IgnoredRules = append(p.IgnoredRules, id)
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return "", err
	}
	data = append(data, '\n')
	temp, err := os.CreateTemp(filepath.Dir(name), ".markata-prose-*")
	if err != nil {
		return "", err
	}
	defer os.Remove(temp.Name())
	if _, err := temp.Write(data); err != nil {
		temp.Close()
		return "", err
	}
	if err := temp.Close(); err != nil {
		return "", err
	}
	if err := os.Rename(temp.Name(), name); err != nil {
		return "", err
	}
	return name, nil
}
