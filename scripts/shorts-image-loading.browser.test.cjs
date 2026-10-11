// Browser regression tests for the Shorts image request policy and blur-up.
// Run with: node --test scripts/shorts-image-loading.browser.test.cjs
const {test} = require('node:test');
const assert = require('node:assert/strict');
const http = require('node:http');
const fs = require('node:fs');
const path = require('node:path');
// Install a pinned Playwright version via the dedicated Shorts CI workflow.
// PLAYWRIGHT_MODULE remains available for developer environments using playwright-core.
const playwright = require(process.env.PLAYWRIGHT_MODULE || 'playwright');

const player = fs.readFileSync(path.join(__dirname, '../pkg/themes/default/static/js/feed-shorts.js'));
const css = fs.readFileSync(path.join(__dirname, '../pkg/themes/default/static/css/feed-shorts.css'));
const requests = [];
const ids = Array.from({length: 8}, (_, i) => 'shots/p' + i);
const items = ids.map((id, index) => {
  if (index === 2) return {
    id, kind: 'video', title: 'Video ' + index, alt: '',
    src: '/media/clip-' + index + '.mp4',
    poster: '/media/poster-' + index + '.webp?w=720',
    thumb: '/media/poster-' + index + '.webp?w=240',
    placeholder: '/media/poster-' + index + '.webp?w=72', mime: 'video/mp4',
  };
  const source = '/media/photo-' + index + '.webp';
  return {
    id, kind: 'image', title: 'Photo ' + index, alt: 'Test photograph ' + index,
    src: source + '?w=1280', thumb: source + '?w=240', placeholder: source + '?w=72',
    poster: '',
  };
});
items[6].src = '/media/broken.webp?w=1280';

function pageHTML() {
  return `<!doctype html><html><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><link rel="stylesheet" href="/feed-shorts.css"></head><body class="shorts-body">
  <section class="shorts-viewer" id="shorts-viewer" data-manifest="/shorts/data/index.json" aria-label="Shorts media viewer">
    <div id="shorts-slides" class="shorts-slides" aria-hidden="true"></div>
    <header class="shorts-top"><a href="/shots/">Shots</a><span id="shorts-counter" class="shorts-counter"></span></header>
    <div id="shorts-message" class="shorts-message">Loading shorts…</div>
    <div class="shorts-hud"><div class="shorts-hud-copy"><h1 id="shorts-title" class="shorts-title"></h1><p id="shorts-description" class="shorts-description"></p><button id="shorts-more" class="shorts-more">More</button><a id="shorts-permalink" class="shorts-permalink"></a></div>
      <div class="shorts-actions"><button id="shorts-play" class="shorts-action" hidden></button><button id="shorts-sound" class="shorts-action" hidden></button><button id="shorts-up" class="shorts-action">Previous short</button><button id="shorts-down" class="shorts-action">Next short</button></div></div>
    <div id="shorts-end" class="shorts-end" hidden><button id="shorts-restart">Restart</button></div>
  </section><script defer src="/feed-shorts.js"></script></body></html>`;
}

function serve() {
  return http.createServer((req, res) => {
    const requestURL = new URL(req.url, 'http://localhost');
    requests.push(requestURL.pathname + requestURL.search);
    if (requestURL.pathname === '/') {
      res.setHeader('Content-Type', 'text/html; charset=utf-8');
      return res.end(pageHTML());
    }
    if (requestURL.pathname === '/shorts/data/index.json') {
      res.setHeader('Content-Type', 'application/json');
      return res.end(JSON.stringify({version: 1, total: ids.length, chunk_size: 128, ids}));
    }
    if (requestURL.pathname === '/shorts/data/0000.json') {
      res.setHeader('Content-Type', 'application/json');
      return res.end(JSON.stringify(items));
    }
    if (requestURL.pathname === '/feed-shorts.js') {
      res.setHeader('Content-Type', 'text/javascript');
      return res.end(player);
    }
    if (requestURL.pathname === '/feed-shorts.css') {
      res.setHeader('Content-Type', 'text/css');
      return res.end(css);
    }
    if (requestURL.pathname.startsWith('/media/')) {
      const width = Number(requestURL.searchParams.get('w') || 1280);
      if (requestURL.pathname.includes('broken.webp') && width === 1280) {
        res.writeHead(404); return res.end('missing image');
      }
      res.setHeader('Content-Type', 'image/svg+xml');
      const height = Math.round(width * 0.66);
      const body = `<svg xmlns="http://www.w3.org/2000/svg" width="${width}" height="${height}" viewBox="0 0 ${width} ${height}"><rect width="100%" height="100%" fill="#23647a"/><text x="12" y="32" fill="white">${width}</text></svg>`;
      const delay = width === 1280 ? 450 : width === 240 ? 100 : 15;
      return setTimeout(() => { if (!res.destroyed) res.end(body); }, delay);
    }
    if (requestURL.pathname === '/favicon.ico') {res.writeHead(204); return res.end();}
    res.writeHead(404); res.end('not found');
  });
}

