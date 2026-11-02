# Editor file picker

The shared Go and Python workspace editor provides Open file in the command
palette and Ctrl/Cmd+P. Ctrl/Cmd+Shift+P continues to open commands.
A modal combobox searches discovered workspace paths and available titles without
fetching source or waiting for metadata. It includes the standalone draft.
The workspace sidebar has a visible Open file button for mouse and touch users.

Use bundled MIT-licensed fuzzysort, with immutable prepared search snapshots.
Rank basename matches above path and title matches. Show at most 60 rows,
with safe DOM match highlighting, full paths, a result count and empty state.
Arrow keys select; Enter opens through the existing workspace selection handler;
Escape cancels. Focus returns to the editor. Composition input must not trigger
shortcuts. Read-only files still permit file navigation and command discovery.
Each result separates its filename, full path, and optional title; the selected
result is visibly highlighted and kept in view while navigating by keyboard.
Opening a file preserves the existing draft, source access and stale-request rules.
Metadata updates rebuild the index and refresh an open picker.
Files with a root-relative image path in frontmatter show a thumbnail after
metadata loads. Image loading MUST not delay path search or file selection.

Research: fuzzysort offers Sublime-style scoring and prepared targets with no
runtime dependencies. uFuzzy favors configurable word matching; fzf-for-js ports
fzf but documents browser performance tradeoffs. Choose fuzzysort for this small
plain-JavaScript runtime. Validate ranking against ambiguous filenames, path
queries and titles, and measure repeated queries over at least 10,000 paths.
