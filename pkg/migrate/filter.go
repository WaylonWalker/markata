package migrate

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/WaylonWalker/markata/pkg/filter"
)

// Filter migrates a Python markata filter expression to markata-go syntax.
// It returns the migrated expression and a list of changes made.
func Filter(expr string) (migrated string, changes []string) {
	if expr == "" {
		return expr, nil
	}

	migrated = expr

	// 1. Migrate boolean literals: 'True'/'False' -> True/False
	migrated, boolChanges := migrateBooleanLiterals(migrated)
	changes = append(changes, boolChanges...)

	// 2. Migrate 'in' operator with lists
	migrated, inChanges := migrateInOperator(migrated)
	changes = append(changes, inChanges...)

	// 3. Fix operator spacing
	migrated, spaceChanges := fixOperatorSpacing(migrated)
	changes = append(changes, spaceChanges...)

	// 4. Migrate None comparisons
	migrated, noneChanges := migrateNoneComparisons(migrated)
	changes = append(changes, noneChanges...)

	return migrated, changes
}

// ValidateFilter validates a filter expression using the markata-go filter parser.
func ValidateFilter(expr string) error {
	if expr == "" {
		return nil
	}

	lexer := filter.NewLexer(expr)
	_, err := lexer.Tokenize()
	if err != nil {
		return fmt.Errorf("invalid filter expression: %w", err)
	}

	return nil
}

// migrateBooleanLiterals converts quoted boolean strings to unquoted.
// Examples:
//   - published == 'True' -> published == True
//   - draft == "false" -> draft == False
func migrateBooleanLiterals(expr string) (result string, changes []string) {
	result = expr

	// Patterns for quoted booleans (Go regex doesn't support backreferences)
	patterns := []struct {
		pattern     *regexp.Regexp
		replacement string
		description string
	}{
		{
			pattern:     regexp.MustCompile(`'True'`),
			replacement: "True",
			description: "Boolean literal: 'True' -> True",
		},
		{
			pattern:     regexp.MustCompile(`"True"`),
			replacement: "True",
			description: "Boolean literal: \"True\" -> True",
		},
		{
			pattern:     regexp.MustCompile(`'False'`),
			replacement: "False",
			description: "Boolean literal: 'False' -> False",
		},
		{
			pattern:     regexp.MustCompile(`"False"`),
			replacement: "False",
			description: "Boolean literal: \"False\" -> False",
		},
		{
			pattern:     regexp.MustCompile(`'true'`),
			replacement: "True",
			description: "Boolean literal: 'true' -> True",
		},
		{
			pattern:     regexp.MustCompile(`"true"`),
			replacement: "True",
			description: "Boolean literal: \"true\" -> True",
		},
		{
			pattern:     regexp.MustCompile(`'false'`),
			replacement: "False",
			description: "Boolean literal: 'false' -> False",
		},
		{
			pattern:     regexp.MustCompile(`"false"`),
			replacement: "False",
			description: "Boolean literal: \"false\" -> False",
		},
	}

	for _, p := range patterns {
		if p.pattern.MatchString(result) {
			result = p.pattern.ReplaceAllString(result, p.replacement)
			changes = append(changes, p.description)
		}
	}

	return result, changes
}

// migrateInOperator converts 'in' operator with lists to 'or' expressions.
// Example: templateKey in ['blog-post', 'til'] -> templateKey == 'blog-post' or templateKey == 'til'
func migrateInOperator(expr string) (result string, changes []string) {
	// Pattern to match: identifier in [list]
	// This handles: field in ['a', 'b', 'c'] or field in ["a", "b", "c"]
	inPattern := regexp.MustCompile(`(\w+)\s+in\s+\[((?:[^]]+))\]`)

	matches := inPattern.FindAllStringSubmatch(expr, -1)
	if len(matches) == 0 {
		return expr, changes
	}

	result = expr
	for _, match := range matches {
		if len(match) < 3 {
			continue
		}

		field := match[1]
		listContent := match[2]

		// Parse list items (handle both single and double quotes)
		items := parseListItems(listContent)
		if len(items) == 0 {
			continue
		}

		// Build or expression
		var orParts []string
		for _, item := range items {
			orParts = append(orParts, fmt.Sprintf("%s == %s", field, item))
		}

		replacement := strings.Join(orParts, " or ")

		// Replace in the expression
		fullMatch := match[0]
		result = strings.Replace(result, fullMatch, replacement, 1)

		changes = append(changes, fmt.Sprintf("in operator: %s -> %s", fullMatch, replacement))
	}

	return result, changes
}

