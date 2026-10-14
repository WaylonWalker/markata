// Browser regressions for page-scoped assets and calendar view transitions.
// Run: node --test scripts/navigation-media-calendar.browser.test.cjs
const { test } = require('node:test');
const assert = require('node:assert/strict');
const http = require('node:http');
const fs = require('node:fs');
const path = require('node:path');
const playwright = require(process.env.PLAYWRIGHT_MODULE || 'playwright');

const scripts = {
  '/view.js': fs.readFileSync(path.join(__dirname, '../pkg/themes/default/static/js/view-transitions.js')),
  '/conditional.js': fs.readFileSync(path.join(__dirname, '../pkg/themes/default/static/js/conditional-css.js')),
  '/calendar.js': fs.readFileSync(path.join(__dirname, '../pkg/themes/default/static/js/calendar-feed.js')),
  '/lite.js': Buffer.from(
    "window.__liteLoads=(window.__liteLoads||0)+1;" +
    "if (!customElements.get('lite-youtube')) {" +
    "customElements.define('lite-youtube',class extends HTMLElement {" +
    "connectedCallback(){this.setAttribute('data-upgraded','true');}});}"
  ),
};

function pageHTML(name) {
  const calendar = name.startsWith('calendar-');
  const video = name === 'video';
  const body = video
    ? '<h1>Video</h1><lite-youtube videoid="example"></lite-youtube>'
    : calendar
      ? '<h1>Calendar ' + name + '</h1>' +
        '<div class="feed calendar-feed" data-calendar-feed data-calendar-default="list" ' +
        'data-calendar-url-state data-calendar-src="/' + name + '/calendar.json">' +
        '<nav data-calendar-switch hidden>' +
        '<button data-calendar-mode="list" aria-pressed="true">List</button>' +
        '<button data-calendar-mode="calendar" aria-pressed="false">Calendar</button>' +
        '</nav><div data-calendar-primary><a href="/article/">Example post</a></div>' +
        '<div data-calendar-years hidden aria-label="Publishing calendar"></div>' +
        '<p data-calendar-error hidden role="status">Calendar data could not load</p></div>'
      : '<h1>Plain article</h1>';
  const hints = video
    ? ' data-markata-lite-youtube-css="/lite.css" data-markata-lite-youtube-js="/lite.js"'
    : '';
  return [
    '<!doctype html><html data-loading-mode="auto"><head><meta charset="utf-8">',
    '<title>' + name + '</title>',
    calendar ? '<script src="/calendar.js" defer></script>' : '',
    '</head><body>',
    '<div id="view-transition-page"' + hints + '><main>',
    '<a id="to-plain" href="/plain/">Plain</a> ',
    '<a id="to-video" href="/video/">Video</a> ',
    '<a id="to-calendar-a" href="/calendar-a/">Calendar A</a> ',
    '<a id="to-calendar-b" href="/calendar-b/">Calendar B</a>',
    body, '</main></div>',
    '<script>window.VIEW_TRANSITIONS_CONFIG={enabled:true};</script>',
    '<script src="/view.js" defer></script>',
    '<script src="/conditional.js" defer></script>',
    '</body></html>',
  ].join('');
}

function fixture() {
  const requests = [];
  const server = http.createServer((req, res) => {
    const route = new URL(req.url, 'http://localhost').pathname;
    requests.push(route);
    if (route in scripts) {
      res.setHeader('content-type', 'text/javascript');
      return res.end(scripts[route]);
    }
    if (route === '/lite.css') {
      res.setHeader('content-type', 'text/css');
      return res.end('lite-youtube{display:block;min-height:10px}');
    }
    if (/^\/calendar-[ab]\/calendar\.json$/.test(route)) {
      res.setHeader('content-type', 'application/json');
      // Delay data to expose click-order races.
      return setTimeout(() => res.end(JSON.stringify([
        {date:'2026-10-12', href:'/article/', title:'Example post'}
      ])), 450);
    }
    const match = route.match(/^\/(plain|video|calendar-a|calendar-b)\/$/);
    if (match) {
      res.setHeader('content-type', 'text/html; charset=utf-8');
      return res.end(pageHTML(match[1]));
    }
    res.writeHead(404);
    res.end('Not found');
  });
  return { server, requests };
}

async function withBrowser(run) {
  const { server, requests } = fixture();
  await new Promise(resolve => server.listen(0, '127.0.0.1', resolve));
  const browser = await playwright.chromium.launch({headless:true, args:['--no-sandbox']});
  try {
    const page = await browser.newPage();
    const errors = [];
    page.on('pageerror', error => errors.push(error.message));
    const origin = 'http://127.0.0.1:' + server.address().port;
    await run(page, origin, requests, errors);
    assert.deepEqual(errors, [], 'No uncaught errors in page scripts');
  } finally {
    await browser.close();
    await new Promise(resolve => server.close(resolve));
  }
}

