package cmd

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/PuerkitoBio/goquery"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/plugins"
)

// The subscription editor is serve-only; it is never copied to the output site.
//
//go:embed devui/subscriptions.html
var subscriptionPage []byte

const (
	serveSubscriptionsPage   = "/_markata/subscriptions"
	serveSubscriptionsAPI    = "/_markata/subscriptions/api"
	serveSubscriptionsFind   = "/_markata/subscriptions/discover"
	subscriptionActionAdd    = "add"
	subscriptionActionUpdate = "update"
	subscriptionActionRemove = "remove"
	subscriptionTOMLExt      = ".toml"
	subscriptionHTTPS        = "https"
	subscriptionHTTP         = "http"
)

type subscriptionView struct {
	models.ExternalFeedConfig
	Source   string `json:"source"`
	Editable bool   `json:"editable"`
}

type subscriptionResponse struct {
	Feeds   []subscriptionView `json:"feeds,omitempty"`
	Target  string             `json:"target,omitempty"`
	Enabled bool               `json:"enabled"`
	Error   string             `json:"error,omitempty"`
	Message string             `json:"message,omitempty"`
}

type subscriptionRequest struct {
	Action      string `json:"action"`
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Type        string `json:"type"`
	Tags        string `json:"tags"`
	Handle      string `json:"handle"`
	Active      bool   `json:"active"`
}

type discoveredSubscription struct {
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	SiteURL     string `json:"site_url"`
	Error       string `json:"error,omitempty"`
}

func subscriptionJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		warnf("subscription response write failed: %v", err)
	}
}

func handleSubscriptionsPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet || !sameOriginRequest(r) {
		http.Error(w, "local development access only", http.StatusForbidden)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if _, err := w.Write(subscriptionPage); err != nil {
		warnf("subscription page write failed: %v", err)
	}
}

func subscriptionReadRequest(w http.ResponseWriter, r *http.Request, dst any) bool {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		subscriptionJSON(w, http.StatusUnsupportedMediaType, subscriptionResponse{Error: "expected application/json"})
		return false
	}
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 16<<10))
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		subscriptionJSON(w, http.StatusBadRequest, subscriptionResponse{Error: err.Error()})
		return false
	}
	return true
}

// blogrollRawFeeds reads the *local* entries in each source, before composition.
// Editing an included file rather than the root is crucial to avoid misplaced changes.
func blogrollRawFeeds(raw map[string]any) []map[string]any {
	root := raw
	if nested, ok := root["markata-go"].(map[string]any); ok {
		root = nested
	}
	blogroll, ok := root["blogroll"].(map[string]any)
	if !ok {
		return nil
	}
	entries, ok := blogroll["feeds"].([]any)
	if !ok {
		return nil
	}
	result := make([]map[string]any, 0, len(entries))
	for _, entry := range entries {
		if feed, ok := entry.(map[string]any); ok {
			result = append(result, feed)
		}
	}
	return result
}

func feedOwner(state *settingsState, address string) (string, error) {
	// Arrays replace across config layers. The last source defining feeds
	// owns the effective array; never edit a shadowed earlier layer.
	owner := ""
	matches := 0
	for i, raw := range state.raws {
		entries := blogrollRawFeeds(raw)
		if entries == nil {
			continue
		}
		owner = state.sources[i]
		matches = 0
		for _, entry := range entries {
			if value, ok := entry["url"].(string); ok && value == address {
				matches++
			}
		}
	}
	if matches != 1 {
		return "", fmt.Errorf("could not uniquely identify feed %q in the effective config file (found %d)", address, matches)
	}
	if state.kinds[owner] != "" {
		return "", fmt.Errorf("feed belongs to a %s config (%s); edit that source manually", state.kinds[owner], displayConfigPath(owner))
	}
	if strings.ToLower(filepath.Ext(owner)) != subscriptionTOMLExt {
		return "", fmt.Errorf("editing %s is not yet supported in the browser; use the config file directly", filepath.Ext(owner))
	}
	return owner, nil
}

