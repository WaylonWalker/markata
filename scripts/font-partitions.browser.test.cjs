// Verify real bundled font transfers and shaping, not only manifest metadata.
const {test} = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const http = require('node:http');
const path = require('node:path');
const {chromium} = require(process.env.PLAYWRIGHT_MODULE || 'playwright');
const yaml = require(process.env.YAML_MODULE || 'yaml');
const root = path.join(__dirname, '../internal/fontcatalog/source-sans-3');
const manifest = yaml.parse(fs.readFileSync(path.join(root, 'manifest.yaml'), 'utf8'));
const css = ['prose-basic', 'prose-extra'].map(tier => {
  const t = manifest.tiers[tier];
  return `@font-face{font-family:Partitioned;src:url('/${t.file}');font-weight:200 900;unicode-range:${t.unicode_range.join(',')};}`;
}).join('') + `@font-face{font-family:Original;src:url('/${manifest.tiers['prose-core'].file}');font-weight:200 900;}`;

test('prose partitions defer extra glyphs and retain bounded text geometry', async () => {
  const requests = [];
  const files = new Set(Object.values(manifest.tiers).map(t => '/' + t.file));
  const server = http.createServer((req, res) => {
    if (req.url === '/') {
      res.setHeader('Content-Type', 'text/html; charset=utf-8');
      return res.end(`<meta charset="utf-8"><style>${css}#sample{font:32px Partitioned}</style><div id="sample">AVATAR façade curly “quotes” and ↗ arrows</div>`);
    }
    if (files.has(req.url)) {
      requests.push(req.url);
      res.setHeader('Content-Type', 'font/woff2');
      return res.end(fs.readFileSync(path.join(root, req.url.slice(1))));
    }
    res.writeHead(404);
    res.end();
  });
  await new Promise(resolve => server.listen(0, '127.0.0.1', resolve));
  let browser;
  try {
    browser = await chromium.launch({headless:true, executablePath:process.env.CHROMIUM_EXECUTABLE || undefined, args:['--no-sandbox']});
    const page = await browser.newPage();
    await page.goto(`http://127.0.0.1:${server.address().port}`);
    await page.evaluate(() => document.fonts.ready);
    assert.deepEqual(requests, ['/' + manifest.tiers['prose-basic'].file], 'common text must fetch only basic');
    await page.locator('#sample').evaluate(e => {e.textContent='Ā Č Ǎ Việt Nam Ȁ € ™';});
    await page.evaluate(() => document.fonts.ready);
    await page.waitForLoadState('networkidle');
    assert.ok(requests.includes('/' + manifest.tiers['prose-extra'].file), 'extended text must fetch extra');
    assert.ok(!requests.includes('/' + manifest.tiers['prose-core'].file), 'the broad face must remain idle');
    const metrics = await page.evaluate(async () => {
      const texts = ['AVATAR façade curly “quotes” ↗ arrows', 'office affinity fi fl ffi', 'ĀV Č Ǎ Việt Nam Ȁ € ™'];
      const ctx = document.createElement('canvas').getContext('2d');
      const result = [];
      for (const weight of [200,400,550,700,900]) for (const size of [16,24,32]) for (const text of texts) {
        await document.fonts.load(`${weight} ${size}px Original`, text);
        await document.fonts.load(`${weight} ${size}px Partitioned`, text);
        ctx.font = `${weight} ${size}px Original`;
        const before = ctx.measureText(text).width;
        ctx.font = `${weight} ${size}px Partitioned`;
        result.push({weight,size,text,before,delta:ctx.measureText(text).width-before});
      }
      return result;
    });
    for (const metric of metrics) {
      const delta = Math.abs(metric.delta);
      if (metric.text.startsWith('Ā')) {
        // Separate font faces cannot kern across the Unicode boundary.
        // Bound that documented tradeoff; retain exact common-text checks.
        assert.ok(delta < 1 && delta / metric.before < .004, JSON.stringify(metric));
      } else {
        assert.ok(delta < .01, 'common text/ligature width changed: ' + JSON.stringify(metric));
      }
    }
  } finally {
    await browser?.close();
    await new Promise(resolve => server.close(resolve));
  }
});
