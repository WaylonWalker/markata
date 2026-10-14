// Focused browser coverage for the stale View Transition document lifecycle.
// Uses the Playwright installation bundled with agent-browser unless overridden.
const {test} = require('node:test');
const assert = require('node:assert/strict');
const http = require('node:http');
const fs = require('node:fs');
const path = require('node:path');
const {execFileSync} = require('node:child_process');
const playwright = require(process.env.PLAYWRIGHT_MODULE || path.join(
  execFileSync('npm', ['root', '-g'], {encoding: 'utf8'}).trim(),
  'agent-browser/node_modules/playwright-core',
));
const viewTransitions = fs.readFileSync(path.join(__dirname, '../pkg/themes/default/static/js/view-transitions.js'));

const state = {versions: {b: 1, etag: 1, noval: 1, fail: 1}, failedRevalidation: false, failedPrefetch: false, requests: []};
function pageHTML(route) {
  const content = route === 'a'
    ? '<h1>A</h1><a id="to-b" href="/b/">B</a><a id="to-etag" href="/etag/">ETag</a><a id="to-noval" href="/noval/">No validators</a><a id="to-fail" href="/fail/">Failure</a><a id="hash" href="/a/#section">Hash</a><a id="external" href="https://example.invalid/" target="_blank">External</a><a id="download" href="/file.pdf" download>Download</a><a id="blank" href="/b/?case=blank" target="_blank">New tab</a><div id="section">Section</div>'
    : `<h1 id="version">${route.toUpperCase()} v${state.versions[route] || 1}</h1><a id="to-a" href="/a/">A</a>`;
  return `<!doctype html><html data-loading-mode="auto" data-loading-policy="full-quality"><head><title>${route}</title></head><body><div id="view-transition-page"><main>${content}</main></div>${route === 'a' ? '<script>window.VIEW_TRANSITIONS_CONFIG={enabled:true};</script><script src="/view.js"></script>' : ''}</body></html>`;
}
function startServer() {
  state.requests = [];
  state.failedRevalidation = false;
  state.failedPrefetch = false;
  const server = http.createServer((req, res) => {
    const url = new URL(req.url, 'http://localhost');
    const requestRecord = {path: url.pathname, search: url.search, at: Date.now(), mode: req.headers['sec-fetch-mode'], ifNoneMatch: req.headers['if-none-match'], cacheControl: req.headers['cache-control']};
    state.requests.push(requestRecord);
    if (url.pathname === '/view.js') {res.setHeader('content-type', 'text/javascript'); return res.end(viewTransitions);}
    if (url.pathname === '/prefetch-fail/') {
      res.setHeader('content-type', 'text/html; charset=utf-8');
      res.setHeader('cache-control', 'no-store');
      if (req.headers['sec-fetch-mode'] === 'cors' && !state.failedPrefetch) {
        state.failedPrefetch = true;
        requestRecord.status = 503;
        res.statusCode = 503;
        return res.end('speculative request failed');
      }
      requestRecord.status = 200;
      return res.end(pageHTML('b'));
    }
    if (url.pathname === '/file.pdf') {res.setHeader('content-type', 'application/pdf');res.setHeader('content-disposition', 'attachment; filename="file.pdf"');return res.end('%PDF-1.4');}
    const match = url.pathname.match(/^\/(b|etag|noval|fail)\/$/);
    if (match) {
      const route = match[1];
      const version = state.versions[route];
      const body = pageHTML(route);
      res.setHeader('content-type', 'text/html; charset=utf-8');
      if (route !== 'noval') {
        const etag = `"${route}-${version}"`;
        res.setHeader('etag', etag);
        res.setHeader('last-modified', 'Wed, 01 Jan 2025 00:00:00 GMT');
        if (req.headers['if-none-match'] === etag) {requestRecord.status = 304; res.statusCode = 304; return res.end();}
      }
      res.setHeader('cache-control', 'public, max-age=60');
      if (route === 'fail' && version > 1 && req.headers['sec-fetch-mode'] === 'cors' && !state.failedRevalidation) {
        state.failedRevalidation = true;
        res.setHeader('cache-control', 'no-store');
        requestRecord.status = 503;
        res.statusCode = 503;
        return res.end('revalidation failed');
      }
      requestRecord.status = 200;
      if (route === 'b') return setTimeout(() => res.end(body), 100);
      return res.end(body);
    }
    if (url.pathname === '/a/' || url.pathname === '/a') {
      res.setHeader('content-type', 'text/html; charset=utf-8');
      return res.end(pageHTML('a'));
    }
    res.statusCode = 404;
    res.end('not found');
  });
  return server;
}
async function waitForRequest(page, route) {
  return page.waitForResponse(response => new URL(response.url()).pathname === `/${route}/`, {timeout: 5000});
}
async function waitForNavigation(page, route) {
  await page.waitForFunction(route => location.pathname === `/${route}/` && document.querySelector('#version'), route);
  await page.waitForFunction(() => !!window.__lastViewTransitionMetrics || !window.VIEW_TRANSITIONS_CONFIG);
}

