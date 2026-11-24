package plugins

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/palettes"
	"github.com/WaylonWalker/markata/pkg/terminalpage"
)

type postCopyPayloads struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Markdown    string `json:"markdown"`
	Text        string `json:"text"`
	MarkdownURL string `json:"markdown_url"`
	TextURL     string `json:"text_url"`
	ANSICurl    string `json:"ansi_curl"`
}

func (p postCopyPayloads) JSON() string {
	data, err := json.Marshal(p)
	if err != nil {
		return "{}"
	}
	return string(data)
}

func buildPostCopyPayloads(post *models.Post, config *lifecycle.Config, baseURL string) postCopyPayloads {
	postURL := buildAbsolutePostURL(baseURL, post)
	title := resolvePostTitle(post, "")
	postFormats := models.NewPostFormatsConfig()
	if config != nil {
		postFormats = resolvePostFormats(post, config)
	}

	markdownURL := ""
	textURL := ""
	ansiCurl := ""
	if post != nil && !post.Private {
		if postFormats.Markdown {
			markdownURL = buildAbsolutePostFormatURL(baseURL, post, ".md")
		}
		if postFormats.Text {
			textURL = buildAbsolutePostFormatURL(baseURL, post, ".txt")
		}
		if postFormats.ANSI {
			ansiCurl = buildPostCurlCommand(baseURL, post, ".ansi")
		}
	}

	// Generated export files are fetched only if the reader requests a copy.
	// Keep the inline payload as a fallback only when the format is disabled.
	markdown := ""
	if markdownURL == "" {
		markdown = buildPostCopyMarkdown(post, postURL)
	}
	text := ""
	if textURL == "" {
		text = buildPostCopyText(post, config, postURL)
	}

	return postCopyPayloads{
		Title:       title,
		URL:         postURL,
		Markdown:    markdown,
		Text:        text,
		MarkdownURL: markdownURL,
		TextURL:     textURL,
		ANSICurl:    ansiCurl,
	}
}

func buildAbsolutePostFormatURL(baseURL string, post *models.Post, ext string) string {
	if post == nil {
		return ""
	}

	href := strings.TrimSpace(post.Href)
	if href == "" {
		return ""
	}

	formattedHref := strings.TrimSuffix(href, "/") + ext
	if baseURL == "" {
		return formattedHref
	}

	base := strings.TrimRight(baseURL, "/")
	return base + formattedHref
}

func buildPostCurlCommand(baseURL string, post *models.Post, ext string) string {
	url := buildAbsolutePostFormatURL(baseURL, post, ext)
	if url == "" {
		return ""
	}
	return "curl " + url
}

func buildAbsolutePostURL(baseURL string, post *models.Post) string {
	if post == nil {
		return ""
	}
	href := post.Href
	if href == "" {
		href = "/"
	}
	if baseURL == "" {
		return href
	}
	base := strings.TrimRight(baseURL, "/")
	if href == "/" {
		return base + "/"
	}
	return base + href
}

func buildPostCopyMarkdown(post *models.Post, postURL string) string {
	if post == nil {
		return ""
	}
	markdown := strings.TrimSpace(post.Content)
	if markdown == "" {
		markdown = strings.TrimSpace(buildMarkdownContent(post))
	}

	var buf strings.Builder
	title := resolvePostTitle(post, "")
	if title != "" {
		buf.WriteString("# ")
		buf.WriteString(title)
		buf.WriteString("\n\n")
	}
	if postURL != "" {
		buf.WriteString("Source: ")
		buf.WriteString(postURL)
		buf.WriteString("\n\n")
	}
	buf.WriteString(markdown)

	return strings.TrimSpace(buf.String())
}

func buildPostCopyText(post *models.Post, config *lifecycle.Config, postURL string) string {
	body := buildTerminalPage(post, config, false)
	if postURL == "" {
		return body
	}
	if body == "" {
		return "Source: " + postURL
	}
	return body + "\n\nSource: " + postURL
}

func resolvePostTitle(post *models.Post, fallback string) string {
	if post != nil && post.PlainTitle() != "" {
		return post.PlainTitle()
	}
	if fallback != "" {
		return fallback
	}
	if post != nil {
		return post.Slug
	}
	return ""
}

func buildTerminalPage(post *models.Post, config *lifecycle.Config, ansi bool) string {
	if post == nil {
		return ""
	}

	var buf strings.Builder

	if post.PlainTitle() != "" {
		buf.WriteString(post.PlainTitle())
		buf.WriteString("\n")
		buf.WriteString(strings.Repeat(terminalpage.DoubleRule, len([]rune(post.PlainTitle()))))
		buf.WriteString("\n\n")
	}
	if post.Description != nil && *post.Description != "" {
		buf.WriteString(*post.Description)
		buf.WriteString("\n\n")
	}
	if post.Date != nil {
		buf.WriteString("Date: ")
		buf.WriteString(post.Date.Format("January 2, 2006"))
		buf.WriteString("\n\n")
	}

	paletteName, variant := resolveTerminalPalette(config)
	chromaStyle := palettes.ChromaTheme(paletteName)
	if chromaStyle == "" {
		chromaStyle = palettes.ChromaThemeForVariant(variant)
	}

	source := post.ArticleHTML
	if strings.TrimSpace(source) == "" {
		source = post.HTML
	}
	if strings.TrimSpace(source) == "" {
		source = post.Content
	}

	body := terminalpage.RenderHTML(source, terminalpage.Options{
		ANSI:        ansi,
		Palette:     paletteName,
		ChromaStyle: chromaStyle,
	})

	mediaLinks := terminalMediaLinks(post, body)
	if len(mediaLinks) > 0 {
		if body != "" {
			buf.WriteString(strings.Join(mediaLinks, "\n"))
			buf.WriteString("\n\n")
		} else {
			buf.WriteString(strings.Join(mediaLinks, "\n"))
		}
	}
	buf.WriteString(body)

	return strings.TrimSpace(buf.String())
}

func terminalMediaLinks(post *models.Post, body string) []string {
	if post == nil {
		return nil
	}

	links := []string{}
	appendMediaLink := func(label, url string) {
		if url == "" || strings.Contains(body, url) {
			return
		}
		links = append(links, label+": "+url)
	}

	appendMediaLink("Image", getPostExtraString(post, "image", "cover_image", "og_image"))
	appendMediaLink("Video", getPostExtraString(post, "video"))

	return links
}

func buildMarkdownContent(post *models.Post) string {
	if post == nil {
		return ""
	}

	var buf strings.Builder
	buf.WriteString("---\n")

	if post.Title != nil {
		buf.WriteString(fmt.Sprintf("title: %q\n", *post.Title))
	}
	if post.Description != nil {
		buf.WriteString(fmt.Sprintf("description: %q\n", *post.Description))
	}
	if post.Date != nil {
		buf.WriteString(fmt.Sprintf("date: %s\n", post.Date.Format("2006-01-02")))
	}
	buf.WriteString(fmt.Sprintf("published: %t\n", post.Published))
	if post.Draft {
		buf.WriteString(fmt.Sprintf("draft: %t\n", post.Draft))
	}
	if len(post.Tags) > 0 {
		buf.WriteString("tags:\n")
		for _, tag := range post.Tags {
			buf.WriteString(fmt.Sprintf("  - %s\n", tag))
		}
	}
	if post.Template != "" && post.Template != defaultTemplate {
		buf.WriteString(fmt.Sprintf("template: %s\n", post.Template))
	}

	buf.WriteString("---\n\n")
	buf.WriteString(post.Content)

	return buf.String()
}
