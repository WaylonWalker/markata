// The visible navbar must keep its geometry while deferred search styles arrive.
const {test} = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const {chromium} = require(process.env.PLAYWRIGHT_MODULE || 'playwright');
const cssRoot = path.join(__dirname, '../pkg/themes/default/static/css');
const readCSS = name => fs.readFileSync(path.join(cssRoot, name + '.css'), 'utf8');
const coreCSS = ['variables', 'main', 'components'].map(readCSS).join('\n');
const searchCSS = readCSS('search');

for (const width of [412, 1440]) {
  test(`navbar search reserves layout at ${width}px`, async () => {
    const browser = await chromium.launch({headless: true, executablePath: process.env.CHROMIUM_EXECUTABLE || undefined, args: ['--no-sandbox']});
    try {
      const page = await browser.newPage({viewport: {width, height: 900}});
      await page.setContent(`<!doctype html><html><head><style>${coreCSS}</style></head>
        <body><header class="site-header"><div class="container">
        <a class="site-title" href="/">Waylon Walker</a>
        <nav class="site-nav site-nav--header site-nav--horizontal"><a href="/archive/">Archive</a></nav>
        <div class="header-controls"><div class="search-container search-container--navbar">
        <div id="pagefind-search"></div></div><div class="header-action-controls"><button>Theme</button></div></div>
        </div></header><main id="main-content"><h1>Archive</h1></main></body></html>`);
      const geometry = () => page.evaluate(() => ({
        header: document.querySelector('.site-header').getBoundingClientRect().height,
        main: document.querySelector('main').getBoundingClientRect().top,
      }));
      const before = await geometry();
      await page.addStyleTag({content: searchCSS});
      assert.deepEqual(await geometry(), before, 'deferred CSS must not shift the header or content');
      for (const markup of [
        '<div class="bleve-search"><input class="bleve-search__input" aria-label="Search"></div>',
        '<div class="pagefind-ui"><form class="pagefind-ui__form"><input class="pagefind-ui__search-input" aria-label="Search"></form></div>',
      ]) {
        await page.locator('#pagefind-search').evaluate((e, html) => {e.innerHTML = html;}, markup);
        assert.deepEqual(await geometry(), before, 'search hydration must not shift the header or content');
      }
    } finally {
      await browser.close();
    }
  });
}
