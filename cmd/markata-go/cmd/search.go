package cmd

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/search"
	"github.com/WaylonWalker/markata/pkg/services"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(searchCmd())
}

func searchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search posts and manage bleve indexes",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return searchRunCmd().RunE(cmd, args)
		},
	}

	cmd.AddCommand(searchRunCmd())
	cmd.AddCommand(searchBuildIndexCmd())
	return cmd
}

func searchRunCmd() *cobra.Command {
	var (
		format string
		sortBy string
		order  string
		filter string
		fields string
		fuzzy  bool
		limit  int
	)

	cmd := &cobra.Command{
		Use:     "run <query>",
		Short:   "Full-text search across posts",
		Aliases: []string{"query"},
		Long: `Search post content, titles, descriptions, and tags.

Uses a bleve full-text index for BM25-ranked results with fuzzy matching.
The index is built on first search and cached for subsequent queries.

Examples:
  markata-go search golang
  markata-go search "error handling" --format json
  markata-go search docker --filter "published == True" --limit 10
  markata-go search cli --sort date --order desc
  markata-go search kubernetes --format path
  markata-go search golang --fields title,tags
  markata-go search tutoral --fuzzy
  markata-go search tutorial --filter '"go" in tags and date >= "2024-01-01"'`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			queryStr := strings.Join(args, " ")

			app, err := loadListApp(cmd.Context())
			if err != nil {
				return err
			}

			outputFormat, err := parseListFormat(format)
			if err != nil {
				return err
			}

			// Get all posts (possibly filtered)
			var posts []*models.Post
			if filter != "" {
				listOpts := services.ListOptions{Filter: filter}
				posts, err = app.Posts.List(cmd.Context(), listOpts)
				if err != nil {
					return err
				}
			} else {
				posts = app.Manager.Posts()
			}

			// Try bleve-powered ranked search
			results, bleveErr := searchWithBleve(posts, queryStr, search.QueryOptions{
				Limit: limit,
				Fuzzy: fuzzy,
			})

			if bleveErr != nil {
				verbosef("bleve search unavailable, falling back to substring: %v", bleveErr)
				results = searchSubstring(posts, queryStr, fields, limit)
			}

			// Apply sorting — default to "score" for bleve, "date" for substring
			if sortBy == "" {
				if bleveErr == nil {
					sortBy = "score"
				} else {
					sortBy = "date"
				}
			}
			if sortBy != "score" {
				sortOrder, err := parseSortOrder(order)
				if err != nil {
					return err
				}
				if !isValidPostSort(sortBy) {
					return invalidChoiceError("--sort", sortBy, []string{"score", "date", "title", "words", "path", "reading_time", "tags"}, cmd.CommandPath())
				}
				sortSearchResults(results, sortBy, sortOrder)
			}

			if outputFormat == listFormatTable {
				return renderSearchTable(results, queryStr)
			}
			return renderPosts(outputFormat, cliSearchResultPosts(results))
		},
	}

	cmd.Flags().StringVar(&format, "format", listFormatTable, "output format: table, json, csv, path")
	cmd.Flags().StringVar(&sortBy, "sort", "", "sort field: score, date, title, words, path, reading_time, tags (default: score)")
	cmd.Flags().StringVar(&order, "order", "desc", "sort order: asc or desc")
	cmd.Flags().StringVar(&filter, "filter", "", "filter expression (applied before search)")
	cmd.Flags().StringVar(&fields, "fields", "", "fields to search: title,content,description,tags (default: all)")
	cmd.Flags().BoolVar(&fuzzy, "fuzzy", false, "enable fuzzy matching (tolerates typos)")
	cmd.Flags().IntVar(&limit, "limit", 0, "maximum number of results (0 = no limit)")

	return cmd
}

func searchBuildIndexCmd() *cobra.Command {
	var (
		indexDir  string
		hashPath  string
		indexName string
		force     bool
	)

	cmd := &cobra.Command{
		Use:   "build-index",
		Short: "Build a reusable bleve index artifact",
		Long: `Build a bleve index without starting the search server.

This command is intended for builder jobs, CI, and container workflows that
publish a search index artifact for later use by read-only search servers.

Examples:
  markata-go search build-index
  markata-go search build-index --index-dir /data/search.bleve
  markata-go search build-index --index-name web-1 --force`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			app, err := loadListApp(cmd.Context())
			if err != nil {
				return err
			}

			posts := app.Manager.Posts()
			cacheDir := filepath.Join(".markata", "cache")
			if indexDir == "" {
				if indexName != "" {
					indexDir = search.NamedDir(cacheDir, indexName)
				} else {
					indexDir = search.DefaultDir(cacheDir)
				}
			}
			if hashPath == "" {
				if indexName != "" {
					hashPath = search.NamedHashFile(cacheDir, indexName)
				} else {
					hashPath = filepath.Join(cacheDir, "search.hash")
				}
			}

			if force {
				idx, buildErr := search.Build(indexDir, posts)
				if buildErr != nil {
					return buildErr
				}
				defer idx.Close()
			} else {
				idx, buildErr := search.BuildIfNeededAt(indexDir, hashPath, posts)
				if buildErr != nil {
					return buildErr
				}
				defer idx.Close()
			}

			outlnf("index_dir=%s", indexDir)
			outlnf("hash_path=%s", hashPath)
			outlnf("posts=%d", len(posts))
			return nil
		},
	}

	cmd.Flags().StringVar(&indexDir, "index-dir", "", "directory to write the bleve index")
	cmd.Flags().StringVar(&hashPath, "hash-path", "", "path for the content hash file")
	cmd.Flags().StringVar(&indexName, "index-name", "", "named index suffix inside the default cache directory")
	cmd.Flags().BoolVar(&force, "force", false, "rebuild the index even if the content hash is unchanged")

	return cmd
}

