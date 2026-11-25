import {test} from 'node:test';
import assert from 'node:assert/strict';
import {createSearch,remember,clearHistory,setEnabled} from './search-history.js';
const memory=new Map();
globalThis.localStorage={getItem:k=>memory.get(k)??null,setItem:(k,v)=>memory.set(k,v),removeItem:k=>memory.delete(k)};
const entries=[{id:'a',label:'Open Settings'},{id:'b',label:'Open Sidebar'},{id:'c',label:'Publish Current Post'}];
const make=scope=>createSearch(entries,['label'],[1],scope,x=>x.id);
test('fuzzy matching, case, exact precedence and empty-query history',()=>{
  clearHistory('commands');setEnabled(true);
  const search=make('commands');
  assert.equal(search('pblsh')[0].obj.id,'c');
  remember('commands','b');assert.equal(search('')[0].obj.id,'b');
  remember('commands','b');remember('commands','b');remember('commands','a');
  assert.equal(search('')[0].obj.id,'a');
  assert.equal(search('OPEN SETTINGS')[0].obj.id,'a');
  assert.equal(make('commands')('')[0].obj.id,'a');
  assert.equal(search('zzzzzz').length,0);
});
test('history scoped, clearable, disabled and malformed-safe',()=>{
  remember('files:one','c');assert.equal(make('files:one')('')[0].obj.id,'c');
  assert.notEqual(make('files:two')('')[0].obj.id,'c');
  clearHistory('files:one');assert.notEqual(make('files:one')('')[0].obj.id,'c');
  setEnabled(false);remember('files:one','c');setEnabled(true);
  assert.notEqual(make('files:one')('')[0].obj.id,'c');
  memory.set('markata-picks-v1:files:one','{"bad":true}');assert.doesNotThrow(()=>make('files:one')(''));
  memory.set('markata-picks-v1:files:one','[null,{"id":"c","at":"bad","count":3}]');assert.doesNotThrow(()=>make('files:one')(''));
});
test('blocked storage retains bounded session history',()=>{
  globalThis.localStorage={getItem(){throw Error('blocked');},setItem(){throw Error('blocked');},removeItem(){throw Error('blocked');}};
  remember('blocked','c');assert.equal(make('blocked')('')[0].obj.id,'c');
  setEnabled(false);assert.notEqual(make('blocked')('')[0].obj.id,'c');setEnabled(true);
  clearHistory('blocked');assert.notEqual(make('blocked')('')[0].obj.id,'c');
});
test('punctuation in a command label does not let history beat an exact query',()=>{
  const commands=[{id:'vim',label:'Editor keymap: Vim'},{id:'emacs',label:'Editor keymap: Emacs'}];
  const search=createSearch(commands,['label'],[1],'keymaps',x=>x.id);
  clearHistory('keymaps');setEnabled(true);
  remember('keymaps','vim');
  assert.equal(search('editor keymap emacs')[0].obj.id,'emacs');
});
