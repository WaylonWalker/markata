/* Fuzzy workspace file picker, bound to Ctrl/Cmd+P. */
(() => {
  'use strict';

  window.MarkataFileFinder = {attach};

  function attach(editor) {
    const workspace = window.MarkataSiteTree;
    if (!workspace?.files || !workspace?.openFile) return;

    const dialog = document.createElement('dialog');
    dialog.className = 'file-finder';
    dialog.setAttribute('aria-labelledby', 'file-finder-title');
    dialog.innerHTML = '<div class="file-finder-heading"><h2 id="file-finder-title">Go to File</h2><button type="button" class="file-finder-close" aria-label="Close file finder">Close</button></div><label for="file-finder-query">Search site files</label><input id="file-finder-query" type="search" placeholder="Type a path, title, slug, or tag" autocomplete="off" role="combobox" aria-autocomplete="list" aria-expanded="true" aria-controls="file-finder-results"><p id="file-finder-count" role="status"></p><div id="file-finder-results" role="listbox" aria-label="Workspace files"></div><p class="file-finder-help">↑ ↓ Choose · Enter Open · Esc Close</p>';
    document.body.append(dialog);

    const query = dialog.querySelector('#file-finder-query');
    const results = dialog.querySelector('#file-finder-results');
    const count = dialog.querySelector('#file-finder-count');
    const storageKey = `markata-file-finder-recent-v1:${workspace.workspaceId?.() || location.origin}`;
    let files = workspace.files();
    let matches = [];
    let active = 0;

    function recentPaths() {
      try { return JSON.parse(localStorage.getItem(storageKey) || '[]'); }
      catch { return []; }
    }

    function score(file, terms) {
      const fields = [file.path, file.title, file.slug, ...(file.tags || [])]
        .filter(Boolean).map(value => String(value).toLowerCase());
      let total = 0;
      for (const term of terms) {
        let best = Infinity;
        for (let fieldIndex = 0; fieldIndex < fields.length; fieldIndex++) {
          const field = fields[fieldIndex];
          const found = field.indexOf(term);
          if (found >= 0) {
            best = Math.min(best, fieldIndex * 18 + found + (found === 0 ? -8 : 0));
            continue;
          }
          let cursor = 0, gaps = 0, first = -1;
          for (const char of term) {
            const position = field.indexOf(char, cursor);
            if (position < 0) { first = -1; break; }
            if (first < 0) first = position;
            gaps += position - cursor;
            cursor = position + 1;
          }
          if (first >= 0) best = Math.min(best, fieldIndex * 18 + first + gaps * 3 + 12);
        }
        if (!Number.isFinite(best)) return Infinity;
        total += best;
      }
      return total;
    }

    function select(index) {
      active = index;
      [...results.children].forEach((row, i) => row.setAttribute('aria-selected', String(i === active)));
      const row = results.children[active];
      if (row) {
        query.setAttribute('aria-activedescendant', row.id);
        row.scrollIntoView({block:'nearest'});
      } else query.removeAttribute('aria-activedescendant');
    }

    function remember(path) {
      try {
        localStorage.setItem(storageKey, JSON.stringify([path, ...recentPaths().filter(item => item !== path)].slice(0, 12)));
      } catch { /* Opening files must still work when storage is unavailable. */ }
    }

    function openFile(file) {
      if (!file) return;
      remember(file.path);
      workspace.openFile(file.path);
      dialog.close();
    }

    function refresh() {
      const terms = query.value.toLowerCase().trim().split(/\s+/).filter(Boolean);
      if (!terms.length) {
        const byPath = new Map(files.map(file => [file.path, file]));
        const recent = [workspace.activeFile?.(), ...recentPaths()];
        matches = [...new Set(recent)].map(path => byPath.get(path)).filter(Boolean);
        count.textContent = matches.length ? 'Current and recently opened' : `Search ${files.length.toLocaleString()} site files`;
      } else {
        matches = files.map(file => ({file, score:score(file, terms)}))
          .filter(item => Number.isFinite(item.score))
          .sort((a, b) => a.score - b.score || a.file.path.localeCompare(b.file.path))
          .slice(0, 50).map(item => item.file);
        count.textContent = matches.length ? `${matches.length} matching files${matches.length === 50 ? ' · showing top 50' : ''}` : 'No matching files';
      }
      results.replaceChildren();
      matches.forEach((file, index) => {
        const row = document.createElement('button');
        row.type = 'button'; row.id = `file-finder-result-${index}`; row.setAttribute('role', 'option');
        const title = document.createElement('strong'); title.textContent = file.title || file.path;
        const path = document.createElement('span'); path.textContent = file.path;
        row.append(title, path);
        row.addEventListener('pointermove', () => select(index));
        row.addEventListener('click', () => openFile(file));
        results.append(row);
      });
      select(0);
    }

    function open() {
      files = workspace.files();
      query.value = '';
      dialog.showModal();
      refresh();
      query.focus();
    }

    query.addEventListener('input', refresh);
    dialog.querySelector('.file-finder-close').addEventListener('click', () => dialog.close());
    dialog.addEventListener('close', () => editor.focus());
    dialog.addEventListener('keydown', event => {
      if (event.isComposing) return;
      if (event.key === 'Escape') { event.preventDefault(); event.stopPropagation(); dialog.close(); return; }
      if (event.target !== query || !['ArrowDown', 'ArrowUp', 'Enter'].includes(event.key)) return;
      event.preventDefault();
      if (event.key === 'Enter') openFile(matches[active]);
      else if (matches.length) select((active + (event.key === 'ArrowDown' ? 1 : -1) + matches.length) % matches.length);
    });
    dialog.addEventListener('click', event => {
      if (event.target !== dialog) return;
      const bounds = dialog.getBoundingClientRect();
      if (event.clientX < bounds.left || event.clientX > bounds.right || event.clientY < bounds.top || event.clientY > bounds.bottom) dialog.close();
    });
    document.addEventListener('keydown', event => {
      if (event.isComposing || event.altKey || event.shiftKey || !(event.ctrlKey || event.metaKey) || event.key.toLowerCase() !== 'p') return;
      event.preventDefault(); event.stopPropagation();
      if (!dialog.open) open();
    }, true);
  }
})();
