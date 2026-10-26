const {test} = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const {chromium} = require(process.env.PLAYWRIGHT_MODULE || 'playwright');
const root = path.join(__dirname, '../pkg/themes/default/static');
const css = ['variables', 'main', 'components'].map(n => fs.readFileSync(path.join(root, `css/${n}.css`), 'utf8')).join('\n');
const script = fs.readFileSync(path.join(root, 'js/view-transitions.js'), 'utf8');
for (const width of [412, 1440]) {
  test(`drawers defer closed layout and preserve opening at ${width}px`, async () => {
    const browser = await chromium.launch({headless: true, executablePath: process.env.CHROMIUM_EXECUTABLE || undefined, args: ['--no-sandbox']});
    try {
      const page = await browser.newPage({viewport: {width, height: 823}});
      if (width > 1200) await page.addInitScript(() => localStorage.setItem('markata-sidebar:left', 'open'));
      await page.route('http://archive.test/**', route => route.fulfill({contentType: 'text/html', body: `<!doctype html><html><head><style>${css}</style></head><body>
        <header class="site-header" style="height:100px">Header</header><div class="page-wrapper" style="--feed-sidebar-width:250px">
        <aside class="feed-sidebar feed-sidebar--left"><button class="sidebar-toggle" data-sidebar-label="feeds">Feeds</button><nav class="feed-nav">
        <div class="feed-nav-header" role="button" aria-controls="feed-nav-collapsible">Feeds</div><div id="feed-nav-collapsible"><ul class="feed-nav-list">
        ${Array.from({length:50},(_,i)=>`<li class="feed-nav-item ${i===35?'feed-nav-item--active':''}"><a href="/feed-${i}/">Feed ${i}</a></li>`).join('')}
        </ul></div></nav></aside><main class="main-content" style="height:3000px">Stories</main></div></body></html>`}));
      await page.goto('http://archive.test/');
      await page.evaluate(() => {
        window.layoutReads = {header: 0, active: 0};
        const rect = Element.prototype.getBoundingClientRect;
        Element.prototype.getBoundingClientRect = function() {
          if (this.matches('.site-header')) window.layoutReads.header++;
          return rect.call(this);
        };
        const offset = Object.getOwnPropertyDescriptor(HTMLElement.prototype, 'offsetTop');
        Object.defineProperty(HTMLElement.prototype, 'offsetTop', {get() {
          if (this.matches('.feed-nav-item--active')) window.layoutReads.active++;
          return offset.get.call(this);
        }});
      });
      await page.addScriptTag({content: script});
      if (width <= 1200) {
        assert.deepEqual(await page.evaluate(() => window.layoutReads), {header: 0, active: 0}, 'closed overlay startup must not read geometry');
        await page.locator('.sidebar-toggle').click();
      }
      await page.waitForFunction(() => document.querySelector('.feed-sidebar').classList.contains('sidebar--pinned'));
      const geometry = await page.evaluate(() => {
        const drawer = document.querySelector('.feed-sidebar');
        const active = document.querySelector('.feed-nav-item--active');
        const container = document.getElementById('feed-nav-collapsible');
        const a = active.getBoundingClientRect(), c = container.getBoundingClientRect();
        return {top: drawer.getBoundingClientRect().top, header: document.querySelector('.site-header').getBoundingClientRect().bottom,
          rootOffset: document.documentElement.style.getPropertyValue('--sidebar-top'), drawerOffset: drawer.style.getPropertyValue('--sidebar-top'),
          activeTop:a.top,activeBottom:a.bottom,containerTop:c.top,containerBottom:c.bottom,scroll:container.scrollTop};
      });
      assert.ok(Math.abs(geometry.top-(width<=1200 ? 0 : geometry.header))<=1, JSON.stringify(geometry));
      assert.equal(geometry.rootOffset, '', 'drawer offset must not invalidate the document root');
      assert.equal(geometry.drawerOffset, width<=1200 ? '' : '100px');
      assert.ok(geometry.scroll > 0 && geometry.activeTop>=geometry.containerTop && geometry.activeBottom<=geometry.containerBottom, 'active feed must be visible: '+JSON.stringify(geometry));
      if (width <= 1200) {
        await page.waitForFunction(() => document.activeElement.matches('.feed-nav-item--active a'));
        await page.keyboard.press('Escape');
        assert.equal(await page.locator('.sidebar-toggle').getAttribute('aria-expanded'), 'false');
        assert.equal(await page.locator('.sidebar-toggle').evaluate(e=>document.activeElement===e), true);
        await page.evaluate(() => window.scrollTo(0,600));
        await page.locator('.sidebar-toggle').click();
        await page.waitForFunction(() => document.querySelector('.feed-sidebar').classList.contains('sidebar--pinned'));
        assert.equal(await page.locator('.feed-sidebar').evaluate(e=>e.getBoundingClientRect().top), 0);
      }
    } finally { await browser.close(); }
  });
}
