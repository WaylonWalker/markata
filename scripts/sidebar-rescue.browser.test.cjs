const {test} = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const {chromium} = require(process.env.PLAYWRIGHT_MODULE || 'playwright');

const root = path.join(__dirname, '../pkg/themes/default/static');
const css = ['variables', 'main', 'components']
  .map(n => fs.readFileSync(path.join(root, 'css/' + n + '.css'), 'utf8')).join('\n');
const rescue = fs.readFileSync(path.join(root, 'js/sidebar-rescue.js'), 'utf8');
const viewTransitions = fs.readFileSync(path.join(root, 'js/view-transitions.js'), 'utf8');

function documentHtml() {
  return `<!doctype html><html><head><style>${css}</style></head><body>
    <header class="site-header" style="height:70px">Navigation</header>
    <div class="page-wrapper">
      <aside class="feed-sidebar feed-sidebar--left" aria-label="Feed navigation">
        <button class="sidebar-toggle" type="button" aria-expanded="false" data-sidebar-label="series" data-sidebar-hotkey="b">Series</button>
        <nav class="feed-nav"><div id="feed-nav-collapsible"><ul class="feed-nav-list">
          <li class="feed-nav-item feed-nav-item--active"><a href="/a/">Current post</a></li>
        </ul></div></nav>
      </aside>
      <main class="main-content" style="min-height:2500px"><h1>Readable static article</h1></main>
      <aside class="doc-sidebar doc-sidebar--right" aria-label="Table of Contents">
        <button class="sidebar-toggle" type="button" aria-expanded="false" data-sidebar-label="on this page" data-sidebar-hotkey="Shift+B">Contents</button>
        <nav class="toc"><a class="toc-link" href="#heading">Heading</a></nav>
      </aside>
    </div></body></html>`;
}

async function setupPage(browser, width) {
  const page = await browser.newPage({viewport: {width, height: 823}});
  await page.route('http://sidebar.test/**', route => route.fulfill({
    contentType: 'text/html', body: documentHtml()
  }));
  await page.goto('http://sidebar.test/');
  return page;
}

for (const width of [412, 1440]) {
  test('fallback drawer controls work without view transitions at ' + width + 'px', async () => {
    const browser = await chromium.launch({headless: true, executablePath: process.env.CHROMIUM_EXECUTABLE || undefined, args: ['--no-sandbox']});
    try {
      const page = await setupPage(browser, width);
      await page.addScriptTag({content: rescue});
      const left = page.locator('.feed-sidebar');
      const control = left.locator('.sidebar-toggle');
      assert.equal(await control.getAttribute('aria-expanded'), 'false');
      assert.equal(await left.evaluate(el => el.classList.contains('sidebar--pinned')), false);
      await control.click();
      assert.equal(await left.evaluate(el => el.classList.contains('sidebar--pinned')), true);
      assert.equal(await control.getAttribute('aria-expanded'), 'true');
      if (width <= 1200) {
        await page.keyboard.press('Escape');
        assert.equal(await left.evaluate(el => el.classList.contains('sidebar--pinned')), false);
        assert.equal(await control.getAttribute('aria-expanded'), 'false');
      } else {
        assert.equal(await page.evaluate(() => localStorage.getItem('markata-sidebar:left')), 'open');
        await page.reload();
        await page.addScriptTag({content: rescue});
        assert.equal(await page.locator('.feed-sidebar').evaluate(el => el.classList.contains('sidebar--pinned')), true);
      }
    } finally {
      await browser.close();
    }
  });
}

test('fallback does not double-toggle when the normal controller is healthy', async () => {
  const browser = await chromium.launch({headless: true, executablePath: process.env.CHROMIUM_EXECUTABLE || undefined, args: ['--no-sandbox']});
  try {
    const page = await setupPage(browser, 412);
    await page.addScriptTag({content: viewTransitions});
    await page.addScriptTag({content: rescue});
    const control = page.locator('.feed-sidebar > .sidebar-toggle');
    await control.click();
    assert.equal(await page.locator('.feed-sidebar').evaluate(el => el.classList.contains('sidebar--pinned')), true);
    await control.click();
    assert.equal(await page.locator('.feed-sidebar').evaluate(el => el.classList.contains('sidebar--pinned')), false);
  } finally {
    await browser.close();
  }
});
