package plugins

import (
	"sync"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
)

// Chroma searches every filename pattern when Get does not recognize a fence
// language. Mermaid and Chart.js use the plain-code fallback before their
// plugins transform the HTML, so give those names a constant-time lookup.
const (
	featureFenceMermaid = "mermaid"
	featureFenceChartJS = "chartjs"
)

var registerFeatureFenceLexersOnce sync.Once

type plainFeatureFenceLexer struct {
	config chroma.Config
}

func (l *plainFeatureFenceLexer) Config() *chroma.Config { return &l.config }

func (l *plainFeatureFenceLexer) Tokenise(_ *chroma.TokeniseOptions, source string) (chroma.Iterator, error) {
	return chroma.Literator(chroma.Token{Type: chroma.Text, Value: source}), nil
}

func (l *plainFeatureFenceLexer) SetRegistry(*chroma.LexerRegistry) chroma.Lexer { return l }

func (l *plainFeatureFenceLexer) SetAnalyser(func(string) float32) chroma.Lexer { return l }

func (l *plainFeatureFenceLexer) AnalyseText(string) float32 { return 0 }

func registerFeatureFenceLexers() {
	registerFeatureFenceLexersOnce.Do(func() {
		lexers.Register(&plainFeatureFenceLexer{config: chroma.Config{Name: "Markata Mermaid fence", Aliases: []string{featureFenceMermaid}}})
		lexers.Register(&plainFeatureFenceLexer{config: chroma.Config{Name: "Markata Chart.js fence", Aliases: []string{featureFenceChartJS}}})
	})
}
