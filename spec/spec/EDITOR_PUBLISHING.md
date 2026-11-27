# Editor publishing and review

Status: incrementally implementing on the Go workbench, stacked on #1905.

## User-facing contract

The writing interface MUST show a persistent **Review changes** action in the top bar
regardless of split/editor-only/preview-only layout. The command palette MUST
offer "Publishing: Review changes" with publish, go live, deploy, release,
and review aliases. Ctrl/Cmd+Shift+Enter and Ctrl/Cmd+S MUST invoke the same
review flow from all editor keymaps. In Vim mode, `:w` MUST invoke that flow;
Emacs mode's Ctrl-X Ctrl-S MUST do the same. Ctrl/Cmd+P and Ctrl/Cmd+Shift+P
retain file-finder and palette behavior. Opening the dialog MUST NOT commit,
push, build, or deploy.

A scratch file MUST NOT be treated as a site post. Its action MUST indicate that
it needs to be saved to the site first. The editor MUST distinguish autosaved
workspace content from content that has been deployed successfully.

The quick-publish action MUST show a disabled **Up to date** state only when a
verified saved-source review confirms that the active, public, non-draft post
matches its local Git commit. Editing the source MUST restore **Publish post**.
Unknown Git state MUST NOT be presented as up to date. The status MUST NOT
imply remote push or live deployment. The existing publish menu MUST retain
access to detailed comparison. Reuse the change gutter review request rather
than adding a second polling loop.

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

## Workspace-wide Changes review queue

The workbench MUST enumerate all Git worktree changes independently of which
files have ever been opened in the editor. A read-only, loopback-only
`GET /workspace/changes` endpoint returns staged, unstaged, and untracked
files; it preserves rename/copy source paths and reports changes outside the
current site rather than silently filtering them away. On failure the UI
shows an error, never an empty success state. The desktop Changes panel
remains available while writing; mobile uses a dismissible bottom sheet.

The Changes view MUST let writers filter a review queue by path and distinguish
posts from other files. Selecting a supported text file shows the current
saved source against the last committed Git baseline. The interface MUST show
Git status and staged/unstaged state and MUST clearly explain when it cannot
show a diff, such as for deleted, conflicted, binary, oversized, or out-of-site
files. Git HEAD is the last commit, not a verified deployment revision.

**Keep in workspace** is a session-only review decision. It leaves the file and
Git index unchanged and does not stage, commit, push, publish, or remove the
entry from the Git inventory. **Review and reset** opens the existing
single-document Review flow only for eligible workspace Markdown files. Reset
must retain its explicit confirmation, exact saved-source and Git-baseline
preconditions, and stale-draft refusal. It remains unavailable for staged
changes, deleted files, conflicts, and out-of-site or non-Markdown files, with
a reason shown to the writer. Reset MUST NOT change unrelated paths or Git
history.

The Changes queue MUST keep its quick source comparison readable without the
optional full-diff renderer. Its Markdown review action MUST expose the
existing rendered, side-by-side, source, and full-diff modes.

The review queue is an inspection and triage surface; multi-file commit and
publish remain separate workflows. Closing review leaves all workspace changes
untouched. Never equate an empty worktree with a fully pushed or live site;
local commits may be ahead of the remote. See `GO_WORKBENCH.md` for the
keyboard, responsive layout, refresh, and session-progress requirements.

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

## Review action and discard contract

**Review changes** is the only primary editor entry point. Review shows a frozen,
exact saved draft against its clearly labeled committed baseline. The review
must offer **Publish**, **Discard changes…**, and **Keep editing**. Keep editing
must never mutate any draft. Publish may guide the writer through explicit
commit, push, and verified deployment stages; it must not declare the page live
until deployment verification succeeds.

Discard is a separate destructive action, not a synonym for closing Review or
undoing a single edit. It MUST require confirmation naming the document and the
baseline being restored. An existing tracked document reverts its workspace
draft to the exact committed blob; a never-committed document removes only that
new draft from the active workspace. Discard MUST NOT mutate unrelated files,
the Git index, HEAD, staged work, external drafts, or builder deployments.

