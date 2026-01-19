"""
The `markata.plugins.md_it_wikilinks` plugin adds support for wiki-style links using
double brackets (`[[link]]`). It automatically resolves links to other posts in your
site using file names or slugs.

## Installation

This plugin is built-in and enabled by default through the 'default' plugin.
If you want to be explicit, you can add it to your list of plugins:

```toml
hooks = [
    "markata.plugins.md_it_wikilinks",
]
```

## Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

```toml
disabled_hooks = [
    "markata.plugins.md_it_wikilinks",
]
```

## Configuration

This plugin supports comprehensive configuration for wikilink resolution.

Configuration Options:
```toml
[plugins.md_it_wikilinks]

# Resolution strategy for duplicate matches
# "priority" (default): Use priority scoring to select best match
# "first": Use first match found
# "warn": Always warn about duplicates
resolution_strategy = "priority"

# Score difference threshold for clear winner selection (default: 20)
clear_winner_threshold = 20

# Suppress warnings for links matching these patterns
suppress_patterns = ["tag/*", "category/*", "archive/*"]

# Custom priority rules (higher priority = higher score)
priority_rules = [
    { pattern = "pages/*", priority = 100 },
    { pattern = "posts/*", priority = 90 },
    { pattern = "blog/*", priority = 85 },
    { pattern = "tutorials/*", priority = 80 },
    { pattern = "docs/*", priority = 75 },
    { pattern = "tag/*", priority = 60 },
    { pattern = "category/*", priority = 55 },
    { pattern = "archive/*", priority = 50 },
    { pattern = "feed/*", priority = 45 },
]

# Behavior for broken links (default: "warn", "silent", "error")
fallback_behavior = "warn"

# Enable/disable logging (default: true)
enable_logging = true
```

## Resolution Priority System

Built-in Scores (configurable via priority_rules):
```toml
# Default scoring system (overridable via priority_rules)
exact_slug_match = 100    # Exact slug match: highest priority
path_match = 80              # Path match: [[folder/page]] -> high priority
feed_slug_match = 60           # Feed slug match: tag/*, category/* patterns
basename_match = 40             # Basename match: default fallback
```

Examples:
```toml
# Override scores for entire patterns
priority_rules = [
    { pattern = "tag/python", priority = 70 },        # Higher than default 60
    { pattern = "docs/*", priority = 95 },        # Higher than exact match!
]

# Change entire scoring algorithm
resolution_strategy = "first"      # Disable scoring, use first match
resolution_strategy = "warn"     # Always warn about duplicates
clear_winner_threshold = 15      # Lower threshold = more warnings
```

## Smart Slug Resolution

The plugin:
1. Looks up the target file in your content
2. Finds its generated slug
3. Creates a link to the final URL

## Link Formats

Supports multiple link styles:
- Basic: `[[filename]]`
- With text: `[[filename|Link Text]]`
- With path: `[[folder/file]]`
- With extension: `[[file.md]]` (extension stripped in output)
- With anchors: `[[filename#anchor]]` or `[[filename#anchor|Display Text]]`
- Complex: `[[folder/file#anchor|Display Text]]`

**Edge Cases Handled:**
- **Spaces**: `[[my tag]]` → matches slug "my-tag" (spaces normalized to hyphens by Markata)
- **Empty display text**: `[[page|]]` → falls back to "page"
- **Multiple pipes**: `[[page|Text with | pipes]]` → splits only on first pipe
- **Quotes**: `[[page|"Display Text"]]` → supports quoted display text
- **Mixed anchors**: `[[page#anchor|Text]]` and `[[page|Text#anchor]]`
- **Complex expressions**: `[[tag/my tag|Posts about "my tag"]]` (full syntax support with quotes and spaces)

## Duplicate Resolution System

The wikilinks system uses intelligent priority-based resolution:

How it works:
1. Creates mapping of all possible matches in markata.possible_wikilink
2. Scores each candidate using configurable priority rules
3. Selects highest-scoring match (clear winner if score difference > threshold)
4. Only warns for truly ambiguous cases

## HTML Output

Generated HTML structure:
```html
<a class="wikilink" href="/target-slug">Link Text</a>
```

## Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

```toml
disabled_hooks = [
    "markata.plugins.md_it_wikilinks",
]
```

## Configuration

This plugin supports comprehensive configuration for wikilink resolution:

```toml
[plugins.md_it_wikilinks]
# Resolution strategy for duplicate matches
# "priority" (default): Use priority scoring to select best match
# "first": Use first match found
# "warn": Always warn about duplicates
resolution_strategy = "priority"

# Score difference threshold for clear winner selection (default: 20)
clear_winner_threshold = 20

# Suppress warnings for links matching these patterns
suppress_patterns = ["tag/*", "category/*", "archive/*"]

# Custom priority rules (higher priority = higher score)
# If no rules match, defaults to built-in scoring system
priority_rules = [
    { pattern = "pages/*", priority = 100 },
    { pattern = "posts/*", priority = 90 },
    { pattern = "blog/*", priority = 85 },
    { pattern = "tutorials/*", priority = 80 },
    { pattern = "docs/*", priority = 75 },
    { pattern = "tag/*", priority = 60 },
    { pattern = "category/*", priority = 55 },
    { pattern = "archive/*", priority = 50 },
    { pattern = "feed/*", priority = 45 },
]

## Priority Scoring System

**Built-in Scores:**
- **Exact slug match**: 100 points (highest priority)
- **Path match**: 80 points (when original link includes path structure)
- **Feed slug match**: 60 points (for hierarchical feeds like `tag/*`)
- **Basename match**: 40 points (default fallback)

**Resolution Examples:**
```toml
# For site with feed "tag/python" and page "python":
[[python]]          # Matches: ["python", "tag/python"] → selects "python" (100 vs 60)
[[tag/python]]      # Matches: ["python", "tag/python"] → selects "tag/python" (100 vs 40)

# For site with multiple tag feeds:
priority_rules = [
    { pattern = "tag/python", priority = 90 },  # Higher than default tag/*
    { pattern = "tag/javascript", priority = 85 },
]
```

# Behavior for broken links (default: "warn")
# "warn": Log warning and use fallback link
# "silent": Use fallback link silently
# "error": Log error and use fallback link
fallback_behavior = "warn"

# Enable/disable logging (default: true)
enable_logging = true
```

### Priority Rules

Priority rules allow you to customize how links are resolved when there are multiple matches:

- **Pattern**: Glob pattern matching page slugs (supports `*` wildcard)
- **Priority**: Score value (higher = more likely to be selected)

The plugin evaluates rules in order and uses the first matching rule. If no custom rules match, it falls back to the built-in scoring system:

- Exact slug match: 100 points
- Path match (`[[folder/page]]`): 80 points
- Feed slug match: 60 points
- Basename match: 40 points

### Warning Suppression

Use `suppress_patterns` to reduce warning noise for expected duplicates:

```toml
suppress_patterns = [
    "tag/*",          # Suppress all tag-related warnings
    "category/*",     # Suppress category warnings
    "*/index",        # Suppress index page conflicts
]
```

By default, the plugin uses priority-based resolution that automatically selects
the best match and only warns for truly ambiguous cases.

## Functionality

## Basic Wikilinks

Simple file-based linking:
```markdown
[[nav]]              -> links to docs/nav.md as /nav
[[blog/post]]        -> links to blog/post.md as /blog/post
[[about|About Me]]   -> links to about.md with "About Me" as text
```

## Smart Slug Resolution

The plugin:
1. Looks up the target file in your content
2. Finds its generated slug
3. Creates a link to the final URL

Example:
```markdown
# File: posts/2024-01-my-post.md
slug: /blog/my-post

# In another file:
[[2024-01-my-post]]  -> links to /blog/my-post
```

## Link Formats

Supports multiple link styles:
- Basic: `[[filename]]`
- With text: `[[filename|Link Text]]`
- With path: `[[folder/file]]`
- With extension: `[[file.md]]` (extension stripped in output)
- With anchors: `[[filename#anchor]]` or `[[filename#anchor|Display Text]]`
- Complex: `[[folder/file#anchor|Display Text]]`

**Edge Cases Handled:**
- **Spaces**: `[[my tag]]` → matches slug "my-tag" (spaces normalized to hyphens by Markata)
- **Empty display text**: `[[page|]]` → falls back to "page"
- **Multiple pipes**: `[[page|Text with | pipes]]` → splits only on first pipe (first `|` is separator, rest is content)
- **Quotes**: `[[page|"Display Text"]]` → supports quoted display text (quotes preserved in display text)
- **Mixed anchors**: `[[page#anchor|Text]]` and `[[page|Text#anchor]]` (anchors and display text work together)
- **Complex expressions**: `[[tag/my tag|Posts about "my tag"]]` (full syntax support with quotes and spaces)

**Normalization Behavior:**
- **Link target normalization**: Spaces and special characters handled according to Markata slug conventions
- **Display text preservation**: Exact display text (including quotes and pipes) preserved
- **Case sensitivity**: Resolution is case-sensitive (matches Markata's slug handling)

**Common Use Cases:**
```markdown
# Basic usage
[[thoughts]]                    # → /thoughts
[[tag/python]]                 # → /tag/python
[[docs/getting-started]]        # → /docs/getting-started

# With display text
[[thoughts|My Thoughts]]        # → /thoughts (text: "My Thoughts")
[[tag/python|Python Posts]]     # → /tag/python (text: "Python Posts")

# With anchors
[[thoughts#intro]]             # → /thoughts#intro
[[thoughts#intro|Introduction]] # → /thoughts#intro (text: "Introduction")

# Complex with spaces and quotes
[[my tag|Posts about "my tag"]]  # → /my-tag (text: 'Posts about "my tag"')
```

## Resolution Priority System

**Built-in Scores (configurable via priority_rules):**
- **Exact slug match**: 100 points (highest priority)
- **Path match**: 80 points (when original link includes path structure)
- **Feed slug match**: 60 points (for hierarchical feeds like `tag/*`)
- **Basename match**: 40 points (default fallback)

## Duplicate Resolution System

The wikilinks system uses intelligent priority-based resolution to eliminate warning noise for common hierarchical patterns:

### How Matches Are Created

**Feed Slugs**: Feed configurations (e.g., `slug = "tag/python"`)
**Page Slugs**: Regular page slugs (e.g., `slug = "python"`)
**Mapping Creation**: Both are mapped in `markata.possible_wikilink` dictionary:

```python
# For feeds in markata.feeds:
for slug in [v.config.slug for v in markata.feeds.values()]:
    wikilink = slug.split("/")[-1]  # Extract basename: "python"
    markata.possible_wikilink[wikilink].append(slug)  # Maps "python" → ["python", "tag/python"]

# For regular pages:
for slug in markata.map("slug"):
    wikilink = slug.split("/")[-1]  # Extract basename: "python"
    if wikilink not in markata.possible_wikilink:
        markata.possible_wikilink[wikilink] = [slug]
```

### Resolution Priority System

## Resolution Priority System

**Built-in Scores (configurable via priority_rules):**
```toml
# Default scoring system (overridable via priority_rules)
exact_slug_match = 100    # Exact slug match: highest priority
path_match = 80              # Path match: [[folder/page]] -> high priority
feed_slug_match = 60           # Feed slug match: tag/*, category/* patterns
basename_match = 40             # Basename match: default fallback
```

**Configuration Options:**
```toml
# Method 1: Override default scores for entire patterns
priority_rules = [
    { pattern = "tag/*", priority = 70 },        # Higher than default 60
    { pattern = "category/*", priority = 80 },    # Higher than default 55
    { pattern = "docs/*", priority = 95 },        # Very high priority
    { pattern = "posts/*", priority = 90 },        # High priority for content
]

# Method 2: Change entire scoring algorithm
resolution_strategy = "first"    # Disable scoring, use first match
resolution_strategy = "warn"     # Always warn about duplicates
clear_winner_threshold = 15      # Lower threshold = more warnings
```

**Resolution Examples:**
```toml
# Example 1: Default behavior with feed "tag/python" and page "python"
# markata.possible_wikilink = {"python": ["python", "tag/python"]}

[[python]]          # → selects "python" (100 vs 60 points, clear winner)
[[tag/python]]      # → selects "tag/python" (100 vs 40 points, clear winner)

# Example 2: Custom priority rules
priority_rules = [
    { pattern = "tag/python", priority = 75 },  # Lower than default 60
    { pattern = "docs/getting-started", priority = 120 }, # Higher than exact match!
]

[[python]]          # Against "tag/python": 75 vs 60 → selects "tag/python" (custom rule wins)
[[tag/python]]      # Against "tag/python": 100 vs 75 → selects "tag/python" (exact match still wins)

# Example 3: Different resolution strategies
resolution_strategy = "first"   # Ignores scoring, always first match
resolution_strategy = "warn"     # Always warns about duplicates
```

**Advanced Configuration:**
```toml
# Fine-tune duplicate resolution behavior
clear_winner_threshold = 30    # Require larger score difference for clear winner
suppress_patterns = ["tag/*"] # Suppress warnings for all tag/* matches
fallback_behavior = "silent"      # No warnings for broken links
```

**Example 2: Multiple Feeds with Same Basename**
```toml
# Site structure:
# posts/thoughts.md      (slug: "thoughts")
# feed: tag/thoughts    (slug: "tag/thoughts")
# feed: category/thoughts (slug: "category/thoughts")

# Result in markata.possible_wikilink:
{
    "thoughts": ["thoughts", "tag/thoughts", "category/thoughts"]  # Three matches!
}

[[thoughts]]           # → selects "thoughts" (100 vs 60 vs 55 points)
[[tag/thoughts]]      # → selects "tag/thoughts" (100 vs 40 vs 55 points)
[[category/thoughts]] # → selects "category/thoughts" (100 vs 40 vs 55 points)
```

**Example 3: Custom Priority Rules**
```toml
# If you want specific feed patterns to have higher priority:
priority_rules = [
    { pattern = "tag/python", priority = 70 },  # Higher than default tag/* (60)
    { pattern = "docs/*", priority = 95 },     # Documentation gets highest priority
    { pattern = "posts/*", priority = 90 },      # Posts get high priority
]

# Result: Custom rules override built-in scoring for matching patterns
```

**Configuration Priority Rules Override:**
```toml
# Custom scoring for specific patterns
priority_rules = [
    { pattern = "pages/*", priority = 100 },  # Pages get highest priority
    { pattern = "tag/python", priority = 90 },  # Specific tag gets boost
    { pattern = "docs/*", priority = 95 },     # Documentation gets high priority
]

# Feed patterns with custom priorities
priority_rules = [
    { pattern = "tag/python", priority = 70 },  # Higher than default tag/* (60)
    { pattern = "category/javascript", priority = 75 }, # Boost specific categories
]
```

## HTML Output

Generated HTML structure:
```html
<a class="wikilink" href="/target-slug">Link Text</a>
```

## Error Handling

For broken links:
- Maintains the wikilink syntax
- Adds a 'broken-link' class
- Optionally logs warnings

## Dependencies

This plugin depends on:
- markdown-it-py for markdown parsing
- The `render_markdown` plugin for final HTML output
"""

