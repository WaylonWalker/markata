# Issue #001: Batch Parse Error Reporting

## Summary

When loading markdown files, the build currently fails on the first parsing error encountered. Users with multiple malformed files must fix them one at a time, running the build repeatedly to discover each error.

## Current Behavior

```
$ markata-go build -v
Verbose mode enabled
Starting build...
Configuration loaded (output: public, patterns: [posts/**/*.md pages/**/*.md])
Error: build failed: [error] load plugin "load": failed to parse pages/blog/about.md: invalid frontmatter: yaml: unmarshal errors:
  line 4: mapping key "tags" already defined at line 2
```

The build stops at the first error. If there are 10 files with issues, the user must run the build 10 times to discover all of them.

## Desired Behavior

Report all parsing errors at once, then fail:

```
$ markata-go build -v
Verbose mode enabled
Starting build...
Configuration loaded (output: public, patterns: [posts/**/*.md pages/**/*.md])
  [configure] running...
  [validate] running...
  [glob] running...
  [glob] discovered 25 files
  [load] running...

Parse errors found in 3 files:

  pages/blog/about.md:
    - invalid frontmatter: yaml: unmarshal errors:
        line 4: mapping key "tags" already defined at line 2

  posts/draft-post.md:
    - invalid date: unable to parse date: "not-a-date"

  posts/old-post.md:
    - failed to read file: permission denied

Error: build failed: 3 files failed to parse (see errors above)
```

## Implementation Notes

### Option A: Collect errors, fail at end (Recommended)

Modify `LoadPlugin.Load()` to:
1. Continue processing all files even when errors occur
2. Collect all errors in a slice
3. After processing all files, if any errors occurred:
   - Print a summary of all errors with file paths and line numbers
   - Return a combined error

```go
func (p *LoadPlugin) Load(m *lifecycle.Manager) error {
    files := m.Files()
    var parseErrors []ParseError

    for _, file := range files {
        post, err := p.loadFile(file)
        if err != nil {
            parseErrors = append(parseErrors, ParseError{
                File:  file,
                Error: err,
            })
            continue // Continue to next file
        }
        m.AddPost(post)
    }

    if len(parseErrors) > 0 {
        // Print detailed error report
        printParseErrors(parseErrors)
        return fmt.Errorf("%d files failed to parse", len(parseErrors))
    }

    return nil
}
```

### Option B: Add `--fail-fast` flag

Keep current behavior as default but add a flag:
- `--fail-fast` (or `-f`): Stop on first error (current behavior)
- Default: Collect and report all errors

### Option C: Warning mode with `--strict`

- Default: Log parse errors as warnings, skip problematic files, continue build
- `--strict`: Fail on any parse error (current behavior)

## Additional Considerations

1. **Error formatting**: Use consistent formatting that's easy to scan:
   - Group by file
   - Include line numbers where available
   - Use color coding if terminal supports it (red for errors)

2. **Exit code**: Should still be non-zero when errors occur

3. **Partial builds**: Consider allowing partial builds where valid files are processed and only problematic ones are skipped (with warnings)

4. **JSON output**: Consider `--format=json` for CI/tooling integration:
   ```json
   {
     "errors": [
       {"file": "pages/about.md", "line": 4, "message": "duplicate key 'tags'"}
     ]
   }
   ```

## Files to Modify

- `pkg/plugins/load.go` - Main implementation
- `cmd/markata-go/cmd/build.go` - Add flags if implementing Option B/C
- `pkg/plugins/frontmatter.go` - Enhance error messages with line numbers

## Related

- Error reporting was previously silent due to `SilenceErrors: true` in cobra (fixed)
- Verbose mode now shows stage-by-stage progress

## Priority

Medium - Quality of life improvement for content authors

## Labels

`enhancement`, `ux`, `error-handling`
