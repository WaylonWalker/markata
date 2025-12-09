---
date: 2025-12-09
description: "The plugin adds performance profiling capabilities using pyinstrument.
  It generates detailed HTML reports showing where your build spends time. Installation\u2026"
published: false
slug: markata/plugins/pyinstrument
title: pyinstrument.py


---

---

The `markata.plugins.pyinstrument` plugin adds performance profiling capabilities using
pyinstrument. It generates detailed HTML reports showing where your build spends time.

## Installation

This plugin is built-in but not enabled by default. Add it to your plugins list:

```toml
hooks = [
    "markata.plugins.pyinstrument",
]
```

You must also install pyinstrument:
```bash
pip install pyinstrument
```

## Uninstallation

Remove the plugin from your hooks list in `markata.toml`:

```toml
hooks = [
    # Remove or comment out the line below
    # "markata.plugins.pyinstrument",
]
```

## Configuration

Configure profiling in `markata.toml`:

```toml
[markata.profiler]
# Enable/disable profiling
should_profile = true

# Output location (relative to output_dir)
output_file = "_profile/index.html"

# Profile options
interval = 0.001
async_mode = "enabled"
show_all = false
timeline = false
```

## Functionality

## Profiling Features

The plugin:
1. Profiles the entire build process
2. Generates HTML reports
3. Shows time distribution
4. Identifies bottlenecks

## Report Generation

Creates reports with:
- Call tree visualization
- Time percentages
- Function details
- Stack traces

### Configuration Options

Supports:
- Custom output paths
- Sampling intervals
- Async mode settings
- Display options
- Timeline view

## Performance Impact

Note:
- Minimal overhead
- Configurable precision
- Optional async profiling
- Selective profiling

## Dependencies

This plugin depends on:
- pyinstrument for profiling
- pydantic for configuration

---

!!! function
    <h2 id="save" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">save <em class="small">function</em></h2>

    stop the profiler and save as late as possible

???+ source "save <em class='small'>source</em>"
    ```python
    def save(markata: Markata) -> None:
        "stop the profiler and save as late as possible"
        if markata.config.profiler.should_profile:
            if markata.config.profiler.profiler is not None:
                if markata.config.profiler.profiler.is_running:
                    try:
                        markata.config.profiler.profiler.stop()
                        html = markata.config.profiler.profiler.output_html()
                        markata.config.profiler.output_file.write_text(html)
                        markata.console.print(
                            markata.config.profiler.profiler.output_text()
                        )

                    except AttributeError:
                        markata.console.log(
                            "profiler not available, skipping save pyinstrument save",
                        )
                        markata.console.log(
                            r"[red]to enable profiler [wheat1][itallic]pip install 'markata[pyinstrument]'",
                        )
    ```
!!! function
    <h2 id="teardown" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">teardown <em class="small">function</em></h2>

    stop the profiler on exit

???+ source "teardown <em class='small'>source</em>"
    ```python
    def teardown(markata: Markata) -> None:
        "stop the profiler on exit"
        # import logging

        # logger = logging.getLogger()
        # logger.handlers.clear()
        if markata.config.profiler.should_profile:
            if markata.config.profiler.profiler is not None:
                if markata.config.profiler.profiler.is_running:
                    markata.config.profiler.profiler.stop()
    ```
!!! method
    <h2 id="ensure_output_dir_exists" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">ensure_output_dir_exists <em class="small">method</em></h2>

    Ensure output directory exists, creating it if necessary.

???+ source "ensure_output_dir_exists <em class='small'>source</em>"
    ```python
    def ensure_output_dir_exists(cls, v: Union[str, Path]) -> Path:
            """Ensure output directory exists, creating it if necessary."""
            if isinstance(v, str):
                v = Path(v)
            v.mkdir(parents=True, exist_ok=True)
            return v
    ```