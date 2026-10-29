"""Shared browser checks for the Python and Go example site trees."""
from playwright.sync_api import expect


def check_site_tree(page, url):
    errors, rendered_sources = [], []
    page.on('pageerror', lambda error: errors.append(str(error)))
    page.on('request', lambda request: rendered_sources.append(request.post_data_json.get('markdown', ''))
            if request.url.endswith('/render') and request.method == 'POST' else None)
    # Pierre uses a closed shadow root. Open it only in the browser test.
    page.add_init_script('''const attach = Element.prototype.attachShadow;
        Element.prototype.attachShadow = function(options) {
          return attach.call(this, {...options, mode:'open'});
        };''')
    page.goto(url)
    preview = page.frame_locator('#preview')
    expect(preview.locator('h1')).to_have_text('Welcome to the field notes')
    expect(preview.locator('body')).not_to_contain_text('published: true')
    page.evaluate("document.querySelector('.CodeMirror').CodeMirror.setValue('# Saved welcome')")
    expect(preview.locator('h1')).to_have_text('Saved welcome')
    page.locator('#site-files').get_by_role('treeitem', name='growing.md', exact=True).click()
    expect(preview.locator('h1')).to_have_text('A growing garden')
    page.locator('#site-view').select_option('feeds')
    page.locator('#site-feeds').get_by_role('treeitem', name='welcome.md', exact=True).click()
    expect(preview.locator('h1')).to_have_text('Saved welcome')
    page.locator('#site-feeds').get_by_role('treeitem', name='growing.md', exact=True).first.click()
    expect(preview.locator('h1')).to_have_text('A growing garden')
    page.locator('#site-view').select_option('files')
    page.locator('#site-files').get_by_role('treeitem', name='welcome.md', exact=True).click()
    expect(preview.locator('h1')).to_have_text('Saved welcome')
    page.locator('#site-view').select_option('feeds')
    page.locator('#site-feeds').get_by_role('treeitem', name='growing.md', exact=True).first.click()
    expect(preview.locator('h1')).to_have_text('A growing garden')
    page.reload()
    expect(preview.locator('h1')).to_have_text('A growing garden')
    page.locator('#site-files').get_by_role('treeitem', name='markata-go.toml', exact=True).click()
    assert page.evaluate("document.querySelector('.CodeMirror').CodeMirror.getOption('readOnly')")
    expect(page.locator('#commands')).to_be_disabled()
    original = page.evaluate("document.querySelector('.CodeMirror').CodeMirror.getValue()")
    page.evaluate("document.querySelector('.CodeMirror').CodeMirror.replaceSelection('attempted edit')")
    assert page.evaluate("document.querySelector('.CodeMirror').CodeMirror.getValue()") == original
    expect(preview.locator('pre')).to_contain_text('[markata-go]')
    page.locator('#import').set_input_files({'name':'import.md','mimeType':'text/markdown','buffer':b'# Imported file'})
    expect(preview.locator('h1')).to_have_text('Imported file')
    expect(page.locator('#active-file')).to_have_text('draft.md')
    page.reload()
    expect(preview.locator('h1')).to_have_text('Imported file')
    page.evaluate("() => {Storage.prototype.setItem = () => {throw Error('full')}}")
    page.evaluate("document.querySelector('.CodeMirror').CodeMirror.setValue('# Unsaved edit')")
    expect(page.locator('#saved')).to_contain_text('Draft not saved')
    page.locator('#site-files').get_by_role('treeitem', name='about.md',exact=True).click()
    expect(preview.locator('h1')).to_have_text('About this site')
    page.locator('#site-files').get_by_role('treeitem', name='draft.md',exact=True).click()
    expect(preview.locator('h1')).to_have_text('Unsaved edit')
    for width in [390, 850, 1280]:
        page.set_viewport_size({'width':width,'height':850})
        assert page.evaluate('document.documentElement.scrollWidth <= innerWidth'), width
    # Rapid file changes must leave both panes on the last selection.
    page.locator('#site-view').select_option('files')
    for name in ['growing.md', 'about.md', 'welcome.md'] * 5:
        page.locator('#site-files').get_by_role('treeitem', name=name, exact=True).click()
        page.wait_for_timeout(300)
    expect(page.locator('#active-file')).to_have_text('posts/welcome.md')
    expect(preview.locator('h1')).to_have_text('Saved welcome', timeout=10000)
    assert rendered_sources[-1].startswith('# Saved welcome'), rendered_sources[-1][:100]
    assert not errors, errors


if __name__ == '__main__':
    import argparse
    import sys
    from pathlib import Path
    from playwright.sync_api import sync_playwright

    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--url', help='Running Go or Python workbench URL')
    parser.add_argument('--chromium', default='/usr/bin/chromium')
    args = parser.parse_args()
    with sync_playwright() as driver:
        browser = driver.chromium.launch(executable_path=args.chromium, headless=True)
        with browser.new_context() as context:
            page = context.new_page()
            if args.url:
                check_site_tree(page, args.url)
            else:
                sys.path.insert(0, str(Path(__file__).resolve().parents[2] / 'python'))
                from editor_test_support import load_editor, serve_editor
                with serve_editor(load_editor().app) as url:
                    check_site_tree(page, url)
        browser.close()
    print('Site tree browser checks passed')
