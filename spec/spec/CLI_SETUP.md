# Optional workspace bootstrap (`markata setup`)

Status: minimum noninteractive implementation contract, issue #2088.

`markata setup [workspace]` is a small convenience, **not** a hosting,
identity, template-selection, or deployment wizard. For agent-chosen setup use
`markata --skill`. The core must remain usable without running this command.

## Behavior

- An omitted workspace means the current directory. A supplied workspace
  directory may be created if missing, without touching other projects.
- Existing supported local config filenames (currently `markata-go.toml`,
  `markata-go.yaml`, `markata-go.yml`, `markata-go.json`) are preserved
  unchanged on repeated invocations.
- An unrecognized `markata.toml`/YAML/JSON config triggers a clear diagnostic;
  never silently create a competing config that Markata itself cannot discover.
- In a new empty workspace, create the built-in default TOML file atomically
  with exclusive creation; never overwrite a concurrent or existing file.
- Print the workspace path and immediate `markata editor` and `markata serve`
  commands. Do not claim `markata up` is available until it is shipped.
- No public publish/deploy, agent installation, network access or credential
  changes. Run without a TUI and in noninteractive agent environments.

## Tests

Create/idempotence, alternate config preservation, unsupported config rejection,
invalid file-as-workspace, lint and cross-platform CI.
