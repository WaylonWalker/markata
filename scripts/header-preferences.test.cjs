// Regression contract for the minimal header and its mirrored theme templates.
const { test } = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');

const root = path.join(__dirname, '..');
const read = (relative) => fs.readFileSync(path.join(root, relative), 'utf8');
const paths = ['templates', 'pkg/themes/default/templates'];

for (const prefix of paths) {
  test(prefix + ': adaptive preferences are included exactly once', () => {
    const base = read(prefix + '/base.html');
    const palette = read(prefix + '/partials/palette-switcher.html');
    const loading = read(prefix + '/partials/adaptive-loading.html');

    assert.match(base, /if not config\.theme\.switcher\.enabled/);
    assert.match(base, /include "partials\/adaptive-loading\.html"/);
    assert.match(palette, /include "partials\/adaptive-loading\.html"/);
    assert.doesNotMatch(base, /<div class="adaptive-loading"/);
    assert.match(palette, /theme-switcher--combined/);
    assert.match(palette, /aria-label="Appearance and loading preferences"/);
    assert.match(palette, /data-picker-mode="light"/);
    assert.match(palette, /data-picker-mode="dark"/);

    for (const mode of ['auto', 'save-data', 'full-quality']) {
      assert.match(loading, new RegExp('name="loading-mode" value="' + mode + '"'));
    }
    for (const action of ['data-loading-label', 'data-loading-mismatch', 'data-loading-status', 'data-loading-use-auto', 'data-loading-keep']) {
      assert.ok(loading.includes(action), action + ' is required by adaptive-loading.js');
    }
  });
}

test('picker leaves adaptive input keys alone and relocates fallback controls', () => {
  const js = read('pkg/themes/default/static/js/palette-switcher.js');
  assert.match(js, /target\.closest\('\.adaptive-loading'\)/);
  assert.match(js, /actions\.appendChild\(loading\)/);
  assert.match(js, /panel\.insertBefore\(fallback/);
});

test('both search implementations use the same compact sizing', () => {
  const css = read('pkg/themes/default/static/css/search.css');
  assert.match(css, /\.pagefind-ui__search-input/);
  assert.match(css, /\.bleve-search__input/);
  assert.match(css, /height: 2\.125rem/);
});
