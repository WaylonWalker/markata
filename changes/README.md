# Release-note fragments

Every user-visible PR adds one Markdown fragment under `changes/unreleased/` named `<PR-number>.<category>.md`, where category is `feature`, `fix`, `performance`, `breaking`, or `other`. Documentation-only, tests-only, and internal CI changes can omit a fragment.

Write one concise, user-facing sentence, with no heading. Explain impact, not implementation. For breaking changes, include migration instructions in the fragment.

Example: `changes/unreleased/1748.fix.md`:

```md
Fix the archive next-page navigation when editorial presentation is enabled.
```

The release workflow should assemble fragments deterministically, sorted by category and PR number, into version-specific notes. Nightly releases preview accumulated fragments; stable releases publish the same curated notes, then a follow-up PR archives consumed fragments. **Never delete or move fragments automatically on a release runner**: release tags should remain immutable, and the archival change should be reviewed as a normal PR.

A missing fragment for a user-visible change should fail the release-notes validation gate. Historical releases can fall back to GoReleaser's GitHub changelog until fragments are populated.

See `scripts/release-notes.py` for deterministic assembly.
