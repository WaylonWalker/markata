import {test} from 'node:test';
import assert from 'node:assert/strict';
import {createFileSearch} from './file-search.js';
test('ranks filenames, paths and titles and caps results', () => {
  const search = createFileSearch([{path:'docs/quickstart.md',title:'Getting started'}, {path:'quickstart/archive.md'}, {path:'posts/other.md',title:'Quickstart'}]);
  assert.equal(search('qkst')[0].obj.path, 'docs/quickstart.md');
  assert.equal(search('docs/quick')[0].obj.path, 'docs/quickstart.md');
  assert.equal(search('getting started')[0].obj.path, 'docs/quickstart.md');
  assert.equal(search('zzzzzz').length, 0);
  assert.equal(search('QUICKSTART')[0].obj.path, 'docs/quickstart.md');
});
test('10,000 files: bounded output and repeated search timing', () => {
  const search = createFileSearch(Array.from({length:10000}, (_,i) => ({path:`content/topic-${i}/article-${i}.md`, title:`Article ${i}`})));
  const start = performance.now();
  for (let i=0;i<100;i++) assert.ok(search('art' + i).length <= 60);
  console.log(`10,000 files, 100 queries: ${(performance.now()-start).toFixed(1)}ms`);
  assert.equal(search('').length,60);
});
