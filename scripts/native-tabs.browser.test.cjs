const {test} = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const {chromium} = require(process.env.PLAYWRIGHT_MODULE || 'playwright');

const css = fs.readFileSync(path.join(__dirname, '../pkg/themes/default/static/css/components.css'), 'utf8');
const js = fs.readFileSync(path.join(__dirname, '../pkg/themes/default/static/js/md-tabs.js'), 'utf8');

const markup = `<!doctype html><html><head><meta name="viewport" content="width=device-width,initial-scale=1">
<style>${css}</style></head><body><article class="post-content">
<div class="md-tabs" data-md-tabs data-style="boxed">
<div class="md-tabs-nav" role="tablist" aria-label="Install options" hidden>
<button type="button" role="tab" id="md-tab-1-0-tab" aria-controls="md-tab-1-0-panel" aria-selected="true" tabindex="0">curl</button>
<button type="button" role="tab" id="md-tab-1-1-tab" aria-controls="md-tab-1-1-panel" aria-selected="false" tabindex="-1">mise</button>
</div><div class="md-tabs-panels">
<section class="md-tab" role="tabpanel" id="md-tab-1-0-panel" aria-labelledby="md-tab-1-0-tab"><h3 class="md-tab-fallback-label">curl</h3><pre>curl example.com</pre></section>
<section class="md-tab" role="tabpanel" id="md-tab-1-1-panel" aria-labelledby="md-tab-1-1-tab"><h3 class="md-tab-fallback-label">mise</h3><pre>mise use</pre></section>
</div></div></article></body></html>`;

for (const width of [320, 1440]) {
  test(`native tabs keyboard, print and responsive at ${width}px`, async () => {
    const browser = await chromium.launch({headless:true,executablePath:process.env.CHROMIUM_EXECUTABLE || undefined,args:['--no-sandbox']});
    try {
      const page = await browser.newPage({viewport:{width,height:800}});
      await page.setContent(markup);
      // Before hydration the accessible fallback exposes every panel.
      assert.equal(await page.locator('.md-tab:visible').count(), 2);
      assert.equal(await page.getByRole('tab').count(), 0);
      await page.addScriptTag({content:js});
      const tabs = page.getByRole('tab');
      assert.equal(await tabs.count(), 2);
      assert.equal(await page.locator('.md-tab:visible').count(), 1);
      await tabs.first().focus();
      await page.keyboard.press('ArrowRight');
      assert.equal(await tabs.nth(1).getAttribute('aria-selected'), 'true');
      assert.equal(await tabs.nth(1).getAttribute('tabindex'), '0');
      assert.equal(await page.locator('#md-tab-1-0-panel').isHidden(), true);
      assert.equal(await page.locator('#md-tab-1-1-panel').isVisible(), true);
      await page.keyboard.press('Home');
      assert.equal(await tabs.first().getAttribute('aria-selected'), 'true');
      await page.keyboard.press('End');
      assert.equal(await tabs.nth(1).getAttribute('aria-selected'), 'true');
      await tabs.first().click();
      assert.equal(await tabs.first().getAttribute('aria-selected'), 'true');
      assert.equal(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth + 1), true);
      await page.emulateMedia({media:'print'});
      assert.equal(await page.locator('.md-tab:visible').count(), 2);
      assert.equal(await page.locator('.md-tab-fallback-label:visible').count(), 2);
    } finally {
      await browser.close();
    }
  });
}

test('native tabs remain readable without JavaScript', async () => {
  const browser = await chromium.launch({headless:true,executablePath:process.env.CHROMIUM_EXECUTABLE || undefined,args:['--no-sandbox']});
  try {
    const page = await browser.newPage({javaScriptEnabled:false});
    await page.setContent(markup);
    assert.equal(await page.locator('.md-tab:visible').count(), 2);
    assert.equal(await page.locator('.md-tab-fallback-label:visible').count(), 2);
  } finally {
    await browser.close();
  }
});
