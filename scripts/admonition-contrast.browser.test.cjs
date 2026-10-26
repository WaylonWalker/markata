// Regression checks for generic container classes sharing names with admonitions.
const {test} = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const {chromium} = require(process.env.PLAYWRIGHT_MODULE || 'playwright');

const staticRoot = path.join(__dirname, '../pkg/themes/default/static');
const componentsCSS = fs.readFileSync(path.join(staticRoot, 'css/components.css'), 'utf8');
const admonitionsCSS = fs.readFileSync(path.join(staticRoot, 'css/admonitions.css'), 'utf8');

function luminance(rgb) {
  const channels = rgb.match(/[\d.]+/g).slice(0, 3).map(Number).map(value => {
    const normalized = value / 255;
    return normalized <= 0.04045 ? normalized / 12.92 : ((normalized + 0.055) / 1.055) ** 2.4;
  });
  return channels[0] * 0.2126 + channels[1] * 0.7152 + channels[2] * 0.0722;
}

function contrastRatio(foreground, background) {
  const values = [luminance(foreground), luminance(background)].sort((a, b) => b - a);
  return (values[0] + 0.05) / (values[1] + 0.05);
}

test('palette admonitions keep readable surfaces without changing generic containers', async () => {
  const browser = await chromium.launch({headless: true, args: ['--no-sandbox']});
  try {
    const page = await browser.newPage({colorScheme: 'light'});
    await page.setContent(`<!doctype html>
      <html data-palette="ayu-dark"><head><style>
        :root {
          --color-text: #bfbdb6;
          --admonition-warn-bg: #131721;
          --admonition-warn-border: #ff8f40;
          --admonition-note-bg: #131721;
          --admonition-note-border: #59c2ff;
        }
        body { color: var(--color-text); }
      </style><style>${componentsCSS}</style><style>${admonitionsCSS}</style></head>
      <body>
        <div id="warning" class="admonition warning"><p class="admonition-title">Warning</p><p>Warning text</p></div>
        <div id="note" class="admonition note"><p class="admonition-title">Note</p><p>Note text</p></div>
        <div id="generic-warning" class="warning">Generic warning</div>
        <div id="generic-note" class="note">Generic note</div>
        <p id="generic-error" class="error">Generic error</p>
      </body></html>`);

    for (const [id, label] of [['warning', 'warning'], ['note', 'note']]) {
      const colors = await page.locator(`#${id}`).evaluate(element => {
        const body = element.querySelector('p:not(.admonition-title)');
        const title = element.querySelector('.admonition-title');
        return {
          foreground: getComputedStyle(body).color,
          titleForeground: getComputedStyle(title).color,
          surface: getComputedStyle(element).backgroundColor,
        };
      });
      const ratio = contrastRatio(colors.foreground, colors.surface);
      const titleRatio = contrastRatio(colors.titleForeground, colors.surface);
      assert.ok(ratio >= 4.5, `${label} admonition body contrast is ${ratio.toFixed(2)}:1 (${JSON.stringify(colors)})`);
      assert.ok(titleRatio >= 4.5, `${label} admonition title contrast is ${titleRatio.toFixed(2)}:1 (${JSON.stringify(colors)})`);
      assert.equal(colors.surface, 'rgb(19, 23, 33)', `${label} admonition must use its palette surface`);
    }

    assert.equal(await page.locator('#generic-warning').evaluate(element => getComputedStyle(element).backgroundColor), 'rgb(255, 243, 205)');
    assert.equal(await page.locator('#generic-note').evaluate(element => getComputedStyle(element).backgroundColor), 'rgb(231, 233, 235)');
    assert.equal(await page.locator('#generic-error').evaluate(element => getComputedStyle(element).backgroundColor), 'rgb(248, 215, 218)');
    assert.equal(await page.locator('#generic-error').evaluate(element => getComputedStyle(element).paddingTop), '16px');
  } finally {
    await browser.close();
  }
});