import logging
import re
from typing import TYPE_CHECKING
from typing import Dict
from typing import List

from markdown_it import MarkdownIt
from markdown_it.rules_inline import StateInline

from markata.hookspec import hook_impl
from markata.hookspec import register_attr

if TYPE_CHECKING:
    from markata import Markata

logger = logging.getLogger("markata")


def get_default_config() -> Dict:
    """
    Get default configuration for wikilinks plugin.

    Returns:
        Dictionary with default configuration values
    """
    return {
        "resolution_strategy": "priority",
        "clear_winner_threshold": 20,
        "suppress_patterns": [],
        "priority_rules": [
            {"pattern": "pages/*", "priority": 100},
            {"pattern": "posts/*", "priority": 90},
            {"pattern": "blog/*", "priority": 85},
            {"pattern": "tutorials/*", "priority": 80},
            {"pattern": "docs/*", "priority": 75},
            {"pattern": "tag/*", "priority": 60},
            {"pattern": "category/*", "priority": 55},
            {"pattern": "archive/*", "priority": 50},
            {"pattern": "feed/*", "priority": 45},
        ],
        "fallback_behavior": "warn",  # "warn", "silent", "error"
        "enable_logging": True,
    }


def get_plugin_config(markata: "Markata") -> Dict:
    """
    Get plugin configuration with defaults merged.

    Args:
        markata: Markata instance

    Returns:
        Merged configuration dictionary
    """
    default_config = get_default_config()
    user_config = (
        getattr(markata, "config", {}).get("plugins", {}).get("md_it_wikilinks", {})
    )

    # Deep merge user config with defaults
    merged_config = default_config.copy()
    merged_config.update(user_config)

    # Handle nested priority_rules merging
    if "priority_rules" in user_config:
        merged_config["priority_rules"] = user_config["priority_rules"]

    return merged_config


