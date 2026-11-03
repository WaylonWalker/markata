const assert = require('node:assert/strict');
const fs = require('node:fs');
const http = require('node:http');
const path = require('node:path');
const test = require('node:test');
const { chromium } = require(process.env.PLAYWRIGHT_MODULE || 'playwright');

const siteDir = path.resolve(__dirname, '../docs-site/_site');
const screenshotDir = process.env.DOCS_SCREENSHOT_DIR || path.join(siteDir, '.screenshots');
const mimeTypes = {
  '.css': 'text/css', '.html': 'text/html; charset=utf-8', '.ico': 'image/x-icon',
  '.js': 'text/javascript', '.json': 'application/json', '.png': 'image/png',
  '.svg': 'image/svg+xml', '.woff': 'font/woff', '.woff2': 'font/woff2',
};

function serveSite() {
  const server = http.createServer((request, response) => {
    let pathname;
    try { pathname = decodeURIComponent(new URL(request.url, 'http://localhost').pathname); }
    catch { response.writeHead(400).end(); return; }
    if (pathname.endsWith('/')) pathname += 'index.html';
    const file = path.resolve(siteDir, `.${pathname}`);
    if (!file.startsWith(`${siteDir}${path.sep}`)) { response.writeHead(403).end(); return; }
    fs.readFile(file, (error, data) => {
      if (error) { response.writeHead(404).end(); return; }
      response.writeHead(200, { 'content-type': mimeTypes[path.extname(file)] || 'application/octet-stream' });
      response.end(data);
    });
  });
  return new Promise((resolve) => server.listen(0, '127.0.0.1', () => {
    resolve({ server, origin: `http://127.0.0.1:${server.address().port}` });
  }));
}

