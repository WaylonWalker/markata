# Markata — agent-first site setup

You are helping someone set up or customize a **Markata static Markdown site**. The user chose their agent; you do **not** need to install or launch another agent.

Markata is a deterministic site builder and authoring toolkit. You may shape the user's site with Markdown, configuration, templates, assets, and plugins. **Do not assume the hosting platform, content design, or workflow.** Discover the user's goals and existing files before editing anything.

## Start without side effects

1. Run `markata version` and `markata --help`. If the binary is missing, ask the user to install it from a trusted source.
2. Inspect the working directory and check for `markata-go.toml`, `markata-go.yaml`, `markata-go.yml`, and `markata-go.json`. Do not clobber or rename existing configurations.
3. If a Markata site exists, run `markata config show` and inspect `posts/`, `pages/`, `docs/`, `static/`, `templates/`, `palettes/`, and any custom plugin source that exists. Some commands or plugin code may execute user-defined code: use appropriate trust boundaries.
4. If it is a new site, ask the user which kind of site they want, inspect an empty or chosen directory, and use `markata config init` to create a minimal starting point. Avoid unnecessary scaffolding.
5. Ask about **publishing destination** only when the user wants deployment. Examples: GitHub Pages, Cloudflare, a VPS, or their own Kubernetes cluster. Do not choose or provision a provider unprompted.

## Choose the smallest extension

- **Content** — Markdown with frontmatter, feeds, post types.
- **Branding** — prefer palette/theme settings and CSS before adding a whole template system.
- **Layout** — override a template only where the default needs changing.
- **Behavior** — prefer existing configuration and hooks; implement a custom plugin only if native capabilities are insufficient.
- **Deployment** — build static files and use the target's existing deployment method. Markata does not need to manage an account, VPN, or Kubernetes cluster.

For the complete version-matched skill and its references, run `markata agent install` (or `markata agent install --agent codex`, `--agent claude-code`, `--agent amp`). Then read the installed skill's `SKILL.md` and the relevant topic file, especially `topics/template-management.md`, `topics/theme-creation.md`, `topics/plugin-creation.md`, `topics/build-deployment.md`, and `topics/faster-builds.md`. Check `markata agent doctor` for drift; do not overwrite customized files without consent.

## Build and verify deterministically

1. Inspect the current config and plugin dependencies. Pin versions and inputs for repeatability; do not fetch changing remote content silently during release builds.
2. Iterate with `markata build --fast` only when speed is the goal; **fast mode skips some production processing** and cannot substitute for a full release build.
3. Before proposing deployment, run a full `markata build` with the intended production configuration. If you are validating reproducibility, compare fresh outputs and hashes on the same pinned inputs in separate clean directories. Report any nondeterminism instead of promising byte-identical output without proof.
4. Show build warnings, execution time, relevant output paths, and what changed. Use `markata build --benchmark-summary-json` when structured build timings are useful; confirm the installed version accepts the flag.
5. Preview locally with `markata serve` or `markata editor` when needed. Running a preview is **not** authorization to publish.
6. Publish/deploy **only after the user explicitly approves the target and operation**. Keep secrets in the existing environment/secret store; never paste tokens into Markdown, prompts, logs, or repository files. Verify the actual deployed revision before claiming success.

## Handoff to the user

Present a concise summary of: site changes, commands executed, reproducibility evidence, build timing, deployment target (or none), and what the user can change next. Prefer leaving editable source and simple commands over a bespoke control plane.

If the user's request is ambiguous, show a minimal plan with the most important question, then adapt to their answer. **Do not build an infrastructure platform on their behalf.**
