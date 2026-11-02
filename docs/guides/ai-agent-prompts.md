---
title: "AI Agent Prompt Templates"
description: "Generate OpenCode GPT-5.4 prompt files from one or more GitHub issues"
date: 2026-03-12
published: true
tags:
  - documentation
  - ai
  - prompts
  - automation
---


markata-go includes reusable GPT-5.4 prompt templates for OpenCode workflows.
Use the generator script to fetch issue metadata from GitHub and render prompt
files for either a single issue or a group of issues.

## Files

- `continue-prompt-gpt-54.template.txt` - multi-issue template
- `continue-prompt-gpt-54-single-issue.template.txt` - single-issue template
- `scripts/fill_continue_prompts.py` - generator script

## Generate Prompt Files

Generate a multi-issue prompt and one single-issue prompt:

```bash
python scripts/fill_continue_prompts.py 942
```

Generate a multi-issue prompt for several issues plus one single-issue prompt
per issue:

```bash
python scripts/fill_continue_prompts.py 935 936 937 938 939
```

Write the rendered content back to the root prompt filenames too:

```bash
python scripts/fill_continue_prompts.py 942 --write-root-prompts
```

## Options

Use a different worktree root hint:

```bash
python scripts/fill_continue_prompts.py 942 --worktree-root ../worktrees
```

Write outputs to a different directory:

```bash
python scripts/fill_continue_prompts.py 942 --output-dir /tmp/prompts
```

Use a different repository:

```bash
python scripts/fill_continue_prompts.py 942 --repo owner/repo
```

## Output Naming

For issue `942`, the script writes:

- `continue-prompt-gpt-54-issues-942.txt`
- `continue-prompt-gpt-54-issue-942.txt`

For issues `935 936 937`, the script writes:

- `continue-prompt-gpt-54-issues-935-936-937.txt`
- `continue-prompt-gpt-54-issue-935.txt`
- `continue-prompt-gpt-54-issue-936.txt`
- `continue-prompt-gpt-54-issue-937.txt`

## Notes

- The script requires the GitHub CLI and access to the target repository.
- Single-issue prompts include a suggested branch name and worktree path hint.
- The templates stay editable as plain text, so prompt tuning does not require
  Python changes.
