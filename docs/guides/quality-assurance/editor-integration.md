---
title: "Editor Integration"
description: "Set up markata-go LSP features and lint feedback in your editor"
date: 2026-01-24
published: true
tags:
  - guides
  - quality-assurance
  - lint
  - vim
  - neovim
  - vscode
---


Use the Language Server Protocol (LSP) for live wikilink completions, broken-link diagnostics, hover details, and go-to-definition. Use `markata-go lint` when you also need content-quality checks in an editor's quickfix or Problems panel.

## Start With A Health Check

From the root of your site, run:

```bash
markata-go lsp doctor
```

The command verifies the local binary and validates any discovered
markata-go configuration. It does not modify your editor configuration. Your
editor—not the configuration file—supplies the LSP workspace root.

It also detects supported editors installed on your `PATH`. For Neovim, Helix,
and Emacs, it loads the editor configuration for a deeper check by default:

```bash
markata-go lsp doctor
```

This can run editor startup and plugin code. Skip it when needed while still
detecting editors:

```bash
markata-go lsp doctor --no-verify-editor
```

The report labels each detected editor as configured, unconfigured, or
inconclusive. VS Code, Cursor, and Zed are detected but do not have stable
headless LSP verification APIs.

When verification is skipped and a headless-verifiable editor is detected, the
report prints the default doctor follow-up command to stderr. This keeps the
primary diagnostic report on stdout while making the next interactive action
easy to spot.

In an interactive terminal, the report uses your configured site palette for
pass, information, warning, and failure states. The textual labels remain
present for accessibility and plain output.

Every LSP client uses the same stable values:

| Setting | Value |
| --- | --- |
| Command | `markata-go lsp` |
| File type | Markdown |
| Workspace root | Directory containing `markata-go.toml` (or the site root) |

Print a maintained starting snippet with:

```bash
markata-go lsp setup --editor neovim
```

Omit `--editor` to print snippets for every supported editor installed on your
`PATH`:

```bash
markata-go lsp setup
```

When more than one editor is found, each snippet has a clear heading. Interactive
terminal output is syntax highlighted; redirected output remains plain text for
copying into configuration files.

Available editors are `generic`, `neovim`, `helix`, `emacs`, `zed`, and
`vscode`. The command only prints guidance: it deliberately does not rewrite
editor settings because editor versions, extensions, and configuration layouts
vary. When using an agent, have it inspect the editor configuration already in
the repository or home directory, then adapt the generated snippet.

## Language Server Setup

### Neovim

For Neovim 0.11 or later, add the output of this command to `init.lua` or a
loaded Lua module:

```bash
markata-go lsp setup --editor neovim
```

### Helix

Append the generated TOML to `~/.config/helix/languages.toml`:

```bash
markata-go lsp setup --editor helix
```

Restart Helix after changing the file.

### Emacs

With Eglot enabled, add the generated Elisp to your Emacs init file:

```bash
markata-go lsp setup --editor emacs
```

### VS Code, Cursor, And Zed

These editors need an installed LSP client or a Markata-aware extension before
they can launch an arbitrary stdio language server. Print the required command,
arguments, file type, and workspace-root values, then apply them using the
configuration format required by the installed extension:

```bash
markata-go lsp setup --editor vscode
markata-go lsp setup --editor zed
```

Do not use an editor task as an LSP replacement: tasks run once and cannot
provide completion, hover, or navigation.

### Other Editors

Use the generic contract when your editor is not listed:

```bash
markata-go lsp setup --editor generic
```

Configure the editor's LSP client to launch `markata-go lsp` over standard I/O
for Markdown documents. Use its native workspace-root setting to point at the
site root.

## Lint Integration

The `markata-go lint` command outputs issues in a format compatible with many editors' error navigation features. The remaining sections show how to integrate those checks with popular editors.

## Vim / Neovim Quickfix

Vim and Neovim have built-in quickfix support for navigating errors and warnings.

### Basic Usage

```vim
" Set markata-go as the make program
:set makeprg=markata-go\ lint

" Run lint and populate quickfix
:make docs/**/*.md

" Open quickfix window
:cwindow

" Navigate errors
:cnext      " Go to next error
:cprev      " Go to previous error
:cfirst     " Go to first error
:clast      " Go to last error
```

### Custom Error Format

For perfect parsing of markata-go lint output, add this to your `~/.vimrc` or `init.vim`:

```vim
" Add markata-go lint error format
set errorformat+=%f:
set errorformat+=\ \ %tarning\ [line\ %l\\,\ col\ %c]:\ %m
set errorformat+=\ \ %trror\ [line\ %l\\,\ col\ %c]:\ %m
```

### One-liner Command

Quick way to run lint and populate quickfix:

```vim
:cgetexpr system('markata-go lint docs/**/*.md')
:cwindow
```

### Neovim Lua Configuration

For Neovim users with Lua config:

```lua
-- ~/.config/nvim/lua/lint.lua or init.lua

-- Function to run markata-go lint
local function markata_lint()
  vim.cmd('cgetexpr system("markata-go lint docs/**/*.md")')
  vim.cmd('cwindow')
end

-- Create a user command
vim.api.nvim_create_user_command('MarkataLint', markata_lint, {})

-- Optional: keybinding
vim.keymap.set('n', '<leader>ml', markata_lint, { desc = 'Run markata-go lint' })
```

