package plugins

import (
	"sort"
	"strings"
	"time"

	"github.com/WaylonWalker/markata-go/pkg/models"
	"github.com/WaylonWalker/markata-go/pkg/templates"
)

const (
	presentationClassic      = "classic"
	presentationEditorial    = "editorial"
	presentationGallery      = "gallery"
	presentationReading      = "reading"
	presentationTimeline     = "timeline"
	presentationResources    = "resources"
	presentationAuto         = "auto"
	presentationFeatured     = "featured"
	presentationOther        = "other"
	presentationUndated      = "undated"
	presentationQuiet        = "quiet"
	presentationFeedTemplate = "feed.html"
)

// EditorialFeedItem is a post or a collection of same-day shots. Grouped shots
// retain individual permalinks and render through the normal card router.
type EditorialFeedItem struct {
	Kind       string
	Importance string
	TypeLabel  string
	DateLabel  string
	Lead       bool
	Mosaic     bool
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
		return presentationOther
	}
	tpl := strings.ToLower(strings.TrimSuffix(post.Template, ".html"))
	path := strings.ToLower(post.Path)
	switch {
	case strings.Contains(path, "/shots/"), tpl == templateTypePhoto, tpl == "shot":
		return templateTypePhoto
	case strings.Contains(path, "/ping/"), strings.Contains(path, "/thoughts/"),
		tpl == "ping", tpl == templateTypeNote, tpl == "inline", tpl == "thought":
		return templateTypeNote
	case strings.Contains(path, "/stars/"), strings.Contains(path, "/link/"),
		tpl == templateTypeLink, tpl == "bookmark":
		return templateTypeLink
	case tpl == templateTypeVideo, strings.Contains(path, "/video/"):
		return templateTypeVideo
	case strings.Contains(path, "/blog/"), strings.Contains(path, "/til/"),
		tpl == templateTypePost, tpl == templateTypeArticle, tpl == "guide", tpl == templateTypeTutorial:
		return templateTypeArticle
	default:
		return presentationOther
	}
}

// Auto mode uses feed-wide distribution for stable page-to-page results.
// The default is classic, so older sites and custom templates retain their UI.
func resolveFeedPresentation(fc *models.FeedConfig) string {
	if fc == nil {
		return presentationClassic
	}
	switch fc.Presentation {
	case presentationClassic, presentationEditorial, presentationGallery, presentationReading, presentationTimeline, presentationResources:
		return fc.Presentation
	case presentationAuto:
		// Custom layout templates, such as feed-photo-grid.html, own rendering.
		if fc.Templates.HTML != "" && fc.Templates.HTML != presentationFeedTemplate {
			return presentationClassic
		}
	default:
		return presentationClassic
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
		return presentationClassic
	}
	for _, option := range []struct{ kind, view string }{
		{templateTypePhoto, presentationGallery},
		{templateTypeNote, presentationTimeline},
		{templateTypeLink, presentationResources},
		{templateTypeArticle, presentationReading},
	} {
		if counts[option.kind]*100 >= total*80 {
			return option.view
		}
	}
	distinct := 0
	for _, kind := range []string{templateTypePhoto, templateTypeNote, templateTypeLink, templateTypeVideo, templateTypeArticle} {
		if counts[kind] >= 2 {
			distinct++
		}
	}
	if total >= 12 && distinct >= 3 {
		return presentationEditorial
	}
	return presentationClassic
}

func editorialRank(post *models.Post) int {
	if post != nil && post.Extra != nil {
		if featured, ok := post.Extra[presentationFeatured].(bool); ok && featured {
			return 0
		}
	}
	switch presentationKind(post) {
	case templateTypeArticle:
		return 1
	case templateTypeVideo, templateTypeLink:
		return 2
	case templateTypePhoto:
		return 3
	case templateTypeNote:
		return 4
	default:
		return 3
	}
}

func editorialWeek(post *models.Post) (key, label string) {
	if post == nil || post.Date == nil {
		return presentationUndated, strings.ToUpper(presentationUndated[:1]) + presentationUndated[1:]
	}
	date := *post.Date
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	start = start.AddDate(0, 0, -((int(start.Weekday()) + 6) % 7))
	return start.Format("2006-01-02"), start.Format("Jan 2") + " – " + start.AddDate(0, 0, 6).Format("Jan 2, 2006")
}

// editorialTypeLabel is display-only: shared card routing and syndication stay unchanged.
func editorialTypeLabel(post *models.Post) string {
	switch presentationKind(post) {
	case templateTypeArticle:
		return "Writing"
	case templateTypeVideo:
		return strings.ToUpper(templateTypeVideo[:1]) + templateTypeVideo[1:]
	case templateTypeLink:
		return "Link"
	case templateTypePhoto:
		return "Photograph"
	case templateTypeNote:
		return "Note"
	default:
		return "From the archive"
	}
}

func editorialPhotoDate(post *models.Post) string {
	if post == nil || post.Date == nil {
		return strings.ToUpper(presentationUndated[:1]) + presentationUndated[1:]
	}
	return post.Date.Format("January 2, 2006")
}

func editorialDay(post *models.Post) string {
	if post == nil || post.Date == nil {
		return presentationUndated
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
			if presentationKind(post) == templateTypePhoto {
				day := editorialDay(post)
				photos[day] = append(photos[day], post)
			}
		}
		emitted := make(map[string]bool)
		for position, post := range postsInWeek {
			if presentationKind(post) == templateTypePhoto {
				day := editorialDay(post)
				if len(photos[day]) > 1 {
					if emitted[day] {
						continue
					}
					emitted[day] = true
					sections[i].Items = append(sections[i].Items, EditorialFeedItem{
						Kind: embedOptionShots, Importance: presentationGallery,
						Posts: templates.PostsToMaps(photos[day]), rank: 3, order: position,
						DateLabel: editorialPhotoDate(post), Mosaic: len(photos[day]) == 3,
					})
					continue
				}
			}

			rank := editorialRank(post)
			importance := "standard"
			if rank <= 1 {
				importance = presentationFeatured
			} else if rank == 4 {
				importance = presentationQuiet
			}
			sections[i].Items = append(sections[i].Items, EditorialFeedItem{
				Kind: "post", Importance: importance, TypeLabel: editorialTypeLabel(post),
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
		// Every week gets at most one leading article/link/video. Short
		// notes never become artificial lead stories; explicit featured
		// content already sorts before ordinary articles and links.
		for j := range sections[i].Items {
			item := &sections[i].Items[j]
			if item.Kind == "post" && item.rank <= 2 {
				item.Lead = true
				break
			}
		}
	}
	return sections
}
