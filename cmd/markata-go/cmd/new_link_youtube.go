package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func fetchNewLinkYouTubeMetadata(ctx context.Context, client *http.Client, rawURL string) (newLinkMetadata, error) {
	metadata := newLinkMetadata{Title: newLinkFallbackTitle(rawURL)}
	endpoint, err := url.Parse(youtubeOEmbedEndpoint)
	if err != nil {
		return metadata, fmt.Errorf("parse YouTube oEmbed endpoint: %w", err)
	}
	query := endpoint.Query()
	query.Set("url", rawURL)
	query.Set("format", "json")
	endpoint.RawQuery = query.Encode()
	body, _, err := fetchNewLinkData(ctx, client, endpoint.String(), "application/json")
	if err != nil {
		return metadata, fmt.Errorf("fetch YouTube oEmbed metadata: %w", err)
	}
	var response youtubeOEmbedResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return metadata, fmt.Errorf("parse YouTube oEmbed metadata: %w", err)
	}
	title := strings.Join(strings.Fields(response.Title), " ")
	if usableNewLinkTitle(title, rawURL) {
		metadata.Title = title
		metadata.TrustedTitle = true
	}
	if response.ThumbnailURL != "" {
		if err := validateNewLinkURL(response.ThumbnailURL); err != nil {
			return metadata, fmt.Errorf("invalid YouTube oEmbed thumbnail: %w", err)
		}
		metadata.Image = response.ThumbnailURL
	}
	if !metadata.TrustedTitle {
		return metadata, fmt.Errorf("YouTube oEmbed metadata has no usable video title")
	}
	return metadata, nil
}
