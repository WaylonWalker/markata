import {FileTree} from '@pierre/trees';

let files = SITE_FILES;
let workspaceKey = "example", metadata = [], selection = 0;
const storageKey = path => 'markata-workspace-v1:' + workspaceKey + ':' + path;
let active = 'draft.md', drafts = {}, currentEditor, warning, label;
let fileTree, feedTree, syncing = false;
const aliases = {};
function syncSelection(tree, path) {
  if (!tree || !path) return;
  for (const previous of tree.getSelectedPaths()) if (previous !== path) tree.getItem(previous)?.deselect();
  if (!tree.getSelectedPaths().includes(path)) tree.getItem(path)?.select();
}
const editable = path => /\.(md|markdown|mdown|mkd)$/i.test(path);
function readDraft(path, fallback) {
  const cached = localStorage.getItem(storageKey(path));
  if (cached === null) return fallback;
  try {
    const record = JSON.parse(cached);
    if (record?.format === 'markata-draft-v2' && typeof record.source === 'string') return record.source;
  } catch { /* Legacy drafts contain plain Markdown. */ }
  // Earlier builds could persist an empty placeholder before loading source.
  // Preserve all nonempty legacy drafts; new records also preserve intentional blanks.
  return cached === '' && path !== 'draft.md' ? fallback : cached;
}
function writeDraft(path, source) {
  localStorage.setItem(storageKey(path), JSON.stringify({format:'markata-draft-v2', source}));
}
function remember() {
  if (!currentEditor || !editable(active)) return;
  drafts[active] = currentEditor.getValue();
  try { writeDraft(active, drafts[active]); }
  catch { warning(); }
}
async function open(path) {
  if (syncing) return;
  const request = ++selection;
  if (path in files && files[path] === null && !(path in drafts)) {
    try {
      const response = await fetch("/workspace/file?path=" + encodeURIComponent(path));
      const data = await response.json();
      if (!response.ok) throw new Error(data.error || data.detail || "Cannot open source");
      if (request !== selection) return;
      files[path] = data.source;
    } catch (error) {
      const box = document.querySelector("#error"); box.textContent = error.message; box.hidden = false; return;
    }
  }
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
  try { localStorage.setItem('markata-workspace-active-v1:' + workspaceKey, path); } catch { warning(); }
  if (!(path in drafts)) {
    drafts[path] = files[path] || '';
    if (editable(path)) {
      try { drafts[path] = readDraft(path, drafts[path]); }
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
    try { writeDraft(active, drafts[active]);
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
  async attach(editor, storageWarning) {
    currentEditor = editor; warning = storageWarning;
    try {
      const response = await fetch('/workspace');
      const data = await response.json();
      if (!response.ok) throw new Error(data.error || data.detail || 'Cannot load workspace');
      metadata = data.files; workspaceKey = data.rootUri || data.root || location.origin;
      files = Object.fromEntries(metadata.map(file => [file.path, null]));
    } catch (error) {
      const box = document.querySelector('#error'); box.textContent = error.message; box.hidden = false;
      files = {};
    }
    editor.on('beforeChange', (instance, change) => {
      if (!editable(active) && change.origin !== 'setValue') change.cancel();
    });
    let selected = Object.keys(files)[0] || 'draft.md';
    try {
      selected = localStorage.getItem('markata-workspace-active-v1:' + workspaceKey) || selected;
      drafts['draft.md'] = readDraft('draft.md', editor.getValue());
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

      const tags = metadata.find(file => file.path === path)?.tags || [];
      for (const feed of tags) {
        if (feed) aliases[`${feed}/${path}`] = path;
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
    await open(selected in files || selected === 'draft.md' ? selected : (Object.keys(files)[0] || 'draft.md'));
  }
};
