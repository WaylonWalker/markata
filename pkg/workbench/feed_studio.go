package workbench

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/WaylonWalker/markata/pkg/config"
	"github.com/WaylonWalker/markata/pkg/filter"
	"github.com/WaylonWalker/markata/pkg/lsp"
	"github.com/WaylonWalker/markata/pkg/models"
)

// Feed Studio uses the same configured filters as the feed publisher. Its
// previews are membership previews, not an approximation of template output.
type feedStudioItem struct {
	Path  string `json:"path"`
	Slug  string `json:"slug"`
	Title string `json:"title"`
	Date  string `json:"date,omitempty"`
	Image string `json:"image,omitempty"`
}

type feedStudioFeed struct {
	Slug         string `json:"slug"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Filter       string `json:"filter"`
	Sort         string `json:"sort"`
	Reverse      bool   `json:"reverse"`
	Presentation string `json:"presentation"`
	Template     string `json:"template"`
	Count        int    `json:"count"`
	Editable     bool   `json:"editable"`
}

type feedStudioUpdate struct {
	Slug            string `json:"slug"`
	ExpectedVersion string `json:"expected_version"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	Filter          string `json:"filter"`
	Sort            string `json:"sort"`
	Reverse         bool   `json:"reverse"`
	Presentation    string `json:"presentation"`
	Template        string `json:"template"`
}

func feedStudioFile(root string) (string, error) {
	for _, name := range []string{"markata-go.toml", "markata-go.yaml", "markata-go.yml", "markata-go.json"} {
		filename := filepath.Join(root, name)
		info, err := os.Lstat(filename)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return "", err
		}
		if !info.Mode().IsRegular() {
			return "", fmt.Errorf("site config must be a regular file")
		}
		return filename, nil
	}
	return "", os.ErrNotExist
}

// feedStudioIndex gives filter evaluation the same metadata, slug and
// workspace-relative path values used by the existing Feeds navigation.
func (h *Handler) feedStudioIndex() []*models.Post {
	var index *lsp.Index
	if snapshot, ok := h.workspace.(interface{ Snapshot() (*lsp.Index, uint64) }); ok {
		index, _ = snapshot.Snapshot()
	}
	if index == nil {
		return nil
	}
	posts := make([]*models.Post, 0)
	for _, info := range index.AllPosts() {
		var post models.Post
		data, err := json.Marshal(info.Metadata)
		if err != nil || json.Unmarshal(data, &post) != nil {
			continue
		}
		post.Path, post.Slug = workspaceRelativePath(h.root, info.Path), info.Slug
		post.Tags = append([]string(nil), getMetadataStrings(info.Metadata["tags"])...)
		title := info.Title
		post.Title = &title
		if !post.Published || post.Private || post.Draft {
			continue
		}
		posts = append(posts, &post)
	}
	return posts
}

func feedStudioMatch(posts []*models.Post, expression string) ([]*models.Post, error) {
	if strings.TrimSpace(expression) == "" {
		return posts, nil
	}
	parsed, err := filter.Parse(expression)
	if err != nil {
		return nil, err
	}
	matches := make([]*models.Post, 0)
	for _, post := range posts {
		ok, err := parsed.Match(post)
		if err != nil {
			return nil, err
		}
		if ok {
			matches = append(matches, post)
		}
	}
	return matches, nil
}

func feedStudioSort(posts []*models.Post, sortField string, reverse bool) {
	sort.SliceStable(posts, func(i, j int) bool {
		left, right := posts[i], posts[j]
		switch sortField {
		case "title":
			a, b := "", ""
			if left.Title != nil {
				a = *left.Title
			}
			if right.Title != nil {
				b = *right.Title
			}
			if a != b {
				return strings.ToLower(a) < strings.ToLower(b)
			}
		case "slug":
			if left.Slug != right.Slug {
				return left.Slug < right.Slug
			}
		case "path":
			return left.Path < right.Path
		case "guide_order":
			la, lok := left.Extra["guide_order"]
			ra, rok := right.Extra["guide_order"]
			if lok && rok && fmt.Sprint(la) != fmt.Sprint(ra) {
				return fmt.Sprint(la) < fmt.Sprint(ra)
			}
		case "date", "":
			if left.Date != nil && right.Date != nil && !left.Date.Equal(*right.Date) {
				return left.Date.Before(*right.Date)
			}
		}
		return left.Path < right.Path
	})
	if reverse {
		for i, j := 0, len(posts)-1; i < j; i, j = i+1, j-1 {
			posts[i], posts[j] = posts[j], posts[i]
		}
	}
}

func feedStudioVersion(source []byte) string {
	sum := sha256.Sum256(source)
	return hex.EncodeToString(sum[:])
}

func feedStudioConfigVersion(root string) (string, error) {
	filename, err := feedStudioFile(root)
	if err != nil {
		return "", err
	}
	source, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return feedStudioVersion(source), nil
}