func newFeedTarget(state *settingsState) (string, error) {
	// If a later included/override file defines feeds, the earlier lists are
	// shadowed. Choose the effective definition, not the first blogroll table.
	var owner string
	for i, raw := range state.raws {
		if blogrollRawFeeds(raw) != nil {
			owner = state.sources[i]
		}
	}
	if owner == "" {
		target, err := state.target([]string{"blogroll", "feeds"})
		if err != nil {
			return "", err
		}
		owner = target
	}
	if state.kinds[owner] != "" {
		return "", fmt.Errorf("blogroll lives in a %s config (%s); choose an editable site config", state.kinds[owner], displayConfigPath(owner))
	}
	if strings.ToLower(filepath.Ext(owner)) != subscriptionTOMLExt {
		return "", fmt.Errorf("browser subscriptions require an owning TOML file, found %s", displayConfigPath(owner))
	}
	return owner, nil
}

func describeSubscriptions(state *settingsState) subscriptionResponse {
	resp := subscriptionResponse{Enabled: state.cfg.Blogroll.Enabled}
	if target, err := newFeedTarget(state); err == nil {
		resp.Target = displayConfigPath(target)
	} else {
		resp.Error = err.Error()
	}
	for i := range state.cfg.Blogroll.Feeds {
		item := &state.cfg.Blogroll.Feeds[i]
		src, err := feedOwner(state, item.URL)
		view := subscriptionView{ExternalFeedConfig: *item, Editable: err == nil}
		if err == nil {
			view.Source = displayConfigPath(src)
		} else {
			view.Source = "Read only: " + err.Error()
		}
		resp.Feeds = append(resp.Feeds, view)
	}
	return resp
}

func handleSubscriptionsAPI(w http.ResponseWriter, r *http.Request) {
	if !sameOriginRequest(r) {
		subscriptionJSON(w, http.StatusForbidden, subscriptionResponse{Error: "local same-origin requests only"})
		return
	}
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		w.Header().Set("Allow", "GET, POST")
		subscriptionJSON(w, http.StatusMethodNotAllowed, subscriptionResponse{Error: "method not allowed"})
		return
	}
	serveConfigMu.Lock()
	defer serveConfigMu.Unlock()
	state, err := loadSettingsState()
	if err != nil {
		subscriptionJSON(w, http.StatusInternalServerError, subscriptionResponse{Error: err.Error()})
		return
	}
	if r.Method == http.MethodGet {
		subscriptionJSON(w, http.StatusOK, describeSubscriptions(state))
		return
	}
	var request subscriptionRequest
	if !subscriptionReadRequest(w, r, &request) {
		return
	}
	if err := changeSubscription(state, request); err != nil {
		subscriptionJSON(w, http.StatusUnprocessableEntity, subscriptionResponse{Error: err.Error()})
		return
	}
	state, err = loadSettingsState()
	if err != nil {
		subscriptionJSON(w, http.StatusInternalServerError, subscriptionResponse{Error: err.Error()})
		return
	}
	resp := describeSubscriptions(state)
	resp.Message = "Saved subscription; updating your local reader."
	subscriptionJSON(w, http.StatusOK, resp)
}

var feedHeader = regexp.MustCompile(`(?m)^[ \t]*\[\[(?:markata-go\.)?blogroll\.feeds\]\][ \t]*(?:\r?\n|$)`)
var anyTableHeader = regexp.MustCompile(`(?m)^[ \t]*\[`)
var feedURLLine = regexp.MustCompile(`(?m)^[ \t]*url[ \t]*=[ \t]*(?:"([^"]+)"|'([^']+)')[ \t]*(?:#.*)?$`)
var blogrollTable = regexp.MustCompile(`(?m)^[ \t]*\[(?:markata-go\.)?blogroll\][ \t]*(?:\r?\n|$)`)

