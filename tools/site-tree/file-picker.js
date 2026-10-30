import {createFileSearch} from './file-search.js';
export function attachFilePicker(editor, openFile) {
  const dialog = document.createElement('dialog');
  dialog.className = 'markdown-commands';
  dialog.setAttribute('aria-labelledby', 'files-title');
  dialog.innerHTML = '<div class="commands-heading"><h2 id="files-title">Open file</h2><button type="button" aria-label="Close file picker">Close</button></div><label for="files-query">Find a workspace file</label><input id="files-query" type="search" placeholder="Filename, path, or title" autocomplete="off" role="combobox" aria-autocomplete="list" aria-expanded="true" aria-controls="files-results"><p role="status"></p><div id="files-results" role="listbox" aria-label="Workspace files"></div><p class="commands-help">↑ ↓ Choose · Enter Open · Esc Close</p>';
  document.body.append(dialog);
  const query = dialog.querySelector('input'), list = dialog.querySelector('[role=listbox]');
  let search = createFileSearch([]), matches = [], active = 0;
  function select(index) {
    active = index;
    [...list.children].forEach((row,i) => row.setAttribute('aria-selected', String(i === index)));
    const row = list.children[index];
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
  function refresh() {
    matches = search(query.value); list.replaceChildren();
    dialog.querySelector('[role=status]').textContent = matches.length ? `${matches.length} of ${matches.total} matching files` : 'No matching files. Try another name or path.';
    matches.forEach((match,i) => {
      const row = document.createElement('div'); row.id = `file-result-${i}`; row.setAttribute('role','option');
      const name = document.createElement('strong'), path = document.createElement('code'), title = document.createElement('span');
      highlighted(name, match[0], match.obj.name); highlighted(path, match[1], match.obj.path);
      highlighted(title, match[2], match.obj.title);
      row.append(name,path,title);
      row.addEventListener('pointermove', () => select(i));
      row.addEventListener('click', () => choose(i)); list.append(row);
    }); select(0);
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
  return {open:show, update(files) { search = createFileSearch(files); if (dialog.open) refresh(); }};
}