test('Chromium reuses fresh prefetches and revalidates stale documents safely', {timeout: 90000}, async t => {
  const server = startServer();
  await new Promise(resolve => server.listen(0, '127.0.0.1', resolve));
  const browser = await playwright.chromium.launch({headless: true, executablePath: '/usr/bin/chromium', args: ['--no-sandbox']});
  try {
    const origin = `http://127.0.0.1:${server.address().port}`;
    const context = await browser.newContext();
    const page = await context.newPage();
    page.on('console', message => t.diagnostic(`browser console: ${message.type()} ${message.text()}`));
    page.on('pageerror', error => t.diagnostic(`browser error: ${error.message}`));
    await page.goto(`${origin}/a/`);
    t.diagnostic(`navigator: ${JSON.stringify(await page.evaluate(() => ({transition: !!document.startViewTransition, prefetch: typeof window.prefetchViewTransitionUrl, href: document.querySelector('#to-b')?.href, config: window.VIEW_TRANSITIONS_CONFIG})))}`);
    await page.mouse.move(0, 0);
    const firstPrefetchResponse = waitForRequest(page, 'b');
    await page.locator('#to-b').hover();
    await page.locator('#to-b').focus();
    await page.mouse.move(0, 0);
    await page.locator('#to-b').hover();
    await firstPrefetchResponse;
    await page.waitForTimeout(150);
    const warmHoverToClickMs = Date.now() - state.requests.find(request => request.path === '/b/').at;
    const warmStart = Date.now();
    await page.locator('#to-b').click();
    await waitForNavigation(page, 'b');
    const warmDuration = Date.now() - warmStart;
    assert.match(await page.locator('#version').innerText(), /B v1/);
    assert.equal(state.requests.filter(request => request.path === '/b/').length, 1, 'fresh click must not request the page again');
    assert.equal(await page.getAttribute('html', 'data-loading-mode'), 'auto');
    assert.equal(await page.getAttribute('html', 'data-loading-policy'), 'full-quality');
    const warmSamples = [{prefetchToClick: warmHoverToClickMs, navigation: warmDuration}];
    for (let i = 1; i <= 2; i++) {
      await page.goto(`${origin}/a/`);
      await page.locator('#to-b').evaluate((link, i) => {link.href += `?case=warm-${i}`;}, i);
      await page.mouse.move(0, 0);
      const response = waitForRequest(page, 'b');
      await page.locator('#to-b').hover();
      await page.locator('#to-b').focus();
      await response;
      await page.waitForTimeout(150);
      const request = state.requests.filter(item => item.path === '/b/' && item.search === `?case=warm-${i}`)[0];
      const clickStarted = Date.now();
      await page.locator('#to-b').click();
      await waitForNavigation(page, 'b');
      warmSamples.push({prefetchToClick: clickStarted - request.at, navigation: Date.now() - clickStarted});
      assert.equal(state.requests.filter(item => item.path === '/b/' && item.search === `?case=warm-${i}`).length, 1, 'fresh click reuses the resolved prefetch');
    }
    const median = key => warmSamples.map(sample => sample[key]).sort((a, b) => a - b)[Math.floor(warmSamples.length / 2)];
    t.diagnostic(`three fresh samples (prefetch-to-click ms: ${warmSamples.map(sample => sample.prefetchToClick).join(', ')}; click-to-complete ms: ${warmSamples.map(sample => sample.navigation).join(', ')}); medians ${median('prefetchToClick')} ms / ${median('navigation')} ms`);

    await page.goto(`${origin}/a/`);
    await page.locator('#to-b').evaluate(link => {link.href += '?case=stale';});
    await page.mouse.move(0, 0);
    await Promise.all([waitForRequest(page, 'b'), page.locator('#to-b').hover()]);
    await page.waitForTimeout(5200);
    state.versions.b = 2;
    await page.locator('#to-b').click();
    await waitForNavigation(page, 'b');
    assert.match(await page.locator('#version').innerText(), /B v2/);
    const bRequests = state.requests.filter(request => request.path === '/b/' && request.search === '?case=stale');
    t.diagnostic(`B requests: ${JSON.stringify(bRequests)}`);
    assert.equal(bRequests.length, 2, 'stale speculative document is followed by one conditional revalidation');
    assert.equal(bRequests[1].ifNoneMatch, '"b-1"');
    assert.equal(bRequests[1].mode, 'cors');

    await page.goto(`${origin}/a/`);
    await page.locator('#to-etag').evaluate(link => {link.href += '?case=etag';});
    await page.mouse.move(0, 0);
    await Promise.all([waitForRequest(page, 'etag'), page.locator('#to-etag').hover()]);
    await page.waitForTimeout(5200);
    await page.locator('#to-etag').click();
    await waitForNavigation(page, 'etag');
    assert.match(await page.locator('#version').innerText(), /ETAG v1/);
    assert.ok(state.requests.some(request => request.path === '/etag/' && request.ifNoneMatch === '"etag-1"'), 'stale fetch sends ETag validator');
    assert.ok(state.requests.some(request => request.path === '/etag/' && request.status === 304), 'unchanged stale content is served from a successful 304 validation');
    assert.equal(state.requests.filter(request => request.path === '/etag/').length, 2, 'prefetch and conditional validation occurred');

    await page.goto(`${origin}/a/`);
    await page.locator('#to-noval').evaluate(link => {link.href += '?case=noval';});
    await page.mouse.move(0, 0);
    await Promise.all([waitForRequest(page, 'noval'), page.locator('#to-noval').hover()]);
    await page.waitForTimeout(5200);
    await page.locator('#to-noval').click();
    await waitForNavigation(page, 'noval');
    assert.match(await page.locator('#version').innerText(), /NOVAL v1/);
    assert.equal(state.requests.filter(request => request.path === '/noval/').length, 2, 'no-validator entry is fetched again');
    assert.ok(state.requests.filter(request => request.path === '/noval/').every(request => !request.ifNoneMatch));

    await page.goto(`${origin}/a/`);
    await page.locator('#to-fail').evaluate(link => {link.href += '?case=fail';});
    await page.mouse.move(0, 0);
    await Promise.all([waitForRequest(page, 'fail'), page.locator('#to-fail').hover()]);
    await page.waitForTimeout(5200);
    state.versions.fail = 2;
    await page.locator('#to-fail').click();
    await page.waitForURL(url => url.pathname === '/fail/', {waitUntil: 'domcontentloaded', timeout: 10000});
    t.diagnostic(`after fallback URL: ${page.url()}, requests: ${JSON.stringify(state.requests.filter(request => request.path === '/fail/'))}`);
    await page.waitForSelector('#version', {timeout: 5000});
    assert.match(await page.locator('#version').innerText(), /FAIL v2/);
    assert.ok(state.requests.some(request => request.path === '/fail/' && request.mode === 'cors' && request.ifNoneMatch === '"fail-1"'), 'failed validation request happened');
    assert.ok(state.requests.some(request => request.path === '/fail/' && request.mode === 'navigate'), 'ordinary document navigation recovered the failure');
    t.diagnostic(`failed-validation requests: ${JSON.stringify(state.requests.filter(request => request.path === '/fail/'))}`);

    await page.goto(`${origin}/a/`);
    const evictionURLs = Array.from({length: 8}, (_, i) => `${origin}/b/?case=evict-${i}`);
    const evictionResponses = evictionURLs.map(url =>
      page.waitForResponse(response => response.url() === url, {timeout: 5000}),
    );
    await page.evaluate(urls => urls.forEach(url => window.prefetchViewTransitionUrl(url)), evictionURLs);
    await Promise.all(evictionResponses);
    await page.waitForTimeout(250);
    await page.waitForTimeout(5200);

    const ninthEvictionURL = `${origin}/b/?case=evict-8`;
    const ninthEvictionResponse = page.waitForResponse(
      response => response.url() === ninthEvictionURL,
      {timeout: 5000},
    );
    await page.evaluate(url => window.prefetchViewTransitionUrl(url), ninthEvictionURL);
    await ninthEvictionResponse;
    assert.equal(
      state.requests.filter(request => request.path === '/b/' && request.search.startsWith('?case=evict-')).length,
      9,
      'expired entries are evicted so a ninth URL can be prefetched',
    );

    await page.goto(`${origin}/a/`);
    const bRequestCount = state.requests.filter(request => request.path === '/b/').length;
    for (const selector of ['#external', '#download', '#blank']) {
      await page.mouse.move(0, 0);
      await page.locator(selector).hover();
    }
    await page.mouse.move(0, 0);
    await page.locator('#hash').click();
    await page.waitForFunction(() => location.hash === '#section');
    assert.equal(state.requests.filter(request => request.path === '/b/').length, bRequestCount, 'external, download, target-blank and hash links do not prefetch documents');
    await page.evaluate(() => history.replaceState(null, '', '/a/'));
    const [externalPopup] = await Promise.all([
      context.waitForEvent('page'),
      page.locator('#external').click(),
    ]);
    await externalPopup.close();
    const [download] = await Promise.all([
      page.waitForEvent('download'),
      page.locator('#download').click(),
    ]);
    assert.equal(download.suggestedFilename(), 'file.pdf');
    const [popup] = await Promise.all([
      context.waitForEvent('page'),
      page.locator('#blank').click(),
    ]);
    await popup.waitForURL(url => url.pathname === '/b/');
    assert.equal(new URL(popup.url()).pathname, '/b/', 'target-blank click opens normal browser navigation');
    await popup.close();

    const [modifierPopup] = await Promise.all([
      context.waitForEvent('page'),
      page.locator('#to-b').click({modifiers: ['Control']}),
    ]);
    await modifierPopup.waitForURL(url => url.pathname === '/b/');
    assert.equal(new URL(modifierPopup.url()).pathname, '/b/', 'modifier click opens normal browser navigation');
    await modifierPopup.close();

    await page.goto(`${origin}/a/`);
    await context.addInitScript(() => Object.defineProperty(navigator, 'connection', {value: {saveData: true}, configurable: true}));
    await page.reload();
    await page.locator('#to-b').evaluate(link => {link.href += '?case=save-data';});
    const beforeSaveDataHover = state.requests.filter(request => request.path === '/b/').length;
    await page.mouse.move(0, 0);
    await page.locator('#to-b').hover();
    await page.waitForTimeout(200);
    assert.equal(state.requests.filter(request => request.path === '/b/').length, beforeSaveDataHover, 'Save Data suppresses speculative prefetch');
    await page.locator('#to-b').click();
    await page.waitForFunction(() => location.pathname === '/b/' && !!document.querySelector('#version'));
    assert.match(await page.locator('#version').innerText(), /B v2/);

    await page.goto(`${origin}/a/`);
    const historyLink = page.locator('#to-b');
    await historyLink.evaluate(link => {link.href += '?case=history';});
    await historyLink.click();
    await waitForNavigation(page, 'b');
    await page.goBack();
    await page.waitForFunction(() => location.pathname === '/a/' && !!document.querySelector('#to-b'));
    await page.goForward();
    await page.waitForFunction(() => location.pathname === '/b/' && !!document.querySelector('#version'));
  } finally {
    await browser.close();
    await new Promise(resolve => server.close(resolve));
  }
});