async function clickThrough(page, target, marker) {
  await page.locator('#to-' + target).click();
  await page.waitForFunction(
    ({target, marker}) => location.pathname === '/' + target + '/' &&
      Boolean(document.querySelector(marker)),
    {target, marker}
  );
  // History and the new DOM can update before the transition's finish hook
  // publishes metrics, so wait for completion rather than racing that hook.
  await page.waitForFunction(() => Boolean(window.__lastViewTransitionMetrics), null, { timeout: 6000 });
}

test('Lite YouTube initializes only when needed, including after an SPA navigation', {timeout:60000}, async () => {
  await withBrowser(async (page, origin, requests) => {
    await page.goto(origin + '/plain/');
    await page.waitForLoadState('load');
    assert.equal(requests.filter(x => x === '/lite.js').length, 0);
    assert.equal(requests.filter(x => x === '/lite.css').length, 0);

    await clickThrough(page, 'video', 'lite-youtube');
    await page.waitForFunction(() => document.querySelector('lite-youtube')?.getAttribute('data-upgraded') === 'true');
    assert.equal(await page.evaluate(() => window.__liteLoads), 1);
    assert.equal(requests.filter(x => x === '/lite.css').length, 1);

    await clickThrough(page, 'plain', 'h1');
    await clickThrough(page, 'video', 'lite-youtube');
    await page.waitForFunction(() => document.querySelector('lite-youtube')?.getAttribute('data-upgraded') === 'true');
    assert.equal(await page.evaluate(() => window.__liteLoads), 1, 'Custom element must not be registered twice');
    assert.equal(requests.filter(x => x === '/lite.js').length, 1, 'No redundant script download');

    await page.goto(origin + '/video/');
    await page.waitForFunction(() => document.querySelector('lite-youtube')?.getAttribute('data-upgraded') === 'true');
    assert.equal(requests.filter(x => x === '/lite.js').length, 2, 'Full page load initializes normally');
  });
});

test('Calendar latest user intent wins over delayed JSON', {timeout:60000}, async () => {
  await withBrowser(async (page, origin, requests) => {
    await page.goto(origin + '/calendar-a/');
    await page.locator('[data-calendar-switch]').waitFor({state:'visible'});
    await page.getByRole('button', {name:'Calendar'}).click();
    await page.getByRole('button', {name:'List'}).click();
    await page.waitForResponse(response => response.url().includes('/calendar-a/calendar.json'));
    await page.waitForTimeout(100);
    assert.equal(await page.locator('[data-calendar-years]').isHidden(), true, 'Stale response cannot switch to Calendar');
    assert.equal(await page.locator('[data-calendar-primary]').isVisible(), true);
    assert.equal(await page.getByRole('button', {name:'List'}).getAttribute('aria-pressed'), 'true');
    await page.getByRole('button', {name:'Calendar'}).click();
    await page.locator('[data-calendar-years]').waitFor({state:'visible'});
    assert.equal(await page.getByRole('button', {name:'Calendar'}).getAttribute('aria-pressed'), 'true');
    assert.equal(requests.filter(x => x === '/calendar-a/calendar.json').length, 1, 'Metadata is cached across view switches');
  });
});

test('Calendar controls rehydrate on first and repeated View Transitions', {timeout:60000}, async () => {
  await withBrowser(async (page, origin) => {
    await page.goto(origin + '/plain/');
    await clickThrough(page, 'calendar-a', '[data-calendar-feed]');
    await page.locator('[data-calendar-switch]').waitFor({state:'visible'});
    await page.getByRole('button', {name:'Calendar'}).click();
    await page.locator('[data-calendar-years]').waitFor({state:'visible'});

    await clickThrough(page, 'plain', 'h1');
    await clickThrough(page, 'calendar-b', '[data-calendar-feed]');
    await page.locator('[data-calendar-switch]').waitFor({state:'visible'});
    await page.getByRole('button', {name:'Calendar'}).click();
    await page.locator('[data-calendar-years]').waitFor({state:'visible'});

    await clickThrough(page, 'calendar-a', '[data-calendar-feed]');
    await page.locator('[data-calendar-switch]').waitFor({state:'visible'});
    await page.getByRole('button', {name:'Calendar'}).click();
    await page.locator('[data-calendar-years]').waitFor({state:'visible'});
  });
});