def matches_pattern(path: str, pattern: str) -> bool:
    """
    Check if a path matches a glob-like pattern.

    Args:
        path: The path to check
        pattern: The pattern (supports * wildcard)

    Returns:
        True if path matches pattern
    """
    # Convert glob pattern to regex
    regex_pattern = pattern.replace("*", ".*")
    return re.match(f"^{regex_pattern}$", path) is not None


def should_suppress_warning(link_target: str, suppress_patterns: List[str]) -> bool:
    """
    Check if warning should be suppressed for a link target.

    Args:
        link_target: The link target to check
        suppress_patterns: List of patterns to suppress

    Returns:
        True if warning should be suppressed
    """
    for pattern in suppress_patterns:
        if matches_pattern(link_target, pattern):
            return True
    return False


def calculate_match_score(
    link_target: str, candidate_slug: str, original_link: str, config: Dict
) -> int:
    """
    Calculate priority score for a wikilink match.
    Higher scores indicate better matches.

    Args:
        link_target: The target link text
        candidate_slug: The candidate slug to score
        original_link: The original link text from markdown
        config: Plugin configuration

    Returns:
        Score for the candidate (higher is better)
    """
    # Check custom priority rules first
    priority_rules = config.get("priority_rules", [])
    for rule in priority_rules:
        pattern = rule.get("pattern", "")
        priority = rule.get("priority", 50)
        if matches_pattern(candidate_slug, pattern):
            return priority

    # Default scoring system if no custom rules match
    # Exact slug match (highest priority)
    if link_target == candidate_slug:
        return 100

    # Path match - when original link includes path structure
    if "/" in original_link and original_link.strip("/") == candidate_slug:
        return 80

    # Check if this is a basename match
    basename = candidate_slug.split("/")[-1]
    if link_target == basename:
        # Heuristic: feeds typically have paths like "tag/", "category/", "archive/"
        feed_prefixes = ["tag/", "category/", "archive/", "feed/", "topic/"]
        if any(candidate_slug.startswith(prefix) for prefix in feed_prefixes):
            return 60  # Feed slug match
        else:
            return 40  # Regular basename match

    return 0