The backend MUST validate the authorized site and document, exact reviewed
source identity, expected committed baseline, and any newer editor/outbox edits
before applying a discard. Stale or concurrent changes MUST fail closed and
require re-review; an offline client MUST retain pending edits instead of
silently discarding them. The hosted sync backend MUST record a new durable
revision or tombstone, with idempotent retry semantics and change-feed updates,
rather than erasing history or conflict copies. A successful local-only discard
MUST NOT claim remote synchronization. The loopback-only workbench endpoint
MUST NOT be exposed on hosted ingress.

The local workbench now has an experimental, revision-checked restore for
Git-tracked files, with explicit confirmation and a server-side lock against
other writes made through the same handler. The editor must refresh recovery
state and verify the saved baseline before reporting local success. This is
not a hosted atomic discard transaction, and other filesystem writers are not
covered by the in-process lock. New/untracked-file deletion remains disabled
until it has its own safe workflow. Do not enable a hosted Discard action until
the durable journal and authorization contract above are implemented.

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

After the remote confirms a push, the review dialog MUST close and the editor
MUST announce **Pushed to origin · live status not verified** outside the
dialog. The user MUST still be able to check the pushed commit's deployment
status from the editor. A lost or failed push response MUST NOT be reported as a
confirmed push.

If the remote branch moves after the local one-page commit, the editor MUST keep
that commit and the saved source. It MUST NOT offer an identical blind retry
when the publisher reports that the remote branch changed. The current recovery
message MUST say that the local commit is retained and that this editor cannot
sync the branch automatically yet. A later sync flow may fast-forward or
reapply the reviewed one-page change only after checking for same-file
overlap. If the remote changed that file, the editor MUST show both versions
before any overwrite. A forced remote update MUST never trigger a forced push
or a silent hard reset.

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
The browser suite MUST include an isolated end-to-end path that edits boolean and
tag frontmatter through Page details, opens review with Rendered diff selected,
switches comparison tabs, edits the commit message, commits the selected page,
and pushes to a temporary bare Git remote. It MUST verify the exact remote commit
message and revision, and MUST NOT use a developer site or production remote.


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
  Initialize it only after the Full diff panel is visible so it can measure
  its container. The UI must report an initialization failure instead of
  leaving an empty viewer. Build the asset with tools/publish-diffs. The
  lightweight source diff must continue working without that asset.

The code is intentionally marked experimental until the exact Git workflow,
bundled asset, Playwright and deployment recovery paths are verified.

## Hosted Git credential isolation

Hosted deployments MUST mount the Git SSH Secret only into a short-lived init
container, which copies it to an ephemeral volume with owner-only permissions
for the dedicated Git publisher process. The publisher runtime UID MUST resolve
to a passwd entry for OpenSSH, and the private key MUST be owned by that UID.
The editor MUST NOT receive either key volume, `GIT_SSH_COMMAND`, or an SSH
agent socket. It may request remote readiness and push only through the
publisher's local Unix socket.

The publisher MUST use an operator-configured HTTPS GitHub origin and branch,
and MUST reject any checkout whose configured fetch origin differs from those
values. Network commands MUST use the fixed allowlisted SSH URL, regardless of
mutable checkout push settings. Readiness MUST verify that the requested local
HEAD and branch still match the checkout before reading the remote branch. The
push source MUST be the reviewed commit hash, not a mutable local branch name.
A push MUST be a non-forced fast-forward of exactly one commit whose sole
changed path is a regular Markdown blob; the remote branch MUST equal that
commit's parent. The publisher MUST verify the remote points to the requested
commit after pushing. If a retry finds that the remote already points to that
exact reviewed commit, it MUST acknowledge the existing push as successful
without pushing again.

The publisher MUST log each accepted or rejected push with the branch, short
commit ID, and outcome. Logs MUST NOT include document source, authentication
material, or request bodies.

The Unix socket MUST remain private to the pod and MUST NOT be exposed through
an HTTP or Kubernetes Service. Hlab Auth continues to authorize editor
requests. The publisher is a credential boundary, not an alternate browser
authentication path; a request to publish is accepted only through the editor's
existing authenticated review, commit, and push flow.

Hosted editor startup probes MUST allow at least ten minutes for initial
workspace discovery and filesystem watcher registration before restarting the
container. Readiness probes remain short after the HTTP listener starts.
