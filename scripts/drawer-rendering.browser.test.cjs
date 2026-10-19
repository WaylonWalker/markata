const {test} = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const {chromium} = require(process.env.PLAYWRIGHT_MODULE || 'playwright');
const root = path.join(__dirname, '../pkg/themes/default/static');
const css = ['variables', 'main', 'components']
  .map(name => fs.readFileSync(path.join(root, `css/${name}.css`), 'utf8')).join('\n');
const script = fs.readFileSync(path.join(root, 'js/view-transitions.js'), 'utf8');

for (const width of [412, 1440]) {
  for (const reducedMotion of ['no-preference', 'reduce']) {
    test(`drawer rendering and focus at ${width}px, motion ${reducedMotion}`, async () => {
      const browser = await chromium.launch({headless: true,
        executablePath: process.env.CHROMIUM_EXECUTABLE || undefined, args: ['--no-sandbox']});
      try {
        const page = await browser.newPage({viewport: {width, height: 823}, reducedMotion});
        await page.route('http://drawer.test/**', route => route.fulfill({contentType: 'text/html', body: `
          <!doctype html><html><head><style>${css}</style></head><body>
          <header class="site-header" style="height:70px">Navigation</header>
          <div class="page-wrapper"><aside class="feed-sidebar feed-sidebar--left">
          <button class="sidebar-toggle" data-sidebar-label="feeds">Feeds</button>
          <nav class="feed-nav"><div id="feed-nav-collapsible"><ul class="feed-nav-list">
          ${Array.from({length: 50}, (_, i) => `<li class="feed-nav-item ${i === 35 ? 'feed-nav-item--active' : ''}"><a href="/feed-${i}/">Feed ${i}</a></li>`).join('')}
          </ul></div></nav></aside><main class="main-content" style="height:3000px">Article</main>
          <aside class="doc-sidebar doc-sidebar--right">
          <button class="sidebar-toggle" data-sidebar-label="contents">Contents</button>
          <nav class="toc"><a href="#heading">Heading</a></nav></aside></div></body></html>`}));
        await page.goto('http://drawer.test/');
        await page.addScriptTag({content: script});
        assert.equal(await page.evaluate(() => CSS.supports('transition-behavior', 'allow-discrete')), true);
        for (const selector of ['.feed-sidebar', '.doc-sidebar']) {
          const drawer = page.locator(selector);
          const contents = drawer.locator(':scope > nav');
          const toggle = drawer.locator(':scope > .sidebar-toggle');
          assert.equal(await contents.evaluate(el => getComputedStyle(el).contentVisibility), 'hidden');
          await toggle.focus();
          await page.keyboard.press('Enter');
          await page.waitForFunction(selector => {
            const el = document.querySelector(selector);
            return el.classList.contains('sidebar--pinned') && getComputedStyle(el.querySelector('nav')).contentVisibility === 'visible';
          }, selector);
          if (width <= 1200) {
            await page.waitForFunction(selector => document.querySelector(selector + ' nav').contains(document.activeElement), selector);
          }
          await drawer.evaluate(el => Promise.all(el.getAnimations({subtree: true}).map(animation => animation.finished)));
          // Sample during the closing animation without relying on wall-clock sleeps.
          const during = await drawer.evaluate(el => {
            el.querySelector('.sidebar-toggle').click();
            const animations = el.getAnimations({subtree: true});
            for (const animation of animations) {
              animation.pause();
              animation.currentTime = Number(animation.effect.getTiming().duration) / 2;
            }
            const result = {visibility: getComputedStyle(el.querySelector('nav')).contentVisibility,
              animations: animations.length};
            for (const animation of animations) animation.finish();
            return result;
          });
          if (reducedMotion === 'no-preference') {
            assert.ok(during.animations > 0, 'closing must animate');
            assert.equal(during.visibility, 'visible', 'contents must survive the closing slide');
          } else {
            assert.equal(during.animations, 0, 'reduced motion must disable content transitions');
          }
          await page.waitForFunction(selector => getComputedStyle(document.querySelector(selector + ' nav')).contentVisibility === 'hidden', selector);
          if (width <= 1200) assert.equal(await toggle.evaluate(el => el === document.activeElement), true);
        }
        if (width > 1200) {
          await page.locator('.feed-sidebar > .sidebar-toggle').click();
          await page.reload();
          await page.addScriptTag({content: script});
          assert.equal(await page.locator('.feed-sidebar > nav').evaluate(el => getComputedStyle(el).contentVisibility), 'visible', 'saved-open drawer must render immediately');
        }
      } finally {
        await browser.close();
      }
    });
  }
}
