# Markata autonomous release train

Tracking: #1747

## Operating model

- **Daytime CI:** keep fast, required Linux checks; measure existing Windows/macOS critical-path times before removing required checks.
- **Nightly, 02:00 America/Chicago:** GitHub Actions owns full matrix, prerelease builds and publication. A coding agent is *optional* and may only prepare narrow PRs for failures; release publication requires deterministic checks.
- **05:00 America/Chicago:** ChatGPT scheduled triage reviews nightly results, identifies blockers and escalates uncertain problems.
- **09:00 America/Chicago:** ChatGPT scheduled issue-to-PR pass validates candidate issues and attempts fixes, with review gates.
- **Weekly hardened release:** full compatibility, package install checks, provenance and signed artifacts before stable publication.

## Human-vetting policy

Issues should have a reproducible problem statement, expected behavior, affected version, scope and risk assessment. Agents may label/triage and propose PRs; ambiguous reports should receive a question rather than a speculative patch.

Low-risk changes can be considered for auto-merge only after required checks, review, and branch protections pass. Never auto-merge authentication, permissions, publishing, release workflows, or dependency supply-chain changes.

## Copilot authentication

This is a personal repository, so GitHub's organization-only `copilot-requests: write` token billing path is **not** appropriate. The GitHub Agentic Workflows documentation requires a fine-grained personal token stored as `COPILOT_GITHUB_TOKEN` for personal repositories. Confirm account quotas and PAT permissions before enabling unattended inference.

## Implementation gates

1. Benchmark 20+ recent runs (queue time, execution time, critical path, p50/p90) and prove or disprove the Windows bottleneck.
2. Split fast required checks from full optional/nightly matrix; preserve all branch protections.
3. Add scheduled full-matrix workflow with manual dry-run and concurrency guard; support Central DST by explicit timezone handling.
4. Add gh-aw Markdown agent workflow and compile its lockfile using `gh aw compile`. Require read-only access plus narrowly scoped safe outputs.
5. Verify Copilot auth, run a dry-run issue-to-PR test, and audit permissions/usage.
6. Add prerelease publishing after successful full tests; separate stable weekly promotion and protected approval.
7. Document rollback and verify no duplicate publishing or runaway repair loops.

Do not enable publication or auto-merge until each gate is verified.