def resolve_best_match(
    link_target: str,
    possible_pages: list,
    original_link: str,
    markata: "Markata",
    md=None,
) -> str:
    """
    Resolve the best match from possible pages using priority scoring and configuration.

    Args:
        link_target: The target link text
        possible_pages: List of possible page slugs
        original_link: The original link text from markdown
        markata: Markata instance
        md: Markdown-it instance (optional)

    Returns:
        Best matching page slug
    """
    # Get full plugin configuration
    config = get_plugin_config(markata)
    resolution_strategy = config.get("resolution_strategy", "priority")
    threshold = config.get("clear_winner_threshold", 20)
    suppress_patterns = config.get("suppress_patterns", [])
    fallback_behavior = config.get("fallback_behavior", "warn")
    enable_logging = config.get("enable_logging", True)

    if len(possible_pages) == 1:
        return possible_pages[0]

    # Check if warning should be suppressed
    should_suppress = should_suppress_warning(link_target, suppress_patterns)

    # For non-priority strategies, fallback to simple behavior
    if resolution_strategy == "first":
        return possible_pages[0]
    elif resolution_strategy == "warn":
        # Always warn and use first match (unless suppressed)
        if not should_suppress and enable_logging and fallback_behavior == "warn":
            if md is None or md.options.get("article") is None:
                debug_value = "UNKNOWN"
            else:
                debug_value = md.options["article"].get(
                    "path",
                    md.options["article"].get(
                        "title", md.options["article"].get("slug", "")
                    ),
                )
            logger.warning(
                f"wikilink [[{original_link}]] has duplicate matches ({possible_pages}) in file '{debug_value}', defaulting to the first match ({possible_pages[0]})",
            )
        return possible_pages[0]

    # Priority-based resolution (default)
    # Calculate scores for all candidates using custom configuration
    scored_candidates = []
    for candidate in possible_pages:
        score = calculate_match_score(link_target, candidate, original_link, config)
        scored_candidates.append((score, candidate))

    # Sort by score (descending) and return the highest scoring match
    scored_candidates.sort(key=lambda x: x[0], reverse=True)

    # Check if we have a clear winner (score difference > threshold)
    if len(scored_candidates) >= 2:
        top_score, top_candidate = scored_candidates[0]
        second_score, second_candidate = scored_candidates[1]

        # If clear winner, return it without warning
        if top_score - second_score > threshold:
            return top_candidate

    # If no clear winner, return top choice but log warning for ambiguity (unless suppressed)
    top_score, top_candidate = scored_candidates[0]

    if (
        not should_suppress
        and enable_logging
        and fallback_behavior in ["warn", "error"]
    ):
        if md is None or md.options.get("article") is None:
            debug_value = "UNKNOWN"
        else:
            debug_value = md.options["article"].get(
                "path",
                md.options["article"].get(
                    "title", md.options["article"].get("slug", "")
                ),
            )

        message = (
            f"wikilink [[{original_link}]] has ambiguous matches ({possible_pages}) "
            f"in file '{debug_value}', selecting highest priority match ({top_candidate})"
        )

        if fallback_behavior == "error":
            logger.error(message)
        else:
            logger.warning(message)

    return top_candidate


