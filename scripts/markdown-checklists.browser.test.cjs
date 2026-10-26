const {test} = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const {chromium} = require(process.env.PLAYWRIGHT_MODULE || 'playwright');

const code = fs.readFileSync(path.join(__dirname, '../pkg/themes/default/static/js/markdown-checklists.js'), 'utf8');
const markup = `<!doctype html><html><head><meta name="viewport" content="width=device-width,initial-scale=1"></head><body>
<article class="post-content">
<div class="md-checklist" data-id="release">
<ul class="task-list">
<li class="task-list-item"><input type="checkbox" disabled> Verify CI</li>
<li class="task-list-item"><input type="checkbox" checked disabled> Review diff</li>
</ul>
</div>
<ul class="task-list"><li class="task-list-item"><input type="checkbox" disabled> Read only outside container</li></ul>
</article></body></html>`;

for (const width of [320, 1440]) {
  test(`opt-in checklists persist, remain scoped and survive SPA swaps at ${width}px`, async () => {
    const browser = await chromium.launch({headless:true,executablePath:process.env.CHROMIUM_EXECUTABLE || undefined,args:['--no-sandbox']});
    try {
      const context = await browser.newContext({viewport:{width,height:900}});
      const page = await context.newPage();
      await page.route('https://markata.test/**', route => route.fulfill({status:200,contentType:'text/html',body:markup}));
      await page.goto('https://markata.test/tutorial/');
      await page.addScriptTag({content:code});
      const checkboxes = page.locator('.md-checklist input[type="checkbox"]');
      assert.equal(await checkboxes.count(), 2);
      assert.equal(await checkboxes.first().isEnabled(), true);
      assert.equal(await checkboxes.nth(1).isChecked(), true);
      assert.equal(await page.locator('.post-content > ul input').isDisabled(), true);
      await checkboxes.first().check();
      assert.equal(await page.evaluate(() => JSON.parse(localStorage.getItem('markata:checklist:v1:/tutorial/:release'))[0]), true);
      await page.reload();
      await page.addScriptTag({content:code});
      assert.equal(await page.locator('.md-checklist input').first().isChecked(), true);
      await page.evaluate(() => {
        const node = document.querySelector('.md-checklist');
        node.innerHTML = '<ul class="task-list"><li class="task-list-item"><input type="checkbox" disabled> New task</li></ul>';
        window.dispatchEvent(new Event('view-transition-complete'));
      });
      assert.equal(await page.locator('.md-checklist input').isEnabled(), true);
      assert.equal(await page.locator('.md-checklist input').count(), 1);
      const sizes = await page.evaluate(() => ({page:document.documentElement.scrollWidth,viewport:innerWidth}));
      assert.ok(sizes.page <= sizes.viewport+1, 'checklist should not overflow: '+JSON.stringify(sizes));
      await context.close();
    } finally {
      await browser.close();
    }
  });
}

test('no JS retains readable disabled checkboxes', async () => {
  const browser = await chromium.launch({headless:true,executablePath:process.env.CHROMIUM_EXECUTABLE || undefined,args:['--no-sandbox']});
  try {
    const context = await browser.newContext({javaScriptEnabled:false});
    const page = await context.newPage();
    await page.setContent(markup);
    assert.equal(await page.locator('.md-checklist input').first().isDisabled(), true);
    assert.equal(await page.getByText('Verify CI').count(), 1);
    await context.close();
  } finally {
    await browser.close();
  }
});
