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

// IgnoreRule persists a rule preference in the nearest project configuration.
// This discovery is intended for trusted local CLI and native LSP workflows.
func IgnoreRule(path, id string) (string, error) {
	name, err := FindPreferences(path)
	if err != nil {
		return "", err
	}
	if name == "" {
		name = filepath.Join(filepath.Dir(path), PreferencesFile)
	}
	return writeIgnoredRule(name, id)
}

// IgnoreRuleAtRoot confines browser preference writes to a server-selected
// workspace root. Unlike IgnoreRule, this never searches parent directories.
func IgnoreRuleAtRoot(root, id string) (string, error) {
	dir, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(dir)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("prose preference root is not a directory: %s", dir)
	}
	return writeIgnoredRule(filepath.Join(dir, PreferencesFile), id)
}

func writeIgnoredRule(name, id string) (string, error) {
	if !KnownRule(id) {
		return "", fmt.Errorf("unknown prose rule %q", id)
	}
	var p Preferences
	data, err := os.ReadFile(name)
	switch {
	case err == nil:
		if err := json.Unmarshal(data, &p); err != nil {
			return "", fmt.Errorf("parse prose preferences %s: %w", name, err)
		}
	case !os.IsNotExist(err):
		return "", fmt.Errorf("read prose preferences %s: %w", name, err)
	}
	for _, existing := range p.IgnoredRules {
		if existing == id {
			return name, nil
		}
	}
	p.IgnoredRules = append(p.IgnoredRules, id)
	output, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return "", err
	}
	output = append(output, '\n')
	temp, err := os.CreateTemp(filepath.Dir(name), ".markata-prose-*")
	if err != nil {
		return "", err
	}
	defer os.Remove(temp.Name())
	if _, err := temp.Write(output); err != nil {
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
