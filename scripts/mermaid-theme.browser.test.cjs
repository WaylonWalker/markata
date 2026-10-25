const {test, before, after} = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const os = require('node:os');
const path = require('node:path');
const http = require('node:http');
const {execFileSync} = require('node:child_process');
const {chromium} = require(process.env.PLAYWRIGHT_MODULE || 'playwright');

const repoRoot = path.join(__dirname, '..');
const gallery = path.join(repoRoot, 'examples/gallery-site');
const tempRoot = fs.mkdtempSync(path.join(os.tmpdir(), 'markata-mermaid-browser-'));
const siteRoot = path.join(tempRoot, 'gallery-site');
const outputRoot = path.join(siteRoot, 'output');
let server;
let baseURL;

function cssColor(hex) {
  const value = hex.slice(1);
  const parts = [0, 2, 4].map(offset => Number.parseInt(value.slice(offset, offset + 2), 16));
  return `rgb(${parts.join(', ')})`;
}

async function resolvedCssColor(page, variable) {
  return page.evaluate(variable => {
    const probe = document.createElement('span');
    probe.style.backgroundColor = `var(${variable})`;
    document.documentElement.appendChild(probe);
    const color = getComputedStyle(probe).backgroundColor;
    probe.remove();
    const canvas = document.createElement('canvas');
    const context = canvas.getContext('2d');
    context.fillStyle = color;
    context.fillRect(0, 0, 1, 1);
    const [red, green, blue] = context.getImageData(0, 0, 1, 1).data;
    return `rgb(${red}, ${green}, ${blue})`;
  }, variable);
}

before(async () => {
  fs.cpSync(gallery, siteRoot, {recursive: true});
  fs.appendFileSync(path.join(siteRoot, 'markata-go.toml'), '\n[markata-go.icons]\nauto_vendor = false\n');
  execFileSync(process.env.MARKATA_GO_BIN || 'go', [
    'run', './cmd/markata-go', 'build', '-c', path.join(siteRoot, 'markata-go.toml'),
  ], {cwd: repoRoot, stdio: 'inherit'});

  server = http.createServer((request, response) => {
    const pathname = decodeURIComponent(new URL(request.url, 'http://localhost').pathname);
    const filePath = path.resolve(outputRoot, `.${pathname}`);
    if (!filePath.startsWith(outputRoot + path.sep) && filePath !== outputRoot) {
      response.writeHead(403).end();
      return;
    }
    const resolvedPath = fs.existsSync(filePath) && fs.statSync(filePath).isDirectory()
      ? path.join(filePath, 'index.html')
      : filePath;
    if (!fs.existsSync(resolvedPath)) {
      response.writeHead(404).end('Not found');
      return;
    }
    const extension = path.extname(resolvedPath);
    const contentType = extension === '.css' ? 'text/css'
      : extension === '.js' ? 'text/javascript'
      : extension === '.svg' ? 'image/svg+xml'
      : extension === '.json' ? 'application/json'
      : 'text/html';
    response.writeHead(200, {'content-type': contentType});
    fs.createReadStream(resolvedPath).pipe(response);
  });
  await new Promise(resolve => server.listen(0, '127.0.0.1', resolve));
  baseURL = `http://127.0.0.1:${server.address().port}`;
});

after(async () => {
  if (server) await new Promise(resolve => server.close(resolve));
  fs.rmSync(tempRoot, {recursive: true, force: true});
});