type subscriptionBlock struct {
	start, end int
	url        string
}

func subscriptionBlocks(content string) []subscriptionBlock {
	headers := feedHeader.FindAllStringIndex(content, -1)
	result := make([]subscriptionBlock, 0, len(headers))
	for _, header := range headers {
		end := len(content)
		if following := anyTableHeader.FindStringIndex(content[header[1]:]); following != nil {
			end = header[1] + following[0]
		}
		url := ""
		if match := feedURLLine.FindStringSubmatch(content[header[1]:end]); match != nil {
			url = match[1]
			if url == "" {
				url = match[2]
			}
		}
		result = append(result, subscriptionBlock{header[0], end, url})
	}
	return result
}

func uniqueBlock(content, address string) (subscriptionBlock, error) {
	var selected subscriptionBlock
	found := 0
	for _, block := range subscriptionBlocks(content) {
		if block.url == address {
			selected = block
			found++
		}
	}
	if found != 1 {
		return selected, fmt.Errorf("cannot safely identify subscription %q in its TOML file (%d matching blocks)", address, found)
	}
	return selected, nil
}

func subscriptionField(body, name, value string) string {
	pattern := regexp.MustCompile(`(?m)^[ \t]*` + regexp.QuoteMeta(name) + `[ \t]*=[^\r\n]*`)
	replacement := name + " = " + value
	if pattern.MatchString(body) {
		return pattern.ReplaceAllStringFunc(body, func(string) string { return replacement })
	}
	return strings.TrimRight(body, "\r\n") + "\n" + replacement + "\n"
}

func subscriptionFields(req subscriptionRequest) (map[string]string, error) {
	if len(req.URL) > 2048 || len(req.Title) > 200 || len(req.Category) > 100 ||
		len(req.Description) > 1000 || len(req.Handle) > 100 || len(req.Tags) > 600 {
		return nil, errors.New("subscription fields exceed their maximum length")
	}
	if req.Type != "written" && req.Type != "video" && req.Type != "podcast" {
		return nil, errors.New("feed type must be written, video, or podcast")
	}
	return map[string]string{
		"title":       strconv.Quote(strings.TrimSpace(req.Title)),
		"description": strconv.Quote(strings.TrimSpace(req.Description)),
		"category":    strconv.Quote(strings.TrimSpace(req.Category)),
		"type":        strconv.Quote(req.Type),
		"tags":        "[" + formatTOMLStringArray(parseTagsBlogroll(req.Tags)) + "]",
		"handle":      strconv.Quote(strings.TrimSpace(req.Handle)),
		"active":      strconv.FormatBool(req.Active),
	}, nil
}

func enableSubscriptions(content string) string {
	if loc := blogrollTable.FindStringIndex(content); loc != nil {
		end := len(content)
		if next := anyTableHeader.FindStringIndex(content[loc[1]:]); next != nil {
			end = loc[1] + next[0]
		}
		body := subscriptionField(content[loc[1]:end], "enabled", "true")
		return content[:loc[1]] + body + content[end:]
	}
	return strings.TrimRight(content, "\r\n") + "\n\n[markata-go.blogroll]\nenabled = true\n"
}

