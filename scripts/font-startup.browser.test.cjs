const {test} = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const {chromium} = require(process.env.PLAYWRIGHT_MODULE || 'playwright');
const root = path.join(__dirname, '..');
const script = name => fs.readFileSync(path.join(root, 'pkg/themes/default/static/js', name), 'utf8');
const font = fs.readFileSync(path.join(root, 'internal/fontcatalog/bitter/bitter-full.woff2'));

for (const observerAvailable of [true, false]) {
  test(`title and highlight geometry survive font loading, resize and navigation, ResizeObserver=${observerAvailable}`, async () => {
    const browser = await chromium.launch({headless: true,
      executablePath: process.env.CHROMIUM_EXECUTABLE || undefined, args: ['--no-sandbox']});
    let releaseFont;
    const fontGate = new Promise(resolve => { releaseFont = resolve; });
    try {
      const page = await browser.newPage({viewport: {width: 412, height: 823}});
      if (!observerAvailable) await page.addInitScript(() => { window.ResizeObserver = undefined; });
      await page.route('http://startup.test/**', route => route.fulfill({contentType: 'text/html', body: `
        <!doctype html><html data-fontpack="test"><head><style>
        @font-face {font-family: Delayed; src: url('/display.woff2'); font-display: swap}
        html {font-size:16px} body {margin:0; font-family:Delayed,Arial,sans-serif}
        .post-header__title-row {width:90vw} h1 {font-size:64px; line-height:1.1; margin:0}
        .post-content {width:70vw} h2 {font-size:32px}
        .heading-highlight {position:relative} .heading-highlight__contour {position:absolute; inset:0; width:100%; height:100%; pointer-events:none}
        </style></head><body><header class="post-header"><div class="post-header__title-row">
        <h1 data-title-size="large">A practical title</h1></div></header>
        <article class="post-content"><h2><span class="heading-highlight"><mark>A highlighted heading that wraps</mark></span></h2></article>
        </body></html>`}));
      await page.route('http://startup.test/display.woff2', async route => {
        await fontGate;
        await route.fulfill({contentType: 'font/woff2', body: font});
      });
      await page.goto('http://startup.test/', {waitUntil: 'domcontentloaded'});
      await page.addScriptTag({content: script('heading-highlights.js')});
      await page.addScriptTag({content: script('fit-post-titles.js')});
      const title = page.locator('h1');
      const contour = page.locator('.heading-highlight__contour path');
      await page.waitForFunction(() => !!document.querySelector('h1').style.fontSize && !!document.querySelector('.heading-highlight__contour path'));
      const fallbackSize = await title.evaluate(el => getComputedStyle(el).fontSize);
      const fallbackPath = await contour.getAttribute('d');
      // Give the font-ready subscription a rendering opportunity before the
      // delayed font completes, reproducing a font swap after first paint.
      await page.evaluate(() => new Promise(resolve => requestAnimationFrame(() => requestAnimationFrame(resolve))));
      releaseFont();
      await page.evaluate(() => document.fonts.ready);
      await page.waitForFunction(size => getComputedStyle(document.querySelector('h1')).fontSize !== size, fallbackSize);
      await page.waitForFunction(old => document.querySelector('.heading-highlight__contour path').getAttribute('d') !== old, fallbackPath);
      async function assertTitleFits() {
        const geometry = await title.evaluate(el => {
          const range = document.createRange(); range.selectNodeContents(el);
          return {text: range.getBoundingClientRect().width, available: el.clientWidth,
            size: parseFloat(getComputedStyle(el).fontSize), whiteSpace: el.style.whiteSpace};
        });
        assert.ok(geometry.text <= geometry.available + 1, JSON.stringify(geometry));
        assert.ok(geometry.size >= 20, JSON.stringify(geometry));
        assert.equal(geometry.whiteSpace, '');
      }
      await assertTitleFits();
      const smallSize = await title.evaluate(el => getComputedStyle(el).fontSize);
      const smallPath = await contour.getAttribute('d');
      await page.setViewportSize({width: 900, height: 823});
      await page.waitForFunction(size => getComputedStyle(document.querySelector('h1')).fontSize !== size, smallSize);
      await page.waitForFunction(old => document.querySelector('.heading-highlight__contour path').getAttribute('d') !== old, smallPath);
      await assertTitleFits();
      const oldPath = await contour.getAttribute('d');
      await page.locator('mark').evaluate(el => { el.textContent = 'Edited heading'; });
      await page.waitForFunction(old => document.querySelector('.heading-highlight__contour path').getAttribute('d') !== old, oldPath);
      await page.evaluate(() => {
        document.querySelector('h1').outerHTML = '<h1 data-title-size="large">A much longer navigated article title</h1>';
        window.dispatchEvent(new Event('view-transition-complete'));
      });
      await page.waitForFunction(() => !!document.querySelector('h1').style.fontSize);
      await assertTitleFits();
      await page.evaluate(() => {
        document.querySelector('h1').outerHTML = '<h1 data-title-size="large">' + 'A very long readable title '.repeat(20) + '</h1>';
        window.dispatchEvent(new Event('view-transition-complete'));
      });
      await page.evaluate(() => new Promise(resolve => requestAnimationFrame(() => requestAnimationFrame(resolve))));
      await page.waitForFunction(() => {
        const el = document.querySelector('h1');
        return el.style.whiteSpace === '' && el.style.fontSize === '' && el.getBoundingClientRect().height > 100;
      });
      assert.equal(await title.evaluate(el => getComputedStyle(el).fontSize), '64px', 'unfittable titles retain their readable wrapping size');
    } finally {
      releaseFont();
      await browser.close();
    }
  });
}

