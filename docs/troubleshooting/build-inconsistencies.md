---
title: "Troubleshooting Build Inconsistencies"
description: "How to fix issues when builds produce different results on different machines"
date: 2026-01-30
published: true
tags:
  - troubleshooting
  - build
  - cache
---


If you're getting different build results on different machines, follow this diagnostic guide.

## Quick Fix

On the machine with issues, run:

```bash
cd ~/git/markata-go  # or your repo location
./scripts/clean-rebuild.sh
markata-go build --clean
```

## Root Causes

### 1. Go Build Cache (Most Common)

**Problem**: Templates are embedded into the binary at compile time via `//go:embed`. Go's build cache can serve stale compiled objects even after template files change.

**Symptoms**:
- Old templates appear in output
- Changes to `pkg/themes/default/templates/` don't take effect
- Rebuilding with `go build` doesn't help

**Solution**:
```bash
go clean -cache
go install -a ./cmd/markata-go
```

The `-a` flag forces Go to rebuild **all** packages, ignoring the cache.

### 2. Markata Build Cache

**Problem**: The `.markata/` directory caches rendered HTML to speed up incremental builds.

Unchanged pages, including published frontmatter-only pages, restore their
cached full-page HTML before the write stage. If you recreate the output
directory while keeping `.markata/`, a normal or `--fast` build can materialize
those pages again without clearing the build cache.

The cache also hashes the embedded default theme templates. A new builder image
that changes those templates invalidates cached pages on its next build.

**Symptoms**:
- Removing templates but still seeing cached output
- Template changes not reflected in output

**Solution**:
```bash
rm -rf .markata
markata-go build --clean
```

Current releases should still detect moved or deleted content files during `markata-go build --fast`.
If a moved-file error only goes away after clearing `.markata/`, you are likely hitting an older
build-cache bug rather than expected behavior.

### 3. Config File Differences

**Problem**: `markata-go.toml` points to wrong templates or has different settings.

**Example Issue**: Config pointed to old `partials/card.html` instead of new `partials/cards/card-router.html`.

**Solution**:
- Check git status: `git status`
- Compare configs: `diff markata-go.toml <(git show HEAD:markata-go.toml)`
- Use diagnostic script: `./scripts/diagnose-build.sh`

### 4. Different Git Commits

**Problem**: Machines are on different branches or commits.

**Symptoms**:
- Binary checksums differ
- Template checksums differ

**Solution**:
```bash
git status
git log -1
git pull origin main  # or your branch
```

## Diagnostic Tools

### Full Diagnostic

Run this on both machines and compare output:

```bash
./scripts/diagnose-build.sh > ~/diagnostic.txt
cat ~/diagnostic.txt
```

**Key things to compare**:
- Git commit hash
- Binary SHA256 checksum
- Template checksums (embedded vs project overrides)
- Config file checksum

### Clean Rebuild

Use the automated clean rebuild script:

```bash
./scripts/clean-rebuild.sh
```

This script:
1. Checks git status
2. Cleans Go cache
3. Removes old binaries
4. Removes markata caches
5. Shows template checksums
6. Rebuilds and installs binary
7. Verifies installation

## Expected Checksums (as of commit bfa19d0)

**Binary**:
```
SHA256: 72393651937484c022ced092a80f5171353d8bc49320453a58a187ed1b168bcb
```

**Key Templates** (embedded in `pkg/themes/default/templates/`):
```
article-card.html: b9a26a3dbbec42e00aaea791772de29c4c1060cb1cd6047b4a8bf9e30e158d23
feed.html:         b4a5cab1c989446d62240e9e98126dc2dfc62accadc86775a2757c69788ad63a
card-router.html:  12c02e186b2ca95f9c65187f91661a2da994984895bcaf4d2cfea89658835627
```

**Config**:
```
markata-go.toml: a01e140d4788d732dddc1eb90f57b183781ebcbcbda78c90389eadf3e0691c9d
```

## Card System

### Why No Card Types Appear

The card router (`partials/cards/card-router.html`) selects card templates based on the post's `template` field:

```yaml
---
title: My Note
template: note  # This determines the card type!
---
```

**Card types**:
- `article`, `blog-post`, `post`, `essay`, `tutorial` → article-card.html
- `note`, `ping`, `thought`, `status`, `tweet` → note-card.html
- `photo`, `shot`, `shots`, `image`, `gallery` → photo-card.html
- `video`, `clip`, `cast`, `stream` → video-card.html
- `link`, `bookmark`, `stars` → link-card.html
- `quote`, `quotation` → quote-card.html
- `guide`, `series`, `step`, `chapter` → guide-card.html
- `gratitude`, `inline`, `micro`, `til` → inline-card.html
- (no template field) → default-card.html

**To get different card types**, add `template:` to your post frontmatter:

```yaml
---
title: "Quick Thought"
template: note
date: 2026-01-30
---
```

## Workflow for Updates

When you update templates:

```bash
# 1. Edit templates in pkg/themes/default/templates/
vim pkg/themes/default/templates/feed.html

# 2. Force rebuild to pick up embedded file changes
go clean -cache
go install -a ./cmd/markata-go

# 3. Clean build
rm -rf .markata public
markata-go build --clean
```

## Related Files

- `pkg/themes/embed.go` - Embeds default templates into binary
- `pkg/templates/engine.go` - Template loading order
- `markata-go.toml` - Config including feed card template
- `templates/partials/cards/card-router.html` - Card type router

## See Also

- [Templates Guide](/docs/guides/templates/)
- [Feeds Guide](/docs/guides/feeds/)
- [Frontmatter Reference](/docs/guides/frontmatter/)
