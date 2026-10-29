"""Behavior checks shared by the three real editor browser integrations."""
from playwright.sync_api import expect


def check_commands(page, preview):
    dialog = page.get_by_role('dialog', name='Insert Markdown')
    query = page.get_by_role('combobox', name='Find a command')
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

    # Click execution, default placeholder and actual chart rendering.
    page.locator('#commands').click()
    query.fill('chartjs')
    dialog.get_by_role('option').click()
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
