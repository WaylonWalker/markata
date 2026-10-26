# repo-steward · Workspace lifecycle

Scan what exists, make a deliberate decision for every project, and apply only
the changes you have reviewed. This example is a **visual Markdown document**,
not a live status report; it invents no counts, completion percentages, or repo states.

::: md-layout

## One workspace. Two deliberate choices.

:::: md-columns {data-gap="lg"}

::::: md-card {data-align="center"}

:::::: md-icon {data-size="lg"}
:lucide-folder-search:
::::::

### Inspect before acting

Discover branches, worktree changes, ignored files, and missing Git metadata.
No writes are required to understand the current state.

:::::

::::: md-card {data-align="center"}

:::::: md-icon {data-size="lg"}
:lucide-shield-check:
::::::

### Apply only after review

Keep/archive decisions persist. Preparation and ignored-file cleanup have
explicit preview and apply stages.

:::::

::::

:::: md-status {data-tone="caution"}

Archiving or cleaning can remove local files. Read the preview and confirmation
before applying either operation.

::::

:::

## The common route

::: road steps

- [Discover](#discover) 01
- [Decide](#decide) 02
- [Prepare](#prepare) 03
- [Clean up](#clean-up) 04
- [Publish](#publish) 05

:::

::: md-feature

:::: md-feature-icon
:lucide-scan-search:
::::

:::: md-feature-body

### Discover

`repo-steward scan` inspects each immediate project directory. The interactive
dashboard provides searchable project status and persisted decisions.

::::

:::

::: md-feature

:::: md-feature-icon
:lucide-list-checks:
::::

:::: md-feature-body

### Decide

Mark each project **Keep**, **Archive**, or leave it undecided. The review state
survives restarts. Archive only after reviewing the requested move and confirming
the shared destination is mounted and available.

::::

:::

## After the decision

::: md-columns {data-gap="lg"}

:::: md-card

::::: md-icon {data-size="lg"}
:lucide-git-branch:
:::::

### Keep the project

Use `repo-steward prepare` to preview Git initialization, `.gitignore`
changes, and Forgejo remotes. `prepare --apply` performs the approved setup.

**Next:** review `cleanup` before using `cleanup --apply` to remove ignored
files, then commit intentional changes before `publish` pushes `main`.

::::

:::: md-card

::::: md-icon {data-size="lg"}
:lucide-archive:
:::::

### Archive the project

`repo-steward review` can copy the project to the Walkershare archive,
checksum-verify the copy, then remove the local copy. It requires the mount and
`rsync`; existing archive destinations are not overwritten.

**Next:** retain the archive path in the saved review state for later recovery.

::::

:::

::: md-feature

:::: md-feature-icon
:lucide-folder-git:
::::

:::: md-feature-body

### Prepare

For kept projects, `repo-steward prepare` shows the planned changes; only
`prepare --apply` changes Git metadata or ignore rules. It never commits or
pushes for you.

::::

:::

::: md-feature

:::: md-feature-icon
:lucide-file-check:
::::

:::: md-feature-body

### Clean up

`repo-steward cleanup` previews branch switches and ignored paths. The
`--apply` form re-previews on `main` and uses `git clean -fdX` after
confirmation. Treat ignored-file deletion as an explicit destructive gate.

::::

:::

::: md-feature

:::: md-feature-icon
:lucide-upload:
::::

:::: md-feature-body

### Publish

Commit the content you actually want to keep, then use `repo-steward publish`
to push committed `main` branches to Forgejo. Publishing is optional.

::::

:::

## Useful side routes

::: md-columns

:::: md-card

### Clone a remote

`repo-steward clone` lists GitHub and Forgejo repositories, identifies known
local copies, and clones a selected remote into your project root.

::::

:::: md-card

### Resume later

Start the dashboard or `repo-steward review` again. Saved decisions remain and
undecided projects can be revisited without repeating completed review work.

::::

:::

> **Rendering note:** This is one portable Markdown source file. Markata turns
> the `md-*`, `road`, and `:lucide-*:` syntax into theme-aware layout in a
> Markata site. A generic Markdown viewer may display the container fences
> literally. Do not claim a visual preview was tested without building a site.
