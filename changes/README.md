# Release-note fragments

For user-visible PRs, add one Markdown fragment under `changes/unreleased/` named `<PR-number>.<category>.md`. Categories: `feature`, `fix`, `performance`, `breaking`, `other`. Documentation-only, test-only, and internal CI changes can omit a fragment.

Write one concise, user-facing sentence (not a heading). Describe the user impact. A breaking change should include migration guidance.

Example `changes/unreleased/1748.fix.md`:

```md
Fix the archive next-page navigation when editorial presentation is enabled.
```

`python3 scripts/release-notes.py --version v0.15.0` assembles the available fragments in deterministic category/PR order. The assembler validates filenames and content. **It does not yet verify that every merged user-visible PR has a fragment.** The current backlog predates this convention, so releases must also review merged changes and write curated highlights.

The scheduled release workflow produces disposable snapshot artifacts and note previews only. Its daily nightly and Tuesday/Friday stable-candidate runs never create tags or publish to GitHub or PyPI. Humans must review candidate artifacts and tag/publish using the existing controlled release processes after the gates pass.

After a stable release, archive consumed fragments in a reviewed PR. Never mutate release tags or remove fragments from a release runner. Historical releases can use GoReleaser's GitHub changelog as a starting point.

See `scripts/release-notes.py` for the deterministic assembler.
