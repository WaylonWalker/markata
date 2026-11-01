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

    page.keyboard.press('Control+K')
    page.keyboard.press('Control+S')
    expect(page.locator('#shortcuts-dialog[open]')).to_be_visible()
    expect(page.locator('#shortcuts-list')).to_contain_text('Review & Publish')
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
