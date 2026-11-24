'use strict';
const test = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const vm = require('node:vm');
const path = require('node:path');

const scope = {window: {}, TextEncoder};
vm.runInNewContext(fs.readFileSync(path.join(__dirname, 'preview-contract.js'), 'utf8'), scope);
const {fromGo, fromWasm, chooseRenderer, wasmEligibility} = scope.window.MarkataPreviewContract;
const go = () => ({
  schema_version: 1, renderer_fingerprint: 'build-abc',
  html: '<h1>Title</h1>', document: '<!doctype html>',
  source_lines: [0, 2, 3], source_anchors: [1, 4],
  feature_head_html: '', body_end_html: '',
  render_ms: 2.5, queue_ms: 0, input_bytes: 21, html_bytes: 14, cached: false,
});

test('Go payload keeps first-frame HTML and offsets source lines exactly once', () => {
  const original = go();
  const normalized = fromGo(original, 5, true);
  assert.equal(normalized.document, '<!doctype html>');
  assert.deepEqual(Array.from(normalized.source_lines), [5, 7, 8]);
  assert.deepEqual(Array.from(normalized.source_anchors), [6, 9]);
  assert.equal(normalized.source_offset, 5);
  assert.deepEqual(original.source_lines, [0, 2, 3], 'normalization must not mutate source');
  assert.deepEqual(Array.from(fromGo(original, 0).source_lines), [0, 2, 3]);
});

test('Version skew and missing initial shell fail before iframe update', () => {
  assert.throws(() => fromGo({...go(), schema_version: 2}), /Incompatible Go/);
  assert.throws(() => fromGo({...go(), renderer_fingerprint: undefined}), /Incomplete Go/);
  assert.throws(() => fromGo({...go(), document: null}, 0, true), /Initial Go/);
  assert.throws(() => fromGo({...go(), source_lines: [0, -1]}), /source lines/);
  assert.throws(() => fromGo({...go(), source_anchors: [NaN]}), /source anchors/);
  assert.throws(() => fromGo({...go(), render_ms: -3}), /timings/);
  assert.throws(() => fromGo(go(), -1), /offset/);
});

test('WASM source maps adapt a base-style standalone fragment', () => {
  const normalized = fromWasm({
    schema_version: 1, html: '<h1>Hello</h1>', source_lines: [0, 2],
    features: {}, assets: [
      {id: 'css/variables.css', kind: 'style', content: ''},
      {id: 'css/main.css', kind: 'style', content: ''},
      {id: 'css/components.css', kind: 'style', content: ''},
    ],
  }, {render_ms: 4.8, input_bytes: 18});
  assert.equal(normalized.document, null);
  assert.deepEqual(Array.from(normalized.source_lines), [0, 2]);
  assert.equal(normalized.html, '<h1>Hello</h1>');
  assert.equal(normalized.body_end_html, '');
  assert.equal(normalized.render_ms, 4.8);
  assert.equal(normalized.queue_ms, 0);
  assert.equal(normalized.input_bytes, 18);
  assert.equal(normalized.html_bytes, new TextEncoder().encode('<h1>Hello</h1>').length);
});

test('WASM requires Go asset orchestration for feature-rich previews', () => {
  const base = {schema_version: 1, html: 'hello', source_lines: [], features: {}, assets: []};
  assert.throws(() => fromWasm({...base, assets: [{id: 'css/code.css'}]}), /Go-managed/);
  assert.throws(() => fromWasm({...base, features: {mermaid: true}}), /Go-managed/);
  assert.throws(() => fromWasm({...base, schema_version: 2}), /Incompatible/);
  assert.throws(() => fromWasm({...base, source_lines: 'wrong'}), /source lines/);
  assert.throws(() => fromWasm({error: 'broken'}), /broken/);
});

test('WASM cannot activate just because a worker loaded', () => {
  const same = {workerReady: true, previewFrameReady: true,
    siteFingerprint: 'build-abc', workerFingerprint: 'build-abc'};
  assert.equal(chooseRenderer({}).renderer, 'go');
  assert.equal(chooseRenderer({...same}).reason, 'parity-not-certified');
  assert.equal(chooseRenderer({...same, parityCertified: true,
    workerFingerprint: 'build-def'}).reason, 'renderer-fingerprint-mismatch');
  assert.equal(chooseRenderer({...same, parityCertified: true,
    requiresSiteRender: true}).reason, 'site-dependent-content');
  assert.equal(chooseRenderer({...same, parityCertified: true}).renderer, 'wasm');
});

test('Missing shell and worker readiness always keep Go available', () => {
  const opts = {siteFingerprint: 'v1', workerFingerprint: 'v1', parityCertified: true};
  assert.equal(chooseRenderer({...opts, workerReady: false, previewFrameReady: true}).renderer, 'go');
  assert.equal(chooseRenderer({...opts, workerReady: true, previewFrameReady: false}).renderer, 'go');
});

test('WASM avoids browsers with limited or unknown resources and mismatched site renderers', () => {
  const capable = {workerSupported: true, wasmSupported: true,
    siteFingerprint: 'wasm-v1', workerFingerprint: 'wasm-v1',
    deviceMemory: 8, hardwareConcurrency: 8};
  assert.equal(wasmEligibility(capable).allowed, true);
  assert.equal(wasmEligibility({...capable, siteFingerprint: ''}).reason,
    'renderer-fingerprint-mismatch');
  assert.equal(wasmEligibility({...capable, deviceMemory: 4}).reason,
    'device-memory-limited');
  assert.equal(wasmEligibility({...capable, deviceMemory: undefined}).reason,
    'device-memory-limited');
  assert.equal(wasmEligibility({...capable, hardwareConcurrency: 4}).reason,
    'processor-capacity-limited');
  assert.equal(wasmEligibility({...capable, saveData: true}).reason, 'reduced-data');
  assert.equal(wasmEligibility({...capable, effectiveType: '2g'}).reason, 'reduced-data');
  assert.equal(wasmEligibility({...capable, workerSupported: false}).reason,
    'browser-unsupported');
});
