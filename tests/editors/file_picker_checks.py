"""Run against a workspace editor: python file_picker_checks.py URL."""
import sys
from playwright.sync_api import sync_playwright, expect

with sync_playwright() as driver:
    browser = driver.chromium.launch(executable_path='/usr/bin/chromium', headless=True)
    page = browser.new_page()
    errors = []
    page.on('pageerror', lambda error: errors.append(str(error)))
    page.goto(sys.argv[1])
    page.wait_for_function("document.querySelector('#active-file').textContent !== 'draft.md'")
    original = page.locator('#active-file').inner_text()
    page.evaluate("document.querySelector('.CodeMirror').CodeMirror.setValue('picker draft preserved')")
    page.keyboard.press('Control+p')
    query = page.locator('#files-query')
    expect(query).to_be_focused()
    query.fill('draft.md')
    expect(page.locator('#files-results [role=option]').first).to_contain_text('draft.md')
    page.keyboard.press('Enter')
    expect(page.locator('#active-file')).to_have_text('draft.md')
    page.keyboard.press('Control+k')
    page.locator('#commands-query').fill('open file')
    page.keyboard.press('Enter')
    expect(query).to_be_focused()
    query.fill(original)
    page.keyboard.press('Enter')
    expect(page.locator('#active-file')).to_have_text(original)
    assert page.evaluate("document.querySelector('.CodeMirror').CodeMirror.getValue()") == 'picker draft preserved'
    page.keyboard.press('Control+p')
    query.fill('zzzzzzzzzzz')
    expect(page.locator('#files-results')).to_be_empty()
    page.keyboard.press('Escape')
    expect(query).not_to_be_visible()
    assert not errors, errors
    browser.close()
