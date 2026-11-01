---
title: "Keyboard shortcuts and app keymaps"
description: "Customize Markata workbench navigation, writing, and publishing shortcuts."
date: 2026-11-01
published: true
tags: [documentation, editor, keyboard]
---

# Application keyboard shortcuts

In the Go workbench, open **Editor settings → Keyboard shortcuts…** or press **Ctrl/Cmd+K**, then **Ctrl/Cmd+S** with the VS Code preset. Search for commands, select a shortcut to record a new one, or reset your overrides. Markata rejects conflicting bindings.

**Application keymaps** are separate from the CodeMirror **editor keymap** (Default, Vim, Emacs, Sublime). Select Writer, VS Code, Vim, or Emacs for application actions without replacing your Markdown text-editing motions. App presets default to VS Code.

| Action | VS Code app preset | Writer app preset |
| --- | --- | --- |
| Command palette | Ctrl/Cmd+Shift+P | Ctrl/Cmd+K |
| Find file | Ctrl/Cmd+P | Ctrl/Cmd+P |
| Toggle sidebar | Ctrl/Cmd+B | Ctrl/Cmd+Backslash |
| Focus explorer | Ctrl/Cmd+Shift+E | Ctrl/Cmd+Shift+E |
| Focus editor | Ctrl/Cmd+Alt+1 | Ctrl/Cmd+Alt+1 |
| Focus preview | Ctrl/Cmd+Alt+2 | Ctrl/Cmd+Alt+2 |
| Preview-only toggle | Ctrl/Cmd+Alt+V | Ctrl/Cmd+Alt+V |
| Toggle line wrapping | Alt+Z | Alt+Z |
| Toggle Problems panel | Ctrl/Cmd+Shift+M | Ctrl/Cmd+Shift+M |
| Save now | Ctrl/Cmd+S | Ctrl/Cmd+S |
| New post from template | Ctrl/Cmd+Alt+N | Ctrl/Cmd+Alt+N |
| Review & Publish | Ctrl/Cmd+Shift+Enter | Ctrl/Cmd+Shift+Enter |
| Shortcut settings | Ctrl/Cmd+K, Ctrl/Cmd+S | Ctrl/Cmd+Alt+K |

A comma between key combinations denotes a **sequence**, not a simultaneous key press. On some platforms, an operating-system shortcut may take precedence. You can remap Markata actions through settings.

Vim application leader sequences require both **Vim app preset** and **Vim editor keymap**, in normal mode. Examples: Space E (explorer), Space F E (focus explorer), Space F N (new post), Space P V (preview), Space X (problems), and Space ? (shortcuts). They do not fire in insert mode or input controls.

The Emacs app preset provides Alt+X for the command palette and editor-independent actions while preserving Emacs editing commands. Normal text editing, including bold/link controls in Writer mode, remains the responsibility of CodeMirror.

Keyboard bindings are stored in browser local storage. The shortcut dialog reports conflicts and supports disabling an action with Backspace while recording. Markata never captures a shortcut while a modal is open, during IME composition, or when typing in a form field.

The **Review & Publish** shortcut opens the same read-only dialog as the Publish button. It cannot commit, push, deploy, or bypass publication confirmation. See [Editor publishing](../../spec/spec/EDITOR_PUBLISHING.md) for the deployment contract.

The Python and offline WASM workbenches retain their existing command palette shortcuts; this app-wide registry currently targets the Go site-backed workbench.
