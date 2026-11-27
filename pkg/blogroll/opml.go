package blogroll

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
)

// OPMLFeed is the interoperable subset of a subscription shared with other readers.
type OPMLFeed struct {
	URL         string
	Title       string
	SiteURL     string
	Description string
	Category    string
}

type opmlDocument struct {
	XMLName xml.Name      `xml:"opml"`
	Version string        `xml:"version,attr"`
	Title   string        `xml:"head>title"`
	Body    []opmlOutline `xml:"body>outline"`
}

type opmlOutline struct {
	Text        string        `xml:"text,attr"`
	Title       string        `xml:"title,attr,omitempty"`
	Type        string        `xml:"type,attr,omitempty"`
	URL         string        `xml:"xmlUrl,attr,omitempty"`
	SiteURL     string        `xml:"htmlUrl,attr,omitempty"`
	Description string        `xml:"description,attr,omitempty"`
	Category    string        `xml:"category,attr,omitempty"`
	Children    []opmlOutline `xml:"outline"`
}

// ReadOPML parses a bounded subscription export. Nested folders become
// slash-separated categories; duplicate URLs are preserved for the importer
// to handle explicitly. Non-feed outlines are ignored.
func ReadOPML(r io.Reader) ([]OPMLFeed, error) {
	const maxBytes = 4 << 20
	data, err := io.ReadAll(io.LimitReader(r, maxBytes+1))
	if err != nil {
		return nil, fmt.Errorf("read OPML: %w", err)
	}
	if len(data) > maxBytes {
		return nil, errors.New("OPML exceeds 4 MiB")
	}
	var doc opmlDocument
	if err := xml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("parse OPML: %w", err)
	}
	feeds := make([]OPMLFeed, 0)
	for _, outline := range doc.Body {
		if err := collectOPMLFeeds(outline, "", &feeds); err != nil {
			return nil, err
		}
	}
	return feeds, nil
}

func collectOPMLFeeds(outline opmlOutline, parent string, feeds *[]OPMLFeed) error {
	title := strings.TrimSpace(outline.Title)
	if title == "" {
		title = strings.TrimSpace(outline.Text)
	}
	if address := strings.TrimSpace(outline.URL); address != "" {
		if err := validateOPMLURL(address); err != nil {
			return fmt.Errorf("invalid feed URL in OPML outline %q: %w", title, err)
		}
		category := strings.TrimSpace(outline.Category)
		if category == "" {
			category = parent
		}
		*feeds = append(*feeds, OPMLFeed{
			URL: address, Title: title, SiteURL: strings.TrimSpace(outline.SiteURL),
			Description: outline.Description, Category: category,
		})
	} else if title != "" {
		if parent == "" {
			parent = title
		} else {
			parent += "/" + title
		}
	}
	for _, child := range outline.Children {
		if err := collectOPMLFeeds(child, parent, feeds); err != nil {
			return err
		}
	}
	return nil
}

func validateOPMLURL(address string) error {
	u, err := url.Parse(address)
	if err != nil || u.Hostname() == "" || (u.Scheme != "http" && u.Scheme != "https") || u.User != nil {
		return errors.New("expected an absolute HTTP(S) feed URL without embedded credentials")
	}
	return nil
}

// WriteOPML exports subscriptions in OPML 2.0. Each category is represented
// by a folder. Feed URLs, site URLs, titles and descriptions are XML-escaped.
func WriteOPML(w io.Writer, title string, feeds []OPMLFeed) error {
	doc := opmlDocument{Version: "2.0", Title: title}
	folders := make(map[string]int)
	for _, feed := range feeds {
		if err := validateOPMLURL(feed.URL); err != nil {
			return fmt.Errorf("cannot export subscription %q: %w", feed.Title, err)
		}
		name := feed.Title
		if name == "" {
			name = feed.URL
		}
		outline := opmlOutline{
			Text: name, Title: name, Type: "rss", URL: feed.URL,
			SiteURL: feed.SiteURL, Description: feed.Description,
		}
		if feed.Category == "" {
			doc.Body = append(doc.Body, outline)
			continue
		}
		index, ok := folders[feed.Category]
		if !ok {
			index = len(doc.Body)
			folders[feed.Category] = index
			doc.Body = append(doc.Body, opmlOutline{Text: feed.Category, Title: feed.Category})
		}
		doc.Body[index].Children = append(doc.Body[index].Children, outline)
	}
	if _, err := io.WriteString(w, xml.Header); err != nil {
		return err
	}
	encoder := xml.NewEncoder(w)
	encoder.Indent("", "  ")
	return encoder.Encode(doc)
}
