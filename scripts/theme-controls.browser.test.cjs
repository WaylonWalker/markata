// Real browser regressions for the lazy theme controls.
// Run: node --test scripts/theme-controls.browser.test.cjs
const { test } = require('node:test');
const assert = require('node:assert/strict');
const http = require('node:http');
const fs = require('node:fs');
const path = require('node:path');
const playwright = require(process.env.PLAYWRIGHT_MODULE || 'playwright');

const staticRoot = path.join(__dirname, '../pkg/themes/default');
const picker = fs.readFileSync(path.join(staticRoot, 'static/js/palette-switcher.js'));
const template = fs.readFileSync(path.join(staticRoot, 'templates/partials/theme-controls-loader.html'), 'utf8');
const loader = template.replace(/\{#[\s\S]*?#\}/g, '')
  .replace("{{ 'js/palette-switcher.js' | theme_asset_hashed }}", '/js/palette-switcher.hash.js');

const manifest = [
  { name: 'ayu-dark', variant: 'dark', displayName: 'Ayu Dark', counterpart: 'ayu-light' },
  { name: 'ayu-light', variant: 'light', displayName: 'Ayu Light', counterpart: 'ayu-dark' },
  { name: 'nord-dark', variant: 'dark', displayName: 'Nord Dark', counterpart: 'nord-light' },
  { name: 'nord-light', variant: 'light', displayName: 'Nord Light', counterpart: 'nord-dark' }
];
const css = [
  ':root {',
  '  --palette-switcher-enabled: 1;',
  '  --palette-light: "ayu-light";',
  '  --palette-dark: "ayu-dark";',
  "  --palette-manifest: '" + JSON.stringify(manifest) + "';",
  '}',
  '[data-palette="ayu-dark"] { --color-background: #111827; }',
  '[data-palette="ayu-light"] { --color-background: #fafafa; }',
  '[data-palette="nord-dark"] { --color-background: #161c2c; }',
  '[data-palette="nord-light"] { --color-background: #eeeeee; }',
  'html { background: var(--color-background); }',
].join('\n');

function controls() {
  return `
    <button class="palette-mode-toggle" type="button" aria-label="Toggle dark/light mode">Mode</button>
    <div class="theme-picker" data-theme-picker>
      <button class="theme-picker-toggle" type="button" aria-expanded="false">Change theme</button>
      <div class="theme-picker-panel" hidden>
        <div role="tablist">
          <button data-picker-tab="colors" role="tab">Colors</button>
          <button data-picker-tab="style" role="tab">Style</button>
          <button data-picker-tab="font" role="tab">Font</button>
        </div>
        <button data-picker-step="prev">Previous</button>
        <button data-picker-step="next">Next</button>
        <span data-picker-current></span>
        <div data-picker-pane="colors">
          <input class="theme-picker-search" type="search">
          <button data-picker-mode="light">Light</button>
          <button data-picker-mode="dark">Dark</button>
          <button data-picker-shuffle>Random</button>
          <button data-picker-schedule hidden>Schedule</button>
          <div class="theme-picker-grid"></div>
          <p class="theme-picker-empty" hidden>No themes match.</p>
          <div data-picker-schedule-view hidden></div>
        </div>
        <div data-picker-pane="style" hidden><div data-picker-style-grid></div></div>
        <div data-picker-pane="font" hidden><div data-picker-font-grid></div></div>
        <button data-picker-reset>Reset</button>
      </div>
    </div>`;
}

function pageHTML(options = {}) {
  return `<!doctype html>
    <html data-theme="dark" data-palette="ayu-dark" data-aesthetic="minimal" data-fontpack="default">
    <head><meta charset="UTF-8">
      <script>
        (function () {
          window.__markataThemeFallbackMode = 'dark';
          window.__markataThemeDefaults = { light: 'ayu-light', dark: 'ayu-dark' };
          var saved = localStorage.getItem('color-mode');
          if (saved === 'light' || saved === 'dark') {
            document.documentElement.dataset.theme = saved;
            document.documentElement.dataset.palette = window.__markataThemeDefaults[saved];
          }
        })();
      </script>
      ${options.lazyCSS ? `<link rel="stylesheet" href="/palette-core.css"><script>
        window.__markataPaletteStylesReady = false;
        window.__markataLoadPaletteStyles = function () {
          if (window.palettePromise) return window.palettePromise;
          window.palettePromise = new Promise(function (resolve, reject) {
            var link = document.createElement('link');
            link.rel = 'stylesheet'; link.href = '/palette.css';
            link.onload = function () {window.__markataPaletteStylesReady = true; resolve();};
            link.onerror = reject; document.head.appendChild(link);
          });
          return window.palettePromise;
        };
      </script>` : `<link rel="stylesheet" href="/palette.hash.css" onerror="this.onerror=null;this.href='/palette.css'">`}
    </head><body>
      <main><div id="view-transition-page">${controls()}</div></main>
      ${loader}
    </body></html>`;
}

async function browserCase(options, fn) {
  const requests = [];
  const server = http.createServer((req, res) => {
    const pathname = new URL(req.url, 'http://localhost').pathname;
    requests.push(pathname);
    if (pathname === '/palette-core.css') {
      res.setHeader('content-type', 'text/css');
      return res.end(css.replace(JSON.stringify(manifest), JSON.stringify(manifest.slice(0, 2))));
    }
    if (pathname === '/palette.css' || (pathname === '/palette.hash.css' && !options.missingCSS)) {
      res.setHeader('content-type', 'text/css');
      return res.end(css);
    }
    if (pathname === '/js/palette-switcher.js' ||
        (pathname === '/js/palette-switcher.hash.js' && !options.missingHash)) {
      res.setHeader('content-type', 'text/javascript');
      return res.end(picker);
    }
    if (pathname === '/') {
      res.setHeader('content-type', 'text/html; charset=utf-8');
      return res.end(pageHTML(options));
    }
    res.writeHead(404);
    return res.end('Not found');
  });
  await new Promise(resolve => server.listen(0, '127.0.0.1', resolve));
  const browser = await playwright.chromium.launch({ headless: true, args: ['--no-sandbox'] });
  try {
    const page = await browser.newPage();
    const errors = [];
    page.on('pageerror', error => errors.push(error.message));
    await fn(page, 'http://127.0.0.1:' + server.address().port, requests);
    assert.deepEqual(errors, [], 'No unhandled JavaScript errors');
  } finally {
    await browser.close();
    await new Promise(resolve => server.close(resolve));
  }
}

function pickerRequests(requests) {
  return requests.filter(route => route.startsWith('/js/palette-switcher'));
}

test('mode switching, persistence and shortcut stay eager; picker stays lazy', {timeout: 60000}, async () => {
  await browserCase({}, async (page, origin, requests) => {
    await page.goto(origin + '/');
    assert.deepEqual(pickerRequests(requests), [], 'No picker download on initial load');
    assert.equal(await page.locator('html').getAttribute('data-theme'), 'dark');

    await page.locator('.palette-mode-toggle').click();
    assert.equal(await page.locator('html').getAttribute('data-theme'), 'light');
    assert.equal(await page.locator('html').getAttribute('data-palette'), 'ayu-light');
    assert.equal(await page.evaluate(() => localStorage.getItem('color-mode')), 'light');
    assert.match(await page.locator('.palette-mode-toggle').getAttribute('aria-label'), /Current: Light/);
    assert.deepEqual(pickerRequests(requests), [], 'Mode click should not fetch picker');

    await page.keyboard.press('Backslash');
    assert.equal(await page.locator('html').getAttribute('data-theme'), 'dark');
    assert.deepEqual(pickerRequests(requests), [], 'Keyboard mode toggle should not fetch picker');
    await page.reload();
    assert.equal(await page.locator('html').getAttribute('data-theme'), 'dark');
    assert.deepEqual(pickerRequests(requests), [], 'Restoring mode should not fetch picker');

    await page.keyboard.press('t');
    await page.locator('.theme-picker-panel').waitFor({ state: 'visible' });
    assert.equal(await page.locator('.palette-card').count(), 2, 'Only active mode cards rendered');
    assert.deepEqual(pickerRequests(requests), ['/js/palette-switcher.hash.js']);
    await page.locator('[data-picker-mode="light"]').click();
    assert.equal(await page.locator('html').getAttribute('data-theme'), 'light');
    await page.locator('.palette-card[data-name="nord-light"]').click();
    assert.equal(await page.locator('html').getAttribute('data-palette'), 'nord-light');
    assert.equal(await page.evaluate(() => localStorage.getItem('theme-palette-light')), 'nord-light');
    await page.reload();
    // The minimal fixture bootstrap uses the default palette; production's
    // theme-bootstrap.html restores the visitor's stored palette as well.
    assert.equal(await page.locator('html').getAttribute('data-theme'), 'light');
  });
});

test('picker hash 404 falls back to unhashed JS; controls survive a DOM swap', {timeout: 60000}, async () => {
  await browserCase({ missingHash: true }, async (page, origin, requests) => {
    await page.goto(origin + '/');
    await page.locator('.theme-picker-toggle').click();
    await page.locator('.theme-picker-panel').waitFor({ state: 'visible' });
    assert.deepEqual(pickerRequests(requests), [
      '/js/palette-switcher.hash.js', '/js/palette-switcher.js'
    ]);
    await page.evaluate((html) => {
      document.querySelector('#view-transition-page').innerHTML = html;
      window.dispatchEvent(new Event('view-transition-complete'));
    }, controls());
    await page.locator('.palette-mode-toggle').click();
    assert.equal(await page.locator('html').getAttribute('data-theme'), 'light');
    await page.locator('.theme-picker-toggle').click();
    await page.locator('.theme-picker-panel').waitFor({ state: 'visible' });
    assert.equal(pickerRequests(requests).length, 2, 'No duplicate picker downloads on navigation');
  });
});

test('missing fingerprinted palette CSS recovers without fetching picker JS', {timeout: 60000}, async () => {
  await browserCase({ missingCSS: true }, async (page, origin, requests) => {
    await page.goto(origin + '/');
    assert.ok(requests.includes('/palette.hash.css'));
    assert.ok(requests.includes('/palette.css'), 'Unhashed stylesheet fallback must load');
    assert.equal(await page.locator('html').evaluate(el => getComputedStyle(el).backgroundColor), 'rgb(17, 24, 39)');
    await page.locator('.palette-mode-toggle').click();
    assert.equal(await page.locator('html').evaluate(el => getComputedStyle(el).backgroundColor), 'rgb(250, 250, 250)');
    assert.deepEqual(pickerRequests(requests), [], 'CSS recovery and mode click remain picker-JS-free');
  });
});


test('first picker opening upgrades the cached core catalog after lazy CSS', {timeout: 60000}, async () => {
  await browserCase({lazyCSS: true}, async (page, origin, requests) => {
    await page.goto(origin + '/');
    const countCatalog = () => page.evaluate(() => {
      const raw = getComputedStyle(document.documentElement).getPropertyValue('--palette-manifest').trim();
      return JSON.parse(raw.slice(1, -1)).length;
    });
    assert.equal(await countCatalog(), 2, 'initial catalog contains only the configured pair');
    assert.equal(requests.includes('/palette.css'), false);
    await page.locator('.palette-mode-toggle').click();
    assert.equal(requests.includes('/palette.css'), false, 'dark/light remains eager without the catalog');
    await page.locator('.theme-picker-toggle').click();
    await page.locator('.theme-picker-panel').waitFor({state: 'visible'});
    assert.equal(await countCatalog(), 4);
    assert.equal(await page.locator('.palette-card').count(), 2, 'both light palettes appear on the first opening');
    await page.locator('.palette-card[data-name="nord-light"]').click();
    assert.equal(await page.locator('html').getAttribute('data-palette'), 'nord-light');
    assert.equal(requests.filter(p => p === '/palette.css').length, 1, 'full catalog loads once on intent');
  });
});