for (const observerAvailable of [true, false]) {
  test(`progress follows article height and navigation, ResizeObserver=${observerAvailable}`, async () => {
    const paths = ['pkg/themes/default/templates/base.html', 'templates/base.html'];
    const scripts = paths.map(file => {
      const template = fs.readFileSync(path.join(root, file), 'utf8');
      return [...template.matchAll(/<script[^>]*>([\s\S]*?)<\/script>/g)]
        .map(match => match[1]).find(source => source.includes('function initArticleProgressIndicator'));
    });
    assert.ok(scripts[0]);
    assert.equal(scripts[0], scripts[1], 'both template trees must use the same progress behavior');
    const browser = await chromium.launch({headless: true,
      executablePath: process.env.CHROMIUM_EXECUTABLE || undefined, args: ['--no-sandbox']});
    try {
      const page = await browser.newPage({viewport: {width: 412, height: 823}});
      if (!observerAvailable) await page.addInitScript(() => { window.ResizeObserver = undefined; });
      await page.route('http://progress.test/', route => route.fulfill({contentType: 'text/html', body: `
        <!doctype html><style>body{margin:0}.article-progress{position:fixed;top:0;height:3px;width:100%}.article-progress--hidden{display:none}</style>
        <div class="article-progress"><div class="article-progress__indicator"></div></div>
        <article class="post" style="height:3000px"></article>`}));
      await page.goto('http://progress.test/');
      await page.addScriptTag({content: scripts[0]});
      await page.waitForFunction(() => document.querySelector('.article-progress__indicator').style.transform === 'scaleX(0)');
      await page.evaluate(() => window.scrollTo(0, 1000));
      await page.waitForFunction(() => parseFloat(document.querySelector('.article-progress__indicator').style.transform.slice(7)) > 0.4);
      const initial = await page.locator('.article-progress__indicator').evaluate(el => parseFloat(el.style.transform.slice(7)));
      await page.locator('article').evaluate(el => { el.style.height = '6000px'; });
      if (!observerAvailable) await page.evaluate(() => window.dispatchEvent(new Event('resize')));
      await page.waitForFunction(previous => parseFloat(document.querySelector('.article-progress__indicator').style.transform.slice(7)) < previous / 2, initial);
      await page.locator('article').evaluate(el => { el.style.height = '500px'; });
      if (!observerAvailable) await page.evaluate(() => window.dispatchEvent(new Event('resize')));
      await page.waitForFunction(() => document.querySelector('.article-progress').classList.contains('article-progress--hidden'));
      await page.evaluate(() => {
        document.querySelector('article').outerHTML = '<article class="post" style="height:4000px"></article>';
        window.dispatchEvent(new Event('view-transition-complete'));
      });
      await page.waitForFunction(() => !document.querySelector('.article-progress').classList.contains('article-progress--hidden'));
      await page.evaluate(() => window.scrollTo(0, document.documentElement.scrollHeight));
      await page.waitForFunction(() => document.querySelector('.article-progress__indicator').style.transform === 'scaleX(1)');
    } finally {
      await browser.close();
    }
  });
}
