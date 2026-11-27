package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/WaylonWalker/markata/pkg/blogroll"
	"github.com/WaylonWalker/markata/pkg/models"
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

func resolveBlogrollAddFeed(cmd *cobra.Command, cfg *models.Config, address string) (string, *blogroll.Metadata, error) {
	parent := cmd.Context()
	if parent == nil {
		parent = context.Background()
	}
	ctx, cancel := context.WithTimeout(parent, 60*time.Second)
	defer cancel()
	client := blogrollAddHTTPClientFactory(blogrollAddTimeout(cfg))
	defer client.CloseIdleConnections()
	directErr := verifyDirectBlogrollFeed(ctx, client, address)
	if directErr != nil {
		if ctx.Err() != nil {
			return "", nil, ctx.Err()
		}
		discovered, err := discoverSubscription(ctx, address)
		if err != nil {
			if ctx.Err() != nil {
				return "", nil, ctx.Err()
			}
			metadata, recoveryErr := recoverBlogrollMetadata(fmt.Errorf("direct feed: %w; website discovery: %w", directErr, err))
			return address, metadata, recoveryErr
		}
		interactive := !noInput && !blogrollAddNoPrompt && inputIsTerminal() && outputIsTerminal()
		address, err = chooseBlogrollCandidate(discovered.Candidates, interactive)
		if err != nil {
			return "", nil, err
		}
		publicClient := safeSubscriptionClient()
		defer publicClient.CloseIdleConnections()
		if err := verifyDirectBlogrollFeed(ctx, publicClient, address); err != nil {
			return "", nil, fmt.Errorf("verify discovered feed: %w", err)
		}
	}
	// The feed is verified already. Website metadata is enrichment, not a
	// reason to reject a usable subscription.
	metadata, err := fetchFeedMetadataForAdd(cfg, address)
	if err != nil {
		warnf("Feed verified, but metadata enrichment failed: %v", err)
		metadata = &blogroll.Metadata{}
	}
	return address, metadata, nil
}
