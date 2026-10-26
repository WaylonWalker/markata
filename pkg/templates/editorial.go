package templates

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/WaylonWalker/markata/pkg/models"
)

const (
	editorialHTTPS    = "https"
	editorialEdited   = "edited"
	editorialReviewed = "reviewed"
)

// editorialRecord is the portable v1 frontmatter event, not a freshness claim.
type editorialRecord struct {
	SchemaVersion int    `json:"schema_version"`
	Author        string `json:"author"`
	Role          string `json:"role"`
	Action        string `json:"action"`
	RecordedAt    string `json:"recorded_at"`
	RecordedBy    string `json:"recorded_by"`
	BaseCommit    string `json:"base_commit"`
	Commit        string `json:"commit"`
	Path          string `json:"path"`
	FileURL       string `json:"file_url"`
	Summary       string `json:"summary"`
}

func editorialHistoryToMaps(p *models.Post) []map[string]interface{} {
	if p.Extra == nil {
		return nil
	}
	encoded, err := json.Marshal(p.Extra["editorial_history"])
	if err != nil {
		return nil
	}
	var records []json.RawMessage
	if err := json.Unmarshal(encoded, &records); err != nil {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(records))
	for i := len(records) - 1; i >= 0; i-- {
		var record editorialRecord
		if err := json.Unmarshal(records[i], &record); err != nil || record.SchemaVersion != 1 ||
			record.Author == "" || (record.Action != editorialEdited && record.Action != editorialReviewed) {
			continue
		}
		name := record.Author
		for j := range p.AuthorObjects {
			if p.AuthorObjects[j].ID == record.Author {
				name = p.AuthorObjects[j].Name
				break
			}
		}
		item := map[string]interface{}{
			"author": record.Author, "name": name, "role": strings.ReplaceAll(record.Role, "-", " "),
			"action": record.Action, "recorded_at": record.RecordedAt, "recorded_by": record.RecordedBy,
			"commit": record.Commit, "summary": record.Summary,
		}
		if len(record.Commit) >= 8 {
			item["short_commit"] = record.Commit[:8]
		} else {
			item["short_commit"] = record.Commit
		}
		if parsed, err := time.Parse(time.RFC3339, record.RecordedAt); err == nil {
			item["recorded_date"] = parsed
		}
		addEditorialRevisionURLs(item, record)
		result = append(result, item)
	}
	return result
}

func fullEditorialSHA(value string) bool {
	if len(value) != 40 {
		return false
	}
	_, err := hex.DecodeString(value)
	return err == nil
}

func validEditorialPath(value string) bool {
	return value != "" && !strings.ContainsAny(value, "\\\x00") &&
		!strings.HasPrefix(value, "/") && path.Clean(value) == value && value != "." &&
		value != ".." && !strings.HasPrefix(value, "../")
}

func validEditorialRepositoryPart(value string) bool {
	return value != "" && value != "." && value != ".."
}

func editorialRepository(r editorialRecord) (string, bool) {
	u, err := url.Parse(r.FileURL)
	if err != nil || u.Scheme != editorialHTTPS || u.Host != "github.com" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return "", false
	}
	parts := strings.Split(strings.TrimPrefix(u.Path, "/"), "/")
	if len(parts) < 5 || !validEditorialRepositoryPart(parts[0]) || !validEditorialRepositoryPart(parts[1]) ||
		parts[2] != "blob" || parts[3] != r.Commit || strings.Join(parts[4:], "/") != r.Path {
		return "", false
	}
	return parts[0] + "/" + parts[1], true
}

func addEditorialRevisionURLs(item map[string]interface{}, r editorialRecord) {
	if !fullEditorialSHA(r.Commit) || !validEditorialPath(r.Path) {
		return
	}
	repository, valid := editorialRepository(r)
	if !valid {
		return
	}
	fileURL := &url.URL{Scheme: editorialHTTPS, Host: "github.com", Path: "/" + repository + "/blob/" + r.Commit + "/" + r.Path}
	diffURL := &url.URL{Scheme: editorialHTTPS, Host: "github.com", Path: "/" + repository + "/commit/" + r.Commit}
	anchor := sha256.Sum256([]byte(r.Path))
	diffURL.Fragment = "diff-" + hex.EncodeToString(anchor[:])
	item["file_url"] = fileURL.String()
	item["diff_url"] = diffURL.String()
	if r.Action != editorialEdited || !fullEditorialSHA(r.BaseCommit) {
		return
	}
	before := &url.URL{Scheme: editorialHTTPS, Host: "raw.githubusercontent.com", Path: "/" + repository + "/" + r.BaseCommit + "/" + r.Path}
	after := &url.URL{Scheme: editorialHTTPS, Host: "raw.githubusercontent.com", Path: "/" + repository + "/" + r.Commit + "/" + r.Path}
	item["before_url"] = before.String()
	item["after_url"] = after.String()
}

func isEditorialContributor(a *models.Author) bool {
	if a.IsPrimaryContributor() || (a.Role != nil && strings.EqualFold(*a.Role, "author")) {
		return false
	}
	role := strings.ToLower(a.GetRoleDisplay())
	return strings.Contains(role, "editor") || strings.Contains(role, "editing") || strings.Contains(role, "review")
}