//nolint:gocyclo // A single atomic operation validates, edits, verifies, and rolls back the config.
func changeSubscription(state *settingsState, req subscriptionRequest) error {
	if req.Action != subscriptionActionAdd && req.Action != subscriptionActionUpdate && req.Action != subscriptionActionRemove {
		return errors.New("invalid subscription action")
	}
	u, err := url.Parse(req.URL)
	if err != nil || u.Scheme != subscriptionHTTPS && u.Scheme != subscriptionHTTP || u.Host == "" || u.User != nil {
		return errors.New("expected an absolute HTTP(S) feed URL")
	}
	var target string
	if req.Action == subscriptionActionAdd {
		target, err = newFeedTarget(state)
		if err == nil {
			for i := range state.cfg.Blogroll.Feeds {
				existing := &state.cfg.Blogroll.Feeds[i]
				if existing.URL == req.URL {
					return fmt.Errorf("already subscribed to %s", req.URL)
				}
			}
		}
	} else {
		target, err = feedOwner(state, req.URL)
	}
	if err != nil {
		return err
	}
	fields := map[string]string{}
	if req.Action != subscriptionActionRemove {
		fields, err = subscriptionFields(req)
		if err != nil {
			return err
		}
		if req.Action == subscriptionActionAdd {
			for i := range state.cfg.Blogroll.Feeds {
				f := &state.cfg.Blogroll.Feeds[i]
				if req.Handle != "" && f.Handle == req.Handle {
					return fmt.Errorf("handle %q already belongs to %s", req.Handle, f.Title)
				}
			}
		}
	}
	original, err := os.ReadFile(target)
	if errors.Is(err, os.ErrNotExist) && req.Action == subscriptionActionAdd {
		original = nil
		err = nil
	}
	if err != nil {
		return err
	}
	content := string(original)
	switch req.Action {
	case subscriptionActionAdd:
		if !state.cfg.Blogroll.Enabled {
			content = enableSubscriptions(content)
		}
		prefix := detectBlogrollTablePrefix(content)
		var block strings.Builder
		block.WriteString("\n\n[[" + prefix + ".feeds]]\n")
		block.WriteString("url = " + strconv.Quote(req.URL) + "\n")
		for _, name := range []string{"title", "description", "category", "type", "tags", "handle", "active"} {
			block.WriteString(name + " = " + fields[name] + "\n")
		}
		content = strings.TrimRight(content, "\r\n") + block.String()
	case subscriptionActionUpdate, subscriptionActionRemove:
		selected, findErr := uniqueBlock(content, req.URL)
		if findErr != nil {
			return findErr
		}
		replacement := ""
		if req.Action == subscriptionActionUpdate {
			replacement = content[selected.start:selected.end]
			for _, name := range []string{"title", "description", "category", "type", "tags", "handle", "active"} {
				replacement = subscriptionField(replacement, name, fields[name])
			}
		}
		content = content[:selected.start] + replacement + content[selected.end:]
	}
	if _, err := toml.Decode(content, &map[string]any{}); err != nil {
		return fmt.Errorf("generated TOML was invalid: %w", err)
	}
	// Retain the original on a verification failure. Rebuilds share
	// serveConfigMu so they cannot observe a half-applied transaction.
	stat, statErr := os.Stat(target)
	mode := os.FileMode(0o600)
	if statErr == nil {
		mode = stat.Mode().Perm()
	}
	if err := writeSubscriptionFile(target, []byte(content), mode); err != nil {
		return err
	}
	updated, verifyErr := loadSettingsState()
	if verifyErr == nil {
		var present bool
		for i := range updated.cfg.Blogroll.Feeds {
			f := &updated.cfg.Blogroll.Feeds[i]
			if f.URL == req.URL {
				present = true
				if req.Action == subscriptionActionUpdate && (f.Title != req.Title || f.Category != req.Category || string(f.Type) != req.Type || f.IsActive() != req.Active) {
					verifyErr = errors.New("saved subscription did not take effect in composed configuration")
				}
			}
		}
		if (req.Action == subscriptionActionRemove && present) || (req.Action != subscriptionActionRemove && !present) {
			verifyErr = errors.New("saved subscriptions were shadowed by another config file")
		}
		if req.Action == subscriptionActionAdd && !updated.cfg.Blogroll.Enabled {
			verifyErr = errors.New("blogroll remains disabled")
		}
	}
	if verifyErr != nil {
		var rollbackErr error
		if len(original) == 0 && errors.Is(statErr, os.ErrNotExist) {
			rollbackErr = os.Remove(target)
		} else {
			rollbackErr = writeSubscriptionFile(target, original, mode)
		}
		if rollbackErr != nil {
			return fmt.Errorf("subscription change failed: %w (rollback failed: %w)", verifyErr, rollbackErr)
		}
		return fmt.Errorf("subscription change rolled back: %w", verifyErr)
	}
	// Keep fast serve builds lazy, but fetch the affected source explicitly.
	// The normal blogroll build deliberately avoids remote refresh by default.
	requestSubscriptionRefresh(updated.cfg.Blogroll, req)
	return nil
}

