// Behavioral checks for deferred archive layout and intent-based pagination.
const {test} = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const http = require('node:http');
const {chromium} = require(process.env.PLAYWRIGHT_MODULE || 'playwright');
const root = path.join(__dirname, '../pkg/themes/default/static');

test('archive leaves adjacent pages idle and preserves deferred weeks', async () => {
  const requests = [];
  const server = http.createServer((req, res) => {
    requests.push(req.url);
    if (req.url === '/navigation.js') {
      res.setHeader('content-type', 'text/javascript');
      return res.end(fs.readFileSync(path.join(root, 'js/navigation-shortcuts.js')));
    }
    if (req.url === '/feeds.css') {
      res.setHeader('content-type', 'text/css');
      return res.end(fs.readFileSync(path.join(root, 'css/feeds.css')));
    }
    res.setHeader('content-type', 'text/html');
    res.end(`<!doctype html><html><head><link rel="stylesheet" href="/feeds.css"></head>
      <body><nav><a class="pagination-next" href="/older/">Older</a></nav>
      <main class="feed feed--editorial"><div class="editorial-feed">
      ${Array.from({length: 7}, (_, i) => `<section class="editorial-week" id="week-${i}" style="--editorial-week-items:${i===6 ? 10 : 3}">
        <h2>Week ${i}</h2><div class="editorial-week-entries">
        <div class="editorial-entry" style="height:1800px"><a href="/post-${i}/">Story ${i}</a></div>
        <div class="editorial-entry" style="height:1800px"><a href="/extra-${i}/">Extra ${i}</a></div>
        </div></section>`).join('')}
      </div></main><script>
      window.shortcutsRegistry={register(){}};
      window.prefetchViewTransitionUrl=url=>fetch(url);
      </script><script src="/navigation.js"></script></body></html>`);
  });
  await new Promise(resolve => server.listen(0, '127.0.0.1', resolve));
  const browser = await chromium.launch({headless: true, executablePath: process.env.CHROMIUM_EXECUTABLE || undefined, args: ['--no-sandbox']});
  try {
    const origin = `http://127.0.0.1:${server.address().port}`;
    const page = await browser.newPage({viewport: {width: 412, height: 823}});
    await page.goto(origin);
    await page.waitForFunction(() => typeof window.initNavigationShortcuts === 'function');
    await page.evaluate(() => window.initNavigationShortcuts());
    await page.waitForTimeout(100);
    assert.equal(requests.includes('/older/'), false, 'setup must not fetch another page');
    assert.equal(await page.locator('#week-0').evaluate(e => getComputedStyle(e).contentVisibility), 'visible');
    assert.equal(await page.locator('#week-6').evaluate(e => getComputedStyle(e).contentVisibility), 'auto');
    const reserved = await page.evaluate(() => [5,6].map(i => document.getElementById('week-'+i).getBoundingClientRect().height));
    assert.ok(reserved[1] > reserved[0], 'larger skipped weeks must reserve more room than smaller weeks');
    // Finding skipped content can schedule its layout for a subsequent frame.
    // Wait for a real selection instead of requiring synchronous discovery.
    await page.waitForFunction(() => window.find('Extra 6', false, false, true), null, {timeout: 5000});
    assert.equal(await page.evaluate(() => getSelection().toString()), 'Extra 6');
    // Test focus from a fresh document with skipped weeks, independently of
    // Find's selection and scroll anchoring. Allow one CSS pixel of rounding.
    await page.goto(origin + '/?focus-check');
    await page.locator('#week-6 a').first().focus();
    await page.waitForFunction(() => {
      const link = document.querySelector('#week-6 a');
      const rect = link.getBoundingClientRect();
      const visible = document.activeElement === link && rect.top >= -1 && rect.bottom <= innerHeight + 1;
      window.focusVisibleFrames = visible ? (window.focusVisibleFrames || 0) + 1 : 0;
      return window.focusVisibleFrames >= 3;
    }, null, {timeout: 5000});
    const focused = await page.locator('#week-6 a').first().evaluate(e => {
      const rect = e.getBoundingClientRect();
      return {top: rect.top, bottom: rect.bottom, viewport: innerHeight};
    });
    assert.ok(focused.top >= -1 && focused.bottom <= focused.viewport + 1,
      'keyboard focus must keep stories visible: ' + JSON.stringify(focused));
    await page.emulateMedia({media: 'print'});
    assert.equal(await page.locator('#week-6').evaluate(e => getComputedStyle(e).contentVisibility), 'visible');
    await page.evaluate(() => window.preloadNavigationShortcuts());
    await page.waitForFunction(() => performance.getEntriesByType('resource').some(e => e.name.endsWith('/older/')));
    assert.equal(requests.includes('/older/'), true, 'explicit navigation intent may prefetch');
    const noJS = await browser.newContext({javaScriptEnabled: false});
    const plain = await noJS.newPage();
    await plain.goto(origin);
    assert.equal(await plain.locator('.editorial-week a').count(), 14);
    await plain.locator('#week-6 a').first().scrollIntoViewIfNeeded();
    assert.equal(await plain.locator('#week-6 a').first().isVisible(), true);
    await noJS.close();
  } finally {
    await browser.close();
    await new Promise(resolve => server.close(resolve));
  }
});

