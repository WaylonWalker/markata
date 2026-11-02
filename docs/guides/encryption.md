---
title: "Encryption and Private Posts"
description: "Protect sensitive content with client-side encryption using AES-256-GCM"
date: 2024-01-15
published: true
slug: /docs/guides/encryption/
tags:
  - documentation
  - encryption
  - privacy
  - security
---


markata-go can encrypt post content so that only visitors with the correct password can read it. Encryption uses AES-256-GCM and runs entirely client-side -- your server never sees the decrypted content.

> **Key principle:** Private posts are **never published in plaintext**. If the encryption key is missing, the build fails rather than exposing your content.

## Quick Start

1. Create a `.env` file in your project root:

```bash
MARKATA_GO_ENCRYPTION_KEY_DEFAULT=your-secret-password
```

2. Mark a post as private in its frontmatter:

```yaml
---
title: My Private Post
private: true
---
Your secret content here.
```

3. Build your site:

```bash
markata-go build
```

The post's HTML content is replaced with an encrypted blob and a password prompt. Visitors enter the password in their browser to decrypt and view the content.

## How It Works

1. During the build, the encryption plugin finds all posts with `private: true`
2. It encrypts the rendered HTML using AES-256-GCM with a password derived via PBKDF2
3. The encrypted content is embedded as a base64 string in a `data-encrypted` attribute
4. A password form is rendered in place of the content
5. Client-side JavaScript uses the Web Crypto API to decrypt on password entry

**The build fails if any private post cannot be encrypted.** This prevents accidentally publishing sensitive content.

## Configuration

Encryption is **enabled by default** with `default_key = "default"`. You only need to set the environment variable.

### Config File Options

```toml
[encryption]
enabled = true                           # default: true
default_key = "default"                  # default: "default"
enforce_strength = true                  # default: true
min_estimated_crack_time = "10y"        # default: "10y"
min_password_length = 14                  # default: 14
decryption_hint = "DM me for access"     # optional hint shown to visitors

[encryption.private_tags]
diary = "personal"                       # tag "diary" encrypts with key "personal"
draft-ideas = "default"                  # tag "draft-ideas" encrypts with key "default"
```

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `enabled` | bool | `true` | Enable/disable encryption processing |
| `default_key` | string | `"default"` | Key name used when a post doesn't specify one |
| `decryption_hint` | string | `""` | Help text shown next to the password prompt |
| `private_tags` | map | `{}` | Maps tag names to encryption key names |
| `enforce_strength` | bool | `true` | Require keys to meet the configured strength policy before encrypting private posts |
| `min_estimated_crack_time` | string | `"10y"` | Minimum estimated crack time for each key (supports `y`, `d`, `h`, `m`, `s`) |
| `min_password_length` | int | `14` | Minimum password length required for every encryption key |

### Environment Variables

Encryption keys are loaded from environment variables with the prefix `MARKATA_GO_ENCRYPTION_KEY_`:

```bash
# The default key (used when no specific key is set on a post)
MARKATA_GO_ENCRYPTION_KEY_DEFAULT=my-password

# Named keys for different access levels
MARKATA_GO_ENCRYPTION_KEY_PERSONAL=personal-password
MARKATA_GO_ENCRYPTION_KEY_PREMIUM=premium-password
```

You can also override config options via environment:

```bash
MARKATA_GO_ENCRYPTION_ENABLED=true
MARKATA_GO_ENCRYPTION_DEFAULT_KEY=default
MARKATA_GO_ENCRYPTION_DECRYPTION_HINT="Contact me for access"
MARKATA_GO_ENCRYPTION_ENFORCE_STRENGTH=false
MARKATA_GO_ENCRYPTION_MIN_ESTIMATED_CRACK_TIME=5d
MARKATA_GO_ENCRYPTION_MIN_PASSWORD_LENGTH=20
```

### .env File Support

Place a `.env` file in your project root. It is loaded automatically during the build:

```bash
# .env
MARKATA_GO_ENCRYPTION_KEY_DEFAULT=my-secret-password
MARKATA_GO_ENCRYPTION_KEY_PERSONAL=another-password
```

Rules:
- Lines starting with `#` are comments
- Values can be quoted with single or double quotes
- Real environment variables take precedence over `.env` values
- The `.env` file should be in your `.gitignore`

## Password Strength Policy

Every encryption key used by a private post must satisfy the configured strength policy before the plugin encrypts the content. The defaults are strict:

