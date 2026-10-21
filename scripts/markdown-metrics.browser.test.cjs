// Static-first article metrics: responsive, semantic, no-JS browser coverage.
const {test} = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const {chromium} = require(process.env.PLAYWRIGHT_MODULE || 'playwright');

const cssRoot = path.join(__dirname, '../pkg/themes/default/static/css');
const styles = ['variables', 'main', 'components']
  .map(name => fs.readFileSync(path.join(cssRoot, name + '.css'), 'utf8'))
  .join('\n');

const markup = `<!doctype html><html><head><meta name="viewport" content="width=device-width, initial-scale=1">
<style>${styles}</style></head><body>
<article class="post-content" style="margin:0 auto; width:calc(100% - 2rem); max-width:960px">
<div class="md-status" data-tone="success"><p>NVMe cleanup verified</p></div>
<div class="md-layout">
 <p>Falcon3 NVMe — before and after</p>
 <div class="md-columns" id="metrics">
  <div class="md-stat" data-tone="caution"><p>Before</p><p><strong>83 GiB</strong></p><p>18% free</p>
   <div class="md-meter" data-value="18" data-max="100" data-label="Before: 18 percent free" data-tone="caution">
    <meter value="18" max="100" aria-label="Before: 18 percent free">18 / 100</meter></div>
  </div>
  <div class="md-stat" data-tone="success"><p>After</p><p><strong>212 GiB</strong></p><p>45% free</p>
   <div class="md-meter" data-value="45" data-max="100" data-label="After: 45 percent free" data-tone="success">
    <meter value="45" max="100" aria-label="After: 45 percent free">45 / 100</meter></div>
  </div>
 </div>
 <div class="md-delta" data-tone="success"><p>Space recovered</p><p><strong>+129 GiB</strong></p></div>
</div>
<table><thead><tr><th>Component</th><th>Before</th><th>After</th></tr></thead>
<tbody><tr><td>Snapshots</td><td>158 GiB</td><td>76 GiB</td></tr></tbody></table>
<div class="md-progress" data-value="7" data-max="10" data-label="Build steps completed">
<progress value="7" max="10" aria-label="Build steps completed">7 / 10</progress><p>7 of 10 steps complete</p></div>
</article></body></html>`;

for (const width of [320, 412, 1440]) {
  test(`static Markdown gauges and before/after cards at ${width}px without JavaScript`, async () => {
    const browser = await chromium.launch({headless: true, executablePath: process.env.CHROMIUM_EXECUTABLE || undefined, args: ['--no-sandbox']});
    try {
      const context = await browser.newContext({viewport:{width,height:900}, javaScriptEnabled:false, colorScheme:'dark'});
      const page = await context.newPage();
      await page.setContent(markup);
      assert.equal(await page.getByText('NVMe cleanup verified').count(), 1);
      assert.equal(await page.locator('meter').count(), 2);
      assert.equal(await page.locator('progress').count(), 1);
      assert.equal(await page.locator('meter').nth(1).evaluate(e => e.value), 45);
      assert.equal(await page.locator('progress').evaluate(e => e.value), 7);
      assert.equal(await page.locator('meter').nth(1).getAttribute('aria-label'), 'After: 45 percent free');

      const layout = await page.evaluate(() => {
        const [first, second] = document.querySelectorAll('#metrics > .md-stat');
        const a = first.getBoundingClientRect();
        const b = second.getBoundingClientRect();
        const meter = document.querySelector('meter');
        return {a:{x:a.x,y:a.y,bottom:a.bottom},b:{x:b.x,y:b.y},
          viewport:innerWidth,scroll:document.documentElement.scrollWidth,
          meterWidth:meter.getBoundingClientRect().width};
      });
      assert.ok(layout.scroll <= layout.viewport + 1, 'metrics must not cause horizontal page scroll: ' + JSON.stringify(layout));
      assert.ok(layout.meterWidth > 16, 'native meter must be visible');
      if (width < 500) assert.ok(layout.b.y >= layout.a.bottom - 1, 'narrow metrics must stack: ' + JSON.stringify(layout));
      else assert.ok(layout.b.x > layout.a.x + 10, 'wide metrics must be side by side: ' + JSON.stringify(layout));
      await context.close();
    } finally {
      await browser.close();
    }
  });
}