async function openPage(browser, origin, connection, reducedMotion = false) {
  const context = await browser.newContext({viewport: {width: 390, height: 844}, deviceScaleFactor: 1, isMobile: true, hasTouch: true});
  await context.addInitScript((connectionValue) => {
    Object.defineProperty(navigator, 'connection', {value: connectionValue, configurable: true});
  }, connection);
  const page = await context.newPage();
  if (reducedMotion) await page.emulateMedia({reducedMotion: 'reduce'});
  await page.goto(origin, {waitUntil: 'domcontentloaded'});
  return {context, page};
}

test('Shorts uses a decoded blur-up and a bounded preview window', async (t) => {
  const server = serve();
  await new Promise((resolve) => server.listen(0, '127.0.0.1', resolve));
  const origin = `http://127.0.0.1:${server.address().port}/`;
  const browserOptions = {headless: true, args: ['--no-sandbox', '--disable-dev-shm-usage']};
  if (process.env.CHROMIUM_PATH) browserOptions.executablePath = process.env.CHROMIUM_PATH;
  const browser = await playwright.chromium.launch(browserOptions);
  const evidence = process.env.SHORTS_EVIDENCE_DIR || '/tmp/markata-shorts-browser-evidence';
  fs.mkdirSync(evidence, {recursive: true});
  t.after(async () => {await browser.close(); await new Promise((resolve) => server.close(resolve));});

  requests.length = 0;
  const normal = await openPage(browser, origin, {saveData: false, effectiveType: '3g', downlink: 1.5});
  const pageErrors = [];
  normal.page.on('pageerror', (error) => pageErrors.push(error.message));
  await normal.page.waitForFunction(() => {
    const img = document.querySelector('.shorts-image-placeholder');
    return img && img.complete && img.naturalWidth === 72;
  }, null, {timeout: 8000}).catch(async (error) => {
    const state = await normal.page.evaluate(() => ({title: document.querySelector('#shorts-title')?.textContent,
      slides: document.querySelectorAll('.shorts-slide').length,
      placeholder: document.querySelector('.shorts-image-placeholder')?.outerHTML,
      image: document.querySelector('.shorts-image--full')?.outerHTML}));
    throw new Error(`${error.message}; requests=${JSON.stringify(requests)}; pageErrors=${JSON.stringify(pageErrors)}; state=${JSON.stringify(state)}`);
  });
  const beforeSharp = await normal.page.locator('.shorts-slide[aria-hidden="false"]').boundingBox();
  await normal.page.screenshot({path: path.join(evidence, 'shorts-placeholder.png')});
  await normal.page.waitForFunction(() => document.querySelector('.shorts-media-wrap--sharp .shorts-image--full')?.naturalWidth === 1280);
  await normal.page.waitForTimeout(160);
  await normal.page.screenshot({path: path.join(evidence, 'shorts-sharp.png')});
  const afterSharp = await normal.page.locator('.shorts-slide[aria-hidden="false"]').boundingBox();
  assert.deepEqual(afterSharp, beforeSharp, 'image decode must not move the slide geometry');
  assert.equal(await normal.page.locator('.shorts-image--full').getAttribute('alt'), 'Test photograph 0');
  assert.equal(await normal.page.locator('.shorts-image-placeholder').getAttribute('aria-hidden'), 'true');
  assert.deepEqual(requests.filter((url) => url.startsWith('/media/')).sort(), [
    '/media/photo-0.webp?w=1280', '/media/photo-0.webp?w=72', '/media/photo-1.webp?w=240',
  ].sort(), 'normal mobile should load active image, active placeholder, and one upcoming preview');
  // Promoting a 240px preloaded slide must reuse its already-fetched preview,
  // rather than issuing a second 72px placeholder request.
  await normal.page.waitForFunction(() => {
    const preview = document.querySelector('.shorts-slide[data-index="1"] img.shorts-media');
    return preview && preview.complete && preview.naturalWidth === 240;
  });
  requests.length = 0;
  await normal.page.evaluate(() => window.__markataShorts.goTo(1));
  await normal.page.waitForFunction(() =>
    document.querySelector('.shorts-slide[aria-hidden="false"] .shorts-media-wrap--sharp .shorts-image--full')?.naturalWidth === 1280
  );
  const promotedPreview = await normal.page.locator('.shorts-slide[aria-hidden="false"] .shorts-image-placeholder').getAttribute('src');
  assert.equal(promotedPreview, '/media/photo-1.webp?w=240', 'reuses the loaded preview as the blurred placeholder');
  assert.ok(!requests.some((url) => url.includes('photo-1.webp?w=72')), 'must not download a redundant placeholder on swipe');
  await normal.context.close();

  requests.length = 0;
  const saveData = await openPage(browser, origin, {saveData: true, effectiveType: '4g', downlink: 20});
  await saveData.page.waitForFunction(() => document.querySelector('.shorts-media-wrap--sharp .shorts-image--full')?.naturalWidth === 1280);
  await saveData.page.waitForTimeout(100);
  assert.deepEqual(requests.filter((url) => url.startsWith('/media/')).sort(), [
    '/media/photo-0.webp?w=1280', '/media/photo-0.webp?w=72',
  ].sort(), 'Data Saver should not fetch adjacent images');
  await saveData.page.screenshot({path: path.join(evidence, 'shorts-save-data.png')});
  await saveData.context.close();

  requests.length = 0;
  const fast = await openPage(browser, origin, {saveData: false, effectiveType: '4g', downlink: 20});
  await fast.page.waitForFunction(() => document.querySelector('.shorts-media-wrap--sharp .shorts-image--full')?.naturalWidth === 1280);
  await fast.page.waitForSelector('#shorts-slides article[data-index="2"] .shorts-media[src*="poster-2.webp"]');
  assert.ok(!requests.some((url) => url.includes('.mp4')), 'inactive video bodies must not be prefetched');
  assert.equal(await fast.page.evaluate(() => matchMedia('(prefers-reduced-motion: reduce)').matches), false);

  requests.length = 0;
  await fast.page.evaluate(() => {window.__markataShorts.goTo(4); window.__markataShorts.goTo(5);});
  await fast.page.waitForFunction(() => document.querySelector('#shorts-title')?.textContent === 'Photo 5');
  await fast.page.waitForFunction(() => document.querySelector('.shorts-media-wrap--sharp .shorts-image--full')?.getAttribute('src')?.includes('photo-5.webp'));
  assert.ok(!requests.some((url) => url.includes('photo-4.webp')), 'stale async slides must not start image requests after rapid navigation');

  await fast.page.evaluate(() => window.__markataShorts.goTo(6));
  await fast.page.waitForFunction(() => document.querySelector('.shorts-slide[aria-hidden="false"].shorts-slide--empty') && !document.querySelector('.shorts-image-placeholder'));
  assert.equal(await fast.page.locator('#shorts-message').isVisible(), false, 'a broken image should leave the viewer navigable without a permanent placeholder');
  await fast.context.close();

  const reduced = await openPage(browser, origin, undefined, true);
  const transitionDuration = await reduced.page.locator('.shorts-image-placeholder').evaluate((el) => getComputedStyle(el).transitionDuration);
  assert.equal(transitionDuration, '0s', 'reduced motion should disable the image crossfade');
  await reduced.context.close();
});
