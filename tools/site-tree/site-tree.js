import {FileTree} from '@pierre/trees';

const files = SITE_FILES;
const storageKey = path => 'markata-example-site-v1:' + path;
let active = 'draft.md', drafts = {}, currentEditor, warning, label;
let fileTree, feedTree, syncing = false;
const aliases = {};
function syncSelection(tree, path) {
  if (!tree || !path) return;
  for (const previous of tree.getSelectedPaths()) if (previous !== path) tree.getItem(previous)?.deselect();
  if (!tree.getSelectedPaths().includes(path)) tree.getItem(path)?.select();
}
const editable = path => path.endsWith('.md');
function remember() {
  if (!currentEditor || !editable(active)) return;
  drafts[active] = currentEditor.getValue();
  try { localStorage.setItem(storageKey(active), drafts[active]); }
  catch { warning(); }
}
function open(path) {
  if (syncing || (!(path in files) && path !== 'draft.md')) return;
  syncing = true;
  try {
    syncSelection(fileTree, path);
    const selectedAlias = feedTree?.getSelectedPaths().find(alias => aliases[alias] === path);
    const alias = selectedAlias || Object.keys(aliases).find(alias => aliases[alias] === path);
    for (const previous of feedTree?.getSelectedPaths() || []) {
      if (previous !== alias) feedTree.getItem(previous)?.deselect();
    }
    syncSelection(feedTree, alias);
  } finally { syncing = false; }
  remember();
  active = path;
  try { localStorage.setItem('markata-example-site-active-v1', path); } catch { warning(); }
  if (!(path in drafts)) {
    drafts[path] = files[path] || '';
    if (editable(path)) {
      try { drafts[path] = localStorage.getItem(storageKey(path)) ?? drafts[path]; }
      catch { warning(); }
    }
  }
  currentEditor.setOption('readOnly', !editable(path));
  document.querySelector('#commands').disabled = !editable(path);
  currentEditor.setOption('mode', editable(path) ? 'gfm' : 'text/plain');
  label.textContent = path + (editable(path) ? '' : ' · read only');
  currentEditor.setValue(drafts[path]);
  currentEditor.clearHistory();
  document.querySelector('#download').textContent = editable(path) ? 'Download .md' : 'Download .toml';
}
window.MarkataSiteTree = {
  save() {
    if (!currentEditor) return false;
    if (!editable(active)) return true;
    drafts[active] = currentEditor.getValue();
    try { localStorage.setItem(storageKey(active), drafts[active]);
      document.querySelector('#saved').textContent = 'Draft saved on this device'; }
    catch { warning(); }
    return true;
  },
  filename: () => active.split('/').pop(),
  markdown(text) {
    if (!editable(active)) return '````toml\n' + text + '\n````';
    return text.replace(/^---\r?\n[\s\S]*?\r?\n---(?:\r?\n|$)/, '');
  },
  import(text) { open('draft.md'); currentEditor.setValue(text); },
  attach(editor, storageWarning) {
    currentEditor = editor; warning = storageWarning;
    editor.on('beforeChange', (instance, change) => {
      if (!editable(active) && change.origin !== 'setValue') change.cancel();
    });
    let selected = 'posts/welcome.md';
    try {
      selected = localStorage.getItem('markata-example-site-active-v1') || selected;
      drafts['draft.md'] = localStorage.getItem(storageKey('draft.md')) ?? editor.getValue();
    } catch { drafts['draft.md'] = editor.getValue(); warning(); }
    editor.setValue(drafts['draft.md']);
    label = document.querySelector('#active-file');
    label.textContent = 'draft.md';
    fileTree = new FileTree({paths:['draft.md', ...Object.keys(files)], initialExpansion:'open',
      search:true,
      onSelectionChange(paths) { if (paths.length) open(paths[paths.length - 1]); }});
    fileTree.render({containerWrapper:document.querySelector('#site-files')});
    // Feed aliases resolve to canonical source paths, so every view shares a draft.
    for (const [path, text] of Object.entries(files)) {
      if (!path.startsWith('posts/')) continue;
      const tags = text.match(/^tags:\s*\[([^\]]*)\]/m)?.[1].split(',').map(tag => tag.trim()) || [];
      for (const feed of ['notes', 'garden']) {
        if (tags.includes(feed)) aliases[`${feed}/${path}`] = path;
      }
    }
    feedTree = new FileTree({paths:Object.keys(aliases), initialExpansion:'open', search:true,
      onSelectionChange(paths) { const path = aliases[paths[paths.length - 1]]; if (path) open(path); }});
    feedTree.render({containerWrapper:document.querySelector('#site-feeds')});
    document.querySelector('#site-view').addEventListener('change', event => {
      document.querySelector('#site-files').hidden = event.target.value !== 'files';
      document.querySelector('#site-feeds').hidden = event.target.value !== 'feeds';
    });
    // Start with a real site post. The previous standalone draft stays in the tree.
    fileTree.getItem(selected in files || selected === 'draft.md' ? selected : 'posts/welcome.md')?.select();
  }
};