test('Firefox smoke: transition runtime loads and hash/special links remain browser navigation', {timeout: 30000}, async () => {
  state.versions.b = 1;
  const server = startServer();
  await new Promise(resolve => server.listen(0, '127.0.0.1', resolve));
  const browser = await playwright.firefox.launch({headless: true, executablePath: process.env.ADAPTIVE_FIREFOX_EXECUTABLE || '/home/waylon/.cache/ms-playwright/firefox-1511/firefox/firefox'});
  try {
    const page = await browser.newPage();
    const pageErrors = [];
    page.on('pageerror', error => pageErrors.push(error.message));
    const origin = `http://127.0.0.1:${server.address().port}`;
    await page.goto(`${origin}/a/`);
    assert.ok(await page.evaluate(() => typeof window.prefetchViewTransitionUrl === 'function') || !await page.evaluate(() => !!document.startViewTransition));
    const before = state.requests.length;
    await page.locator('#hash').click();
    await page.waitForFunction(() => location.hash === '#section');
    assert.equal(state.requests.length, before, 'same-document hash does not fetch a document');
    await page.locator('#to-b').click();
    await page.waitForSelector('#version');
    assert.match(await page.locator('#version').innerText(), /B v1/);
    assert.deepEqual(pageErrors, [], 'Firefox smoke has no uncaught script exceptions');
  } finally {
    await browser.close();
    await new Promise(resolve => server.close(resolve));
  }
});