@hook_impl()
@register_attr("possible_wikilink")
def pre_render(markata: "Markata") -> None:
    markata.possible_wikilink = {}

    for slug in markata.map("slug"):
        # register both final slug and full path slug
        wikilink = slug
        if wikilink in markata.possible_wikilink:
            if slug not in markata.possible_wikilink[wikilink]:
                markata.possible_wikilink[wikilink].append(slug)
        else:
            markata.possible_wikilink[wikilink] = [slug]

        wikilink = slug.split("/")[-1]
        if wikilink in markata.possible_wikilink:
            if slug not in markata.possible_wikilink[wikilink]:
                markata.possible_wikilink[wikilink].append(slug)
        else:
            markata.possible_wikilink[wikilink] = [slug]
    markata.possible_wikilink["index"] = ["index"]

    for slug in [v.config.slug for v in markata.feeds.values()]:
        # Register the full slug (e.g., "tag/python")
        wikilink = slug
        if wikilink in markata.possible_wikilink:
            if slug not in markata.possible_wikilink[wikilink]:
                markata.possible_wikilink[wikilink].append(slug)
        else:
            markata.possible_wikilink[wikilink] = [slug]

        # Register the basename (e.g., "python")
        wikilink = slug.split("/")[-1]
        if wikilink in markata.possible_wikilink:
            if slug not in markata.possible_wikilink[wikilink]:
                markata.possible_wikilink[wikilink].append(slug)
        else:
            markata.possible_wikilink[wikilink] = [slug]