var subscriptionRefreshMu sync.Mutex

func requestSubscriptionRefresh(cfg models.BlogrollConfig, req subscriptionRequest) {
	rebuild := serveRequestFullRebuild
	if rebuild == nil {
		return
	}
	if req.Action == subscriptionActionRemove || !req.Active {
		rebuild()
		return
	}
	var feed *models.ExternalFeedConfig
	for i := range cfg.Feeds {
		if cfg.Feeds[i].URL == req.URL {
			feed = &cfg.Feeds[i]
			break
		}
	}
	if feed == nil {
		rebuild()
		return
	}
	one := *feed
	go func() {
		subscriptionRefreshMu.Lock()
		defer subscriptionRefreshMu.Unlock()
		cfg.Feeds = []models.ExternalFeedConfig{one}
		if _, err := plugins.RefreshBlogrollCache(cfg); err != nil {
			warnf("Subscription refresh failed for %s: %v", req.URL, err)
		}
		rebuild()
	}()
}

func writeSubscriptionFile(path string, data []byte, mode os.FileMode) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".markata-subscriptions-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	defer tmp.Close()
	if err := tmp.Chmod(mode); err != nil {
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		return err
	}
	if err := tmp.Sync(); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}

// Discovery never allows HTTP(S) requests to private/loopback/link-local IPs,
// including redirect targets. DNS is resolved once and the vetted IP is pinned
// to the socket to avoid DNS-rebinding between validation and connection.
func safeSubscriptionClient() *http.Client {
	transport := &http.Transport{DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, err
		}
		if port != "443" && port != "80" {
			return nil, errors.New("only standard HTTP(S) ports are supported")
		}
		resolved, err := net.DefaultResolver.LookupIPAddr(ctx, host)
		if err != nil {
			return nil, err
		}
		for _, item := range resolved {
			ip := item.IP
			if !ip.IsGlobalUnicast() || ip.IsPrivate() || ip.IsLoopback() || ip.IsLinkLocalUnicast() {
				continue
			}
			return (&net.Dialer{Timeout: 5 * time.Second}).DialContext(ctx, network, net.JoinHostPort(ip.String(), port))
		}
		return nil, errors.New("website has no public IP address")
	}}
	return &http.Client{Timeout: 12 * time.Second, Transport: transport, CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) >= 4 {
			return errors.New("too many redirects")
		}
		return validSubscriptionURL(req.URL)
	}}
}

func validSubscriptionURL(u *url.URL) error {
	if u == nil || (u.Scheme != subscriptionHTTP && u.Scheme != subscriptionHTTPS) || u.Hostname() == "" || u.User != nil {
		return errors.New("use a public HTTP(S) website or feed URL")
	}
	if u.Port() != "" && u.Port() != "443" && u.Port() != "80" {
		return errors.New("feed discovery supports ports 80 and 443 only")
	}
	return nil
}

func fetchSubscriptionURL(ctx context.Context, client *http.Client, address string) (body []byte, finalURL string, err error) {
	u, err := url.Parse(address)
	if err != nil {
		return nil, "", err
	}
	if err := validSubscriptionURL(u); err != nil {
		return nil, "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, address, http.NoBody)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Accept", "application/rss+xml,application/atom+xml,application/feed+json,text/html,application/xml")
	res, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, "", fmt.Errorf("website returned %s", res.Status)
	}
	body, err = io.ReadAll(io.LimitReader(res.Body, 1<<20+1))
	if len(body) > 1<<20 {
		return nil, "", errors.New("feed discovery response is too large")
	}
	return body, res.Request.URL.String(), err
}

