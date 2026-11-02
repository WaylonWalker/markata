import {FileTree} from '@pierre/trees';

let files = SITE_FILES;
let workspaceKey = "example", metadata = [], metadataByPath = new Map(), selection = 0;
let scratchPath = 'draft.md';
const legacyStorageKey = path => 'markata-workspace-v1:' + workspaceKey + ':' + path;
const storageKey = path => legacyStorageKey('buffer:') + (path === scratchPath ? 'scratch' : 'file:' + path);
let active = 'draft.md', drafts = {}, currentEditor, warning, label;
let fileTree, feedTree, tagTree, syncing = false;
let siteAutosave = false;
const aliases = {}, feedAliases = {};
const baselines = {}, versions = {}, restored = {};
let workspaceCounts = {files:0,feeds:0,tags:0};
const saveTimers = {}, saveQueues = {}, saveRevisions = {};
function saveToSite(path, source) {
  const revision = saveRevisions[path] = (saveRevisions[path] || 0) + 1;
  clearTimeout(saveTimers[path]);
  saveTimers[path] = setTimeout(() => {
    saveQueues[path] = (saveQueues[path] || Promise.resolve()).catch(() => {}).then(async () => {
      if (revision !== saveRevisions[path]) return;
      const response = await fetch('/workspace/file?path=' + encodeURIComponent(path), {
        method:'PUT', headers:{'Content-Type':'text/plain;charset=utf-8'}, body:source
      });
      const data = await response.json();
      if (!response.ok) throw new Error(data.error || 'Cannot save to the site');
      baselines[path] = source;
      versions[path] = data.version;
      if (active === path && currentEditor.getValue() === source) {
        document.querySelector('#saved').textContent = 'Saved to site';
      }
    }).catch(error => {
      if (active === path) document.querySelector('#saved').textContent = `Site save failed · ${error.message} · recovery copy kept`;
    });
  },350);
}
function syncSelection(tree, path) {
  if (!tree || !path) return;
  for (const previous of tree.getSelectedPaths()) if (previous !== path) tree.getItem(previous)?.deselect();
  if (!tree.getSelectedPaths().includes(path)) tree.getItem(path)?.select();
}
const editable = path => /\.(md|markdown|mdown|mkd)$/i.test(path);
function readDraft(path, fallback) {
  delete restored[path];
  let cached = localStorage.getItem(storageKey(path));
  if (cached === null && (path === scratchPath || path !== 'draft.md')) cached = localStorage.getItem(legacyStorageKey(path === scratchPath ? 'draft.md' : path));
  if (cached === null) return fallback;
  try {
    const record = JSON.parse(cached);
    if (['markata-draft-v2','markata-draft-v3'].includes(record?.format) && typeof record.source === 'string') { restored[path] = record; return record.source; }
  } catch { /* Legacy drafts contain plain Markdown. */ }
  // Earlier builds could persist an empty placeholder before loading source.
  // Preserve all nonempty legacy drafts; new records also preserve intentional blanks.
  return cached === '' && path !== scratchPath ? fallback : cached;
}
function writeDraft(path, source) {
  localStorage.setItem(storageKey(path), JSON.stringify({format:'markata-draft-v3', source, dirty:path === scratchPath || source !== baselines[path], version:versions[path]}));
}
function remember() {
  if (!currentEditor || !editable(active)) return;
  drafts[active] = currentEditor.getValue();
  try { writeDraft(active, drafts[active]); }
  catch { warning(); }
}
async function open(path) {
  if (syncing) return;
  if (!Object.hasOwn(files,path) && path !== scratchPath) return;
  const request = ++selection;
  let disk = null;
  if (Object.hasOwn(files,path)) {
    try {
      const response = await fetch('/workspace/file?path=' + encodeURIComponent(path));
      const data = await response.json();
      if (!response.ok) throw new Error(data.error || data.detail || 'Cannot open source');
      if (request !== selection) return;
      disk = data;
    } catch (error) {
      const box = document.querySelector('#error'); box.textContent = error.message; box.hidden = false; return;
    }
  }
  remember();
  let diskChanged = false;
  if (disk) {
    let previous, dirty, previousVersion;
    if (Object.hasOwn(drafts,path)) {
      previous = drafts[path]; dirty = previous !== baselines[path]; previousVersion = versions[path];
    } else {
      try { previous = readDraft(path,disk.source); }
      catch { previous = disk.source; warning(); }
      dirty = restored[path]?.format === 'markata-draft-v3' ? restored[path].dirty !== false : previous !== disk.source;
      previousVersion = restored[path]?.version;
    }
    diskChanged = dirty && previousVersion && disk.version && previousVersion !== disk.version;
    drafts[path] = dirty ? previous : disk.source;
    files[path] = disk.source;
    baselines[path] = disk.source; versions[path] = disk.version;
  }
  syncing = true;
  try {
    syncSelection(fileTree, path);
    for (const [tree, lookup] of [[feedTree,feedAliases],[tagTree,aliases]]) {
      const selectedAlias = tree?.getSelectedPaths().find(alias => lookup[alias] === path);
      const alias = selectedAlias || Object.keys(lookup).find(alias => lookup[alias] === path);
      for (const previous of tree?.getSelectedPaths() || []) {
        if (previous !== alias) tree.getItem(previous)?.deselect();
      }
      syncSelection(tree, alias);
    }
  } finally { syncing = false; }
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
  if (diskChanged) document.querySelector('#saved').textContent = 'File changed on disk · keeping your browser draft';
}
window.MarkataSiteTree = {
  save() {
    if (!currentEditor) return false;
    if (!editable(active)) return true;
    drafts[active] = currentEditor.getValue();
    try { writeDraft(active, drafts[active]);
      document.querySelector('#saved').textContent = 'Draft saved on this device'; }
    catch { warning(); }
    if (siteAutosave && Object.hasOwn(baselines, active)) saveToSite(active, drafts[active]);
    return true;
  },
  filename: () => active === scratchPath ? 'draft.md' : active.split('/').pop(),
  markdown(text) {
    if (!editable(active)) return '````toml\n' + text + '\n````';
    return text.replace(/^---\r?\n[\s\S]*?\r?\n---(?:\r?\n|$)/, '');
  },
  async import(text) { await open(scratchPath); currentEditor.setValue(text); this.save(); },
  openPost(slug) {
    const post = metadata.find(file => file.slug?.toLowerCase() === slug.toLowerCase());
    if (post) open(post.path);
  },
  files: () => metadata.map(file => ({...file})),
  openFile(path) { if (Object.hasOwn(files, path)) open(path); },
  activeFile: () => active,
  workspaceId: () => workspaceKey,
  activate(path) { localStorage.setItem('markata-workspace-active-v1:' + workspaceKey,path); },
  async attach(editor, storageWarning, options = {}) {
    currentEditor = editor; warning = storageWarning;
    siteAutosave = options.saveToSite === true;
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
    const feedOption = document.querySelector('#site-view option[value="feeds"]');
    const tagOption = document.querySelector('#site-view option[value="tags"]');
    feedOption.disabled = true; tagOption.disabled = true;
    feedOption.textContent = 'Feeds (loading…)'; tagOption.textContent = 'Tags (loading…)';
    const buildTags = () => {
      if (feedTree) return;
      for (const file of metadata) {
        for (const tag of file.tags || []) aliases[`${tag}/${file.path}`] = file.path;
      }
      tagTree = new FileTree({paths:Object.keys(aliases),initialExpansion:'open',search:true,
        onSelectionChange(paths) {const path=aliases[paths[paths.length-1]];if(path)open(path);}});
      tagTree.render({containerWrapper:document.querySelector('#site-tags')});
      syncing = true;
      try { syncSelection(tagTree, Object.keys(aliases).find(alias => aliases[alias] === active)); }
      finally { syncing = false; }
    };
    const buildFeeds = feeds => {
      if (feedTree || !feeds?.length) return;
      let rootPath = '';
      try { rootPath = decodeURIComponent(new URL(workspaceKey).pathname).replace(/\\/g,'/').replace(/\/+$/,''); }
      catch { /* Display paths without an absolute host root when the root is not a URI. */ }
      for (const feed of feeds || []) for (const path of feed.paths || []) {
        let displayPath = String(path).replace(/\\/g,'/');
        if (rootPath && displayPath.startsWith(`${rootPath}/`)) displayPath = displayPath.slice(rootPath.length + 1);
        displayPath = displayPath.replace(/^\/+/, '');
        const file = metadataByPath.get(displayPath);
        let leaf = String(file?.slug || displayPath.replace(/\//g,' — ')).replace(/\//g,' — ');
        const root = `${feed.title} (${feed.paths.length.toLocaleString()})`;
        const base = leaf;
        for (let suffix=2; Object.hasOwn(feedAliases,`${root}/${leaf}`); suffix++) leaf = `${base} (${suffix})`;
        feedAliases[`${root}/${leaf}`] = displayPath;
      }
      feedTree = new FileTree({paths:Object.keys(feedAliases),initialExpansion:0,search:true,
        onSelectionChange(paths) {const path=feedAliases[paths[paths.length-1]];if(path)open(path);}});
      feedTree.render({containerWrapper:document.querySelector('#site-feeds')});
      syncing = true;
      try { syncSelection(feedTree, Object.keys(feedAliases).find(alias => feedAliases[alias] === active)); }
      finally { syncing = false; }
    };
    const updateWorkspaceContext = view => {
      const context = document.querySelector('#site-context');
      const count = workspaceCounts[view] || 0;
      context.textContent = view === 'files' ? `${count} source files in this site.`
        : view === 'feeds' ? `${count} configured feeds · select a post to open it.`
        : `${count} tags from loaded frontmatter.`;
    };
    document.querySelector('#site-view').addEventListener('change', event => {
      if (event.target.value === 'feeds') buildFeeds(window.MarkataSiteTree.feeds || []);
      if (event.target.value === 'tags') buildTags();
      document.querySelector('#site-files').hidden = event.target.value !== 'files';
      document.querySelector('#site-feeds').hidden = event.target.value !== 'feeds';
      document.querySelector('#site-tags').hidden = event.target.value !== 'tags';
      updateWorkspaceContext(event.target.value);
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
        const tags = new Set(metadata.flatMap(file => file.tags || []));
        workspaceCounts.files = metadata.length;
        workspaceCounts.tags = tags.size;
        document.querySelector('#site-view option[value="files"]').textContent = `Files (${metadata.length})`;
        tagOption.textContent = `Tags (${tags.size})`; tagOption.disabled = false;
        const feeds = data.feeds || [];
        workspaceCounts.feeds = feeds.length;
        window.MarkataSiteTree.feeds = feeds;
      feedOption.textContent = `Feeds (${feeds.length})`; feedOption.disabled = feeds.length === 0;
        if (document.querySelector('#site-view').value === 'feeds') buildFeeds(feeds);
        updateWorkspaceContext(document.querySelector('#site-view').value);
      } catch { tagOption.textContent = 'Tags unavailable'; feedOption.textContent = 'Feeds unavailable'; }
    },0);
  }
};
