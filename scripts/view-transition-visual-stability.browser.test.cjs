// Regression: new page CSS and sidebar chrome must be ready in the
// destination snapshot, without animations or flashes across navigation.
// Run: node --test scripts/view-transition-visual-stability.browser.test.cjs
const {test} = require('node:test');
const assert = require('node:assert/strict');
const http = require('node:http');
const fs = require('node:fs');
const path = require('node:path');
const {chromium} = require(process.env.PLAYWRIGHT_MODULE || 'playwright');

const themeRoot = path.join(__dirname, '../pkg/themes/default/static');
const runtime = fs.readFileSync(path.join(themeRoot, 'js/view-transitions.js'));
const css = fs.readFileSync(path.join(themeRoot, 'css/components.css'), 'utf8');

function html(name) {
  const next = name === 'start' ? 'target' : 'start';
  return [
    '<!doctype html><html><head><meta charset="utf-8">',
    '<title>' + name + '</title>',
    '<link rel="stylesheet" href="/' + name + '.css">',
    '<style>' + css + '</style>',
    '</head><body><div id="view-transition-page" class="page-wrapper">',
    '<aside class="feed-sidebar"><nav class="feed-nav">Sidebar ' + name + '</nav></aside>',
    '<main class="main-content"><div id="surface">Page ' + name + '</div>',
    '<a id="next" href="/' + next + '/">Next</a></main>',
    '</div><script>window.VIEW_TRANSITIONS_CONFIG={enabled:true};</script>',
    '<script src="/view.js" defer></script></body></html>',
  ].join('');
}

test('delayed destination CSS is ready before the new snapshot, without hidden drawers or animation', {timeout: 30000}, async () => {
  const requests = [];
  const server = http.createServer((req, res) => {
    const route = new URL(req.url, 'http://localhost').pathname;
    requests.push(route);
    if (route === '/view.js') {
      res.setHeader('content-type', 'text/javascript');
      return res.end(runtime);
    }
    if (route === '/start.css' || route === '/target.css') {
      res.setHeader('content-type', 'text/css');
      const background = route === '/start.css' ? 'rgb(180, 0, 0)' : 'rgb(0, 120, 0)';
      const body = '#surface{background:' + background + ';min-height:120px;}';
      // Deliberately expose the old FOUC: the target stylesheet is NOT
      // available by the time the destination HTML request completes.
      if (route === '/target.css') return setTimeout(() => res.end(body), 500);
      return res.end(body);
    }
    const match = route.match(/^\/(start|target)\/$/);
    if (match) {
      res.setHeader('content-type', 'text/html; charset=utf-8');
      return res.end(html(match[1]));
    }
    res.writeHead(404);
    res.end('not found');
  });
  await new Promise(resolve => server.listen(0, '127.0.0.1', resolve));
  const browser = await chromium.launch({
    headless: true,
    executablePath: process.env.CHROMIUM_EXECUTABLE || undefined,
    args: ['--no-sandbox'],
  });
  try {
    const page = await browser.newPage();
    const errors = [];
    page.on('pageerror', error => errors.push(error.message));
    const origin = 'http://127.0.0.1:' + server.address().port;
    await page.goto(origin + '/start/');
    assert.equal(await page.locator('#surface').evaluate(el => getComputedStyle(el).backgroundColor), 'rgb(180, 0, 0)');

    // Capture exactly what the browser will use as its new snapshot.
    await page.evaluate(() => {
      const original = document.startViewTransition.bind(document);
      document.startViewTransition = update => {
        const transition = original(async () => {
          const result = await update();
          window.__destinationLinks = Array.from(document.querySelectorAll('link[rel="stylesheet"]')).map(link => ({
            href: link.href, media: link.media, loaded: !!link.sheet,
            pending: link.getAttribute('data-markata-staged-stylesheet'),
          }));
          window.__destinationSnapshot = {
            pathname: location.pathname,
            hiddenDrawer: !!document.querySelector('[data-sidebar-hydrating]'),
            incomingStylesheetActive: Array.from(document.querySelectorAll('link[rel="stylesheet"]')).some(link =>
              link.href.endsWith('/target.css') && !link.hasAttribute('data-markata-staged-stylesheet') &&
              (!link.media || link.media === 'all') && !!link.sheet
            ),
          };
          return result;
        });
        // The browser has completed its layout and created the destination
        // snapshot when ready resolves. Reading computed style in the update
        // callback itself is premature during a stylesheet media switch.
        transition.ready.then(() => {
          window.__readyFrameBackground = getComputedStyle(document.querySelector('#surface')).backgroundColor;
        });
        return transition;
      };
    });

    const requested = page.waitForRequest(req => req.url().endsWith('/target.css'));
    await page.locator('#next').click();
    await requested;
    await page.waitForTimeout(100);
    assert.equal(new URL(page.url()).pathname, '/start/', 'the previous page stays visible while destination styles load');
    assert.equal(await page.locator('#surface').evaluate(el => getComputedStyle(el).backgroundColor), 'rgb(180, 0, 0)', 'staged CSS cannot restyle the outgoing page');

    await page.waitForFunction(() => location.pathname === '/target/' && !!window.__destinationSnapshot);
    const captured = await page.evaluate(() => window.__destinationSnapshot);
    assert.deepEqual(captured, {
      pathname: '/target/',
      hiddenDrawer: false,
      incomingStylesheetActive: true,
    }, 'the incoming stylesheet is active and sidebar visible before snapshot capture: ' +
      JSON.stringify(await page.evaluate(() => window.__destinationLinks)));
    await page.waitForFunction(() => !!window.__readyFrameBackground);
    assert.equal(await page.evaluate(() => window.__readyFrameBackground),
      'rgb(0, 120, 0)', 'the first rendered destination frame has correct CSS');
    await page.waitForFunction(() => !!window.__lastViewTransitionMetrics);
    assert.equal(await page.locator('[data-sidebar-hydrating]').count(), 0);
    assert.equal(requests.filter(r => r === '/target.css').length, 1);
    assert.deepEqual(errors, []);

    // No root, main-content, or shared-element screenshot animations.
    assert.match(css, /::view-transition-group\(\*\),[\s\S]*?::view-transition-new\(\*\)\s*\{\s*animation:\s*none !important;/);
    await page.goBack();
    await page.waitForFunction(() => location.pathname === '/start/' && document.querySelector('#surface')?.textContent === 'Page start');
    assert.equal(await page.locator('#surface').evaluate(el => getComputedStyle(el).backgroundColor), 'rgb(180, 0, 0)');
    assert.deepEqual(errors, []);
  } finally {
    await browser.close();
    await new Promise(resolve => server.close(resolve));
  }
});
