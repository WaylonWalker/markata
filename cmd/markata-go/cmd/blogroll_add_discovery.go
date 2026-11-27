package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/charmbracelet/huh"

	"github.com/WaylonWalker/markata/pkg/blogroll"
)

// verifyDirectBlogrollFeed is for an explicitly entered, trusted CLI URL.
// Unlike website discovery, it permits private feeds and nonstandard ports.
// Do not use this client for hosted editor or browser discovery.
func verifyDirectBlogrollFeed(ctx context.Context, client *http.Client, address string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, address, http.NoBody)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/rss+xml, application/atom+xml, application/xml, text/xml")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	const maxBytes = 8 << 20
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBytes+1))
	if err != nil {
		return err
	}
	if len(body) > maxBytes {
		return errors.New("feed document exceeds 8 MiB; use --manual if you trust this feed")
	}
	format, err := blogroll.ValidateFeedDocument(body)
	if err != nil {
		return err
	}
	if format == "json" {
		return errors.New("JSON Feed is not supported by the external reader; use RSS or Atom")
	}
	return nil
}

// chooseBlogrollCandidate never guesses when multiple subscriptions are
// advertised. Scripts must explicitly supply the desired feed URL.
func chooseBlogrollCandidate(candidates []string, interactive bool) (string, error) {
	switch len(candidates) {
	case 0:
		return "", errors.New("no validated feed candidates found")
	case 1:
		return candidates[0], nil
	}
	if !interactive {
		return "", &Diagnostic{
			Problem: "multiple feeds are advertised by this website",
			Details: candidates,
			Fix:     "Choose an exact feed URL and rerun blogroll add with that URL.",
		}
	}
	options := make([]huh.Option[string], 0, len(candidates))
	for _, address := range candidates {
		options = append(options, huh.NewOption(address, address))
	}
	selected := candidates[0]
	form := huh.NewForm(huh.NewGroup(huh.NewSelect[string]().
		Title("Choose a subscription").
		Options(options...).
		Value(&selected))).WithTheme(createHuhThemeForPalette(currentCLIPalette))
	if err := form.Run(); err != nil {
		return "", err
	}
	return selected, nil
}