def wikilinks_plugin(
    md: MarkdownIt,
    start_delimiter: str = "[",
    end_delimiter: str = "]",
    markata=None,
):
    """A plugin to create wikilinks tokens.
    These, token should be handled by the renderer.

    ???+ example

        ```md title=markdown
        [[nav]]
        ```

        ```html title=html
        <a class="wikilink" href="/nav">load</a>
        ```
    """

    start_char = ord(start_delimiter)
    end_char = ord(end_delimiter)

    def _wikilinks_inline(state: StateInline, silent: bool):
        try:
            if (
                ord(state.src[state.pos]) != start_char
                or ord(state.src[state.pos + 1]) != start_char
            ):
                return False
        except IndexError:
            return False

        pos = state.pos + 2
        found_closing = False
        while True:
            try:
                end = state.src.find(chr(end_char), pos)
            except ValueError:
                return False
            try:
                if state.src[end + 1] == chr(end_char):
                    found_closing = True
                    break
            except IndexError:
                return False
            pos = end + 2

        if not found_closing:
            return False

        text = state.src[state.pos + 2 : end].strip()
        state.pos = end + 2

        if silent:
            return True

        token = state.push("link_open", "a", 1)
        token.block = False
        token.attrSet("class", "wikilink")

        # Parse display text override syntax: [[page|Display Text]]
        if "|" in text:
            # Split only on first pipe to allow pipes in display text
            link_part, display_text = text.split("|", 1)
            display_text = display_text.strip()
            # Fall back to link_part if display_text is empty
            if not display_text:
                display_text = link_part
        else:
            link_part, display_text = text, None

        # Handle anchor in link part: [[page#anchor]] or [[page#anchor|Display Text]]
        if "#" in link_part:
            link, id = link_part.split("#", 1)
            link = link.strip("/")
        else:
            link, id = link_part.strip("/"), None

        # Get configuration for handling broken links
        if markata:
            config = get_plugin_config(markata)
            suppress_patterns = config.get("suppress_patterns", [])
            fallback_behavior = config.get("fallback_behavior", "warn")
            enable_logging = config.get("enable_logging", True)
        else:
            config = get_default_config()
            suppress_patterns = config.get("suppress_patterns", [])
            fallback_behavior = config.get("fallback_behavior", "warn")
            enable_logging = config.get("enable_logging", True)

        # possible_pages = markata.filter(
        #     f'str(path).split("/")[-1].split(".")[0].replace("_", "-") == "{link.replace("_", "-")}"',
        # )
        possible_pages = markata.possible_wikilink.get(link, []) if markata else []
        if len(possible_pages) == 1:
            link = possible_pages[0]
        elif len(possible_pages) > 1:
            # Use priority-based resolution instead of simple first match
            link = resolve_best_match(link, possible_pages, text, markata, md)
        else:
            # No matches found - handle according to configuration
            should_suppress = should_suppress_warning(link_part, suppress_patterns)

            if (
                not should_suppress
                and enable_logging
                and fallback_behavior in ["warn", "error"]
            ):
                if md.options.get("article") is None:
                    debug_value = "UNKNOWN"
                else:
                    debug_value = md.options["article"].get(
                        "path",
                        md.options["article"].get(
                            "title", md.options["article"].get("slug", "")
                        ),
                    )

                message = f"wikilink [[{text}]] no matches in file '{debug_value}', defaulting to '/{link_part}'"

                if fallback_behavior == "error":
                    logger.error(message)
                else:
                    logger.warning(message)

            # Fallback to original link text
            link = link_part

        if id and not link.endswith(f"#{id}"):
            link = f"{link}#{id}"

        token.attrSet("href", f"/{link}")
        content_token = state.push("text", "", 0)
        # Use display text if available, otherwise fall back to the link part
        content_token.content = display_text if display_text is not None else link_part

        token = state.push("link_close", "a", -1)
        token.content = display_text if display_text is not None else link_part

        return True

    md.inline.ruler.before("escape", "wikilinks_inline", _wikilinks_inline)
