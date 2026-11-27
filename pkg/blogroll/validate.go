package blogroll

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
)

// ValidateFeedDocument recognizes supported subscription document formats
// without fetching links, enriching entries or accepting HTML as RSS.
func ValidateFeedDocument(body []byte) (string, error) {
	body = bytes.TrimSpace(bytes.TrimPrefix(body, []byte{0xef, 0xbb, 0xbf}))
	if len(body) == 0 {
		return "", errors.New("empty feed document")
	}
	if body[0] == '{' {
		return validateJSONFeed(body)
	}
	var document struct {
		XMLName xml.Name
		Channel *struct{} `xml:"channel"`
	}
	if err := xml.Unmarshal(body, &document); err != nil {
		return "", fmt.Errorf("invalid feed XML: %w", err)
	}
	switch document.XMLName.Local {
	case "rss":
		if document.Channel == nil {
			return "", errors.New("RSS document is missing its channel")
		}
		return "rss", nil
	case "feed":
		if document.XMLName.Space != "" && document.XMLName.Space != "http://www.w3.org/2005/Atom" {
			return "", errors.New("unrecognized Atom namespace")
		}
		return "atom", nil
	default:
		return "", fmt.Errorf("not an RSS or Atom feed (XML root %q)", document.XMLName.Local)
	}
}

func validateJSONFeed(body []byte) (string, error) {
	var document struct {
		Version string            `json:"version"`
		Items   []json.RawMessage `json:"items"`
	}
	if err := json.Unmarshal(body, &document); err != nil {
		return "", fmt.Errorf("invalid JSON Feed: %w", err)
	}
	if (document.Version != "https://jsonfeed.org/version/1" && document.Version != "https://jsonfeed.org/version/1.1") || document.Items == nil {
		return "", errors.New("not a JSON Feed document")
	}
	return "json", nil
}