// cliSearchResult pairs a post with its relevance score.
type cliSearchResult struct {
	post  *models.Post
	score float64
}

// searchWithBleve uses the bleve index for BM25-ranked results.
func searchWithBleve(posts []*models.Post, queryStr string, opts search.QueryOptions) ([]cliSearchResult, error) {
	idx, err := search.BuildIfNeeded(".markata/cache", posts)
	if err != nil {
		return nil, err
	}
	defer idx.Close()

	postsByPath := search.PostsByPath(posts)
	hits, err := idx.Search(queryStr, opts, postsByPath)
	if err != nil {
		return nil, err
	}

	results := make([]cliSearchResult, len(hits))
	for i := range hits {
		hit := &hits[i]
		results[i] = cliSearchResult{post: hit.Post, score: hit.Score}
	}
	return results, nil
}

// searchSubstring falls back to simple substring matching.
func searchSubstring(posts []*models.Post, queryStr, fieldsStr string, limit int) []cliSearchResult {
	q := strings.ToLower(queryStr)
	var fieldList []string
	if fieldsStr != "" {
		fieldList = strings.Split(fieldsStr, ",")
	}

	var results []cliSearchResult
	for _, p := range posts {
		if matchesSearchLocal(p, q, fieldList) {
			results = append(results, cliSearchResult{post: p, score: 1.0})
		}
	}
	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}
	return results
}

// matchesSearchLocal checks if a post matches the query, respecting field restrictions.
func matchesSearchLocal(p *models.Post, query string, fields []string) bool {
	searchAll := len(fields) == 0

	if (searchAll || fieldIn(fields, "title")) && strings.Contains(strings.ToLower(p.PlainTitle()), query) {
		return true
	}
	if (searchAll || fieldIn(fields, "description")) && p.Description != nil && strings.Contains(strings.ToLower(*p.Description), query) {
		return true
	}
	if (searchAll || fieldIn(fields, "content")) && strings.Contains(strings.ToLower(p.Content), query) {
		return true
	}
	if searchAll || fieldIn(fields, "tags") {
		for _, tag := range p.Tags {
			if strings.Contains(strings.ToLower(tag), query) {
				return true
			}
		}
	}
	return false
}

func fieldIn(fields []string, name string) bool {
	for _, f := range fields {
		if strings.EqualFold(strings.TrimSpace(f), name) {
			return true
		}
	}
	return false
}

func cliSearchResultPosts(results []cliSearchResult) []*models.Post {
	posts := make([]*models.Post, len(results))
	for i, r := range results {
		posts[i] = r.post
	}
	return posts
}

// sortSearchResults sorts results by the given field and order.
func sortSearchResults(results []cliSearchResult, field string, order services.SortOrder) {
	sort.SliceStable(results, func(i, j int) bool {
		a, b := results[i].post, results[j].post
		var cmp int
		switch strings.ToLower(field) {
		case "date":
			cmp = compareDatePtrs(a.Date, b.Date)
		case "title":
			cmp = strings.Compare(
				strings.ToLower(derefStr(a.Title)),
				strings.ToLower(derefStr(b.Title)),
			)
		case "path":
			cmp = strings.Compare(a.Path, b.Path)
		case "words":
			cmp = compareInts(postWordCount(a), postWordCount(b))
		case "reading_time":
			cmp = compareInts(postReadingTime(a), postReadingTime(b))
		case "tags":
			cmp = strings.Compare(
				strings.ToLower(strings.Join(a.Tags, ", ")),
				strings.ToLower(strings.Join(b.Tags, ", ")),
			)
		}
		if cmp == 0 {
			cmp = strings.Compare(a.Path, b.Path)
		}
		if order == services.SortDesc {
			return cmp > 0
		}
		return cmp < 0
	})
}

func compareDatePtrs(a, b *time.Time) int {
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return -1
	}
	if b == nil {
		return 1
	}
	if a.Before(*b) {
		return -1
	}
	if a.After(*b) {
		return 1
	}
	return 0
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// renderSearchTable renders search results with a match count header and optional scores.
func renderSearchTable(results []cliSearchResult, query string) error {
	outlnf("%d results for %q\n", len(results), query)
	if len(results) == 0 {
		return nil
	}

	// Check if we have meaningful scores (bleve results)
	hasScores := len(results) > 0 && results[0].score != 1.0

	w := tabwriter.NewWriter(outWriter(), 0, 0, 2, ' ', 0)
	if hasScores {
		fmt.Fprintln(w, "SCORE\tTITLE\tDATE\tWORDS\tREAD\tTAGS\tPATH")
		fmt.Fprintln(w, "-----\t-----\t----\t-----\t----\t----\t----")
	} else {
		fmt.Fprintln(w, "TITLE\tDATE\tWORDS\tREAD\tTAGS\tPATH")
		fmt.Fprintln(w, "-----\t----\t-----\t----\t----\t----")
	}
	for _, result := range results {
		row := postToRow(result.post)
		if hasScores {
			fmt.Fprintf(w, "%.3f\t%s\t%s\t%s\t%s\t%s\t%s\n",
				result.score,
				row.Title,
				row.Date,
				formatWordCount(row.Words),
				formatReadingTime(row.ReadingTime),
				strings.Join(row.Tags, ", "),
				row.Path,
			)
		} else {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
				row.Title,
				row.Date,
				formatWordCount(row.Words),
				formatReadingTime(row.ReadingTime),
				strings.Join(row.Tags, ", "),
				row.Path,
			)
		}
	}
	return w.Flush()
}
