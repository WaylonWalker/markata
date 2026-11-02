"""Browser checks for Go workbench application shortcuts and dialogs."""
from playwright.sync_api import expect


def check_app_shortcuts(page, url):
    page.goto(url)
    expect(page.locator('#site-sidebar')).to_be_visible()
    expect(page.locator('.file-finder')).to_have_count(1, timeout=15000)
    page.keyboard.press('Control+Shift+P')
    expect(page.locator('.markdown-commands[open]')).to_be_visible()
    page.keyboard.press('Escape')
    expect(page.locator('.markdown-commands[open]')).to_have_count(0)

    page.keyboard.press('Control+P')
    expect(page.locator('.file-finder[open]')).to_be_visible()
    page.keyboard.press('Escape')
    expect(page.locator('.file-finder[open]')).to_have_count(0)

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
    expect(page.locator('#shortcuts-list')).to_contain_text('Review & Publish')
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

    page.keyboard.press('Control+Shift+Enter')
    expect(page.locator('#publish-review[open]')).to_be_visible()
    expect(page.locator('#publish-confirm')).to_be_disabled()
    page.keyboard.press('Escape')
    expect(page.locator('#publish-review[open]')).to_have_count(0)
