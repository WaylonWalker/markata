// Regression test for image-library.js DOM ordering without browser dependencies.
// Run: node --test scripts/image-library-order.test.cjs
const {test} = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const vm = require('node:vm');

test('hidden cards do not trigger repeated image reparenting', () => {
  const source = fs.readFileSync(
    path.join(__dirname, '../pkg/themes/default/static/js/image-library.js'), 'utf8'
  );
  const cards = ['a', 'b', 'c'].map((name, i) => ({
    hidden: false,
    getAttribute(key) {
      const attrs = {
        'data-name': name,
        'data-search': name,
        'data-src': '/images/' + name + '.webp',
        'data-added': String((i + 1) * 10),
        'data-uses': '1',
        'data-last-used': '0',
      };
      return attrs[key] ?? null;
    },
  }));
  const handlers = {};
  const search = {value: '', addEventListener(type, fn) { handlers['search:' + type] = fn; }};
  const sort = {value: 'name', addEventListener(type, fn) { handlers['sort:' + type] = fn; }};
  const grid = {
    children: [...cards],
    moves: 0,
    appendChild(fragment) {
      for (const card of fragment.items) {
        const oldIndex = this.children.indexOf(card);
        assert.notEqual(oldIndex, -1);
        this.children.splice(oldIndex, 1);
        this.children.push(card);
        this.moves++;
      }
    },
  };
  const counter = {textContent: ''};
  const empty = {hidden: false};
  const root = {
    querySelector(selector) {
      return ({
        '[data-image-grid]': grid,
        '[data-image-search]': search,
        '[data-image-sort]': sort,
        '[data-image-count]': counter,
        '[data-image-empty]': empty,
      })[selector] ?? null;
    },
    querySelectorAll(selector) {
      if (selector === '[data-image-card]') return cards;
      return [];
    },
    addEventListener() {},
  };
  const document = {
    readyState: 'complete',
    querySelector: () => root,
    createDocumentFragment() {
      return {items: [], appendChild(card) {this.items.push(card);}};
    },
  };
  const window = {
    requestAnimationFrame(fn) {fn(); return 1;},
    cancelAnimationFrame() {},
  };
  vm.runInNewContext(source, {document, window, navigator: {}}, {
    filename: 'image-library.js',
  });

  assert.equal(grid.moves, 0, 'initial sorted order should not be reparented');
  search.value = 'b';
  handlers['search:input']();
  assert.equal(counter.textContent, '1 of 3 media items');
  assert.equal(grid.moves, 0, 'hidden preceding siblings must not force reordering');
  handlers['search:input']();
  assert.equal(grid.moves, 0, 'repeat filtering must not move unchanged cards');

  search.value = '';
  handlers['search:input']();
  sort.value = 'recent';
  handlers['sort:change']();
  assert.deepEqual(grid.children, [cards[2], cards[1], cards[0]]);
  const firstSortMoves = grid.moves;
  assert.equal(firstSortMoves, 3, 'new sort should reorder visible cards exactly once');
  handlers['sort:change']();
  assert.equal(grid.moves, firstSortMoves, 'same sort must never reparent cards');

  search.value = 'a';
  handlers['search:input']();
  assert.equal(grid.moves, firstSortMoves, 'filtering later card does not require moves');
  search.value = '';
  handlers['search:input']();
  assert.equal(grid.moves, firstSortMoves, 'clearing filter retains current sort');
  sort.value = 'name';
  handlers['sort:change']();
  assert.deepEqual(grid.children, cards, 'name sort restores original order');
});
