// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestBridgingDetector_DetectSource(t *testing.T) {
	tests := []struct {
		name     string
		config   models.BridgesConfig
		url      string
		content  string
		expected MentionSource
	}{
		{
			name: "bridgy bluesky URL",
			config: models.BridgesConfig{
				Enabled:         true,
				BridgyFediverse: true,
				Bluesky:         true,
			},
			url:      "https://brid.gy/publish/bluesky/123",
			content:  "",
			expected: SourceBridgyBluesky,
		},
		{
			name: "direct bsky.app URL",
			config: models.BridgesConfig{
				Enabled:         true,
				BridgyFediverse: true,
				Bluesky:         true,
			},
			url:      "https://bsky.app/profile/user.bsky.social/post/123",
			content:  "",
			expected: SourceBridgyBluesky,
		},
		{
			name: "bridgy twitter URL",
			config: models.BridgesConfig{
				Enabled:         true,
				BridgyFediverse: true,
				Twitter:         true,
			},
			url:      "https://brid.gy/publish/twitter/user/123",
			content:  "",
			expected: SourceBridgyTwitter,
		},
		{
			name: "direct twitter URL",
			config: models.BridgesConfig{
				Enabled:         true,
				BridgyFediverse: true,
				Twitter:         true,
			},
			url:      "https://twitter.com/user/status/123",
			content:  "",
			expected: SourceBridgyTwitter,
		},
		{
			name: "x.com URL",
			config: models.BridgesConfig{
				Enabled:         true,
				BridgyFediverse: true,
				Twitter:         true,
			},
			url:      "https://x.com/user/status/123",
			content:  "",
			expected: SourceBridgyTwitter,
		},
		{
			name: "mastodon URL",
			config: models.BridgesConfig{
				Enabled:         true,
				BridgyFediverse: true,
				Mastodon:        true,
			},
			url:      "https://mastodon.social/@user/123",
			content:  "",
			expected: SourceBridgyMastodon,
		},
		{
			name: "fosstodon URL",
			config: models.BridgesConfig{
				Enabled:         true,
				BridgyFediverse: true,
				Mastodon:        true,
			},
			url:      "https://fosstodon.org/@user/123",
			content:  "",
			expected: SourceBridgyMastodon,
		},
		{
			name: "github URL",
			config: models.BridgesConfig{
				Enabled:         true,
				BridgyFediverse: true,
				GitHub:          true,
			},
			url:      "https://github.com/user/repo/issues/123",
			content:  "",
			expected: SourceBridgyGitHub,
		},
		{
			name: "regular web URL",
			config: models.BridgesConfig{
				Enabled:         true,
				BridgyFediverse: true,
				Bluesky:         true,
			},
			url:      "https://example.com/post/123",
			content:  "",
			expected: SourceWeb,
		},
		{
			name: "bluesky disabled",
			config: models.BridgesConfig{
				Enabled:         true,
				BridgyFediverse: true,
				Bluesky:         false,
			},
			url:      "https://bsky.app/profile/user/post/123",
			content:  "",
			expected: SourceWeb,
		},
		{
			name: "content-based detection",
			config: models.BridgesConfig{
				Enabled:         true,
				BridgyFediverse: true,
				Bluesky:         true,
			},
			url:      "https://brid.gy/comment/123",
			content:  "Liked on bsky.app",
			expected: SourceBridgyBluesky,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detector := NewBridgingDetector(tt.config)
			result := detector.DetectSource(tt.url, tt.content)
			if result != tt.expected {
				t.Errorf("DetectSource() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestBridgingDetector_EnrichMention(t *testing.T) {
	config := models.BridgesConfig{
		Enabled:         true,
		BridgyFediverse: true,
		Bluesky:         true,
		Twitter:         true,
		Mastodon:        true,
		GitHub:          true,
	}
	detector := NewBridgingDetector(config)

	tests := []struct {
		name             string
		mention          ReceivedWebMention
		expectedPlatform string
		expectedHandle   string
	}{
		{
			name: "bluesky mention",
			mention: ReceivedWebMention{
				Source: "https://bsky.app/profile/alice.bsky.social/post/123",
				Author: MentionAuthor{
					URL: "https://bsky.app/profile/alice.bsky.social",
				},
			},
			expectedPlatform: "bluesky",
			expectedHandle:   "@alice.bsky.social",
		},
		{
			name: "twitter mention",
			mention: ReceivedWebMention{
				Source: "https://twitter.com/alice/status/123",
				Author: MentionAuthor{
					URL: "https://twitter.com/alice",
				},
			},
			expectedPlatform: "twitter",
			expectedHandle:   "@alice",
		},
		{
			name: "mastodon mention",
			mention: ReceivedWebMention{
				Source: "https://mastodon.social/@alice/123",
				Author: MentionAuthor{
					URL: "https://mastodon.social/@alice",
				},
			},
			expectedPlatform: "mastodon",
			expectedHandle:   "@alice@mastodon.social",
		},
		{
			name: "github mention",
			mention: ReceivedWebMention{
				Source: "https://github.com/alice/repo/issues/123",
				Author: MentionAuthor{
					URL: "https://github.com/alice",
				},
			},
			expectedPlatform: "github",
			expectedHandle:   "@alice",
		},
		{
			name: "web mention",
			mention: ReceivedWebMention{
				Source: "https://example.com/post",
				Author: MentionAuthor{
					URL: "https://example.com/about",
				},
			},
			expectedPlatform: "web",
			expectedHandle:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mention := tt.mention
			detector.EnrichMention(&mention)

			if mention.Platform != tt.expectedPlatform {
				t.Errorf("Platform = %q, want %q", mention.Platform, tt.expectedPlatform)
			}
			if mention.Handle != tt.expectedHandle {
				t.Errorf("Handle = %q, want %q", mention.Handle, tt.expectedHandle)
			}
		})
	}
}

func TestBridgingDetector_ShouldAccept(t *testing.T) {
	tests := []struct {
		name     string
		config   models.BridgesConfig
		mention  ReceivedWebMention
		expected bool
	}{
		{
			name: "no filters - accept all",
			config: models.BridgesConfig{
				Enabled: true,
				Filters: models.BridgeFiltersConfig{},
			},
			mention: ReceivedWebMention{
				Platform:   "bluesky",
				WMProperty: "like-of",
				Source:     "https://bsky.app/test",
			},
			expected: true,
		},
		{
			name: "platform filter - accept matching",
			config: models.BridgesConfig{
				Enabled: true,
				Filters: models.BridgeFiltersConfig{
					Platforms: []string{"bluesky", "mastodon"},
				},
			},
			mention: ReceivedWebMention{
				Platform: "bluesky",
				Source:   "https://bsky.app/test",
			},
			expected: true,
		},
		{
			name: "platform filter - reject non-matching",
			config: models.BridgesConfig{
				Enabled: true,
				Filters: models.BridgeFiltersConfig{
					Platforms: []string{"bluesky", "mastodon"},
				},
			},
			mention: ReceivedWebMention{
				Platform: "twitter",
				Source:   "https://twitter.com/test",
			},
			expected: false,
		},
		{
			name: "interaction type filter - accept matching",
			config: models.BridgesConfig{
				Enabled: true,
				Filters: models.BridgeFiltersConfig{
					InteractionTypes: []string{"like", "repost"},
				},
			},
			mention: ReceivedWebMention{
				Platform:   "bluesky",
				WMProperty: "like-of",
				Source:     "https://bsky.app/test",
			},
			expected: true,
		},
		{
			name: "interaction type filter - reject non-matching",
			config: models.BridgesConfig{
				Enabled: true,
				Filters: models.BridgeFiltersConfig{
					InteractionTypes: []string{"like", "repost"},
				},
			},
			mention: ReceivedWebMention{
				Platform:   "bluesky",
				WMProperty: "in-reply-to",
				Source:     "https://bsky.app/test",
			},
			expected: false,
		},
		{
			name: "content length filter - accept long enough",
			config: models.BridgesConfig{
				Enabled: true,
				Filters: models.BridgeFiltersConfig{
					MinContentLength: 10,
				},
			},
			mention: ReceivedWebMention{
				Platform: "bluesky",
				Source:   "https://bsky.app/test",
				Content: MentionContent{
					Text: "This is a longer reply text",
				},
			},
			expected: true,
		},
		{
			name: "content length filter - reject too short",
			config: models.BridgesConfig{
				Enabled: true,
				Filters: models.BridgeFiltersConfig{
					MinContentLength: 100,
				},
			},
			mention: ReceivedWebMention{
				Platform: "bluesky",
				Source:   "https://bsky.app/test",
				Content: MentionContent{
					Text: "Short",
				},
			},
			expected: false,
		},
		{
			name: "blocked domain filter - reject blocked",
			config: models.BridgesConfig{
				Enabled: true,
				Filters: models.BridgeFiltersConfig{
					BlockedDomains: []string{"spam.com", "bad.net"},
				},
			},
			mention: ReceivedWebMention{
				Platform: "web",
				Source:   "https://spam.com/post/123",
			},
			expected: false,
		},
		{
			name: "blocked domain filter - accept non-blocked",
			config: models.BridgesConfig{
				Enabled: true,
				Filters: models.BridgeFiltersConfig{
					BlockedDomains: []string{"spam.com", "bad.net"},
				},
			},
			mention: ReceivedWebMention{
				Platform: "web",
				Source:   "https://good.com/post/123",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detector := NewBridgingDetector(tt.config)
			result := detector.ShouldAccept(&tt.mention)
			if result != tt.expected {
				t.Errorf("ShouldAccept() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMentionSource_String(t *testing.T) {
	tests := []struct {
		source   MentionSource
		expected string
	}{
		{SourceWeb, "web"},
		{SourceBridgyBluesky, "bluesky"},
		{SourceBridgyTwitter, "twitter"},
		{SourceBridgyMastodon, "mastodon"},
		{SourceBridgyGitHub, "github"},
		{SourceBridgyFlickr, "flickr"},
		{SourceCustomBridge, "custom"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.source.String(); got != tt.expected {
				t.Errorf("String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestReceivedWebMention_InteractionType(t *testing.T) {
	tests := []struct {
		property string
		expected string
	}{
		{"like-of", "like"},
		{"repost-of", "repost"},
		{"in-reply-to", "reply"},
		{"bookmark-of", "bookmark"},
		{"mention-of", "mention"},
		{"unknown", "mention"},
		{"", "mention"},
	}

	for _, tt := range tests {
		t.Run(tt.property, func(t *testing.T) {
			m := &ReceivedWebMention{WMProperty: tt.property}
			if got := m.InteractionType(); got != tt.expected {
				t.Errorf("InteractionType() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestExtractHandles(t *testing.T) {
	tests := []struct {
		name      string
		extractor func(string, string) string
		authorURL string
		sourceURL string
		expected  string
	}{
		{
			name:      "bluesky from author URL",
			extractor: extractBlueskyHandle,
			authorURL: "https://bsky.app/profile/alice.bsky.social",
			sourceURL: "",
			expected:  "@alice.bsky.social",
		},
		{
			name:      "twitter from author URL",
			extractor: extractTwitterHandle,
			authorURL: "https://twitter.com/alice",
			sourceURL: "",
			expected:  "@alice",
		},
		{
			name:      "twitter from x.com",
			extractor: extractTwitterHandle,
			authorURL: "https://x.com/bob",
			sourceURL: "",
			expected:  "@bob",
		},
		{
			name:      "github from author URL",
			extractor: extractGitHubHandle,
			authorURL: "https://github.com/alice",
			sourceURL: "",
			expected:  "@alice",
		},
		{
			name:      "mastodon with domain",
			extractor: extractMastodonHandle,
			authorURL: "https://mastodon.social/@alice",
			sourceURL: "",
			expected:  "@alice@mastodon.social",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.extractor(tt.authorURL, tt.sourceURL)
			if got != tt.expected {
				t.Errorf("extract = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestPlatformHelpers(t *testing.T) {
	// Test that all helpers return expected platforms
	colors := PlatformColors()
	emojis := PlatformEmoji()
	names := PlatformName()

	platforms := []string{"bluesky", "twitter", "mastodon", "github", "flickr", "web"}

	for _, p := range platforms {
		if _, ok := colors[p]; !ok {
			t.Errorf("PlatformColors missing %q", p)
		}
		if _, ok := emojis[p]; !ok {
			t.Errorf("PlatformEmoji missing %q", p)
		}
		if _, ok := names[p]; !ok {
			t.Errorf("PlatformName missing %q", p)
		}
	}
}
