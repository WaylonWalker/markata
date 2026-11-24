const test = require('node:test');
const assert = require('node:assert/strict');
const policy = require('../pkg/themes/default/static/js/adaptive-loading.js');

test('visible article autoplay respects constrained mode, motion, and Save Data', () => {
  assert.equal(policy.canAutoplayInViewport('unknown', false, false), true);
  assert.equal(policy.canAutoplayInViewport('normal', false, false), true);
  assert.equal(policy.canAutoplayInViewport('full-quality', false, false), true);
  assert.equal(policy.canAutoplayInViewport('constrained', false, false), false);
  assert.equal(policy.canAutoplayInViewport('normal', true, false), false);
  assert.equal(policy.canAutoplayInViewport('full-quality', false, true), false);
});

test('missing Network Information API leaves Auto usable', () => {
  const current = policy.networkHint(policy.state(), undefined);
  assert.equal(current.recommendation, 'unknown');
  assert.equal(policy.effectiveMode('auto', current.recommendation), 'unknown');
});

test('strong Network Information signal recommends a matching policy', () => {
  const slow = policy.networkHint(policy.state(), { saveData: true });
  const fast = policy.networkHint(policy.state(), { effectiveType: '4g', downlink: 8 });
  assert.equal(slow.recommendation, 'constrained');
  assert.equal(fast.recommendation, 'normal');
});

test('one weak hint does not change policy, repeated poor resources do', () => {
  let current = policy.networkHint(policy.state(), { effectiveType: '2g' });
  assert.equal(current.recommendation, 'constrained');
  current = policy.vote(current, 'poor');
  assert.equal(current.recommendation, 'constrained');
});

test('fast resource evidence requires four observations', () => {
  let current = policy.state();
  current = policy.vote(current, 'good');
  current = policy.vote(current, 'good');
  current = policy.vote(current, 'good');
  assert.notEqual(current.recommendation, 'full-quality');
  current = policy.vote(current, 'good');
  assert.equal(current.recommendation, 'full-quality');
});

test('manual Save Data always wins and mismatch waits for high confidence', () => {
  let current = policy.state();
  for (let i = 0; i < 4; i++) current = policy.vote(current, 'good');
  assert.equal(policy.effectiveMode('save-data', current.recommendation), 'constrained');
  assert.equal(policy.mismatch('save-data', current), '');
  current = policy.vote(current, 'good');
  assert.equal(policy.mismatch('save-data', current), 'full-quality');
});

test('manual Full Quality wins while sustained poor evidence raises a mismatch', () => {
  let current = policy.state();
  for (let i = 0; i < 4; i++) current = policy.vote(current, 'poor');
  assert.equal(policy.effectiveMode('full-quality', current.recommendation), 'full-quality');
  assert.equal(policy.mismatch('full-quality', current), 'constrained');
});

test('hysteresis prevents recommendation thrashing', () => {
  let current = policy.state();
  current = policy.vote(current, 'poor');
  current = policy.vote(current, 'poor');
  assert.equal(current.recommendation, 'constrained');
  current = policy.vote(current, 'good');
  assert.equal(current.recommendation, 'constrained');
  current = policy.vote(current, 'poor');
  assert.equal(current.recommendation, 'constrained');
  for (let i = 0; i < 4; i++) current = policy.vote(current, 'good');
  assert.equal(current.recommendation, 'full-quality');
});

test('session evidence persists briefly and expires', () => {
  let current = policy.state();
  current = policy.vote(current, 'poor');
  current = policy.vote(current, 'poor');
  const serialized = policy.savedState(current, 1000);
  assert.equal(policy.restore(serialized, 2000).recommendation, 'constrained');
  assert.equal(policy.restore(serialized, 1000 + 10 * 60 * 1000 + 1).recommendation, 'unknown');
});

test('manual preference persists and Auto clears the override', () => {
  const values = new Map();
  const storage = {
    getItem: key => values.has(key) ? values.get(key) : null,
    setItem: (key, value) => values.set(key, value),
    removeItem: key => values.delete(key),
  };
  policy.persistMode(storage, 'mode', 'save-data');
  assert.equal(policy.readMode(storage, 'mode'), 'save-data');
  policy.persistMode(storage, 'mode', 'auto');
  assert.equal(policy.readMode(storage, 'mode'), 'auto');
});

test('offline applies immediately and online recovery still uses hysteresis', () => {
  let current = policy.vote(policy.state(), 'offline');
  assert.equal(current.recommendation, 'constrained');
  current = policy.vote(current, 'online');
  assert.equal(current.recommendation, 'constrained');
  current = policy.vote(current, 'good');
  assert.equal(current.recommendation, 'constrained');
});

test('sample count is bounded', () => {
  let current = policy.state();
  for (let i = 0; i < 20; i++) current = policy.vote(current, 'poor');
  assert.equal(current.samples, 12);
  assert.equal(current.poor, 8);
});

 test('generic fast hints cannot erase measured poor evidence', () => {
  let current = policy.state();
  for (let i = 0; i < 4; i++) current = policy.vote(current, 'poor');
  current = policy.networkHint(current, {effectiveType: '4g', downlink: 10});
  assert.equal(current.poor, 4);
  assert.equal(current.recommendation, 'constrained');
  for (let i = 0; i < 3; i++) current = policy.vote(current, 'good');
  assert.equal(current.recommendation, 'constrained');
  assert.equal(policy.restore(policy.savedState(current, 1000), 1001).recommendation, 'constrained');
  current = policy.vote(current, 'good');
  assert.equal(current.recommendation, 'full-quality');
 });
 test('offline discards old fast confidence', () => {
  let current = policy.state();
  for (let i = 0; i < 5; i++) current = policy.vote(current, 'good');
  current = policy.vote(current, 'offline');
  current = policy.vote(current, 'online');
  current = policy.networkHint(current, {effectiveType:'4g', downlink:10});
  current = policy.vote(current, 'good');
  assert.equal(current.recommendation, 'constrained');
 });
