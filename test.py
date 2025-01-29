#!/usr/bin/env -S uv run --script
# /// script
# requires-python = ">=3.12"
# dependencies = [
#     "markata==0.8.1",
#     "setuptools",
# ]
# ///

from markata.cli import cli

cli()
