import {FileTree} from '@pierre/trees';

let files = SITE_FILES;
let workspaceKey = "example", metadata = [], metadataByPath = new Map(), selection = 0;
let scratchPath = 'draft.md';
const legacyStorageKey = path => 'markata-workspace-v1:' + workspaceKey + ':' + path;
const storageKey = path => legacyStorageKey('buffer:') + (path === scratchPath ? 'scratch' : 'file:' + path);
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
  let cached = localStorage.getItem(storageKey(path));
  if (cached === null && (path === scratchPath || path !== 'draft.md')) cached = localStorage.getItem(legacyStorageKey(path === scratchPath ? 'draft.md' : path));
  if (cached === null) return fallback;
  try {
    const record = JSON.parse(cached);
    if (record?.format === 'markata-draft-v2' && typeof record.source === 'string') return record.source;
  } catch { /* Legacy drafts contain plain Markdown. */ }
  // Earlier builds could persist an empty placeholder before loading source.
  // Preserve all nonempty legacy drafts; new records also preserve intentional blanks.
  return cached === '' && path !== scratchPath ? fallback : cached;
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
  if (Object.hasOwn(files, path) && files[path] === null && !(path in drafts)) {
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
  if (syncing || (!(Object.hasOwn(files, path)) && path !== scratchPath)) return;
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
  label.title = metadataByPath.get(path)?.title || path;
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
  filename: () => active === scratchPath ? 'draft.md' : active.split('/').pop(),
  markdown(text) {
    if (!editable(active)) return '````toml\n' + text + '\n````';
    return text.replace(/^---\r?\n[\s\S]*?\r?\n---(?:\r?\n|$)/, '');
  },
  import(text) { open(scratchPath); currentEditor.setValue(text); },
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
    scratchPath = 'draft.md';
    for (let suffix = 1; Object.hasOwn(files, scratchPath); suffix++) scratchPath = suffix === 1 ? 'Unsaved draft.md' : `Unsaved draft ${suffix}.md`;
    active = scratchPath;
    let selected = Object.keys(files)[0] || scratchPath;
    try {
      selected = localStorage.getItem('markata-workspace-active-v1:' + workspaceKey) || selected;
      drafts[scratchPath] = readDraft(scratchPath, editor.getValue());
    } catch { drafts[scratchPath] = editor.getValue(); warning(); }
    editor.setValue(drafts[scratchPath]);
    label = document.querySelector('#active-file');
    label.textContent = scratchPath;
    fileTree = new FileTree({paths:[scratchPath, ...Object.keys(files)], initialExpansion:'open',
      search:true,
      onSelectionChange(paths) { if (paths.length) open(paths[paths.length - 1]); }});
    fileTree.render({containerWrapper:document.querySelector('#site-files')});
    const tagOption = document.querySelector('#site-view option[value="feeds"]');
    tagOption.disabled = true; tagOption.textContent = 'Tags (loading…)';
    const buildTags = () => {
      if (feedTree) return;
      for (const file of metadata) {
        for (const tag of file.tags || []) aliases[`${tag}/${file.path}`] = file.path;
      }
      feedTree = new FileTree({paths:Object.keys(aliases),initialExpansion:'open',search:true,
        onSelectionChange(paths) {const path=aliases[paths[paths.length-1]];if(path)open(path);}});
      feedTree.render({containerWrapper:document.querySelector('#site-feeds')});
      syncing = true;
      try { syncSelection(feedTree, Object.keys(aliases).find(alias => aliases[alias] === active)); }
      finally { syncing = false; }
    };
    document.querySelector('#site-view').addEventListener('change', event => {
      if (event.target.value === 'feeds') buildTags();
      document.querySelector('#site-files').hidden = event.target.value !== 'files';
      document.querySelector('#site-feeds').hidden = event.target.value !== 'feeds';
    });
    // Start with a real site post. The previous standalone draft stays in the tree.
    await open(Object.hasOwn(files, selected) || selected === scratchPath ? selected : (Object.keys(files)[0] || scratchPath));
    // Yield to the initial preview; metadata is never part of attach's promise.
    setTimeout(async () => {
      try {
        const response = await fetch('/workspace?metadata=true');
        const data = await response.json();
        if (!response.ok) throw new Error(data.error || data.detail || 'Metadata unavailable');
        metadata = data.files;
        metadataByPath = new Map(metadata.map(file => [file.path,file]));
        label.title = metadataByPath.get(active)?.title || active;
        tagOption.textContent = 'Tags'; tagOption.disabled = false;
      } catch { tagOption.textContent = 'Tags unavailable'; }
    },0);
  }
};
