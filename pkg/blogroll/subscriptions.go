package blogroll

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/WaylonWalker/markata/pkg/models"
)

// SubscriptionSource describes a resolved config layer, in merge order.
// Kind is nonempty for read-only layers (global configuration or overrides).
// Raw must contain the same normalized keys used by the config loader.
type SubscriptionSource struct {
	Path        string
	DisplayPath string
	Kind        string
	Raw         map[string]any
}

// SourceSet resolves ownership independently of CLI flags, HTTP handlers,
// terminal prompts, and process-wide workspace state.
type SourceSet []SubscriptionSource

// rawFeeds distinguishes an absent array from an explicitly empty array.
// The latter shadows earlier lists just like a nonempty array.
func rawFeeds(raw map[string]any) ([]map[string]any, bool, error) {
	root := raw
	if nested, ok := root["markata-go"].(map[string]any); ok {
		root = nested
	}
	section, ok := root["blogroll"].(map[string]any)
	if !ok {
		return nil, false, nil
	}
	value, present := section["feeds"]
	if !present {
		return nil, false, nil
	}
	switch entries := value.(type) {
	case []map[string]any:
		return entries, true, nil
	case []any:
		result := make([]map[string]any, 0, len(entries))
		for _, entry := range entries {
			feed, ok := entry.(map[string]any)
			if !ok {
				return nil, true, errors.New("subscription array contains a non-table entry")
			}
			result = append(result, feed)
		}
		return result, true, nil
	default:
		return nil, true, errors.New("subscription configuration must contain an array of feeds")
	}
}

// effective returns the last layer that defines the array. It never falls
// through a malformed or explicitly empty higher-precedence definition.
func (sources SourceSet) effective() (*SubscriptionSource, []map[string]any, error) {
	for i := len(sources) - 1; i >= 0; i-- {
		source := &sources[i]
		feeds, present, err := rawFeeds(source.Raw)
		if err != nil {
			return nil, nil, fmt.Errorf("read subscriptions in %s: %w", source.label(), err)
		}
		if present {
			return source, feeds, nil
		}
	}
	return nil, nil, nil
}

func (s SubscriptionSource) label() string {
	if s.DisplayPath != "" {
		return s.DisplayPath
	}
	return s.Path
}

// checkEditable deliberately limits the current transactional writer to TOML.
// A future writer can extend this policy without changing callers.
func (s SubscriptionSource) checkEditable() error {
	if s.Kind != "" {
		return fmt.Errorf("subscriptions belong to a %s config (%s); edit that source manually", s.Kind, s.label())
	}
	if s.Path == "" {
		return errors.New("no editable subscription configuration file")
	}
	if !strings.EqualFold(filepath.Ext(s.Path), ".toml") {
		return fmt.Errorf("editing %s subscriptions is not yet supported; use the config file directly", s.label())
	}
	return nil
}

// Owner returns an editable source containing exactly one matching URL in the
// effective array. It does not reveal the URL in errors.
func (sources SourceSet) Owner(address string) (SubscriptionSource, error) {
	source, entries, err := sources.effective()
	if err != nil {
		return SubscriptionSource{}, err
	}
	count := 0
	for _, entry := range entries {
		if value, ok := entry["url"].(string); ok && value == address {
			count++
		}
	}
	if source == nil || count != 1 {
		return SubscriptionSource{}, fmt.Errorf("could not uniquely identify subscription in the effective config file (found %d)", count)
	}
	if err := source.checkEditable(); err != nil {
		return SubscriptionSource{}, err
	}
	return *source, nil
}

// AddTarget selects the effective array owner, or invokes the caller's
// workspace-specific fallback if no layer defines an array.
func (sources SourceSet) AddTarget(fallback func() (SubscriptionSource, error)) (SubscriptionSource, error) {
	source, _, err := sources.effective()
	if err != nil {
		return SubscriptionSource{}, err
	}
	if source == nil {
		if fallback == nil {
			return SubscriptionSource{}, errors.New("no subscription configuration target")
		}
		target, err := fallback()
		if err != nil {
			return SubscriptionSource{}, err
		}
		source = &target
	}
	if err := source.checkEditable(); err != nil {
		return SubscriptionSource{}, err
	}
	return *source, nil
}

// SelectSubscription selects by exact URL or case-insensitive handle. It
// never mutates a guessed or substring match and rejects ambiguity.
func SelectSubscription(feeds []models.ExternalFeedConfig, selector string) (models.ExternalFeedConfig, error) {
	selector = strings.TrimSpace(selector)
	var match *models.ExternalFeedConfig
	for i := range feeds {
		feed := &feeds[i]
		if selector == "" || (feed.URL != selector && (feed.Handle == "" || !strings.EqualFold(feed.Handle, selector))) {
			continue
		}
		if match != nil {
			return models.ExternalFeedConfig{}, errors.New("subscription selector is ambiguous; use an exact feed URL")
		}
		match = feed
	}
	if match == nil {
		return models.ExternalFeedConfig{}, errors.New("no subscription matches; use its exact handle or feed URL")
	}
	return *match, nil
}

// SubscriptionView is a read-only inventory item usable by any frontend.
type SubscriptionView struct {
	models.ExternalFeedConfig
	Source   string `json:"source"`
	Editable bool   `json:"editable"`
}

// Inventory describes the effective subscriptions without network or disk I/O.
func (sources SourceSet) Inventory(feeds []models.ExternalFeedConfig) []SubscriptionView {
	result := make([]SubscriptionView, 0, len(feeds))
	for i := range feeds {
		item := &feeds[i]
		owner, err := sources.Owner(item.URL)
		view := SubscriptionView{ExternalFeedConfig: *item, Editable: err == nil}
		if err == nil {
			view.Source = owner.label()
		} else {
			view.Source = "Read only: " + err.Error()
		}
		result = append(result, view)
	}
	return result
}