// parseListItems parses items from a list string like "'a', 'b', 'c'" or '"a", "b"'.
func parseListItems(listContent string) []string {
	var items []string

	// Match single-quoted strings
	singleQuotePattern := regexp.MustCompile(`'([^']*)'`)
	singleMatches := singleQuotePattern.FindAllStringSubmatch(listContent, -1)
	for _, match := range singleMatches {
		if len(match) >= 2 {
			items = append(items, fmt.Sprintf("'%s'", match[1]))
		}
	}

	// If no single-quoted items, try double-quoted
	if len(items) == 0 {
		doubleQuotePattern := regexp.MustCompile(`"([^"]*)"`)
		doubleMatches := doubleQuotePattern.FindAllStringSubmatch(listContent, -1)
		for _, match := range doubleMatches {
			if len(match) >= 2 {
				items = append(items, fmt.Sprintf("'%s'", match[1]))
			}
		}
	}

	return items
}

// fixOperatorSpacing ensures operators have surrounding whitespace.
// Examples:
//   - date<=today -> date <= today
//   - count>=10 -> count >= 10
func fixOperatorSpacing(expr string) (result string, changes []string) {
	result = expr

	// Operators that need spacing
	operators := []string{"<=", ">=", "==", "!=", "<", ">"}

	for _, op := range operators {
		// Pattern for operator without proper spacing (not inside quotes)
		// We need to be careful not to modify operators inside string literals

		// First, handle cases where there's no space before the operator
		pattern := regexp.MustCompile(`(\w)` + regexp.QuoteMeta(op))
		if pattern.MatchString(result) {
			newResult := pattern.ReplaceAllString(result, "${1} "+op)
			if newResult != result {
				changes = append(changes, fmt.Sprintf("Added space before '%s'", op))
				result = newResult
			}
		}

		// Handle cases where there's no space after the operator
		pattern = regexp.MustCompile(regexp.QuoteMeta(op) + `(\w)`)
		if pattern.MatchString(result) {
			newResult := pattern.ReplaceAllString(result, op+" ${1}")
			if newResult != result {
				changes = append(changes, fmt.Sprintf("Added space after '%s'", op))
				result = newResult
			}
		}
	}

	// Clean up any double spaces we might have created
	result = regexp.MustCompile(`\s{2,}`).ReplaceAllString(result, " ")

	return result, changes
}

// migrateNoneComparisons converts Python-style None comparisons.
// Examples:
//   - image is None -> image == None
//   - image is not None -> image != None
func migrateNoneComparisons(expr string) (result string, changes []string) {
	result = expr

	// is not None -> != None
	isNotNonePattern := regexp.MustCompile(`(\w+)\s+is\s+not\s+None`)
	if isNotNonePattern.MatchString(result) {
		result = isNotNonePattern.ReplaceAllString(result, "${1} != None")
		changes = append(changes, "None comparison: 'is not None' -> '!= None'")
	}

	// is None -> == None
	isNonePattern := regexp.MustCompile(`(\w+)\s+is\s+None`)
	if isNonePattern.MatchString(result) {
		result = isNonePattern.ReplaceAllString(result, "${1} == None")
		changes = append(changes, "None comparison: 'is None' -> '== None'")
	}

	return result, changes
}

// AnalyzeFilter analyzes a filter expression and returns potential issues.
func AnalyzeFilter(expr string) []Warning {
	var warnings []Warning

	// Check for Python-specific patterns that might cause issues
	if strings.Contains(expr, "lambda") {
		warnings = append(warnings, Warning{
			Category:   "filter",
			Message:    "Lambda expressions are not supported",
			Path:       expr,
			Suggestion: "Rewrite the filter without lambda expressions",
		})
	}

	if strings.Contains(expr, ".lower()") || strings.Contains(expr, ".upper()") {
		warnings = append(warnings, Warning{
			Category:   "filter",
			Message:    "String methods like .lower()/.upper() are not supported",
			Path:       expr,
			Suggestion: "Use case-insensitive comparisons or pre-process data",
		})
	}

	if strings.Contains(expr, "len(") {
		warnings = append(warnings, Warning{
			Category:   "filter",
			Message:    "len() function is not supported",
			Path:       expr,
			Suggestion: "Use a dedicated field for length if needed",
		})
	}

	return warnings
}