test('Mermaid follows active palettes at responsive widths and remains readable without JavaScript', async () => {
  const browser = await chromium.launch({
    headless: true,
    executablePath: process.env.CHROMIUM_EXECUTABLE || undefined,
    args: ['--no-sandbox'],
  });
  try {
    const context = await browser.newContext({
      viewport: {width: 1440, height: 1000},
      colorScheme: 'dark',
      reducedMotion: 'reduce',
    });
    const page = await context.newPage();
    await page.goto(`${baseURL}/rich-markdown-gallery/`);
    await page.addStyleTag({content: '*{animation:none!important;transition:none!important}'});
    await page.locator('.post-content .mermaid svg').waitFor();

    const themes = [
      {name: 'dark', background: '#111827', surface: '#1f2937', text: '#f9fafb', border: '#f59e0b'},
      {name: 'light', background: '#ffffff', surface: '#f3f4f6', text: '#111827', border: '#4b5563'},
      {name: 'custom', background: 'color-mix(in srgb, #291a3d 100%, #ffffff)', surface: 'color-mix(in srgb, #3b2455 100%, #ffffff)', text: 'color-mix(in srgb, #f6edff 100%, #000000)', border: 'color-mix(in srgb, #d98fff 100%, #000000)'},
    ];
    const screenshotDir = process.env.SCREENSHOT_DIR || path.join(repoRoot, 'test-results/mermaid-theme');
    fs.mkdirSync(screenshotDir, {recursive: true});

    for (const width of [320, 360, 412, 768, 1440]) {
      await page.setViewportSize({width, height: 1000});
      for (const theme of themes) {
        await page.evaluate(async ({theme}) => {
          const root = document.documentElement;
          root.dataset.theme = theme.name === 'custom' ? 'dark' : theme.name;
          root.style.setProperty('--color-background', theme.background);
          root.style.setProperty('--color-surface', theme.surface);
          root.style.setProperty('--color-text', theme.text);
          root.style.setProperty('--color-primary', theme.border);
          root.style.setProperty('--color-border', theme.border);
          const eventName = theme.name === 'dark' ? 'theme-change'
            : theme.name === 'light' ? 'color-mode-change' : 'palette-change';
          window.dispatchEvent(new CustomEvent(eventName));
        }, {theme});

        const expected = {
          background: await resolvedCssColor(page, '--color-background'),
          surface: await resolvedCssColor(page, '--color-surface'),
          text: await resolvedCssColor(page, '--color-text'),
          border: await resolvedCssColor(page, '--color-primary'),
        };

        await page.waitForFunction(surface => {
          const node = document.querySelector('.post-content .mermaid svg .node rect, .post-content .mermaid svg .node polygon');
          return node && getComputedStyle(node).fill === surface;
        }, expected.surface);

        const result = await page.evaluate(() => {
          const svg = document.querySelector('.post-content .mermaid svg');
          const node = svg?.querySelector('.node rect, .node polygon');
          const cluster = svg?.querySelector('.cluster rect');
          const edgeLabel = svg?.querySelector('.edgeLabel span');
          return {
            pageWidth: document.documentElement.scrollWidth,
            viewportWidth: innerWidth,
            nodeFill: node && getComputedStyle(node).fill,
            clusterFill: cluster && getComputedStyle(cluster).fill,
            edgeLabelColor: edgeLabel && getComputedStyle(edgeLabel).color,
            border: node && getComputedStyle(node).stroke,
            labels: svg?.textContent || '',
            edges: svg?.querySelectorAll('path.flowchart-link').length || 0,
            reducedMotion: matchMedia('(prefers-reduced-motion: reduce)').matches,
          };
        });
        assert.ok(result.pageWidth <= result.viewportWidth + 1, `page overflow at ${width}px (${theme.name}): ${JSON.stringify(result)}`);
        assert.equal(result.nodeFill, expected.surface);
        assert.equal(result.clusterFill, expected.background);
        assert.equal(result.edgeLabelColor, expected.text);
        assert.equal(result.border, expected.border);
        assert.match(result.labels, /Raw water/);
        assert.match(result.labels, /Ready\?/);
        assert.match(result.labels, /yes/);
        assert.ok(result.edges >= 4, `arrows/edges missing: ${JSON.stringify(result)}`);
        assert.equal(result.reducedMotion, true);

        await page.locator('.post-content .mermaid svg').first().screenshot({
          path: path.join(screenshotDir, `${width}-${theme.name}.png`),
          animations: 'disabled',
        });
      }
    }

    await page.emulateMedia({media: 'print'});
    assert.equal(await page.locator('.post-content .mermaid svg').count(), 1, 'print retains the rendered diagram');
    await context.close();

    const noJS = await browser.newContext({viewport: {width: 320, height: 900}, javaScriptEnabled: false});
    const noJSPage = await noJS.newPage();
    await noJSPage.goto(`${baseURL}/rich-markdown-gallery/`);
    assert.equal(await noJSPage.locator('.post-content pre.mermaid svg').count(), 0);
    assert.match(await noJSPage.locator('.post-content pre.mermaid').innerText(), /Raw water/);
    assert.equal(await noJSPage.locator('.md-code-copy').count(), 0);
    await noJS.close();
  } finally {
    await browser.close();
  }
});
