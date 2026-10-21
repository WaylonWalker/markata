const {test} = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const {chromium} = require(process.env.PLAYWRIGHT_MODULE || 'playwright');

const root = path.join(__dirname, '../pkg/themes/default/static');
const codeJS = fs.readFileSync(path.join(root, 'js/code-chrome.js'), 'utf8');
const codeCSS = fs.readFileSync(path.join(root, 'css/code.css'), 'utf8');

const markup = `<!doctype html><html><head><meta name="viewport" content="width=device-width,initial-scale=1">
<style>:root{--font-mono:monospace;--color-text:#222;--color-surface:#eee;--color-background:white;--color-border:#bbb;--color-text-muted:#777;--color-primary:blue;--radius-lg:.5rem;--radius:.25rem;--space-2:.5rem;--space-3:.75rem;--space-4:1rem;--text-sm:.9rem;--text-xs:.75rem;}</style>
<style>${codeCSS}</style></head><body>
<article class="post-content" style="width:calc(100% - 2rem);max-width:600px;margin:0 auto">
  <div id="source"><pre><code class="language-bash">kubectl get pods -A\nkubectl get nodes\n</code></pre></div>
  <pre class="mermaid"><code class="language-mermaid">graph TD;A-->B;</code></pre>
</article></body></html>`;

for (const width of [320, 1440]) {
  test(`code copy enhancement is accessible and idempotent at ${width}px`, async () => {
    const browser = await chromium.launch({headless:true,executablePath:process.env.CHROMIUM_EXECUTABLE || undefined,args:['--no-sandbox']});
    try {
      const context = await browser.newContext({viewport:{width,height:900},javaScriptEnabled:true});
      const page = await context.newPage();
      await page.setContent(markup);
      await page.evaluate(() => {
        window.copiedText = null;
        Object.defineProperty(navigator, 'clipboard', {configurable:true,value:{writeText:async text => {window.copiedText=text;}}});
      });
      await page.addScriptTag({content:codeJS});
      assert.equal(await page.locator('.md-code-copy').count(), 1, 'only code blocks, never Mermaid, get controls');
      assert.equal(await page.locator('#source pre').getAttribute('data-language'), 'bash');
      assert.equal(await page.locator('pre.mermaid .md-code-copy').count(), 0);

      const button = page.getByRole('button',{name:'Copy code block'});
      await button.click();
      assert.equal(await page.evaluate(() => window.copiedText), 'kubectl get pods -A\nkubectl get nodes');
      assert.equal(await page.getByRole('button',{name:'Code copied'}).count(), 1);
      const metrics = await page.evaluate(() => ({page:document.documentElement.scrollWidth,viewport:innerWidth}));
      assert.ok(metrics.page <= metrics.viewport+1, 'code toolbar must not overflow narrow article: '+JSON.stringify(metrics));

      await page.evaluate(() => {
        document.querySelector('#source').innerHTML = '<pre><code class="language-go">fmt.Println("hello")</code></pre>';
        window.dispatchEvent(new Event('view-transition-complete'));
        window.markataInitCodeChrome();
      });
      assert.equal(await page.locator('.md-code-copy').count(), 1, 'SPA navigation gets one fresh control');
      assert.equal(await page.locator('#source pre').getAttribute('data-language'), 'go');
      assert.equal(await page.locator('.md-code-copy').getAttribute('data-markata-ui'), 'code-copy');
      await context.close();
    } finally {
      await browser.close();
    }
  });
}

test('code remains readable without JavaScript', async () => {
  const browser = await chromium.launch({headless:true,executablePath:process.env.CHROMIUM_EXECUTABLE || undefined,args:['--no-sandbox']});
  try {
    const context = await browser.newContext({javaScriptEnabled:false});
    const page = await context.newPage();
    await page.setContent(markup);
    assert.equal(await page.locator('.md-code-copy').count(), 0);
    assert.equal(await page.locator('#source pre code').textContent(), 'kubectl get pods -A\nkubectl get nodes\n');
    await context.close();
  } finally {
    await browser.close();
  }
});
