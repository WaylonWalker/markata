# Editor publishing and review

Status: incrementally implementing on the Go workbench, stacked on #1905.

## User-facing contract

The writing interface MUST show a persistent Publish action in the top bar
regardless of split/editor-only/preview-only layout. The command palette MUST
offer "Publishing: Review & Publish" with publish, go live, deploy, release,
and review aliases. Ctrl/Cmd+Shift+Enter and Ctrl/Cmd+S MUST invoke the same
review flow from all editor keymaps. In Vim mode, `:w` MUST invoke that flow;
Emacs mode's Ctrl-X Ctrl-S MUST do the same. Ctrl/Cmd+P and Ctrl/Cmd+Shift+P
retain file-finder and palette behavior. Opening the dialog MUST NOT commit,
push, build, or deploy.

A scratch file MUST NOT be treated as a site post. Its action MUST indicate that
it needs to be saved to the site first. The editor MUST distinguish autosaved
workspace content from content that has been deployed successfully.

The desktop review MUST support side-by-side rendered versions, while mobile
MUST show one version at a time with tap targets and optional horizontal swiping.
Keyboard and screen-reader access MUST not depend on swiping. The dialog MUST
restore focus when dismissed and MUST not lose editor source.
On desktop, the review dialog MUST use the available viewport height for its
comparison panes while keeping its header and actions visible. Optional full-diff
loading failures MUST replace the loading state with a readable error.
Rendered diff MUST be the initial view. Comparison modes MUST be presented as
keyboard-operable tabs. Side-by-side rendered panes MUST synchronize around
matching rendered blocks when possible and fall back to relative scroll position
when no block matches. Publication status summaries MUST look like informational
chips, not action buttons.

The review header MUST summarize publish metadata differences as textual chips
(not color alone). Especially show changes to draft/published flags, slug changes,
and contradictory publication flags. YAML boolean values MUST be recognized
case-insensitively. An already public, non-draft page MUST NOT offer Mark as
public. Content changes need a full source comparison.
The review MUST provide a Commit message tab showing the exact message used for
the local Git commit. The user MUST be able to edit it, and the server MUST use
the supplied message for that commit.

## Revision identity

At least three independently identified states exist:
- Unsaved editor text (in-memory CodeMirror)
- Saved workspace source (possibly not yet committed)
- The source revision included in the *currently deployed* site.

Git HEAD is only "last committed". It MUST NOT be labeled "live" unless a
deployed build record proves that the running site corresponds to that revision.

Before comparing, the client MUST wait for the exact editor snapshot to reach
workspace storage. It MUST refuse or retry if a newer edit or file switch occurs.
Mutating operations MUST verify source digest and expected Git ref again to
avoid a time-of-check/time-of-use publish.

## Version 1 foundation (read-only)

GET /workspace/publish/review?path=posts/example.md returns the last committed
Git blob and the current saved workspace file. Requests must use the existing
loopback-only workbench gate and workspace path restrictions.

It MUST NOT write Git or site files or claim deployment support. No Git checkout
means a useful "review unavailable" state, never a fabricated baseline. A new
untracked post MUST be identified as having no previous committed revision.
Markup renders in isolated iframe previews; source diff is shown via text nodes
and cannot execute user HTML.

## Next steps: publishing back end

Prepare should return an opaque review token bound to the selected path, saved
source digest, baseline commit, deployment revision, and publication settings.

The confirmation request MUST be scoped to explicitly selected paths, not
"git add .". It MUST preserve unrelated worktree changes, staged files,
branches, and index state. Reject conflicting HEAD/source versions. Publish
metadata should be validated before commit and include explicit, intelligible
redirect consequences for slug/path changes.

A Git commit is not deployment success. Push, build, deployment, verification,
and optional PR review are separately observable states. A failed deployment
MUST keep the revision recoverable and avoid showing "Published" prematurely.
The UI should support publishing a single post, a selected group, or an explicit
site-wide release, each displaying its scope.

Power users MAY opt into Git details and a full Pierre source diff. The default
experience must be a rendered content review. A rich unified rendered diff must
respect paragraph, heading, list, table, code, image, container and frontmatter
semantics, with fallbacks for unrepresentable changes. Important link-target,
script, asset and metadata modifications MUST remain discoverable.

## Tests

Cover command-button-shortcut parity, all editor keymaps, interrupted autosave,
workspace change while reviewing, scratch files, new/untracked posts, missing Git,
large diffs, injected HTML, hidden metadata changes, mobile tap and swipe,
keyboard focus, exact source bytes, preview rendering failures, unrelated staged
files, Git conflicts, failed builds, stale deployments, and explicit confirmation.


## Publishing implementation progress (experimental stack)

- The review surface can compare exact last-committed and saved versions, show
  semantic metadata warnings, render side-by-side, and show a source diff.
- Committing and pushing are **separate** explicit operations. The commit
  endpoint checks expected HEAD and source SHA-256 and stages the approved
  source in an isolated Git index. Publication flags must already indicate
  public content before a Git commit. The UI offers an explicit Mark as
  public action, then reopens review.
- A successful Git push is still **not** a verified public deployment. Live
  status needs a deployment-specific integration and content revision receipt.
- The rendered-diff comparison aligns structural blocks and annotates plain
  text word differences; tables, code, diagrams, and complex containers use
  whole-block comparisons. The side-by-side rendering remains authoritative.
- Full source review uses Pierre's vanilla FileDiff, compiled to a local,
  lazy-loaded optional asset. Never fetch third-party code at publish time.
  Build the asset with tools/publish-diffs. The lightweight source diff must
  continue working without that asset.

The code is intentionally marked experimental until the exact Git workflow,
bundled asset, Playwright and deployment recovery paths are verified.
