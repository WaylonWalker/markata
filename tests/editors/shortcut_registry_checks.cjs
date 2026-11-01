/* Run with node tests/editors/shortcut_registry_checks.cjs */
'use strict';
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const vm = require('node:vm');

const state = {dialog:false, saved:{}};
const document = {
  addEventListener() {},
  querySelector(query) { return query === 'dialog[open]' && state.dialog ? {} : null; }
};
const localStorage = {
  getItem(key) { return state.saved[key] || null; },
  setItem(key,value) { state.saved[key] = value; }
};
const context = {window:{}, document, localStorage, setTimeout, clearTimeout};
vm.runInNewContext(fs.readFileSync(path.join(__dirname, '../../pkg/workbench/web/shortcuts.js'),'utf8'),
  context, {filename:'shortcuts.js'});
const shortcuts = context.window.MarkataShortcuts;
let palette = 0, settings = 0, publish = 0, explorer = 0;
shortcuts.register({id:'palette',label:'Palette',bindings:{vscode:'Mod+Shift+P',writer:'Mod+K',vim:'Mod+Shift+P',emacs:'Alt+X'},run:()=>palette++});
shortcuts.register({id:'settings',label:'Settings',bindings:{vscode:'Mod+K Mod+S',writer:'Mod+Alt+K',vim:'Space ?',emacs:'Alt+Shift+K'},run:()=>settings++});
shortcuts.register({id:'publish',label:'Publish review',bindings:{vscode:'Mod+Shift+Enter',writer:'Mod+Shift+Enter'},run:()=>publish++});
shortcuts.register({id:'explorer',label:'Explorer',bindings:{vscode:'Mod+B',writer:'Mod+\\',vim:'Space E'},run:()=>explorer++});
assert.throws(()=>shortcuts.register({id:'palette',run:()=>{}}), /duplicate/);

function fire(key, opts={}) {
  let prevented=false, stopped=false;
  const event={key,ctrlKey:false,metaKey:false,altKey:false,shiftKey:false,
    isComposing:false,repeat:false,defaultPrevented:false,target:{matches:()=>false,closest:()=>false},
    preventDefault(){prevented=true;},stopPropagation(){stopped=true;},...opts};
  const used=shortcuts.dispatch(event);
  return {used,prevented,stopped};
}
assert.equal(shortcuts.getPreset(),'vscode');
assert.equal(fire('P',{ctrlKey:true,shiftKey:true}).used,true);
assert.equal(palette,1);
assert.equal(fire('k',{ctrlKey:true}).used,true);
assert.equal(settings,0);
assert.equal(fire('s',{ctrlKey:true}).used,true);
assert.equal(settings,1);
assert.equal(fire('Enter',{ctrlKey:true,shiftKey:true}).used,true);
assert.equal(publish,1);
assert.equal(fire('b',{ctrlKey:true}).used,true);
assert.equal(explorer,1);
state.dialog=true;
assert.equal(fire('Enter',{ctrlKey:true,shiftKey:true}).used,false);
assert.equal(publish,1);
state.dialog=false;
assert.equal(fire('Enter',{ctrlKey:true,shiftKey:true,isComposing:true}).used,false);

shortcuts.setPreset('writer');
assert.equal(fire('k',{ctrlKey:true}).used,true);
assert.equal(palette,2);
assert.equal(shortcuts.bindingFor('explorer'),'Mod+\\');
assert.equal(shortcuts.setBinding('publish','Mod+K').ok,false);
assert.equal(shortcuts.bindingFor('publish'),'Mod+Shift+Enter');
assert.equal(shortcuts.setBinding('publish','Alt+P').ok,true);
assert.equal(fire('p',{altKey:true}).used,true);
assert.equal(publish,2);
assert.equal(shortcuts.bindingFor('publish'),'Alt+P');
shortcuts.reset();
assert.equal(shortcuts.bindingFor('publish'),'Mod+Shift+Enter');

shortcuts.setPreset('vim');
const editor={state:{vim:{insertMode:false}},getOption:()=> 'vim'};
shortcuts.init(editor);
const vimTarget={matches:()=>false,closest:selector=>selector==='.CodeMirror' ? {} : null};
assert.equal(fire(' ',{target:vimTarget}).used,true);
assert.equal(fire('e',{target:vimTarget}).used,true);
assert.equal(explorer,2);
editor.state.vim.insertMode=true;
assert.equal(fire(' ',{target:vimTarget}).used,false);
assert.equal(fire('e',{target:vimTarget}).used,false);
assert.equal(explorer,2);
editor.state.vim.insertMode=false;
const inputTarget={matches:()=>true,closest:()=>null};
assert.equal(fire(' ',{target:inputTarget}).used,false);
assert.equal(fire('e',{target:inputTarget}).used,false);
assert.equal(explorer,2);
assert.equal(fire('P',{ctrlKey:true,shiftKey:true}).used,true);
assert.equal(palette,3);
assert.equal(shortcuts.getPreset(),'vim');

console.log('Shortcut registry tests passed');
