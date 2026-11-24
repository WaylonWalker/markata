const {test, before, after} = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const os = require('node:os');
const path = require('node:path');
const http = require('node:http');
const {execFileSync} = require('node:child_process');
const {chromium} = require(process.env.PLAYWRIGHT_MODULE || 'playwright');

const repoRoot = path.join(__dirname, '..');
const siteRoot = fs.mkdtempSync(path.join(os.tmpdir(), 'markata-native-mermaid-theme-'));
const gallery = path.join(repoRoot, 'examples/gallery-site');
let server;
let baseURL;

before(async () => {
  fs.cpSync(gallery, siteRoot, {recursive: true});
  fs.appendFileSync(path.join(siteRoot, 'markata-go.toml'),
    '\n[markata-go.icons]\nauto_vendor = false\n[markata-go.mermaid]\nmode = "go"\n');
  execFileSync(process.env.MARKATA_GO_BIN || 'go', [
    'run', './cmd/markata-go', 'build', '-c', path.join(siteRoot, 'markata-go.toml'),
  ], {cwd: repoRoot, stdio: 'inherit'});

  const root = path.join(siteRoot, 'output');
  server = http.createServer((request, response) => {
    const pathname = decodeURIComponent(new URL(request.url, 'http://localhost').pathname);
    const filePath = path.resolve(root, '.' + pathname);
    if (!filePath.startsWith(root + path.sep) && filePath !== root) {
      response.writeHead(403).end();
      return;
    }
    const resolved = fs.existsSync(filePath) && fs.statSync(filePath).isDirectory()
      ? path.join(filePath, 'index.html')
      : filePath;
    if (!fs.existsSync(resolved)) {
      response.writeHead(404).end('Not found');
      return;
    }
    const extension = path.extname(resolved);
    const contentType = extension === '.css' ? 'text/css'
      : extension === '.js' ? 'text/javascript'
      : extension === '.svg' ? 'image/svg+xml'
      : 'text/html';
    response.writeHead(200, {'content-type': contentType});
    fs.createReadStream(resolved).pipe(response);
  });
  await new Promise(resolve => server.listen(0, '127.0.0.1', resolve));
  baseURL = 'http://127.0.0.1:' + server.address().port;
});

after(async () => {
  if (server) await new Promise(resolve => server.close(resolve));
  fs.rmSync(siteRoot, {recursive: true, force: true});
});

test('native Mermaid SVG responds to live palette updates without Mermaid.js', async () => {
  const browser = await chromium.launch({
    headless: true,
    executablePath: process.env.CHROMIUM_EXECUTABLE || undefined,
    args: ['--no-sandbox'],
  });
  try {
    const page = await browser.newPage({viewport: {width: 360, height: 900}});
    await page.goto(baseURL + '/rich-markdown-gallery/');
    const svg = page.locator('.post-content pre.mermaid svg').first();
    await svg.waitFor();
    assert.match(await svg.textContent(), /Raw water/);

    const paint = svg.locator('[fill*="var(--color-"]').first();
    assert.ok(await paint.count(), 'native Go SVG should use palette-dependent fill colors');
    const variable = await paint.getAttribute('fill');
    const match = variable.match(/var\((--color-[\w-]+),/);
    assert.ok(match, 'expected CSS variable with static fallback: ' + variable);
    const colorVariable = match[1];

    for (const [hex, expected] of [
      ['#ff0000', 'rgb(255, 0, 0)'],
      ['#00ff00', 'rgb(0, 255, 0)'],
      ['#0000ff', 'rgb(0, 0, 255)'],
    ]) {
      await page.evaluate(([name, value]) => {
        document.documentElement.style.setProperty(name, value);
        dispatchEvent(new Event('palette-change'));
      }, [colorVariable, hex]);
      assert.equal(await paint.evaluate(element => getComputedStyle(element).fill), expected,
        'native SVG must react to palette changes without re-rendering');
    }
    assert.ok(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth + 1),
      'native SVG must not overflow narrow screens');
    await page.close();

    const noJS = await browser.newPage({javaScriptEnabled: false, viewport: {width: 320, height: 900}});
    await noJS.goto(baseURL + '/rich-markdown-gallery/');
    assert.equal(await noJS.locator('.post-content pre.mermaid svg').count(), 1,
      'native diagram must be visible without JavaScript');
    assert.match(await noJS.locator('.post-content pre.mermaid svg').textContent(), /Raw water/);
    await noJS.close();
  } finally {
    await browser.close();
  }
});
