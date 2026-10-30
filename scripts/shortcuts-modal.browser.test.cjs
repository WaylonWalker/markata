#!/usr/bin/env node
// Browser acceptance for both shipped shortcuts template copies.
// Uses actual template DOM and shipped shortcut JS, without CSS or a network server.
// A full generated-site legacy-template smoke remains a separate deployment gate.
'use strict';

const assert = require('node:assert/strict');
const fs = require('node:fs/promises');
const path = require('node:path');
const { chromium } = require('playwright');

const root = path.resolve(__dirname, '..');
const templates = [
  'templates/partials/shortcuts-modal.html',
  'pkg/themes/default/templates/partials/shortcuts-modal.html',
];

function templateBody(template) {
  // Only remove template control-flow directives; keep the actual checked HTML.
  return template
    .replace(/\{#[\s\S]*?#\}/g, '')
    .replace(/\{%[\s\S]*?%\}/g, '')
    .replace(/\{\{[\s\S]*?\}\}/g, '');
}

async function checkTemplate(browser, filename) {
  const body = templateBody(await fs.readFile(path.join(root, filename), 'utf8'));
  const markup = `<!doctype html><html lang="en"><head><meta charset="utf-8"></head>
    <body><button id="opener" type="button" data-shortcuts-trigger
    aria-controls="shortcuts-modal">Keyboard shortcuts</button>
    <main><p id="article">Sample article</p></main>
    ${body}</body></html>`;

  // Browser should respect native hidden/inert even if CSS and JS never arrive.
  const noJS = await browser.newContext({ javaScriptEnabled: false });
  const staticPage = await noJS.newPage();
  await staticPage.setContent(markup);
  const staticModal = staticPage.locator('#shortcuts-modal');
  assert.equal(await staticModal.count(), 1, `${filename}: exactly one dialog`);
  assert.equal(await staticModal.isVisible(), false, `${filename}: visible without CSS/JS`);
  assert.equal(await staticModal.getAttribute('aria-hidden'), 'true');
  assert.equal(await staticModal.getAttribute('inert'), '');
  await staticPage.keyboard.press('Tab');
  assert.notEqual(await staticPage.evaluate(() => document.activeElement?.className),
    'shortcuts-modal-close', `${filename}: closed dialog received focus`);
  await noJS.close();

  const page = await browser.newPage();
  const errors = [];
  page.on('pageerror', error => errors.push(error.message));
  await page.setContent(markup);
  await page.addScriptTag({ path: path.join(root, 'pkg/themes/default/static/js/shortcuts-registry.js') });
  await page.addScriptTag({ path: path.join(root, 'pkg/themes/default/static/js/search-shortcuts.js') });
  await page.waitForFunction(() => window.shortcutsRegistry?.getAll().some(item => item.key === '?'));

  const modal = page.locator('#shortcuts-modal');
  const close = page.locator('.shortcuts-modal-close');
  assert.equal(await modal.isVisible(), false, `${filename}: initially visible`);
  await page.locator('#opener').click();
  assert.equal(await modal.isVisible(), true, `${filename}: opener failed`);
  assert.equal(await modal.getAttribute('hidden'), null);
  assert.equal(await modal.getAttribute('inert'), null);
  assert.equal(await modal.getAttribute('aria-hidden'), 'false');
  assert.equal(await close.evaluate(el => document.activeElement === el), true,
    `${filename}: close button not focused`);

  await page.keyboard.press('Escape');
  assert.equal(await modal.isVisible(), false, `${filename}: Escape did not hide`);
  assert.equal(await modal.getAttribute('inert'), '');
  assert.equal(await modal.getAttribute('aria-hidden'), 'true');
  assert.equal(await page.locator('#opener').evaluate(el => document.activeElement === el), true,
    `${filename}: opener focus not restored`);

  await page.locator('#opener').click();
  await close.click();
  assert.equal(await modal.isVisible(), false, `${filename}: close button did not hide`);
  assert.deepEqual(errors, [], `${filename}: JavaScript errors`);
  await page.close();
  console.log(`PASS ${filename}: no-JS hidden, opener, focus, Escape, close`);
}

(async () => {
  const browser = await chromium.launch({ headless: true });
  try {
    for (const filename of templates) await checkTemplate(browser, filename);
  } finally {
    await browser.close();
  }
})().catch(error => {
  console.error(error);
  process.exitCode = 1;
});