### Lint Current File

To lint only the current file:

```vim
:cgetexpr system('markata-go lint ' . expand('%'))
```

Or in Lua:

```lua
local function lint_current_file()
  local file = vim.fn.expand('%')
  vim.cmd('cgetexpr system("markata-go lint ' .. file .. '")')
  vim.cmd('cwindow')
end
```

## VS Code Integration

VS Code can use tasks with problem matchers to integrate with markata-go lint.

### Task Configuration

Create `.vscode/tasks.json` in your project:

```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "Lint Markdown",
      "type": "shell",
      "command": "markata-go lint docs/**/*.md",
      "group": "build",
      "problemMatcher": {
        "owner": "markata-go",
        "fileLocation": ["relative", "${workspaceFolder}"],
        "pattern": [
          {
            "regexp": "^(.+):$",
            "file": 1
          },
          {
            "regexp": "^\\s+(warning|error)\\s+\\[line\\s+(\\d+),\\s+col\\s+(\\d+)\\]:\\s+(.+)$",
            "severity": 1,
            "line": 2,
            "column": 3,
            "message": 4,
            "loop": true
          }
        ]
      }
    },
    {
      "label": "Lint Current File",
      "type": "shell",
      "command": "markata-go lint ${file}",
      "group": "build",
      "problemMatcher": {
        "owner": "markata-go",
        "fileLocation": ["relative", "${workspaceFolder}"],
        "pattern": [
          {
            "regexp": "^(.+):$",
            "file": 1
          },
          {
            "regexp": "^\\s+(warning|error)\\s+\\[line\\s+(\\d+),\\s+col\\s+(\\d+)\\]:\\s+(.+)$",
            "severity": 1,
            "line": 2,
            "column": 3,
            "message": 4,
            "loop": true
          }
        ]
      }
    }
  ]
}
```

### Running the Task

1. Press `Ctrl+Shift+B` (or `Cmd+Shift+B` on Mac)
2. Select "Lint Markdown" from the task list
3. Errors appear in the Problems panel
4. Click on errors to jump to the file and line

### Keyboard Shortcut

Add a keybinding in `keybindings.json`:

```json
{
  "key": "ctrl+shift+l",
  "command": "workbench.action.tasks.runTask",
  "args": "Lint Markdown"
}
```

## Emacs Integration

Emacs can use `compile-mode` to run lint and navigate errors.

### Basic Usage

```elisp
;; Run markata-go lint
M-x compile RET markata-go lint docs/**/*.md RET

;; Navigate errors
M-g n  ; next-error
M-g p  ; previous-error
```

### Custom Compilation Regexp

Add to your Emacs config:

```elisp
(add-to-list 'compilation-error-regexp-alist-alist
             '(markata-go
               "^\\(.*\\):\n  \\(warning\\|error\\) \\[line \\([0-9]+\\), col \\([0-9]+\\)\\]: \\(.*\\)$"
               1 3 4 2))

(add-to-list 'compilation-error-regexp-alist 'markata-go)
```

### Flycheck Integration

For on-the-fly linting with Flycheck:

```elisp
(flycheck-define-checker markata-go
  "A markdown linter using markata-go."
  :command ("markata-go" "lint" source)
  :error-patterns
  ((warning line-start (file-name) ":\n"
            "  warning [line " line ", col " column "]: " (message) line-end)
   (error line-start (file-name) ":\n"
          "  error [line " line ", col " column "]: " (message) line-end))
  :modes (markdown-mode gfm-mode))

(add-to-list 'flycheck-checkers 'markata-go)
```

## Sublime Text Integration

Create a build system for markata-go lint.

### Build System Configuration

Create `Packages/User/markata-go.sublime-build`:

```json
{
  "cmd": ["markata-go", "lint", "$file"],
  "working_dir": "$project_path",
  "file_regex": "^(.+):\n  (warning|error) \\[line (\\d+), col (\\d+)\\]: (.*)$",
  "selector": "text.html.markdown"
}
```

### Usage

1. Open a Markdown file
2. Press `Ctrl+B` (or `Cmd+B` on Mac)
3. Use `F4` / `Shift+F4` to navigate errors

## Tips and Best Practices

### Lint on Save

Most editors support running commands on file save. Configure your editor to run `markata-go lint` on the current file when saving for immediate feedback.

**Vim/Neovim autocommand:**

```vim
augroup MarkataLint
  autocmd!
  autocmd BufWritePost *.md cgetexpr system('markata-go lint ' . expand('%'))
augroup END
```

**VS Code setting (settings.json):**

```json
{
  "emeraldwalk.runonsave": {
    "commands": [
      {
        "match": "\\.md$",
        "cmd": "markata-go lint ${file}"
      }
    ]
  }
}
```

### Project-wide Lint

For CI/CD or pre-commit hooks, lint all files:

```bash
markata-go lint
```

This uses the glob patterns from your `markata-go.toml` configuration.

### Exit Codes

The lint command uses standard exit codes:

| Code | Meaning |
|------|---------|
| 0 | No issues found |
| 1 | Issues found |
| 2 | Error running lint |

Use these in scripts:

```bash
if markata-go lint; then
  echo "All good!"
else
  echo "Issues found, check the output"
fi
```
