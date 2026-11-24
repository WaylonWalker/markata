// Package lint provides linting functionality for markata-go.
// This file contains blogroll-specific lint checks.
package lint

import (
	"fmt"

	"github.com/WaylonWalker/markata/pkg/models"
)

// BlogrollIssue represents a linting issue found in blogroll configuration.
type BlogrollIssue struct {
	Code     string   // Issue code (e.g., "LBL001")
	Severity Severity // Severity level
	Message  string   // Human-readable message
	FeedURL  string   // Feed URL where issue was found (if applicable)
	Handle   string   // Handle involved (if applicable)
}

// BlogrollResult contains the linting results for blogroll configuration.
type BlogrollResult struct {
	Issues []BlogrollIssue
}

// HasErrors returns true if any issues are errors.
func (r *BlogrollResult) HasErrors() bool {
	for _, issue := range r.Issues {
		if issue.Severity == SeverityError {
			return true
		}
	}
	return false
}

// ErrorCount returns the number of error-severity issues.
func (r *BlogrollResult) ErrorCount() int {
	count := 0
	for _, issue := range r.Issues {
		if issue.Severity == SeverityError {
			count++
		}
	}
	return count
}

// WarningCount returns the number of warning-severity issues.
func (r *BlogrollResult) WarningCount() int {
	count := 0
	for _, issue := range r.Issues {
		if issue.Severity == SeverityWarning {
			count++
		}
	}
	return count
}

// Blogroll checks blogroll configuration for common issues.
// Returns a BlogrollResult containing all issues found.
//
// Supported checks:
//   - LBL001: Duplicate Handle Detection - errors if same handle appears multiple times
//   - LBL002: Duplicate URL Detection - errors if same feed URL appears multiple times
//   - LBL003: Primary Person Validation - errors if primary_person references non-existent handle
func Blogroll(config *models.BlogrollConfig) *BlogrollResult {
	result := &BlogrollResult{}

	if config == nil || len(config.Feeds) == 0 {
		return result
	}

	// LBL001: Check for duplicate handles
	result.Issues = append(result.Issues, checkDuplicateHandles(config.Feeds)...)

	// LBL002: Check for duplicate URLs
	result.Issues = append(result.Issues, checkDuplicateURLs(config.Feeds)...)

	// LBL003: Check primary_person references
	result.Issues = append(result.Issues, checkPrimaryPersonRefs(config.Feeds)...)

	return result
}

// checkDuplicateHandles detects duplicate handles in feed configurations.
// LBL001: Duplicate Handle Detection
func checkDuplicateHandles(feeds []models.ExternalFeedConfig) []BlogrollIssue {
	var issues []BlogrollIssue
	handleSeen := make(map[string]string) // handle -> first feed URL

	for i := range feeds {
		feed := &feeds[i]
		if feed.Handle == "" {
			continue
		}

		if firstURL, exists := handleSeen[feed.Handle]; exists {
			issues = append(issues, BlogrollIssue{
				Code:     "LBL001",
				Severity: SeverityError,
				Message:  fmt.Sprintf("duplicate handle '%s' (first used by %s)", feed.Handle, firstURL),
				FeedURL:  feed.URL,
				Handle:   feed.Handle,
			})
		} else {
			handleSeen[feed.Handle] = feed.URL
		}
	}

	return issues
}

// checkDuplicateURLs detects duplicate feed URLs in configuration.
// LBL002: Duplicate URL Detection
func checkDuplicateURLs(feeds []models.ExternalFeedConfig) []BlogrollIssue {
	var issues []BlogrollIssue
	urlSeen := make(map[string]string) // URL -> handle (or title if no handle)

	for i := range feeds {
		feed := &feeds[i]
		if feed.URL == "" {
			continue
		}

		identifier := feed.Handle
		if identifier == "" {
			identifier = feed.Title
		}
		if identifier == "" {
			identifier = feed.URL
		}

		if firstIdentifier, exists := urlSeen[feed.URL]; exists {
			issues = append(issues, BlogrollIssue{
				Code:     "LBL002",
				Severity: SeverityError,
				Message:  fmt.Sprintf("duplicate feed URL '%s' (first used by %s)", feed.URL, firstIdentifier),
				FeedURL:  feed.URL,
				Handle:   feed.Handle,
			})
		} else {
			urlSeen[feed.URL] = identifier
		}
	}

	return issues
}

// checkPrimaryPersonRefs validates that primary_person references exist.
// LBL003: Primary Person Validation
func checkPrimaryPersonRefs(feeds []models.ExternalFeedConfig) []BlogrollIssue {
	var issues []BlogrollIssue

	// Build set of valid handles
	validHandles := make(map[string]bool)
	for i := range feeds {
		if feeds[i].Handle != "" {
			validHandles[feeds[i].Handle] = true
		}
	}

	// Check primary_person references
	for i := range feeds {
		feed := &feeds[i]
		if feed.PrimaryPerson == "" {
			continue
		}

		if !validHandles[feed.PrimaryPerson] {
			issues = append(issues, BlogrollIssue{
				Code:     "LBL003",
				Severity: SeverityError,
				Message:  fmt.Sprintf("primary_person '%s' references non-existent handle", feed.PrimaryPerson),
				FeedURL:  feed.URL,
				Handle:   feed.Handle,
			})
		}
	}

	return issues
}
