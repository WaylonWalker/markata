'use strict';
const {test} = require('node:test');
const assert = require('node:assert/strict');
const {parseArgs, isScriptResource, summarize} = require('./js-loading-audit.cjs');

test('CLI accepts multiple routes and validates input', () => {
  const opts = parseArgs(['--url', 'https://example.com/post/', '--url', 'http://localhost:8000/archive/', '--wait-ms', '250']);
  assert.equal(opts.urls.length, 2);
  assert.equal(opts.waitMs, 250);
  assert.throws(() => parseArgs(['--url', 'file:///etc/passwd']), /Unsupported protocol/);
  assert.throws(() => parseArgs(['--wait-ms', '-5', '--url', 'https://example.com']), /wait-ms/);
  assert.throws(() => parseArgs([]), /Provide at least one/);
});

test('identifies JavaScript resources but excludes stylesheets', () => {
  assert.ok(isScriptResource({name: 'https://a.test/js/script.abcdef12.js', initiatorType: 'script'}));
  assert.ok(isScriptResource({name: 'https://a.test/js/core.mjs?v=2', initiatorType: 'script'}));
  assert.equal(isScriptResource({name: 'https://a.test/css/style.css', initiatorType: 'link'}), false);
});

test('accounts for zero transfer-size entries without claiming they are cached', () => {
  const summary = summarize({url: 'https://example.com/posts/', title: 'Posts', resources: [
    {name: 'https://example.com/js/core.js', initiatorType: 'script', transferSize: 1500, encodedBodySize: 1200, decodedBodySize: 7000, duration: 3.2},
    {name: 'https://cdn.example.org/library.js', initiatorType: 'script', transferSize: 0, encodedBodySize: 0, decodedBodySize: 0, duration: 8},
    {name: 'https://example.com/css/site.css', initiatorType: 'link', transferSize: 999},
  ], inlineExecutableScriptCount: 2, inlineExecutableScriptCharacters: 700});
  assert.equal(summary.scriptRequests, 2);
  assert.equal(summary.sameOriginScriptRequests, 1);
  assert.equal(summary.crossOriginScriptRequests, 1);
  assert.equal(summary.observedTransferBytes, 1500);
  assert.equal(summary.observedDecodedBytes, 7000);
  assert.equal(summary.zeroTransferEntries, 1);
  assert.equal(summary.inlineExecutableScriptCharacters, 700);
});
