# Site-owned custom plugin starter

Issue #2098; agent-first initiative #2087.

## Contract

The stock Markata CLI cannot load an arbitrary `plugins/` directory at
runtime. A site may build its own Go executable with a pinned dependency on
a Markata release, register an implementation via
`plugins.RegisterPluginConstructor(name, func() lifecycle.Plugin)`, and
invoke the normal `cmd.Execute()` CLI. The site's `hooks` list opts in by
name. Markata retains its existing build, config, serve, and output formats.

The starter under `examples/custom-plugin/` must compile on the repository's
required Go version, define one deterministic Transform-stage plugin, skip
skipped posts, and have regression tests. It should need no Node or separate
deployment service, and should not silently alter the stock binary.

## Agent behavior

- Prefer config, frontmatter, CSS and local template overrides before code.
- Pin module and plugin dependencies; add site-owned tests.
- Test full production output and inspect changed files before asking for
  publish approval; do not access credentials or deploy without consent.
- For rollback, remove the configured hook and return to the stock binary.
- Do not claim compile-time Go registration is dynamic site-local plugin
  loading; design any future runtime protocol separately, with a security and
  determinism contract.