- `enforce_strength = true` (can be disabled, but keys will no longer be validated)
- `min_estimated_crack_time = "10y"` (supports `y`, `d`, `h`, `m`, `s` units)
- `min_password_length = 14`

Strength estimation uses the zxcvbn algorithm (dictionary/pattern-aware), so common passwords like `password12345` are treated as weak even when they include numbers.

If any key violates those thresholds, the build halts with an `EncryptionBuildError` that lists the affected posts and key names. Passwords and hints never appear in the error text.

### CLI Password Generator

Generate a compliant password without running the full build:

```
markata-go encryption generate-password
markata-go encryption generate-password --length 20
markata-go encryption check
markata-go encryption check --key default
markata-go encryption encrypt-posts --dry-run
markata-go encryption encrypt-posts
markata-go encryption encrypt-posts --workers 4
```

The command prints only the password to stdout, making it easy to pipe into your `.env` file or a password manager. The optional `--length` flag requests a longer password (it must be at least the configured `min_password_length`). The generated password already meets the default crack-time and length thresholds.

`encryption check` validates configured keys against your active policy and exits non-zero when a key is missing or weak. By default it checks every key referenced by `default_key` and `private_tags`.

### Source-Encrypt Private Posts

Use `encrypt-posts` to rewrite private Markdown files so the body is encrypted at rest in your source repository:

```
markata-go encryption encrypt-posts --dry-run
markata-go encryption encrypt-posts
```

The command scans files matched by your active `glob.patterns`, then encrypts posts marked `private: true` or matched by `encryption.private_tags`. It skips draft, skipped, public, and already source-encrypted posts. `--dry-run` reports what would be encrypted without modifying files.

Missing or weak keys fail before any files are written. The encrypted source keeps frontmatter readable and replaces only the Markdown body with a marker and ciphertext:

```markdown
---
title: My Secret Post
private: true
---

<!-- markata-encrypted-source:v1 key=default -->
BASE64_AES_GCM_CIPHERTEXT
```

### Decrypt Posts Back To Plaintext

Use `decrypt-posts` when you need to edit a source-encrypted post. It is the inverse of `encrypt-posts`:

```
export MARKATA_GO_ENCRYPTION_KEY_DEFAULT='your-password'

# preview every encrypted post that would be decrypted
markata-go encryption decrypt-posts --dry-run

# decrypt one file
markata-go encryption decrypt-posts pages/my-secret-post.md

# decrypt a whole directory tree
markata-go encryption decrypt-posts pages/private/

# decrypt everything matched by glob.patterns
markata-go encryption decrypt-posts

# use four concurrent workers (0, the default, uses available Go CPU parallelism)
markata-go encryption decrypt-posts --workers 4
```

The key name comes from the `key=` value in the encrypted marker, falling back to `encryption.default_key`. The
password is read from `MARKATA_GO_ENCRYPTION_KEY_<KEY>` (uppercased). Files that are not source-encrypted are
skipped, frontmatter is preserved exactly, and a missing env var or wrong password fails before anything is
written.

A typical edit cycle is `decrypt-posts <file>`, edit the Markdown, then `encrypt-posts` before committing.

The commands keep authenticated ciphertext in
`.markata/source-encryption-cache.json`. When a body is unchanged,
`encrypt-posts` verifies and restores the existing ciphertext instead of
generating a new random nonce. This keeps decrypt/re-encrypt cycles from
changing every encrypted file. The cache stores ciphertext only, never
plaintext or password-derived hashes, and is written with owner-only
permissions.

To opt one post out of a private tag or template default, set an explicit
boolean override:

```yaml
private: false
```

An explicit `private: false` takes precedence over `encryption.private_tags`.

Both bulk commands prepare documents in parallel by default, then begin writing only after every candidate is successfully prepared. A missing key, invalid password, malformed source file, cryptographic error, or detected edit during preparation leaves the repository unchanged. Replacements use atomic file writes, and a write failure restores earlier files. Use `--workers 1` if you need serial processing; `--workers 0` is the default and uses `GOMAXPROCS` workers.

When run in a terminal, encryption progress uses your active CLI color theme. Use `--color` to force colors or `--no-color` for plain output.

## Lint Rule

`markata-go lint` now includes an encryption policy check. When encryption is enabled, lint reports an error if configured keys are missing or fail strength thresholds.

## Making Posts Private

There are three ways to make a post private:

### 1. Frontmatter `private: true`

The simplest approach. The post is encrypted with the default key:

```yaml
---
title: My Secret Post
private: true
---
```

### 2. Frontmatter with a specific key