func (h *Handler) serveFeedStudio(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		h.feedStudioRead(w, r)
	case http.MethodPost:
		if !sameOriginWorkspaceMutation(r) {
			fail(w, http.StatusForbidden, "same-origin feed changes are required")
			return
		}
		h.feedStudioSave(w, r)
	default:
		w.Header().Set("Allow", "GET, POST")
		fail(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) feedStudioRead(w http.ResponseWriter, r *http.Request) {
	filename, err := feedStudioFile(h.root)
	if errors.Is(err, os.ErrNotExist) {
		_ = json.NewEncoder(w).Encode(map[string]any{"feeds": []feedStudioFeed{}, "posts": []feedStudioItem{}})
		return
	}
	if err != nil {
		fail(w, http.StatusUnprocessableEntity, "cannot read site configuration: "+err.Error())
		return
	}
	var cfg *models.Config
	// The initial feed selection must reflect newly saved configuration even
	// before the asynchronous workspace file watcher catches up.
	if !r.URL.Query().Has("slug") {
		cfg, err = config.Load(filename)
		if err != nil {
			fail(w, http.StatusUnprocessableEntity, "cannot load feed settings: "+err.Error())
			return
		}
	} else {
		cfg = h.currentSiteConfig()
		if cfg == nil {
			cfg, err = config.Load(filename)
			if err != nil {
				fail(w, http.StatusUnprocessableEntity, "cannot load feed settings: "+err.Error())
				return
			}
		}
	}
	version, err := feedStudioConfigVersion(h.root)
	if err != nil {
		fail(w, http.StatusUnprocessableEntity, "cannot read site configuration: "+err.Error())
		return
	}
	posts := h.feedStudioIndex()
	feeds := make([]feedStudioFeed, 0, len(cfg.Feeds))
	selectedSlug := r.URL.Query().Get("slug")
	selectedFilter, filterOverride := r.URL.Query()["filter"]
	selectedSort := r.URL.Query().Get("sort")
	var selected *models.FeedConfig
	for i := range cfg.Feeds {
		feed := &cfg.Feeds[i]
		matches, matchErr := feedStudioMatch(posts, feed.Filter)
		count := len(matches)
		if matchErr != nil {
			count = 0
		}
		editable := feedStudioEditable(h.root, i, feed.Slug)
		feeds = append(feeds, feedStudioFeed{
			Slug: feed.Slug, Title: feed.Title, Description: feed.Description,
			Filter: feed.Filter, Sort: feed.Sort, Reverse: feed.Reverse,
			Presentation: feed.Presentation, Template: feed.Templates.HTML,
			Count: count, Editable: editable,
		})
		if feed.Slug == selectedSlug {
			selected = feed
		}
	}
	if selected == nil {
		if err := json.NewEncoder(w).Encode(map[string]any{"feeds": feeds, "posts": []feedStudioItem{}, "version": version}); err != nil {
			return
		}
		return
	}
	expression := selected.Filter
	if filterOverride {
		expression = selectedFilter[0]
	}
	matches, err := feedStudioMatch(posts, expression)
	if err != nil {
		if err := json.NewEncoder(w).Encode(map[string]any{"feeds": feeds, "posts": []feedStudioItem{}, "version": version, "filter_error": err.Error()}); err != nil {
			return
		}
		return
	}
	sortField := selected.Sort
	if r.URL.Query().Has("sort") {
		sortField = selectedSort
	}
	reverse := selected.Reverse
	if r.URL.Query().Has("reverse") {
		reverse = r.URL.Query().Get("reverse") == "true"
	}
	if sortField == "" {
		if selected.Type == models.FeedTypeGuide {
			sortField = "guide_order"
			reverse = false
		} else {
			sortField = "date"
			reverse = true
		}
	}
	feedStudioSort(matches, sortField, reverse)
	matchedCount := len(matches)
	offset := selected.Offset
	if offset > len(matches) {
		offset = len(matches)
	}
	matches = matches[offset:]
	if selected.Limit > 0 && len(matches) > selected.Limit {
		matches = matches[:selected.Limit]
	}
	visibleCount := len(matches)
	if len(matches) > 60 {
		matches = matches[:60]
	}
	items := make([]feedStudioItem, 0, len(matches))
	for _, post := range matches {
		title := post.Slug
		if post.Title != nil && *post.Title != "" {
			title = *post.Title
		}
		item := feedStudioItem{Path: post.Path, Slug: post.Slug, Title: title}
		if post.Date != nil {
			item.Date = post.Date.Format("2006-01-02")
		}
		items = append(items, item)
	}
	if err := json.NewEncoder(w).Encode(map[string]any{
		"feeds": feeds, "posts": items, "version": version,
		"rendered_available": feedStudioSnapshotAvailable(h.root, cfg, selected.Slug),
		"rendered_url": feedSnapshotPrefix + strings.Trim(selected.Slug, "/") + "/",
		"matched_count": matchedCount, "visible_count": visibleCount,
		"partial": len(items) < visibleCount,
	}); err != nil {
		return
	}
}

func feedStudioEditable(root string, index int, slug string) bool {
	filename, err := feedStudioFile(root)
	if err != nil {
		return false
	}
	value, err := config.GetValueFromFile(filename, fmt.Sprintf("feeds[%d].slug", index))
	return err == nil && value == slug
}

func feedStudioValidate(update feedStudioUpdate) error {
	if len(update.Title) > 200 {
		return fmt.Errorf("feed title exceeds 200 characters")
	}
	if len(update.Description) > 4000 || len(update.Filter) > 4000 || len(update.Sort) > 100 || len(update.Template) > 255 {
		return fmt.Errorf("feed settings exceed size limits")
	}
	if strings.ContainsAny(update.Template, "\r\n") || strings.ContainsAny(update.Presentation, "\r\n") || len(update.Presentation) > 40 {
		return fmt.Errorf("invalid feed presentation or template")
	}
	if _, err := feedStudioMatch(nil, update.Filter); err != nil {
		return err
	}
	return nil
}

func (h *Handler) feedStudioSave(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 16<<10)
	defer r.Body.Close()
	var update feedStudioUpdate
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&update); err != nil {
		fail(w, http.StatusBadRequest, "invalid feed settings payload: "+err.Error())
		return
	}
	if err := feedStudioValidate(update); err != nil {
		fail(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	if len(update.ExpectedVersion) != 64 {
		fail(w, http.StatusPreconditionRequired, "expected configuration version is required")
		return
	}
	h.workspaceWriteMu.Lock()
	defer h.workspaceWriteMu.Unlock()
	filename, err := feedStudioFile(h.root)
	if err != nil {
		fail(w, http.StatusUnprocessableEntity, "cannot locate site configuration")
		return
	}
	info, err := os.Lstat(filename)
	if err != nil || !info.Mode().IsRegular() {
		fail(w, http.StatusUnprocessableEntity, "site configuration is not a regular file")
		return
	}
	original, err := os.ReadFile(filename)
	if err != nil {
		fail(w, http.StatusUnprocessableEntity, "cannot read site configuration")
		return
	}
	if feedStudioVersion(original) != update.ExpectedVersion {
		fail(w, http.StatusConflict, "site configuration changed on disk; reload and compare before saving")
		return
	}
	cfg := h.currentSiteConfig()
	index := -1
	if cfg != nil {
		for i := range cfg.Feeds {
			if cfg.Feeds[i].Slug == update.Slug {
				index = i
				break
			}
		}
	}
	if index < 0 || !feedStudioEditable(h.root, index, update.Slug) {
		fail(w, http.StatusUnprocessableEntity, "feed is not directly editable in the site config (it may come from an included file)")
		return
	}
	temporary, err := os.CreateTemp(h.root, ".markata-feed-*"+filepath.Ext(filename))
	if err != nil {
		fail(w, http.StatusUnprocessableEntity, "cannot create configuration staging file")
		return
	}
	tempPath := temporary.Name()
	defer os.Remove(tempPath)
	err = temporary.Chmod(info.Mode().Perm())
	if err == nil {
		_, err = temporary.Write(original)
	}
	if err == nil {
		err = temporary.Sync()
	}
	if closeErr := temporary.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		fail(w, http.StatusUnprocessableEntity, "cannot stage configuration: "+err.Error())
		return
	}
	settings := []struct {
		key string
		value any
	}{
		{"title", update.Title}, {"description", update.Description}, {"filter", update.Filter},
		{"sort", update.Sort}, {"reverse", update.Reverse},
		{"presentation", update.Presentation}, {"templates.html", update.Template},
	}
	for _, setting := range settings {
		key := fmt.Sprintf("feeds[%d].%s", index, setting.key)
		if err := config.SetValueInFile(tempPath, key, setting.value); err != nil {
			fail(w, http.StatusUnprocessableEntity, "cannot update "+setting.key+": "+err.Error())
			return
		}
	}
	parsed, err := config.Load(tempPath)
	if err != nil {
		fail(w, http.StatusUnprocessableEntity, "invalid site configuration after edit: "+err.Error())
		return
	}
	if len(parsed.Feeds) <= index || parsed.Feeds[index].Slug != update.Slug {
		fail(w, http.StatusUnprocessableEntity, "feed identity changed during configuration update")
		return
	}
	if err := os.Rename(tempPath, filename); err != nil {
		fail(w, http.StatusUnprocessableEntity, "cannot save site configuration: "+err.Error())
		return
	}
	saved, err := os.ReadFile(filename)
	if err != nil {
		fail(w, http.StatusInternalServerError, "could not read saved configuration")
		return
	}
	if err := json.NewEncoder(w).Encode(map[string]string{
		"version": feedStudioVersion(saved), "status": "saved", "file": filepath.Base(filename),
	}); err != nil {
		return
	}
}


func feedStudioSnapshotAvailable(root string, cfg *models.Config, slug string) bool {
	_, err := feedSnapshotFile(root, cfg, strings.Trim(slug, "/")+"/index.html")
	return err == nil
}
