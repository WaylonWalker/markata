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
    # Metadata cannot finish until the selected source and preview are usable.
    metadata_checked = []
    def hold_metadata(route):
        expect(page.frame_locator('#preview').locator('h1')).to_have_text('About this site')
        expect(page.locator('#site-view option[value=feeds]')).to_be_disabled()
        metadata_checked.append(True)
        route.fulfill(status=200, content_type='application/json', body='{"files":[]}')
    page.route('**/workspace?metadata=true', hold_metadata)
    with page.expect_response('**/workspace?metadata=true'):
        page.goto(url)
    assert metadata_checked
    page.unroute('**/workspace?metadata=true', hold_metadata)
    # Simulate an empty placeholder persisted by the first workspace build.
    workspace = page.request.get(url.rstrip('/') + '/workspace').json()
    key = workspace.get('rootUri') or workspace.get('root')
    page.evaluate("key => {localStorage.removeItem('markata-workspace-v1:' + key + ':buffer:file:pages/about.md'); localStorage.setItem('markata-workspace-v1:' + key + ':pages/about.md', '');}", key)
    page.reload()
    preview = page.frame_locator('#preview')

    def open_feed_post(slug):
        group = page.locator('#site-feeds').get_by_role('treeitem', name='Field notes (2)', exact=True)
        if group.get_attribute('aria-expanded') != 'true':
            group.click()
        page.locator('#site-feeds').get_by_role('treeitem', name=slug, exact=True).click()

    expect(preview.locator('h1')).to_have_text('About this site')
    page.locator('#site-files').get_by_role('treeitem', name='welcome.md', exact=True).click()
    expect(preview.locator('h1')).to_have_text('Welcome to the field notes')
    expect(preview.locator('body')).not_to_contain_text('published: true')
    page.evaluate("document.querySelector('.CodeMirror').CodeMirror.setValue('# Saved welcome')")
    expect(preview.locator('h1')).to_have_text('Saved welcome')
    page.locator('#site-files').get_by_role('treeitem', name='growing.md', exact=True).click()
    expect(preview.locator('h1')).to_have_text('A growing garden')
    feeds_available = bool(page.request.get(url.rstrip('/') + '/workspace?metadata=true').json().get('feeds'))
    if feeds_available:
        expect(page.locator('#site-view option[value=feeds]')).to_be_enabled(timeout=30000)
        page.locator('#site-view').select_option('feeds')
        open_feed_post('welcome')
        expect(preview.locator('h1')).to_have_text('Saved welcome')
        open_feed_post('growing')
        expect(preview.locator('h1')).to_have_text('A growing garden')
    else:
        expect(page.locator('#site-view option[value=feeds]')).to_be_disabled()
    page.locator('#site-view').select_option('files')
    page.locator('#site-files').get_by_role('treeitem', name='welcome.md', exact=True).first.click()
    expect(preview.locator('h1')).to_have_text('Saved welcome')
    if feeds_available:
        expect(page.locator('#site-view option[value=feeds]')).to_be_enabled(timeout=30000)
        page.locator('#site-view').select_option('feeds')
        open_feed_post('growing')
        expect(preview.locator('h1')).to_have_text('A growing garden')
    else:
        page.locator('#site-files').get_by_role('treeitem', name='growing.md', exact=True).click()
        expect(preview.locator('h1')).to_have_text('A growing garden')
    page.reload()
    expect(preview.locator('h1')).to_have_text('A growing garden')
    page.locator('#import').set_input_files({'name':'import.md','mimeType':'text/markdown','buffer':b'# Imported file'})
    expect(preview.locator('h1')).to_have_text('Imported file')
    expect(page.locator('#active-file')).to_have_text('draft.md')
    page.reload()
    expect(preview.locator('h1')).to_have_text('Imported file')
    page.evaluate("() => {Storage.prototype.setItem = () => {throw Error('full')}}")
    page.locator('.CodeMirror').click()
    page.keyboard.press('Control+A')
    page.keyboard.type('# Unsaved edit')
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
    # Reload restores Storage after the deliberate failure injection above.
    page.reload()
    expect(page.frame_locator('#preview').locator('h1')).to_be_visible()
    page.locator('#site-files').get_by_role('treeitem', name='welcome.md', exact=True).click()
    expect(page.locator('#active-file')).to_have_text('posts/welcome.md')
    # An intentionally empty draft in the new format must stay empty on reload.
    page.locator('.CodeMirror').click()
    page.keyboard.press('Control+A')
    page.keyboard.press('Backspace')
    page.reload()
    expect(page.locator('#active-file')).to_have_text('posts/welcome.md')
    assert page.evaluate("document.querySelector('.CodeMirror').CodeMirror.getValue()") == ''
    page.route('**/workspace?metadata=true', lambda route: route.abort())
    page.locator('.CodeMirror').click()
    page.keyboard.type('# Metadata offline')
    page.reload()
    expect(preview.locator('h1')).to_have_text('Metadata offline')
    expect(page.locator('#site-view option[value=feeds]')).to_have_text('Feeds unavailable')
    expect(page.locator('#site-view option[value=tags]')).to_have_text('Tags unavailable')
    assert not errors, errors


if __name__ == '__main__':
    import argparse
    import shutil
    import sys
    import tempfile
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
                with tempfile.TemporaryDirectory(prefix='markata-site-tree-') as directory:
                    site = Path(directory) / 'site'
                    shutil.copytree(Path(__file__).resolve().parents[2] / 'examples/editor-site', site)
                    module = load_editor()
                    module.LSP_WORKSPACE = site
                    module.presentation = module.appearance(asset_base='/appearance', workspace=str(site))
                    with serve_editor(module.app) as url:
                        check_site_tree(page, url)
        browser.close()
    print('Site tree browser checks passed')
