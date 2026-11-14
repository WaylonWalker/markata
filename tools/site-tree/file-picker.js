import {createFileSearch, parseFileQuery, fileFacets, suggestFor, acceptSuggestion} from './file-search.js';
export function attachFilePicker(editor, openFile) {
  const dialog = document.createElement('dialog');
  dialog.className = 'markdown-commands file-picker';
  dialog.setAttribute('aria-labelledby', 'files-title');
  dialog.innerHTML = '<div class="commands-heading"><h2 id="files-title">Open file</h2><button type="button" class="commands-close" aria-label="Close file picker">Close</button></div><label for="files-query">Find a workspace file</label><input id="files-query" type="search" placeholder="Find a file" autocomplete="off" role="combobox" aria-autocomplete="list" aria-expanded="true" aria-controls="files-results"><p role="status"></p><div id="files-results" role="listbox" aria-label="Workspace files"></div><p class="commands-help">↑ ↓ Choose · Enter Open · Esc Close</p>';
  document.body.append(dialog);
  const sidebarHeading = document.querySelector('#site-sidebar .pane-label');
  if (sidebarHeading) {
    const button = document.createElement('button');
    button.type = 'button'; button.className = 'files-open'; button.textContent = 'Open file';
    button.title = 'Open file (Ctrl/Cmd+P)'; button.addEventListener('click', show);
    sidebarHeading.after(button);
  }
  const query = dialog.querySelector('input'), list = dialog.querySelector('[role=listbox]');
  const configuredFeeds = () => window.MarkataSiteTree?.feeds || [];
  let search = createFileSearch([], configuredFeeds), facets = fileFacets([]), matches = [], active = 0;
  const hint = document.createElement('p'); hint.className = 'files-hint';
  hint.innerHTML = 'Name, path, or title · <kbd>#</kbd>tag · <kbd>status:</kbd>draft or published · <kbd>in:</kbd>folder · <kbd>Tab</kbd> completes';
  const suggest = document.createElement('div'); suggest.className = 'files-suggest'; suggest.hidden = true;
  query.after(hint); hint.after(suggest);
  function applySuggestion(value) { query.value = acceptSuggestion(query.value, value); query.focus(); refresh(); }
  function renderSuggestions() {
    const feedSlugs = configuredFeeds().filter(feed => feed.paths?.length).map(feed => feed.slug);
    const {items} = suggestFor(query.value, {...facets, feeds: feedSlugs});
    suggest.replaceChildren(...items.map(item => {
      const chip = document.createElement('button'); chip.type = 'button'; chip.textContent = item.label;
      chip.addEventListener('click', () => applySuggestion(item.value));
      return chip;
    }));
    suggest.hidden = !items.length;
    return items;
  }
  const options = () => [...list.querySelectorAll('[role=option]')];
  function select(index) {
    active = index;
    options().forEach((row,i) => row.setAttribute('aria-selected', String(i === index)));
    const row = options()[index];
    if (row) { query.setAttribute('aria-activedescendant', row.id); row.scrollIntoView({block:'nearest'}); }
    else query.removeAttribute('aria-activedescendant');
  }
  function highlighted(node, result, fallback) {
    // Construct text nodes; filenames and metadata are never parsed as HTML.
    const indexes = new Set(result?.indexes || []);
    Array.from(result?.target || fallback).forEach((char,i) => {
      if (indexes.has(i)) { const mark = document.createElement('mark'); mark.textContent = char; node.append(mark); }
      else node.append(document.createTextNode(char));
    });
  }
  function recentPaths() {
    try {
      const key = 'markata-recent-files-v1:' + (window.MarkataSiteTree?.workspaceId?.() || location.origin);
      const saved = JSON.parse(localStorage.getItem(key) || '[]');
      return Array.isArray(saved) ? saved.filter(path => typeof path === 'string') : [];
    } catch { return []; }
  }
  function metaRow(match) {
    const file = match.obj;
    const meta = document.createElement('div'); meta.className = 'file-result-meta';
    if (file.status) {
      const badge = document.createElement('span'); badge.className = 'file-status'; badge.dataset.status = file.status;
      badge.textContent = file.status; meta.append(badge);
    }
    if (file.date) { const date = document.createElement('span'); date.className = 'file-date'; date.textContent = file.date; meta.append(date); }
    for (const tag of (file.tags || []).slice(0, 2)) {
      const chip = document.createElement('span'); chip.className = 'file-tag'; chip.textContent = '#' + tag; meta.append(chip);
    }
    return meta.childElementCount ? meta : null;
  }
  function row(match, index) {
    const node = document.createElement('div'); node.id = `file-result-${index}`; node.setAttribute('role','option');
    const name = document.createElement('strong'), path = document.createElement('code'), title = document.createElement('span');
    name.className = 'file-result-name'; path.className = 'file-result-path'; title.className = 'file-result-title';
    highlighted(name, match[0], match.obj.name); highlighted(path, match[1], match.obj.path);
    highlighted(title, match[2], match.obj.title);
    const image = match.obj.image;
    if (typeof image === 'string' && /^\/images\/[^?#]+\.(?:jpe?g|png|webp|gif|avif)$/i.test(image)) {
      const thumb = document.createElement('img');
      thumb.className = 'file-result-image'; thumb.alt = ''; thumb.loading = 'lazy';
      thumb.src = '/workspace/image?path=' + encodeURIComponent(image.slice(1));
      node.classList.add('has-image'); node.append(thumb);
    }
    const details = document.createElement('div'); details.className = 'file-result-details';
    details.append(name,path);
    if (match.obj.title && match.obj.title !== match.obj.name) details.append(title);
    const meta = metaRow(match); if (meta) details.append(meta);
    node.append(details);
    node.addEventListener('pointermove', () => select(index));
    node.addEventListener('click', () => choose(index));
    return node;
  }
  function heading(text, count) {
    const node = document.createElement('div'); node.className = 'file-group-heading'; node.setAttribute('role','presentation');
    node.textContent = `${text} · ${count}`; return node;
  }
  function refresh() {
    const found = search(query.value);
    const browsing = !parseFileQuery(query.value).text && !parseFileQuery(query.value).filters.length;
    let groups = [];
    if (browsing) {
      const seen = new Set();
      const take = list => list.filter(match => !seen.has(match.obj.path) && seen.add(match.obj.path));
      const recent = recentPaths().map(path => found.find(match => match.obj.path === path)).filter(Boolean).slice(0,5);
      const recentTaken = take(recent);
      const drafts = take(found.filter(match => match.obj.status === 'draft'));
      const rest = take(found);
      groups = [['Recent', recentTaken], ['Drafts', drafts], ['All files', rest]].filter(([, items]) => items.length);
    } else {
      groups = [[null, found]];
    }
    matches = groups.flatMap(([, items]) => items);
    list.replaceChildren();
    dialog.querySelector('[role=status]').textContent = matches.length ? `${matches.length} of ${found.total} matching files` : 'No matching files. Try another name, path, #tag, or status:draft.';
    let index = 0;
    for (const [name, items] of groups) {
      if (name) list.append(heading(name, items.length));
      for (const match of items) list.append(row(match, index++));
    }
    renderSuggestions();
    select(0);
  }
  async function choose(index) {
    const match = matches[index]; if (!match) return;
    dialog.close(); await openFile(match.obj.path); editor.focus();
  }
  function show() { if (dialog.open) { query.focus(); return; } query.value = ''; dialog.showModal(); refresh(); query.focus(); }
  dialog.querySelector('button').addEventListener('click', () => dialog.close());
  dialog.addEventListener('close', () => editor.focus());
  query.addEventListener('input', refresh);
  dialog.addEventListener('keydown', event => {
    if (event.isComposing) return;
    if (event.key === 'Escape') { event.preventDefault(); event.stopPropagation(); dialog.close(); return; }
    if (event.key === 'Tab' && !event.shiftKey) {
      const [first] = suggestFor(query.value, {...facets, feeds: configuredFeeds().filter(feed => feed.paths?.length).map(feed => feed.slug)}).items;
      if (first) { event.preventDefault(); applySuggestion(first.value); }
    }
    if (event.key === 'Enter') { event.preventDefault(); choose(active); }
    if (['ArrowDown','ArrowUp'].includes(event.key)) {
      event.preventDefault(); if (matches.length) select((active + (event.key === 'ArrowDown' ? 1 : -1) + matches.length) % matches.length);
    }
  });
  document.addEventListener('keydown', event => {
    if (!event.isComposing && (event.ctrlKey || event.metaKey) && !event.shiftKey && !event.altKey && event.key.toLowerCase() === 'p') {
      event.preventDefault(); event.stopPropagation(); show();
    }
  }, true);
  return {open:show, update(files) { search = createFileSearch(files, configuredFeeds); facets = fileFacets(files); if (dialog.open) refresh(); }};
}
