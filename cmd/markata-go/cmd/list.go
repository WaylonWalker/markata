package cmd

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/WaylonWalker/markata/pkg/filter"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/listcache"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/services"

	"github.com/spf13/cobra"
)

const (
	listFormatTable = "table"
	listFormatJSON  = "json"
	listFormatCSV   = "csv"
	listFormatPath  = "path"
	listSortDate    = "date"
	listSortTitle   = "title"
	listSortName    = "name"
	listSortPosts   = "posts"
	listSortCount   = "count"
	listSortWords   = "words"
	listSortReading = "reading_time"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List posts, tags, or feeds",
	Long: `List posts, tags, or feeds for quick inspection and scripting.

Use subcommands to select the data source:
  markata-go list posts
  markata-go list tags
  markata-go list feeds`,
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.AddCommand(listPostsCmd())
	listCmd.AddCommand(listTagsCmd())
	listCmd.AddCommand(listFeedsCmd())
}

type listCommonOptions struct {
	format string
	sortBy string
	order  string
}

func listPostsCmd() *cobra.Command {
	var opts listCommonOptions
	var filter string
	var feed string

	cmd := &cobra.Command{
		Use:   "posts",
		Short: "List posts",
		RunE: func(cmd *cobra.Command, _ []string) error {
			app, err := loadListApp(cmd.Context())
			if err != nil {
				return err
			}

			format, err := parseListFormat(opts.format)
			if err != nil {
				return err
			}

			sortBy := opts.sortBy
			if sortBy == "" {
				sortBy = listSortDate
			}

			order, err := parseSortOrder(opts.order)
			if err != nil {
				return err
			}

			if !isValidPostSort(sortBy) {
				return invalidChoiceError("--sort", sortBy, []string{"date", "title", "words", "path", "reading_time", "tags"}, cmd.CommandPath())
			}

			listOpts := services.ListOptions{
				Filter:    filter,
				SortBy:    sortBy,
				SortOrder: order,
			}

			var posts []*models.Post
			if feed != "" {
				posts, err = postsForFeed(cmd.Context(), app, feed, listOpts, filter)
				if err != nil {
					return err
				}
			} else {
				posts, err = app.Posts.List(cmd.Context(), listOpts)
				if err != nil {
					return err
				}
			}

			return renderPosts(format, posts)
		},
	}

	cmd.Flags().StringVar(&opts.format, "format", listFormatTable, "output format: table, json, csv, path")
	cmd.Flags().StringVar(&opts.sortBy, "sort", listSortDate, "sort field: date, title, words, path, reading_time, tags")
	cmd.Flags().StringVar(&opts.order, "order", "desc", "sort order: asc or desc")
	cmd.Flags().StringVar(&filter, "filter", "", "filter expression for posts")
	cmd.Flags().StringVar(&feed, "feed", "", "limit posts to a feed by name")

	return cmd
}

func listTagsCmd() *cobra.Command {
	var opts listCommonOptions

	cmd := &cobra.Command{
		Use:   "tags",
		Short: "List tags",
		RunE: func(cmd *cobra.Command, _ []string) error {
			app, err := loadListApp(cmd.Context())
			if err != nil {
				return err
			}

			format, err := parseListFormat(opts.format)
			if err != nil {
				return err
			}

			sortBy := opts.sortBy
			if sortBy == "" {
				sortBy = listSortCount
			}

			order, err := parseSortOrder(opts.order)
			if err != nil {
				return err
			}

			if !isValidTagSort(sortBy) {
				return invalidChoiceError("--sort", sortBy, []string{"name", "count", "words", "reading_time"}, cmd.CommandPath())
			}

			tags, err := app.Tags.List(cmd.Context())
			if err != nil {
				return err
			}

			posts := app.Manager.Posts()
			stats := buildTagStats(posts)
			rows := buildTagRows(tags, stats)
			sortTagRows(rows, sortBy, order)

			return renderTags(format, rows)
		},
	}

	cmd.Flags().StringVar(&opts.format, "format", listFormatTable, "output format: table, json, csv, path")
	cmd.Flags().StringVar(&opts.sortBy, "sort", "count", "sort field: name, count, words, reading_time")
	cmd.Flags().StringVar(&opts.order, "order", "desc", "sort order: asc or desc")

	return cmd
}

