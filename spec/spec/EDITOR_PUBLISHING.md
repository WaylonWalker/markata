# Editor publishing and review

Status: incrementally implementing on the Go workbench, stacked on #1905.

## User-facing contract

The writing interface MUST show a persistent Publish action in the top bar
regardless of split/editor-only/preview-only layout. The command palette MUST
offer "Publishing: Review & Publish" with publish, go live, deploy, release,
and review aliases. Ctrl/Cmd+Shift+Enter MUST invoke the same action from all
editor keymaps; Ctrl/Cmd+P and Ctrl/Cmd+Shift+P retain file-finder and palette
behavior. Opening the dialog MUST NOT commit, push, build, or deploy.

A scratch file MUST NOT be treated as a site post. Its action MUST indicate that
it needs to be saved to the site first. The editor MUST distinguish autosaved
workspace content from content that has been deployed successfully.

The desktop review MUST support side-by-side rendered versions, while mobile
MUST show one version at a time with tap targets and optional horizontal swiping.
Keyboard and screen-reader access MUST not depend on swiping. The dialog MUST
restore focus when dismissed and MUST not lose editor source.

The review header MUST summarize publish metadata differences as textual chips
(not color alone). Especially show changes to draft/published flags, slug changes,
and contradictory publication flags. Content changes need a full source comparison.

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
