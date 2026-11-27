package prose

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"
)

// PreferencesFile is discovered by walking upward from a Markdown file.
const PreferencesFile = ".markata-prose.json"

// IgnoredPhrase identifies one accepted phrase without suppressing its entire rule.
type IgnoredPhrase struct {
	Rule   string `json:"rule"`
	Phrase string `json:"phrase"`
}

// Preferences holds project-level rules, phrase exceptions and accepted terms.
// A private per-user dictionary is a separate future feature.
type Preferences struct {
	IgnoredRules   []string        `json:"ignored_rules"`
	IgnoredPhrases []IgnoredPhrase `json:"ignored_phrases,omitempty"`
	AcceptedWords  []string        `json:"accepted_words,omitempty"`
}

func (p Preferences) Options() Options {
	ignored := make(map[string]bool, len(p.IgnoredRules))
	for _, id := range p.IgnoredRules {
		ignored[id] = true
	}
	phrases := make(map[string]map[string]bool)
	for _, entry := range p.IgnoredPhrases {
		if phrases[entry.Rule] == nil {
			phrases[entry.Rule] = make(map[string]bool)
		}
		phrases[entry.Rule][normalizePhrase(entry.Phrase)] = true
	}
	words := make(map[string]bool, len(p.AcceptedWords))
	for _, word := range p.AcceptedWords {
		words[strings.ToLower(word)] = true
	}
	return Options{IgnoredRules: ignored, IgnoredPhrases: phrases, AcceptedWords: words}
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
	return readPreferences(name)
}

func preferencesPath(path string) (string, error) {
	name, err := FindPreferences(path)
	if err != nil {
		return "", err
	}
	if name == "" {
		name = filepath.Join(filepath.Dir(path), PreferencesFile)
	}
	return name, nil
}

func rootPreferencesPath(root string) (string, error) {
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
	return filepath.Join(dir, PreferencesFile), nil
}

// IgnoreRule persists a rule-wide exception in a local project configuration.
func IgnoreRule(path, id string) (string, error) {
	name, err := preferencesPath(path)
	if err != nil {
		return "", err
	}
	return writeIgnoredRule(name, id)
}

// IgnoreRuleAtRoot confines browser rule writes to a server-selected root.
func IgnoreRuleAtRoot(root, id string) (string, error) {
	name, err := rootPreferencesPath(root)
	if err != nil {
		return "", err
	}
	return writeIgnoredRule(name, id)
}

// IgnorePhrase accepts a specific phrase in the nearest local project config.
func IgnorePhrase(path, id, phrase string) (string, error) {
	name, err := preferencesPath(path)
	if err != nil {
		return "", err
	}
	return writeIgnoredPhrase(name, id, phrase)
}

// IgnorePhraseAtRoot accepts a specific phrase at a server-selected root.
func IgnorePhraseAtRoot(root, id, phrase string) (string, error) {
	name, err := rootPreferencesPath(root)
	if err != nil {
		return "", err
	}
	return writeIgnoredPhrase(name, id, phrase)
}

// AddWord adds an accepted spelling term to the nearest project preferences.
func AddWord(path, word string) (string, error) {
	name, err := preferencesPath(path)
	if err != nil {
		return "", err
	}
	return writeAcceptedWord(name, word, true)
}

// RemoveWord revokes an accepted word without changing rule/phrase ignores.
func RemoveWord(path, word string) (string, error) {
	name, err := preferencesPath(path)
	if err != nil {
		return "", err
	}
	return writeAcceptedWord(name, word, false)
}

// AddWordAtRoot confines browser dictionary writes to a server-selected root.
func AddWordAtRoot(root, word string) (string, error) {
	name, err := rootPreferencesPath(root)
	if err != nil {
		return "", err
	}
	return writeAcceptedWord(name, word, true)
}

// CanonicalWord validates user-entered technical vocabulary and folds case.
// Single tokens may include common identifier characters but no separators,
// controls, whitespace, or paths. This is not a grammar-rule ignore.
func CanonicalWord(word string) (string, error) {
	if word == "" || len(word) > 128 || !utf8.ValidString(word) || strings.TrimSpace(word) != word {
		return "", fmt.Errorf("invalid dictionary word")
	}
	for _, r := range word {
		switch {
		case unicode.IsLetter(r), unicode.IsNumber(r), unicode.IsMark(r):
		case strings.ContainsRune("-_'+.#", r):
		default:
			return "", fmt.Errorf("invalid dictionary character")
		}
	}
	return strings.ToLower(word), nil
}

func writeAcceptedWord(name, word string, add bool) (string, error) {
	normalized, err := CanonicalWord(word)
	if err != nil {
		return "", err
	}
	return updatePreferences(name, func(p *Preferences) bool {
		index := -1
		for i, existing := range p.AcceptedWords {
			if strings.ToLower(existing) == normalized {
				index = i
				break
			}
		}
		if add {
			if index >= 0 {
				return false
			}
			p.AcceptedWords = append(p.AcceptedWords, normalized)
			return true
		}
		if index < 0 {
			return false
		}
		p.AcceptedWords = append(p.AcceptedWords[:index], p.AcceptedWords[index+1:]...)
		return true
	})
}

func writeIgnoredRule(name, id string) (string, error) {
	if !KnownRule(id) {
		return "", fmt.Errorf("unknown prose rule %q", id)
	}
	return updatePreferences(name, func(p *Preferences) bool {
		for _, existing := range p.IgnoredRules {
			if existing == id {
				return false
			}
		}
		p.IgnoredRules = append(p.IgnoredRules, id)
		return true
	})
}

func writeIgnoredPhrase(name, id, phrase string) (string, error) {
	if !ValidPhrase(id, phrase) {
		return "", fmt.Errorf("unknown or invalid prose phrase for %q", id)
	}
	normalized := normalizePhrase(phrase)
	return updatePreferences(name, func(p *Preferences) bool {
		for _, existing := range p.IgnoredPhrases {
			if existing.Rule == id && normalizePhrase(existing.Phrase) == normalized {
				return false
			}
		}
		p.IgnoredPhrases = append(p.IgnoredPhrases, IgnoredPhrase{Rule: id, Phrase: normalized})
		return true
	})
}

func readPreferences(name string) (Preferences, error) {
	info, statErr := os.Lstat(name)
	if statErr == nil && info.Mode()&os.ModeSymlink != 0 {
		return Preferences{}, fmt.Errorf("prose preference file must not be a symlink: %s", name)
	}
	if statErr != nil && !os.IsNotExist(statErr) {
		return Preferences{}, statErr
	}
	data, err := os.ReadFile(name)
	if os.IsNotExist(err) {
		return Preferences{}, nil
	}
	if err != nil {
		return Preferences{}, fmt.Errorf("read prose preferences %s: %w", name, err)
	}
	var p Preferences
	if err := json.Unmarshal(data, &p); err != nil {
		return Preferences{}, fmt.Errorf("parse prose preferences %s: %w", name, err)
	}
	return p, nil
}

func updatePreferences(name string, apply func(*Preferences) bool) (string, error) {
	p, err := readPreferences(name)
	if err != nil {
		return "", err
	}
	if !apply(&p) {
		return name, nil
	}
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
		_ = temp.Close()
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