func listFeedsCmd() *cobra.Command {
	var opts listCommonOptions

	cmd := &cobra.Command{
		Use:   "feeds",
		Short: "List feeds",
		RunE: func(cmd *cobra.Command, _ []string) error {
			app, err := loadListApp(cmd.Context())
			if err != nil {
				return err
			}

			format, err := parseListFormat(opts.format)
			if err != nil {
				return err
			}

			sortBy := opts.sortBy
			if sortBy == "" {
				sortBy = listSortName
			}

			order, err := parseSortOrder(opts.order)
			if err != nil {
				return err
			}

			if !isValidFeedSort(sortBy) {
				return invalidChoiceError("--sort", sortBy, []string{"name", "posts", "words", "reading_time", "avg_reading_time"}, cmd.CommandPath())
			}

			feeds, err := app.Feeds.List(cmd.Context())
			if err != nil {
				return err
			}

			rows := buildFeedRows(feeds)
			sortFeedRows(rows, sortBy, order)

			return renderFeeds(format, rows)
		},
	}

	cmd.Flags().StringVar(&opts.format, "format", listFormatTable, "output format: table, json, csv, path")
	cmd.Flags().StringVar(&opts.sortBy, "sort", "name", "sort field: name, posts, words, reading_time, avg_reading_time")
	cmd.Flags().StringVar(&opts.order, "order", "asc", "sort order: asc or desc")

	cmd.AddCommand(listFeedsPostsCmd())

	return cmd
}

func listFeedsPostsCmd() *cobra.Command {
	var format string

	cmd := &cobra.Command{
		Use:   "posts <feed>",
		Short: "List posts for a feed",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			app, err := loadListApp(cmd.Context())
			if err != nil {
				return err
			}

			outputFormat, err := parseListFormat(format)
			if err != nil {
				return err
			}

			posts, err := postsForFeed(cmd.Context(), app, args[0], services.ListOptions{
				SortBy:    listSortDate,
				SortOrder: services.SortDesc,
			}, "")
			if err != nil {
				return err
			}

			return renderPosts(outputFormat, posts)
		},
	}

	cmd.Flags().StringVar(&format, "format", listFormatTable, "output format: table, json, csv, path")

	return cmd
}

