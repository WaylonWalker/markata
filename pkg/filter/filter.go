package filter

import (
	"fmt"
	"sync"

	"github.com/WaylonWalker/markata/pkg/models"
)

// Filter represents a compiled filter expression
type Filter struct {
	expression string
	ast        Expr
	context    *EvalContext
}

// astCache caches parsed ASTs to avoid re-parsing identical expressions.
// The AST is immutable after parsing, so this is safe for concurrent use.
var astCache = struct {
	sync.RWMutex
	m map[string]Expr
}{m: make(map[string]Expr)}

// Parse parses a filter expression and returns a Filter.
// Parsed ASTs are cached to avoid re-parsing identical expressions.
func Parse(expression string) (*Filter, error) {
	// Check cache first
	astCache.RLock()
	cachedAST, ok := astCache.m[expression]
	astCache.RUnlock()

	if ok {
		// Return new Filter with cached AST and fresh context
		return &Filter{
			expression: expression,
			ast:        cachedAST,
			context:    NewEvalContext(),
		}, nil
	}

	// Parse and cache
	ast, err := ParseExpression(expression)
	if err != nil {
		return nil, fmt.Errorf("failed to parse filter expression %q: %w", expression, err)
	}

	// Cache the AST
	astCache.Lock()
	astCache.m[expression] = ast
	astCache.Unlock()

	return &Filter{
		expression: expression,
		ast:        ast,
		context:    NewEvalContext(),
	}, nil
}

// MustParse parses a filter expression and panics on error
func MustParse(expression string) *Filter {
	f, err := Parse(expression)
	if err != nil {
		panic(err)
	}
	return f
}

// Expression returns the original filter expression string
func (f *Filter) Expression() string {
	return f.expression
}

// SetContext sets a custom evaluation context
func (f *Filter) SetContext(ctx *EvalContext) {
	f.context = ctx
}

// RefreshContext refreshes the context with current time values
func (f *Filter) RefreshContext() {
	f.context = NewEvalContext()
}

// Match evaluates the filter against a single post
func (f *Filter) Match(post *models.Post) (bool, error) {
	if f.ast == nil {
		return true, nil
	}
	return Evaluate(f.ast, post, f.context)
}

// MustMatch evaluates the filter against a single post and panics on error
func (f *Filter) MustMatch(post *models.Post) bool {
	result, err := f.Match(post)
	if err != nil {
		panic(err)
	}
	return result
}

// MatchAll filters a slice of posts and returns only those that match
func (f *Filter) MatchAll(posts []*models.Post) []*models.Post {
	result := make([]*models.Post, 0, len(posts))
	for _, post := range posts {
		if match, err := f.Match(post); err == nil && match {
			result = append(result, post)
		}
	}
	return result
}

// MatchAllWithErrors filters a slice of posts and returns matches along with any errors
func (f *Filter) MatchAllWithErrors(posts []*models.Post) ([]*models.Post, []error) {
	var result []*models.Post
	var errors []error

	for _, post := range posts {
		match, err := f.Match(post)
		if err != nil {
			errors = append(errors, fmt.Errorf("error evaluating post %s: %w", post.Path, err))
			continue
		}
		if match {
			result = append(result, post)
		}
	}

	return result, errors
}

// Posts is a convenience function that parses an expression and filters posts
func Posts(expression string, posts []*models.Post) ([]*models.Post, error) {
	f, err := Parse(expression)
	if err != nil {
		return nil, err
	}
	return f.MatchAll(posts), nil
}

// MatchPost is a convenience function that parses an expression and matches a single post
func MatchPost(expression string, post *models.Post) (bool, error) {
	f, err := Parse(expression)
	if err != nil {
		return false, err
	}
	return f.Match(post)
}

// And combines multiple filters with AND logic
func And(filters ...*Filter) *Filter {
	if len(filters) == 0 {
		return &Filter{
			expression: "True",
			ast:        &Literal{Value: true},
			context:    NewEvalContext(),
		}
	}

	if len(filters) == 1 {
		return filters[0]
	}

	// Build combined AST
	combined := filters[0].ast
	expression := filters[0].expression

	for i := 1; i < len(filters); i++ {
		combined = &BinaryExpr{
			Left:  combined,
			Op:    opAnd,
			Right: filters[i].ast,
		}
		expression += " and " + filters[i].expression
	}

	return &Filter{
		expression: expression,
		ast:        combined,
		context:    NewEvalContext(),
	}
}

// Or combines multiple filters with OR logic
func Or(filters ...*Filter) *Filter {
	if len(filters) == 0 {
		return &Filter{
			expression: "False",
			ast:        &Literal{Value: false},
			context:    NewEvalContext(),
		}
	}

	if len(filters) == 1 {
		return filters[0]
	}

	// Build combined AST
	combined := filters[0].ast
	expression := filters[0].expression

	for i := 1; i < len(filters); i++ {
		combined = &BinaryExpr{
			Left:  combined,
			Op:    opOr,
			Right: filters[i].ast,
		}
		expression += " or " + filters[i].expression
	}

	return &Filter{
		expression: expression,
		ast:        combined,
		context:    NewEvalContext(),
	}
}

// Not creates a negated filter
func Not(f *Filter) *Filter {
	return &Filter{
		expression: "not (" + f.expression + ")",
		ast: &UnaryExpr{
			Op:   "not",
			Expr: f.ast,
		},
		context: NewEvalContext(),
	}
}

// Always returns a filter that always matches
func Always() *Filter {
	return &Filter{
		expression: "True",
		ast:        &Literal{Value: true},
		context:    NewEvalContext(),
	}
}

// Never returns a filter that never matches
func Never() *Filter {
	return &Filter{
		expression: "False",
		ast:        &Literal{Value: false},
		context:    NewEvalContext(),
	}
}
