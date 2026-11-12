package cmd

import (
	"bufio"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type newLinkMetadata struct {
	Title        string
	Description  string
	Image        string
	TrustedTitle bool
}

const (
	newTitleKey         = "title"
	newLinkHTMLSuffix   = ".html"
	newLinkActionCreate = "create"
	newLinkActionEdit   = "edit"
)

func validateNewLinkURL(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil || u.Hostname() == "" || (u.Scheme != subscriptionHTTP && u.Scheme != subscriptionHTTPS) {
		return fmt.Errorf("invalid link URL %q; enter an absolute HTTP(S) URL", rawURL)
	}
	return nil
}

func isNewLinkURL(value string) bool {
	lower := strings.ToLower(value)
	return strings.HasPrefix(lower, "http:") || strings.HasPrefix(lower, "https:")
}

func newLinkFallbackTitle(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	slug := path.Base(strings.TrimRight(u.Path, "/"))
	if slug != "." && slug != "/" && slug != "" {
		for _, suffix := range []string{newLinkHTMLSuffix, ".htm", ".php", ".aspx"} {
			slug = strings.TrimSuffix(slug, suffix)
		}
		return strings.Join(strings.Fields(strings.NewReplacer("-", " ", "_", " ").Replace(slug)), " ")
	}
	return u.Hostname()
}

func usableNewLinkTitle(title, rawURL string) bool {
	title = strings.TrimSpace(title)
	u, err := url.Parse(rawURL)
	if err != nil || title == "" || isNewLinkURL(title) {
		return false
	}
	host := strings.TrimPrefix(strings.ToLower(u.Hostname()), "www.")
	candidate := strings.TrimPrefix(strings.ToLower(strings.TrimSuffix(title, "/")), "www.")
	return candidate != host && !strings.EqualFold(candidate, u.Host)
}

func fetchNewLinkMetadata(ctx context.Context, rawURL string) (newLinkMetadata, error) {
	metadata := newLinkMetadata{Title: newLinkFallbackTitle(rawURL)}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, http.NoBody)
	if err != nil {
		return metadata, fmt.Errorf("create metadata request: %w", err)
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml")
	req.Header.Set("User-Agent", "markata (link post metadata)")
	client := &http.Client{Timeout: 8 * time.Second}
	response, err := client.Do(req)
	if err != nil {
		return metadata, fmt.Errorf("fetch link metadata: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return metadata, fmt.Errorf("fetch link metadata: HTTP %d", response.StatusCode)
	}
	const maxPageSize = 1024 * 1024
	body, err := io.ReadAll(io.LimitReader(response.Body, maxPageSize+1))
	if err != nil {
		return metadata, fmt.Errorf("read link metadata: %w", err)
	}
	if len(body) > maxPageSize {
		return metadata, fmt.Errorf("link metadata response exceeds 1 MiB")
	}
	document, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return metadata, fmt.Errorf("parse link metadata: %w", err)
	}
	values := make(map[string]string)
	document.Find("meta").Each(func(_ int, selection *goquery.Selection) {
		name, ok := selection.Attr("property")
		if !ok {
			name, _ = selection.Attr("name")
		}
		content, _ := selection.Attr("content")
		name = strings.ToLower(strings.TrimSpace(name))
		if values[name] == "" {
			values[name] = strings.Join(strings.Fields(content), " ")
		}
	})
	ogTitle := values["og:title"]
	metadata.TrustedTitle = usableNewLinkTitle(ogTitle, rawURL) && usableNewLinkTitle(ogTitle, response.Request.URL.String())
	if metadata.TrustedTitle {
		metadata.Title = ogTitle
	} else {
		for _, title := range []string{values["twitter:title"], document.Find("title").First().Text(), document.Find("h1").First().Text()} {
			title = strings.Join(strings.Fields(title), " ")
			if usableNewLinkTitle(title, rawURL) {
				metadata.Title = title
				break
			}
		}
	}
	metadata.Description = firstNewLinkValue(values["og:description"], values["twitter:description"], values["description"])
	image := firstNewLinkValue(values["og:image"], values["twitter:image"], values["twitter:image:src"])
	if image != "" {
		imageURL, parseErr := url.Parse(image)
		if parseErr != nil {
			return metadata, fmt.Errorf("invalid preview image URL: %w", parseErr)
		}
		resolved := response.Request.URL.ResolveReference(imageURL)
		if resolved.Scheme == subscriptionHTTP || resolved.Scheme == subscriptionHTTPS {
			metadata.Image = resolved.String()
		} else {
			return metadata, fmt.Errorf("preview image must use an HTTP(S) URL")
		}
	}
	return metadata, nil
}

func firstNewLinkValue(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func resolveNewLinkMetadata(cmd *cobra.Command, rawURL string) (newLinkMetadata, error) {
	if err := validateNewLinkURL(rawURL); err != nil {
		return newLinkMetadata{}, err
	}
	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	metadata, err := fetchNewLinkMetadata(ctx, rawURL)
	if ctx.Err() != nil {
		return metadata, ctx.Err()
	}
	if err != nil {
		errlnf("Warning: %v. You can enter missing fields manually.", err)
	}
	if flag := cmd.Flags().Lookup(newTitleKey); flag != nil && flag.Changed {
		title := strings.TrimSpace(flag.Value.String())
		if title == "" {
			return metadata, fmt.Errorf("--title must not be empty")
		}
		metadata.Title = title
		metadata.TrustedTitle = true
	}
	return metadata, nil
}

func newLinkTitleError(cmd *cobra.Command, rawURL, title string) error {
	quote := func(value string) string {
		return "'" + strings.ReplaceAll(value, "'", "'\"'\"'") + "'"
	}
	parts := []string{"markata", "new"}
	seen := make(map[string]bool)
	addFlag := func(flag *pflag.Flag) {
		if flag.Name == newTitleKey || flag.Name == "no-input" || seen[flag.Name] {
			return
		}
		seen[flag.Name] = true
		parts = append(parts, "--"+flag.Name+"="+quote(flag.Value.String()))
	}
	cmd.Flags().Visit(addFlag)
	cmd.InheritedFlags().Visit(addFlag)
	if !seen["template"] {
		parts = append(parts, "--template=link")
	}
	parts = append(parts, "--no-input", "--title="+quote(title), quote(rawURL))
	return fmt.Errorf("link title needs confirmation; no file was created. Review the suggested title, then run:\n  %s", strings.Join(parts, " "))
}

func applyNewLinkMetadata(template ContentTemplate, rawURL string, metadata newLinkMetadata) ContentTemplate {
	template.Frontmatter = maps.Clone(template.Frontmatter)
	if template.Frontmatter == nil {
		template.Frontmatter = make(map[string]interface{})
	}
	template.Frontmatter[newTitleKey] = metadata.Title
	template.Frontmatter["url"] = rawURL
	if metadata.Description != "" {
		template.Frontmatter["description"] = metadata.Description
	}
	if metadata.Image != "" {
		template.Frontmatter["image"] = metadata.Image
	}
	return template
}

func newLinkSlug(title, rawURL string) string {
	if slug := generateSlug(title); slug != "" {
		return slug
	}
	if slug := generateSlug(newLinkFallbackTitle(rawURL)); slug != "" {
		return slug
	}
	hash := sha256.Sum256([]byte(rawURL))
	return fmt.Sprintf("link-%x", hash[:6])
}

func promptNewLinkTitle(reader *bufio.Reader, suggestion string) (string, error) {
	errf("Title [%s]: ", suggestion)
	answer, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("title confirmation required: %w", err)
	}
	answer = strings.TrimSpace(answer)
	if answer == "" {
		answer = suggestion
	}
	if answer == "" {
		return "", fmt.Errorf("title is required")
	}
	return answer, nil
}

func runPlainNewLink(cmd *cobra.Command, template ContentTemplate, reader *bufio.Reader) (*interactiveInput, error) {
	rawURL := promptNew(reader, "Link URL", "")
	metadata, err := resolveNewLinkMetadata(cmd, rawURL)
	if err != nil {
		return nil, err
	}
	if !metadata.TrustedTitle {
		if noInput || !inputIsTerminal() {
			return nil, newLinkTitleError(cmd, rawURL, metadata.Title)
		}
		metadata.Title, err = promptNewLinkTitle(reader, metadata.Title)
		if err != nil {
			return nil, err
		}
	}
	template = applyNewLinkMetadata(template, rawURL, metadata)
	if !cmd.Flags().Changed("draft") {
		newDraft = starterDraftDefault(template)
	}
	action := promptNew(reader, "Action (create/edit)", newLinkActionCreate)
	if action != newLinkActionCreate && action != newLinkActionEdit {
		return nil, fmt.Errorf("unknown action %q; enter create or edit", action)
	}
	if action == newLinkActionEdit {
		if !cmd.Flags().Changed(newTitleKey) {
			metadata.Title, err = promptNewLinkTitle(reader, metadata.Title)
			if err != nil {
				return nil, err
			}
		}
		metadata.Description = promptNew(reader, "Description", metadata.Description)
		metadata.Image = promptNew(reader, "Preview Image", metadata.Image)
		template = applyNewLinkMetadata(template, rawURL, metadata)
		if !cmd.Flags().Changed("dir") {
			newDir = promptNew(reader, "Directory", template.Directory)
		}
		if !cmd.Flags().Changed("tags") {
			newTags = promptNew(reader, "Tags (comma-separated)", "")
		}
		if !cmd.Flags().Changed("draft") {
			newDraft = promptYesNoNew(reader, "Create as draft?", starterDraftDefault(template))
		}
	}
	var tags []string
	if newTags != "" {
		tags = parseTags(newTags)
	} else if cmd.Flags().Changed("tags") {
		tags = []string{}
	}
	return &interactiveInput{title: metadata.Title, template: template, tags: tags}, nil
}
