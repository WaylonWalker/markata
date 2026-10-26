// Card permalinks must not manufacture scrollbars; real wide content still scrolls.
const {test} = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const {chromium} = require(process.env.PLAYWRIGHT_MODULE || 'playwright');

const cssRoot = path.join(__dirname, '../pkg/themes/default/static/css');
const styles = ['variables', 'main', 'components', 'code']
  .map(name => fs.readFileSync(path.join(cssRoot, `${name}.css`), 'utf8'))
  .join('\n');

function heading(level, id, title) {
  return `<h${level} id="${id}">${title}<a class="heading-anchor" href="#${id}" aria-label="Link to ${title}">#</a></h${level}>`;
}

const markup = `<!doctype html><html><head>
<meta name="viewport" content="width=device-width, initial-scale=1">
<style>${styles}</style></head><body>
<article class="post-content" style="margin:0 auto;width:calc(100% - 2rem);max-width:960px">
<div style="padding-right:2rem">${heading(2, 'outside', 'Ordinary article heading')}</div>
<div class="md-grid">
 <div class="md-card">${heading(4, 'short', 'Title')}<p>Short card</p></div>
 <div class="md-card">${heading(4, 'long', 'A longer card heading that wraps across several lines on narrow screens')}
  <div class="md-stack">${heading(5, 'nested', 'Nested heading')}</div>
 </div>
 <div class="md-card">${heading(4, 'code', 'Wide code')}<pre><code>${'wide_code_'.repeat(40)}</code></pre></div>
</div>
<div class="md-grid"><div class="md-card" id="table-card">${heading(4, 'table', 'Wide comparison table')}
<table><thead><tr><th>Component</th><th>Before</th><th>After</th></tr></thead>
<tbody><tr>${['Component', 'Before', 'After'].map(text => `<td>${text.repeat(30)}</td>`).join('')}</tr></tbody></table>
</div></div>
</article></body></html>`;

for (const width of [390, 1024, 1440, 1700]) {
  test(`card heading anchors and wide content at ${width}px without JavaScript`, async () => {
    const browser = await chromium.launch({
      headless: true,
      executablePath: process.env.CHROMIUM_EXECUTABLE || undefined,
      args: ['--no-sandbox'],
    });
    try {
      const context = await browser.newContext({
        viewport: {width, height: 900}, javaScriptEnabled: false, reducedMotion: 'reduce',
      });
      const page = await context.newPage();
      await page.setContent(markup);
      const cards = await page.locator('.md-card').evaluateAll(elements => elements.map(element => ({
        client: element.clientWidth, scroll: element.scrollWidth,
      })));
      for (const card of cards.slice(0, 3)) {
        assert.equal(card.scroll, card.client, `heading/code card must not scroll: ${JSON.stringify(cards)}`);
      }
      // The existing direct-table contract deliberately scrolls the card.
      assert.ok(cards[3].scroll > cards[3].client, 'wide table must still scroll within its card');
      assert.equal(await page.locator('pre').evaluate(element => {
        element.scrollLeft = 100;
        return element.scrollWidth > element.clientWidth && element.scrollLeft > 0;
      }), true, 'wide code must scroll in its own element');
      assert.ok(await page.locator('#table-card').evaluate(element => {
        element.scrollLeft = 100;
        return element.scrollLeft > 0;
      }), 'table scroll container must remain usable');
      await page.locator('#table-card').evaluate(element => { element.scrollLeft = 0; });

      const anchors = page.locator('.md-card .heading-anchor');
      for (let index = 0; index < await anchors.count(); index++) {
        const anchor = anchors.nth(index);
        await anchor.focus();
        // Poll from the driver: waitForFunction's injected polling script does
        // not run when page JavaScript is disabled.
        for (let attempt = 0; attempt < 20; attempt++) {
          if (await anchor.evaluate(element => getComputedStyle(element).opacity) === '1') break;
          await page.waitForTimeout(50);
        }
        assert.equal(await anchor.evaluate(element => getComputedStyle(element).opacity), '1');
        assert.equal(await anchor.evaluate(element => {
          const link = element.getBoundingClientRect();
          const card = element.closest('.md-card').getBoundingClientRect();
          const heading = element.parentElement.getBoundingClientRect();
          const style = getComputedStyle(element);
          return document.activeElement === element && link.left >= card.left &&
            link.right <= card.right && link.right <= heading.right && style.outlineStyle !== 'none';
        }), true, 'keyboard-focused permalink must be visible inside the card');
      }
      assert.equal(await page.locator('#outside .heading-anchor').evaluate(element => getComputedStyle(element).position), 'absolute');
      assert.ok(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth), 'cards must not widen the page');
      console.log(`${width}px cards (client/scroll): ${cards.map(card => `${card.client}/${card.scroll}`).join(', ')}`);
      await context.close();
    } finally {
      await browser.close();
    }
  });
}