test('failed hover prefetch does not emit unhandled rejections and click retries', {timeout: 30000}, async t => {
  const server = startServer();
  await new Promise(resolve => server.listen(0, '127.0.0.1', resolve));
  const browser = await playwright.chromium.launch({headless: true, executablePath: '/usr/bin/chromium', args: ['--no-sandbox']});
  try {
    const page = await browser.newPage();
    const pageErrors = [];
    page.on('pageerror', err => pageErrors.push(err.message));
    const origin = `http://127.0.0.1:${server.address().port}`;
    await page.goto(`${origin}/a/`);
    await page.evaluate(() => {
      const a = document.createElement('a');
      a.id = 'to-prefetch-fail';
      a.href = '/prefetch-fail/';
      a.textContent = 'Fail then recover';
      document.body.appendChild(a);
    });
    const response = page.waitForResponse(r => new URL(r.url()).pathname === '/prefetch-fail/');
    await page.locator('#to-prefetch-fail').hover();
    await response;
    await page.waitForTimeout(150);
    assert.deepEqual(pageErrors, [], 'speculative network failures must not reject unobserved promises');
    await page.locator('#to-prefetch-fail').click();
    await page.waitForFunction(() => location.pathname === '/prefetch-fail/' && !!document.querySelector('#version'));
    assert.match(await page.locator('#version').innerText(), /B v1/);
    assert.ok(state.requests.filter(x => x.path === '/prefetch-fail/').length >= 2, 'click must retry failed speculative request');
  } finally {
    await browser.close();
    await new Promise(resolve => server.close(resolve));
  }
});
