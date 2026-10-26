# Versioning and Release Specification

This document specifies the versioning scheme and release process for markata-go.

## Versioning Scheme

markata-go follows [Semantic Versioning 2.0.0](https://semver.org/):

```
MAJOR.MINOR.PATCH[-PRERELEASE][+BUILD]
```

### Version Components

| Component | Description | Example |
|-----------|-------------|---------|
| MAJOR | Breaking changes to CLI, config, or public API | `1.0.0` → `2.0.0` |
| MINOR | New features, backward-compatible | `0.1.0` → `0.2.0` |
| PATCH | Bug fixes, backward-compatible | `0.1.0` → `0.1.1` |
| PRERELEASE | Pre-release identifier | `0.1.0-alpha.1`, `0.1.0-beta.2`, `0.1.0-rc.1` |

### Zero-Version (0.x.y)

During initial development (version `0.x.y`):

- The public API is not considered stable
- MINOR version bumps may include breaking changes
- PATCH versions are for bug fixes only
- Users should expect rapid iteration

### Stability Promise

Once `1.0.0` is released:

- PATCH releases never break existing functionality
- MINOR releases are backward-compatible
- MAJOR releases may include breaking changes (with migration guides)

## Version Information

### Embedded Version Data

The binary embeds version information via ldflags at build time:

```go
// cmd/markata-go/cmd/version.go
var (
    Version = "dev"     // Semantic version (e.g., "0.1.0")
    Commit  = "none"    // Git commit SHA
    Date    = "unknown" // Build timestamp (RFC3339)
)
```

Every build entry point MUST target these variables using their canonical Go
import path, `github.com/WaylonWalker/markata/cmd/markata-go/cmd`. This includes
the Docker builder image and local `just` builds. A path from the retired
`github.com/WaylonWalker/markata-go` module is not valid; Go silently ignores
linker flags that do not match a package variable. Production binaries MUST
report the supplied version, commit, and build date through `markata-go
version`.

### Version Command

```bash
$ markata-go version
markata-go 0.1.0
  commit:    abc1234
  built:     2024-01-15T12:00:00Z
  go:        go1.22.2
  os/arch:   linux/amd64

$ markata-go version --short
0.1.0

$ markata-go --version
markata-go version 0.1.0
```

## Release Process

### Git Tags

Releases are triggered by Git tags matching `v*`:

```bash
# Create annotated tag
git tag -a v0.1.0 -m "Release v0.1.0"

# Push tag to trigger release
git push origin v0.1.0
```

### Pre-release Tags

Pre-release versions use suffixes:

```bash
git tag -a v0.1.0-alpha.1 -m "Alpha release"
git tag -a v0.1.0-beta.1 -m "Beta release"
git tag -a v0.1.0-rc.1 -m "Release candidate"
```

### Release Workflow

1. **Version Bump** - Update any version references in docs
2. **Changelog** - Update CHANGELOG.md (if manual)
3. **Tag** - Create annotated Git tag
4. **Push** - Push tag to trigger GitHub Actions
5. **Verify** - Check release artifacts on GitHub

### Automated Release

The GitHub Actions release workflow:

1. Runs tests
2. Builds binaries for all platforms via GoReleaser
3. Creates GitHub release with:
   - Platform-specific archives
   - Checksums file
   - SBOM (Software Bill of Materials)
   - Signed artifacts (via cosign)
4. Generates changelog from conventional commits

## Build Artifacts

### Supported Platforms

| OS | Architecture | Binary Name |
|----|--------------|-------------|
| Linux | amd64 | `markata-go` |
| Linux | arm64 | `markata-go` |
| Linux | armv7 | `markata-go` |
| macOS | amd64 (Intel) | `markata-go` |
| macOS | arm64 (Apple Silicon) | `markata-go` |
| Windows | amd64 | `markata-go.exe` |
| FreeBSD | amd64 | `markata-go` |

### Archive Naming

```
markata-go_{version}_{os}_{arch}.{ext}
```

Examples:
- `markata-go_0.1.0_linux_x86_64.tar.gz`
- `markata-go_0.1.0_darwin_arm64.tar.gz`
- `markata-go_0.1.0_windows_x86_64.zip`

### Release Assets

Each release includes:

| File | Description |
|------|-------------|
| `markata-go_*_{os}_{arch}.{tar.gz,zip}` | Platform binary + README |
| `checksums.txt` | SHA256 checksums for all archives |
| `checksums.txt.pem` | Cosign certificate |
| `checksums.txt.sig` | Cosign signature |
| `markata-go_*_{os}_{arch}.tar.gz.sbom.json` | SBOM in SPDX format |

## Installation Methods

### Supported Installers

| Method | Command |
|--------|---------|
| jpillora/installer | `curl -sL https://i.jpillora.com/WaylonWalker/markata-go \| bash` |
| eget | `eget WaylonWalker/markata-go` |
| mise | `mise use -g github:WaylonWalker/markata-go` |
| go install | `go install github.com/WaylonWalker/markata-go/cmd/markata-go@latest` |

### Compatibility Requirements

For jpillora/installer and eget compatibility:

1. Archive naming must follow `{project}_{version}_{os}_{arch}` pattern
2. Binary must be at archive root (not in subdirectory)
3. Checksums file must be named `checksums.txt`
4. GitHub releases must be public

## Local Development

### Building with Version Info

```bash
# Using just
just build

# Manual with ldflags
go build -ldflags "-s -w \
  -X github.com/WaylonWalker/markata/cmd/markata-go/cmd.Version=$(git describe --tags) \
  -X github.com/WaylonWalker/markata/cmd/markata-go/cmd.Commit=$(git rev-parse --short HEAD) \
  -X github.com/WaylonWalker/markata/cmd/markata-go/cmd.Date=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
  -o markata-go ./cmd/markata-go
```

### Snapshot Builds

Test the full release process locally:

```bash
# Using just
just snapshot

# Using goreleaser directly
goreleaser release --snapshot --clean --skip=publish
```

## Changelog

Changelog is auto-generated from conventional commits:

| Prefix | Section |
|--------|---------|
| `feat:` | New Features |
| `fix:` | Bug Fixes |
| `perf:` | Performance |
| `docs:` | (excluded) |
| `test:` | (excluded) |
| `ci:` | (excluded) |
| `chore:` | (excluded) |

### Commit Message Format

```
type(scope): description

[optional body]

[optional footer]
```

Examples:
```
feat(feeds): add JSON Feed 1.1 support
fix(template): handle nil values in conditionals
perf(render): parallelize markdown processing
docs(readme): update installation instructions
```