Use `secret_key` (or its aliases `private_key`, `encryption_key`) to encrypt with a named key:

```yaml
---
title: Premium Content
private: true
secret_key: premium
---
```

This looks for `MARKATA_GO_ENCRYPTION_KEY_PREMIUM` in the environment.

All three frontmatter fields are equivalent -- use whichever name you prefer:

| Field | Example |
|-------|---------|
| `secret_key` | `secret_key: premium` |
| `private_key` | `private_key: premium` |
| `encryption_key` | `encryption_key: premium` |

If multiple are set, `secret_key` takes priority, then `private_key`, then `encryption_key`.

### 3. Private tags

Configure tags that automatically mark posts as private:

```toml
[encryption.private_tags]
diary = "personal"
journal = "personal"
```

Any post tagged `diary` or `journal` is automatically treated as private and encrypted with the `personal` key. You don't need to set `private: true` in the frontmatter.

The `private_tags` check matches against both the post's `tags` list and its `templateKey` (or `template`) frontmatter field. This is useful for content that uses `templateKey` as its primary categorization, such as gratitude journals or diary entries that may not have explicit tags.

```toml
[encryption.private_tags]
gratitude = "default"    # Matches posts with tag "gratitude" OR templateKey "gratitude"
```

**Priority rules:**
- If a tag matches, the tag's key is used
- If only `templateKey` matches, its key is used
- **Frontmatter key overrides both:** If a post has a frontmatter `secret_key`, it takes priority over any tag or templateKey match

## Build Behavior

### Missing Keys Fail the Build

If any private post has no available encryption key, the build **fails early with a critical error** that summarizes the failures by key and reason:

```
encryption error: private posts found without available encryption keys. Build halted to prevent exposing private content
key "personal": set MARKATA_GO_ENCRYPTION_KEY_PERSONAL in environment or .env (posts: 12, sample: diary/2024-01-15.md, diary/2024-01-16.md, diary/2024-01-17.md)
```

This is intentional. Private content must never be published unencrypted.

### Draft and Skipped Posts

Posts with `draft: true` or `skip: true` are excluded from encryption checks. They are not published at all, so they don't need encryption.

### Disabling Encryption

To disable encryption entirely:

```toml
[encryption]
enabled = false
```

When disabled, the encryption plugin does nothing. Private posts pass through unmodified (they are still rendered but not encrypted). Use this only for local development.

## Multiple Access Levels

You can use different keys for different audiences:

```toml
[encryption]
default_key = "default"
decryption_hint = "Contact me for the password"

[encryption.private_tags]
diary = "personal"
premium = "subscribers"
```

```bash
# .env
MARKATA_GO_ENCRYPTION_KEY_DEFAULT=general-password
MARKATA_GO_ENCRYPTION_KEY_PERSONAL=my-eyes-only
MARKATA_GO_ENCRYPTION_KEY_SUBSCRIBERS=subscriber-password
```

Then in your posts:

```yaml
# Uses default key
---
private: true
---

# Uses personal key (via tag)
---
tags: [diary]
---

# Uses subscribers key (explicit)
---
private: true
secret_key: subscribers
---
```

## Client-Side Decryption

The encrypted content includes:
- A lock icon and "Encrypted Content" heading
- The decryption hint (if configured)
- A password input field
- A "Remember for this session" checkbox (uses sessionStorage)

When the correct password is entered, JavaScript decrypts the content in-browser using the Web Crypto API with matching PBKDF2 parameters. The decrypted HTML replaces the password form.

Encrypted controls do not autofocus when a page loads, so multiple encrypted
blocks do not move the reader to the last password field. A remembered session
password can decrypt content automatically without taking focus or changing
the current scroll position. When a user submits a password, the revealed
content may receive accessibility focus without moving the viewport.

If decrypted HTML contains another encrypted block, markata-go initializes its
controls immediately. Canonical wrappers that use the same key keep the
existing same-key unlock behavior. Feed projections intentionally remove the
key name, so feed entries remain isolated and require manual decryption.

If "Remember for this session" is checked, the password is stored in sessionStorage (cleared when the browser tab closes). This allows navigating between encrypted posts without re-entering the password for posts using the same key.

Opted-in private feed entries intentionally omit the key name. They can be
decrypted after manual password entry, but they do not use session storage or
unlock other feed entries.

## Privacy Boundary

Encryption protects the **post body**, not metadata. Frontmatter fields like title, description, tags, and dates remain in cleartext by design.

### What stays public

