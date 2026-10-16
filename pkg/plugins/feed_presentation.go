package plugins

import (
	"sort"
	"strings"
	"time"

	"github.com/WaylonWalker/markata-go/pkg/models"
	"github.com/WaylonWalker/markata-go/pkg/templates"
)

// EditorialFeedItem is a post or a collection of same-day shots. Grouped shots
// retain individual permalinks and render through the normal card router.
type EditorialFeedItem struct {
	Kind       string
	Importance string
	Posts      []map[string]interface{}
	order      int
	rank       int
}

type EditorialFeedSection struct {
	Label string
	Items []EditorialFeedItem
}

func presentationKind(post *models.Post) string {
	if post == nil {
		return "other"
	}
	tpl := strings.ToLower(strings.TrimSuffix(post.Template, ".html"))
	path := strings.ToLower(post.Path)
	switch {
	case strings.Contains(path, "/shots/"), tpl == "photo", tpl == "shot":
		return "photo"
	case strings.Contains(path, "/ping/"), tpl == "ping", tpl == "note", tpl == "inline":
		return "note"
	case strings.Contains(path, "/stars/"), strings.Contains(path, "/link/"), strings.Contains(path, "/thoughts/"),
		tpl == "link", tpl == "bookmark", tpl == "thought":
		return "link"
	case tpl == "video", strings.Contains(path, "/video/"):
		return "video"
	case strings.Contains(path, "/blog/"), strings.Contains(path, "/til/"),
		tpl == "post", tpl == "article", tpl == "guide", tpl == "tutorial":
		return "article"
	default:
		return "other"
	}
}

// Auto mode uses feed-wide distribution for stable page-to-page results.
// The default is classic, so older sites and custom templates retain their UI.
func resolveFeedPresentation(fc *models.FeedConfig) string {
	if fc == nil {
		return "classic"
	}
	switch fc.Presentation {
	case "classic", "editorial", "gallery", "reading", "timeline", "resources":
		return fc.Presentation
	case "auto":
		// Custom layout templates, such as feed-photo-grid.html, own rendering.
		if fc.Templates.HTML != "" && fc.Templates.HTML != "feed.html" {
			return "classic"
		}
	default:
		return "classic"
	}

	counts := make(map[string]int)
	total := 0
	for _, post := range fc.Posts {
		if post == nil || post.Skip || post.Draft || (post.Private && !fc.IncludesPrivate()) {
			continue
		}
		counts[presentationKind(post)]++
		total++
	}
	if total < 8 {
		return "classic"
	}
	for _, option := range []struct{ kind, view string }{
		{"photo", "gallery"},
		{"note", "timeline"},
		{"link", "resources"},
		{"article", "reading"},
	} {
		if counts[option.kind]*100 >= total*80 {
			return option.view
		}
	}
	distinct := 0
	for _, kind := range []string{"photo", "note", "link", "video", "article"} {
		if counts[kind] >= 2 {
			distinct++
		}
	}
	if total >= 12 && distinct >= 3 {
		return "editorial"
	}
	return "classic"
}

func editorialRank(post *models.Post) int {
	if post != nil && post.Extra != nil {
		if featured, ok := post.Extra["featured"].(bool); ok && featured {
			return 0
		}
	}
	switch presentationKind(post) {
	case "article":
		return 1
	case "video", "link":
		return 2
	case "photo":
		return 3
	case "note":
		return 4
	default:
		return 3
	}
}

func editorialWeek(post *models.Post) (key, label string) {
	if post == nil || post.Date == nil {
		return "undated", "Undated"
	}
	date := *post.Date
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	start = start.AddDate(0, 0, -((int(start.Weekday()) + 6) % 7))
	return start.Format("2006-01-02"), start.Format("Jan 2") + " – " + start.AddDate(0, 0, 6).Format("Jan 2, 2006")
}

func editorialDay(post *models.Post) string {
	if post == nil || post.Date == nil {
		return "undated"
	}
	return post.Date.Format("2006-01-02")
}

// Only regroup the visible HTML page. Publishing metadata, syndication,
// Simple, Calendar, and existing pagination are never reordered or replaced.
func buildEditorialFeedSections(posts []*models.Post) []EditorialFeedSection {
	sections := make([]EditorialFeedSection, 0)
	index := make(map[string]int)
	weeks := make([][]*models.Post, 0)

	for _, post := range posts {
		if post == nil {
			continue
		}
		key, label := editorialWeek(post)
		i, exists := index[key]
		if !exists {
			i = len(sections)
			index[key] = i
			sections = append(sections, EditorialFeedSection{Label: label})
			weeks = append(weeks, nil)
		}
		weeks[i] = append(weeks[i], post)
	}

	for i, postsInWeek := range weeks {
		photos := make(map[string][]*models.Post)
		for _, post := range postsInWeek {
			if presentationKind(post) == "photo" {
				day := editorialDay(post)
				photos[day] = append(photos[day], post)
			}
		}
		emitted := make(map[string]bool)
		for position, post := range postsInWeek {
			if presentationKind(post) == "photo" {
				day := editorialDay(post)
				if len(photos[day]) > 1 {
					if emitted[day] {
						continue
					}
					emitted[day] = true
					sections[i].Items = append(sections[i].Items, EditorialFeedItem{
						Kind: "shots", Importance: "gallery",
						Posts: templates.PostsToMaps(photos[day]), rank: 3, order: position,
					})
					continue
				}
			}

			rank := editorialRank(post)
			importance := "standard"
			if rank <= 1 {
				importance = "featured"
			} else if rank == 4 {
				importance = "quiet"
			}
			sections[i].Items = append(sections[i].Items, EditorialFeedItem{
				Kind: "post", Importance: importance,
				Posts: templates.PostsToMaps([]*models.Post{post}), rank: rank, order: position,
			})
		}

		sort.SliceStable(sections[i].Items, func(a, b int) bool {
			left, right := sections[i].Items[a], sections[i].Items[b]
			if left.rank != right.rank {
				return left.rank < right.rank
			}
			return left.order < right.order
		})
	}
	return sections
}
