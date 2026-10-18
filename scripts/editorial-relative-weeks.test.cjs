'use strict';
const test = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const vm = require('node:vm');
const path = require('node:path');

const script = fs.readFileSync(path.join(__dirname, '../pkg/themes/default/static/js/editorial-relative-weeks.js'), 'utf8');

function renderLabels(dateParts, starts) {
  const headings = starts.map(weekStart => ({dataset: {weekStart}, textContent: 'fallback'}));
  const listeners = {};
  const doc = {
    readyState: 'complete',
    querySelectorAll() { return headings; },
    addEventListener(name, callback) { listeners[name] = callback; }
  };
  const win = {addEventListener(name, callback) { listeners[name] = callback; }};
  const FixedDate = class extends Date {
    constructor(...args) { super(...(args.length ? args : dateParts)); }
  };
  vm.runInNewContext(script, {window: win, document: doc, Date: FixedDate});
  return {headings, listeners, doc};
}

test('calendar-relative editorial week headings', () => {
  const {headings} = renderLabels([2026, 9, 18, 12], [
    '2026-10-12', '2026-10-05', '2026-09-28', '2026-08-17',
    '2025-10-13', '2024-10-14', '2026-10-19', 'undated'
  ]);
  assert.deepEqual(headings.map(x => x.textContent), [
    'This week', 'Last week', '2 weeks ago', '2 months ago',
    '1 year ago', '2 years ago', 'Next week', 'fallback'
  ]);
});

test('a cached archive moves to last week on Monday and updates after HTMX swaps', () => {
  const {headings, listeners} = renderLabels([2026, 9, 19, 12], ['2026-10-12']);
  assert.equal(headings[0].textContent, 'Last week');
  const injected = {dataset: {weekStart: '2026-10-05'}, textContent: 'fallback'};
  const target = {querySelectorAll() { return [injected]; }};
  listeners['htmx:afterSwap']({target, detail: {target}});
  assert.equal(injected.textContent, '2 weeks ago');
});
