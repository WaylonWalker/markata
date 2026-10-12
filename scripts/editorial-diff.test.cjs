const test = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const vm = require('node:vm');
const context = { window: {} };
vm.runInNewContext(fs.readFileSync('pkg/themes/default/static/js/editorial-diff.js', 'utf8'), context);
const { compare, body } = context.window.MarkataEditorialDiff;

test('word comparisons reconstruct both exact revisions, including Unicode and punctuation', () => {
  for (const [before, after] of [
    ['A tpyο.', 'A typo.'], ['', 'Added\ntext'], ['Removed\ntext', ''],
    ['same', 'same'], ['<script>alert(1)</script>', '<img src=x onerror=alert(2)>'],
    ['café mañana\nこんにちは', 'café demain\nこんばんは'],
    ['one two three', 'three two one'],
  ]) {
    const pieces = compare(before, after);
    assert.equal(pieces.filter(p => p.kind !== 'added').map(p => p.text).join(''), before);
    assert.equal(pieces.filter(p => p.kind !== 'removed').map(p => p.text).join(''), after);
  }
});
test('bounded large comparison preserves exact text', () => {
  const before = 'before '.repeat(1500), after = 'after '.repeat(1500);
  const pieces = compare(before, after);
  assert.equal(pieces.filter(p => p.kind !== 'added').map(p => p.text).join(''), before);
  assert.equal(pieces.filter(p => p.kind !== 'removed').map(p => p.text).join(''), after);
});
test('frontmatter is omitted, without treating later thematic breaks as metadata', () => {
  assert.equal(body('---\r\ntitle: Old\r\n---\r\n\r\nBody\r\n---\r\nEnd'), '\nBody\n---\nEnd');
  assert.equal(body('No frontmatter\n---\nEnd'), 'No frontmatter\n---\nEnd');
  assert.throws(() => body('---\ntitle: Missing delimiter'));
});
