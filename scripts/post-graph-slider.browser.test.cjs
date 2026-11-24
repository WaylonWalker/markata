'use strict';
const {test} = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const {chromium} = require(process.env.PLAYWRIGHT_MODULE || 'playwright');

const script = fs.readFileSync(path.join(__dirname, '../pkg/themes/default/static/js/post-graph.js'), 'utf8');

// Exercise real graph initialization with deterministic graph data and a D3
// force stub, avoiding external requests and physics timing dependencies.
test('node-limit slider is hidden only at or below its configured minimum', async () => {
  const browser = await chromium.launch({args: ['--no-sandbox'], executablePath: process.env.CHROMIUM_EXECUTABLE || undefined});
  try {
    for (const [connections, minimum, hidden] of [
      [7, 8, true], [8, 8, true], [9, 8, false],
      [23, 8, false], [24, 8, false],
      [12, 12, true], [13, 12, false],
    ]) {
      const page = await browser.newPage();
      try {
        const errors = [];
        page.on('pageerror', error => errors.push(error.message));
        await page.setContent(`<!doctype html><section class="post-graph" data-graph="/graph.json" data-post="/post/">
          <div style="width:320px;height:180px"><canvas class="post-graph__canvas-el"></canvas></div>
          <div class="post-graph__tooltip"></div>
          <label class="post-graph__control"><input class="post-graph__limit" type="range" min="${minimum}" max="120" value="100" data-default="100"><span class="post-graph__limit-value"></span></label>
          <input class="post-graph__animate" type="checkbox" checked>
        </section>`);
        await page.evaluate(count => {
          const chain = () => ({id() {return this;}, distance() {return this;}, strength() {return this;}, distanceMax() {return this;}, radius() {return this;}});
          window.d3 = {
            forceSimulation: () => ({force() {return this;}, on() {return this;}, alpha() {return this;}, restart() {return this;}, stop() {return this;}, tick() {return this;}}),
            forceLink: chain, forceManyBody: chain, forceCenter: () => ({}), forceCollide: chain,
          };
          const nodes = [{id: 'post', href: '/post/', label: 'Post', type: 'post'}];
          const edges = [];
          for (let i = 0; i < count; i++) {
            nodes.push({id: `node-${i}`, href: `/node-${i}/`, label: `Node ${i}`, type: 'post'});
            edges.push({source: 'post', target: `node-${i}`, type: 'link'});
          }
          window.fetch = async () => ({json: async () => ({nodes, edges})});
        }, connections);
        await page.addScriptTag({content: script});
        await page.waitForFunction(() => document.querySelector('.post-graph__limit')?.dataset.defaultSet === 'true');
        const actual = await page.locator('.post-graph__control').evaluate(el => ({hidden: el.hidden, display: el.style.display}));
        assert.equal(actual.hidden, hidden, `connections=${connections}, minimum=${minimum}`);
        assert.equal(actual.display, hidden ? 'none' : '', `connections=${connections}, minimum=${minimum}`);
        assert.deepEqual(errors, []);
      } finally {
        await page.close();
      }
    }
  } finally {
    await browser.close();
  }
});
