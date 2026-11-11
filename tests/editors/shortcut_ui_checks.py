"""Browser checks for Go workbench application shortcuts and dialogs."""
from playwright.sync_api import expect


def check_app_shortcuts(page, url):
    page.goto(url)
    expect(page.locator('#site-sidebar')).to_be_visible()
    expect(page.locator('.file-picker, .file-finder')).to_have_count(1, timeout=15000)
    page.keyboard.press('Control+Shift+P')
    expect(page.locator('.markdown-commands[open]')).to_be_visible()
    page.keyboard.press('Escape')
    expect(page.locator('.markdown-commands[open]')).to_have_count(0)

    # The current site tree supplies a richer fuzzy picker; the legacy file
    # finder remains a fallback for older workspaces.
    page.keyboard.press('Control+P')
    expect(page.locator('.file-picker[open], .file-finder[open]')).to_be_visible()
    if page.locator('#files-query').count():
        expect(page.locator('#files-query')).to_be_focused()
        page.locator('#files-query').fill('welcome')
        expect(page.locator('#files-results [role=option]').first).to_contain_text('welcome')
    page.keyboard.press('Escape')
    expect(page.locator('.file-picker[open], .file-finder[open]')).to_have_count(0)

    # The new picker also remains reachable through the shared command palette.
    page.keyboard.press('Control+Shift+P')
    expect(page.locator('.markdown-commands[open]')).to_be_visible()
    page.locator('#commands-query').fill('open file')
    page.keyboard.press('Enter')
    expect(page.locator('#files-query')).to_be_focused()
    page.keyboard.press('Escape')

    page.keyboard.press('Control+B')
    expect(page.locator('#workbench')).to_have_class('sidebar-hidden', timeout=5000)
    page.keyboard.press('Control+B')
    expect(page.locator('#workbench')).not_to_have_class('sidebar-hidden', timeout=5000)

    # Hiding a focused explorer must move focus to visible editor content.
    page.locator('#site-view').focus()
    page.keyboard.press('Control+B')
    expect(page.locator('#site-sidebar')).to_have_attribute('aria-hidden', 'true')
    assert page.evaluate("document.activeElement.closest('.CodeMirror') !== null")
    page.keyboard.press('Control+B')
    expect(page.locator('#site-sidebar')).to_have_attribute('aria-hidden', 'false')

    page.keyboard.press('Control+Alt+2')
    expect(page.locator('#preview-pane')).to_be_focused()
    page.keyboard.press('Control+Alt+1')
    assert page.evaluate("document.activeElement.closest('.CodeMirror') !== null")

    # Layout select lives in a modal; open it before interacting.
    page.locator('#settings-open').click()
    page.locator('#setting-layout').select_option('editor')
    page.locator('#settings-close').click()
    page.set_viewport_size({'width': 960, 'height': 800})
    page.keyboard.press('Control+B')
    assert len(page.locator('#workbench').evaluate("e => getComputedStyle(e).gridTemplateColumns.split(' ')")) == 2
    page.keyboard.press('Control+B')
    page.locator('#settings-open').click()
    page.locator('#setting-layout').select_option('split')
    page.locator('#settings-close').click()
    page.set_viewport_size({'width': 1280, 'height': 720})

    # Re-establish in-page editor focus after closing a modal. Browser-level
    # Ctrl+K shortcuts can otherwise move focus to the address bar.
    page.locator('.CodeMirror').click()
    page.evaluate("""() => {
      window.__shortcutKeyTrace = [];
      window.__shortcutErrors = [];
      window.addEventListener('error', e => window.__shortcutErrors.push(e.message));
      window.addEventListener('unhandledrejection', e => window.__shortcutErrors.push(String(e.reason)));
      window.addEventListener('keydown', e => {
        const state = {
          key:e.key,ctrl:e.ctrlKey,meta:e.metaKey,shift:e.shiftKey,
          alt:e.altKey,repeat:e.repeat,composing:e.isComposing,
          initialPrevented:e.defaultPrevented,
          focus:document.activeElement?.tagName, active:document.activeElement?.id
        };
        window.__shortcutKeyTrace.push(state);
        queueMicrotask(() => {
          state.finalPrevented = e.defaultPrevented;
          state.stopped = e.cancelBubble;
          state.dialogs = [...document.querySelectorAll('dialog[open]')].map(d => d.id);
        });
      }, true);
    }""")
    page.keyboard.press('Control+K')
    page.keyboard.press('Control+S')
    assert page.locator('#shortcuts-dialog[open]').is_visible(), page.evaluate("""() => ({
      preset:window.MarkataShortcuts?.getPreset(),
      binding:window.MarkataShortcuts?.bindingFor('workbench.shortcutSettings'),
      events:window.__shortcutKeyTrace,
      dialogs:[...document.querySelectorAll('dialog[open]')].map(e=>e.id),
      errors:window.__shortcutErrors,
      active:document.activeElement?.outerHTML?.slice(0,200)
    })""")
    expect(page.locator('#shortcuts-list')).to_contain_text('Publishing: Review changes')
    recorder = page.locator('[data-command-id="workbench.save"]')
    recorder.click()
    page.keyboard.press('Control+Alt+J')
    expect(recorder).to_have_text('Mod+Alt+J')
    expect(recorder).to_be_focused()
    page.locator('#shortcuts-reset').click()
    expect(recorder).to_have_text('Mod+S')
    page.locator('#shortcuts-preset').select_option('writer')
    page.locator('#shortcuts-close').click()

    page.keyboard.press('Control+K')
    expect(page.locator('.markdown-commands[open]')).to_be_visible()
    page.keyboard.press('Escape')
    page.locator('#settings-open').click()
    page.locator('.shortcuts-settings-link').click()
    expect(page.locator('#shortcuts-dialog[open]')).to_be_visible()
    page.locator('#shortcuts-preset').select_option('vscode')
    page.locator('#shortcuts-close').click()

    # Editor text keymaps are separate from the application-shortcut preset
    # and can be selected directly from the shared command palette.
    app_preset = page.evaluate("window.MarkataShortcuts.getPreset()")
    for keymap_name in ('vim', 'emacs', 'sublime', 'default'):
        page.locator('#commands').click()
        page.locator('#commands-query').fill('editor keymap ' + keymap_name)
        option = page.locator('#commands-results [role="option"]')
        expect(option).to_have_count(1)
        expect(option.first).to_contain_text('Editor keymap: ' + keymap_name.capitalize())
        page.keyboard.press('Enter')
        expect(page.locator('#keymap')).to_have_value(keymap_name)
        # Switching into Vim focuses the editor in insert mode. CodeMirror
        # intentionally changes the active map to "vim-insert", while the
        # persisted user preference remains "vim".
        page.wait_for_function(
            """(key) => {
              const active = document.querySelector('.CodeMirror').CodeMirror.getOption('keyMap');
              return key === 'vim' ? (active === 'vim' || active === 'vim-insert') : active === key;
            }""",
            arg=keymap_name,
            timeout=5000,
        )
        assert page.evaluate("localStorage.getItem('markata-editor-keymap')") == keymap_name
        assert page.evaluate("window.MarkataShortcuts.getPreset()") == app_preset

    page.keyboard.press('Control+Shift+Enter')
    expect(page.locator('#publish-review[open]')).to_be_visible()
    expect(page.locator('#publish-confirm')).to_be_disabled()
    page.keyboard.press('Escape')
    expect(page.locator('#publish-review[open]')).to_have_count(0)
