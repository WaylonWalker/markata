import {test} from 'node:test';
import assert from 'node:assert/strict';
import {createFileSearch, fileFacets, suggestFor, acceptSuggestion} from './file-search.js';
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

test('operators narrow results before fuzzy matching', () => {
  const search = createFileSearch([
    {path:'pages/blog/a.md', title:'Alpha', tags:['goals'], status:'draft'},
    {path:'pages/blog/b.md', title:'Beta', tags:['goals'], status:'published'},
    {path:'pages/notes/c.md', title:'Gamma', tags:['goals'], status:'draft'},
    {path:'pages/blog/d.md', title:'Delta', tags:['travel'], status:'draft'},
  ]);
  assert.deepEqual(search('#goals').map(m => m.obj.path).sort(), ['pages/blog/a.md','pages/blog/b.md','pages/notes/c.md']);
  assert.deepEqual(search('tag:travel').map(m => m.obj.path), ['pages/blog/d.md']);
  assert.deepEqual(search('status:draft').map(m => m.obj.path).sort(), ['pages/blog/a.md','pages/blog/d.md','pages/notes/c.md']);
  assert.deepEqual(search('#goals status:draft in:pages/blog').map(m => m.obj.path), ['pages/blog/a.md']);
  assert.equal(search('#goals alp')[0].obj.path, 'pages/blog/a.md');
  assert.equal(search('#missing').length, 0);
});

test('operators keep matching when the match list is larger than the result cap', () => {
  const files = Array.from({length:200}, (_,i) => ({path:`pages/blog/post-${i}.md`, title:`Post ${i}`, tags: i === 199 ? ['rare'] : []}));
  const search = createFileSearch(files);
  assert.deepEqual(search('#rare').map(m => m.obj.path), ['pages/blog/post-199.md']);
});

test('completes tags, statuses, and folders as operators are typed', () => {
  const facets = fileFacets([
    {path:'pages/blog/a.md', tags:['goals','travel']},
    {path:'pages/blog/b.md', tags:['goals']},
    {path:'pages/notes/c.md', tags:['gardening']},
    {path:'readme.md'},
  ]);
  assert.deepEqual(facets.tags, ['goals','gardening','travel']);
  assert.deepEqual(facets.folders, ['pages','pages/blog','pages/notes']);
  assert.deepEqual(suggestFor('#go', facets).items.map(item => item.value), ['#goals']);
  assert.deepEqual(suggestFor('status:p', facets).items.map(item => item.value), ['status:published']);
  assert.deepEqual(suggestFor('in:pages/b', facets).items.map(item => item.value), ['in:pages/blog']);
  assert.deepEqual(suggestFor('tag:tr', facets).items.map(item => item.value), ['tag:travel']);
  assert.deepEqual(suggestFor('plain words', facets).items, []);
  assert.equal(acceptSuggestion('alp #go', '#goals'), 'alp #goals ');
  assert.equal(acceptSuggestion('in:pag', 'in:pages'), 'in:pages ');
});

test('feed: restricts results to a feed and completes feed names', () => {
  const files = [
    {path:'pages/blog/a.md', title:'Alpha'},
    {path:'pages/blog/b.md', title:'Beta'},
    {path:'pages/notes/c.md', title:'Gamma'},
  ];
  const feeds = () => [{slug:'Blog', paths:['pages/blog/a.md']}, {slug:'empty', paths:[]}];
  const search = createFileSearch(files, feeds);
  assert.deepEqual(search('feed:blog').map(m => m.obj.path), ['pages/blog/a.md']);
  assert.deepEqual(search('feed:blog gamma').map(m => m.obj.path), []);
  assert.deepEqual(search('feed:missing').map(m => m.obj.path), []);
  assert.deepEqual(suggestFor('feed:b', {tags:[], folders:[], feeds:['blog','empty']}).items.map(item => item.value), ['feed:blog']);
});

test('$feed is a shorthand for feed:', () => {
  const files = [{path:'pages/thoughts/a.md'}, {path:'pages/blog/b.md'}];
  const search = createFileSearch(files, () => [{slug:'thoughts-feed', paths:['pages/thoughts/a.md']}]);
  assert.deepEqual(search('$thoughts-feed').map(m => m.obj.path), ['pages/thoughts/a.md']);
  assert.deepEqual(suggestFor('$thou', {tags:[], folders:[], feeds:['thoughts-feed','blog']}).items.map(item => item.value), ['$thoughts-feed']);
  assert.equal(acceptSuggestion('$th', '$thoughts-feed'), '$thoughts-feed ');
});
