const {test} = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const {chromium} = require(process.env.PLAYWRIGHT_MODULE || 'playwright');

const css = fs.readFileSync(path.join(__dirname, '../pkg/themes/default/static/css/components.css'), 'utf8');
const icon = '<span class="icon" aria-hidden="true" data-icon="lucide/rocket"><svg viewBox="0 0 24 24" width="1em" height="1em"><path d="M2 12h20" stroke="currentColor"/></svg></span>';
const html = `<!doctype html><html><head><meta name="viewport" content="width=device-width, initial-scale=1">
<style>:root{--color-text:#ddd;--color-text-muted:#888;--color-surface:#282d33;--color-background:#191c20;--space-3:.75rem;--space-4:1rem;--radius-lg:.75rem}
body{margin:0;background:var(--color-background);color:var(--color-text);font:16px system-ui}
.post-content{width:calc(100% - 2rem);max-width:680px;margin:0 auto}
</style><style>${css}</style></head><body><article class="post-content">
<div class="md-feature" id="feature">
  <div class="md-feature-icon"><p>${icon}</p></div>
  <div class="md-feature-body"><h3>1. Verify production</h3><p>Confirm the live image and JavaScript-free article rendering.</p></div>
</div>
<div class="md-columns">
  <div class="md-card" data-align="center"><div class="md-icon" data-size="lg"><p>${icon}</p></div><h3>Falcon2</h3><p>K3s node</p></div>
  <div class="md-card" data-align="center"><div class="md-icon" data-size="lg"><p>${icon}</p></div><h3>Falcon3</h3><p>NAS + ZFS</p></div>
</div>
</article></body></html>`;

for (const width of [320, 412, 1440]) {
  test(`icon-first Markdown feature row and centered cards fit ${width}px without JS`, async () => {
    const browser = await chromium.launch({headless:true, executablePath:process.env.CHROMIUM_EXECUTABLE || undefined,args:['--no-sandbox']});
    try {
      const context = await browser.newContext({viewport:{width,height:900},javaScriptEnabled:false});
      const page = await context.newPage();
      await page.setContent(html);
      assert.equal(await page.getByText('1. Verify production').count(), 1);
      assert.equal(await page.locator('.md-feature-icon svg').count(), 1);
      assert.equal(await page.locator('.md-card .md-icon svg').count(), 2);
      const metrics = await page.evaluate(() => {
        const feature = document.querySelector('#feature').getBoundingClientRect();
        const tile = document.querySelector('.md-feature-icon').getBoundingClientRect();
        const body = document.querySelector('.md-feature-body').getBoundingClientRect();
        const icon = document.querySelector('.md-icon').getBoundingClientRect();
        const cards = Array.from(document.querySelectorAll('.md-card')).map(el => el.getBoundingClientRect());
        return {viewport:innerWidth,scroll:document.documentElement.scrollWidth,
          feature:{x:feature.x,right:feature.right},tile:{x:tile.x,right:tile.right},body:{x:body.x,right:body.right},
          iconSize:icon.width,cards:cards.map(x=>({x:x.x,y:x.y,bottom:x.bottom}))};
      });
      assert.ok(metrics.scroll <= metrics.viewport+1,'icon layouts must not scroll the page: '+JSON.stringify(metrics));
      assert.ok(metrics.body.x >= metrics.tile.right-1,'icon tile and text must not overlap');
      assert.ok(metrics.body.right <= metrics.feature.right+1,'body must fit in feature row');
      assert.ok(metrics.iconSize > 20,'large icon must have enough visual space');
      if(width<=412) assert.ok(metrics.cards[1].y >= metrics.cards[0].bottom-1,'cards should stack on phones');
      else assert.ok(metrics.cards[1].x > metrics.cards[0].x,'desktop cards should sit beside each other');
      await context.close();
    } finally {
      await browser.close();
    }
  });
}