//nolint:gocyclo // Feed discovery intentionally handles redirects, HTML alternatives, and response validation together.
func discoverSubscription(ctx context.Context, address string) (discoveredSubscription, error) {
	client := safeSubscriptionClient()
	body, finalURL, err := fetchSubscriptionURL(ctx, client, address)
	if err != nil {
		return discoveredSubscription{}, err
	}
	out := discoveredSubscription{URL: finalURL, SiteURL: finalURL}
	first := bytes.TrimSpace(body)
	if bytes.HasPrefix(first, []byte("<?xml")) || bytes.HasPrefix(first, []byte("<rss")) ||
		bytes.HasPrefix(first, []byte("<feed")) || bytes.HasPrefix(first, []byte("{")) {
		return out, nil // Already a feed, not an HTML page.
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return out, fmt.Errorf("parse website: %w", err)
	}
	out.Title = strings.TrimSpace(doc.Find("head title").First().Text())
	out.Description, _ = doc.Find("meta[name='description']").First().Attr("content")
	if parsed, err := url.Parse(finalURL); err == nil && isYouTubeURL(parsed) {
		channelID := extractYouTubeChannelIDFromURL(finalURL)
		if channelID == "" {
			channelID = extractYouTubeChannelIDFromHTML(string(body))
		}
		if channelID != "" {
			out.URL = youtubeFeedBaseURL + channelID
			return out, nil
		}
	}
	base, err := url.Parse(finalURL)
	if err != nil {
		return out, fmt.Errorf("parse final site URL: %w", err)
	}
	var feeds []string
	doc.Find("link[rel~='alternate']").Each(func(_ int, s *goquery.Selection) {
		typ, _ := s.Attr("type")
		if !strings.Contains(typ, "rss+xml") && !strings.Contains(typ, "atom+xml") && !strings.Contains(typ, "feed+json") {
			return
		}
		ref, ok := s.Attr("href")
		if !ok {
			return
		}
		parsed, err := url.Parse(strings.TrimSpace(ref))
		if err != nil {
			return
		}
		target := base.ResolveReference(parsed)
		if validSubscriptionURL(target) == nil {
			feeds = append(feeds, target.String())
		}
	})
	if len(feeds) == 0 {
		return out, errors.New("no RSS, Atom, or JSON Feed advertised by this website; paste its feed URL directly")
	}
	for _, feed := range feeds {
		if _, resolved, err := fetchSubscriptionURL(ctx, client, feed); err == nil {
			out.URL = resolved
			return out, nil
		}
	}
	return out, errors.New("the website advertises feeds, but none could be fetched")
}

func handleSubscriptionsDiscover(w http.ResponseWriter, r *http.Request) {
	if !sameOriginRequest(r) {
		subscriptionJSON(w, http.StatusForbidden, discoveredSubscription{Error: "local same-origin requests only"})
		return
	}
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		subscriptionJSON(w, http.StatusMethodNotAllowed, discoveredSubscription{Error: "method not allowed"})
		return
	}
	var input struct {
		URL string `json:"url"`
	}
	if !subscriptionReadRequest(w, r, &input) {
		return
	}
	if len(input.URL) > 2048 {
		subscriptionJSON(w, http.StatusBadRequest, discoveredSubscription{Error: "URL is too long"})
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	d, err := discoverSubscription(ctx, strings.TrimSpace(input.URL))
	if err != nil {
		subscriptionJSON(w, http.StatusUnprocessableEntity, discoveredSubscription{Error: err.Error()})
		return
	}
	subscriptionJSON(w, http.StatusOK, d)
}