for (const width of [320, 412, 1440]) {
  test(`editorial pagination stays put while calendar controls initialize at ${width}px`, async () => {
    const browser = await chromium.launch({headless: true,
      executablePath: process.env.CHROMIUM_EXECUTABLE || undefined, args: ['--no-sandbox']});
    try {
      const css = ['variables', 'main', 'components', 'feeds', 'calendar-feed']
        .map(name => fs.readFileSync(path.join(root, `css/${name}.css`), 'utf8')).join('\n');
      const bootstrap = fs.readFileSync(path.join(root, '../templates/partials/theme-bootstrap.html'), 'utf8');
      const marker = bootstrap.match(/document\.documentElement\.classList\.add\('editorial-calendar-js'\);/);
      assert.ok(marker, 'the real bootstrap must install its pre-paint calendar marker');
      for (const javaScriptEnabled of [false, true]) {
        const context = await browser.newContext({viewport: {width, height: 823}, javaScriptEnabled});
        try {
          const page = await context.newPage();
          await page.setContent(`<!doctype html><html><head><script>${marker[0]}</script><style>${css}</style></head>
            <body><div class="feed feed--editorial calendar-feed" style="width:calc(100vw - 64px);margin:32px">
            <header class="feed-header"><div class="editorial-feed-tools">
            <nav class="calendar-view-switch" data-calendar-switch hidden aria-label="Feed views">
            <button type="button">Posts</button><a href="/simple/">Simple</a><button type="button">Calendar</button></nav>
            <nav class="pagination pagination--top"><span class="pagination-page current">Page 1 of 49</span>
            <a class="pagination-next" href="/page/2/">Older →</a></nav>
            </div></header></div></body></html>`);
          const state = await page.locator('.calendar-view-switch').evaluate(el => {
            el.querySelector('button').focus();
            return {display: getComputedStyle(el).display, visibility: getComputedStyle(el).visibility,
              focused: el.contains(document.activeElement)};
          });
          if (!javaScriptEnabled) {
            assert.equal(state.display, 'none', 'no-JavaScript controls retain their original hidden layout');
            continue;
          }
          assert.equal(state.display, 'flex', 'inline flex containers are blockified as toolbar flex items');
          assert.equal(state.visibility, 'hidden');
          assert.equal(state.focused, false, 'pending controls must stay non-interactive');
          const geometry = () => page.locator('.pagination--top').evaluate(el => {
            const r = el.getBoundingClientRect();
            return [r.x, r.y, r.width, r.height];
          });
          const before = await geometry();
          // The deferred calendar controller removes hidden after wiring controls.
          await page.locator('.calendar-view-switch').evaluate(el => { el.hidden = false; });
          assert.deepEqual(await geometry(), before, 'initialization must not move pagination');
          assert.equal(await page.locator('.calendar-view-switch').evaluate(el => getComputedStyle(el).visibility), 'visible');
          await page.evaluate(() => document.documentElement.classList.remove('editorial-calendar-js'));
          assert.deepEqual(await geometry(), before, 'the marker must not change the final toolbar layout');
        } finally {
          await context.close();
        }
      }
    } finally {
      await browser.close();
    }
  });
}
