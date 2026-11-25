'use strict';
const test = require('node:test');
const assert = require('node:assert/strict');
const vm = require('node:vm');
const fs = require('node:fs');
const path = require('node:path');

const window = {};
vm.runInNewContext(fs.readFileSync(path.join(__dirname,'wasm-preview.js'),'utf8'),
  {window, Worker: undefined, URL, setTimeout, clearTimeout, console});
const {simpleTextOnly, compareParity} = window.MarkataWasmPreview;

test('quick preview grammar includes ordinary paragraphs and headings', () => {
  for (const source of ['# Title\n\nText with **bold**.', 'One\n\nTwo', '> Quoted text\n\n- Item']) {
    assert.equal(simpleTextOnly(source), true, source);
  }
});
test('site-specific and generated-feature Markdown always chooses native fallback', () => {
  for (const source of [
    '---\ntitle: Test\n---', '\x60\x60\x60mermaid\nA-->B\n\x60\x60\x60',
    '!!! note\n    Hi', '::: md-section\nHello\n:::',
    '[Link](https://example.org)', '![photo](/images/1.png)',
    '[[wikilink]]', '<div>custom</div>', '{{ plugin }}',
    '| A | B |\n|---|---|', '- [x] Checklist', '    indented code',
  ]) {
    assert.equal(simpleTextOnly(source), false, source);
  }
});
test('matched Go HTML/source maps yield only a verified frontmatter prefix', () => {
  const native = {html: '<header>Title</header><h1>Text</h1>',
    source_lines: [5, 7], source_anchors: []};
  const wasm = {html: '<h1>Text</h1>', source_lines: [0,2]};
  assert.equal(compareParity(native,wasm,5).prefix, '<header>Title</header>');
  assert.equal(compareParity({...native, html:'Different HTML'},wasm,5), null);
  assert.equal(compareParity({...native, source_anchors:[1]},wasm,5), null);
  assert.equal(compareParity({...native, source_lines:[5,8]},wasm,5), null);
});
test('worker initialization is gated on an actual version and Worker capability', () => {
  assert.throws(() => window.MarkataWasmPreview.create(''), /unavailable/);
  assert.throws(() => window.MarkataWasmPreview.create('v1'), /unavailable/);
});

test('worker readiness and render results retain linear-memory measurements', async () => {
  let activeWorker;
  class FakeWorker {
    constructor() { activeWorker = this; }
    postMessage(request) {
      this.onmessage({data: {type: 'result', revision: request.revision,
        result: {schema_version: 1, html: '<p>ok</p>', source_lines: [], assets: []},
        render_ms: 1, memory_bytes: 131072}});
    }
    terminate() {}
  }
  const memoryWindow = {Worker: FakeWorker};
  vm.runInNewContext(fs.readFileSync(path.join(__dirname,'wasm-preview.js'),'utf8'),
    {window: memoryWindow, Worker: FakeWorker, URL, setTimeout, clearTimeout, console});
  const client = memoryWindow.MarkataWasmPreview.create('v1');
  activeWorker.onmessage({data: {type: 'ready', schema_version: 1, version: 'v1', memory_bytes: 65536}});
  assert.equal((await client.ready).memory_bytes, 65536);
  const result = await client.render('# sample');
  assert.equal(result.memory_bytes, 131072);
  client.close();
});
