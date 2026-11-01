"""Behavior checks shared by the three real editor browser integrations."""
from playwright.sync_api import expect


def check_commands(page, preview):
    dialog = page.get_by_role('dialog', name='Command Palette')
    query = page.get_by_role('combobox', name='Find a command or action')
    value = lambda: page.evaluate("document.querySelector('.CodeMirror').CodeMirror.getValue()")

    def select_text(text, start=0, end=None):
        page.evaluate("""({text,start,end}) => {
          const cm = document.querySelector('.CodeMirror').CodeMirror;
          cm.setValue(text); cm.clearHistory();
          cm.setSelection(cm.posFromIndex(start),cm.posFromIndex(end ?? text.length));
          cm.focus();
        }""", {'text': text, 'start': start, 'end': end})

    # Each keymap must permit opening the palette. Cancellation preserves selection.
    for keymap in ['default', 'vim', 'emacs', 'sublime']:
        page.locator('#keymap').select_option(keymap)
        select_text('Keep this text', 5, 9)
        page.keyboard.press('Control+Shift+P')
        expect(dialog).to_be_visible()
        expect(query).to_be_focused()
        query.fill('not-a-real-command')
        expect(dialog.get_by_role('status')).to_contain_text('No matching commands')
        page.keyboard.press('Enter')
        expect(dialog).to_be_visible()
        page.keyboard.press('Escape')
        expect(dialog).not_to_be_visible()
        assert value() == 'Keep this text'
        assert page.evaluate("document.querySelector('.CodeMirror').CodeMirror.getSelection()") == 'this'
        assert page.evaluate("document.querySelector('.CodeMirror').CodeMirror.hasFocus()")

    page.locator('#keymap').select_option('default')
    page.keyboard.press('Control+Shift+P')
    expect(dialog).to_be_visible()
    page.keyboard.press('Escape')
    expect(dialog).not_to_be_visible()
    if not page.evaluate('Boolean(window.MarkataShortcuts)'):
        # Standalone Python and WASM editors retain their Ctrl/Cmd+K shortcut.
        page.keyboard.press('Meta+K')
        expect(dialog).to_be_visible()
        page.keyboard.press('Escape')
        expect(dialog).not_to_be_visible()
    select_text('Important')
    page.locator('#commands').click()
    query.fill('HIGHLIGHT text')
    expect(dialog.get_by_role('option')).to_have_count(1)
    page.keyboard.press('Enter')
    assert value() == '==Important=='
    expect(preview.locator('mark')).to_have_text('Important')
    page.evaluate("document.querySelector('.CodeMirror').CodeMirror.undo()")
    assert value() == 'Important'
    page.evaluate("document.querySelector('.CodeMirror').CodeMirror.redo()")
    assert value() == '==Important=='

    # Indentation and blank-line boundaries retain both surrounding paragraphs.
    select_text('Before\nfirst\nsecond\nAfter', 7, 19)
    page.locator('#commands').click()
    query.fill('note box')
    page.keyboard.press('Enter')
    assert value() == 'Before\n\n!!! note "Note"\n    first\n    second\n\nAfter'
    expect(preview.locator('.admonition.note')).to_contain_text('second')
    page.evaluate("document.querySelector('.CodeMirror').CodeMirror.undo()")
    assert value() == 'Before\nfirst\nsecond\nAfter'

    # Keyboard navigation updates the accessible active descendant.
    select_text('')
    page.locator('#commands').click()
    first = query.get_attribute('aria-activedescendant')
    page.keyboard.press('ArrowDown')
    assert query.get_attribute('aria-activedescendant') != first
    page.keyboard.press('ArrowUp')
    assert query.get_attribute('aria-activedescendant') == first
    page.keyboard.press('ArrowUp')
    expect(dialog.get_by_role('option').last).to_have_attribute('aria-selected', 'true')
    # Native modal tab traversal stays in the dialog.
    page.keyboard.press('Tab')
    assert page.evaluate("document.activeElement.closest('dialog') !== null")
    dialog.get_by_role('button', name='Close commands').click()
    assert value() == ''

    # Custom element families insert runnable standalone syntax in every runtime.
    for search, command_id, selector in [
        ('road steps', 'road-steps', '.road.steps[data-road-ready]'),
        ('Article card', 'md-card', '.md-card'),
        ('Icon feature row', 'md-feature', '.md-feature .icon svg'),
        ('Measurement meter', 'md-meter', 'meter[value="45"]'),
        ('Markdown source disclosure', 'md-source', 'details.md-source'),
        ('Interactive checklist', 'md-checklist', '.md-checklist input:not([disabled])'),
        ('Lucide: zodiac-aries', 'lucide-zodiac-aries', '[data-icon="lucide/zodiac-aries"] svg'),
        ('Mermaid: sequence', 'mermaid-sequence', 'pre.mermaid svg'),
    ]:
        select_text('Replace this with an icon' if command_id.startswith('lucide-') else '')
        page.locator('#commands').click()
        query.fill(search)
        dialog.locator('#markdown-command-' + command_id).click()
        expect(preview.locator(selector).first).to_be_visible(timeout=30000)

    # Icon names beyond the initial result window remain searchable, while the
    # unfiltered modal stays bounded. Site requirements are visible before insert.
    select_text('')
    page.locator('#commands').click()
    expect(dialog.get_by_role('status')).to_contain_text('Showing 60 of')
    expect(dialog.get_by_role('option')).to_have_count(60)
    query.fill('Web Awesome:')
    expect(dialog.get_by_role('option')).to_have_count(11)
    expect(dialog.locator('#markdown-command-wa-details')).to_contain_text('Requires the webawesome plugin and runtime')
    dialog.locator('#markdown-command-wa-details').click()
    assert value().startswith('::: wa-details {summary="Read more"}')
    page.evaluate("document.querySelector('.CodeMirror').CodeMirror.undo()")
    assert value() == ''

    # Click execution, default placeholder and actual chart rendering.
    page.locator('#commands').click()
    query.fill('chartjs')
    dialog.locator('#markdown-command-chart').click()
    assert value().startswith('```chartjs\n')
    expect(preview.locator('canvas')).to_be_visible(timeout=30000)
    expect(page.locator('#saved')).to_contain_text('Draft saved')
    page.reload()
    expect(preview.locator('canvas')).to_be_visible(timeout=60000)
    assert value().startswith('```chartjs\n')

    # Small-screen access and backdrop cancellation.
    page.set_viewport_size({'width': 390, 'height': 844})
    page.locator('#commands').click()
    expect(dialog).to_be_visible()
    box = dialog.bounding_box()
    assert box['x'] >= 0 and box['x'] + box['width'] <= 390
    page.mouse.click(2, 2)
    expect(dialog).not_to_be_visible()
    page.set_viewport_size({'width':1280, 'height':720})
    # Restore the expected initial fixture for the existing parity checks.
    select_text('# A place to think')
    expect(preview.locator('h1')).to_have_text('A place to think')
