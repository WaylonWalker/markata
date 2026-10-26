package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/WaylonWalker/markata/pkg/models"
)

func validateFeedViews(field string, views []string) []error {
	if views == nil {
		return nil
	}
	var errs []error
	hasDefault := false
	for i, view := range views {
		if view == models.FeedViewDefault {
			hasDefault = true
		}
		if !models.IsKnownFeedView(view) {
			errs = append(errs, ValidationError{
				Field:   fmt.Sprintf("%s[%d]", field, i),
				Message: fmt.Sprintf("unknown feed view %q; supported views are %q, %q, %q, %q", view, models.FeedViewDefault, models.FeedViewSimple, models.FeedViewCalendar, models.FeedViewShorts),
			})
		}
	}
	if !hasDefault {
		errs = append(errs, ValidationError{
			Field:   field,
			Message: fmt.Sprintf("must include %q; the primary feed presentation cannot be disabled", models.FeedViewDefault),
		})
	}
	return errs
}

func validateFeedViewsWithPositions(field string, views []string, tracker *PositionTracker, configErrors *ConfigErrors) {
	if views == nil {
		return
	}
	for _, err := range validateFeedViews(field, views) {
		var validationErr ValidationError
		if !errors.As(err, &validationErr) {
			continue
		}
		configErrors.Add(NewConfigErrorWithFix(
			tracker,
			validationErr.Field,
			strings.Join(views, ", "),
			validationErr.Message,
			fmt.Sprintf("Use a subset of [%s, %s, %s, %s] and always keep %s.", models.FeedViewDefault, models.FeedViewSimple, models.FeedViewCalendar, models.FeedViewShorts, models.FeedViewDefault),
			false,
		))
	}
}