test('built documentation works across layouts, themes, sizes, and viewport widths', async (t) => {
  assert.ok(fs.existsSync(path.join(siteDir, 'index.html')), 'build the docs site before browser verification');
  fs.mkdirSync(screenshotDir, { recursive: true });
  const { server, origin } = await serveSite();
  const browser = await chromium.launch({ headless: true, executablePath: process.env.CHROMIUM_EXECUTABLE || undefined, args: ['--no-sandbox'] });
  t.after(async () => { await browser.close(); server.close(); });

  for (const route of ['/', '/docs/quickstart/']) {
    for (const width of [1440, 390]) {
      for (const mode of ['light', 'dark']) {
        const context = await browser.newContext({ viewport: { width, height: 1000 }, colorScheme: mode });
        await context.addInitScript((theme) => localStorage.setItem('color-mode', theme), mode);
        const page = await context.newPage();
        const failedAssets = [];
        const consoleErrors = [];
        page.on('response', (response) => {
          if (response.url().startsWith(origin) && response.status() >= 400) failedAssets.push(`${response.status()} ${response.url()}`);
        });
        page.on('pageerror', (error) => consoleErrors.push(error.message));
        page.on('console', (message) => {
          if (message.type() === 'error') consoleErrors.push(message.text());
        });
        const response = await page.goto(`${origin}${route}`, { waitUntil: 'networkidle' });
        assert.equal(response.status(), 200, `${route} responds successfully`);
        await page.evaluate(() => document.fonts.ready);
        await page.waitForTimeout(100);
        assert.equal(await page.locator('html').getAttribute('data-theme'), mode, `${route} applies ${mode} theme`);
        assert.equal(await page.locator('html').getAttribute('data-text-size'), 'medium', 'docs use medium text size');
        assert.ok(await page.locator('link[rel="stylesheet"]').count() > 0, 'page links stylesheets');
        assert.equal(await page.locator('link[rel="icon"]').getAttribute('href'), '/favicon.ico', 'page links the site favicon');
        assert.deepEqual(failedAssets, [], 'same-origin assets load successfully');
        assert.deepEqual(consoleErrors, [], `${route} has no browser console errors`);
        const overflow = await page.evaluate(() => ({
          width: document.documentElement.scrollWidth,
          offenders: [...document.querySelectorAll('body *')].filter((element) => element.getBoundingClientRect().right > innerWidth + 1).slice(0, 8).map((element) => {
            const rect = element.getBoundingClientRect();
            const parent = element.parentElement;
            const box = (node) => { const r = node.getBoundingClientRect(); return `${node.tagName}.${typeof node.className === 'string' ? node.className : ''}=${Math.round(r.left)}..${Math.round(r.right)}`; };
            return `${box(element)} <- ${parent ? box(parent) : ''}`;
          }),
          header: (() => { const element = document.querySelector('.site-header .container'); const rect = element.getBoundingClientRect(); return {rect: [rect.left,rect.right],display:getComputedStyle(element).display,columns:getComputedStyle(element).gridTemplateColumns,html:getComputedStyle(document.documentElement).overflowX,body:getComputedStyle(document.body).overflowX}; })(),
        }));
        const contentBounds = await page.locator('main').first().evaluate((element) => {
          const rect = element.getBoundingClientRect();
          return { left: rect.left, right: rect.right, clientWidth: element.clientWidth, scrollWidth: element.scrollWidth };
        });
        assert.ok(contentBounds.left >= -1 && contentBounds.right <= width + 1, `${route} main content fits ${width}px: ${JSON.stringify(contentBounds)}`);
        assert.ok(overflow.width <= width + 1, `${route} has no page overflow at ${width}px: ${JSON.stringify(overflow)}; main: ${JSON.stringify(contentBounds)}`);

        if (route === '/') {
          assert.equal((await page.locator('.hero-title').innerText()).replace(/\s+/g, ' ').trim(), 'Your words. On your web.');
          assert.equal((await page.locator('.hero-eyebrow').innerText()).toLowerCase(), 'markdown to website');
          assert.equal(await page.locator('.hero-title-line').count(), 2);
          assert.equal(await page.locator('[data-md-tabs] [role="tab"]').count(), 6);
          await page.screenshot({ path: path.join(screenshotDir, `home-${width}-${mode}.png`), fullPage: width > 600 });
        } else {
          assert.ok(await page.locator('nav[aria-label="Documentation navigation"]').count(), 'docs use the documentation sidebar');
          assert.ok(await page.locator('.post-content pre').count(), 'quickstart includes rendered code blocks');
          if (width <= 600) {
            const sidebar = page.locator('#docs-sidebar');
            const menu = page.locator('.mobile-menu-toggle');
            assert.equal(await sidebar.getAttribute('data-open'), 'false', 'mobile docs start with the sidebar closed');
            await menu.click();
            assert.equal(await sidebar.getAttribute('data-open'), 'true', 'mobile navigation opens from its menu button');
            assert.equal(await menu.getAttribute('aria-expanded'), 'true', 'mobile navigation exposes its expanded state');
            await page.keyboard.press('Escape');
            assert.equal(await sidebar.getAttribute('data-open'), 'false', 'Escape closes mobile navigation');
            await menu.evaluate((button) => button.blur());
          }
          const surfaces = await page.locator('.post-content pre').first().evaluate((pre) => ({
            background: getComputedStyle(pre).backgroundColor,
            page: getComputedStyle(document.body).backgroundColor,
          }));
          assert.notEqual(surfaces.background, surfaces.page, `code blocks have a distinct ${mode} surface`);
          const copyBox = await page.locator('.post-content .code-block__copy').first().evaluate((button) => {
            const buttonRect = button.getBoundingClientRect();
            const headerRect = button.closest('.code-block__header').getBoundingClientRect();
            return { fits: buttonRect.left >= headerRect.left - 1 && buttonRect.right <= headerRect.right + 1 && buttonRect.top >= headerRect.top - 1 && buttonRect.bottom <= headerRect.bottom + 1, button: [buttonRect.left, buttonRect.top, buttonRect.right, buttonRect.bottom], header: [headerRect.left, headerRect.top, headerRect.right, headerRect.bottom] };
          });
          assert.ok(copyBox.fits, `copy control fits within its code-block header: ${JSON.stringify(copyBox)}`);
          await page.screenshot({ path: path.join(screenshotDir, `quickstart-${width}-${mode}.png`), fullPage: width > 600 });
        }
        await context.close();
      }
    }
  }

  const context = await browser.newContext({ viewport: { width: 1280, height: 900 } });
  await context.addInitScript(() => {
    window.__clipboardWrites = [];
    Object.defineProperty(navigator, 'clipboard', { configurable: true, value: { writeText: async (text) => window.__clipboardWrites.push(text) } });
  });
  const page = await context.newPage();
  await page.goto(`${origin}/`, { waitUntil: 'networkidle' });
  page.setDefaultTimeout(5000);
  const tabs = page.locator('[data-md-tabs] [role="tab"]');
  await tabs.first().focus();
  await page.keyboard.press('ArrowRight');
  assert.equal(await tabs.nth(1).getAttribute('aria-selected'), 'true', 'ArrowRight selects the next install method');
  assert.equal(await tabs.nth(1).getAttribute('tabindex'), '0');
  assert.equal(await page.locator('#md-tab-1-1-panel').getAttribute('hidden'), null, 'selected tab panel is visible');
  const copyButton = page.locator('[data-md-tabs] [role="tabpanel"]:not([hidden]) .code-block__copy').first();
  await copyButton.click();
  await page.waitForFunction(() => [...document.querySelectorAll('.code-block__copy')].some((button) => button.getAttribute('aria-label') === 'Code copied' || button.getAttribute('aria-label') === 'Select code to copy manually'));
  const copyResult = await copyButton.getAttribute('aria-label');
  assert.equal(copyResult, 'Code copied', `copy control should succeed; clipboard writes: ${JSON.stringify(await page.evaluate(() => window.__clipboardWrites))}`);
  assert.match(await page.evaluate(() => window.__clipboardWrites[0]), /curl -fsSL/);
  await context.close();

  const noScriptContext = await browser.newContext({ javaScriptEnabled: false, viewport: { width: 390, height: 844 } });
  const noScriptPage = await noScriptContext.newPage();
  await noScriptPage.goto(`${origin}/`, { waitUntil: 'networkidle' });
  assert.equal(await noScriptPage.locator('[data-md-tabs] .md-tab').count(), 6, 'all install choices remain in the HTML without JavaScript');
  assert.equal(await noScriptPage.locator('[data-md-tabs] .md-tab:visible').count(), 6, 'all install choices remain readable without JavaScript');
  assert.equal(await noScriptPage.locator('[data-md-tabs] [role="tablist"]').getAttribute('hidden'), '');
  await noScriptContext.close();
});