| Field | Why |
|-------|-----|
| Title | Shown in page cards, feed listings, navigation, and HTML `<title>` |
| Description | Only if explicitly set in frontmatter -- you chose to make it public |
| Tags and dates | Used for site structure, filtering, and feed membership |
| Slug / URL | The page needs to be routable and linkable |
| Avatar | Shown in mentions and author cards |

### What is private

| Field | Protection |
|-------|-----------|
| Post body (Markdown) | Cleared from output; never written to any file |
| Article HTML | Encrypted with AES-256-GCM; only the ciphertext is published |
| Auto-generated descriptions | Suppressed entirely for private posts |
| Inlinks / outlinks text | Cleared during metadata scrubbing |

If you put sensitive information in your title or description frontmatter, it **will** be visible in the built site. Keep sensitive content in the post body.

## Privacy Protection

When a post is marked private (by any method), markata-go suppresses it across all output types -- not just the HTML article page. This prevents content from leaking through alternate channels.

### What is protected

| Output | Behavior |
|--------|----------|
| HTML page | Content encrypted with password prompt |
| `.md` / `.txt` alternates | Not generated for private posts |
| OG image cards | Not generated for private posts |
| RSS / Atom / JSON feeds | Private posts are excluded unless the feed explicitly sets `include_private = true`; included entries contain encrypted HTML and safe metadata only |
| Feed pages | Private posts are excluded from public feed pages; only feeds that explicitly set `include_private = true` can include safe metadata and encrypted content |
| Embed cards (`![[slug]]`) | Shows a "Private Content" card with no title, description, or date |
| Wikilinks (`[[slug]]`) | Link text is rendered but `data-title`, `data-description`, `data-date` attributes are omitted |
| Wikilink hover previews | No hover preview is shown for private posts |
| Auto-generated descriptions | Not generated from private post content |

### How it works

Privacy marking happens at the very start of the Transform stage -- before any other plugin processes the post. This ensures that every downstream plugin (description generation, embed cards, wikilinks, feeds, etc.) sees `private: true` and acts accordingly.

The encrypted article HTML is the **only** representation of your private content in the built site.

### Warm builds and key rotation

Encrypted wrappers are cached using the rendered article, key name, resolved
password, public hint, source path (used by accessibility controls), and a
maintained wrapper/browser-crypto revision. Wrapper or browser crypto format
changes bump that revision; old entries miss once without a build-cache schema
migration. A fresh randomized source envelope does not force new browser
ciphertext when plaintext and the key contract are unchanged.

When a wrapper really regenerates—including password/key-name rotation, hint
changes, or encrypted-cache file loss—the page and its dependent pages are
rendered and published again, even if the decrypted input hash is unchanged.
This includes private homepages with an explicitly empty slug. Regeneration
clears the previous full-page cache reference before fresh cache writes. If a
best-effort full-page write fails, saving cache metadata preserves that cleared
reference: the next build re-renders instead of restoring old-password HTML,
even when its encrypted wrapper is a valid hit or its output file is missing.
The old page-cache file may remain orphaned but is no longer usable through
that post's cache reference.

Valid warm hits keep the identical encrypted wrapper and full page and still
scrub private metadata. Wrapper hits and encryption errors do not clear a valid
full-page reference. Best-effort wrapper-cache write failure still invalidates
the old full page; this is not a general transactional cache-recovery guarantee.

When checking a key rotation, first confirm at least two equivalent warm builds
with the old key. Cold-to-warm asset or configuration resets can otherwise mask
stale-page behavior. Verify decryption of the actual published HTML file using
the new key, rather than checking only an in-memory wrapper.

Source-encrypted posts still never cache decrypted parsed posts or plaintext
articles. Inferred titles use a transient hash-only prior-build baseline for
canonical comparison; that baseline stores no plaintext or password. Load's
conservative feed/tag/garden dirty signals are preserved.

### Feed pages and private tags

`private_tags` marks matching posts as private and encrypted. Auto-generated feeds
for private tags opt into those posts so their encrypted entries can render;
other auto-generated tag, category, and archive feeds exclude private posts.

If you intentionally need a private-aware archive or admin page, create an
explicit feed with `include_private = true` and restrict where you publish it.
Private RSS, Atom, and JSON Feed entries still contain no raw body, key name,
derived description, or private media reference.

## Security Notes

- **AES-256-GCM** encryption with random IVs
- **PBKDF2** key derivation with 100,000 iterations
- Encryption happens at build time; decryption happens client-side
- The server only ever serves encrypted content
- Passwords are never transmitted to the server
- Session storage is opt-in and per-tab only
