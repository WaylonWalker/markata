package plugins

import (
	"bytes"
	"strings"

	"github.com/WaylonWalker/markata/pkg/models"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// safeFeedPost returns a private-post copy suitable for public feed rendering.
// Public posts are returned unchanged. The copy prevents feed templates and
// serializers from reaching raw content or secret runtime fields on the source
// post while retaining the encrypted wrapper when one is available.
func safeFeedPost(post *models.Post) *models.Post {
	if post == nil || !post.Private {
		return post
	}

	safe := *post
	safe.Content = ""
	safe.HTML = ""
	safe.SecretKey = ""
	safe.InputHash = ""
	safe.RawFrontmatter = ""
	safe.PrevNextFeed = ""
	safe.PrivateOverride = nil
	safe.Toc = nil
	safe.TocPlaceholder = false
	safe.Hrefs = nil
	safe.Inlinks = nil
	safe.Outlinks = nil
	safe.Dependencies = nil
	safe.Prev = nil
	safe.Next = nil
	safe.PrevNextContext = nil
	safe.AuthorRoleOverrides = nil
	safe.AuthorDetailsOverrides = nil
	safe.AuthorObjects = nil
	safe.Templates = nil

	if !postExtraBool(post, "_title_explicit") {
		safe.Title = nil
		safe.TitleHTML = ""
		safe.TitleText = ""
		safe.TitleTextDerived = false
	}
	if !postExtraBool(post, "_description_explicit") {
		safe.Description = nil
	}
	if encryptedHTML, ok := safeEncryptedFeedContent(post); ok {
		safe.ArticleHTML = encryptedHTML
	} else {
		safe.ArticleHTML = ""
	}

	safe.Tags = append([]string(nil), post.Tags...)
	safe.Authors = append([]string(nil), post.Authors...)
	safe.Extra = safeFeedExtra(post)
	return &safe
}

func postExtraBool(post *models.Post, key string) bool {
	if post == nil {
		return false
	}
	value, ok := post.Get(key).(bool)
	return ok && value
}

// safeEncryptedFeedContent returns an encrypted wrapper after removing the
// page-only key name. Feed entries may be rendered outside the canonical page,
// so they must not disclose the key-selection metadata used by the page
// decryption UI. Invalid or ambiguous wrappers fail closed and are omitted.
func safeEncryptedFeedContent(post *models.Post) (string, bool) {
	if post == nil || post.ArticleHTML == "" {
		return "", false
	}
	if value, ok := post.Extra["has_encrypted_content"].(bool); ok && !value {
		return "", false
	}

	root, ok := parseEncryptedFeedWrapper(post.ArticleHTML)
	if !ok {
		return "", false
	}

	removedKeyName := false
	var removeKeyNameAttributes func(*html.Node)
	removeKeyNameAttributes = func(node *html.Node) {
		if node.Type == html.ElementNode {
			attrs := node.Attr[:0]
			for _, attr := range node.Attr {
				if strings.EqualFold(attr.Key, "data-key-name") {
					removedKeyName = true
					continue
				}
				attrs = append(attrs, attr)
			}
			node.Attr = attrs
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			removeKeyNameAttributes(child)
		}
	}
	removeKeyNameAttributes(root)

	// Preserve the original generated wrapper when it was already free of key
	// names. This avoids unnecessary serialization changes for older fixtures.
	if !removedKeyName {
		return post.ArticleHTML, true
	}

	var output bytes.Buffer
	if err := html.Render(&output, root); err != nil {
		return "", false
	}
	if strings.Contains(strings.ToLower(output.String()), "data-key-name") {
		return "", false
	}
	return output.String(), true
}

func parseEncryptedFeedWrapper(articleHTML string) (*html.Node, bool) {
	nodes, err := html.ParseFragment(strings.NewReader(articleHTML), &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Div,
		Data:     "div",
	})
	if err != nil {
		return nil, false
	}
	root := singleFeedWrapper(nodes)
	if root == nil || !isEncryptedFeedWrapper(root) {
		return nil, false
	}
	return root, true
}

func singleFeedWrapper(nodes []*html.Node) *html.Node {
	var root *html.Node
	for _, node := range nodes {
		if node.Type == html.TextNode && strings.TrimSpace(node.Data) == "" {
			continue
		}
		if root != nil {
			return nil
		}
		root = node
	}
	return root
}

func isEncryptedFeedWrapper(node *html.Node) bool {
	if node == nil || node.Type != html.ElementNode || !strings.EqualFold(node.Data, "div") {
		return false
	}

	hasClass := false
	hasCiphertext := false
	for _, attr := range node.Attr {
		switch {
		case strings.EqualFold(attr.Key, "class"):
			for _, className := range strings.Fields(attr.Val) {
				if className == "encrypted-content" {
					hasClass = true
					break
				}
			}
		case strings.EqualFold(attr.Key, "data-encrypted") && attr.Val != "":
			hasCiphertext = true
		}
	}
	return hasClass && hasCiphertext
}

func safeFeedExtra(post *models.Post) map[string]interface{} {
	if post == nil || post.Extra == nil {
		return nil
	}

	// These fields are authored metadata or internal state needed by the feed
	// renderer. Do not copy arbitrary Extra values: that map can contain secret
	// keys, decrypted source data, or content-derived plugin output.
	keys := []string{"_title_explicit", "_description_explicit", "aliases", "avatar", "category", "categories"}
	result := make(map[string]interface{}, len(keys)+1)
	for _, key := range keys {
		if value, ok := post.Extra[key]; ok {
			result[key] = value
		}
	}
	if hasEncryptedFeedContent(post) {
		result["has_encrypted_content"] = true
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func hasEncryptedFeedContent(post *models.Post) bool {
	if post == nil || post.ArticleHTML == "" {
		return false
	}
	if value, ok := post.Extra["has_encrypted_content"].(bool); ok {
		return value
	}
	// This fallback supports callers that provide the serialized encryption
	// wrapper directly instead of running the encryption plugin first.
	_, ok := parseEncryptedFeedWrapper(post.ArticleHTML)
	return ok
}
