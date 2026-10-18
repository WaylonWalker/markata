// Command gen writes pkg/config/settings_docs_gen.go from pkg/models doc
// comments. Run it with `go generate ./pkg/config`.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/WaylonWalker/markata/internal/configdocs"
)

func main() {
	models := flag.String("models", "../models", "directory containing the models package")
	out := flag.String("out", "settings_docs_gen.go", "output file")
	pkg := flag.String("pkg", "config", "package name of the output file")
	flag.Parse()

	docs, err := configdocs.Extract(*models)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	src, err := configdocs.Render(*pkg, docs)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := os.WriteFile(*out, src, 0o600); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