func loadListApp(ctx context.Context) (*services.App, error) {
	manager, err := createManager(cfgFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	app := services.NewApp(manager)
	configHash, err := configFilesHash(cfgFile, mergeConfigFiles)
	if err != nil {
		return nil, err
	}
	listcache.SetOptions(manager, listcache.Options{
		CacheDir:   listcache.DefaultCacheDir,
		ConfigHash: configHash,
	})

	if err := app.Build.LoadForTUI(ctx); err != nil {
		return nil, fmt.Errorf("failed to load posts: %w", err)
	}

	return app, nil
}

func parseListFormat(format string) (string, error) {
	format = strings.ToLower(strings.TrimSpace(format))
	switch format {
	case listFormatTable, listFormatJSON, listFormatCSV, listFormatPath:
		return format, nil
	default:
		return "", invalidChoiceError("--format", format, []string{listFormatTable, listFormatJSON, listFormatCSV, listFormatPath}, currentCommandPath())
	}
}

func parseSortOrder(order string) (services.SortOrder, error) {
	order = strings.ToLower(strings.TrimSpace(order))
	switch order {
	case string(services.SortAsc):
		return services.SortAsc, nil
	case string(services.SortDesc):
		return services.SortDesc, nil
	default:
		return "", invalidChoiceError("--order", order, []string{"asc", "desc"}, currentCommandPath())
	}
}

func currentCommandPath() string {
	if currentCmd != nil {
		return currentCmd.CommandPath()
	}
	return "markata-go list"
}

func isValidPostSort(field string) bool {
	switch strings.ToLower(field) {
	case listSortDate, listSortTitle, listSortWords, "path", listSortReading, "tags":
		return true
	default:
		return false
	}
}

func isValidTagSort(field string) bool {
	switch strings.ToLower(field) {
	case listSortName, listSortCount, listSortWords, listSortReading:
		return true
	default:
		return false
	}
}

func isValidFeedSort(field string) bool {
	switch strings.ToLower(field) {
	case listSortName, listSortPosts, listSortWords, listSortReading, "avg_reading_time":
		return true
	default:
		return false
	}
}

type postRow struct {
	Title       string   `json:"title"`
	Date        string   `json:"date"`
	Words       int      `json:"words"`
	ReadingTime int      `json:"reading_time"`
	Tags        []string `json:"tags"`
	Path        string   `json:"path"`
}

func renderPosts(format string, posts []*models.Post) error {
	rows := make([]postRow, 0, len(posts))
	for _, post := range posts {
		rows = append(rows, postToRow(post))
	}

	switch format {
	case listFormatTable:
		return renderPostsTable(rows)
	case listFormatJSON:
		return renderJSON(rows)
	case listFormatCSV:
		return renderPostsCSV(rows)
	case listFormatPath:
		for _, row := range rows {
			outln(row.Path)
		}
		return nil
	default:
		return fmt.Errorf("unsupported format %q", format)
	}
}

func postToRow(post *models.Post) postRow {
	return postRow{
		Title:       postTitle(post),
		Date:        postDate(post),
		Words:       postWordCount(post),
		ReadingTime: postReadingTime(post),
		Tags:        post.Tags,
		Path:        sourcePathForOutput(post.Path),
	}
}

func renderPostsTable(rows []postRow) error {
	w := tabwriter.NewWriter(outWriter(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TITLE\tDATE\tWORDS\tREAD\tTAGS\tPATH")
	fmt.Fprintln(w, "-----\t----\t-----\t----\t----\t----")
	for _, row := range rows {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			row.Title,
			row.Date,
			formatWordCount(row.Words),
			formatReadingTime(row.ReadingTime),
			strings.Join(row.Tags, ", "),
			row.Path,
		)
	}
	return w.Flush()
}

func renderPostsCSV(rows []postRow) error {
	w := csv.NewWriter(outWriter())
	if err := w.Write([]string{"title", listSortDate, "words", "reading_time", "tags", "path"}); err != nil {
		return err
	}
	for _, row := range rows {
		tags := strings.Join(row.Tags, ", ")
		record := []string{
			row.Title,
			row.Date,
			fmt.Sprintf("%d", row.Words),
			fmt.Sprintf("%d", row.ReadingTime),
			tags,
			row.Path,
		}
		if err := w.Write(record); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}

type tagRow struct {
	Name        string `json:"name"`
	Count       int    `json:"count"`
	Words       int    `json:"words"`
	ReadingTime int    `json:"reading_time"`
	Slug        string `json:"slug"`
}

type tagStats struct {
	Words       int
	ReadingTime int
}

func buildTagStats(posts []*models.Post) map[string]tagStats {
	stats := make(map[string]tagStats)
	for _, post := range posts {
		if len(post.Tags) == 0 {
			continue
		}
		words := postWordCount(post)
		readingTime := postReadingTime(post)
		for _, tag := range post.Tags {
			entry := stats[tag]
			entry.Words += words
			entry.ReadingTime += readingTime
			stats[tag] = entry
		}
	}
	return stats
}

func buildTagRows(tags []services.TagInfo, stats map[string]tagStats) []tagRow {
	rows := make([]tagRow, 0, len(tags))
	for _, tag := range tags {
		stat := stats[tag.Name]
		rows = append(rows, tagRow{
			Name:        tag.Name,
			Count:       tag.Count,
			Words:       stat.Words,
			ReadingTime: stat.ReadingTime,
			Slug:        tag.Slug,
		})
	}
	return rows
}

func sortTagRows(rows []tagRow, field string, order services.SortOrder) {
	sort.SliceStable(rows, func(i, j int) bool {
		var cmp int
		switch strings.ToLower(field) {
		case listSortCount:
			cmp = compareInts(rows[i].Count, rows[j].Count)
		case listSortWords:
			cmp = compareInts(rows[i].Words, rows[j].Words)
		case listSortReading:
			cmp = compareInts(rows[i].ReadingTime, rows[j].ReadingTime)
		case listSortName:
			cmp = strings.Compare(strings.ToLower(rows[i].Name), strings.ToLower(rows[j].Name))
		default:
			cmp = 0
		}

		if cmp == 0 {
			cmp = strings.Compare(strings.ToLower(rows[i].Name), strings.ToLower(rows[j].Name))
		}

		if order == services.SortDesc {
			return cmp > 0
		}
		return cmp < 0
	})
}

func renderTags(format string, rows []tagRow) error {
	switch format {
	case listFormatTable:
		return renderTagsTable(rows)
	case listFormatJSON:
		return renderJSON(rows)
	case listFormatCSV:
		return renderTagsCSV(rows)
	case listFormatPath:
		for _, row := range rows {
			outln(row.Name)
		}
		return nil
	default:
		return fmt.Errorf("unsupported format %q", format)
	}
}

func renderTagsTable(rows []tagRow) error {
	w := tabwriter.NewWriter(outWriter(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TAG\tCOUNT\tWORDS\tREAD\tSLUG")
	fmt.Fprintln(w, "---\t-----\t-----\t----\t----")
	for _, row := range rows {
		fmt.Fprintf(w, "%s\t%d\t%s\t%s\t%s\n",
			row.Name,
			row.Count,
			formatWordCount(row.Words),
			formatReadingTime(row.ReadingTime),
			row.Slug,
		)
	}
	return w.Flush()
}

func renderTagsCSV(rows []tagRow) error {
	w := csv.NewWriter(outWriter())
	if err := w.Write([]string{"name", "count", "words", "reading_time", "slug"}); err != nil {
		return err
	}
	for _, row := range rows {
		record := []string{
			row.Name,
			fmt.Sprintf("%d", row.Count),
			fmt.Sprintf("%d", row.Words),
			fmt.Sprintf("%d", row.ReadingTime),
			row.Slug,
		}
		if err := w.Write(record); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}

type feedRow struct {
	Name           string `json:"name"`
	Posts          int    `json:"posts"`
	Words          int    `json:"words"`
	ReadingTime    int    `json:"reading_time"`
	AvgReadingTime int    `json:"avg_reading_time"`
	Path           string `json:"output"`
}

func buildFeedRows(feeds []*lifecycle.Feed) []feedRow {
	rows := make([]feedRow, 0, len(feeds))
	for _, feed := range feeds {
		totalWords, totalReadingTime := calculateFeedStats(feed.Posts)
		avgReadingTime := 0
		if len(feed.Posts) > 0 {
			avgReadingTime = totalReadingTime / len(feed.Posts)
		}

		rows = append(rows, feedRow{
			Name:           feed.Name,
			Posts:          len(feed.Posts),
			Words:          totalWords,
			ReadingTime:    totalReadingTime,
			AvgReadingTime: avgReadingTime,
			Path:           feed.Path,
		})
	}
	return rows
}

func sortFeedRows(rows []feedRow, field string, order services.SortOrder) {
	sort.SliceStable(rows, func(i, j int) bool {
		var cmp int
		switch strings.ToLower(field) {
		case listSortPosts:
			cmp = compareInts(rows[i].Posts, rows[j].Posts)
		case listSortWords:
			cmp = compareInts(rows[i].Words, rows[j].Words)
		case listSortReading:
			cmp = compareInts(rows[i].ReadingTime, rows[j].ReadingTime)
		case "avg_reading_time":
			cmp = compareInts(rows[i].AvgReadingTime, rows[j].AvgReadingTime)
		case listSortName:
			cmp = strings.Compare(strings.ToLower(rows[i].Name), strings.ToLower(rows[j].Name))
		default:
			cmp = 0
		}

		if cmp == 0 {
			cmp = strings.Compare(strings.ToLower(rows[i].Name), strings.ToLower(rows[j].Name))
		}

		if order == services.SortDesc {
			return cmp > 0
		}
		return cmp < 0
	})
}

func renderFeeds(format string, rows []feedRow) error {
	switch format {
	case listFormatTable:
		return renderFeedsTable(rows)
	case listFormatJSON:
		return renderJSON(rows)
	case listFormatCSV:
		return renderFeedsCSV(rows)
	case listFormatPath:
		for _, row := range rows {
			value := row.Path
			if value == "" {
				value = row.Name
			}
			outln(value)
		}
		return nil
	default:
		return fmt.Errorf("unsupported format %q", format)
	}
}

func renderFeedsTable(rows []feedRow) error {
	w := tabwriter.NewWriter(outWriter(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tPOSTS\tWORDS\tTOTAL READ\tAVG READ\tOUTPUT")
	fmt.Fprintln(w, "----\t-----\t-----\t----------\t--------\t------")
	for _, row := range rows {
		fmt.Fprintf(w, "%s\t%d\t%s\t%s\t%s\t%s\n",
			row.Name,
			row.Posts,
			formatWordCount(row.Words),
			formatReadingTime(row.ReadingTime),
			formatReadingTime(row.AvgReadingTime),
			row.Path,
		)
	}
	return w.Flush()
}

func renderFeedsCSV(rows []feedRow) error {
	w := csv.NewWriter(outWriter())
	if err := w.Write([]string{"name", listSortPosts, "words", "reading_time", "avg_reading_time", "output"}); err != nil {
		return err
	}
	for _, row := range rows {
		record := []string{
			row.Name,
			fmt.Sprintf("%d", row.Posts),
			fmt.Sprintf("%d", row.Words),
			fmt.Sprintf("%d", row.ReadingTime),
			fmt.Sprintf("%d", row.AvgReadingTime),
			row.Path,
		}
		if err := w.Write(record); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}

func renderJSON(v interface{}) error {
	enc := json.NewEncoder(outWriter())
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func postTitle(post *models.Post) string {
	if post == nil || post.PlainTitle() == "" {
		return "(untitled)"
	}
	return post.PlainTitle()
}

func postDate(post *models.Post) string {
	if post == nil || post.Date == nil {
		return ""
	}
	return post.Date.Format("2006-01-02")
}

func postWordCount(post *models.Post) int {
	if post == nil || post.Extra == nil {
		return 0
	}
	if wc, ok := post.Extra["word_count"].(int); ok {
		return wc
	}
	if post.Content != "" {
		return len(strings.Fields(post.Content))
	}
	return 0
}

func postReadingTime(post *models.Post) int {
	if post == nil || post.Extra == nil {
		return 0
	}
	if rt, ok := post.Extra["reading_time"].(int); ok {
		return rt
	}
	return 0
}

func formatWordCount(count int) string {
	if count < 1000 {
		return fmt.Sprintf("%d", count)
	}
	if count < 10000 {
		return fmt.Sprintf("%.1fk", float64(count)/1000)
	}
	return fmt.Sprintf("%dk", count/1000)
}

func formatReadingTime(minutes int) string {
	if minutes == 0 {
		return "<1 min"
	}
	if minutes == 1 {
		return "1 min"
	}
	if minutes < 60 {
		return fmt.Sprintf("%d min", minutes)
	}
	hours := minutes / 60
	mins := minutes % 60
	if mins == 0 {
		return fmt.Sprintf("%dh", hours)
	}
	return fmt.Sprintf("%dh%dm", hours, mins)
}

func calculateFeedStats(posts []*models.Post) (totalWords, totalReadingTime int) {
	for _, post := range posts {
		totalWords += postWordCount(post)
		totalReadingTime += postReadingTime(post)
	}
	return totalWords, totalReadingTime
}

func compareInts(a, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

func postsForFeed(ctx context.Context, app *services.App, feedName string, opts services.ListOptions, filterExpr string) ([]*models.Post, error) {
	if feedName == "" {
		return nil, fmt.Errorf("feed name is required")
	}

	feed, err := resolveFeed(ctx, app, feedName)
	if err != nil {
		return nil, err
	}
	if feed == nil {
		return nil, fmt.Errorf("feed %q not found", feedName)
	}

	posts, err := app.Feeds.GetPosts(ctx, feed.Name, opts)
	if err != nil {
		return nil, err
	}

	if filterExpr == "" {
		return posts, nil
	}

	filtered, err := filter.Posts(filterExpr, posts)
	if err != nil {
		return nil, err
	}
	return filtered, nil
}

func resolveFeed(ctx context.Context, app *services.App, feedName string) (*lifecycle.Feed, error) {
	feeds, err := app.Feeds.List(ctx)
	if err != nil {
		return nil, err
	}

	needle := strings.TrimSpace(feedName)
	if needle == "" {
		return nil, nil
	}

	for _, feed := range feeds {
		if strings.EqualFold(feed.Name, needle) {
			return feed, nil
		}
	}

	for _, feed := range feeds {
		if strings.EqualFold(feed.Path, needle) {
			return feed, nil
		}
	}

	for _, feed := range feeds {
		if strings.EqualFold(feed.Title, needle) {
			return feed, nil
		}
	}

	return nil, nil
}
