package plugins

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/templates"
	"golang.org/x/net/html"
)

// A versioned asset is immutable. Previous releases may still refer to it,
// so purge/minification must only process the unversioned working stylesheet.
var versionedCSSFilename = regexp.MustCompile(`\.[0-9a-fA-F]{8}\.css$`)

func isVersionedCSSFile(path string) bool {
	return versionedCSSFilename.MatchString(filepath.Base(path))
}

type finalCSSReplacement struct {
	matcher     *regexp.Regexp
	destination string
}

// rewriteCSSLinkReferences only edits actual <link href> attributes. Scanning
// the HTML as tokens preserves authored prose, code examples, and comments
// that happen to mention an older versioned CSS filename.
func rewriteCSSLinkReferences(input []byte, changes []finalCSSReplacement) ([]byte, error) {
	tokenizer := html.NewTokenizer(bytes.NewReader(input))
	var output bytes.Buffer
	for {
		kind := tokenizer.Next()
		if kind == html.ErrorToken {
			if errors.Is(tokenizer.Err(), io.EOF) {
				return output.Bytes(), nil
			}
			return nil, fmt.Errorf("tokenize HTML for CSS references: %w", tokenizer.Err())
		}
		raw := tokenizer.Raw()
		if kind != html.StartTagToken && kind != html.SelfClosingTagToken {
			output.Write(raw)
			continue
		}
		token := tokenizer.Token()
		if token.Data != "link" {
			output.Write(raw)
			continue
		}
		changed := false
		for i := range token.Attr {
			if token.Attr[i].Key != "href" {
				continue
			}
			for _, change := range changes {
				updated := change.matcher.ReplaceAllString(token.Attr[i].Val, change.destination)
				if updated != token.Attr[i].Val {
					token.Attr[i].Val = updated
					changed = true
				}
			}
		}
		if changed {
			output.WriteString(token.String())
		} else {
			output.Write(raw)
		}
	}
}

// reconcileFinalCSSAssets reconciles final published CSS with already-rendered
// HTML after purge and minification, including HTML restored from build cache.
func reconcileFinalCSSAssets(m *lifecycle.Manager, outputFS *os.Root, outputDir string) error {
	hashes := m.AssetHashes()
	paths := make([]string, 0, len(hashes))
	for path := range hashes {
		if strings.EqualFold(filepath.Ext(path), ".css") {
			paths = append(paths, path)
		}
	}
	sort.Strings(paths)
	changes := make([]finalCSSReplacement, 0, len(paths))
	finalHashes := make(map[string]string)
	for _, path := range paths {
		relative, err := outputRelativePath(outputDir, filepath.Join(outputDir, path))
		if err != nil {
			return err
		}
		data, err := outputFS.ReadFile(relative)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return fmt.Errorf("reading final stylesheet %s: %w", path, err)
		}
		hash := fmt.Sprintf("%x", sha256.Sum256(data))[:8]
		ext := filepath.Ext(path)
		base := strings.TrimSuffix(filepath.ToSlash(path), ext)
		matcher, err := regexp.Compile(regexp.QuoteMeta(base) + `\.[0-9a-fA-F]{8}` + regexp.QuoteMeta(ext))
		if err != nil {
			return err
		}
		changes = append(changes, finalCSSReplacement{matcher: matcher, destination: base + "." + hash + ext})
		if hash != hashes[path] {
			finalHashes[path] = hash
		}
	}
	if len(changes) == 0 {
		return nil
	}

	// Cached HTML may reference either the original Configure hash or a final
	// hash from an earlier build. Normalize every matching URL, not only the
	// source-hash reference from this build.
	err := filepath.WalkDir(outputDir, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		if entry.Type()&os.ModeSymlink != 0 || !strings.EqualFold(filepath.Ext(path), ".html") {
			return nil
		}
		relative, err := outputRelativePath(outputDir, path)
		if err != nil {
			return err
		}
		data, err := outputFS.ReadFile(relative)
		if err != nil {
			return err
		}
		after, err := rewriteCSSLinkReferences(data, changes)
		if err != nil {
			return fmt.Errorf("reconcile CSS links in %s: %w", relative, err)
		}
		if bytes.Equal(data, after) {
			return nil
		}
		return outputFS.WriteFile(relative, after, 0o644)
	})
	if err != nil {
		return fmt.Errorf("updating HTML CSS URLs: %w", err)
	}
	for path, hash := range finalHashes {
		m.SetAssetHash(path, hash)
	}
	templates.SetAssetHashes(finalHashes)
	return nil
}
