package cmd

import (
	"slices"
	"strings"

	"github.com/WaylonWalker/markata/pkg/models"
)

const (
	subscriptionFieldDescription = "description"
	subscriptionFieldCategory    = "category"
	subscriptionFieldType        = "type"
	subscriptionFieldTags        = "tags"
	subscriptionFieldHandle      = "handle"
)

func subscriptionUpdateMatches(feed *models.ExternalFeedConfig, req subscriptionRequest) bool {
	names := req.selectedFields
	if len(names) == 0 {
		names = []string{listSortTitle, subscriptionFieldDescription, subscriptionFieldCategory, subscriptionFieldType, subscriptionFieldTags, subscriptionFieldHandle, blogrollStatusActive}
	}
	actual := map[string]bool{
		listSortTitle:                feed.Title == strings.TrimSpace(req.Title),
		subscriptionFieldDescription: feed.Description == strings.TrimSpace(req.Description),
		subscriptionFieldCategory:    feed.Category == strings.TrimSpace(req.Category),
		subscriptionFieldType:        string(feed.Type) == req.Type,
		subscriptionFieldTags:        slices.Equal(feed.Tags, parseTagsBlogroll(req.Tags)),
		subscriptionFieldHandle:      feed.Handle == strings.TrimSpace(req.Handle),
		blogrollStatusActive:         feed.IsActive() == req.Active,
	}
	for _, name := range names {
		if !actual[name] {
			return false
		}
	}
	return true
}
